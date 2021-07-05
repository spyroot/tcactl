// Package request
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
package request

type ClusterTemplate struct {
	ClusterType   string `json:"clusterType" yaml:"name"`
	ClusterConfig struct {
		Cni []struct {
			Name       string `json:"name" yaml:"name"`
			Properties struct {
			} `json:"properties" yaml:"name"`
		} `json:"cni" yaml:"name"`
		Csi []struct {
			Name       string `json:"name" yaml:"name"`
			Properties struct {
				Name      string `json:"name" yaml:"name"`
				IsDefault bool   `json:"isDefault" yaml:"name"`
				Timeout   string `json:"timeout" yaml:"name"`
			} `json:"properties" yaml:"name"`
		} `json:"csi" yaml:"name"`
		KubernetesVersion string `json:"kubernetesVersion" yaml:"name"`
		Tools             []struct {
			Name    string `json:"name" yaml:"name"`
			Version string `json:"version" yaml:"name"`
		} `json:"tools" yaml:"name"`
	} `json:"clusterConfig" yaml:"name"`
	Description string `json:"description" yaml:"name"`
	MasterNodes []struct {
		Cpu      int    `json:"cpu" yaml:"name"`
		Memory   int    `json:"memory" yaml:"name"`
		Name     string `json:"name" yaml:"name"`
		Networks []struct {
			Label string `json:"label"`
		} `json:"networks"`
		Storage   int      `json:"storage"`
		Replica   int      `json:"replica"`
		Labels    []string `json:"labels"`
		CloneMode string   `json:"cloneMode"`
	} `json:"masterNodes"`
	Name string `json:"name"`
	Id   string `json:"id"`
	Tags []struct {
		AutoCreated bool   `json:"autoCreated"`
		Name        string `json:"name"`
	} `json:"tags"`
	WorkerNodes []struct {
		Cpu      int    `json:"cpu"`
		Memory   int    `json:"memory"`
		Name     string `json:"name"`
		Networks []struct {
			Label string `json:"label"`
		} `json:"networks"`
		Storage   int      `json:"storage"`
		Replica   int      `json:"replica"`
		Labels    []string `json:"labels"`
		CloneMode string   `json:"cloneMode"`
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
			HealthCheck struct {
				NodeStartupTimeout  string `json:"nodeStartupTimeout"`
				UnhealthyConditions []struct {
					Type    string `json:"type"`
					Status  string `json:"status"`
					Timeout string `json:"timeout"`
				} `json:"unhealthyConditions"`
			} `json:"healthCheck"`
		} `json:"config"`
	} `json:"workerNodes"`
}
