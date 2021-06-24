package response

// LcmInfo information about current state
type LcmInfo struct {
	Id                        string `json:"id" yaml:"id"`
	VnfInstanceName           string `json:"vnfInstanceName" yaml:"vnf_instance_name"`
	VnfInstanceDescription    string `json:"vnfInstanceDescription" yaml:"vnf_instance_description"`
	VnfdId                    string `json:"vnfdId" yaml:"vnfd_id"`
	VnfProvider               string `json:"vnfProvider" yaml:"vnf_provider"`
	VnfProductName            string `json:"vnfProductName" yaml:"vnf_product_name"`
	VnfSoftwareVersion        string `json:"vnfSoftwareVersion" yaml:"vnf_software_version"`
	VnfdVersion               string `json:"vnfdVersion" yaml:"vnfd_version"`
	VnfConfigurableProperties struct {
	} `json:"vnfConfigurableProperties,omitempty" yaml:"vnf_configurable_properties"`
	VimConnectionInfo  []VimConnectionInfo `json:"vimConnectionInfo,omitempty" yaml:"vim_connection_info"`
	InstantiationState string              `json:"instantiationState" yaml:"instantiation_state"`
	VnfInfo            InstantiatedVnfInfo `json:"instantiatedVnfInfo,omitempty" yaml:"instantiatedVnfInfo"`
	Metadata           ExtendedMetadata    `json:"metadata" yaml:"metadata"`
	Extensions         struct {
	} `json:"extensions" yaml:"extensions"`
	Links PolicyLinks `json:"_links" yaml:"links"`
}

// Cnfs - list of CNF LCM respond
type Cnfs struct {
	CnfLcms []LcmInfo
}

// InstantiatedVnfInfo - Extended Info
type InstantiatedVnfInfo struct {
	FlavourId   string `json:"flavourId" yaml:"flavourId"`
	VnfState    string `json:"vnfState" yaml:"vnfState"`
	ScaleStatus []struct {
		AspectId   string `json:"aspectId" yaml:"aspect_id"`
		ScaleLevel int    `json:"scaleLevel" yaml:"scale_level"`
	} `json:"scaleStatus" yaml:"scaleStatus"`
	MaxScaleLevels []struct {
		AspectId   string `json:"aspectId" yaml:"aspectId"`
		ScaleLevel int    `json:"scaleLevel" yaml:"scaleLevel"`
	} `json:"maxScaleLevels" yaml:"maxScaleLevels"`
	ExtCpInfo []struct {
		Id             string `json:"id" yaml:"id"`
		CpdId          string `json:"cpdId" yaml:"cpdId"`
		CpProtocolInfo []struct {
			LayerProtocol  string `json:"layerProtocol" yaml:"layerProtocol"`
			IpOverEthernet struct {
				MacAddress  string `json:"macAddress" yaml:"mac_address"`
				IpAddresses []struct {
					Type         string   `json:"type" yaml:"type"`
					Addresses    []string `json:"addresses" yaml:"addresses"`
					IsDynamic    bool     `json:"isDynamic" yaml:"isDynamic"`
					AddressRange struct {
						MinAddress string `json:"minAddress" yaml:"minAddress"`
						MaxAddress string `json:"maxAddress" yaml:"maxAddress"`
					} `json:"addressRange" yaml:"addressRange"`
					SubnetId string `json:"subnetId" yaml:"subnetId"`
				} `json:"ipAddresses" yaml:"ipAddresses"`
			} `json:"ipOverEthernet" yaml:"ipOverEthernet"`
		} `json:"cpProtocolInfo" yaml:"cpProtocolInfo"`
		ExtLinkPortId string `json:"extLinkPortId" yaml:"ext_link_port_id"`
		Metadata      struct {
		} `json:"metadata" yaml:"metadata"`
		AssociatedVnfcCpId         string `json:"associatedVnfcCpId" yaml:"associated_vnfc_cp_id"`
		AssociatedVnfVirtualLinkId string `json:"associatedVnfVirtualLinkId" yaml:"associated_vnf_virtual_link_id"`
	} `json:"extCpInfo" yaml:"ext_cp_info"`
	ExtVirtualLinkInfo []struct {
		Id             string `json:"id" yaml:"id"`
		ResourceHandle struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
		} `json:"resourceHandle" yaml:"resource_handle"`
		ExtLinkPorts []struct {
			Id             string `json:"id" yaml:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
				ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
				ResourceId           string `json:"resourceId" yaml:"resource_id"`
				VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
			} `json:"resourceHandle" yaml:"resource_handle"`
			CpInstanceId string `json:"cpInstanceId" yaml:"cp_instance_id"`
		} `json:"extLinkPorts" yaml:"ext_link_ports"`
	} `json:"extVirtualLinkInfo" yaml:"ext_virtual_link_info"`
	ExtManagedVirtualLinkInfo []struct {
		Id                   string `json:"id" yaml:"id"`
		VnfVirtualLinkDescId string `json:"vnfVirtualLinkDescId" yaml:"vnf_virtual_link_desc_id"`
		NetworkResource      struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
		} `json:"networkResource" yaml:"network_resource"`
		VnfLinkPorts []struct {
			Id             string `json:"id" yaml:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
				ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
				ResourceId           string `json:"resourceId" yaml:"resource_id"`
				VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
			} `json:"resourceHandle" yaml:"resource_handle"`
			CpInstanceId   string `json:"cpInstanceId" yaml:"cp_instance_id"`
			CpInstanceType string `json:"cpInstanceType" yaml:"cp_instance_type"`
		} `json:"vnfLinkPorts" yaml:"vnf_link_ports"`
	} `json:"extManagedVirtualLinkInfo" yaml:"ext_managed_virtual_link_info"`
	MonitoringParameters []struct {
		Id                string `json:"id" yaml:"id"`
		Name              string `json:"name" yaml:"name"`
		PerformanceMetric string `json:"performanceMetric" yaml:"performance_metric"`
	} `json:"monitoringParameters" yaml:"monitoring_parameters"`
	LocalizationLanguage string `json:"localizationLanguage" yaml:"localization_language"`
	VnfcResourceInfo     []struct {
		Id              string `json:"id" yaml:"id"`
		VduId           string `json:"vduId" yaml:"vdu_id"`
		ComputeResource struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
		} `json:"computeResource" yaml:"compute_resource"`
		ZoneId             string   `json:"zoneId" yaml:"zone_id"`
		StorageResourceIds []string `json:"storageResourceIds" yaml:"storage_resource_ids"`
		ReservationId      string   `json:"reservationId" yaml:"reservation_id"`
		VnfcCpInfo         []struct {
			Id             string `json:"id" yaml:"id"`
			CpdId          string `json:"cpdId" yaml:"cpd_id"`
			VnfExtCpId     string `json:"vnfExtCpId" yaml:"vnf_ext_cp_id"`
			CpProtocolInfo []struct {
				LayerProtocol  string `json:"layerProtocol" yaml:"layer_protocol"`
				IpOverEthernet struct {
					MacAddress  string `json:"macAddress" yaml:"mac_address"`
					IpAddresses []struct {
						Type         string   `json:"type" yaml:"type"`
						Addresses    []string `json:"addresses" yaml:"addresses"`
						IsDynamic    bool     `json:"isDynamic" yaml:"is_dynamic"`
						AddressRange struct {
							MinAddress string `json:"minAddress" yaml:"min_address"`
							MaxAddress string `json:"maxAddress" yaml:"max_address"`
						} `json:"addressRange" yaml:"address_range"`
						SubnetId string `json:"subnetId" yaml:"subnet_id"`
					} `json:"ipAddresses" yaml:"ip_addresses"`
				} `json:"ipOverEthernet" yaml:"ip_over_ethernet"`
			} `json:"cpProtocolInfo" yaml:"cp_protocol_info"`
			VnfLinkPortId string `json:"vnfLinkPortId" yaml:"vnf_link_port_id"`
			Metadata      struct {
			} `json:"metadata" yaml:"metadata"`
		} `json:"vnfcCpInfo" yaml:"vnfc_cp_info"`
		Metadata struct {
		} `json:"metadata" yaml:"metadata"`
	} `json:"vnfcResourceInfo" yaml:"vnfc_resource_info"`
	VirtualLinkResourceInfo []struct {
		Id                   string `json:"id" yaml:"id"`
		VnfVirtualLinkDescId string `json:"vnfVirtualLinkDescId" yaml:"vnf_virtual_link_desc_id"`
		NetworkResource      struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
		} `json:"networkResource" yaml:"network_resource"`
		ZoneId        string `json:"zoneId" yaml:"zone_id"`
		ReservationId string `json:"reservationId" yaml:"reservation_id"`
		VnfLinkPorts  []struct {
			Id             string `json:"id" yaml:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
				ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
				ResourceId           string `json:"resourceId" yaml:"resource_id"`
				VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
			} `json:"resourceHandle" yaml:"resource_handle"`
			CpInstanceId   string `json:"cpInstanceId" yaml:"cp_instance_id"`
			CpInstanceType string `json:"cpInstanceType" yaml:"cp_instance_type"`
		} `json:"vnfLinkPorts" yaml:"vnf_link_ports"`
		Metadata struct {
		} `json:"metadata" yaml:"metadata"`
	} `json:"virtualLinkResourceInfo" yaml:"virtual_link_resource_info"`
	VirtualStorageResourceInfo []struct {
		Id                   string `json:"id" yaml:"id"`
		VirtualStorageDescId string `json:"virtualStorageDescId" yaml:"virtual_storage_desc_id"`
		StorageResource      struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
		} `json:"storageResource" yaml:"storage_resource"`
		ZoneId        string `json:"zoneId" yaml:"zone_id"`
		ReservationId string `json:"reservationId" yaml:"reservation_id"`
		Metadata      struct {
		} `json:"metadata" yaml:"metadata"`
	} `json:"virtualStorageResourceInfo" yaml:"virtual_storage_resource_info"`
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
