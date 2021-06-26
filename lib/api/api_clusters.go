package api

import (
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/response"
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

// GetClusterNodePool -  method retrieve clusters node pool.
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
