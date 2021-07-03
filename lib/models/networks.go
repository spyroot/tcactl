package models

import "strings"

const (
	// NetworkActive is active or not
	NetworkActive = "ACTIVE"

	// DvsPortGroupType type for DVS port group
	DvsPortGroupType = "DistributedVirtualPortgroup"
)

// Network Generic Network json struct.
// it used in many specs
type Network struct {
	Label       string   `json:"label" yaml:"label"`
	NetworkName string   `json:"networkName" yaml:"networkName"`
	Nameservers []string `json:"nameservers" yaml:"nameservers"`
}

// NetworkSpec is network spec returned by API
// Example
// 	"status" : "ACTIVE",
// 	"tenantId" : "20210212053126765-51947d48-91b1-447e-ba40-668eb411f545",
// 	"id" : "dvportgroup-69009",
// 	"name" : "tkg-dhcp-vlan1007-10.241.7.0",
// 	"dvsName" : "core02-services",
// 	"fullNetworkPath" : "/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0",
// 	"networkType" : "vlan",
// 	"isShared" : false,
// 	"type" : "DistributedVirtualPortgroup"
type NetworkSpec struct {
	Status string `json:"status" yaml:"status"`

	//tenantId is in format 20210212053126765-51947d48-91b1-447e-ba40-668eb411f545
	TenantId string `json:"tenantId" yaml:"tenantId"`

	//Id is vc id dvportgroup-69009
	Id string `json:"id" yaml:"id"`

	//Name is port-group name for example tkg-dhcp-vlan100x-10.x.x.0
	Name string `json:"name" yaml:"name"`

	// dvsName for example core02-services
	DvsName string `json:"dvsName" yaml:"dvsName"`

	//FullNetworkPath full vc path "/Datacenter/network/tkg-dhcp-vlan100X-10.x.x.x",

	FullNetworkPath string `json:"fullNetworkPath" yaml:"fullNetworkPath"`

	// NetworkType vlan is vlan type
	NetworkType string `json:"networkType" yaml:"networkType"`

	IsShared bool `json:"isShared" yaml:"isShared"`

	// Type DistributedVirtualPortgroup
	Type string `json:"type" yaml:"type"`
}

// IsActive if network is active
func (s *NetworkSpec) IsActive() bool {
	if s == nil {
		return false
	}

	return s.Status == NetworkActive
}

// IsDvsPortGroup - return true if network backend is port-group
func (s *NetworkSpec) IsDvsPortGroup() bool {
	if s == nil {
		return false
	}

	return s.Type == DvsPortGroupType
}

type CloudNetworks struct {
	Network []NetworkSpec `json:"networks" yaml:"type"`
}

// NetworkNotFound error raised if network not found
type NetworkNotFound struct {
	errMsg string
}

//
func (m *NetworkNotFound) Error() string {
	return " network '" + m.errMsg + "' not found"
}

//GetNetwork return network spec as NetworkSpec
func (n *CloudNetworks) GetNetwork(name string) (*NetworkSpec, error) {

	for _, spec := range n.Network {
		if spec.Name == name {
			return &spec, nil
		}
	}

	return nil, &NetworkNotFound{errMsg: name}
}

//FindFullNetwork return network spec as NetworkSpec
func (n *CloudNetworks) FindFullNetwork(name string) (*NetworkSpec, error) {

	for _, spec := range n.Network {
		if spec.FullNetworkPath == name {
			return &spec, nil
		}
	}

	return nil, &NetworkNotFound{errMsg: name}
}
func (n *CloudNetworks) NormalizeName(name string) (string, error) {

	// if client provide full path , validate
	if strings.HasPrefix(name, "/Datacenter/network/") {
		network, err := n.FindFullNetwork(name)
		if err != nil {
			return "", err
		}

		return network.FullNetworkPath, nil
	}

	// otherwise name , resolve to full network path
	network, err := n.GetNetwork(name)
	if err != nil {
		return "", err
	}

	return network.FullNetworkPath, nil
}
