package models

type VimExtra struct {
	DeploymentProfileId string `json:"deploymentProfileId,omitempty" yaml:"deploymentProfileId,omitempty"`
	NodeProfileName     string `json:"nodeProfileName,omitempty" yaml:"nodeProfileName,omitempty"`
	NodePoolId          string `json:"nodePoolId,omitempty" yaml:"nodePoolId,omitempty"`
	NodePoolName        string `json:"nodePoolName,omitempty" yaml:"nodePoolName,omitempty"`
	VimName             string `json:"vimName,omitempty" yaml:"vimName,omitempty"`
}

// VimConnectionInfo - Contains extra information including vim id , type
type VimConnectionInfo struct {
	Id            string `json:"id" yaml:"id"`
	VimId         string `json:"vimId" yaml:"vimId"`
	VimType       string `json:"vimType" yaml:"vimType"`
	InterfaceInfo struct {
	} `json:"interfaceInfo,omitempty" yaml:"interfaceInfo,omitempty"`
	AccessInfo struct {
	} `json:"accessInfo,omitempty" yaml:"accessInfo,omitempty"`
	Extra *VimExtra `json:"extra,omitempty" yaml:"extra,omitempty"`
}

// CnfPolicyUri CNF policy URI
type CnfPolicyUri struct {
	Href string `json:"href,omitempty"`
}

type PolicyLinks struct {
	Self           CnfPolicyUri `json:"self,omitempty" yaml:"self"`
	Indicators     CnfPolicyUri `json:"indicators,omitempty" yaml:"indicators"`
	Instantiate    CnfPolicyUri `json:"instantiate,omitempty" yaml:"instantiate"`
	Retry          CnfPolicyUri `json:"retry,omitempty" yaml:"retry"`
	Rollback       CnfPolicyUri `json:"rollback,omitempty" yaml:"rollback"`
	UpdateState    CnfPolicyUri `json:"update_state,omitempty" yaml:"update_state"`
	Terminate      CnfPolicyUri `json:"terminate,omitempty" yaml:"terminate"`
	Scale          CnfPolicyUri `json:"scale,omitempty" yaml:"scale"`
	ScaleToLevel   CnfPolicyUri `json:"scaleToLevel,omitempty" yaml:"scaleToLevel"`
	Heal           CnfPolicyUri `json:"heal,omitempty" yaml:"heal"`
	Update         CnfPolicyUri `json:"update,omitempty" yaml:"update"`
	UpgradePackage CnfPolicyUri `json:"upgrade_package,omitempty" yaml:"upgrade_package"`
	Upgrade        CnfPolicyUri `json:"upgrade,omitempty" yaml:"upgrade"`
	Reconfigure    CnfPolicyUri `json:"reconfigure,omitempty" yaml:"reconfigure"`
	ChangeFlavour  CnfPolicyUri `json:"changeFlavour,omitempty" yaml:"changeFlavour"`
	Operate        CnfPolicyUri `json:"operate,omitempty" yaml:"operate"`
	ChangeExtConn  CnfPolicyUri `json:"changeExtConn,omitempty" yaml:"changeExtConn"`
}

type Networks struct {
	Label       string   `json:"label" yaml:"label"`
	NetworkName string   `json:"networkName" yaml:"networkName"`
	Nameservers []string `json:"nameservers" yaml:"nameservers"`
}

func NewNetworks(label string, networkName string, nameservers []string) *Networks {
	return &Networks{Label: label, NetworkName: networkName, Nameservers: nameservers}
}

// Location location details
type Location struct {
	City      string  `json:"city" yaml:"city"`
	Country   string  `json:"country" yaml:"country"`
	CityASCII string  `json:"cityAscii" yaml:"cityAscii"`
	Latitude  float64 `json:"latitude" yaml:"latitude"`
	Longitude float64 `json:"longitude" yaml:"longitude"`
}

// VimConnection Connection status
type VimConnection struct {
	Status              string `json:"status" yaml:"status"`
	RemoteStatus        string `json:"remoteStatus" yaml:"remoteStatus"`
	VimConnectionStatus string `json:"vimConnectionStatus" yaml:"vimConnectionStatus"`
}

// PlacementParams Node Placement
type PlacementParams struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
}

func NewPlacementParams(name string, Type string) *PlacementParams {
	return &PlacementParams{Name: name, Type: Type}
}

// Networks
type TypeNode struct {
	Name            string            `json:"name" yaml:"name"`
	Networks        []Networks        `json:"networks" yaml:"networks"`
	PlacementParams []PlacementParams `json:"placementParams" yaml:"placementParams"`
}

func NewTypeNode(name string, networks []Networks, placementParams []PlacementParams) *TypeNode {
	return &TypeNode{Name: name, Networks: networks, PlacementParams: placementParams}
}
