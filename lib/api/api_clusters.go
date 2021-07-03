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
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/api_errors"
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

// ResolveManagementCluster mgmt cluster
func (a *TcaApi) ResolveManagementCluster(NameOrId string, tenants *response.Clusters) (*response.ClusterSpec, error) {

	if tenants == nil {
		return nil, errors.New("tenant information is nil")
	}

	mcid, err := tenants.GetClusterId(NameOrId)
	if err != nil {
		return nil, err
	}

	spec, err := tenants.GetClusterSpec(mcid)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(spec.ClusterType) != strings.ToLower(string(request.ClusterManagement)) {
		glog.Warning("managementClusterId cluster id is not management cluster")
		return nil, errors.New("managementClusterId cluster id is not management cluster")
	}

	return spec, nil
}

func specFixup(spec *request.Cluster) {

	// fix cluster type
	spec.ClusterType = strings.ToUpper(spec.ClusterType)
	if spec.ClusterConfig != nil {
		// fix url for datastore
		for _, config := range spec.ClusterConfig.Csi {
			normalizeDatastoreName(config.Properties.DatastoreUrl)
		}
	}
}

// CreateClusters - creates new cluster , in Dry Run method only will do
// spec validation. isBlocking will indicate if caller expects cluster
// task to finish. verbose will output status of each task,
// linked to specific cluster creation.
func (a *TcaApi) CreateClusters(spec *request.Cluster,
	isDry bool, isBlocking bool, verbose bool) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if spec == nil {
		return nil, api_errors.NewInvalidSpec("new cluster spec can't be nil")
	}

	// fix spec
	specFixup(spec)

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
	// swap name if needed
	if err == nil {
		spec.Name = spec.Name + "-" + uuid.New().String()
		spec.Name = spec.Name[0:25]
		glog.Infof("Duplicate name regenerated new name '%v'", spec.Name)
	}

	// resolve template id, and cluster type
	spec.ClusterTemplateId, err = a.NormalizeTemplateId(spec.ClusterTemplateId, spec.ClusterType)
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
	spec.HcxCloudUrl = tenant.HcxCloudURL

	err = a.validateTenant(tenant)
	if err != nil {
		return nil, err
	}

	if spec.IsWorkload() {
		glog.Infof("Validating workload cluster spec")

		// resolve template id, in case client used name instead id
		mspec, err := a.ResolveManagementCluster(spec.ManagementClusterId, tenants)
		if err != nil {
			return nil, err
		}

		spec.ManagementClusterId = mspec.Id
		err = a.validatePlacements(spec, tenant)
		if err != nil {
			return nil, err
		}
	}

	if spec.IsManagement() {
		glog.Infof("Validating management cluster spec")
		// ignoring mgmt cluster id
		spec.ManagementClusterId = ""
		err = a.validatePlacements(spec, tenant)
		if err != nil {
			return nil, err
		}
	}

	spec.ClusterPassword = b64.StdEncoding.EncodeToString([]byte(spec.ClusterPassword))

	glog.Infof("Cluster spec validated.")

	if isDry {
		return &models.TcaTask{}, nil
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

// DeleteCluster - Method deletes cluster, note in order delete
// it must not have anything running on it. In order delete management
// cluster, all tenant cluster must be deleted first.
func (a *TcaApi) DeleteCluster(clusterId string, isBlocking bool, verbose bool) (*models.TcaTask, error) {

	var (
		cid = clusterId
		err error
	)

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(cid) == 0 {
		return nil, fmt.Errorf("empty cluster id or name")
	}

	// resolve id if it not UUID
	if !IsValidUUID(cid) {
		cid, err = a.ResolveClusterName(cid)
		if err != nil {
			return nil, err
		}
	}

	task, err := a.rest.DeleteCluster(cid)
	if err != nil {
		return nil, err
	}

	// block and wait task to finish
	if isBlocking {
		err := a.BlockWaitTaskFinish(context.Background(), task, TaskStateSuccess, BlockMaxRetryTimer, verbose)
		if err != nil {
			return task, err
		}
	}

	return task, nil
}
