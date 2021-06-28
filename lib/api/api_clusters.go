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
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
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
func (a *TcaApi) GetClusterTask(clusterId string, showChildren bool) (*models.ClusterTask, error) {

	var err error
	_clusterid := clusterId

	if !IsValidUUID(_clusterid) {
		glog.Infof("Resolving cluster id from name %s", clusterId)
		_clusterid, err = a.ResolveClusterName(_clusterid)
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
