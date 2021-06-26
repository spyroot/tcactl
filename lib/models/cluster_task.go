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

// TaskErrors error for give task
type TaskErrors struct {
	InternalMessage string `json:"internalMessage" yaml:"internal_message"`
	ErrorCode       string `json:"errorCode" yaml:"error_code"`
	Message         string `json:"message" yaml:"message"`
}

// TaskSteps Current execution or executed task list.
type TaskSteps struct {
	Title     string         `json:"title" yaml:"title"`
	Status    string         `json:"status" yaml:"status"`
	Progress  int            `json:"progress" yaml:"progress"`
	StartTime int64          `json:"startTime" yaml:"start_time"`
	EndTime   int64          `json:"endTime" yaml:"end_time"`
	Message   string         `json:"message" yaml:"message"`
	Children  *[]interface{} `json:"children" yaml:"children"`
	Response  struct {
		Warnings []interface{} `json:"warnings" yaml:"warnings"`
	} `json:"response,omitempty" yaml:"response"`
	Errors *[]TaskErrors `json:"errors,omitempty" yaml:"errors,omitempty"`
}

type TaskInterfaceInfo struct {
	Url                string `json:"url" yaml:"url"`
	Description        string `json:"description" yaml:"description"`
	TrustedCertificate string `json:"trustedCertificate" yaml:"trusted_certificate"`
}

type TaskOldEntry struct {
	ExtensionId   string             `json:"extensionId" yaml:"extension_id"`
	Name          string             `json:"name" yaml:"name"`
	Type          string             `json:"type" yaml:"type"`
	ExtensionKey  string             `json:"extensionKey" yaml:"extension_key"`
	Description   string             `json:"description" yaml:"description"`
	InterfaceInfo *TaskInterfaceInfo `json:"interfaceInfo" yaml:"interface_info"`
	AccessInfo    struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	} `json:"accessInfo" yaml:"access_info"`
	AdditionalParameters struct {
		TrustAllCerts bool `json:"trustAllCerts" yaml:"trust_all_certs"`
		RepoSyncVim   struct {
			VimName       string `json:"vimName" yaml:"vim_name"`
			VimId         string `json:"vimId" yaml:"vim_id"`
			VimSystemUUID string `json:"vimSystemUUID" yaml:"vim_system_uuid"`
		} `json:"repoSyncVim" yaml:"repo_sync_vim"`
	} `json:"additionalParameters" yaml:"additional_parameters"`
	State            string        `json:"state" yaml:"state"`
	ExtensionSubtype string        `json:"extensionSubtype" yaml:"extension_subtype"`
	Products         []interface{} `json:"products" yaml:"products"`
	VimInfo          []struct {
		VimName       string `json:"vimName" yaml:"vim_name"`
		VimId         string `json:"vimId" yaml:"vim_id"`
		VimSystemUUID string `json:"vimSystemUUID" yaml:"vim_system_uuid"`
	} `json:"vimInfo" yaml:"vim_info"`
	Version         string `json:"version" yaml:"version"`
	VnfCount        int    `json:"vnfCount" yaml:"vnf_count"`
	VnfCatalogCount int    `json:"vnfCatalogCount" yaml:"vnf_catalog_count"`
	Error           string `json:"error,omitempty" yaml:"error"`
}

// TaskTools tool provision for task
type TaskTools struct {
	Name       string `json:"name" yaml:"name"`
	Properties struct {
		Password string `json:"password" yaml:"password"`
		Type     string `json:"type" yaml:"type"`
		Url      string `json:"url" yaml:"url"`
		Username string `json:"username" yaml:"username"`
	} `json:"properties" yaml:"properties"`
}

// TaskPayload task payload
type TaskPayload struct {
	ClusterPassword   string `json:"clusterPassword" yaml:"cluster_password"`
	ClusterTemplateId string `json:"clusterTemplateId" yaml:"cluster_template_id"`
	ClusterType       string `json:"clusterType" yaml:"cluster_type"`
	ClusterConfig     struct {
		Csi []struct {
			Name       string `json:"name" yaml:"name"`
			Properties struct {
				ServerIP     string `json:"serverIP,omitempty" yaml:"server_ip"`
				MountPath    string `json:"mountPath,omitempty" yaml:"mount_path"`
				DatastoreUrl string `json:"datastoreUrl,omitempty" yaml:"datastore_url"`
			} `json:"properties" yaml:"properties"`
		} `json:"csi" yaml:"csi"`
		Tools []TaskTools `json:"tools" yaml:"tools"`
	} `json:"clusterConfig" yaml:"cluster_config"`
	HcxCloudUrl         string `json:"hcxCloudUrl" yaml:"hcx_cloud_url"`
	EndpointIP          string `json:"endpointIP" yaml:"endpoint_ip"`
	ManagementClusterId string `json:"managementClusterId" yaml:"management_cluster_id"`
	MasterNodes         []struct {
		Name     string `json:"name" yaml:"name"`
		Networks []struct {
			Label       string   `json:"label" yaml:"label"`
			NetworkName string   `json:"networkName" yaml:"networkName"`
			Nameservers []string `json:"nameservers" yaml:"nameservers"`
		} `json:"networks" yaml:"networks"`
		PlacementParams []struct {
			Name string `json:"name" yaml:"name"`
			Type string `json:"type" yaml:"type"`
		} `json:"placementParams" yaml:"placementParams"`
	} `json:"masterNodes" yaml:"master_nodes"`
	Name            string `json:"name" yaml:"name"`
	PlacementParams []struct {
		Name string `json:"name" yaml:"name"`
		Type string `json:"type" yaml:"type"`
	} `json:"placementParams" yaml:"placement_params"`
	VmTemplate  string `json:"vmTemplate" yaml:"vm_template"`
	WorkerNodes []struct {
		Name     string `json:"name" yaml:"name"`
		Networks []struct {
			Label       string   `json:"label" yaml:"label"`
			NetworkName string   `json:"networkName" yaml:"networkName"`
			Nameservers []string `json:"nameservers" yaml:"nameservers"`
		} `json:"networks" yaml:"networks"`
		PlacementParams []struct {
			Name string `json:"name" yaml:"name"`
			Type string `json:"type" yaml:"type"`
		} `json:"placementParams" yaml:"placementParams"`
		Id string `json:"id" yaml:"id"`
	} `json:"workerNodes" yaml:"worker_nodes"`
	Location `json:"location" yaml:"location"`
}

type TaskRepoSyncVim struct {
	VimName       string `json:"vimName" yaml:"vim_name"`
	VimId         string `json:"vimId" yaml:"vim_id"`
	VimSystemUUID string `json:"vimSystemUUID" yaml:"vim_system_uuid"`
}

type TaskAdditionalParameters struct {
	TrustAllCerts bool             `json:"trustAllCerts" yaml:"trust_all_certs"`
	RepoSyncVim   *TaskRepoSyncVim `json:"repoSyncVim,omitempty" yaml:"repo_sync_vim"`
}

type AccessInfo struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type TaskRequest struct {
	Name             string        `json:"name,omitempty" yaml:"name"`
	Version          string        `json:"version,omitempty" yaml:"version"`
	Type             string        `json:"type,omitempty" yaml:"type"`
	ExtensionKey     string        `json:"extensionKey,omitempty" yaml:"extension_key"`
	ExtensionSubtype string        `json:"extensionSubtype,omitempty" yaml:"extension_subtype"`
	Products         []interface{} `json:"products,omitempty" yaml:"products"`
	VimInfo          []struct {
		VimName       string `json:"vimName" yaml:"vim_name"`
		VimId         string `json:"vimId" yaml:"vim_id"`
		VimSystemUUID string `json:"vimSystemUUID" yaml:"vim_system_uuid"`
	} `json:"vimInfo,omitempty" yaml:"vimInfo"`
	InterfaceInfo        *TaskInterfaceInfo        `json:"interfaceInfo,omitempty" yaml:"interface_info"`
	AdditionalParameters *TaskAdditionalParameters `json:"additionalParameters,omitempty" yaml:"additional_parameters"`
	AutoScaleEnabled     bool                      `json:"autoScaleEnabled,omitempty" yaml:"auto_scale_enabled"`
	AutoHealEnabled      bool                      `json:"autoHealEnabled,omitempty" yaml:"auto_heal_enabled"`
	AccessInfo           *AccessInfo
	ExtensionId          string        `json:"extensionId,omitempty" yaml:"extension_id"`
	VimId                string        `json:"vimId,omitempty" yaml:"vim_id"`
	AssociatedVims       []string      `json:"associatedVims,omitempty" yaml:"associated_vims"`
	Id                   string        `json:"id" yaml:"id"`
	HcxUUID              string        `json:"hcxUUID,omitempty" yaml:"hcx_uuid"`
	ClusterName          string        `json:"clusterName,omitempty" yaml:"cluster_name"`
	IntentJobId          string        `json:"intentJobId,omitempty" yaml:"intent_job_id"`
	Intent               string        `json:"intent" yaml:"intent"`
	Request              *TaskRequest  `json:"request,omitempty" yaml:"request"`
	OldEntry             *TaskOldEntry `json:"oldEntry,omitempty" yaml:"old_entry"`
	DissociatedVims      []string      `json:"dissociatedVims,omitempty" yaml:"dissociated_vims"`
	RetainedVims         []interface{} `json:"retainedVims,omitempty" yaml:"retained_vims"`
	IsDelete             bool          `json:"isDelete,omitempty" yaml:"is_delete"`
	HcxCloudUrl          string        `json:"hcxCloudUrl,omitempty" yaml:"hcx_cloud_url"`
	VimUrl               string        `json:"vimUrl,omitempty" yaml:"vim_url"`
	ClusterType          string        `json:"clusterType,omitempty" yaml:"cluster_type"`
	Payload              *TaskPayload  `json:"payload,omitempty" yaml:"payload"`
}

type ClusterTask struct {
	Items []struct {
		EntityDetails struct {
			Id   string `json:"id" yaml:"id" yaml:"id"`
			Type string `json:"type" yaml:"type" yaml:"type"`
			Name string `json:"name" yaml:"name" yaml:"name"`
		} `json:"entityDetails" yaml:"entity_details"`
		Type      string       `json:"type" yaml:"type"`
		Status    string       `json:"status" yaml:"status"`
		Progress  int          `json:"progress" yaml:"progress"`
		Message   string       `json:"message" yaml:"message"`
		StartTime int64        `json:"startTime" yaml:"start_time"`
		EndTime   int64        `json:"endTime" yaml:"end_time"`
		Request   *TaskRequest `json:"request" yaml:"request"`
		Steps     []TaskSteps  `json:"steps" yaml:"steps"`
		TaskId    string       `json:"taskId" yaml:"task_id"`
	} `json:"items"`

	Paging struct {
		PageSize   int `json:"pageSize"`
		Offset     int `json:"offset"`
		TotalSize  int `json:"totalSize"`
		PageNumber int `json:"pageNumber"`
	} `json:"paging"`
}
