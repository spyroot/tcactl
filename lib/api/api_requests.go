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

// InstanceRequestSpec new instance request
type InstanceRequestSpec struct {
	cloudName    string
	clusterName  string
	vimType      string
	nfdName      string
	repo         string
	useAttached  bool
	repoUsername string
	repoPassword string
	instanceName string
	nodePoolName string
	namespace    string
	flavorName   string
	description  string

	disableGrant        bool
	ignoreGrantFailure  bool
	disableAutoRollback bool
}

func (i *InstanceRequestSpec) CloudName() string {
	return i.cloudName
}

func (i *InstanceRequestSpec) SetCloudName(cloudName string) {
	i.cloudName = cloudName
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

func (i *InstanceRequestSpec) DisableGrant() bool {
	return i.disableGrant
}

func (i *InstanceRequestSpec) SetDisableGrant(DisableGrant bool) {
	i.disableGrant = DisableGrant
}

func (i *InstanceRequestSpec) IgnoreGrantFailure() bool {
	return i.ignoreGrantFailure
}

func (i *InstanceRequestSpec) SetIgnoreGrantFailure(IgnoreGrantFailure bool) {
	i.ignoreGrantFailure = IgnoreGrantFailure
}

func (i *InstanceRequestSpec) DisableAutoRollback() bool {
	return i.disableAutoRollback
}

func (i *InstanceRequestSpec) SetDisableAutoRollback(DisableAutoRollback bool) {
	i.disableAutoRollback = DisableAutoRollback
}
