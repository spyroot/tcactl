package response

type VNFInstantiate struct {
	Id                     string `json:"id" yaml:"id"`
	VnfInstanceName        string `json:"vnfInstanceName" yaml:"vnfInstanceName"`
	VnfInstanceDescription string `json:"vnfInstanceDescription" yaml:"vnfInstanceDescription"`
	VnfdId                 string `json:"vnfdId" yaml:"vnfdId"`
	VnfProvider            string `json:"vnfProvider" yaml:"vnfProvider"`
	VnfProductName         string `json:"vnfProductName" yaml:"vnfProductName"`
	VnfSoftwareVersion     string `json:"vnfSoftwareVersion" yaml:"vnfSoftwareVersion"`
	VnfdVersion            string `json:"vnfdVersion" yaml:"vnfdVersion"`
	InstantiationState     string `json:"instantiationState" yaml:"instantiationState"`
	Metadata               struct {
		VnfPkgId       string `json:"vnfPkgId" yaml:"vnf_pkg_id"`
		VnfCatalogName string `json:"vnfCatalogName" yaml:"vnfCatalogName"`
		ManagedBy      struct {
			ExtensionSubtype string `json:"extensionSubtype" yaml:"extensionSubtype"`
			ExtensionName    string `json:"extensionName" yaml:"extensionName"`
		} `json:"managedBy,omitempty" yaml:"managed_by"`
		NfType            string `json:"nfType" yaml:"nf_type"`
		LcmOperation      string `json:"lcmOperation" yaml:"lcm_operation"`
		LcmOperationState string `json:"lcmOperationState" yaml:"lcm_operation_state"`
		IsUsedByNS        string `json:"isUsedByNS" yaml:"is_used_by_ns"`
	} `json:"metadata,omitempty" yaml:"metadata"`
}
