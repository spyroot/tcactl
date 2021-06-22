package response

type VNFInstantiate struct {
	Id                     string `json:"id"`
	VnfInstanceName        string `json:"vnfInstanceName"`
	VnfInstanceDescription string `json:"vnfInstanceDescription"`
	VnfdId                 string `json:"vnfdId"`
	VnfProvider            string `json:"vnfProvider"`
	VnfProductName         string `json:"vnfProductName"`
	VnfSoftwareVersion     string `json:"vnfSoftwareVersion"`
	VnfdVersion            string `json:"vnfdVersion"`
	InstantiationState     string `json:"instantiationState"`
	Metadata               struct {
		VnfPkgId       string `json:"vnfPkgId"`
		VnfCatalogName string `json:"vnfCatalogName"`
		ManagedBy      struct {
			ExtensionSubtype string `json:"extensionSubtype"`
			ExtensionName    string `json:"extensionName"`
		} `json:"managedBy,omitempty"`
		NfType            string `json:"nfType"`
		LcmOperation      string `json:"lcmOperation"`
		LcmOperationState string `json:"lcmOperationState"`
		IsUsedByNS        string `json:"isUsedByNS"`
	} `json:"metadata,omitempty"`
}
