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

type NodeProperties struct {
	KubeReserved struct {
		Cpu         int `json:"cpu,omitempty" yaml:"cpu,omitempty"`
		MemoryInGiB int `json:"memoryInGiB,omitempty" yaml:"memoryInGiB,omitempty"`
	} `json:"kubeReserved,omitempty" yaml:"kubeReserved,omitempty"`
	SystemReserved struct {
		Cpu         int `json:"cpu,omitempty" yaml:"cpu,omitempty"`
		MemoryInGiB int `json:"memoryInGiB" yaml:"memoryInGiB"`
	} `json:"systemReserved,omitempty" yaml:"system_reserved,omitempty"`
}

type K8sCpuManagerPolicy struct {
	Type       string          `json:"type,omitempty" yaml:"type,omitempty"`
	Policy     string          `json:"policy,omitempty" yaml:"policy,omitempty"`
	Properties *NodeProperties `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type NodeConfig struct {
	CpuManagerPolicy *K8sCpuManagerPolicy `json:"cpuManagerPolicy,omitempty" yaml:"cpu_manager_policy,omitempty"`
	HealthCheck      *models.HealthCheck  `json:"healthCheck,omitempty" yaml:"health_check,omitempty"`
}

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
	Config                        *NodeConfig              `json:"config,omitempty" yaml:"config,omitempty"`
	Status                        string                   `json:"status" yaml:"status"`
	ActiveTasksCount              int                      `json:"activeTasksCount" yaml:"active_tasks_count"`
	Nodes                         []models.Nodes           `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	IsNodeCustomizationDeprecated bool                     `json:"isNodeCustomizationDeprecated" yaml:"isNodeCustomizationDeprecated"`
}
