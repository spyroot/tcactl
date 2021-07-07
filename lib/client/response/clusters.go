// Package respons
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

package response

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/netutils"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"strings"
)

type ClusterSuccess struct {
	Id          string `json:"id"`
	OperationId string `json:"operationId"`
}

// ClusterSpecTemplate template id/name.
type ClusterSpecTemplate struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Id      string `json:"id"`
}

type ClusterNetwork struct {
	Label       string   `json:"label"`
	NetworkName string   `json:"networkName"`
	Nameservers []string `json:"nameservers"`
}

type MasterNodesDetails struct {
	Cpu       int              `json:"cpu"`
	Memory    int              `json:"memory"`
	Name      string           `json:"name"`
	Networks  []ClusterNetwork `json:"networks"`
	Storage   int              `json:"storage"`
	Replica   int              `json:"replica"`
	Labels    []string         `json:"labels"`
	CloneMode string           `json:"cloneMode"`
}

type WorkerNodesDetails struct {
	Cpu       int               `json:"cpu"`
	Memory    int               `json:"memory"`
	Name      string            `json:"name"`
	Networks  []models.Networks `json:"networks"`
	Storage   int               `json:"storage"`
	Replica   int               `json:"replica"`
	Labels    []string          `json:"labels"`
	CloneMode string            `json:"cloneMode"`
	Config    struct {
		CpuManagerPolicy struct {
			Type       string `json:"type"`
			Policy     string `json:"policy"`
			Properties struct {
				KubeReserved struct {
					Cpu         int `json:"cpu"`
					MemoryInGiB int `json:"memoryInGiB"`
				} `json:"kubeReserved"`
				SystemReserved struct {
					Cpu         int `json:"cpu"`
					MemoryInGiB int `json:"memoryInGiB"`
				} `json:"systemReserved"`
			} `json:"properties"`
		} `json:"cpuManagerPolicy"`
		HealthCheck *models.HealthCheck `json:"healthCheck"`
	} `json:"config"`
}

// ClusterSpec - hold cluster specs
type ClusterSpec struct {
	// Id cluster ID user internal
	Id string `json:"id"`
	// Cluster name
	ClusterName string `json:"clusterName"`
	// Cluster Type VC or Kube
	ClusterType         string               `json:"clusterType"`
	VsphereClusterName  string               `json:"vsphereClusterName"`
	ManagementClusterId string               `json:"managementClusterId"`
	HcxUUID             string               `json:"hcxUUID"`
	Status              string               `json:"status"`
	ActiveTasksCount    int                  `json:"activeTasksCount"`
	ClusterTemplate     *ClusterSpecTemplate `json:"clusterTemplate"`
	ClusterId           string               `json:"clusterId"`
	ClusterUrl          string               `json:"clusterUrl"`
	KubeConfig          string               `json:"kubeConfig"`
	EndpointIP          string               `json:"endpointIP"`
	MasterNodes         []MasterNodesDetails `json:"masterNodes"`
	WorkerNodes         []WorkerNodesDetails `json:"workerNodes"`
	VimId               string               `json:"vimId"`
	Error               string               `json:"error"`
}

// Clusters - a list of all clusters
type Clusters struct {
	Clusters []ClusterSpec
}

// GetField - return field from Cluster Spec struct
func (c *ClusterSpec) GetField(field string) string {

	r := reflect.ValueOf(c)
	v := reflect.Indirect(r)

	if v.IsValid() {
		f := v.FieldByName(strings.Title(field))
		if f.IsValid() {
			k := f.Kind()
			if k == reflect.Int {
				return k.String()
			}
			if k == reflect.String {
				return f.String()
			}
		} else {
			return ""
		}
	}

	return ""
}

func (c *ClusterSpec) GetFields() (map[string]interface{}, error) {

	var m map[string]interface{}

	b, err := json.Marshal(c)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}

	return m, nil
}

type ClusterNotFound struct {
	ErrMsg string
}

func (m *ClusterNotFound) Error() string {
	return "cluster '" + m.ErrMsg + "' not found"
}

// GetClusterId return cluster ID
func (c *Clusters) GetClusterId(NameOrId string) (string, error) {

	if c == nil {
		return "", fmt.Errorf("uninitialized object")
	}

	for _, it := range c.Clusters {
		if it.ClusterName == NameOrId || it.Id == NameOrId {
			glog.Infof("Found cluster '%v' cluster uuid '%v'", NameOrId, it.Id)
			return it.Id, nil
		}
	}

	return "", &ClusterNotFound{ErrMsg: NameOrId}
}

// GetClusterSpec return cluster information,
// loop up up by name or id, if not found return error
func (c *Clusters) GetClusterSpec(q string) (*ClusterSpec, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized clusters object")
	}

	for _, it := range c.Clusters {
		if it.ClusterName == q || it.Id == q {
			glog.Infof("Found cluster %v cluster id %v", q, it.Id)
			return &it, nil
		}
	}

	return nil, &ClusterNotFound{ErrMsg: q}
}

// GetClusterIds -return list of all cluster ids
func (c *Clusters) GetClusterIds() ([]string, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized clusters object")
	}

	var ids []string
	for _, cluster := range c.Clusters {
		if len(cluster.Id) > 0 {
			ids = append(ids, cluster.Id)
		}
	}
	return ids, nil
}

// ClustersSpecsFromFile - reads tenant spec from file
// and return ClusterSpec instance
func ClustersSpecsFromFile(fileName string) (*Clusters, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadClustersSpec(file)
}

// ClustersSpecsFromString take string that hold entire spec
// passed to reader and return ClusterSpec instance
func ClustersSpecsFromString(str string) (*Clusters, error) {
	r := strings.NewReader(str)
	return ReadClustersSpec(r)
}

// ReadClustersSpec - Read cluster spec from io interface
// detects format and use either yaml or json parse
func ReadClustersSpec(b io.Reader) (*Clusters, error) {

	var spec Clusters

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

	return nil, &InvalidClusterSpec{"unknown format"}
}

//InstanceSpecsFromString method return instance form string
func (c Clusters) InstanceSpecsFromString(s string) (interface{}, error) {
	return ClustersSpecsFromString(s)
}

type ClusterEndpoint struct {
	Cluster string
	IsIP    bool
}

// GetClusterIPs return all cluster IP
func (c *Clusters) GetClusterIPs() map[string]ClusterEndpoint {

	m := make(map[string]ClusterEndpoint)

	for _, cluster := range c.Clusters {
		u := cluster.ClusterUrl
		if strings.HasPrefix(u, "https://") {
			u = strings.TrimPrefix(u, "https://")
		}

		if strings.HasSuffix(u, ":6443") {
			u = strings.TrimSuffix(u, ":6443")
		}

		if netutils.IsDNSName(u) {
			m[u] = ClusterEndpoint{u, false}
		} else {
			ip := net.ParseIP(u)
			if ip != nil {
				m[u] = ClusterEndpoint{ip.String(), true}
			}
		}
	}

	return m
}

// NewClustersSpecs create cluster
// spec from reader
func NewClustersSpecs(r io.Reader) (*Clusters, error) {

	spec, err := ReadClustersSpec(r)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// ClusterSpecsFromFile - reads tenant cluster spec from file
// and return ClusterSpec instance
func ClusterSpecsFromFile(fileName string) (*ClusterSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadClusterSpec(file)
}

// ClusterSpecsFromString take string that hold entire spec
// passed to reader and return ClusterSpec instance
func ClusterSpecsFromString(str string) (*ClusterSpec, error) {
	r := strings.NewReader(str)
	return ReadClusterSpec(r)
}

// ReadClusterSpec - Read cluster spec from io interface
// detects format and use either yaml or json parse
func ReadClusterSpec(b io.Reader) (*ClusterSpec, error) {

	var spec ClusterSpec

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

	return nil, &InvalidClusterSpec{"unknown format"}
}

//InstanceSpecsFromString method return instance form string
func (c ClusterSpec) InstanceSpecsFromString(s string) (interface{}, error) {
	return ClusterSpecsFromString(s)
}

// NewClusterSpecs create spec from reader
func NewClusterSpecs(r io.Reader) (*ClusterSpec, error) {
	spec, err := ReadClusterSpec(r)
	if err != nil {
		return nil, err
	}

	return spec, nil
}
