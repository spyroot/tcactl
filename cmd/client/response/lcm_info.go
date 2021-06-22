package response

// CnfExtended
type LcmInfo struct {
	Id                        string `json:"id"`
	VnfInstanceName           string `json:"vnfInstanceName"`
	VnfInstanceDescription    string `json:"vnfInstanceDescription"`
	VnfdId                    string `json:"vnfdId"`
	VnfProvider               string `json:"vnfProvider"`
	VnfProductName            string `json:"vnfProductName"`
	VnfSoftwareVersion        string `json:"vnfSoftwareVersion"`
	VnfdVersion               string `json:"vnfdVersion"`
	VnfConfigurableProperties struct {
	} `json:"vnfConfigurableProperties,omitempty"`
	VimConnectionInfo  []VimConnectionInfo `json:"vimConnectionInfo,omitempty"`
	InstantiationState string              `json:"instantiationState"`
	VnfInfo            InstantiatedVnfInfo `json:"instantiatedVnfInfo,omitempty"`
	Metadata           ExtendedMetadata    `json:"metadata"`
	Extensions         struct {
	} `json:"extensions"`
	Links PolicyLinks `json:"_links"`
}

// Cnfs - list of CNF LCM respond
type Cnfs struct {
	CnfLcms []LcmInfo
}

// InstantiatedVnfInfo - Extended Info
type InstantiatedVnfInfo struct {
	FlavourId   string `json:"flavourId"`
	VnfState    string `json:"vnfState"`
	ScaleStatus []struct {
		AspectId   string `json:"aspectId"`
		ScaleLevel int    `json:"scaleLevel"`
	} `json:"scaleStatus"`
	MaxScaleLevels []struct {
		AspectId   string `json:"aspectId"`
		ScaleLevel int    `json:"scaleLevel"`
	} `json:"maxScaleLevels"`
	ExtCpInfo []struct {
		Id             string `json:"id"`
		CpdId          string `json:"cpdId"`
		CpProtocolInfo []struct {
			LayerProtocol  string `json:"layerProtocol"`
			IpOverEthernet struct {
				MacAddress  string `json:"macAddress"`
				IpAddresses []struct {
					Type         string   `json:"type"`
					Addresses    []string `json:"addresses"`
					IsDynamic    bool     `json:"isDynamic"`
					AddressRange struct {
						MinAddress string `json:"minAddress"`
						MaxAddress string `json:"maxAddress"`
					} `json:"addressRange"`
					SubnetId string `json:"subnetId"`
				} `json:"ipAddresses"`
			} `json:"ipOverEthernet"`
		} `json:"cpProtocolInfo"`
		ExtLinkPortId string `json:"extLinkPortId"`
		Metadata      struct {
		} `json:"metadata"`
		AssociatedVnfcCpId         string `json:"associatedVnfcCpId"`
		AssociatedVnfVirtualLinkId string `json:"associatedVnfVirtualLinkId"`
	} `json:"extCpInfo"`
	ExtVirtualLinkInfo []struct {
		Id             string `json:"id"`
		ResourceHandle struct {
			VimConnectionId      string `json:"vimConnectionId"`
			ResourceProviderId   string `json:"resourceProviderId"`
			ResourceId           string `json:"resourceId"`
			VimLevelResourceType string `json:"vimLevelResourceType"`
		} `json:"resourceHandle"`
		ExtLinkPorts []struct {
			Id             string `json:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId"`
				ResourceProviderId   string `json:"resourceProviderId"`
				ResourceId           string `json:"resourceId"`
				VimLevelResourceType string `json:"vimLevelResourceType"`
			} `json:"resourceHandle"`
			CpInstanceId string `json:"cpInstanceId"`
		} `json:"extLinkPorts"`
	} `json:"extVirtualLinkInfo"`
	ExtManagedVirtualLinkInfo []struct {
		Id                   string `json:"id"`
		VnfVirtualLinkDescId string `json:"vnfVirtualLinkDescId"`
		NetworkResource      struct {
			VimConnectionId      string `json:"vimConnectionId"`
			ResourceProviderId   string `json:"resourceProviderId"`
			ResourceId           string `json:"resourceId"`
			VimLevelResourceType string `json:"vimLevelResourceType"`
		} `json:"networkResource"`
		VnfLinkPorts []struct {
			Id             string `json:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId"`
				ResourceProviderId   string `json:"resourceProviderId"`
				ResourceId           string `json:"resourceId"`
				VimLevelResourceType string `json:"vimLevelResourceType"`
			} `json:"resourceHandle"`
			CpInstanceId   string `json:"cpInstanceId"`
			CpInstanceType string `json:"cpInstanceType"`
		} `json:"vnfLinkPorts"`
	} `json:"extManagedVirtualLinkInfo"`
	MonitoringParameters []struct {
		Id                string `json:"id"`
		Name              string `json:"name"`
		PerformanceMetric string `json:"performanceMetric"`
	} `json:"monitoringParameters"`
	LocalizationLanguage string `json:"localizationLanguage"`
	VnfcResourceInfo     []struct {
		Id              string `json:"id"`
		VduId           string `json:"vduId"`
		ComputeResource struct {
			VimConnectionId      string `json:"vimConnectionId"`
			ResourceProviderId   string `json:"resourceProviderId"`
			ResourceId           string `json:"resourceId"`
			VimLevelResourceType string `json:"vimLevelResourceType"`
		} `json:"computeResource"`
		ZoneId             string   `json:"zoneId"`
		StorageResourceIds []string `json:"storageResourceIds"`
		ReservationId      string   `json:"reservationId"`
		VnfcCpInfo         []struct {
			Id             string `json:"id"`
			CpdId          string `json:"cpdId"`
			VnfExtCpId     string `json:"vnfExtCpId"`
			CpProtocolInfo []struct {
				LayerProtocol  string `json:"layerProtocol"`
				IpOverEthernet struct {
					MacAddress  string `json:"macAddress"`
					IpAddresses []struct {
						Type         string   `json:"type"`
						Addresses    []string `json:"addresses"`
						IsDynamic    bool     `json:"isDynamic"`
						AddressRange struct {
							MinAddress string `json:"minAddress"`
							MaxAddress string `json:"maxAddress"`
						} `json:"addressRange"`
						SubnetId string `json:"subnetId"`
					} `json:"ipAddresses"`
				} `json:"ipOverEthernet"`
			} `json:"cpProtocolInfo"`
			VnfLinkPortId string `json:"vnfLinkPortId"`
			Metadata      struct {
			} `json:"metadata"`
		} `json:"vnfcCpInfo"`
		Metadata struct {
		} `json:"metadata"`
	} `json:"vnfcResourceInfo"`
	VirtualLinkResourceInfo []struct {
		Id                   string `json:"id"`
		VnfVirtualLinkDescId string `json:"vnfVirtualLinkDescId"`
		NetworkResource      struct {
			VimConnectionId      string `json:"vimConnectionId"`
			ResourceProviderId   string `json:"resourceProviderId"`
			ResourceId           string `json:"resourceId"`
			VimLevelResourceType string `json:"vimLevelResourceType"`
		} `json:"networkResource"`
		ZoneId        string `json:"zoneId"`
		ReservationId string `json:"reservationId"`
		VnfLinkPorts  []struct {
			Id             string `json:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId"`
				ResourceProviderId   string `json:"resourceProviderId"`
				ResourceId           string `json:"resourceId"`
				VimLevelResourceType string `json:"vimLevelResourceType"`
			} `json:"resourceHandle"`
			CpInstanceId   string `json:"cpInstanceId"`
			CpInstanceType string `json:"cpInstanceType"`
		} `json:"vnfLinkPorts"`
		Metadata struct {
		} `json:"metadata"`
	} `json:"virtualLinkResourceInfo"`
	VirtualStorageResourceInfo []struct {
		Id                   string `json:"id"`
		VirtualStorageDescId string `json:"virtualStorageDescId"`
		StorageResource      struct {
			VimConnectionId      string `json:"vimConnectionId"`
			ResourceProviderId   string `json:"resourceProviderId"`
			ResourceId           string `json:"resourceId"`
			VimLevelResourceType string `json:"vimLevelResourceType"`
		} `json:"storageResource"`
		ZoneId        string `json:"zoneId"`
		ReservationId string `json:"reservationId"`
		Metadata      struct {
		} `json:"metadata"`
	} `json:"virtualStorageResourceInfo"`
}

// ExtendedMetadata
type ExtendedMetadata struct {
	VnfPkgId       string `json:"vnfPkgId"`
	VnfCatalogName string `json:"vnfCatalogName"`
	ManagedBy      struct {
		ExtensionSubtype string `json:"extensionSubtype"`
		ExtensionName    string `json:"extensionName"`
	} `json:"managedBy"`
	NfType            string `json:"nfType"`
	LcmOperation      string `json:"lcmOperation"`
	LcmOperationState string `json:"lcmOperationState"`
	IsUsedByNS        string `json:"isUsedByNS"`
	AttachedNSCount   string `json:"attachedNSCount"`
	ExtVirtualLinks   []struct {
		Id                 string `json:"id"`
		VimConnectionId    string `json:"vimConnectionId"`
		ResourceProviderId string `json:"resourceProviderId"`
		ResourceId         string `json:"resourceId"`
		ExtCps             []struct {
			CpdId    string `json:"cpdId"`
			CpConfig []struct {
				CpInstanceId   string `json:"cpInstanceId"`
				LinkPortId     string `json:"linkPortId"`
				CpProtocolData []struct {
					LayerProtocol  string `json:"layerProtocol"`
					IpOverEthernet struct {
						MacAddress  string `json:"macAddress"`
						IpAddresses []struct {
							Type                string   `json:"type"`
							FixedAddresses      []string `json:"fixedAddresses"`
							NumDynamicAddresses int      `json:"numDynamicAddresses"`
							AddressRange        struct {
								MinAddress string `json:"minAddress"`
								MaxAddress string `json:"maxAddress"`
							} `json:"addressRange"`
							SubnetId string `json:"subnetId"`
						} `json:"ipAddresses"`
					} `json:"ipOverEthernet"`
				} `json:"cpProtocolData"`
			} `json:"cpConfig"`
		} `json:"extCps"`
		ExtLinkPorts []struct {
			Id             string `json:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId"`
				ResourceProviderId   string `json:"resourceProviderId"`
				ResourceId           string `json:"resourceId"`
				VimLevelResourceType string `json:"vimLevelResourceType"`
			} `json:"resourceHandle"`
		} `json:"extLinkPorts"`
	} `json:"extVirtualLinks"`
	Tags []struct {
		Name        string `json:"name"`
		AutoCreated bool   `json:"autoCreated"`
	} `json:"tags"`
}
