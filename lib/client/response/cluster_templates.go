package response

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/request"
	"reflect"
	"strings"
)

type TemplateFilterType int32
type TemplateType int32

const (
	// FilterTemplateType by cnf id
	FilterTemplateType TemplateFilterType = 0

	// FilterTemplateKubeVersion filer by vnf instance name
	FilterTemplateKubeVersion TemplateFilterType = 1

	// FilterById filer by vnf instance name
	FilterById TemplateFilterType = 1

	// CniTypeAntrea
	CniTypeAntrea = "antrea"

	// CniTypeMultus
	CniTypeMultus = "multus"

	// CniTypeCalico - calico cni
	CniTypeCalico = "calico"

	// CniTypeWhereAbouts - cni CniTypeWhereAbouts
	CniTypeWhereAbouts = "whereabouts"

	//vsphere-csi
	//nfs_client
	//v1.20.4+vmware.1
	//helm
	// 2.17.0

	// TemplateTypeMgmt cluster template mgmt
	TemplateTypeMgmt TemplateType = 0

	//TemplateTypeWorkload cluster template workload
	TemplateTypeWorkload TemplateType = 1

	TemplateMgmt string = "MANAGEMENT"

	TemplateWorkload string = "WORKLOAD"
)

// ClusterConfigSpec cluster config spec hols CNI/CSI
// and K8S version
type ClusterConfigSpec struct {
	Cni []struct {
		Name       string `json:"name" yaml:"name"`
		Properties struct {
		} `json:"properties" yaml:"properties"`
	} `json:"cni" yaml:"cni"`
	Csi []struct {
		Name       string `json:"name" yaml:"name"`
		Properties struct {
			Name      string `json:"name" yaml:"name"`
			IsDefault bool   `json:"isDefault" yaml:"isDefault"`
			Timeout   string `json:"timeout" yaml:"timeout"`
		} `json:"properties" yaml:"properties"`
	} `json:"csi" yaml:"csi"`
	KubernetesVersion string `json:"kubernetesVersion" yaml:"kubernetesVersion"`
	Tools             []struct {
		Name    string `json:"name" yaml:"name"`
		Version string `json:"version" yaml:"version"`
	} `json:"tools" yaml:"tools"`
}

// GetField - return field from Cluster Spec struct
func (t *ClusterConfigSpec) GetField(field string) string {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	return f.String()
}

// HealthCheckSpec - specs
// Type
type HealthCheckSpec struct {
	NodeStartupTimeout  string `json:"nodeStartupTimeout,omitempty" yaml:"nodeStartupTimeout,omitempty"`
	UnhealthyConditions []struct {
		Type    string `json:"type,omitempty" yaml:"type,omitempty"`
		Status  string `json:"status,omitempty" yaml:"status,omitempty"`
		Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	} `json:"unhealthyConditions,omitempty" yaml:"unhealthyConditions,omitempty"`
}

// Properties -CpuManagerPolicy properties
type Properties struct {
	KubeReserved struct {
		Cpu         int `json:"cpu,omitempty" yaml:"cpu,omitempty"`
		MemoryInGiB int `json:"memoryInGiB,omitempty" yaml:"memoryInGiB,omitempty"`
	} `json:"kubeReserved,omitempty" yaml:"kubeReserved,omitempty"`
	SystemReserved struct {
		Cpu         int `json:"cpu,omitempty" yaml:"cpu,omitempty"`
		MemoryInGiB int `json:"memoryInGiB,omitempty" yaml:"memoryInGiB,omitempty"`
	} `json:"systemReserved,omitempty" yaml:"systemReserved,omitempty"`
}

// CpuManagerPolicy - overwrite CPU mgmt policy
type CpuManagerPolicy struct {
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Policy      string `json:"policy,omitempty" yaml:"policy,omitempty"`
	*Properties `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type Tags struct {
	AutoCreated bool   `json:"autoCreated" yaml:"autoCreated"`
	Name        string `json:"name" yaml:"name"`
}

type TemplateNetworks struct {
	Label string `json:"label" yaml:"label" validate:"required"`
}

// ClusterTemplate - cluster template
type ClusterTemplate struct {
	ClusterType   string             `json:"clusterType" yaml:"clusterType" validate:"required"`
	ClusterConfig *ClusterConfigSpec `json:"clusterConfig,omitempty" yaml:"clusterConfig,omitempty" validate:"required"`
	Description   string             `json:"description" yaml:"description"`
	MasterNodes   []struct {
		Cpu       int                `json:"cpu" yaml:"cpu" validate:"required"`
		Memory    int                `json:"memory" yaml:"memory" validate:"required"`
		Name      string             `json:"name" yaml:"name"`
		Networks  []TemplateNetworks `json:"networks" yaml:"networks" validate:"required"`
		Storage   int                `json:"storage" yaml:"storage" validate:"required"`
		Replica   int                `json:"replica" yaml:"replica" validate:"required"`
		Labels    []string           `json:"labels" yaml:"labels"`
		CloneMode string             `json:"cloneMode" yaml:"cloneMode" validate:"required"`
	} `json:"masterNodes" yaml:"masterNodes"`
	Name        string `json:"name" yaml:"name" validate:"required"`
	Id          string `json:"id" yaml:"id"`
	Tags        []Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	WorkerNodes []struct {
		Cpu      int    `json:"cpu" yaml:"cpu" validate:"required"`
		Memory   int    `json:"memory" yaml:"memory" validate:"required"`
		Name     string `json:"name" yaml:"name" validate:"required"`
		Networks []struct {
			Label string `json:"label" yaml:"label" validate:"required"`
		} `json:"networks" yaml:"networks" validate:"required"`
		Storage   int      `json:"storage" yaml:"storage" validate:"required"`
		Replica   int      `json:"replica" yaml:"replica" validate:"required"`
		Labels    []string `json:"labels" yaml:"labels" validate:"required"`
		CloneMode string   `json:"cloneMode" yaml:"cloneMode" validate:"required"`
		Config    struct {
			CpuManagerPolicy *CpuManagerPolicy `json:"cpuManagerPolicy,omitempty" yaml:"cpuManagerPolicy,omitempty"`
			HealthCheckSpec  *HealthCheckSpec  `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
		} `json:"config,omitempty" yaml:"config,omitempty"`
	} `json:"workerNodes,omitempty" yaml:"workerNodes,omitempty"`
}

// GetField - return field from Cluster Spec struct
func (t *ClusterTemplate) GetField(field string) string {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	return f.String()
}

// ValidateSpec - validate cluster specs contains all required node pool
// based on template spec
func (t *ClusterTemplate) ValidateSpec(spec *request.Cluster) (bool, error) {

	if t == nil {
		return false, fmt.Errorf("cluster template is nil")
	}

	if spec == nil {
		return false, fmt.Errorf("cluster spec is nil")
	}

	workerPools := 0
	masterPools := 0

	for _, node := range t.WorkerNodes {
		if spec.FindNodePoolByName(node.Name, true) {
			workerPools++
		}
	}
	for _, node := range t.MasterNodes {
		if spec.FindNodePoolByName(node.Name, false) {
			masterPools++
		}
	}

	if workerPools != len(t.WorkerNodes) {
		return false, fmt.Errorf("check number of worker node pool in spec, "+
			"must be %d found %d", len(t.WorkerNodes), workerPools)
	}

	if masterPools != len(t.MasterNodes) {
		return false, fmt.Errorf("check number of node pool in spec, "+
			"must be %d found %d", len(t.WorkerNodes), masterPools)
	}

	return workerPools == len(t.WorkerNodes) && masterPools == len(t.MasterNodes), nil
}

type ClusterTemplates struct {
	ClusterTemplates []ClusterTemplate
}

type TemplateNotFound struct {
	errMsg string
}

func (m *TemplateNotFound) Error() string {
	return " cluster template '" + m.errMsg + "' not found."
}

// GetTemplateId return cluster template id
func (t *ClusterTemplates) GetTemplateId(q string) (string, error) {

	if t == nil {
		return "", fmt.Errorf("uninitialized object")
	}

	for _, it := range t.ClusterTemplates {
		if it.Name == q || it.Id == q {
			glog.Infof("Found template %v cluster id %v", q, it.Id)
			return it.Id, nil
		}
	}

	return "", &TemplateNotFound{errMsg: q}
}

// GetTemplate return cluster template , lookup by name or id
func (t *ClusterTemplates) GetTemplate(q string) (*ClusterTemplate, error) {

	if t == nil {
		return nil, fmt.Errorf("uninitialized object")
	}

	for _, it := range t.ClusterTemplates {
		if it.Name == q || it.Id == q {
			glog.Infof("Found template %v cluster id %v", q, it.Id)
			return &it, nil
		}
	}

	return nil, &TemplateNotFound{errMsg: q}
}

// Filter filters respond based on filter type and pass to callback
func (t *ClusterTemplates) Filter(q TemplateFilterType, f func(string) bool) (*ClusterTemplates, error) {

	if t == nil {
		return nil, fmt.Errorf("instance is nil")
	}

	filtered := make([]ClusterTemplate, 0)
	for _, tmpl := range t.ClusterTemplates {
		if q == FilterById && f(tmpl.Id) {
			filtered = append(filtered, tmpl)
		}
		if q == FilterTemplateType && f(tmpl.ClusterType) {
			filtered = append(filtered, tmpl)
		}
		if q == FilterTemplateKubeVersion && f(tmpl.ClusterConfig.KubernetesVersion) {
			filtered = append(filtered, tmpl)
		}
	}
	return &ClusterTemplates{filtered}, nil
}
