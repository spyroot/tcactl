// Package request
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
package request

type TerminateVnfRequest struct {
	TerminationType            string `json:"terminationType" yaml:"terminationType"`
	GracefulTerminationTimeout int    `json:"gracefulTerminationTimeout" yaml:"gracefulTerminationTimeout"`
	AdditionalParams           struct {
		LcmInterfaces []struct {
			InterfaceName string `json:"interfaceName" yaml:"interfaceName"`
			Parameters    []struct {
				Name string `json:"name" yaml:"name"`
				Type string `json:"type" yaml:"type"`
			} `json:"parameters" yaml:"parameters"`
		} `json:"lcmInterfaces" yaml:"lcmInterfaces"`
	} `json:"additionalParams" yaml:"additionalParams"`
}

// TaskFilter Task filter Query filter
type TaskFilter struct {
	Filter struct {
		EntityIds []string `json:"entityIds" yaml:"entityIds"`
	} `json:"filter" yaml:"filter"`
}

// VduParam
type VduParam struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	RepoURL   string `json:"repoUrl" yaml:"repoUrl"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	VduName   string `json:"vduName" yaml:"vduName"`
}

type AdditionalParams struct {
	VduParams           []VduParam `json:"vduParams" yal:"vdu_params"`
	DisableGrant        bool       `json:"disableGrant" yal:"disableGrant"`
	IgnoreGrantFailure  bool       `json:"ignoreGrantFailure" yal:"ignoreGrantFailure"`
	DisableAutoRollback bool       `json:"disableAutoRollback" yal:"disableAutoRollback"`
}

type PoolExtra struct {
	NodePoolId string `json:"nodePoolId" yaml:"nodePoolId"`
}

type VimConInfo struct {
	ID      string    `json:"id" yaml:"id"`
	VimType string    `json:"vimType" yaml:"vimType"`
	Extra   PoolExtra `json:"extra" yaml:"extra"`
}

type InstantiateVnfRequest struct {
	FlavourID           string           `json:"flavourId" yaml:"flavourId"`
	VimConnectionInfo   []VimConInfo     `json:"vimConnectionInfo" yaml:"vimConnectionInfo"`
	AdditionalVduParams AdditionalParams `json:"additionalParams" yaml:"additionalParams"`
}
