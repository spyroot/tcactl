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
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/client/response"
	"net/http"
)

const (
	apiTenants = "/hybridity/api/vims/v1/tenants"
	apiVim     = "/hybridity/api/vims/v1/"
	///hybridity/api/vims/v1/vmware_FB40D3DE2967483FBF9033B451DC7571/tenants
)

// GetVimTenants return list of cloud provider attached to TCA
func (c *RestClient) GetVimTenants() (*response.Tenants, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + apiTenants)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			return nil, fmt.Errorf(errRes.Message)
		}
		return nil, fmt.Errorf("unknown error, status code: %v ", resp.StatusCode())
	}

	var tenants response.Tenants
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		return nil, err
	}

	return &tenants, nil
}

// GetVim return list of cloud provider attached to TCA
func (c *RestClient) GetVim(vimId string) (*response.TenantSpecs, error) {

	c.GetClient()
	// format vmware_FB40D3DE2967483FBF9033B451DC7571
	glog.Infof("Sending request to ", c.BaseURL+apiTenants+"/"+vimId+"/tenants")

	resp, err := c.Client.R().Get(c.BaseURL + apiVim + "/" + vimId + "/tenants")
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var tenants response.TenantSpecs
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		return nil, err
	}

	return &tenants, nil
}

// GetTenantsQuery returns list of cloud provider attached to TCA
func (c *RestClient) GetTenantsQuery(f *request.TenantsNfFilter) (*response.Tenants, error) {

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).SetQueryString("action=query").
		Post(c.BaseURL + "/hybridity/api/vims/v1/tenants")
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(resp.RawResponse.Body).Decode(&errRes); err == nil {
			return nil, fmt.Errorf(errRes.Message)
		}
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var tenants response.Tenants
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		return nil, err
	}

	return &tenants, nil
}
