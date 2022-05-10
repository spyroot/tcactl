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
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	ioutils "github.com/spyroot/tcactl/pkg/io"
	"net/http"
)

// GetClusters returns infrastructure k8s clusters
func (c *RestClient) GetClusters(ctx context.Context) (*response.Clusters, error) {

	glog.Infof("Retrieving cluster list.")

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).Get(c.BaseURL + TcaInfraClusters)
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
func (c *RestClient) GetCluster(ctx context.Context, clusterId string) (*response.ClusterSpec, error) {

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).Get(c.BaseURL + TcaInfraClusters + "/" + clusterId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkErrors(resp)
	}

	var clusters response.ClusterSpec
	if err := json.Unmarshal(resp.Body(), &clusters); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	return &clusters, nil
}

// GetClusterNodePools - returns cluster k8s node pools list
// each list hold specs, caller need indicate valid cluster ID not a name.
func (c *RestClient) GetClusterNodePools(clusterId string) (*response.NodePool, error) {

	glog.Infof("Retrieving node pool for cluster id %v", clusterId)

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

	var pools response.NodePool
	if err := json.Unmarshal(resp.Body(), &pools); err != nil {
		fmt.Println(resp.Body())
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	glog.Infof("return node pool list %d size", len(pools.Pools))
	return &pools, nil
}

// GetNamedClusterNodePools method first resolves cluster name
// and than look up node pool list attached to a given cluster.
func (c *RestClient) GetNamedClusterNodePools(ctx context.Context, clusterName string) (*response.NodePool, string, error) {

	// cluster list.
	cluster, err := c.GetClusters(ctx)
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

type TaskNotFound struct {
	errMsg string
}

func (m *TaskNotFound) Error() string {
	return m.errMsg + " cluster not found"
}

// GetClustersTask - returns infrastructure k8s clusters task list
// Before adjusting cluster task , caller must first check existing task list.
// each task can fail.
func (c *RestClient) GetClustersTask(ctx context.Context, f *specs.ClusterTaskQuery) (*models.ClusterTask, error) {

	c.GetClient()

	resp, err := c.Client.R().SetContext(ctx).SetBody(f).Post(c.BaseURL + TcaInfraClusterTask)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkError(resp)
	}

	var task models.ClusterTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	return &task, nil
}

// GetClusterTask returns k8s execution current task.
func (c *RestClient) GetClusterTask(ctx context.Context, clusterId string) (*models.ClusterTask, error) {

	if len(clusterId) == 0 {
		return nil, fmt.Errorf("cluster id is empty string")
	}

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).Get(c.BaseURL + fmt.Sprintf(TcaClusterTask, clusterId))

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

	var task models.ClusterTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	return &task, nil
}

// CreateCluster - create new management and or tenant cluster
// TCA require first create management cluster before any tenant cluster created.
// TCA also require each type of cluster has valid template.
// raw rest call doesn't do any validation, while API interface does basic validation check
// for spec.
func (c *RestClient) CreateCluster(spec *specs.SpecCluster) (*models.TcaTask, error) {

	c.GetClient()
	glog.Infof("Creating cluster %v", spec)

	if spec == nil {
		glog.Error("cluster spec is nil")
		return nil, errors.New("cluster spec is nil")
	}

	ioutils.PrettyString(spec)

	resp, err := c.Client.R().SetBody(spec).Post(c.BaseURL + TcaInfraClusters)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkErrors(resp)
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	glog.Infof("SpecCluster create task created task id %s op id %s", task.Id, task.OperationId)

	return &task, nil
}

// DeleteCluster - delete  k8s clusters. ( Management and tenant workload)
// note Management cluster can't deleted if there are workload clusters
// already attached.  Method return models.TcaTask that can monitored.
func (c *RestClient) DeleteCluster(ctx context.Context, clusterId string) (*models.TcaTask, error) {

	c.GetClient()
	glog.Infof("Deleting cluster %v", clusterId)

	resp, err := c.Client.R().SetContext(ctx).Delete(c.BaseURL + TcaInfraClusters + "/" + clusterId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkErrors(resp)
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	glog.Infof("SpecCluster delete task created task id %s op id %s", task.Id, task.OperationId)

	return &task, nil
}
