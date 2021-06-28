package request

import (
	"github.com/spyroot/tcactl/lib/models"
)

const (
	CpuManagerPolicy = "kubernetes"

	HealthyCheckReady              = "Ready"
	HealthyCheckMemoryPressure     = "MemoryPressure"
	HealthyCheckDiskPressure       = "DiskPressure"
	HealthyCheckPIDPressure        = "PIDPressure"
	HealthyCheckNetworkUnavailable = "NetworkUnavailable"
	DefaultNodeStartupTimeout      = "20m"

	FullCloneMode = "fullClone"
	LinkedClone   = "linkedClone"
)

// NewNodePoolSpec - a request to create new node pool and attach to a target.
type NewNodePoolSpec struct {
	Id                            string                   `json:"id,omitempty" yaml:"id,omitempty"`
	CloneMode                     string                   `json:"cloneMode,omitempty" yaml:"clone_mode,omitempty"`
	Cpu                           int                      `json:"cpu" yaml:"cpu"`
	Labels                        []string                 `json:"labels,omitempty" yaml:"labels"`
	Memory                        int                      `json:"memory" yaml:"memory"`
	Name                          string                   `json:"name" yaml:"name"`
	Networks                      []models.Network         `json:"networks" yaml:"networks"`
	PlacementParams               []models.PlacementParams `json:"placementParams" yaml:"placementParams"`
	Replica                       int                      `json:"replica" yaml:"replica"`
	Storage                       int                      `json:"storage" yaml:"storage"`
	Config                        *models.NodeConfig       `json:"config,omitempty" yaml:"config,omitempty"`
	Status                        string                   `json:"status" yaml:"status"`
	ActiveTasksCount              int                      `json:"activeTasksCount" yaml:"activeTasksCount"`
	Nodes                         []models.Nodes           `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	IsNodeCustomizationDeprecated bool                     `json:"isNodeCustomizationDeprecated" yaml:"isNodeCustomizationDeprecated"`
}

func (n *NewNodePoolSpec) validateCloneMode() bool {
	return n.CloneMode == FullCloneMode || n.CloneMode == LinkedClone
}
