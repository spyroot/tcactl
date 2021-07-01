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
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	errnos "github.com/spyroot/tcactl/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// ReadNodeSpecFromFile - Read node template from file
func ReadNodeSpecFromFile(fileName string) (*request.NewNodePoolSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadNodeSpecSpec(file)
}

// ReadNodeSpecFromString - Read node template from string
func ReadNodeSpecFromString(str string) (*request.NewNodePoolSpec, error) {
	r := strings.NewReader(str)
	return ReadNodeSpecSpec(r)
}

// ReadNodeSpecSpec - Read node template spec
func ReadNodeSpecSpec(b io.Reader) (*request.NewNodePoolSpec, error) {

	var spec request.NewNodePoolSpec

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		glog.Error(err)
		return &spec, nil
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		glog.Error(err)
		return &spec, nil
	}

	return nil, &InvalidSpec{"Failed to parse input spec."}
}

// GetNodePool return a Node pool for particular cluster
func (a *TcaApi) GetNodePool(clusterId string) (*response.NodePool, error) {

	if a.rest == nil {
		return nil, errnos.RestNilError
	}

	clusters, err := a.rest.GetClusters()
	if err != nil || clusters == nil {
		return nil, err
	}

	clusterId, err = clusters.GetClusterId(clusterId)
	if err != nil {
		return nil, err
	}

	pool, err := a.rest.GetClusterNodePools(clusterId)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// ResolvePoolAndCluster - resolve both cluster and node pool in pair
// All cluster plus node pool require cluster and node pool id.
// This method remove constrain on name but be careful for
// name duplicates.
func (a *TcaApi) ResolvePoolAndCluster(cluster string, nodePool string) (string, string, error) {

	if a.rest == nil {
		return "", "", errnos.RestNilError
	}

	_clusterId := cluster
	_nodepoolId := nodePool

	if IsValidUUID(cluster) {
		c, err := a.GetCluster(_clusterId)
		if err != nil {
			glog.Error(err)
			return "", "", err
		}
		_clusterId = c.ClusterId
	} else {
		c, err := a.ResolveClusterName(_clusterId)
		if err != nil {
			glog.Error(err)
			return "", "", err
		}
		_clusterId = c
	}

	if IsValidUUID(_nodepoolId) {
		pool, err := a.rest.GetClusterNodePool(_clusterId, _nodepoolId)
		if err != nil {
			glog.Error(err)
			return "", "", err
		}
		_nodepoolId = pool.Id
	} else {
		pools, err := a.GetNodePool(_clusterId)
		if err != nil {
			glog.Error(err)
			return "", "", err
		}

		pool, err := pools.GetPool(_nodepoolId)
		if err != nil {
			glog.Error(err)
			return "", "", err
		}
		_nodepoolId = pool.Id
	}

	return _clusterId, _nodepoolId, nil
}

// DeleteNodePool api call deletes node pool from given cluster name or id.
// both cluster and node pool can be named or ids.
// API call returns models.TcaTas that monitored.
func (a *TcaApi) DeleteNodePool(cluster string, nodePool string) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, errnos.RestNilError
	}

	_clusterId, _nodepoolId, err := a.ResolvePoolAndCluster(cluster, nodePool)
	if err != nil {
		return nil, err
	}
	task, err := a.rest.DeleteNodePool(_clusterId, _nodepoolId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// updateNodePoolValidator
// spec Validator
func (a *TcaApi) nodePoolValidator(spec *request.NewNodePoolSpec) error {

	if len(spec.PlacementParams) == 0 {
		return errors.New("spec must contain placement params")
	}

	if len(spec.Networks) == 0 {
		return errors.New("spec must network list")
	}

	if spec.Cpu == 0 {
		return errors.New("spec contain zero cpu. Spec must contain same number of cpu.")
	}
	if spec.Memory == 0 {
		return errors.New("spec contain zero memory. Spec must contain same memory value.")
	}
	if spec.Storage == 0 {
		return errors.New("spec contain zero storage. Spec must contain same storage value")
	}

	err := a.specValidator.Struct(spec)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	return nil
}

// CreateNewNodePool api call creates a new node pool for a given cluster name or id.
// both cluster can be named or ids.  API call returns models.TcaTas that monitored.
// isDry run provide capability only validate specs without creating actual node pool.
func (a *TcaApi) CreateNewNodePool(spec *request.NewNodePoolSpec,
	cluster string, isDry bool, isBlocking bool, verbose bool) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, errnos.RestNilError
	}

	if err := a.nodePoolValidator(spec); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &models.TcaTask{}, validationErrors
	}

	_clusterId := cluster

	if IsValidUUID(cluster) {
		c, err := a.GetCluster(_clusterId)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		_clusterId = c.ClusterId
	} else {
		c, err := a.ResolveClusterName(_clusterId)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		_clusterId = c
	}

	err := a.specValidator.Struct(spec)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &models.TcaTask{}, validationErrors
	}

	// in dry we just return
	if isDry {
		return &models.TcaTask{}, nil
	}

	task, err := a.rest.CreateNewNodePool(spec, _clusterId)
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

// UpdateNodePool api call updates a node pool for a existing cluster or node pool.
// Both cluster and node pool can be a named or ids.
//
// API call returns models.TcaTask that letter can monitored.
// isDry run provide capability only validate specs without
// creating actual node pool.
// isBlocking will block and wait when all task
// finish to run
// Verbose flag will output glog info message during wait time.
func (a *TcaApi) UpdateNodePool(spec *request.NewNodePoolSpec,
	cluster string, nodePool string, isDry bool, isBlocking bool, verbose bool) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, errnos.RestNilError
	}
	if a.rest == nil {
		return nil, errnos.RestNilError
	}

	if err := a.nodePoolValidator(spec); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &models.TcaTask{}, validationErrors
	}

	_clusterId, _notebookId, err := a.ResolvePoolAndCluster(cluster, nodePool)
	if err != nil {
		return nil, err
	}

	// in dry we just return
	if isDry {
		return &models.TcaTask{}, nil
	}

	// update node pool id
	if len(spec.Id) == 0 || IsValidUUID(spec.Id) == false {
		spec.Id = _notebookId
	}

	task, err := a.rest.UpdateNodePool(spec, _clusterId, _notebookId)
	if err != nil {
		return nil, err
	}

	if isBlocking {
		err = a.BlockWaitTaskFinish(context.Background(), task, TaskStateSuccess, BlockMaxRetryTimer, verbose)
		if err != nil {
			return task, err
		}
	}

	return task, err
}
