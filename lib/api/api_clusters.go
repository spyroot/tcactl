package api

import (
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
)

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

// GetClusterNodePool -  API Method retrieve clusters node pool.
// cluster identified either by name or uuid.
func (a *TcaApi) GetClusterNodePool(clusterId string, nodePoolId string) (*response.NodesSpecs, error) {

	var (
		_clusterid  = clusterId
		_nodePoolId = nodePoolId
		err         error
	)

	if !IsValidUUID(_clusterid) {
		_clusterid, err = a.ResolveClusterName(clusterId)
		if err != nil {
			return nil, err
		}
	}

	if !IsValidUUID(_nodePoolId) {
		_nodePoolId, err = a.ResolvePoolId(_nodePoolId)
		if err != nil {
			return nil, err
		}
	}

	return a.rest.GetClusterNodePool(clusterId, nodePoolId)
}

// GetClusterTask method return list task models.ClusterTask
// currently executing on given cluster
func (a *TcaApi) GetClusterTask(clusterid string, showChildren bool) (*models.ClusterTask, error) {

	var err error
	_clusterid := clusterid

	if !IsValidUUID(_clusterid) {
		glog.Infof("Resolving cluster id from name %s", clusterid)
		_clusterid, err = a.ResolveClusterName(_clusterid)
		if err != nil {
			return nil, err
		}
	}

	clusters, err := a.rest.GetClusters()
	if err != nil {
		return nil, nil
	}

	clusterSpec, err := clusters.GetClusterSpec(clusterid)
	if err != nil {
		return nil, err
	}

	r := request.NewClusterTaskQuery(clusterSpec.ManagementClusterId)
	r.IncludeChildTasks = showChildren

	return a.rest.GetClustersTask(r)
}
