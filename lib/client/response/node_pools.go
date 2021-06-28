package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/models"
	"reflect"
	"strings"
)

// NodePoolFilterType - cnf filter types
type NodePoolFilterType int32

//type NewNodePoolSpec struct {
//	Id          string `json:"id"`
//	OperationId string `json:"operationId"`
//}

const (
	//
	FilterPoolByID NodePoolFilterType = 0
	//
	FilterByLabel NodePoolFilterType = 1
	//
	FilterByStatus NodePoolFilterType = 2
)

// NodePoolConfig - hold all metadata about node pools
type NodePoolConfig struct {
	CpuManagerPolicy struct {
		Type       string `json:"type" yaml:"type"`
		Policy     string `json:"policy" yaml:"policy"`
		Properties struct {
			KubeReserved struct {
				Cpu         int `json:"cpu" yaml:"cpu"`
				MemoryInGiB int `json:"memoryInGiB" yaml:"memoryInGiB"`
			} `json:"kubeReserved" yaml:"kube_reserved"`
			SystemReserved struct {
				Cpu         int `json:"cpu" yaml:"cpu"`
				MemoryInGiB int `json:"memoryInGiB" yaml:"memoryInGiB"`
			} `json:"systemReserved" yaml:"system_reserved"`
		} `json:"properties" yaml:"properties"`
	} `json:"cpuManagerPolicy"`
	HealthCheck *models.HealthCheck `json:"healthCheck"`
}

type NodesDetails struct {
	Ip     string `json:"ip" yaml:"ip"`
	VmName string `json:"vmName" yaml:"vm_name"`
}

type NodesSpecs struct {
	// linked or not
	CloneMode       string           `json:"cloneMode" yaml:"cloneMode"`
	Cpu             int              `json:"cpu" yaml:"cpu"`
	Id              string           `json:"id" yaml:"id"`
	Labels          []string         `json:"labels" yaml:"labels"`
	Memory          int              `json:"memory" yaml:"memory"`
	Name            string           `json:"name" yaml:"entity_id" yaml:"name"`
	Networks        []models.Network `json:"networks" yaml:"networks"`
	PlacementParams []struct {
		Name string `json:"name" yaml:"name"`
		Type string `json:"type" yaml:"type"`
	} `json:"placementParams" yaml:"placement_params"`
	Replica          int             `json:"replica" yaml:"replica"`
	Storage          int             `json:"storage" yaml:"storage"`
	Config           *NodePoolConfig `json:"config" yaml:"config"`
	Status           string          `json:"status" yaml:"status"`
	ActiveTasksCount int             `json:"activeTasksCount" yaml:"active_tasks_count"`
	// nodes that part of cluster.
	Nodes                         []NodesDetails `json:"nodes" yaml:"nodes"`
	IsNodeCustomizationDeprecated bool           `json:"isNodeCustomizationDeprecated" yaml:"is_node_customization_deprecated"`
}

// GetField - return struct field value
func (t *NodesSpecs) GetField(field string) string {

	r := reflect.ValueOf(t)
	fields, _ := t.GetFields()
	if _, ok := fields[field]; ok {
		f := reflect.Indirect(r).FieldByName(strings.Title(field))
		return f.String()
	}

	return ""
}

// GetFields return VduPackage fields name as
// map[string], each key is field name
func (t *NodesSpecs) GetFields() (map[string]interface{}, error) {

	var m map[string]interface{}

	b, err := json.Marshal(t)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}

	return m, nil
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

type PoolNotFound struct {
	ErrMsg string
}

func (m *PoolNotFound) Error() string {
	return "pool '" + m.ErrMsg + "' not found"
}

// GetPool - search for particular pool
// by name or id
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

	return nil, &PoolNotFound{ErrMsg: q}
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
