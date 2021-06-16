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
	"github.com/spyroot/hestia/cmd/client/response"
	"net/http"
)

const (
	apiVnfpkgm = "/telco/api/vnfpkgm/v2/vnf_packages"
)

// GetCatalogId return vnf package id and vnfd id
func (c *RestClient) GetCatalogId(q string) (string, string, error) {

	vnfCatalog, err := c.GetVnfPkgm("", "")
	if err != nil || vnfCatalog == nil {
		glog.Errorf("Failed acquire catalog information.")
		return "", "", err
	}

	pkgCnf, err := vnfCatalog.GetVnfdID(q)
	if err != nil || pkgCnf == nil {
		glog.Errorf("Failed acquire VNF information for %v", q)
		return "", "", err
	}

	return pkgCnf.PID, pkgCnf.VnfdID, nil
}

// GetVnfPkgm gets VNF/CNF catalog entry
func (c *RestClient) GetVnfPkgm(filter string, pkgId string) (*response.VnfPackages, error) {

	c.GetClient()
	r := c.Client.R()

	var restReq = c.BaseURL + apiVnfpkgm
	if len(pkgId) != 0 {
		restReq = c.BaseURL + apiVnfpkgm + "/" + pkgId
	}

	// attach filter
	if len(filter) != 0 {
		r.SetQueryParams(map[string]string{
			"filter": filter,
		},
		)
	}

	resp, err := r.Get(restReq)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes response.VnfPackagesError
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			return nil, fmt.Errorf(errRes.Detail)
		}
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var pkgs response.VnfPackages
	if len(pkgId) != 0 {
		var pkg response.VnfPackage
		if err := json.Unmarshal(resp.Body(), &pkg); err != nil {
			return nil, err
		}
		pkgs.Packages = append(pkgs.Packages, pkg)
	} else {
		if err := json.Unmarshal(resp.Body(), &pkgs.Packages); err != nil {
			return nil, err
		}
	}

	return &pkgs, nil
}

func (c *RestClient) GetVnfPkgmVnfd(pkgId string) (*response.VduPackage, error) {

	c.GetClient()
	var restReq string
	if len(pkgId) == 0 {
		restReq = c.BaseURL + "/telco/api/vnfpkgm/v2/vnf_packages"
	} else {
		restReq = c.BaseURL + "/telco/api/vnfpkgm/v2/vnf_packages/" + pkgId + "/vnfd"
	}

	r := c.Client.R()
	resp, err := r.Get(restReq)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes response.VnfPackagesError
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			return nil, fmt.Errorf(errRes.Detail)
		}
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	var pkg response.VduPackage
	if err := json.Unmarshal(resp.Body(), &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}
