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
	"github.com/go-playground/validator/v10"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/api_errors"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	errnos "github.com/spyroot/tcactl/pkg/errors"
)

// GetNodePool api method returns a Node pool response.NodePool
// lookup based on cluster id or cluster name
func (a *TcaApi) GetNodePool(ctx context.Context, NameOrId string) (*response.NodePool, error) {

	resolvedId := NameOrId

	clusters, err := a.rest.GetClusters(ctx)
	if err != nil || clusters == nil {
		return nil, err
	}

	resolvedId, err = clusters.GetClusterId(resolvedId)
	if err != nil {
		return nil, err
	}

	pool, err := a.rest.GetClusterNodePools(resolvedId)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// ResolvePoolAndCluster - api method resolves cluster and node pool in pair
// return cluster id and pool id or error
func (a *TcaApi) ResolvePoolAndCluster(ctx context.Context,
	cluster string, nodePool string) (string, string, error) {

	_clusterId := cluster
	_nodepoolId := nodePool

	if IsValidUUID(cluster) {
		c, err := a.GetCluster(ctx, _clusterId)
		if err != nil {
			glog.Error(err)
			return "", "", err
		}
		_clusterId = c.ClusterId
	} else {
		c, err := a.ResolveClusterName(ctx, _clusterId)
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
		pools, err := a.GetNodePool(ctx, _clusterId)
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

// DeleteNodePool api call deletes node pool from a tenant cluster,
// cluster argument is cluster name or id.
// API call returns models.TcaTask that hold operation id that can be used for
// subsequent monitoring.
func (a *TcaApi) DeleteNodePool(ctx context.Context, cluster string, nodePool string) (*models.TcaTask, error) {

	if len(cluster) == 0 {
		return nil, api_errors.NewInvalidArgument("cluster")
	}

	if len(nodePool) == 0 {
		return nil, api_errors.NewInvalidArgument("nodePool")
	}

	_clusterId, _nodepoolId, err := a.ResolvePoolAndCluster(ctx, cluster, nodePool)
	if err != nil {
		return nil, err
	}
	task, err := a.rest.DeleteNodePool(_clusterId, _nodepoolId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// nodePoolReqValidator validate node pool spec
func (a *TcaApi) nodePoolReqValidator(spec *specs.SpecNodePool) error {

	if len(spec.PlacementParams) == 0 {
		return api_errors.NewInvalidSpec("node pool spec must contain placement params")
	}
	if len(spec.Networks) == 0 {
		return api_errors.NewInvalidSpec("node pool spec must network list")
	}
	if spec.Cpu == 0 {
		return api_errors.NewInvalidSpec("node pool spec contain zero cpu. Spec must contain same number of cpu.")
	}
	if spec.Memory == 0 {
		return api_errors.NewInvalidSpec("node pool spec contain zero memory. Spec must contain same memory value.")
	}
	if spec.Storage == 0 {
		return api_errors.NewInvalidSpec("node pool spec contain zero storage. Spec must contain same storage value")
	}

	if spec == nil {
		return api_errors.NewInvalidSpec("Spec is nil")
	}

	err := spec.Validate()
	if err != nil {
		return err
	}

	return nil
}

// CreateNewNodePool api call creates a new node pool for a given cluster name or id.
// both cluster can be named or ids.  API call returns models.TcaTas that monitored.
// isDry run provide capability only validate specs without creating actual node pool.
func (a *TcaApi) CreateNewNodePool(ctx context.Context, req *NodePoolCreateApiReq) (*models.TcaTask, error) {

	if a == nil {
		return nil, errnos.NilError
	}

	if req == nil {
		return nil, errnos.ReqNil
	}

	if req.Spec == nil {
		return nil, errnos.SpecNil
	}

	if err := a.nodePoolReqValidator(req.Spec); err != nil {
		return &models.TcaTask{}, err
	}

	_clusterId := req.Cluster

	if IsValidUUID(req.Cluster) {
		c, err := a.GetCluster(ctx, _clusterId)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		_clusterId = c.ClusterId
	} else {
		c, err := a.ResolveClusterName(ctx, _clusterId)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		_clusterId = c
	}

	// in dry we just return
	if req.IsDryRun {
		return &models.TcaTask{}, nil
	}

	specCopy := req.Spec
	specCopy.SpecType = ""

	task, err := a.rest.CreateNewNodePool(req.Spec, _clusterId)
	if err != nil {
		return nil, err
	}

	if req.IsBlocking {
		err := a.BlockWaitTaskFinish(context.Background(), task, TaskStateSuccess, BlockMaxRetryTimer, req.IsVerbose)
		if err != nil {
			return task, err
		}
	}

	return task, err
}

// UpdateNodePool api call updates a node pool for a existing cluster or node pool.
// Both cluster and node pool can be a pool name or pool ids.
//
// API call returns models.TcaTask that hold operation id.
// req *NodePoolCreateApiReq field
// 	isDry run provide capability to validate specs without applying any changes.
// 	isBlocking will block and wait when all task finish or fail.
// 	isVerbose flag will output progress message during each pool interval.
func (a *TcaApi) UpdateNodePool(ctx context.Context, req *NodePoolCreateApiReq) (*models.TcaTask, error) {

	if a == nil {
		return nil, errnos.NilError
	}

	if req == nil {
		return nil, errnos.ReqNil
	}

	if req.Spec == nil {
		return nil, errnos.SpecNil
	}

	if err := a.nodePoolReqValidator(req.Spec); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &models.TcaTask{}, validationErrors
	}

	_clusterId, _nodePoolId, err := a.ResolvePoolAndCluster(ctx, req.Cluster, req.Spec.Name)
	if err != nil {
		return nil, err
	}

	// in dry we just return
	if req.IsDryRun {
		return &models.TcaTask{}, nil
	}

	// update node pool id
	if len(req.Spec.Id) == 0 || !IsValidUUID(req.Spec.Id) {
		req.Spec.Id = _nodePoolId
	}

	specCopy := req.Spec
	specCopy.SpecType = ""

	task, err := a.rest.UpdateNodePool(req.Spec, _clusterId, _nodePoolId)
	if err != nil {
		return nil, err
	}

	if req.IsBlocking {
		err = a.BlockWaitTaskFinish(context.Background(), task, TaskStateSuccess, BlockMaxRetryTimer, req.IsVerbose)
		if err != nil {
			return task, err
		}
	}

	return task, err
}
