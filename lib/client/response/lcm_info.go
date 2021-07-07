package response

import (
	"github.com/spyroot/tcactl/lib/models"
	"time"
)

type InstanceUpdate struct {
	Id                    string `json:"id" yaml:"id"`
	TaskId                string `json:"taskId" yaml:"taskId"`
	OperationState        string `json:"operationState" yaml:"operationState"`
	StateEnteredTime      int64  `json:"stateEnteredTime" yaml:"stateEnteredTime"`
	StartTime             int64  `json:"startTime" yaml:"startTime"`
	InstanceId            string `json:"instanceId" yaml:"instanceId"`
	Operation             string `json:"operation" yaml:"operation"`
	IsAutomaticInvocation bool   `json:"isAutomaticInvocation" yaml:"isAutomaticInvocation"`
	OperationParams       struct {
		FlavourId        string `json:"flavourId" yaml:"flavourId"`
		AdditionalParams struct {
			VimId               string `json:"vimId" yaml:"vimId"`
			NodePoolId          string `json:"nodePoolId" yaml:"nodePoolId"`
			SkipGrant           bool   `json:"skipGrant" yaml:"skipGrant"`
			IgnoreGrantFailure  bool   `json:"ignoreGrantFailure" yaml:"ignoreGrantFailure"`
			DisableAutoRollback bool   `json:"disableAutoRollback" yaml:"disableAutoRollback"`
			DisableGrant        bool   `json:"disableGrant" yaml:"disableGrant"`
			UseVAppTemplates    bool   `json:"useVAppTemplates" yaml:"useVAppTemplates"`
		} `json:"additionalParams" yaml:"additional_params"`
		VimId        string `json:"vimId" yaml:"vimId"`
		NfInstanceId string `json:"nfInstanceId" yaml:"nfInstanceId"`
		Id           string `json:"id" yaml:"id"`
	} `json:"operationParams" yaml:"operationParams"`
	IsCancelPending        bool      `json:"isCancelPending" yaml:"is_cancel_pending"`
	EntityName             string    `json:"entityName" yaml:"entityName"`
	EntityType             string    `json:"entityType" yaml:"entityType"`
	LastUpdated            time.Time `json:"lastUpdated" yaml:"lastUpdated"`
	LastUpdateEnterprise   string    `json:"lastUpdateEnterprise" yaml:"lastUpdateEnterprise"`
	LastUpdateOrganization string    `json:"lastUpdateOrganization" yaml:"lastUpdateOrganization"`
	LastUpdateUser         string    `json:"lastUpdateUser" yaml:"last_update_user"`
	CreationDate           time.Time `json:"creationDate" yaml:"creationDate"`
	CreationEnterprise     string    `json:"creationEnterprise" yaml:"creationEnterprise"`
	CreationOrganization   string    `json:"creationOrganization" yaml:"creationOrganization"`
	CreationUser           string    `json:"creationUser" yaml:"creation_user"`
	IsDeleted              bool      `json:"isDeleted" yaml:"is_deleted"`
	EndTime                int64     `json:"endTime" yaml:"end_time"`
	Error                  string    `json:"error" yaml:"error"`
}

// LcmInfo information about current state
type LcmInfo struct {
	Id                        string `json:"id" yaml:"id"`
	VnfInstanceName           string `json:"vnfInstanceName" yaml:"vnf_instance_name"`
	VnfInstanceDescription    string `json:"vnfInstanceDescription" yaml:"vnfInstanceDescription"`
	VnfdId                    string `json:"vnfdId" yaml:"vnfdId"`
	VnfProvider               string `json:"vnfProvider" yaml:"vnfProvider"`
	VnfProductName            string `json:"vnfProductName" yaml:"vnfProductName"`
	VnfSoftwareVersion        string `json:"vnfSoftwareVersion" yaml:"vnfSoftwareVersion"`
	VnfdVersion               string `json:"vnfdVersion" yaml:"vnfdVersion"`
	VnfConfigurableProperties struct {
	} `json:"vnfConfigurableProperties,omitempty" yaml:"vnfConfigurableProperties"`
	VimConnectionInfo  []models.VimConnectionInfo `json:"vimConnectionInfo,omitempty" yaml:"vimConnectionInfo"`
	InstantiationState string                     `json:"instantiationState" yaml:"instantiationState"`
	VnfInfo            *InstantiatedVnfInfo       `json:"instantiatedVnfInfo,omitempty" yaml:"instantiatedVnfInfo"`
	Metadata           *ExtendedMetadata          `json:"metadata" yaml:"metadata"`
	Extensions         struct {
	} `json:"extensions" yaml:"extensions"`
	Links models.PolicyLinks `json:"_links" yaml:"_links"`
}

func (i *LcmInfo) IsInstantiated() bool {
	return i.InstantiationState == "INSTANTIATED"
}

func (i *LcmInfo) IsFailed() bool {
	return i.Metadata.LcmOperationState == "FAILED_TEMP"
}

// Cnfs - list of CNF LCM respond
type Cnfs struct {
	CnfLcms []LcmInfo
}

// InstantiatedVnfInfo - Extended Info
type InstantiatedVnfInfo struct {
	FlavourId   string `json:"flavourId" yaml:"flavourId" yaml:"flavour_id"`
	VnfState    string `json:"vnfState" yaml:"vnfState" yaml:"vnf_state"`
	ScaleStatus []struct {
		AspectId   string `json:"aspectId" yaml:"aspect_id" yaml:"aspect_id"`
		ScaleLevel int    `json:"scaleLevel" yaml:"scale_level" yaml:"scale_level"`
	} `json:"scaleStatus" yaml:"scaleStatus" yaml:"scale_status"`
	MaxScaleLevels []struct {
		AspectId   string `json:"aspectId" yaml:"aspectId" yaml:"aspect_id"`
		ScaleLevel int    `json:"scaleLevel" yaml:"scaleLevel" yaml:"scale_level"`
	} `json:"maxScaleLevels" yaml:"maxScaleLevels" yaml:"max_scale_levels"`
	ExtCpInfo []struct {
		Id             string `json:"id" yaml:"id" yaml:"id"`
		CpdId          string `json:"cpdId" yaml:"cpdId" yaml:"cpd_id"`
		CpProtocolInfo []struct {
			LayerProtocol  string `json:"layerProtocol" yaml:"layerProtocol" yaml:"layer_protocol"`
			IpOverEthernet struct {
				MacAddress  string `json:"macAddress" yaml:"mac_address" yaml:"mac_address"`
				IpAddresses []struct {
					Type         string   `json:"type" yaml:"type" yaml:"type"`
					Addresses    []string `json:"addresses" yaml:"addresses" yaml:"addresses"`
					IsDynamic    bool     `json:"isDynamic" yaml:"isDynamic" yaml:"is_dynamic"`
					AddressRange struct {
						MinAddress string `json:"minAddress" yaml:"minAddress" yaml:"min_address"`
						MaxAddress string `json:"maxAddress" yaml:"maxAddress" yaml:"max_address"`
					} `json:"addressRange" yaml:"addressRange" yaml:"address_range"`
					SubnetId string `json:"subnetId" yaml:"subnetId" yaml:"subnet_id"`
				} `json:"ipAddresses" yaml:"ipAddresses" yaml:"ip_addresses"`
			} `json:"ipOverEthernet" yaml:"ipOverEthernet" yaml:"ip_over_ethernet"`
		} `json:"cpProtocolInfo" yaml:"cpProtocolInfo" yaml:"cp_protocol_info"`
		ExtLinkPortId string `json:"extLinkPortId" yaml:"extLinkPortId" yaml:"ext_link_port_id"`
		Metadata      struct {
		} `json:"metadata" yaml:"metadata" yaml:"metadata"`
		AssociatedVnfcCpId         string `json:"associatedVnfcCpId" yaml:"associatedVnfcCpId" yaml:"associated_vnfc_cp_id"`
		AssociatedVnfVirtualLinkId string `json:"associatedVnfVirtualLinkId" yaml:"associatedVnfVirtualLinkId" yaml:"associated_vnf_virtual_link_id"`
	} `json:"extCpInfo" yaml:"ext_cp_info" yaml:"ext_cp_info"`
	ExtVirtualLinkInfo []struct {
		Id             string `json:"id" yaml:"id" yaml:"id"`
		ResourceHandle struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type" yaml:"vim_level_resource_type"`
		} `json:"resourceHandle" yaml:"resource_handle" yaml:"resource_handle"`
		ExtLinkPorts []struct {
			Id             string `json:"id" yaml:"id" yaml:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId" yaml:"vimConnectionId" yaml:"vim_connection_id"`
				ResourceProviderId   string `json:"resourceProviderId" yaml:"resourceProviderId" yaml:"resource_provider_id"`
				ResourceId           string `json:"resourceId" yaml:"resourceId" yaml:"resource_id"`
				VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vimLevelResourceType" yaml:"vim_level_resource_type"`
			} `json:"resourceHandle" yaml:"resource_handle" yaml:"resource_handle"`
			CpInstanceId string `json:"cpInstanceId" yaml:"cp_instance_id" yaml:"cp_instance_id"`
		} `json:"extLinkPorts" yaml:"ext_link_ports" yaml:"ext_link_ports"`
	} `json:"extVirtualLinkInfo" yaml:"ext_virtual_link_info" yaml:"ext_virtual_link_info"`
	ExtManagedVirtualLinkInfo []struct {
		Id                   string `json:"id" yaml:"id" yaml:"id"`
		VnfVirtualLinkDescId string `json:"vnfVirtualLinkDescId" yaml:"vnf_virtual_link_desc_id" yaml:"vnf_virtual_link_desc_id"`
		NetworkResource      struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type" yaml:"vim_level_resource_type"`
		} `json:"networkResource" yaml:"network_resource" yaml:"network_resource"`
		VnfLinkPorts []struct {
			Id             string `json:"id" yaml:"id" yaml:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id" yaml:"vim_connection_id"`
				ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id" yaml:"resource_provider_id"`
				ResourceId           string `json:"resourceId" yaml:"resource_id" yaml:"resource_id"`
				VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type" yaml:"vim_level_resource_type"`
			} `json:"resourceHandle" yaml:"resource_handle" yaml:"resource_handle"`
			CpInstanceId   string `json:"cpInstanceId" yaml:"cp_instance_id" yaml:"cp_instance_id"`
			CpInstanceType string `json:"cpInstanceType" yaml:"cp_instance_type" yaml:"cp_instance_type"`
		} `json:"vnfLinkPorts" yaml:"vnf_link_ports" yaml:"vnf_link_ports"`
	} `json:"extManagedVirtualLinkInfo" yaml:"ext_managed_virtual_link_info" yaml:"ext_managed_virtual_link_info"`
	MonitoringParameters []struct {
		Id                string `json:"id" yaml:"id" yaml:"id"`
		Name              string `json:"name" yaml:"name" yaml:"name"`
		PerformanceMetric string `json:"performanceMetric" yaml:"performance_metric" yaml:"performance_metric"`
	} `json:"monitoringParameters" yaml:"monitoring_parameters" yaml:"monitoring_parameters"`
	LocalizationLanguage string `json:"localizationLanguage" yaml:"localization_language" yaml:"localization_language"`
	VnfcResourceInfo     []struct {
		Id              string `json:"id" yaml:"id" yaml:"id"`
		VduId           string `json:"vduId" yaml:"vdu_id" yaml:"vdu_id"`
		ComputeResource struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type" yaml:"vim_level_resource_type"`
		} `json:"computeResource" yaml:"compute_resource" yaml:"compute_resource"`
		ZoneId             string   `json:"zoneId" yaml:"zone_id" yaml:"zone_id"`
		StorageResourceIds []string `json:"storageResourceIds" yaml:"storage_resource_ids" yaml:"storage_resource_ids"`
		ReservationId      string   `json:"reservationId" yaml:"reservation_id" yaml:"reservation_id"`
		VnfcCpInfo         []struct {
			Id             string `json:"id" yaml:"id" yaml:"id"`
			CpdId          string `json:"cpdId" yaml:"cpd_id" yaml:"cpd_id"`
			VnfExtCpId     string `json:"vnfExtCpId" yaml:"vnf_ext_cp_id" yaml:"vnf_ext_cp_id"`
			CpProtocolInfo []struct {
				LayerProtocol  string `json:"layerProtocol" yaml:"layer_protocol" yaml:"layer_protocol"`
				IpOverEthernet struct {
					MacAddress  string `json:"macAddress" yaml:"mac_address" yaml:"mac_address"`
					IpAddresses []struct {
						Type         string   `json:"type" yaml:"type" yaml:"type"`
						Addresses    []string `json:"addresses" yaml:"addresses" yaml:"addresses"`
						IsDynamic    bool     `json:"isDynamic" yaml:"is_dynamic" yaml:"is_dynamic"`
						AddressRange struct {
							MinAddress string `json:"minAddress" yaml:"min_address" yaml:"min_address"`
							MaxAddress string `json:"maxAddress" yaml:"max_address" yaml:"max_address"`
						} `json:"addressRange" yaml:"address_range" yaml:"address_range"`
						SubnetId string `json:"subnetId" yaml:"subnet_id" yaml:"subnet_id"`
					} `json:"ipAddresses" yaml:"ip_addresses" yaml:"ip_addresses"`
				} `json:"ipOverEthernet" yaml:"ip_over_ethernet" yaml:"ip_over_ethernet"`
			} `json:"cpProtocolInfo" yaml:"cp_protocol_info" yaml:"cp_protocol_info"`
			VnfLinkPortId string `json:"vnfLinkPortId" yaml:"vnf_link_port_id" yaml:"vnf_link_port_id"`
			Metadata      struct {
			} `json:"metadata" yaml:"metadata" yaml:"metadata"`
		} `json:"vnfcCpInfo" yaml:"vnfc_cp_info" yaml:"vnfc_cp_info"`
		Metadata struct {
		} `json:"metadata" yaml:"metadata" yaml:"metadata"`
	} `json:"vnfcResourceInfo" yaml:"vnfc_resource_info" yaml:"vnfc_resource_info"`
	VirtualLinkResourceInfo []struct {
		Id                   string `json:"id" yaml:"id" yaml:"id"`
		VnfVirtualLinkDescId string `json:"vnfVirtualLinkDescId" yaml:"vnf_virtual_link_desc_id" yaml:"vnf_virtual_link_desc_id"`
		NetworkResource      struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type" yaml:"vim_level_resource_type"`
		} `json:"networkResource" yaml:"network_resource" yaml:"network_resource"`
		ZoneId        string `json:"zoneId" yaml:"zone_id" yaml:"zone_id"`
		ReservationId string `json:"reservationId" yaml:"reservation_id" yaml:"reservation_id"`
		VnfLinkPorts  []struct {
			Id             string `json:"id" yaml:"id" yaml:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id" yaml:"vim_connection_id"`
				ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id" yaml:"resource_provider_id"`
				ResourceId           string `json:"resourceId" yaml:"resource_id" yaml:"resource_id"`
				VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type" yaml:"vim_level_resource_type"`
			} `json:"resourceHandle" yaml:"resource_handle" yaml:"resource_handle"`
			CpInstanceId   string `json:"cpInstanceId" yaml:"cp_instance_id" yaml:"cp_instance_id"`
			CpInstanceType string `json:"cpInstanceType" yaml:"cp_instance_type" yaml:"cp_instance_type"`
		} `json:"vnfLinkPorts" yaml:"vnf_link_ports" yaml:"vnf_link_ports"`
		Metadata struct {
		} `json:"metadata" yaml:"metadata" yaml:"metadata"`
	} `json:"virtualLinkResourceInfo" yaml:"virtual_link_resource_info" yaml:"virtual_link_resource_info"`
	VirtualStorageResourceInfo []struct {
		Id                   string `json:"id" yaml:"id" yaml:"id"`
		VirtualStorageDescId string `json:"virtualStorageDescId" yaml:"virtual_storage_desc_id" yaml:"virtual_storage_desc_id"`
		StorageResource      struct {
			VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id" yaml:"vim_connection_id"`
			ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id" yaml:"resource_provider_id"`
			ResourceId           string `json:"resourceId" yaml:"resource_id" yaml:"resource_id"`
			VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type" yaml:"vim_level_resource_type"`
		} `json:"storageResource" yaml:"storage_resource" yaml:"storage_resource"`
		ZoneId        string `json:"zoneId" yaml:"zone_id" yaml:"zone_id"`
		ReservationId string `json:"reservationId" yaml:"reservation_id" yaml:"reservation_id"`
		Metadata      struct {
		} `json:"metadata" yaml:"metadata" yaml:"metadata"`
	} `json:"virtualStorageResourceInfo" yaml:"virtual_storage_resource_info" yaml:"virtual_storage_resource_info"`
}

// ExtendedMetadata
type ExtendedMetadata struct {
	VnfPkgId       string `json:"vnfPkgId" yaml:"vnf_pkg_id"`
	VnfCatalogName string `json:"vnfCatalogName" yaml:"vnf_catalog_name"`
	ManagedBy      struct {
		ExtensionSubtype string `json:"extensionSubtype" yaml:"extension_subtype"`
		ExtensionName    string `json:"extensionName" yaml:"extension_name"`
	} `json:"managedBy" yaml:"managed_by"`
	NfType            string `json:"nfType" yaml:"nf_type"`
	LcmOperation      string `json:"lcmOperation" yaml:"lcm_operation"`
	LcmOperationState string `json:"lcmOperationState" yaml:"lcm_operation_state"`
	IsUsedByNS        string `json:"isUsedByNS" yaml:"is_used_by_ns"`
	AttachedNSCount   string `json:"attachedNSCount" yaml:"attached_ns_count"`
	ExtVirtualLinks   []struct {
		Id                 string `json:"id" yaml:"id"`
		VimConnectionId    string `json:"vimConnectionId" yaml:"vim_connection_id"`
		ResourceProviderId string `json:"resourceProviderId" yaml:"resource_provider_id"`
		ResourceId         string `json:"resourceId" yaml:"resource_id"`
		ExtCps             []struct {
			CpdId    string `json:"cpdId" yaml:"cpd_id"`
			CpConfig []struct {
				CpInstanceId   string `json:"cpInstanceId" yaml:"cp_instance_id"`
				LinkPortId     string `json:"linkPortId" yaml:"link_port_id"`
				CpProtocolData []struct {
					LayerProtocol  string `json:"layerProtocol" yaml:"layer_protocol"`
					IpOverEthernet struct {
						MacAddress  string `json:"macAddress" yaml:"mac_address"`
						IpAddresses []struct {
							Type                string   `json:"type" yaml:"type"`
							FixedAddresses      []string `json:"fixedAddresses" yaml:"fixed_addresses"`
							NumDynamicAddresses int      `json:"numDynamicAddresses" yaml:"num_dynamic_addresses"`
							AddressRange        struct {
								MinAddress string `json:"minAddress" yaml:"min_address"`
								MaxAddress string `json:"maxAddress" yaml:"max_address"`
							} `json:"addressRange" yaml:"address_range"`
							SubnetId string `json:"subnetId" yaml:"subnet_id"`
						} `json:"ipAddresses" yaml:"ip_addresses"`
					} `json:"ipOverEthernet" yaml:"ip_over_ethernet"`
				} `json:"cpProtocolData" yaml:"cp_protocol_data"`
			} `json:"cpConfig" yaml:"cp_config"`
		} `json:"extCps" yaml:"ext_cps"`
		ExtLinkPorts []struct {
			Id             string `json:"id" yaml:"id"`
			ResourceHandle struct {
				VimConnectionId      string `json:"vimConnectionId" yaml:"vim_connection_id"`
				ResourceProviderId   string `json:"resourceProviderId" yaml:"resource_provider_id"`
				ResourceId           string `json:"resourceId" yaml:"resource_id"`
				VimLevelResourceType string `json:"vimLevelResourceType" yaml:"vim_level_resource_type"`
			} `json:"resourceHandle" yaml:"resource_handle"`
		} `json:"extLinkPorts" yaml:"ext_link_ports"`
	} `json:"extVirtualLinks" yaml:"ext_virtual_links"`
	Tags []struct {
		Name        string `json:"name" yaml:"name"`
		AutoCreated bool   `json:"autoCreated" yaml:"auto_created"`
	} `json:"tags" yaml:"tags"`
}
