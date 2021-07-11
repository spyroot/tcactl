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

// VduParam
type VduParam struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	RepoURL   string `json:"repoUrl" yaml:"repoUrl"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	VduName   string `json:"vduName" yaml:"vduName"`
}

type AdditionalParams struct {
	VduParams           []VduParam `json:"vduParams,omitempty" yaml:"vduParams,omitempty"`
	DisableGrant        bool       `json:"disableGrant,omitempty" yaml:"disableGrant,omitempty"`
	IgnoreGrantFailure  bool       `json:"ignoreGrantFailure,omitempty" yaml:"ignoreGrantFailure,omitempty"`
	DisableAutoRollback bool       `json:"disableAutoRollback,omitempty" yaml:"disableAutoRollback,omitempty"`
}

type PoolExtra struct {
	NodePoolId string `json:"nodePoolId" yaml:"nodePoolId"`
}

type VimConInfo struct {
	ID      string    `json:"id" yaml:"id"`
	VimType string    `json:"vimType" yaml:"vimType"`
	Extra   PoolExtra `json:"extra" yaml:"extra"`
}
