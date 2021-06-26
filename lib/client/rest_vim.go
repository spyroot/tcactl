// Package client
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
package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/models"
	"net/http"
)

const (
	// TcaVmwareClusters VMware clusters
	TcaVmwareClusters = "/hybridity/api/infra/inventory/vc/clusters"

	// TcaVmwareVmTemplates VMware VM templates
	TcaVmwareVmTemplates = "/hybridity/api/infra/inventory/vc/templates"

	// TcaVmwareNetworks VMware virtual networks
	TcaVmwareNetworks = "/hybridity/api/nfv/networks"

	// TcaVmwareVmContainers VMware container view
	TcaVmwareVmContainers = "/hybridity/api/service/inventory/containers"

	// TcaServiceVmwareVmContainers Vmware service container
	TcaServiceVmwareVmContainers = "/api/service/inventory/containers"
)

func (c *RestClient) GetVmwareCluster(f *request.ClusterFilterQuery) (*models.VMwareClusters, error) {

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).Post(c.BaseURL + TcaVmwareClusters)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var tenants models.VMwareClusters
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &tenants, nil
}

// GetVmwareNetworks - return query for vmware network list
func (c *RestClient) GetVmwareNetworks(f *request.VMwareNetworkQuery) (*models.CloudNetworks, error) {

	if f == nil {
		glog.Error("vmware network filter query is nil")
		return nil, fmt.Errorf("vmware network filter query is nil")
	}

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).Post(c.BaseURL + TcaVmwareNetworks)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var networks models.CloudNetworks
	if err := json.Unmarshal(resp.Body(), &networks); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &networks, nil
}

// GetVMwareTemplates - return VMware VM templates
// Typically Query filters based on cloud provider id.
func (c *RestClient) GetVMwareTemplates(f *request.VMwareTemplateQuery) (*models.VcInventory, error) {

	if f == nil {
		glog.Error("vm template filter query is nil")
		return nil, fmt.Errorf("vm template filter query is nil")
	}

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).Post(c.BaseURL + TcaVmwareVmTemplates)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var inventory models.VcInventory
	if err := json.Unmarshal(resp.Body(), &inventory); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &inventory, nil
}

// GetVMwareFolders - return VMware folder view
// Typically Query filters based on cloud provider id.
func (c *RestClient) GetVMwareFolders(f *request.VmwareFolderQuery) (*models.Folders, error) {

	if f == nil {
		glog.Error("folder filter query is nil")
		return nil, fmt.Errorf("folder filter query is nil")
	}

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).Post(c.BaseURL + TcaVmwareVmContainers)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var folder models.Folders
	if err := json.Unmarshal(resp.Body(), &folder); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &folder, nil
}

// GetVMwareResourcePool return VMware resource view
// Typically Query filters based on cloud provider id.
func (c *RestClient) GetVMwareResourcePool(f *request.VMwareResourcePoolQuery) (*models.ResourcePool, error) {

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).Post(c.BaseURL + TcaVmwareVmContainers)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var resourcePool models.ResourcePool
	if err := json.Unmarshal(resp.Body(), &resourcePool); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &resourcePool, nil
}

// GetVMwareInfraContainers - Call for VC cluster container
func (c *RestClient) GetVMwareInfraContainers(clusterId string) (*models.VmwareContainerView, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + TcaServiceVmwareVmContainers)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var vcInfra models.VmwareContainerView
	if err := json.Unmarshal(resp.Body(), &vcInfra); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &vcInfra, nil
}
