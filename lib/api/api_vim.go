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
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/api_errors"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"strings"
)

// GetVim return cloud provider.
func (a *TcaApi) GetVim(ctx context.Context, NameOrId string) (*response.TenantSpecs, error) {

	var (
		providerId = NameOrId
		err        error
	)

	glog.Infof("Retrieving vim specString vim id %s", providerId)

	if len(NameOrId) == 0 {
		return nil, api_errors.NewInvalidSpec("empty cloud provider id or name")
	}

	// vim id is format vmware_numeric
	inputs := strings.Split(NameOrId, "_")

	if len(inputs) != 2 {
		// if we just string it a name
		if len(inputs) == 1 {
			providerId, err = a.ResolveVimId(ctx, NameOrId)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, api_errors.NewInvalidSpec(NameOrId)
		}
	}

	return a.rest.GetVim(ctx, providerId)
}

// GetConsumption - return lic consumption
func (a *TcaApi) GetConsumption(ctx context.Context) (*models.ConsumptionResp, error) {

	consumptionResp, err := a.rest.GetConsumption(ctx)
	if err != nil {
		return nil, err
	}
	return consumptionResp, nil

}

// GetVimComputeClusters - return compute cluster attached to VIM
// For example VMware VIM is vCenter.
func (a *TcaApi) GetVimComputeClusters(ctx context.Context, cloudName string) (*models.VMwareClusters, error) {

	if len(cloudName) == 0 {
		return nil, api_errors.NewInvalidArgument("cloudName is empty")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	if !tenant.IsVMware() {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	glog.Infof("Fetching list for cloud provider %v '%v'", tenant.HcxUUID, tenant.VimURL)
	f := specs.NewClusterFilterQuery(tenant.HcxUUID)
	clusterInventory, err := a.rest.GetVmwareCluster(ctx, f)
	if err != nil {
		return nil, err
	}

	return clusterInventory, nil
}

// GetVimNetworks - method return network attached
// to vim, cloud provider, for a VMware it full path to object
func (a *TcaApi) GetVimNetworks(ctx context.Context, cloudName string) (*models.CloudNetworks, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	glog.Infof("Fetching network list for cloud provider uuid %s, %s",
		tenant.HcxUUID, tenant.VimURL)

	if !tenant.IsVMware() {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	f := specs.NewClusterFilterQuery(tenant.HcxUUID)
	clusterInventory, err := a.rest.GetVmwareCluster(ctx, f)
	if err != nil {
		return nil, err
	}

	// get all network for all clusters
	var networks models.CloudNetworks
	for _, item := range clusterInventory.Items {
		networkFilter := specs.VMwareNetworkQuery{}
		networkFilter.Filter.TenantId = tenant.HcxUUID
		if strings.HasPrefix(item.EntityId, "domain") {
			networkFilter.Filter.ClusterId = item.EntityId
			net, err := a.rest.GetVmwareNetworks(ctx, &networkFilter)
			if err != nil {
				return nil, err
			}

			networks.Network = append(networks.Network, net.Network...)
		}
	}

	return &networks, nil
}
