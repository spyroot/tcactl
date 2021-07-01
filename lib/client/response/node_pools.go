package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/models"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

// NodePoolFilterType - cnf filter types
type NodePoolFilterType int32

const (
	// FilterPoolByID filter node pool by node pool type
	FilterPoolByID NodePoolFilterType = 0

	// FilterByLabel filters node pool by label
	FilterByLabel NodePoolFilterType = 1

	// FilterByStatus filters node pool by status
	FilterByStatus NodePoolFilterType = 2
)

type NodesSpecs struct {
	// linked or not
	CloneMode        string                   `json:"cloneMode,omitempty" yaml:"cloneMode,omitempty"`
	Cpu              int                      `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Id               string                   `json:"id,omitempty" yaml:"id,omitempty"`
	Labels           []string                 `json:"labels,omitempty" yaml:"labels,omitempty"`
	Memory           int                      `json:"memory,omitempty" yaml:"memory,omitempty"`
	Name             string                   `json:"name,omitempty" yaml:"name,omitempty"`
	Networks         []models.Network         `json:"networks,omitempty" yaml:"networks,omitempty"`
	PlacementParams  []models.PlacementParams `json:"placementParams,omitempty" yaml:"placementParams,omitempty"`
	Replica          int                      `json:"replica,omitempty" yaml:"replica,omitempty"`
	Storage          int                      `json:"storage,omitempty" yaml:"storage,omitempty"`
	Config           *models.NodePoolConfig   `json:"config,omitempty" yaml:"config,omitempty"`
	Status           string                   `json:"status,omitempty" yaml:"status,omitempty"`
	ActiveTasksCount int                      `json:"activeTasksCount,omitempty" yaml:"activeTasksCount,omitempty"`
	// nodes that part of cluster.
	Nodes                         []models.Nodes `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	IsNodeCustomizationDeprecated bool           `json:"isNodeCustomizationDeprecated,omitempty" yaml:"isNodeCustomizationDeprecated,omitempty"`
}

// GetField - return struct field value
func (n *NodesSpecs) GetField(field string) string {

	r := reflect.ValueOf(n)
	fields, _ := n.GetFields()
	if _, ok := fields[field]; ok {
		f := reflect.Indirect(r).FieldByName(strings.Title(field))
		return f.String()
	}

	return ""
}

// GetFields return VduPackage fields name as
// map[string], each key is field name
func (n *NodesSpecs) GetFields() (map[string]interface{}, error) {

	var m map[string]interface{}

	b, err := json.Marshal(n)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}

	return m, nil
}

func (n *NodesSpecs) GetNodesSpecs() []models.Nodes {
	return n.Nodes
}

// NodePool - holds a list of NodePool
type NodePool struct {
	// nodes spec
	Pools []NodesSpecs `json:"items" yaml:"items"`
}

func NewNodePool(n *NodesSpecs) *NodePool {
	p := NodePool{}
	p.Pools = append(p.Pools, *n)
	return &p
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

// GetPoolByName - search for particular pool
// by name only
func (n *NodePool) GetPoolByName(name string) (*NodesSpecs, error) {

	if n == nil {
		return nil, errors.New("uninitialized object")
	}

	for _, it := range n.Pools {
		if it.Name == name {
			glog.Infof("found pool %s id %s", it.Name, it.Id)
			return &it, nil
		}
	}

	return nil, &PoolNotFound{ErrMsg: name}
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
				spliced := strings.Split(label, "=")
				if len(spliced) == 2 {
					if f(spliced[1]) {
						filtered = append(filtered, p)
					}
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

// NodeSpecsFromFile - reads node specs from file
// and return TenantSpecs instance
func NodeSpecsFromFile(fileName string) (*NodesSpecs, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadNodeSpecs(file)
}

// NodeSpecsFromString take string that hold entire spec
// passed to reader and return TenantSpecs instance
func NodeSpecsFromString(str string) (*NodesSpecs, error) {
	r := strings.NewReader(str)
	return ReadNodeSpecs(r)
}

// ReadNodeSpecs - Read node spec from io interface
// detects format and use either yaml or json parse
func ReadNodeSpecs(b io.Reader) (*NodesSpecs, error) {

	var spec NodesSpecs

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	return nil, &InvalidTenantsSpec{"unknown format"}
}

//InstanceSpecsFromString method return instance form string
func (n NodesSpecs) InstanceSpecsFromString(s string) (interface{}, error) {
	return TenantsSpecsFromString(s)
}

// AsJsonString return object as json string,
// it mainly used for testing.
func (n *NodesSpecs) AsJsonString() (string, error) {
	if n == nil {
		return "", nil
	}

	b, err := json.Marshal(n)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// AsYamlString return object as json yaml string,
// it mainly used for testing.
func (n *NodesSpecs) AsYamlString() (string, error) {
	if n == nil {
		return "", nil
	}

	b, err := yaml.Marshal(n)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// NewNodePoolSpecs create spec from reader
func NewNodePoolSpecs(r io.Reader) (*NodesSpecs, error) {
	spec, err := ReadNodeSpecs(r)
	if err != nil {
		return nil, err
	}

	return spec, nil
}
