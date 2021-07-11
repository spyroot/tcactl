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
	// Spec is a specs.NodePoolSpec spec
	Spec *specs.NodePoolSpec

	// Cluster is cluster name or cluster id
	Cluster string

	//
	IsDryRun bool

	// if node pool creation needs to block
	IsBlocking bool

	// if blocking request require output progress
	IsVerbose bool
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

// TerminateInstanceApiReq - api request to terminate new cnf or vnf instance
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

	// if blocking request require output progress
	IsVerbose bool
}

type ResetInstanceApiReq struct {
	//InstanceName instance name
	InstanceName string

	//ClusterName target cluster name
	ClusterName string

	//IsBlocking block or async task
	IsBlocking bool

	// if blocking request require output progress
	IsVerbose bool
}
