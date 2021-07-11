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
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
)

// GetVimTenants return list of all cloud provider
// attached to TCA
func (c *RestClient) GetVimTenants(ctx context.Context) (*response.Tenants, error) {

	glog.Infof("Retrieving vim tenants")

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).Get(c.BaseURL + apiTenants)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkError(resp)
	}

	var tenants response.Tenants
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		glog.Error("Failed parse server respond. %v", err)
		return nil, err
	}

	return &tenants, nil
}

// GetVim return vim attached to TCA
// vimId in format vmware_FB40D3DE2967483FBF9033B451DC7571
func (c *RestClient) GetVim(ctx context.Context, vimId string) (*response.TenantSpecs, error) {

	c.GetClient()
	apiReq := fmt.Sprintf(TcaVimTenant, vimId)
	glog.Info("Sending request to ", apiReq)
	resp, err := c.Client.R().SetContext(ctx).Get(c.BaseURL + apiReq)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkError(resp)
	}

	var tenants response.TenantSpecs
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &tenants, nil
}

// GetTenantsQuery returns list of cloud provider attached to TCA.
// tenant filter , filter result
func (c *RestClient) GetTenantsQuery(f *specs.TenantsNfFilter) (*response.Tenants, error) {

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).SetQueryString(apiTenantAction).Post(c.BaseURL + apiTenants)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkError(resp)
	}

	var tenants response.Tenants
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &tenants, nil
}

// RegisterCloudProvider method register new cloud provider.
func (c *RestClient) RegisterCloudProvider(r *specs.SpecCloudProvider) (*models.TcaTask, error) {

	glog.Infof("Cloud provider registration request %s", r.HcxCloudUrl)

	c.GetClient()
	resp, err := c.Client.R().SetBody(r).Post(c.BaseURL + apiTenants)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkErrors(resp)
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &task, nil
}

// DeleteTenant delete tenant cluster
func (c *RestClient) DeleteTenant(tenantCluster string) (*models.TcaTask, error) {

	glog.Infof("Deleting tenant cluster %v", tenantCluster)

	c.GetClient()
	resp, err := c.Client.R().Delete(c.BaseURL + fmt.Sprintf(TcaDeleteTenant, tenantCluster))
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkErrors(resp)
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &task, nil
}
