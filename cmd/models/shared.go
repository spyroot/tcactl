package models

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

// TypeNode
type TypeNode struct {
	Name            string            `json:"name" yaml:"name"`
	Networks        []Networks        `json:"networks" yaml:"networks"`
	PlacementParams []PlacementParams `json:"placementParams" yaml:"placementParams"`
}

func NewTypeNode(name string, networks []Networks, placementParams []PlacementParams) *TypeNode {
	return &TypeNode{Name: name, Networks: networks, PlacementParams: placementParams}
}
