// Package app
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

import "github.com/spyroot/tcactl/lib/client/specs"

// NodePoolCreateApiReq api request issued to create new node pool
type NodePoolCreateApiReq struct {

	// Spec is a specs.SpecNodePool
	Spec *specs.SpecNodePool

	// Cluster is cluster name or cluster id
	Cluster string

	// dry run or apply change
	IsDryRun bool

	// if node pool creation needs to block
	IsBlocking bool

	// if blocking request require output progress
	IsVerbose bool
}

// ClusterDeleteApiReq api request issued to delete cluster
type ClusterDeleteApiReq struct {
	// cluster name or id
	Cluster string
	//
	IsBlocking bool
	//
	IsVerbose bool
}

// ClusterCreateApiReq api request issued to create new cluster
type ClusterCreateApiReq struct {

	// Spec hold a pointer to specs.SpecCluster and must not be nil
	Spec *specs.SpecCluster

	//
	IsDryRun bool

	//
	IsBlocking bool

	//
	IsVerbose bool

	// isFixConflict if true , api create method will try to resolve cluster IP conflict.
	IsFixConflict bool
}

// CreateInstanceApiReq - api request to create new cnf or vnf instance
type CreateInstanceApiReq struct {

	//InstanceName instance name
	InstanceName string

	//PoolName node pool name
	PoolName string

	//VimName target vim name
	VimName string

	//ClusterName target cluster name
	ClusterName string

	//IsBlocking block or async task
	IsBlocking bool

	// if blocking request require output progress
	IsVerbose bool

	// additional param
	AdditionalParam *specs.AdditionalParams

	//Namespace  overwrite name
	Namespace string

	//RepoUsername overwrites Repo username
	RepoUsername string

	//RepoPassword overwrite Repo password
	RepoPassword string

	// RepoUrl overwrite Repo url
	RepoUrl string
}

// TerminateInstanceApiReq - api request
// to terminate CNF or VNF instance.
type TerminateInstanceApiReq struct {

	//InstanceName instance name
	InstanceName string

	//ClusterName target cluster name
	ClusterName string

	//IsBlocking block or async task
	IsBlocking bool

	// if blocking request require output progress
	IsVerbose bool
}

// UpdateInstanceApiReq api request to update a LCM state of existing
// CNF or VNF instance.
type UpdateInstanceApiReq struct {

	//
	UpdateReq *specs.LcmInstantiateRequest

	//InstanceName instance name
	InstanceName string

	//PoolName node pool name
	PoolName string

	//ClusterName target cluster name
	ClusterName string

	//IsBlocking block or async task
	IsBlocking bool

	// if a request is blocking, and caller requires output progress
	IsVerbose bool
}

// ResetInstanceApiReq api request to reset existing CNF or VNF instance.
type ResetInstanceApiReq struct {

	//InstanceName instance name or id
	InstanceName string

	//ClusterName target cluster name
	ClusterName string

	//IsBlocking block or async task
	IsBlocking bool

	// if a request is blocking, and caller requires output progress
	IsVerbose bool
}

// NewInstanceRequestSpec return new instance request spec
func NewInstanceRequestSpec(cloudName string, clusterName string, vimType string, nfdName string,
	repo string, instanceName string, nodePoolName string) (*specs.InstanceRequestSpec, error) {
	i := &specs.InstanceRequestSpec{
		CloudName:        cloudName,
		ClusterName:      clusterName,
		VimType:          vimType,
		NfdName:          nfdName,
		Repo:             repo,
		InstanceName:     instanceName,
		NodePoolName:     nodePoolName,
		UseLinkedRepo:    true,
		AdditionalParams: specs.AdditionalParams{}}

	i.FlavorName = DefaultNamespace
	i.Description = ""
	i.Namespace = DefaultFlavor

	err := i.Validate()
	if err != nil {
		return nil, err
	}

	return i, nil
}
