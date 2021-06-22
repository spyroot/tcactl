package request

//
type TenantFilter struct {
	NfType string `json:"nfType,omitempty" yaml:"nfType"`
	NfdId  string `json:"nfdId,omitempty" yaml:"nfdId"`
}

// TenantsNfFilter filter
type TenantsNfFilter struct {
	Filter TenantFilter `json:"filter"`
}
