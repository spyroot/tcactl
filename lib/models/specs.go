package models

import "time"

type HealthCheck struct {
	NodeStartupTimeout  string `json:"nodeStartupTimeout" yaml:"nodeStartupTimeout"`
	UnhealthyConditions []struct {
		Type    string `json:"type" yaml:"type"`
		Status  string `json:"status" yaml:"status"`
		Timeout string `json:"timeout" yaml:"timeout"`
	} `json:"unhealthyConditions," yaml:"unhealthy_conditions"`
}

type Nodes struct {
	Ip     string `json:"ip,omitempty" yaml:"ip,omitempty"`
	VmName string `json:"vmName,omitempty" yaml:"vmName,omitempty"`
}

type UserDefinedData struct {
	Name                   string            `json:"name" yaml:"name"`
	Tags                   []interface{}     `json:"tags" yaml:"tags"`
	NfType                 string            `json:"nfType" yaml:"nf_type"`
	ManagedBy              InternalManagedBy `json:"managedBy" yaml:"managed_by"`
	LocalFilePath          string            `json:"localFilePath" yaml:"localFilePath"`
	LastUpdated            time.Time         `json:"lastUpdated" yaml:"lastUpdated"`
	LastUpdateEnterprise   string            `json:"lastUpdateEnterprise" yaml:"lastUpdateEnterprise"`
	LastUpdateOrganization string            `json:"lastUpdateOrganization" yaml:"lastUpdateOrganization"`
	LastUpdateUser         string            `json:"lastUpdateUser" yaml:"lastUpdateUser"`
	CreationDate           time.Time         `json:"creationDate" yaml:"creationDate"`
	CreationEnterprise     string            `json:"creationEnterprise" yaml:"creationEnterprise"`
	CreationOrganization   string            `json:"creationOrganization" yaml:"creationOrganization"`
	CreationUser           string            `json:"creationUser" yaml:"creation_user"`
	IsDeleted              bool              `json:"isDeleted" yaml:"isDeleted"`
}

type InternalManagedBy struct {
	ExtensionSubtype string `json:"extensionSubtype" yaml:"extension_subtype"`
	ExtensionName    string `json:"extensionName" yaml:"extension_name"`
}

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
	HealthCheck      *HealthCheck         `json:"healthCheck,omitempty" yaml:"health_check,omitempty"`
}
