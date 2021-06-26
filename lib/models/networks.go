package models

//Example
// "status" : "ACTIVE",
// "tenantId" : "20210212053126765-51947d48-91b1-447e-ba40-668eb411f545",
// "id" : "dvportgroup-69009",
// "name" : "tkg-dhcp-vlan1007-10.241.7.0",
// "dvsName" : "core02-services",
// "fullNetworkPath" : "/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0",
// "networkType" : "vlan",
// "isShared" : false,
// "type" : "DistributedVirtualPortgroup"

type NetworkSpec struct {
	Status          string `json:"status" yaml:"status"`
	TenantId        string `json:"tenantId" yaml:"tenantId"`
	Id              string `json:"id" yaml:"id"`
	Name            string `json:"name" yaml:"name"`
	DvsName         string `json:"dvsName" yaml:"dvsName"`
	FullNetworkPath string `json:"fullNetworkPath" yaml:"fullNetworkPath"`
	NetworkType     string `json:"networkType" yaml:"networkType"`
	IsShared        bool   `json:"isShared" yaml:"isShared"`
	Type            string `json:"type" yaml:"type"`
}

type CloudNetworks struct {
	Network []NetworkSpec `json:"networks" yaml:"type"`
}
