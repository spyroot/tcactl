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

// VmwareTypeFolder - type vmware folder
type VmwareTypeFolder string

// VmwareDatastore - type vmware datastore
type VmwareDatastore string

// VmwareResourcePool - type vmware resource pool
type VmwareResourcePool string

// ClusterComputeResource - type vmware compute resource
type ClusterComputeResource string

const (
	// TypeFolder vc folder
	TypeFolder VmwareTypeFolder = "Folder"

	// TypeDataStore vc datastore
	TypeDataStore VmwareDatastore = "Datastore"
	// TypeResourcePool vc resource pool

	TypeResourcePool VmwareResourcePool = "ResourcePool"

	// TypeClusterComputeResource vc cluster name
	TypeClusterComputeResource ClusterComputeResource = "IsValidClusterCompute"

	// TypeNetworkIdPrefix VC dvs port-group prefix
	TypeNetworkIdPrefix = "dvportgroup"

	// TypeNetworkVlan vc network type
	TypeNetworkVlan = "vlan"

	// TypeNetworkVirtualPortGroup vc dvs pg
	TypeNetworkVirtualPortGroup = "DistributedVirtualPortgroup"

	// EntityTypeFolder VMware VC contain entity folder
	EntityTypeFolder = "folder"

	// EntityTypeResourcePool - VMware VC resource pool
	EntityTypeResourcePool = "resourcePool"

	// EntityTypeCluster - VMware VC resource pool
	EntityTypeCluster = "cluster"

	// EntityTypeComputeResource - VMware VC compute resource
	EntityTypeComputeResource = "computeResource"

	// ChildVirtualTypeMachine VMware VC child element VM
	ChildVirtualTypeMachine = "VirtualMachine"

	// ChildTypeVirtualApp VMware VC child element VApp
	ChildTypeVirtualApp = "VirtualApp"
)

// ResourcePool VMware resource pool container view
type ResourcePool struct {
	Success   bool  `json:"success" yaml:"success"`
	Completed bool  `json:"completed" yaml:"completed"`
	Time      int64 `json:"time" yaml:"time"`
	Data      struct {
		Items []struct {
			EntityId          string `json:"entity_id" yaml:"entity_id"`
			VcenterInstanceId string `json:"vcenter_instanceId" yaml:"vcenter_instance_id"`
			Name              string `json:"name" yaml:"name"`
			Owner             struct {
				Type       string `json:"type" yaml:"type"`
				ServerGuid string `json:"serverGuid" yaml:"server_guid"`
				XsiType    string `json:"@xsi:type" yaml:"xsi_type"`
				Value      string `json:"value" yaml:"value"`
			} `json:"owner" yaml:"owner"`
			Parent struct {
				Type       string `json:"type" yaml:"type"`
				ServerGuid string `json:"serverGuid" yaml:"server_guid"`
				XsiType    string `json:"@xsi:type" yaml:"xsi_type"`
				Value      string `json:"value" yaml:"value"`
			} `json:"parent" yaml:"parent"`
			ResourcePool []struct {
				Type       string `json:"type" yaml:"type"`
				ServerGuid string `json:"serverGuid" yaml:"server_guid"`
				XsiType    string `json:"@xsi:type" yaml:"xsi_type"`
				Value      string `json:"value" yaml:"value"`
			} `json:"resourcePool,omitempty" yaml:"resource_pool"`
			EntityType string `json:"entityType" yaml:"entity_type"`
			Origin     struct {
				EndpointId   string `json:"endpointId" yaml:"endpoint_id"`
				EndpointType string `json:"endpointType" yaml:"endpoint_type"`
				EndpointName string `json:"endpointName" yaml:"endpoint_name"`
				ResourceId   string `json:"resourceId" yaml:"resource_id"`
				ResourceType string `json:"resourceType" yaml:"resource_type"`
				ResourceName string `json:"resourceName" yaml:"resource_name"`
			} `json:"_origin" yaml:"origin"`
			Source struct {
				Version    string `json:"version" yaml:"version"`
				Uuid       string `json:"uuid" yaml:"uuid"`
				HcspUUID   string `json:"hcspUUID" yaml:"hcsp_uuid"`
				SystemType string `json:"systemType" yaml:"system_type"`
			} `json:"_source" yaml:"source"`
		} `json:"items" yaml:"items"`
	} `json:"data" yaml:"data"`
}

// IsValidResource return if name is valid resource pool
func (r *ResourcePool) IsValidResource(name string) bool {

	if r == nil {
		return false
	}

	for _, it := range r.Data.Items {
		if it.Name == name && it.EntityType == EntityTypeResourcePool {
			return true
		}
	}

	return false
}

// Folders VMware VC folder container view.
type Folders struct {
	Success   bool  `json:"success" yaml:"success"`
	Completed bool  `json:"completed" yaml:"completed"`
	Time      int64 `json:"time" yaml:"time"`
	Data      struct {
		Items []struct {
			EntityId          string `json:"entity_id" yaml:"entity_id"`
			VcenterInstanceId string `json:"vcenter_instanceId" yaml:"vcenter_instance_id"`
			Name              string `json:"name" yaml:"name"`
			Parent            struct {
				Type       string `json:"type" yaml:"type"`
				ServerGuid string `json:"serverGuid" yaml:"server_guid"`
				XsiType    string `json:"@xsi:type" yaml:"xsi_type"`
				Value      string `json:"value" yaml:"value"`
			} `json:"parent" yaml:"parent"`
			EntityType       string `json:"entityType" yaml:"entity_type"`
			Id               string `json:"id" yaml:"id"`
			DisplayName      string `json:"displayName" yaml:"display_name"`
			ObjectAttributes struct {
				ChildType []string `json:"childType" yaml:"child_type"`
			} `json:"objectAttributes" yaml:"object_attributes"`
			Origin struct {
				EndpointId   string `json:"endpointId" yaml:"endpoint_id"`
				EndpointType string `json:"endpointType" yaml:"endpoint_type"`
				EndpointName string `json:"endpointName" yaml:"endpoint_name"`
				ResourceId   string `json:"resourceId" yaml:"resource_id"`
				ResourceType string `json:"resourceType" yaml:"resource_type"`
				ResourceName string `json:"resourceName" yaml:"resource_name"`
			} `json:"_origin" yaml:"origin"`
			Source struct {
				Version    string `json:"version" yaml:"version"`
				Uuid       string `json:"uuid" yaml:"uuid"`
				HcspUUID   string `json:"hcspUUID" yaml:"hcsp_uuid"`
				SystemType string `json:"systemType" yaml:"system_type"`
			} `json:"_source" yaml:"source"`
		} `json:"items" yaml:"items"`
	} `json:"data" yaml:"data"`
}

// IsValidFolder return true if name is valid folder
func (f *Folders) IsValidFolder(name string) bool {

	if f == nil {
		return false
	}

	for _, it := range f.Data.Items {
		if it.Name == name && it.EntityType == EntityTypeFolder {
			return true
		}
	}

	return false
}
