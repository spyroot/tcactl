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

// CreateVnfLcm Vnf Lcm Action

const (
	LcmTypeScaleOut = "SCALE_OUT"
	AspectId        = "aspect1"
)

type CreateVnfLcm struct {
	VnfdId                 string `json:"vnfdId" yaml:"vnfdId"`
	VnfInstanceName        string `json:"vnfInstanceName" yaml:"vnfInstanceName"`
	VnfInstanceDescription string `json:"vnfInstanceDescription" yaml:"vnfInstanceDescription"`
	Metadata               struct {
		Tags []interface{} `json:"tags" yaml:"tags"`
	} `json:"metadata" yaml:"metadata"`
}

type VduParams struct {
	VduName             string `json:"vduName" yaml:"vduName"`
	DeploymentProfileId string `json:"deploymentProfileId" yaml:"deploymentProfileId"`
	ChartName           string `json:"chartName" yaml:"chartName"`
	HelmName            string `json:"helmName" yaml:"helmName"`
	Namespace           string `json:"namespace" yaml:"namespace"`
	RepoUrl             string `json:"repoUrl" yaml:"repoUrl"`
	Username            string `json:"username" yaml:"username"`
	Password            string `json:"password" yaml:"password"`
	Overrides           string `json:"overrides" yaml:"overrides"`
	Metadata            []struct {
	} `json:"metadata" yaml:"metadata"`
	ImageName           string `json:"imageName" yaml:"imageName"`
	DisableAutoRollback bool   `json:"disableAutoRollback" yaml:"disableAutoRollback"`
	DisableGrant        bool   `json:"disableGrant" yaml:"disableGrant"`
	IgnoreGrantFailure  bool   `json:"ignoreGrantFailure" yaml:"ignoreGrantFailure"`
	CatalogName         string `json:"catalogName" yaml:"catalogName"`
	CatalogId           string `json:"catalogId" yaml:"catalogId"`
}

type CnfReconfigure struct {
	Type             string `json:"type" yaml:"type"`
	AspectId         string `json:"aspectId" yaml:"aspectId"`
	NumberOfSteps    int    `json:"numberOfSteps" yaml:"numberOfSteps"`
	AdditionalParams struct {
		SkipGrant     bool        `json:"skipGrant" yaml:"skipGrant"`
		VduParams     []VduParams `json:"vduParams" yaml:"vduParams"`
		LcmInterfaces []struct {
			InterfaceName string `json:"interfaceName" yaml:"interfaceName"`
			Parameters    []struct {
				Name string `json:"name" yaml:"name"`
				Type string `json:"type" yaml:"type"`
			} `json:"parameters" yaml:"parameters"`
		} `json:"lcmInterfaces" yaml:"lcm_interfaces"`
	} `json:"additionalParams" yaml:"additional_params"`
}
