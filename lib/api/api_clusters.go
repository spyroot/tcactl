// Package api
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Mustafa mbayramo@vmware.com
package api

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	"strings"
)

func ClusterFields() []string {
	f := response.ClusterSpec{}
	fields, _ := f.GetFields()

	var keys []string
	for s, _ := range fields {
		keys = append(keys, s)
	}

	return keys
}

// GetCluster -  method retrieve cluster information
func (a *TcaApi) GetCluster(clusterId string) (*response.ClusterSpec, error) {

	if IsValidUUID(clusterId) {
		return a.rest.GetCluster(clusterId)
	}

	clusters, err := a.rest.GetClusters()
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	_clusterId, err := clusters.GetClusterId(clusterId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return a.rest.GetCluster(_clusterId)
}

// GetClusterNodePool -  API Method lookup clusters node pool.
// clusterId is identified either a name or UUID.
// nodePoolId is identifier either a name or UUID.
func (a *TcaApi) GetClusterNodePool(clusterId string, nodePoolId string) (*response.NodesSpecs, error) {

	var (
		_clusterid  = clusterId
		_nodePoolId = nodePoolId
		err         error
	)

	if !IsValidUUID(_clusterid) {
		glog.Infof("Resolving cluster name %s to id", clusterId)
		_clusterid, err = a.ResolveClusterName(clusterId)
		if err != nil {
			return nil, err
		}
	}

	if !IsValidUUID(_nodePoolId) {
		glog.Infof("Resolving pool name %s to id", _nodePoolId)
		_nodePoolId, err = a.ResolvePoolId(_nodePoolId, _clusterid)
		if err != nil {
			return nil, err
		}
	}

	return a.rest.GetClusterNodePool(_clusterid, _nodePoolId)
}

// GetClusterTask method return list task models.ClusterTask
// currently executing on given cluster
func (a *TcaApi) GetClusterTask(clusterId string, showChildren bool) (*models.ClusterTask, error) {

	var err error
	_clustered := clusterId

	if !IsValidUUID(_clustered) {
		glog.Infof("Resolving cluster id from name %s", clusterId)
		_clustered, err = a.ResolveClusterName(_clustered)
		if err != nil {
			return nil, err
		}
	}

	clusters, err := a.rest.GetClusters()
	if err != nil {
		return nil, nil
	}

	clusterSpec, err := clusters.GetClusterSpec(clusterId)
	if err != nil {
		return nil, err
	}

	r := request.NewClusterTaskQuery(clusterSpec.ManagementClusterId)
	r.IncludeChildTasks = showChildren

	return a.rest.GetClustersTask(r)
}

// CreateClusters - create new cluster
// in Dry Run we only parse
func (a *TcaApi) CreateClusters(spec *request.Cluster,
	isDry bool, isBlocking bool, verbose bool) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if spec == nil {
		return nil, fmt.Errorf("new cluster initialSpec can't be nil")
	}

	spec.ClusterType = strings.ToUpper(spec.ClusterType)

	// validate cluster
	if err := validateClusterSpec(spec); err != nil {
		return nil, err
	}

	// do all sanity check here.
	tenants, err := a.rest.GetClusters()
	if err != nil {
		return nil, err
	}

	_, err = tenants.GetClusterId(spec.Name)
	// swap name
	if err == nil {
		spec.Name = spec.Name + "-" + uuid.New().String()
		spec.Name = spec.Name[0:25]
		glog.Infof("Duplicate name regenerated new name '%v'", spec.Name)
	}

	// resolve template id, and cluster type
	spec.ClusterTemplateId, err = a.ResolveTypedTemplateId(spec.ClusterTemplateId, spec.ClusterType)
	if err != nil {
		return nil, err
	}

	// get template and validate specs
	t, err := a.rest.GetClusterTemplate(spec.ClusterTemplateId)
	if err != nil {
		return nil, err
	}

	glog.Infof("Validating node pool specs.")
	_, err = t.ValidateSpec(spec)
	if err != nil {
		return nil, err
	}

	glog.Infof("Resolved template id %v", spec.ClusterTemplateId)
	tenant, err := a.validateCloudEndpoint(spec.HcxCloudUrl)
	if err != nil {
		return nil, err
	}

	err = a.validateTenant(tenant)
	if err != nil {
		return nil, err
	}

	if spec.ClusterType == string(request.ClusterWorkload) {
		// resolve template id, in case client used name instead id
		mgmtClusterId, err := tenants.GetClusterId(spec.ManagementClusterId)
		if err != nil {
			return nil, err
		}
		mgmtSpec, err := tenants.GetClusterSpec(mgmtClusterId)
		if err != nil {
			return nil, err
		}
		if mgmtSpec.ClusterType != string(request.ClusterManagement) {
			return nil, fmt.Errorf("managementClusterId cluster id is not management cluster")
		}
		glog.Infof("Resolved cluster name %v", mgmtClusterId)
		spec.ManagementClusterId = mgmtClusterId

		err = a.validatePlacements(spec, tenant)
		if err != nil {
			return nil, err
		}
	}

	if spec.ClusterType == string(request.ClusterManagement) {
		// ignoring mgmt cluster id
		spec.ManagementClusterId = ""
		err = a.validatePlacements(spec, tenant)
		if err != nil {
			return nil, err
		}
	}

	spec.ClusterPassword = b64.StdEncoding.EncodeToString([]byte(spec.ClusterPassword))

	if isDry {
		return nil, nil
	}

	task, err := a.rest.CreateCluster(spec)
	if err != nil {
		return nil, err
	}

	if isBlocking {
		err := a.BlockWaitTaskFinish(context.Background(), task, TaskStateSuccess, BlockMaxRetryTimer, verbose)
		if err != nil {
			return task, err
		}
	}

	return task, err
}

// DeleteCluster method deletes cluster, note you can't
// delete cluster that is active and contains running instances,
// TCA will return error.
// TODO refactor return task
func (a *TcaApi) DeleteCluster(clusterId string) (bool, error) {

	if a.rest == nil {
		return false, fmt.Errorf("rest interface is nil")
	}

	clusters, err := a.rest.GetClusters()
	if err != nil {
		return false, err
	}

	cid, clusterErr := clusters.GetClusterId(clusterId)
	if clusterErr != nil {
		glog.Error(clusterErr)
		return false, err
	}

	ok, err := a.rest.DeleteCluster(cid)
	if err != nil {
		return false, err
	}

	return ok, nil
}
