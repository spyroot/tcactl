package request

import (
	"github.com/spyroot/tcactl/lib/models"
)

const (
	//
	CpuManagerPolicy = "kubernetes"

	//
	HealthyCheckReady = "Ready"
	//
	HealthyCheckMemoryPressure = "MemoryPressure"
	//
	HealthyCheckDiskPressure = "DiskPressure"
	//
	HealthyCheckPIDPressure = "PIDPressure"
	//
	HealthyCheckNetworkUnavailable = "NetworkUnavailable"
	//
	DefaultNodeStartupTimeout = "20m"

	FullCloneMode = "fullClone"
	LinkedClone   = "linkedClone"
)

// NewNodePoolSpec - a request to create new node pool and attach to a target.
type NewNodePoolSpec struct {
	Id                            string                   `json:"id,omitempty" yaml:"id,omitempty" validate:"required"`
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
	Status                        string                   `json:"status" yaml:"status" validate:"required"`
	ActiveTasksCount              int                      `json:"activeTasksCount" yaml:"activeTasksCount"`
	Nodes                         []models.Nodes           `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	IsNodeCustomizationDeprecated bool                     `json:"isNodeCustomizationDeprecated" yaml:"isNodeCustomizationDeprecated"`
}

func (n *NewNodePoolSpec) validateCloneMode() bool {
	return n.CloneMode == FullCloneMode || n.CloneMode == LinkedClone
}
