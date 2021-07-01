package request

import (
	"encoding/json"
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

type ClusterType string

const (
	// ClusterManagement management k8s cluster
	ClusterManagement ClusterType = "MANAGEMENT"

	// ClusterWorkload workload k8s cluster
	ClusterWorkload ClusterType = "WORKLOAD"
)

type ClusterConfig struct {
	Csi []struct {
		Name       string `json:"name" yaml:"name"`
		Properties struct {
			ServerIP      string `json:"serverIP,omitempty" yaml:"serverIP,omitempty"`
			MountPath     string `json:"mountPath,omitempty" yaml:"mountPath,omitempty"`
			DatastoreUrl  string `json:"datastoreUrl,omitempty" yaml:"datastoreUrl,omitempty"`
			DatastoreName string `json:"datastoreName,omitempty" yaml:"datastoreName,omitempty"`
		} `json:"properties" yaml:"properties"`
	} `json:"csi" yaml:"csi,omitempty"`
	Tools []struct {
		Name       string `json:"name,omitempty" yaml:"name,omitempty"`
		Properties struct {
			ExtensionId string `json:"extensionId,omitempty" yaml:"extensionId,omitempty"`
			Password    string `json:"password,omitempty" yaml:"password,omitempty"`
			Type        string `json:"type,omitempty" yaml:"type,omitempty"`
			Url         string `json:"url,omitempty" yaml:"url,omitempty"`
			Username    string `json:"username,omitempty" yaml:"username,omitempty"`
		} `json:"properties" yaml:"properties"`
	} `json:"tools" yaml:"tools"`
	SystemSettings []struct {
		Name       string `json:"name,omitempty" yaml:"name,omitempty"`
		Properties struct {
			Host     string `json:"host,omitempty" yaml:"host,omitempty"`
			Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
			Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
		} `json:"properties" yaml:"properties,omitempty"`
	} `json:"systemSettings,omitempty" yaml:"systemSettings,omitempty"`
}

// Cluster new cluster creation request
type Cluster struct {
	ClusterPassword     string                   `json:"clusterPassword" yaml:"clusterPassword"`
	ClusterTemplateId   string                   `json:"clusterTemplateId" yaml:"clusterTemplateId"`
	ClusterType         string                   `json:"clusterType" yaml:"clusterType"`
	Description         string                   `json:"description,omitempty" yaml:"description,omitempty"`
	Location            *models.Location         `json:"location,omitempty" yaml:"location,omitempty"`
	ClusterConfig       *ClusterConfig           `json:"clusterConfig,omitempty" yaml:"clusterConfig,omitempty"`
	HcxCloudUrl         string                   `json:"hcxCloudUrl" yaml:"hcxCloudUrl"`
	EndpointIP          string                   `json:"endpointIP" yaml:"endpointIP"`
	ManagementClusterId string                   `json:"managementClusterId,omitempty" yaml:"managementClusterId,omitempty"`
	Name                string                   `json:"name" yaml:"name"`
	VmTemplate          string                   `json:"vmTemplate" yaml:"vmTemplate"`
	MasterNodes         []models.TypeNode        `json:"masterNodes" yaml:"masterNodes"`
	WorkerNodes         []models.TypeNode        `json:"workerNodes" yaml:"workerNodes"`
	PlacementParams     []models.PlacementParams `json:"placementParams" yaml:"placementParams"`
}

func (c *Cluster) IsManagement() bool {
	return strings.ToLower(c.ClusterType) == strings.ToLower(string(ClusterManagement))
}

func (c *Cluster) IsWorkload() bool {
	return strings.ToLower(c.ClusterType) == strings.ToLower(string(ClusterWorkload))
}

// InvalidClusterSpec error if specs invalid
type InvalidClusterSpec struct {
	errMsg string
}

//
func (m *InvalidClusterSpec) Error() string {
	return m.errMsg
}

//FindNodePoolByName search for node pool name
// if isWorker will check worker node pool, otherwise Master node pools.
func (c *Cluster) FindNodePoolByName(name string, isWorker bool) bool {

	nodes := c.MasterNodes

	if isWorker {
		nodes = c.WorkerNodes
	}

	for _, node := range nodes {
		if node.Name == name {
			glog.Infof("Found node name %v", name)
			return true
		}
	}

	return false
}

// ClusterSpecsFromFile - reads tenant spec from file
// and return TenantSpecs instance
func ClusterSpecsFromFile(fileName string) (*Cluster, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadClusterSpec(file)
}

// ClusterSpecsFromString take string that holdw entire spec
// passed to reader and return TenantSpecs instance
func ClusterSpecsFromString(str string) (*Cluster, error) {
	r := strings.NewReader(str)
	return ReadClusterSpec(r)
}

// ReadClusterSpec - Read tenants spec from io interface
// detects format and use either yaml or json parse
func ReadClusterSpec(b io.Reader) (*Cluster, error) {

	var spec Cluster

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	} else {
		fmt.Println(reflect.TypeOf(err).String())
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	} else {
		fmt.Println(err)
	}

	return nil, &InvalidClusterSpec{"unknown format"}
}

//InstanceSpecsFromString method return instance form string
func (c Cluster) InstanceSpecsFromString(s string) (interface{}, error) {
	return ClusterSpecsFromString(s)
}

// NewClusterSpecs create spec from reader
func NewClusterSpecs(r io.Reader) (*Cluster, error) {
	spec, err := ReadClusterSpec(r)
	if err != nil {
		return nil, err
	}

	return spec, nil
}
