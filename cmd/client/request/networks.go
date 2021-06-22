package request

//{"filter":{"vimId":"vmware_8BF6253CE6E247018D909605A437B827","nodePoolId":"7f6f853c-7a87-4122-81df-f7a65186674c","type":["DistributedVirtualPortgroup","OpaqueNetwork"]}}

type NetworkFilter struct {
	Filter struct {
		VimId      string   `json:"vimId"`
		NodePoolId string   `json:"nodePoolId"`
		Type       []string `json:"type"`
	} `json:"filter"`
}
