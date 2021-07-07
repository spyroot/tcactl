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

import (
	"encoding/json"
	"github.com/spyroot/tcactl/lib/models"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// CpuManagerPolicy node pool attribute cpu manager
	CpuManagerPolicy = "kubernetes"

	// HealthyCheckReady healthy check ready flag
	HealthyCheckReady = "Ready"

	//HealthyCheckMemoryPressure memory pressure trigger
	HealthyCheckMemoryPressure = "MemoryPressure"

	//HealthyCheckDiskPressure disk pressure trigger
	HealthyCheckDiskPressure = "DiskPressure"

	//HealthyCheckPIDPressure pid trigger
	HealthyCheckPIDPressure = "PIDPressure"

	//HealthyCheckNetworkUnavailable network availability trigger
	HealthyCheckNetworkUnavailable = "NetworkUnavailable"

	//DefaultNodeStartupTimeout default start up timeout trigger
	DefaultNodeStartupTimeout = "20m"

	//FullCloneMode - full clone of VM
	FullCloneMode = "fullClone"

	// LinkedClone - link clone for VM
	LinkedClone = "linkedClone"
)

// NewNodePoolSpec - a request to create new node pool and attach to a target.
type NewNodePoolSpec struct {
	SpecType                      *SpecKind                `json:"kind,omitempty" yaml:"kind,omitempty" validate:"required"`
	Id                            string                   `json:"id,omitempty" yaml:"id,omitempty"`
	CloneMode                     string                   `json:"cloneMode,omitempty" yaml:"clone_mode,omitempty" validate:"required"`
	Cpu                           int                      `json:"cpu" yaml:"cpu" validate:"required"`
	Labels                        []string                 `json:"labels,omitempty" yaml:"labels" validate:"required"`
	Memory                        int                      `json:"memory" yaml:"memory" validate:"required"`
	Name                          string                   `json:"name" yaml:"name" validate:"required"`
	Networks                      []models.Network         `json:"networks" yaml:"networks" validate:"required"`
	PlacementParams               []models.PlacementParams `json:"placementParams" yaml:"placementParams" validate:"required"`
	Replica                       int                      `json:"replica" yaml:"replica" validate:"required"`
	Storage                       int                      `json:"storage" yaml:"storage" validate:"required"`
	Config                        *models.NodeConfig       `json:"config,omitempty" yaml:"config,omitempty"`
	Status                        string                   `json:"status" yaml:"status"`
	ActiveTasksCount              int                      `json:"activeTasksCount" yaml:"activeTasksCount"`
	Nodes                         []models.Nodes           `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	IsNodeCustomizationDeprecated bool                     `json:"isNodeCustomizationDeprecated" yaml:"isNodeCustomizationDeprecated"`
}

func (n *NewNodePoolSpec) validateCloneMode() bool {
	return strings.ToLower(n.CloneMode) == strings.ToLower(FullCloneMode) ||
		strings.ToLower(n.CloneMode) == strings.ToLower(LinkedClone)
}

//GetKind return spec kind
func (n *NewNodePoolSpec) GetKind() *SpecKind {
	return n.SpecType
}

// InvalidNodePoolSpec error if specs invalid
type InvalidNodePoolSpec struct {
	errMsg string
}

//InstanceSpecsFromString method return instance form string
func (n NewNodePoolSpec) InstanceSpecsFromString(s string) (interface{}, error) {
	return ReadNodeSpecFromString(s)
}

// NewNewNodePoolSpec create spec from reader
func NewNewNodePoolSpec(r io.Reader) (*NewNodePoolSpec, error) {
	spec, err := ReadNodeSpecSpec(r)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

//
func (m *InvalidNodePoolSpec) Error() string {
	return m.errMsg
}

// ReadNodeSpecFromFile - Read node template from file
func ReadNodeSpecFromFile(fileName string) (*NewNodePoolSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadNodeSpecSpec(file)
}

// ReadNodeSpecFromString - Read node template from string
func ReadNodeSpecFromString(str string) (*NewNodePoolSpec, error) {
	r := strings.NewReader(str)
	return ReadNodeSpecSpec(r)
}

// ReadNodeSpecSpec - Read node pool template specString
// either from yaml or json
func ReadNodeSpecSpec(b io.Reader) (*NewNodePoolSpec, error) {

	var spec NewNodePoolSpec

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	return nil, &InvalidNodePoolSpec{"Failed to parse input spec."}
}
