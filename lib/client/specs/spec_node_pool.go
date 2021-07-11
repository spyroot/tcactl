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
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spyroot/tcactl/lib/models"
	"io"
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

// InvalidPoolSpec error if specs invalid
type InvalidPoolSpec struct {
	errMsg string
}

//
func (m *InvalidPoolSpec) Error() string {
	return m.errMsg
}

// SpecNodePool - a request to create new node pool and attach to a target.
type SpecNodePool struct {
	SpecType                      SpecType                 `json:"kind,omitempty" yaml:"kind,omitempty" valid:"required~kind is mandatory spec field"`
	Id                            string                   `json:"id,omitempty" yaml:"id,omitempty"`
	Name                          string                   `json:"name" yaml:"name" validate:"required" valid:"required~name is mandatory spec field"`
	CloneMode                     string                   `json:"cloneMode,omitempty" yaml:"clone_mode,omitempty" valid:"required~clone_mode is mandatory spec field"`
	Cpu                           int                      `json:"cpu" yaml:"cpu" validate:"required" valid:"required~cpu is mandatory spec field"`
	Memory                        int                      `json:"memory" yaml:"memory" validate:"required" valid:"required~memory is mandatory spec field"`
	Replica                       int                      `json:"replica" yaml:"replica" validate:"required" valid:"required~replica is mandatory spec field"`
	Storage                       int                      `json:"storage" yaml:"storage" validate:"required" valid:"required~storage is mandatory spec field"`
	Labels                        []string                 `json:"labels,omitempty" yaml:"labels" validate:"required" valid:"required~labels is mandatory spec field"`
	Networks                      []models.Network         `json:"networks" yaml:"networks" validate:"required" valid:"required~networks is mandatory spec field"`
	PlacementParams               []models.PlacementParams `json:"placementParams" yaml:"placementParams" validate:"required" valid:"required~placementParams is mandatory spec field"`
	Config                        *models.NodeConfig       `json:"config,omitempty" yaml:"config,omitempty"`
	Status                        string                   `json:"status" yaml:"status"`
	Nodes                         []models.Nodes           `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	ActiveTasksCount              int                      `json:"activeTasksCount" yaml:"activeTasksCount"`
	IsNodeCustomizationDeprecated bool                     `json:"isNodeCustomizationDeprecated" yaml:"isNodeCustomizationDeprecated"`

	// specError hold spec validator error
	specError error
}

func (t *SpecNodePool) validateCloneMode() bool {
	return strings.ToLower(t.CloneMode) == strings.ToLower(FullCloneMode) ||
		strings.ToLower(t.CloneMode) == strings.ToLower(LinkedClone)
}

//GetKind return spec kind
func (t *SpecNodePool) GetKind() SpecType {
	return t.SpecType
}

// InvalidNodePoolSpec error if specs invalid
type InvalidNodePoolSpec struct {
	errMsg string
}

//
func (m *InvalidNodePoolSpec) Error() string {
	return m.errMsg
}

func (t *SpecNodePool) Kind() SpecType {
	return t.SpecType
}

//Default  sets all all optional parameter to default value
func (t *SpecNodePool) Default() error {

	// set default clone mode
	t.CloneMode = LinkedClone
	t.IsNodeCustomizationDeprecated = false
	t.Replica = 1

	return nil
}

//SpecsFromString method read cluster spec from string
//and return instance
func (t SpecNodePool) SpecsFromString(s string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromString(s, new(SpecNodePool), f...)
}

//SpecsFromFile method reads instance form string
func (t SpecNodePool) SpecsFromFile(fileName string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromFile(fileName, new(SpecNodePool), f...)
}

// SpecsFromReader reads node pool spec from io.Reader
func (t SpecNodePool) SpecsFromReader(r io.Reader, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpec(r, new(SpecNodePool), f...)
}

// IsValid return false if validator set error
func (t *SpecNodePool) IsValid() bool {
	if t.specError != nil {
		return false
	}
	return true
}

//Validate method validate node pool specs
func (t *SpecNodePool) Validate() error {

	if t == nil {
		return &InvalidPoolSpec{errMsg: "nil instance"}
	}

	if t.Kind() != SpecKindNodePool {
		return &InvalidPoolSpec{errMsg: fmt.Sprintf(
			"Invalid spec kind. Node pool must use kind %s", SpecKindNodePool)}
	}

	if !t.validateCloneMode() {
		return &InvalidPoolSpec{errMsg: fmt.Sprintf("invalid clone mode supported types %s, %s", FullCloneMode, LinkedClone)}
	}
	if len(t.Networks) == 0 {
		return &InvalidPoolSpec{errMsg: "node pool spec must contain networks spec section."}
	}
	if len(t.Labels) == 0 {
		return &InvalidPoolSpec{errMsg: "node pool spec must contain labels spec section."}
	}
	if len(t.PlacementParams) == 0 {
		return &InvalidPoolSpec{errMsg: "node pool spec must contain labels spec section."}
	}
	if t.Cpu == 0 {
		return &InvalidPoolSpec{errMsg: "node pool spec cpu value is zero."}
	}
	if t.Memory == 0 {
		return &InvalidPoolSpec{errMsg: "node pool spec memory value is zero."}
	}
	if t.Replica == 0 {
		return &InvalidPoolSpec{errMsg: "node pool spec replica value is zero."}
	}
	if t.Storage == 0 {
		return &InvalidPoolSpec{errMsg: "node pool spec storage value is zero."}
	}
	_, err := govalidator.ValidateStruct(t)
	if err != nil {
		return err
	}

	return nil
}
