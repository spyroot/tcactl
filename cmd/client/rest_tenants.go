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
	"github.com/spyroot/tcactl/cmd/client/request"
	"github.com/spyroot/tcactl/cmd/client/response"
	"net/http"
)

const (
	// apiTenants - list of vim tenant rest call
	apiTenants = "/hybridity/api/vims/v1/tenants"

	// apiVim - attached vim list rest call
	apiVim = "/hybridity/api/vims/v1/"

	// apiTenantAction query action
	apiTenantAction = "action=query"
)

// GetVimTenants return list of all cloud provider attached
// to TCA
func (c *RestClient) GetVimTenants() (*response.Tenants, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + apiTenants)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.dumpRespond {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var tenants response.Tenants
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &tenants, nil
}

// GetVim return list of cloud provider attached to TCA
func (c *RestClient) GetVim(vimId string) (*response.TenantSpecs, error) {

	c.GetClient()
	// format vmware_FB40D3DE2967483FBF9033B451DC7571
	glog.Info("Sending request to ", c.BaseURL+apiTenants+"/"+vimId+"/tenants")
	resp, err := c.Client.R().Get(c.BaseURL + apiVim + "/" + vimId + "/tenants")
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	fmt.Println(resp.StatusCode())

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
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
func (c *RestClient) GetTenantsQuery(f *request.TenantsNfFilter) (*response.Tenants, error) {

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).SetQueryString(apiTenantAction).
		Post(c.BaseURL + apiTenants)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var tenants response.Tenants
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &tenants, nil
}

// DeleteTenant delete tenant cluster
func (c *RestClient) DeleteTenant(tenantCluster string) (bool, error) {

	c.GetClient()
	resp, err := c.Client.R().Delete(c.BaseURL + "/hybridity/api/vims/v1/tenants/" + tenantCluster)
	if err != nil {
		glog.Error(err)
		return false, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return false, c.checkErrors(resp)
	}

	var tenants response.Tenants
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		glog.Error("Failed parse server respond.")
		return false, err
	}

	return resp.StatusCode() == http.StatusOK, nil
}
