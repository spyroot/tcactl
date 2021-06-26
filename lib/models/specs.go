package models

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
