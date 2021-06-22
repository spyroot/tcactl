package models

//"type" : "RESOURCE_CONNECTIVITY",

// VMwareDisks VM template Disk information
type VMwareDisks struct {
	DeviceInfo struct {
		Label string `json:"label" yaml:"label"`
	} `json:"deviceInfo" yaml:"device_info"`
	Capacity     int64  `json:"capacity" yaml:"capacity"`
	DiskObjectId string `json:"diskObjectId" yaml:"diskobjectid"`
	Backing      struct {
		Datastore struct {
			Type  string `json:"type" yaml:"type"`
			Value string `json:"value" yaml:"value"`
		} `json:"datastore" yaml:"datastore"`
		FileName        string `json:"fileName" yaml:"filename"`
		BackingObjectId string `json:"backingObjectId" yaml:"backingObjectId"`
		DiskMode        string `json:"diskMode" yaml:"diskMode"`
		Split           bool   `json:"split" yaml:"split"`
		WriteThrough    bool   `json:"writeThrough" yaml:"writeThrough"`
		ThinProvisioned bool   `json:"thinProvisioned" yaml:"thinProvisioned"`
		Uuid            string `json:"uuid" yaml:"uuid"`
		ContentId       string `json:"contentId" yaml:"contentId"`
		DigestEnabled   bool   `json:"digestEnabled" yaml:"digestEnabled"`
		Sharing         string `json:"sharing" yaml:"sharing"`
	} `json:"backing:" yaml:"backing"`
}

type VmwareSummary struct {
	GuestFullName      string `json:"guestFullName" yaml:"guestFullName"`
	GuestId            string `json:"guestId" yaml:"guest_id"`
	MemorySizeMB       int    `json:"memorySizeMB" yaml:"memorySizeMb"`
	NumCpu             int    `json:"numCpu" yaml:"num_cpu"`
	CpuReservation     int    `json:"cpuReservation" yaml:"cpuReservation"`
	MemoryReservation  int    `json:"memoryReservation" yaml:"memoryReservation"`
	LatencySensitivity string `json:"latencySensitivity" yaml:"latencySensitivity"`
	DiskSize           int64  `json:"diskSize" yaml:"diskSize"`
	Annotation         struct {
		Description string `json:"description" yaml:"description"`
	} `json:"annotation" yaml:"annotation"`
}

type VMwareInventoryItem struct {
	Disks             []VMwareDisks `json:"disks" yaml:"disks"`
	Id                string        `json:"id" yaml:"id"`
	EntityId          string        `json:"entity_id" yaml:"entity_id"`
	EntityType        string        `json:"entityType" yaml:"entity_type"`
	Name              string        `json:"name" yaml:"name"`
	InstanceUuid      string        `json:"instanceUuid" yaml:"instance_uuid"`
	Uuid              string        `json:"uuid" yaml:"uuid"`
	VcenterInstanceId string        `json:"vcenter_instanceId" yaml:"vcenter_instanceId"`
	Parent            struct {
		Type  string `json:"type" yaml:"type"`
		Value string `json:"value" yaml:"value"`
	} `json:"parent" yaml:"parent"`
	Host struct {
		Type  string `json:"type" yaml:"type"`
		Value string `json:"value" yaml:"value"`
	} `json:"host" yaml:"host"`
	Summary         *VmwareSummary       `json:"summary" yaml:"summary"`
	IsTemplate      bool                 `json:"isTemplate" yaml:"isTemplate"`
	SnapshotsCount  int                  `json:"snapshotsCount" yaml:"snapshotsCount"`
	MountedISOCount int                  `json:"mountedISOCount" yaml:"mountedISOCount"`
	PowerState      string               `json:"powerState" yaml:"powerState"`
	HybridityTags   *VmwareHybridityTags `json:"hybridityTags" yaml:"hybridity_tags"`
	Containers      []struct {
		Type       string `json:"type" yaml:"type"`
		ServerGuid string `json:"serverGuid" yaml:"server_guid"`
		XsiType    string `json:"@xsi:type" yaml:"xsi_type"`
		Value      string `json:"value" yaml:"value"`
	} `json:"containers" yaml:"containers"`
	NetworkDevices []struct {
		XsiType    string `json:"@xsi:type" yaml:"xsi:type"`
		Key        int    `json:"key" yaml:"key"`
		DeviceInfo struct {
			Label   string `json:"label" yaml:"label"`
			Summary string `json:"summary" yaml:"summary"`
		} `json:"deviceInfo" yaml:"deviceInfo"`
		Backing struct {
			XsiType string `json:"@xsi:type" yaml:"xsi_type"`
			Port    struct {
				SwitchUuid       string `json:"switchUuid" yaml:"switchUuid"`
				PortgroupKey     string `json:"portgroupKey" yaml:"portgroupKey"`
				PortKey          string `json:"portKey" yaml:"portKey"`
				ConnectionCookie int    `json:"connectionCookie" yaml:"connectionCookie"`
			} `json:"port" yaml:"port"`
		} `json:"backing" yaml:"backing"`
		Connectable struct {
			MigrateConnect    string `json:"migrateConnect" yaml:"migrateConnect"`
			StartConnected    bool   `json:"startConnected" yaml:"startConnected"`
			AllowGuestControl bool   `json:"allowGuestControl" yaml:"allowGuestControl"`
			Connected         bool   `json:"connected" yaml:"connected"`
			Status            string `json:"status" yaml:"status"`
		} `json:"connectable" yaml:"connectable"`
		SlotInfo struct {
			XsiType       string `json:"@xsi:type" yaml:"xsi_type"`
			PciSlotNumber int    `json:"pciSlotNumber" yaml:"pciSlotNumber"`
		} `json:"slotInfo" yaml:"slot_info"`
		ControllerKey      int    `json:"controllerKey" yaml:"controllerKey"`
		UnitNumber         int    `json:"unitNumber" yaml:"unitNumber"`
		AddressType        string `json:"addressType" yaml:"addressType"`
		MacAddress         string `json:"macAddress" yaml:"macAddress"`
		WakeOnLanEnabled   bool   `json:"wakeOnLanEnabled" yaml:"wakeOnLanEnabled"`
		ResourceAllocation struct {
			Reservation int `json:"reservation" yaml:"reservation"`
			Share       struct {
				Shares int    `json:"shares" yaml:"shares"`
				Level  string `json:"level" yaml:"level"`
			} `json:"share" yaml:"share"`
			Limit int `json:"limit" yaml:"limit"`
		} `json:"resourceAllocation" yaml:"resourceAllocation"`
		UptCompatibilityEnabled bool `json:"uptCompatibilityEnabled" yaml:"upt_compatibility_enabled"`
	} `json:"networkDevices" yaml:"network_devices"`
	Product []struct {
		XsiType     string `json:"@xsi:type" yaml:"xsi_type"`
		Key         int    `json:"key" yaml:"key"`
		ClassId     string `json:"classId" yaml:"classid"`
		InstanceId  string `json:"instanceId" yaml:"instanceid"`
		Name        string `json:"name" yaml:"name"`
		Vendor      string `json:"vendor" yaml:"vendor"`
		Version     string `json:"version" yaml:"version"`
		FullVersion string `json:"fullVersion" yaml:"fullversion"`
		VendorUrl   string `json:"vendorUrl" yaml:"vendorurl"`
		ProductUrl  string `json:"productUrl" yaml:"product_url"`
		AppUrl      string `json:"appUrl" yaml:"app_url"`
	} `json:"product" yaml:"product"`
	Property []struct {
		XsiType          string `json:"@xsi:type" yaml:"xsi:type"`
		Key              int    `json:"key" yaml:"key"`
		ClassId          string `json:"classId" yaml:"classId"`
		InstanceId       string `json:"instanceId" yaml:"instanceId"`
		Id               string `json:"id" yaml:"id"`
		Category         string `json:"category" yaml:"category"`
		Label            string `json:"label" yaml:"label"`
		Type             string `json:"type" yaml:"type"`
		TypeReference    string `json:"typeReference" yaml:"typeReference"`
		UserConfigurable bool   `json:"userConfigurable" yaml:"userConfigurable"`
		DefaultValue     string `json:"defaultValue" yaml:"defaultValue"`
		Value            string `json:"value" yaml:"value"`
		Description      string `json:"description" yaml:"description"`
	} `json:"property" yaml:"property"`
	Origin struct {
		EndpointId   string `json:"endpointId" yaml:"endpointId"`
		EndpointType string `json:"endpointType" yaml:"endpointType"`
		EndpointName string `json:"endpointName" yaml:"endpointName"`
		ResourceId   string `json:"resourceId" yaml:"resourceId"`
		ResourceType string `json:"resourceType" yaml:"resourceType"`
		ResourceName string `json:"resourceName" yaml:"resourceName"`
	} `json:"_origin" yaml:"origin"`
	Source struct {
		Version    string `json:"version" yaml:"version"`
		Uuid       string `json:"uuid" yaml:"uuid"`
		HcspUUID   string `json:"hcspUUID" yaml:"hcspUUID"`
		SystemType string `json:"systemType" yaml:"system_type"`
	} `json:"_source" yaml:"source"`
	FullPath string `json:"fullPath" yaml:"full_path"`
}

// VcInventory list of vc inventory
type VcInventory struct {
	Items []VMwareInventoryItem `json:"items"`
}

type VmwareHybridityTags struct {
	VMWITHSNAPSHOTS struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"VM_WITH_SNAPSHOTS" yaml:"vmwithsnapshots"`
	SYSTEMRESOURCE struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"SYSTEM_RESOURCE" yaml:"systemresource"`
	SYSTEMWITHISO struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"SYSTEM_WITH_ISO" yaml:"systemwithiso"`
	VMUNDERTRANSFER struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"VM_UNDER_TRANSFER" yaml:"vmundertransfer"`
	VMUNDERSWITCHOVER struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"VM_UNDER_SWITCHOVER" yaml:"vmunderswitchover"`
	BACKUPVMAFTERMIGRATION struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"BACKUP_VM_AFTER_MIGRATION" yaml:"backupvmaftermigration"`
	VMMIGRATIONREQUESTED struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"VM_MIGRATION_REQUESTED" yaml:"vmmigrationrequested"`
	VMUNDERMIGRATION struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"VM_UNDER_MIGRATION" yaml:"vmundermigration"`
	MULEVMFORHCXRAV struct {
		Value       bool   `json:"value" yaml:"value"`
		Description string `json:"description" yaml:"description"`
	} `json:"MULE_VM_FOR_HCX_RAV" yaml:"mulevmforhcxrav"`
}
