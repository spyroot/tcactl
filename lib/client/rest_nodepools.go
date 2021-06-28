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
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/io"
	"net/http"
)

type SupportedVersion struct {
	Items []struct {
		ClusterType       string        `json:"clusterType" yaml:"cluster_type"`
		SupportedVersions []interface{} `json:"supportedVersions" yaml:"supported_versions"`
	} `json:"items"`
}

type PasswordUpdateSpec struct {
	ExistingClusterPassword string `json:"existingClusterPassword" yaml:"existingClusterPassword"`
	ClusterPassword         string `json:"clusterPassword" yaml:"clusterPassword"`
}

//
//type NodeUpgradeTask struct {
//	Id          string `json:"id" yaml:"id"`
//	OperationId string `json:"operationId" yaml:"operation_id"`
//}
//
//type NodePoolDeleteTask struct {
//	Id          string `json:"id"`
//	OperationId string `json:"operationId"`
//}

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

	if c.isTrace && resp != nil {
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

	var pools response.NodesSpecs
	if err := json.Unmarshal(resp.Body(), &pools); err != nil {
		return nil, err
	}

	return &pools, nil
}

// CreateNewNodePool api method create a new node pool for
// a target cluster, method return models.TcaTask that progress
// can be monitored.
func (c *RestClient) CreateNewNodePool(r *request.NewNodePoolSpec, clusterId string) (*models.TcaTask, error) {

	if len(clusterId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	c.GetClient()

	io.PrettyPrint(r)

	resp, err := c.Client.R().SetBody(r).
		Post(c.BaseURL + fmt.Sprintf(TcaInfraCreatPool, clusterId))

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkErrors(resp)
	}

	var pools models.TcaTask
	if err := json.Unmarshal(resp.Body(), &pools); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &pools, nil
}

// DeleteNodePool - delete a note pool
func (c *RestClient) DeleteNodePool(clusterId string, nodePoolId string) (*models.TcaTask, error) {

	if len(clusterId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	if len(nodePoolId) == 0 {
		return nil, fmt.Errorf("nodePool id is empty string")
	}

	c.GetClient()
	resp, err := c.Client.R().Delete(c.BaseURL + TcaInfraCluster + "/" + clusterId + "/nodepool/" + nodePoolId)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
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

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &task, nil
}

// UpdateNodePool - update a note pool
func (c *RestClient) UpdateNodePool(r *request.NewNodePoolSpec, clusterId string, nodePoolId string) (*models.TcaTask, error) {

	if len(clusterId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	if len(nodePoolId) == 0 {
		return nil, fmt.Errorf("nodePool id is empty string")
	}

	c.GetClient()
	resp, err := c.Client.R().SetBody(r).Put(c.BaseURL + fmt.Sprintf(TcaInfraUpdatePool, clusterId, nodePoolId))

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("server return error %v %v %v", errRes.Details, errRes.Message, string(resp.Body()))
			return nil, fmt.Errorf("server return error %v", errRes.Message)
		}
		glog.Errorf("server return unknown error %v %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &task, nil
}

// NodePoolRetryTask - retry task related to node pool
func (c *RestClient) NodePoolRetryTask(taskId string) (*models.TcaTask, error) {

	if len(taskId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	c.GetClient()
	resp, err := c.Client.R().Post(c.BaseURL + fmt.Sprintf(TcaInfraPoolRetry, taskId))

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkErrors(resp)
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &task, nil
}

// NodePoolAbortTask - retry task related to node pool
func (c *RestClient) NodePoolAbortTask(taskId string) (*models.TcaTask, error) {

	if len(taskId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	c.GetClient()
	resp, err := c.Client.R().Post(c.BaseURL + fmt.Sprintf(TcaInfraPoolAbort, taskId))

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkErrors(resp)
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &task, nil
}

// GetClusterCompatability returns k8s node pool detail.
func (c *RestClient) GetClusterCompatability() (*SupportedVersion, error) {

	c.GetClient()
	resp, err := c.Client.R().
		Get(c.BaseURL + TcaInfraSupportedVer)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
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

	var l SupportedVersion
	if err := json.Unmarshal(resp.Body(), &l); err != nil {
		return nil, err
	}

	return &l, nil
}

// UpdateClusterPassword update cluster password.
func (c *RestClient) UpdateClusterPassword(req *PasswordUpdateSpec, clusterId string) (*models.TcaTask, error) {

	c.GetClient()
	resp, err := c.Client.R().SetBody(req).
		Put(c.BaseURL + fmt.Sprintf(TcaClusterChangePassword, clusterId))

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
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

	var l models.TcaTask
	if err := json.Unmarshal(resp.Body(), &l); err != nil {
		return nil, err
	}

	return &l, nil
}

// UpgradeNodePool - delete a note pool
func (c *RestClient) UpgradeNodePool(clusterId string, nodePoolId string) (*models.TcaTask, error) {

	if len(clusterId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	if len(nodePoolId) == 0 {
		return nil, fmt.Errorf("nodePool id is empty string")
	}

	c.GetClient()
	resp, err := c.Client.R().Post(c.BaseURL + fmt.Sprintf(TcaClustersNodePoolUpgrade, clusterId, nodePoolId))

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkErrors(resp)
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &task, nil
}
