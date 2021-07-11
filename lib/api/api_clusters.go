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
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/netutils"
	"net"
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
func (a *TcaApi) GetCluster(ctx context.Context, clusterId string) (*response.ClusterSpec, error) {

	if IsValidUUID(clusterId) {
		return a.rest.GetCluster(ctx, clusterId)
	}

	clusters, err := a.rest.GetClusters(ctx)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	_clusterId, err := clusters.GetClusterId(clusterId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return a.rest.GetCluster(ctx, _clusterId)
}

// GetClusterNodePool -  API Method lookup clusters node pool.
// clusterId is identified either a name or UUID.
// nodePoolId is identifier either a name or UUID.
func (a *TcaApi) GetClusterNodePool(ctx context.Context, clusterId string, nodePoolId string) (*response.NodesSpecs, error) {

	var (
		_clusterid  = clusterId
		_nodePoolId = nodePoolId
		err         error
	)

	if !IsValidUUID(_clusterid) {
		glog.Infof("Resolving cluster name %s to id", clusterId)
		_clusterid, err = a.ResolveClusterName(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}

	if !IsValidUUID(_nodePoolId) {
		glog.Infof("Resolving pool name %s to id", _nodePoolId)
		_nodePoolId, err = a.ResolvePoolId(ctx, _nodePoolId, _clusterid)
		if err != nil {
			return nil, err
		}
	}

	return a.rest.GetClusterNodePool(_clusterid, _nodePoolId)
}

// GetClusterTask method return list task models.ClusterTask
// currently executing on given cluster
func (a *TcaApi) GetClusterTask(ctx context.Context, clusterId string, showChildren bool) (*models.ClusterTask, error) {

	var err error
	_clustered := clusterId

	if !IsValidUUID(_clustered) {
		glog.Infof("Resolving cluster id from name %s", clusterId)
		_clustered, err = a.ResolveClusterName(ctx, _clustered)
		if err != nil {
			return nil, err
		}
	}

	clusters, err := a.rest.GetClusters(ctx)
	if err != nil {
		return nil, nil
	}

	clusterSpec, err := clusters.GetClusterSpec(clusterId)
	if err != nil {
		return nil, err
	}

	r := specs.NewClusterTaskQuery(clusterSpec.ManagementClusterId)
	r.IncludeChildTasks = showChildren

	return a.rest.GetClustersTask(ctx, r)
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

	if strings.ToLower(spec.ClusterType) != strings.ToLower(string(specs.ClusterManagement)) {
		glog.Warning("managementClusterId cluster id is not management cluster")
		return nil, errors.New("managementClusterId cluster id is not management cluster")
	}

	return spec, nil
}

func normalizeSpec(spec *specs.SpecCluster) {

	// fix cluster type
	spec.ClusterType = strings.ToUpper(spec.ClusterType)
	if spec.ClusterConfig != nil {
		// fix url for datastore
		for _, config := range spec.ClusterConfig.Csi {
			normalizeDatastoreName(config.Properties.DatastoreUrl)
		}
	}

}

// checkClusterAddrConflict - check for cluster IP conflicts
func (a *TcaApi) checkClusterAddrConflict(ctx context.Context, spec *specs.SpecCluster) (bool, *response.ClusterEndpoint) {

	clusters, err := a.GetClusters(ctx)
	if err != nil {
		return false, nil
	}

	ips := clusters.GetClusterIPs()
	v, ok := ips[spec.EndpointIP]
	return ok, &v
}

// allocateNewClusterIp - allocate new IP if spec uses cluster IP
// that already allocated.
func (a *TcaApi) allocateNewClusterIp(ctx context.Context, spec *specs.SpecCluster) error {

	clusters, err := a.GetClusters(ctx)
	if err != nil {
		return nil
	}

	allIps := clusters.GetClusterIPs()

	// check if addr conflicts
	var (
		maxCheck = 0
		ok       = true
		endPoint response.ClusterEndpoint
		addr     = spec.EndpointIP
	)

	for ok != false || maxCheck < 16 {

		endPoint, ok = allIps[addr]
		// if no overlap we done
		if !ok {
			spec.EndpointIP = addr
			return nil
		}

		// otherwise check if it IP or FQDN name conflicts
		// if IP compute next address and check
		if endPoint.IsIP {
			ip := net.ParseIP(addr)
			if ip == nil {
				return fmt.Errorf("failed parse IP address")
			}

			nextAddr := netutils.NextIPv4(ip, 1)
			if nextAddr != nil {
				addr = nextAddr.String()
			}
		} else {
			// TODO when TCA will support FQDN add FQDN support.
		}
		maxCheck++
	}

	return nil
}

// CreateClusters - creates new cluster , in Dry Run method only will do
// specString validation. isBlocking will indicate if caller expects cluster
// task to finish. verbose will output status of each task,
// linked to specific cluster creation.
func (a *TcaApi) CreateClusters(ctx context.Context, req *ClusterCreateApiReq) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if req == nil {
		return nil, api_errors.NewInvalidSpec("nil req")
	}

	if req.Spec == nil {
		return nil, api_errors.NewInvalidSpec("new cluster specString can't be nil")
	}

	// fix specString
	normalizeSpec(req.Spec)

	// validate cluster
	if err := validateClusterSpec(req.Spec); err != nil {
		return nil, err
	}

	// do all sanity check here.
	tenants, err := a.rest.GetClusters(ctx)
	if err != nil {
		return nil, err
	}

	_, err = tenants.GetClusterId(req.Spec.Name)
	// swap name if needed
	if err == nil {
		req.Spec.Name = req.Spec.Name + "-" + uuid.New().String()
		req.Spec.Name = req.Spec.Name[0:25]
		glog.Infof("Duplicate name regenerated new name '%v'", req.Spec.Name)
	}

	if ok, _ := a.checkClusterAddrConflict(ctx, req.Spec); ok {
		return nil, fmt.Errorf("cluster IP %s already in use", req.Spec.EndpointIP)
	}

	// resolve template id, and cluster type
	req.Spec.ClusterTemplateId, err = a.NormalizeTemplateId(req.Spec.ClusterTemplateId, req.Spec.ClusterType)
	if err != nil {
		return nil, err
	}

	// get template and validate specs
	t, err := a.rest.GetClusterTemplate(req.Spec.ClusterTemplateId)
	if err != nil {
		return nil, err
	}

	_, err = t.ValidateSpec(req.Spec)
	if err != nil {
		return nil, err
	}

	if req.IsFixConflict {
		err := a.allocateNewClusterIp(ctx, req.Spec)
		if err != nil {
			return nil, fmt.Errorf("failed resolve cluster ip conflict error %v", err)
		}
	}

	glog.Infof("Resolved template id %v", req.Spec.ClusterTemplateId)
	tenant, err := a.validateCloudEndpoint(ctx, req.Spec.HcxCloudUrl)
	if err != nil {
		return nil, err
	}
	req.Spec.HcxCloudUrl = tenant.HcxCloudURL

	err = a.validateTenant(tenant)
	if err != nil {
		return nil, err
	}

	identified := false
	if req.Spec.IsWorkload() {
		glog.Infof("Validating workload cluster specString")

		// resolve template id, in case client used name instead id
		mspec, err := a.ResolveManagementCluster(req.Spec.ManagementClusterId, tenants)
		if err != nil {
			return nil, err
		}

		req.Spec.ManagementClusterId = mspec.Id
		err = a.validatePlacements(ctx, req.Spec, tenant)
		if err != nil {
			return nil, err
		}

		identified = true
	}

	if req.Spec.IsManagement() {
		glog.Infof("Validating management cluster specString")
		// ignoring mgmt cluster id
		req.Spec.ManagementClusterId = ""
		err = a.validatePlacements(ctx, req.Spec, tenant)
		if err != nil {
			return nil, err
		}
		identified = true
	}

	if !identified {
		return nil, fmt.Errorf("invalid cluster type %s", req.Spec.ClusterType)
	}

	req.Spec.ClusterPassword = b64.StdEncoding.EncodeToString([]byte(req.Spec.ClusterPassword))

	glog.Infof("SpecCluster specString validated.")

	if req.IsDryRun {
		return &models.TcaTask{}, nil
	}

	task, err := a.rest.CreateCluster(req.Spec)
	if err != nil {
		return nil, err
	}

	if req.IsBlocking {
		err := a.BlockWaitTaskFinish(context.Background(), task, TaskStateSuccess, BlockMaxRetryTimer, req.IsBlocking)
		if err != nil {
			return task, err
		}
	}

	return task, err
}

// DeleteCluster - Method deletes cluster, note in order delete
// it must not have anything running on it. In order delete management
// cluster, all tenant cluster must be deleted first.
func (a *TcaApi) DeleteCluster(ctx context.Context, req *ClusterDeleteApiReq) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	var (
		cid = req.Cluster
		err error
	)

	if len(cid) == 0 {
		return nil, fmt.Errorf("empty cluster id or name")
	}

	// resolve name to id , if it not UUID
	if !IsValidUUID(cid) {
		cid, err = a.ResolveClusterName(ctx, cid)
		if err != nil {
			return nil, err
		}
	}

	task, err := a.rest.DeleteCluster(ctx, cid)
	if err != nil {
		return nil, err
	}

	// block and wait task to finish
	if req.IsBlocking {
		err := a.BlockWaitTaskFinish(ctx, task, TaskStateSuccess, BlockMaxRetryTimer, req.IsVerbose)
		if err != nil {
			return task, err
		}
	}

	return task, nil
}
