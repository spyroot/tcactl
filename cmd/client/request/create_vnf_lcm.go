package request

// CreateVnfLcm Vnf Lcm Action
type CreateVnfLcm struct {
	VnfdId                 string `json:"vnfdId" yaml:"vnfd_id"`
	VnfInstanceName        string `json:"vnfInstanceName" yaml:"vnf_instance_name"`
	VnfInstanceDescription string `json:"vnfInstanceDescription" yaml:"vnf_instance_description"`
	Metadata               struct {
		Tags []interface{} `json:"tags" yaml:"tags"`
	} `json:"metadata" yaml:"metadata"`
}
