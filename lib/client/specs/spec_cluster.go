// Package specs
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Mustafa mbayramo@vmware.com
package specs

import (
	"github.com/asaskevich/govalidator"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/models"
	"io"
	"net"
	"strings"
)

type ClusterType string

const (
	// ClusterManagement management k8s cluster
	ClusterManagement ClusterType = "MANAGEMENT"

	// ClusterWorkload workload k8s cluster
	ClusterWorkload ClusterType = "WORKLOAD"
)

// InvalidClusterSpec error if specs invalid
type InvalidClusterSpec struct {
	errMsg string
}

//
func (m *InvalidClusterSpec) Error() string {
	return m.errMsg
}

type ClusterConfigProperties struct {
	ServerIP      string `json:"serverIP,omitempty" yaml:"serverIP,omitempty"`
	MountPath     string `json:"mountPath,omitempty" yaml:"mountPath,omitempty"`
	DatastoreUrl  string `json:"datastoreUrl,omitempty" yaml:"datastoreUrl,omitempty"`
	DatastoreName string `json:"datastoreName,omitempty" yaml:"datastoreName,omitempty"`
}
type CsiConfig struct {
	Name       string                   `json:"name,omitempty" yaml:"name,omitempty"`
	Properties *ClusterConfigProperties `json:"properties,omitempty" yaml:"properties,omitempty"`
}

func (c CsiConfig) IsVsphereCsi() bool {
	return strings.ToLower(c.Name) == "vsphere-csi"
}

func (c CsiConfig) IsNfsCsi() bool {
	return strings.ToLower(c.Name) == "nfs_client"
}

type ClusterConfigToolProperties struct {
	ExtensionId string `json:"extensionId,omitempty" yaml:"extensionId,omitempty"`
	Password    string `json:"password,omitempty" yaml:"password,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Url         string `json:"url,omitempty" yaml:"url,omitempty"`
	Username    string `json:"username,omitempty" yaml:"username,omitempty"`
}

type ClusterConfigTools struct {
	Name       string                       `json:"name,omitempty" yaml:"name,omitempty"`
	Properties *ClusterConfigToolProperties `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type ClusterConfig struct {
	Csi            []CsiConfig          `json:"csi,omitempty" yaml:"csi,omitempty"`
	Tools          []ClusterConfigTools `json:"tools,omitempty" yaml:"tools,omitempty"`
	SystemSettings []struct {
		Name       string `json:"name,omitempty" yaml:"name,omitempty"`
		Properties struct {
			Host     string `json:"host,omitempty" yaml:"host,omitempty"`
			Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
			Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
		} `json:"properties" yaml:"properties,omitempty"`
	} `json:"systemSettings,omitempty" yaml:"systemSettings,omitempty"`
}

// SpecCluster new cluster creation request
type SpecCluster struct {
	SpecType            SpecType                 `json:"kind,omitempty" yaml:"kind,omitempty" valid:"required~kind is mandatory spec field"`
	Name                string                   `json:"name" yaml:"name" valid:"required~name is mandatory spec field"`
	ClusterPassword     string                   `json:"clusterPassword" yaml:"clusterPassword" valid:"required~clusterPassword is mandatory spec field"`
	ClusterTemplateId   string                   `json:"clusterTemplateId" yaml:"clusterTemplateId" valid:"required~clusterTemplateId is mandatory spec field"`
	ClusterType         string                   `json:"clusterType" yaml:"clusterType" valid:"required~clusterType is mandatory spec field"`
	Description         string                   `json:"description,omitempty" yaml:"description,omitempty"`
	Location            *models.Location         `json:"location,omitempty" yaml:"location,omitempty"`
	ClusterConfig       *ClusterConfig           `json:"clusterConfig,omitempty" yaml:"clusterConfig,omitempty"`
	HcxCloudUrl         string                   `json:"hcxCloudUrl" yaml:"hcxCloudUrl" valid:"required~vmTemplate is mandatory spec field"`
	EndpointIP          string                   `json:"endpointIP" yaml:"endpointIP" valid:"required~vmTemplate is mandatory spec field"`
	ManagementClusterId string                   `json:"managementClusterId,omitempty" yaml:"managementClusterId,omitempty"`
	VmTemplate          string                   `json:"vmTemplate" yaml:"vmTemplate" valid:"required~vmTemplate is mandatory spec field"`
	MasterNodes         []models.TypeNode        `json:"masterNodes" yaml:"masterNodes" valid:"required"`
	WorkerNodes         []models.TypeNode        `json:"workerNodes" yaml:"workerNodes" valid:"required"`
	PlacementParams     []models.PlacementParams `json:"placementParams" yaml:"placementParams" valid:"required"`

	specError error
}

//Validate method validate specs
func (c *SpecCluster) Validate() error {

	if c == nil {
		return &InvalidClusterSpec{errMsg: "nil instance"}
	}

	if c.Kind() != SpecKindCluster {
		return &InvalidClusterSpec{errMsg: "spec must contain kind field"}
	}

	ip := net.ParseIP(c.EndpointIP)
	if ip == nil {
		return &InvalidClusterSpec{errMsg: "invalid endpoint ip"}
	}

	if len(c.MasterNodes) == 0 {
		return &InvalidClusterSpec{errMsg: "cluster spec must include master node spec"}
	}
	if len(c.WorkerNodes) == 0 {
		return &InvalidClusterSpec{errMsg: "cluster spec must include worker node spec"}
	}
	if len(c.PlacementParams) == 0 {
		return &InvalidClusterSpec{errMsg: "cluster spec must include placement spec"}
	}

	for _, node := range c.MasterNodes {
		if len(node.PlacementParams) == 0 {
			return &InvalidClusterSpec{errMsg: "cluster spec must include placement spec for master node"}
		}
		if len(node.Name) == 0 {
			return &InvalidClusterSpec{errMsg: "cluster spec must node pool name"}
		}
		if len(node.Networks) == 0 {
			return &InvalidClusterSpec{errMsg: "cluster spec worker node section must contain networks:"}
		}
		// check al network section
		for _, network := range node.Networks {
			if len(network.Label) == 0 {
				return &InvalidClusterSpec{errMsg: "cluster spec networks section has no label"}
			}
			if len(network.NetworkName) == 0 {
				return &InvalidClusterSpec{errMsg: "cluster spec networks section has no network name"}
			}
		}

	}

	for _, node := range c.WorkerNodes {
		if len(node.PlacementParams) == 0 {
			return &InvalidClusterSpec{errMsg: "cluster spec must include placement spec for worker node"}
		}
		if len(node.Name) == 0 {
			return &InvalidClusterSpec{errMsg: "cluster spec must include field name:node_pool_name"}
		}
		if len(node.Networks) == 0 {
			return &InvalidClusterSpec{errMsg: "cluster spec worker node section must contain networks:"}
		}
		// check al network section
		for _, network := range node.Networks {
			if len(network.Label) == 0 {
				return &InvalidClusterSpec{errMsg: "cluster spec networks section has no label"}
			}
			if len(network.NetworkName) == 0 {
				return &InvalidClusterSpec{errMsg: "cluster spec networks section has no network name"}
			}
		}
	}

	if c.IsWorkload() && len(c.ManagementClusterId) == 0 {
		return &InvalidClusterSpec{errMsg: "workload cluster spec must include ManagementClusterId:"}
	}

	if c.IsWorkload() && c.ClusterConfig == nil {
		return &InvalidClusterSpec{errMsg: "workload cluster spec must contain clusterConfig: section"}
	}

	if c.IsWorkload() && c.ClusterConfig != nil {
		for _, config := range c.ClusterConfig.Csi {
			if config.IsNfsCsi() != true && config.IsVsphereCsi() != true {
				return &InvalidClusterSpec{errMsg: "csi name must be either nfs_client or vsphere_csi error on value" + config.Name}
			}
		}
	}

	if c.IsManagement() == false && c.IsWorkload() == false {
		return &InvalidClusterSpec{errMsg: "cluster type must be either workload or management"}
	}
	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return err
	}

	return nil
}

// IsManagement return true if spec if management cluster spec
func (c *SpecCluster) IsManagement() bool {

	if c == nil {
		return false
	}

	return strings.ToLower(c.ClusterType) == strings.ToLower(string(ClusterManagement))
}

// IsWorkload return if spec is for workload cluster
func (c *SpecCluster) IsWorkload() bool {

	if c == nil {
		return false
	}

	return strings.ToLower(c.ClusterType) == strings.ToLower(string(ClusterWorkload))
}

//FindNodePoolByName search for node pool name
// if isWorker will check worker node pool, otherwise Master node pools.
func (c *SpecCluster) FindNodePoolByName(name string, isWorker bool) bool {

	if c == nil {
		return false
	}

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

//SpecsFromString method read cluster spec from string
//and return instance
func (c SpecCluster) SpecsFromString(s string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromString(s, new(SpecCluster), f...)
}

//SpecsFromFile method return instance form string
func (c SpecCluster) SpecsFromFile(fileName string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromFile(fileName, new(SpecCluster), f...)
}

// SpecsFromReader create spec from reader
func (c SpecCluster) SpecsFromReader(r io.Reader, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpec(r, new(SpecCluster), f...)
}

func (c *SpecCluster) Kind() SpecType {
	return c.SpecType
}

// IsValid return false if validator set error
func (c *SpecCluster) IsValid() bool {
	if c.specError != nil {
		return false
	}
	return true
}

// Default TODO
func (c *SpecCluster) Default() error {
	return nil
}
