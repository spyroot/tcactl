package response

type Vims struct {
	Items []struct {
		TenantId    string `json:"tenantId"`
		VimName     string `json:"vimName"`
		TenantName  string `json:"tenantName"`
		HcxCloudUrl string `json:"hcxCloudUrl"`
		Username    string `json:"username"`
		Tags        []struct {
			Name struct {
			} `json:"name"`
			AutoCreated bool `json:"autoCreated"`
		} `json:"tags"`
		VimType  string `json:"vimType"`
		VimUrl   string `json:"vimUrl"`
		Location struct {
			City      string `json:"city"`
			Country   string `json:"country"`
			CityAscii string `json:"cityAscii"`
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"location"`
		VimId string `json:"vimId"`
		Audit struct {
			CreationUser      string `json:"creationUser"`
			CreationTimestamp string `json:"creationTimestamp"`
		} `json:"audit"`
		Connection struct {
			Status       string `json:"status"`
			RemoteStatus string `json:"remoteStatus"`
		} `json:"connection"`
		Compatible         bool `json:"compatible"`
		VnfInstancesCounts []struct {
			TotalSize int    `json:"totalSize"`
			State     string `json:"state"`
		} `json:"vnfInstancesCounts"`
		Stats struct {
			Cpu struct {
				Used      int    `json:"used"`
				Available int    `json:"available"`
				Unit      string `json:"unit"`
				Capacity  int    `json:"capacity"`
			} `json:"cpu"`
			Memory struct {
				Used      int    `json:"used"`
				Available int    `json:"available"`
				Unit      string `json:"unit"`
				Capacity  int    `json:"capacity"`
			} `json:"memory"`
			Storage struct {
				Used      int    `json:"used"`
				Available int    `json:"available"`
				Unit      string `json:"unit"`
				Capacity  int    `json:"capacity"`
			} `json:"storage"`
		} `json:"stats"`
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
}
