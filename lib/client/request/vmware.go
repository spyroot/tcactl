package request

import "github.com/spyroot/tcactl/lib/models"

type ClusterEntityTaskFilter struct {
	EntityIds []string `json:"entityIds" yaml:"entity_ids"`
}

func NewClusterEntityTaskFilter(id string) *ClusterEntityTaskFilter {
	return &ClusterEntityTaskFilter{EntityIds: []string{id}}
}

type ClusterTaskQuery struct {
	Filter            ClusterEntityTaskFilter `json:"filter" yaml:"filter" yaml:"filter"`
	IncludeChildTasks bool                    `json:"includeChildTasks,omitempty" yaml:"includeChildTasks"`
}

func NewClusterTaskQuery(id string) *ClusterTaskQuery {
	return &ClusterTaskQuery{
		Filter: *NewClusterEntityTaskFilter(id),
	}
}

type FilterCloud struct {
	// 20210212053126765-51947d48-91b1-447e-ba40-668eb411f545
	EndpointId string `json:"endpointId,omitempty"`
}

// ClusterFilterQuery - filter based on endpoint cloud id
type ClusterFilterQuery struct {
	Filter struct {
		// EntityTypes
		// "cluster"
		// cluster, computeResource
		EntityTypes []string     `json:"entityTypes" yaml:"entity_types"`
		Cloud       *FilterCloud `json:"cloud" yaml:"cloud"`
	} `json:"filter"`
}

// NewClusterFilterQuery creates new Filter query
// in order to get Folder view.
func NewClusterFilterQuery(cloudId string) *ClusterFilterQuery {
	v := ClusterFilterQuery{}
	v.Filter.Cloud = &FilterCloud{}
	v.Filter.Cloud.EndpointId = cloudId
	v.Filter.EntityTypes = []string{models.EntityTypeCluster, models.EntityTypeComputeResource}
	return &v
}

// VMwareNetworkQuery - filter
type VMwareNetworkQuery struct {
	Filter struct {
		// tenant id 20210212053126765-51947d48-91b1-447e-ba40-668eb411f545
		TenantId string `json:"tenantId"`
		// domain id in VC format domain-c13
		ClusterId string `json:"clusterId"`
	} `json:"filter"`
}

//
type VMwareTemplateQuery struct {
	Filter struct {
		Cloud *FilterCloud `json:"cloud"`
		// templateType k8svm
		TemplateType string `json:"templateType,omitempty"`
		// v1.20.4+vmware.1
		K8SVersion string `json:"k8sVersion,omitempty"`
	} `json:"filter"`
}

//
func NewVMwareTemplateQuery(cloudId string, templateType string, ver string) *VMwareTemplateQuery {

	v := VMwareTemplateQuery{}
	v.Filter.Cloud = &FilterCloud{}
	v.Filter.Cloud.EndpointId = cloudId
	v.Filter.TemplateType = templateType
	v.Filter.K8SVersion = ver

	return &v
}

type VMwareResourcePoolQuery struct {
	Filter struct {
		// resourcePool
		EntityTypes []string     `json:"entityTypes"`
		Cloud       *FilterCloud `json:"cloud,omitempty"`
		// "childTypes":["VirtualMachine","VirtualApp"],"container":"domain-c8"}}
		ChildTypes []string `json:"childTypes,omitempty"`
		Container  string   `json:"container,omitempty"`
	} `json:"filter"`
}

// VmwareFolderQuery folder view query
type VmwareFolderQuery struct {
	Filter struct {
		// folder
		EntityTypes []string     `json:"entityTypes"`
		Cloud       *FilterCloud `json:"cloud,omitempty"`
		//"childTypes":["VirtualMachine","VirtualApp"]
		//childTypes":["VirtualMachine","VirtualApp"]}}
		ChildTypes []string `json:"childTypes"`
	} `json:"filter"`
}

// NewVmwareFolderQuery creates new Filter query
// in order to get Folder view.
func NewVmwareFolderQuery(cloudId string) *VmwareFolderQuery {
	v := VmwareFolderQuery{}
	v.Filter.Cloud = &FilterCloud{}
	v.Filter.Cloud.EndpointId = cloudId
	v.Filter.EntityTypes = []string{models.EntityTypeFolder}
	v.Filter.ChildTypes = []string{models.ChildVirtualTypeMachine, models.ChildTypeVirtualApp}
	return &v
}

// NewVMwareResourcePoolQuery create a new Filter
// to query resource pools.
func NewVMwareResourcePoolQuery(cloudId string) *VMwareResourcePoolQuery {
	v := VMwareResourcePoolQuery{}
	v.Filter.Cloud = &FilterCloud{}
	v.Filter.Cloud.EndpointId = cloudId
	v.Filter.EntityTypes = []string{models.EntityTypeResourcePool}
	v.Filter.ChildTypes = []string{models.ChildVirtualTypeMachine, models.ChildTypeVirtualApp}
	return &v
}
