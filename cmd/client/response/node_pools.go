package response

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"strings"
)

// NodePoolFilterType - cnf filter types
type NodePoolFilterType int32

const (
	//
	FilterPoolByID NodePoolFilterType = 0
	//
	FilterByLabel NodePoolFilterType = 1
	//
	FilterByStatus NodePoolFilterType = 2
)

// NodePoolNetwork - Node Pool contains Network information
type NodePoolNetwork struct {
	Label       string   `json:"label" yaml:"label"`
	NetworkName string   `json:"networkName" yaml:"network_name"`
	Nameservers []string `json:"nameservers" yaml:"nameservers"`
}

// NodePoolConfig - hold all metadata about node pools
type NodePoolConfig struct {
	CpuManagerPolicy struct {
		Type       string `json:"type" yaml:"type"`
		Policy     string `json:"policy" yaml:"policy"`
		Properties struct {
			KubeReserved struct {
				Cpu         int `json:"cpu" yaml:"cpu"`
				MemoryInGiB int `json:"memoryInGiB" yaml:"memory_in_gi_b"`
			} `json:"kubeReserved" yaml:"kube_reserved"`
			SystemReserved struct {
				Cpu         int `json:"cpu" yaml:"cpu"`
				MemoryInGiB int `json:"memoryInGiB" yaml:"memory_in_gi_b"`
			} `json:"systemReserved" yaml:"system_reserved"`
		} `json:"properties" yaml:"properties"`
	} `json:"cpuManagerPolicy"`
	HealthCheck struct {
		NodeStartupTimeout  string `json:"nodeStartupTimeout"`
		UnhealthyConditions []struct {
			Type    string `json:"type"`
			Status  string `json:"status"`
			Timeout string `json:"timeout"`
		} `json:"unhealthyConditions"`
	} `json:"healthCheck"`
}

type NodesDetails struct {
	Ip     string `json:"ip" yaml:"ip"`
	VmName string `json:"vmName" yaml:"vm_name"`
}

type NodesSpecs struct {
	// linked or not
	CloneMode       string            `json:"cloneMode" yaml:"clone_mode"`
	Cpu             int               `json:"cpu" yaml:"cpu"`
	Id              string            `json:"id" yaml:"id"`
	Labels          []string          `json:"labels" yaml:"labels"`
	Memory          int               `json:"memory" yaml:"memory"`
	Name            string            `json:"name" yaml:"entity_id" yaml:"name"`
	Networks        []NodePoolNetwork `json:"networks" yaml:"networks"`
	PlacementParams []struct {
		Name string `json:"name" yaml:"name"`
		Type string `json:"type" yaml:"type"`
	} `json:"placementParams" yaml:"placement_params"`
	Replica          int            `json:"replica" yaml:"replica"`
	Storage          int            `json:"storage" yaml:"storage"`
	Config           NodePoolConfig `json:"config" yaml:"config"`
	Status           string         `json:"status" yaml:"status"`
	ActiveTasksCount int            `json:"activeTasksCount" yaml:"active_tasks_count"`
	// nodes that part of cluster.
	Nodes                         []NodesDetails `json:"nodes" yaml:"nodes"`
	IsNodeCustomizationDeprecated bool           `json:"isNodeCustomizationDeprecated" yaml:"is_node_customization_deprecated"`
}

func (n *NodesSpecs) GetNodesSpecs() []NodesDetails {
	return n.Nodes
}

// NodePool - holds a list of NodePool
type NodePool struct {
	// nodes spec
	Pools []NodesSpecs `json:"items" yaml:"items"`
}

// GetIds - search pool by name or id
func (n *NodePool) GetIds() []string {

	var ids []string

	if n == nil {
		return ids
	}

	for _, it := range n.Pools {
		if len(it.Id) > 0 {
			ids = append(ids, it.Id)
		}
	}

	return ids
}

// GetPool - search pool by name or id
func (n *NodePool) GetPool(q string) (*NodesSpecs, error) {

	if n == nil {
		return nil, errors.New("uninitialized object")
	}

	for _, it := range n.Pools {
		if strings.Contains(it.Name, q) || strings.Contains(it.Id, q) {
			glog.Infof("found pool %s id %s", it.Name, it.Id)
			return &it, nil
		}
	}

	return nil, fmt.Errorf("%v not found", q)
}

// Filter filters respond based on filter type and pass to callback
func (n *NodePool) Filter(q NodePoolFilterType, f func(string) bool) (*NodePool, error) {

	if n == nil {
		return nil, fmt.Errorf("node pool instance is nil")
	}

	var pool NodePool
	filtered := make([]NodesSpecs, 0)

	for _, p := range n.Pools {
		//
		if q == FilterPoolByID && f(p.Id) {
			filtered = append(filtered, p)
		}
		// match by label
		if q == FilterByLabel {
			for _, label := range p.Labels {
				if f(label) {
					filtered = append(filtered, p)
				}
			}
		}
		//
		if q == FilterByStatus && f(p.Status) {
			filtered = append(filtered, p)
		}
	}

	pool.Pools = filtered
	return &pool, nil
}

//
//type T struct {
//	CloneMode string   `json:"cloneMode"`
//	Cpu       int      `json:"cpu"`
//	Id        string   `json:"id"`
//	Labels    []string `json:"labels"`
//	Memory    int      `json:"memory"`
//	Name      string   `json:"name"`
//	Networks  []struct {
//		Label       string   `json:"label"`
//		NetworkName string   `json:"networkName"`
//		Nameservers []string `json:"nameservers"`
//	} `json:"networks"`
//	PlacementParams []struct {
//		Name string `json:"name"`
//		Type string `json:"type"`
//	} `json:"placementParams"`
//	Replica int `json:"replica"`
//	Storage int `json:"storage"`
//	Config  struct {
//		CpuManagerPolicy struct {
//			Type       string `json:"type"`
//			Policy     string `json:"policy"`
//			Properties struct {
//				KubeReserved struct {
//					Cpu         int `json:"cpu"`
//					MemoryInGiB int `json:"memoryInGiB"`
//				} `json:"kubeReserved"`
//				SystemReserved struct {
//					Cpu         int `json:"cpu"`
//					MemoryInGiB int `json:"memoryInGiB"`
//				} `json:"systemReserved"`
//			} `json:"properties"`
//		} `json:"cpuManagerPolicy"`
//		HealthCheck struct {
//			NodeStartupTimeout  string `json:"nodeStartupTimeout"`
//			UnhealthyConditions []struct {
//				Type    string `json:"type"`
//				Status  string `json:"status"`
//				Timeout string `json:"timeout"`
//			} `json:"unhealthyConditions"`
//		} `json:"healthCheck"`
//	} `json:"config"`
//	Status           string `json:"status"`
//	ActiveTasksCount int    `json:"activeTasksCount"`
//	Nodes            []struct {
//		Ip     string `json:"ip"`
//		VmName string `json:"vmName"`
//	} `json:"nodes"`
//	IsNodeCustomizationDeprecated bool `json:"isNodeCustomizationDeprecated"`
//}
