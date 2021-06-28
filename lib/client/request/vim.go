package request

type RegisterVim struct {
	HcxCloudUrl string `json:"hcxCloudUrl"`
	VimName     string `json:"vimName"`
	TenantName  string `json:"tenantName"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}
