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

import "github.com/spyroot/tcactl/lib/client/request"

// NodePoolCreateApiReq api request issued to create new node pool
type NodePoolCreateApiReq struct {
	// Spec is a request.NewNodePoolSpec spec
	Spec *request.NewNodePoolSpec

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
	AdditionalParam *request.AdditionalParams

	//Namespace  overwrite name
	Namespace string

	//RepoUsername overwrites repo username
	RepoUsername string

	//RepoPassword overwrite repo password
	RepoPassword string

	// RepoUrl overwrite repo url
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
	UpdateReq *request.InstantiateVnfRequest

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

// InstanceRequestSpec new instance request
type InstanceRequestSpec struct {
	// target cloud name
	cloudName string
	// target cluster name
	clusterName string
	// vim name
	vimType string
	// catalog name
	nfdName string
	//
	useAttached bool
	//
	repo string
	//
	repoUsername string
	//
	repoPassword string
	//
	instanceName string
	//
	nodePoolName string
	// user linked repo
	useLinkedRepo bool
	// target namespace
	namespace string
	// flavor name
	flavorName string
	//
	description string
	// additional placement details
	AdditionalParams request.AdditionalParams
	// fix name conflict
	doAutoName bool
}

func (i *InstanceRequestSpec) CloudName() string {
	return i.cloudName
}

func (i *InstanceRequestSpec) SetCloudName(cloudName string) {
	i.cloudName = cloudName
}

func (i *InstanceRequestSpec) SetAutoName(f bool) {
	i.doAutoName = f
}

func (i *InstanceRequestSpec) IsAutoName() bool {
	return i.doAutoName
}

func (i *InstanceRequestSpec) ClusterName() string {
	return i.clusterName
}

func (i *InstanceRequestSpec) SetClusterName(clusterName string) {
	i.clusterName = clusterName
}

func (i *InstanceRequestSpec) VimType() string {
	return i.vimType
}

func (i *InstanceRequestSpec) SetVimType(vimType string) {
	i.vimType = vimType
}

func (i *InstanceRequestSpec) NfdName() string {
	return i.nfdName
}

func (i *InstanceRequestSpec) SetNfdName(nfdName string) {
	i.nfdName = nfdName
}

func (i *InstanceRequestSpec) Repo() string {
	return i.repo
}

func (i *InstanceRequestSpec) SetRepo(repo string) {
	i.repo = repo
}

func (i *InstanceRequestSpec) InstanceName() string {
	return i.instanceName
}

func (i *InstanceRequestSpec) SetInstanceName(instanceName string) {
	i.instanceName = instanceName
}

func (i *InstanceRequestSpec) NodePoolName() string {
	return i.nodePoolName
}

func (i *InstanceRequestSpec) SetNodePoolName(nodePoolName string) {
	i.nodePoolName = nodePoolName
}

func (i *InstanceRequestSpec) Namespace() string {
	return i.namespace
}

func (i *InstanceRequestSpec) SetNamespace(namespace string) {
	i.namespace = namespace
}

func (i *InstanceRequestSpec) FlavorName() string {
	return i.flavorName
}

func (i *InstanceRequestSpec) SetFlavorName(flavorName string) {
	i.flavorName = flavorName
}

func (i *InstanceRequestSpec) Description() string {
	return i.description
}

func (i *InstanceRequestSpec) SetDescription(description string) {
	i.description = description
}

func (i *InstanceRequestSpec) UseLinked() bool {
	return i.useLinkedRepo
}
