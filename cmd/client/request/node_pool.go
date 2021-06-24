package request

const (
	CpuManagerPolicy = "kubernetes"

	HealthyCheckReady              = "Ready"
	HealthyCheckMemoryPressure     = "MemoryPressure"
	HealthyCheckDiskPressure       = "DiskPressure"
	HealthyCheckPIDPressure        = "PIDPressure"
	HealthyCheckNetworkUnavailable = "NetworkUnavailable"
	DefaultNodeStartupTimeout      = "20m"
)

type NewNodePool struct {
	Name     string   `json:"name" yaml:"name"`
	Storage  int      `json:"storage" yaml:"storage"`
	Cpu      int      `json:"cpu" yaml:"cpu"`
	Memory   int      `json:"memory" yaml:"memory"`
	Replica  int      `json:"replica" yaml:"replica"`
	Labels   []string `json:"labels" yaml:"labels"`
	Networks []struct {
		Label       string   `json:"label" yaml:"label"`
		NetworkName string   `json:"networkName" yaml:"networkName"`
		Nameservers []string `json:"nameservers" yaml:"nameservers"`
	} `json:"networks" yaml:"networks"`
	PlacementParams []struct {
		Type string `json:"type" yaml:"type"`
		Name string `json:"name" yaml:"name"`
	} `json:"placementParams" yaml:"placement_params"`
	Config struct {
		CpuManagerPolicy struct {
			Type   string `json:"type" yaml:"type"`
			Policy string `json:"policy" yaml:"policy"`
		} `json:"cpuManagerPolicy" yaml:"cpuManagerPolicy"`
		HealthCheck struct {
			NodeStartupTimeout  string `json:"nodeStartupTimeout" yaml:"nodeStartupTimeout"`
			UnhealthyConditions []struct {
				Type    string `json:"type" yaml:"type"`
				Status  string `json:"status" yaml:"status"`
				Timeout string `json:"timeout" yaml:"timeout"`
			} `json:"unhealthyConditions" yaml:"unhealthyConditions"`
		} `json:"healthCheck" yaml:"healthCheck"`
	} `json:"config" yaml:"config"`
}
