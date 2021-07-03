// Package app
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

package api

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/models"
	"strings"
)

// GetVimComputeClusters - return compute cluster attached to VIM
// For example VMware VIM is vCenter.
func (a *TcaApi) GetVimComputeClusters(cloudName string) (*models.VMwareClusters, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(cloudName) == 0 {
		return nil, errors.New("empty cloud provider")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	if tenant.IsVMware() {
		//
		glog.Infof("Retrieving list for cloud provider %v '%v'",
			tenant.HcxUUID, tenant.VimURL)

		f := request.NewClusterFilterQuery(tenant.HcxUUID)
		clusterInventory, err := a.rest.GetVmwareCluster(f)

		if err != nil {
			return nil, err
		}

		return clusterInventory, nil
	} else {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	return nil, &CloudProviderNotFound{errMsg: cloudName}
}

// GetVimNetworks - method return network attached
// to vim, cloud provider, for a VMware it full path to object
func (a *TcaApi) GetVimNetworks(cloudName string) (*models.CloudNetworks, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	glog.Infof("Retrieving network list for cloud provider uuid %s, %s",
		tenant.HcxUUID, tenant.VimURL)

	if !tenant.IsVMware() {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	f := request.NewClusterFilterQuery(tenant.HcxUUID)
	clusterInventory, err := a.rest.GetVmwareCluster(f)
	if err != nil {
		return nil, err
	}

	// get all network for all clusters
	var networks models.CloudNetworks
	for _, item := range clusterInventory.Items {
		networkFilter := request.VMwareNetworkQuery{}
		networkFilter.Filter.TenantId = tenant.HcxUUID
		if strings.HasPrefix(item.EntityId, "domain") {
			networkFilter.Filter.ClusterId = item.EntityId
			net, err := a.rest.GetVmwareNetworks(&networkFilter)
			if err != nil {
				return nil, err
			}

			networks.Network = append(networks.Network, net.Network...)
		}
	}

	return &networks, nil
}

func (a *TcaApi) DeleteCloudProvider(s string) (*models.TcaTask, error) {

	vims, err := a.GetVims()
	if err != nil {
		return nil, err
	}

	provider, err := vims.FindCloudProvider(s)
	if err != nil {
		return nil, err
	}

	task, err := a.rest.DeleteTenant(provider.ID)
	if err != nil {
		return nil, err
	}

	return task, nil
}
