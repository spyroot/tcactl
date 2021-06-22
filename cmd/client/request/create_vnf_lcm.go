package request

type CreateVnfLcm struct {
	VnfdId                 string `json:"vnfdId"`
	VnfInstanceName        string `json:"vnfInstanceName"`
	VnfInstanceDescription string `json:"vnfInstanceDescription"`
	Metadata               struct {
		Tags []interface{} `json:"tags"`
	} `json:"metadata"`
}
