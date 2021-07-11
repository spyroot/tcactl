// Package specs
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
package specs

//{"filter":{"vimId":"vmware_8BF6253CE6E247018D909605A437B827","nodePoolId":"7f6f853c-7a87-4122-81df-f7a65186674c","type":["DistributedVirtualPortgroup","OpaqueNetwork"]}}

type NetworkFilter struct {
	Filter struct {
		VimId      string   `json:"vimId" yaml:"vimId"`
		NodePoolId string   `json:"nodePoolId" yaml:"nodePoolId"`
		Type       []string `json:"type" yaml:"type"`
	} `json:"filter"`
}

// TaskFilter Task filter Query filter
type TaskFilter struct {
	Filter struct {
		EntityIds []string `json:"entityIds" yaml:"entityIds"`
	} `json:"filter" yaml:"filter"`
}

// TenantFilter filter by nfType nad NfdId
type TenantFilter struct {
	NfType string `json:"nfType,omitempty" yaml:"nfType"`
	NfdId  string `json:"nfdId,omitempty" yaml:"nfdId"`
}

// TenantsNfFilter filter
type TenantsNfFilter struct {
	Filter TenantFilter `json:"filter"`
}
