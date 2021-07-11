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

import (
	"github.com/asaskevich/govalidator"
	"io"
	"strings"
)

// InvalidTemplateSpec error if specs invalid
type InvalidTemplateSpec struct {
	errMsg string
}

//
func (m *InvalidTemplateSpec) Error() string {
	return m.errMsg
}

type CniSpec struct {
	Name       string `json:"name" yaml:"name" yaml:"name"`
	Properties struct {
	} `json:"properties" yaml:"name" yaml:"properties"`
}

type CsiSpec struct {
	Name       string `json:"name" yaml:"name" yaml:"name"`
	Properties struct {
		Name      string `json:"name" yaml:"name" yaml:"name"`
		IsDefault bool   `json:"isDefault" yaml:"name" yaml:"isDefault"`
		Timeout   string `json:"timeout" yaml:"name" yaml:"timeout"`
	} `json:"properties" yaml:"name" yaml:"properties"`
}

type ToolsSpec struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

//SpecCpuManagerPolicy cluster template parameter cpu manager policy
type SpecCpuManagerPolicy struct {
	Type       string `json:"type,omitempty" yaml:"type,omitempty"`
	Policy     string `json:"policy,omitempty" yaml:"policy,omitempty"`
	Properties struct {
		KubeReserved struct {
			Cpu         int `json:"cpu" yaml:"cpu"`
			MemoryInGiB int `json:"memoryInGiB" yaml:"memoryInGiB"`
		} `json:"kubeReserved" yaml:"kubeReserved"`
		SystemReserved struct {
			Cpu         int `json:"cpu" yaml:"cpu"`
			MemoryInGiB int `json:"memoryInGiB" yaml:"memoryInGiB"`
		} `json:"systemReserved" yaml:"systemReserved"`
	} `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// SpecHealthCheck cluster template parameter health check
type SpecHealthCheck struct {
	NodeStartupTimeout  string `json:"nodeStartupTimeout" yaml:"nodeStartupTimeout"`
	UnhealthyConditions []struct {
		Type    string `json:"type" yaml:"type"`
		Status  string `json:"status" yaml:"status"`
		Timeout string `json:"timeout" yaml:"timeout"`
	} `json:"unhealthyConditions" yaml:"unhealthyConditions"`
}

type SpecNodeConfig struct {
	KubernetesVersion string                `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	CpuManagerPolicy  *SpecCpuManagerPolicy `json:"cpuManagerPolicy,omitempty" yaml:"cpuManagerPolicy,omitempty"`
	HealthCheck       *SpecHealthCheck      `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
}

type SpecNodeTemplate struct {
	Cpu      int    `json:"cpu" yaml:"cpu"`
	Memory   int    `json:"memory" yaml:"memory"`
	Name     string `json:"name" yaml:"name"`
	Networks []struct {
		Label string `json:"label" yaml:"label"`
	} `json:"networks" yaml:"networks"`
	Storage   int             `json:"storage" yaml:"storage"`
	Replica   int             `json:"replica" yaml:"replica"`
	Labels    []string        `json:"labels" yaml:"labels"`
	CloneMode string          `json:"cloneMode" yaml:"cloneMode"`
	Config    *SpecNodeConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

type SpecClusterConfig struct {
	KubernetesVersion string      `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	Cni               []CniSpec   `json:"cni,omitempty" yaml:"cni,omitempty"`
	Csi               []CsiSpec   `json:"csi,omitempty" yaml:"csi,omitempty"`
	Tools             []ToolsSpec `json:"tools,omitempty" yaml:"tools,omitempty"`
}
type SpecClusterTemplate struct {
	SpecType          SpecType           `json:"kind,omitempty" yaml:"kind,omitempty" valid:"required"`
	Id                string             `json:"id,omitempty" yaml:"id,omitempty"`
	Name              string             `json:"name" yaml:"name" valid:"required"`
	ClusterType       string             `json:"clusterType" yaml:"clusterType" valid:"required"`
	KubernetesVersion string             `json:"kubernetesVersion,omitempty" yaml:"kubernetesVersion,omitempty"`
	ClusterConfig     *SpecNodeConfig    `json:"clusterConfig,omitempty" yaml:"clusterConfig,omitempty"`
	Description       string             `json:"description,omitempty" yaml:"description,omitempty"`
	MasterNodes       []SpecNodeTemplate `json:"masterNodes" yaml:"masterNodes"`
	WorkerNodes       []SpecNodeTemplate `json:"workerNodes" yaml:"workerNodes"`
	Tags              []struct {
		AutoCreated bool   `json:"autoCreated" yaml:"autoCreated"`
		Name        string `json:"name" yaml:"name"`
	} `json:"tags,omitempty" yaml:"tags,omitempty"`
	specError error
}

func (t *SpecClusterTemplate) Kind() SpecType {
	return t.SpecType
}

//SpecsFromString method read cluster spec from string
//and return instance
func (t SpecClusterTemplate) SpecsFromString(s string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromString(s, new(SpecClusterTemplate), f...)
}

//SpecsFromFile method return instance form string
func (t SpecClusterTemplate) SpecsFromFile(fileName string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromFile(fileName, new(SpecClusterTemplate), f...)
}

// SpecsFromReader create spec from reader
func (t SpecClusterTemplate) SpecsFromReader(r io.Reader, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpec(r, new(SpecClusterTemplate), f...)
}

// IsValid return false if validator set error
func (t *SpecClusterTemplate) IsValid() bool {
	if t.specError != nil {
		return false
	}
	return true
}

// Default set optional template values
//TODO
func (t *SpecClusterTemplate) Default() error {
	return nil
}

// validateCloneMode
func (t *SpecClusterTemplate) validateCloneMode() bool {
	for _, node := range t.MasterNodes {
		return strings.ToLower(node.CloneMode) != strings.ToLower(FullCloneMode) ||
			strings.ToLower(node.CloneMode) != strings.ToLower(LinkedClone)
	}

	for _, node := range t.WorkerNodes {
		return strings.ToLower(node.CloneMode) != strings.ToLower(FullCloneMode) ||
			strings.ToLower(node.CloneMode) != strings.ToLower(LinkedClone)
	}

	return true
}

//Validate method validate node pool specs
func (t *SpecClusterTemplate) Validate() error {

	if t == nil {
		return &InvalidTemplateSpec{errMsg: "nil instance"}
	}

	if t.Kind() != SpecKindTemplate {
		return &InvalidTemplateSpec{errMsg: "spec must contain kind field"}
	}

	if !t.validateCloneMode() {
		return &InvalidTemplateSpec{errMsg: "spec contains unknown clone mode type field"}
	}

	_, err := govalidator.ValidateStruct(t)
	if err != nil {
		return err
	}

	return nil
}
