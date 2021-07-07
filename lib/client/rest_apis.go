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

const (

	// DefaultContentType content type used
	defaultContentType = "application/json"

	// uriAuthorize
	uriAuthorize = "/hybridity/api/sessions"

	// authorizationHeader - TCA authorization header
	authorizationHeader = "x-hm-authorization"

	// DefaultVersion api version used
	defaultVersion = "2"

	// DefaultAccept default content type we receive
	defaultAccept = "application/json"

	// apiTenants - list of vim tenant rest call
	apiTenants = "/hybridity/api/vims/v1/tenants"

	// TcaDeleteTenant - api call delete tenant
	TcaDeleteTenant = "/hybridity/api/vims/v1/tenants/%s"

	TcaVimTenant = "/hybridity/api/vims/v1/%s/tenants"

	// apiVim - attached vim list rest call
	apiVim = "/hybridity/api/vims/v1/"

	// apiTenantAction query action
	apiTenantAction = "action=query"

	// TcaInfraCluster api endpoint clusters api.
	TcaInfraCluster = "/hybridity/api/infra/k8s/cluster"

	// TcaInfraClusters - api clusters
	TcaInfraClusters = "/hybridity/api/infra/k8s/clusters"

	// TcaInfraClusterTask - api cluster task
	TcaInfraClusterTask = "/hybridity/api/infra/k8s/tasks"

	// TcaClusterTask - API query current cluster task
	TcaClusterTask = "/hybridity/api/infra/k8s/clusters/%s/tasks"

	// TcaApiVnfLcmExtensionVnfInstance api return vnf instance
	TcaApiVnfLcmExtensionVnfInstance = "/telco/api/vnflcm/v2/extension/vnf_instances"

	TcaApiVnfLcmVnfInstance = "/telco/api/vnflcm/v2/vnf_instances"

	// TcaVmwareExtensions extensions api
	TcaVmwareExtensions = "/hybridity/api/extensions"

	// TcaVmwareGetExtensions  api call deletes extension
	TcaVmwareGetExtensions = "/hybridity/api/extensions/%s"

	// TcaVmwareDeleteExtensions  api call deletes extension
	TcaVmwareDeleteExtensions = "/hybridity/api/extensions/%s"

	// TcaVmwareUpdateExtensions  api call deletes extension
	TcaVmwareUpdateExtensions = "/hybridity/api/extensions/%s"

	//TcaVmwareExtensionsTypes extensions types
	TcaVmwareExtensionsTypes = "/hybridity/api/extensions/types"

	// TcaVmwareNfvNetworks api request to get networks
	TcaVmwareNfvNetworks = "/hybridity/api/nfv/networks"

	// TcaVmwareRepositories - getter for repositories
	TcaVmwareRepositories = "/hybridity/api/repositories/query"

	// TcaVmwareRepos rest call to manipulate repositories
	TcaVmwareRepos = "/hybridity/api/repositories"

	// TcaVmwarePackagesActionContent pacakge content suffix
	TcaVmwarePackagesActionContent = "package_content"

	// TcaVmwareTelcoPackages api endpoint for vnf packages.
	TcaVmwareTelcoPackages = "/telco/api/vnfpkgm/v2/vnf_packages"

	// TcaVmwareVnflcmInstances api endpoint LCM instances
	TcaVmwareVnflcmInstances = "/telco/api/vnflcm/v2/vnf_instances"

	// TcaVmwareVnflcmInstance operation on instance
	TcaVmwareVnflcmInstance = "/telco/api/vnflcm/v2/vnf_instances/%s"

	//TcaVmwareVnflcmInstantiate instantiate
	TcaVmwareVnflcmInstantiate = "/telco/api/vnflcm/v2/vnf_instances/%s/instantiate"

	//TcaVmwareVnflcmUpdate update state
	TcaVmwareVnflcmUpdate = "/hybridity/api/vnflcm/v1/vnf_instances/%s/update_state"

	// TcaInfraPoolRetry Retry task
	TcaInfraPoolRetry = "/hybridity/api/infra/k8s/operations/%s/retry"

	// TcaInfraPoolAbort abort task
	TcaInfraPoolAbort = "/hybridity/api/infra/k8s/operations/%s/abort"

	// TcaInfraCreatPool API call create node pool on target cluster
	TcaInfraCreatPool = "/hybridity/api/infra/k8s/cluster/%s/nodepool"

	TcaClustersNodePool = "/hybridity/api/infra/k8s/cluster/%s/nodepool/%s"

	// TcaInfraDeletePool API call delete Node pool
	TcaInfraDeletePool = TcaClustersNodePool

	// TcaInfraUpdatePool API call update Node pool
	TcaInfraUpdatePool = "/hybridity/api/infra/k8s/cluster/%s/nodepool/%s"

	// TcaClustersNodePoolUpgrade upgrade existing node pool
	TcaClustersNodePoolUpgrade = "/hybridity/infra/k8s/cluster/%s/nodepool/%s/upgrade"

	// TcaInfraSupportedVer return list of supported version
	TcaInfraSupportedVer = "/hybridity/api/infra/k8s/supportedK8sVersions"

	// TcaClusterChangePassword change password for existing cluster
	TcaClusterChangePassword = "/hybridity/api/infra/k8s/clusters/%s/changePassword"
)
