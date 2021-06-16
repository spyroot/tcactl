package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/client/response"
	"net/http"
)

// GetClusters returns infrastructure k8s clusters
func (c *RestClient) GetClusters() (*response.Clusters, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + TcaInfraClusters)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("server return error %v", errRes.Message)
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var clusters response.Clusters
	if err := json.Unmarshal(resp.Body(), &clusters.Clusters); err != nil {
		return nil, err
	}

	return &clusters, nil
}

// GetCluster returns infrastructure k8s clusters
func (c *RestClient) GetCluster(clusterId string) (*response.ClusterSpec, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + TcaInfraClusters + "/" + clusterId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("server return error %v", errRes.Message)
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var clusters response.ClusterSpec
	if err := json.Unmarshal(resp.Body(), &clusters); err != nil {
		return nil, err
	}

	return &clusters, nil
}

// GetClusterNodePools - returns cluster k8s node pools list
// each list hold specs, caller need indicate valid cluster ID not a name.
func (c *RestClient) GetClusterNodePools(clusterId string) (*response.NodePool, error) {

	if len(clusterId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	c.GetClient()

	resp, err := c.Client.R().
		Get(c.BaseURL + TcaInfraCluster + "/" + clusterId + "/nodepools")

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("server return error %v", errRes.Details, errRes.Message, string(resp.Body()))
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var pools response.NodePool
	if err := json.Unmarshal(resp.Body(), &pools); err != nil {
		return nil, err
	}

	return &pools, nil
}

// GetNamedClusterNodePools method first resolves cluster name and than
// look up node pool list.
func (c *RestClient) GetNamedClusterNodePools(clusterName string) (*response.NodePool, string, error) {

	cluster, err := c.GetClusters()

	if err != nil || cluster == nil {
		glog.Error("Failed acquire cluster information for %v", clusterName)
		return nil, "", err
	}

	// get cluster id for a pool
	clusterId, err := cluster.GetClusterId(clusterName)
	if err != nil {
		return nil, "", err
	}

	nodePool, err := c.GetClusterNodePools(clusterId)
	if err != nil {
		return nil, "", err
	}

	return nodePool, clusterId, nil
}

// GetClusterNodePool returns k8s node pool detail.
func (c *RestClient) GetClusterNodePool(clusterId string, nodePoolId string) (*response.NodesSpecs, error) {

	if len(clusterId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	if len(nodePoolId) == 0 {
		return nil, fmt.Errorf("node pool is empty string")
	}

	c.GetClient()
	resp, err := c.Client.R().
		Get(c.BaseURL + TcaInfraCluster + "/" + clusterId + "/nodepool/" + nodePoolId)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("server return error %v", errRes.Details, errRes.Message, string(resp.Body()))
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var pools response.NodesSpecs
	if err := json.Unmarshal(resp.Body(), &pools); err != nil {
		return nil, err
	}

	return &pools, nil
}
