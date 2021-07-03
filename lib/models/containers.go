// Package models
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
package models

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/api_errors"
	"strings"
)

// VmwareContainerView - tca encapsulates all VMware view to very large json
type VmwareContainerView struct {
	Success   bool  `json:"success" yaml:"success"`
	Completed bool  `json:"completed" yaml:"completed"`
	Time      int64 `json:"time" yaml:"time"`
	Data      struct {
		Items []struct {
			EntityId          string `json:"entity_id" yaml:"entity_id"`
			VcenterInstanceId string `json:"vcenter_instanceId" yaml:"vcenter_instanceId"`
			Name              string `json:"name" yaml:"name"`
			Owner             struct {
				Type       string `json:"type" yaml:"type"`
				ServerGuid string `json:"serverGuid" yaml:"serverGuid"`
				XsiType    string `json:"@xsi:type" yaml:"@xsi:type"`
				Value      string `json:"value" yaml:"value"`
			} `json:"owner" yaml:"owner"`
			Parent struct {
				Type       string `json:"type" yaml:"type"`
				ServerGuid string `json:"serverGuid" yaml:"serverGuid"`
				XsiType    string `json:"@xsi:type" yaml:"@xsi:type"`
				Value      string `json:"value" yaml:"value"`
			} `json:"parent" yaml:"parent"`
			EntityType string `json:"entityType" yaml:"entityType"`
			Origin     struct {
				EndpointId   string `json:"endpointId" yaml:"endpointId"`
				EndpointType string `json:"endpointType" yaml:"endpointType"`
				EndpointName string `json:"endpointName" yaml:"endpointName"`
				ResourceId   string `json:"resourceId" yaml:"resourceId"`
				ResourceType string `json:"resourceType" yaml:"resourceType"`
				ResourceName string `json:"resourceName" yaml:"resourceName"`
			} `json:"_origin" yaml:"_origin"`
			Source struct {
				Version    string `json:"version" yaml:"version"`
				Uuid       string `json:"uuid" yaml:"uuid"`
				HcspUUID   string `json:"hcspUUID" yaml:"hcspUUID"`
				SystemType string `json:"systemType" yaml:"systemType"`
			} `json:"_source" yaml:"_source"`
			ResourcePool []struct {
				Type       string `json:"type" yaml:"type"`
				ServerGuid string `json:"serverGuid" yaml:"serverGuid"`
				XsiType    string `json:"@xsi:type" yaml:"@xsi:type"`
				Value      string `json:"value" yaml:"value"`
			} `json:"resourcePool,omitempty" yaml:"resourcePool"`
		} `json:"items" yaml:"items"`
	} `json:"data" yaml:"data"`
}

type VmwareDatastoreSpec struct {
	EntityId string `json:"entity_id" yaml:"entity_id"`
	Name     string `json:"name" yaml:"name"`
	Summary  struct {
		Accessible         string `json:"accessible" yaml:"accessible"`
		Capacity           int64  `json:"capacity" yaml:"capacity"`
		FreeSpace          int64  `json:"freeSpace" yaml:"free_space"`
		MaintenanceMode    string `json:"maintenanceMode" yaml:"maintenanceMode"`
		MultipleHostAccess string `json:"multipleHostAccess" yaml:"multipleHostAccess"`
		Type               string `json:"type" yaml:"type"`
		Url                string `json:"url" yaml:"url"`
		Uncommitted        int64  `json:"uncommitted,omitempty" yaml:"uncommitted"`
	}
}

// VMwareClusters Vmware Clusters Container
type VMwareClusters struct {
	Items []struct {
		EntityId                      string                `json:"entity_id" yaml:"entity_id"`
		Name                          string                `json:"name" yaml:"name"`
		EntityType                    string                `json:"entityType" yaml:"entity_type"`
		NumOfHosts                    int                   `json:"numOfHosts" yaml:"num_of_hosts"`
		Datastore                     []VmwareDatastoreSpec `json:"datastore" yaml:"datastore"`
		Memory                        int64                 `json:"memory" yaml:"memory"`
		Cpu                           int                   `json:"cpu" yaml:"cpu"`
		K8ClusterDeployed             int                   `json:"k8ClusterDeployed" yaml:"k8ClusterDeployed"`
		NumK8SMgmtClusterDeployed     int                   `json:"numK8sMgmtClusterDeployed" yaml:"numK8sMgmtClusterDeployed"`
		NumK8SWorkloadClusterDeployed int                   `json:"numK8sWorkloadClusterDeployed" yaml:"numK8sWorkloadClusterDeployed"`
	} `json:"items"`
}

// IsValidClusterCompute return true if name exists
func (c *VMwareClusters) IsValidClusterCompute(name string) bool {

	if c == nil {
		return false
	}

	for _, item := range c.Items {
		if item.EntityId == EntityTypeCluster && item.Name == name {
			return true
		}
	}

	return false
}

// ResourcePool return true if resource pool exists
func (c *VMwareClusters) ResourcePool(name string) bool {
	if c == nil {
		return false
	}

	for _, item := range c.Items {
		if item.EntityId == "cluster" && item.Name == name {
			glog.Infof("Found cluster %v", item.Name)
			return true
		}
	}

	return false
}

// IsValidDatastore return true if datastore exists
func (c *VMwareClusters) IsValidDatastore(name string) bool {

	if c == nil {
		return false
	}

	for _, it := range c.Items {
		for _, ds := range it.Datastore {
			if ds.Name == name {
				glog.Infof("Found target ds %s %s", ds.Summary.Url, name)
				return true
			}
		}
	}

	return false
}

// GetDatastoreByUrl find datastore VmwareDatastoreSpec based on full datastore ur
// ds:///vmfs/volumes/60254cec-cb11aa98-a02e-e4434bf999aa/
func (c *VMwareClusters) GetDatastoreByUrl(url string) (*VmwareDatastoreSpec, error) {

	if c == nil {
		return nil, fmt.Errorf("nil instance")
	}

	for _, it := range c.Items {
		for _, ds := range it.Datastore {
			if strings.ToLower(ds.Summary.Url) == strings.ToLower(url) {
				glog.Infof("Found target ds uri %s %s", ds.Summary.Url, url)
				return &ds, nil
			}
		}
	}

	return nil, api_errors.NewDatastoreNotFound(url)
}

// IsValidDatastoreUrl return true if datastore exists
func (c *VMwareClusters) IsValidDatastoreUrl(url string) bool {

	if c == nil {
		return false
	}

	ds, err := c.GetDatastoreByUrl(url)
	if err != nil {
		return false
	}

	return len(ds.Name) > 0
}
