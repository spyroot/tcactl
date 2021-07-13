package models

import "time"

type RegistrationRespond struct {
	TenantId    string        `json:"tenantId"`
	VimName     string        `json:"vimName"`
	TenantName  string        `json:"tenantName"`
	HcxCloudUrl string        `json:"hcxCloudUrl"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	Tags        []interface{} `json:"tags"`
	VimType     string        `json:"vimType"`
	VimUrl      string        `json:"vimUrl"`
	HcxUUID     string        `json:"hcxUUID"`
	HcxTenantId string        `json:"hcxTenantId"`
	Location    struct {
		City      string  `json:"city"`
		Country   string  `json:"country"`
		CityAscii string  `json:"cityAscii"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	VimId string `json:"vimId"`
	Audit struct {
		CreationUser      string    `json:"creationUser"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"audit"`
}

func (r *RegistrationRespond) ProviderType() string {
	if r.VimType == "VC" {
		return "VMware"
	}

	return ""
}
