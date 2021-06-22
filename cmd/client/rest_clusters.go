// Package client
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
package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/client/response"
	"github.com/spyroot/hestia/cmd/models"
	"net/http"
)

// GetClusters returns infrastructure k8s clusters
func (c *RestClient) GetClusters() (*response.Clusters, error) {

	glog.Infof("Retrieving cluster list.")

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + TcaInfraClusters)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("server return error %v", errRes.Message)
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var clusters response.Clusters
	if err := json.Unmarshal(resp.Body(), &clusters.Clusters); err != nil {
		glog.Errorf("Failed parse respond body. %v", err)
		return nil, err
	}

	ids, err := clusters.GetClusterIds()
	glog.Infof("Retrieved cluster list. %v", ids)

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

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("server return error %v", errRes.Message)
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v %v", resp.StatusCode(), string(resp.Body()))
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

	glog.Infof("Sending a query node pool for cluster id %v", clusterId)
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

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("server return error %v %v %v", errRes.Details, errRes.Message, string(resp.Body()))
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var pools response.NodePool
	if err := json.Unmarshal(resp.Body(), &pools); err != nil {
		return nil, err
	}

	return &pools, nil
}

// GetNamedClusterNodePools method first resolves cluster name
// and than look up node pool list attached to a given cluster.
func (c *RestClient) GetNamedClusterNodePools(clusterName string) (*response.NodePool, string, error) {

	// cluster list.
	cluster, err := c.GetClusters()
	if err != nil || cluster == nil {
		glog.Errorf("Failed acquire cluster information for %v", clusterName)
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
			glog.Errorf("server return error %v %v %v", errRes.Details, errRes.Message, string(resp.Body()))
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var pools response.NodesSpecs
	if err := json.Unmarshal(resp.Body(), &pools); err != nil {
		return nil, err
	}

	return &pools, nil
}

type TaskNotFound struct {
	errMsg string
}

func (m *TaskNotFound) Error() string {
	return m.errMsg + " cluster not found"
}

// GetClusterTask - returns infrastructure k8s clusters
func (c *RestClient) GetClusterTask(f *request.ClusterTaskQuery) (*models.ClusterTask, error) {

	c.GetClient()

	resp, err := c.Client.R().SetBody(f).Post(c.BaseURL + TcaInfraClusterTask)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var task models.ClusterTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &task, nil
}

// CreateCluster - returns infrastructure k8s clusters
func (c *RestClient) CreateCluster(spec *request.Cluster) (bool, error) {

	c.GetClient()
	glog.Infof("Creating cluster")

	resp, err := c.Client.R().SetBody(spec).Post(c.BaseURL + "/hybridity/api/infra/k8s/clusters")
	if err != nil {
		glog.Error(err)
		return false, err
	}

	fmt.Println(string(resp.Body()))

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return false, c.checkErrors(resp)
	}

	return resp.StatusCode() == http.StatusOK, nil
}

// DeleteCluster - returns infrastructure k8s clusters
func (c *RestClient) DeleteCluster(clusterId string) (bool, error) {

	c.GetClient()

	resp, err := c.Client.R().Delete(c.BaseURL + TcaInfraClusters + "/" + clusterId)
	if err != nil {
		glog.Error(err)
		return false, err
	}

	//fmt.Println(resp.Body())

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return false, c.checkErrors(resp)
	}

	return resp.StatusCode() == http.StatusOK, nil
}
