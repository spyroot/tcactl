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
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"strings"
)

// GetVim return vim
func (a *TcaApi) GetVim(ctx context.Context, NameOrVimId string) (*response.TenantSpecs, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	var (
		vimId = NameOrVimId
		err   error
	)

	// vim id is format vmware_numeric
	inputs := strings.Split(NameOrVimId, "_")

	if len(inputs) != 2 {
		// if we just string it a name
		if len(inputs) == 1 {
			vimId, err = a.ResolveVimId(ctx, NameOrVimId)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, &InvalidVimFormat{errMsg: NameOrVimId}
		}
	}

	glog.Infof("Retrieving vim specString vim id %s", vimId)

	return a.rest.GetVim(ctx, vimId)
}

// GetVimComputeClusters - return compute cluster attached to VIM
// For example VMware VIM is vCenter.
func (a *TcaApi) GetVimComputeClusters(ctx context.Context, cloudName string) (*models.VMwareClusters, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(cloudName) == 0 {
		return nil, errors.New("empty cloud provider")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
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

		f := specs.NewClusterFilterQuery(tenant.HcxUUID)
		clusterInventory, err := a.rest.GetVmwareCluster(ctx, f)

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

	glog.Infof("Retrieving network list for cloud provider uuid %s, %s",
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
