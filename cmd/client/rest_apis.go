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
	// cluster
	TcaInfraCluster = "/hybridity/api/infra/k8s/cluster"

	TcaInfraClusters = "/hybridity/api/infra/k8s/clusters"

	TcaInfraClusterTask = "/hybridity/api/infra/k8s/tasks"

	TcaApiVnfLcmExtensionVnfInstance = "/telco/api/vnflcm/v2/extension/vnf_instances"

	TcaApiVnfLcmVnfInstance = "/telco/api/vnflcm/v2/vnf_instances"

	// extensions
	TcaVmwareExtensions = "/hybridity/api/extensions"

	TcaVmwareExtensionsTypes = "/hybridity/api/extensions/types"

	TcaVmwareNfvNetworks = "/hybridity/api/nfv/networks"

	TcaVmwareRepositories = "/hybridity/api/repositories/query"

	// TcaVmwareRepos rest call fo
	TcaVmwareRepos = "/hybridity/api/repositories"

	//	"/telco/api/vnfpkgm/v2/vnf_packages/22495560-fdd9-4e73-b0b7-774629da2050/package_content"
	//
	TcaVmwarePackages = "/telco/api/vnfpkgm/v2/vnf_packages"
	//
	TcaVmwarePackagesActionContent = "package_content"

	// TcaVmwareTelcoPackages
	TcaVmwareTelcoPackages = "/telco/api/vnfpkgm/v2/vnf_packages"
)
