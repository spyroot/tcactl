package request

import (
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/models"
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
			datastoreName string `json:"datastoreName,omitempty" yaml:"datastoreName,omitempty"`
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
