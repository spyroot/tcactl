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
func (t *ClusterSpec) GetField(field string) string {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	return f.String()
}

func (t *ClusterSpec) GetFields() (map[string]interface{}, error) {

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

type ClusterNotFound struct {
	ErrMsg string
}

func (m *ClusterNotFound) Error() string {
	return "cluster '" + m.ErrMsg + "' not found"
}

// GetClusterId return cluster ID
func (c *Clusters) GetClusterId(q string) (string, error) {

	if c == nil {
		return "", fmt.Errorf("uninitialized object")
	}

	for _, it := range c.Clusters {
		if it.ClusterName == q || it.Id == q {
			glog.Infof("Found cluster '%v' cluster uuid '%v'", q, it.Id)
			return it.Id, nil
		}
	}

	return "", &ClusterNotFound{ErrMsg: q}
}

// GetClusterSpec return single cluster information,
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
