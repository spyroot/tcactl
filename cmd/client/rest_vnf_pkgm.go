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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/client/response"
	_ "mime/multipart"
	"net/http"
)

const (
	apiVnfpkgm = "/telco/api/vnfpkgm/v2/vnf_packages"

	UploadMultipartContentType = "application/zip"

	MultipartFilePart = "file"
)

//PackageUpload - Package Upload REST request.
// it contains name and optional tag
type PackageUpload struct {
	UserDefinedData struct {
		Name string        `json:"name" yaml:"name"`
		Tags []interface{} `json:"tags" yaml:"tags"`
	} `json:"userDefinedData"`
}

// NewPackageUpload creates a new Package Upload REST
// request object
func NewPackageUpload(packageName string) *PackageUpload {
	p := PackageUpload{}
	p.UserDefinedData.Name = packageName
	return &p
}

// PackageCreatedSuccess - when package created,
// TCA respond with Success message.
type PackageCreatedSuccess struct {
	// Id all operation to existing catalog must use this ID
	Id string `json:"id" yaml:"id"`
	// OnboardingState state
	OnboardingState  string        `json:"onboardingState"  yaml:"onboardingState"`
	OperationalState string        `json:"operationalState" yaml:"operationalState"`
	UsageState       string        `json:"usageState" yaml:"usageState"`
	VnfmInfo         []interface{} `json:"vnfmInfo" yaml:"vnfmInfo"`
	UserDefinedData  struct {
		Name string        `json:"name" yaml:"name"`
		Tags []interface{} `json:"tags" yaml:"tags"`
	} `json:"userDefinedData" yaml:"userDefinedData"`
	// Links re-present a link that caller can use to execute update, delete operation
	Links struct {
		Self struct {
			Href string `json:"href" yaml:"href"`
		} `json:"self" yaml:"self"`
	} `json:"_links" yaml:"links"`
}

// GetPackageCatalogId return vnf package id and vnfd id
func (c *RestClient) GetPackageCatalogId(q string) (string, string, error) {

	vnfCatalog, err := c.GetVnfPkgm("", "")
	if err != nil || vnfCatalog == nil {
		glog.Errorf("Failed acquire package catalog information. %v", err)
		return "", "", err
	}

	pkgCnf, err := vnfCatalog.GetVnfdID(q)
	if err != nil || pkgCnf == nil {
		glog.Errorf("Failed acquire CNF/VNF package id for %v, err: %v", q, err)
		return "", "", err
	}

	return pkgCnf.PID, pkgCnf.VnfdID, nil
}

// GetVnfPkgm gets VNF/CNF catalog entity
// pkgId is catalog id and filter is optional argument
// is filter query
func (c *RestClient) GetVnfPkgm(filter string, pkgId string) (*response.VnfPackages, error) {

	c.GetClient()
	r := c.Client.R()

	var restReq = c.BaseURL + apiVnfpkgm
	if len(pkgId) != 0 {
		restReq = c.BaseURL + apiVnfpkgm + "/" + pkgId
	}

	// attach query filter
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

// GetVnfPkgmVnfd return CNF/VNF catalog entity.
func (c *RestClient) GetVnfPkgmVnfd(pkgId string) (*response.VduPackage, error) {

	c.GetClient()
	var restReq string
	if len(pkgId) == 0 {
		restReq = c.BaseURL + TcaVmwareTelcoPackages
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

// DeleteVnfPkgmVnfd - delete package
func (c *RestClient) DeleteVnfPkgmVnfd(pkgId string) (bool, error) {

	c.GetClient()
	resp, err := c.Client.R().Delete(c.BaseURL + TcaVmwareTelcoPackages + "/" + pkgId)
	if err != nil {
		glog.Error(err)
		return false, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return false, c.checkError(resp)

	}

	return resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusNoContent, nil
}

// CreateVnfPkgmVnfd - create a new CNF/VNF package
// Note this call only creates catalog entry,
// entity disabled when it created.
func (c *RestClient) CreateVnfPkgmVnfd(pkg *PackageUpload) (*PackageCreatedSuccess, error) {

	c.GetClient()
	resp, err := c.Client.R().SetBody(pkg).Post(c.BaseURL + TcaVmwareTelcoPackages)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes PackageCreateError
		if err := json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error '%s' details '%s' instance '%s'",
				errRes.Title, errRes.Detail, errRes.Instance)
			return nil, fmt.Errorf("error '%s', details '%s', instance '%s'",
				errRes.Title, errRes.Detail, errRes.Instance)
		} else {
			glog.Errorf("Failed parse server respond.")
		}
		return nil, c.checkError(resp)
	}

	var success PackageCreatedSuccess
	if err := json.Unmarshal(resp.Body(), &success); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &success, nil
}

// UploadVnfPkgmVnfd - Uploads vnf package,
// note method doesn't check if catalog entity
// already create or not.
func (c *RestClient) UploadVnfPkgmVnfd(pkgId string, csar []byte, name string) (bool, error) {

	if len(name) == 0 {
		glog.Error("Received empty name.")
		return false, fmt.Errorf("received empty name")
	}

	if len(csar) == 0 {
		glog.Error("Received empty data bytes.")
		return false, fmt.Errorf("received empty data")
	}

	c.GetClient()
	resp, err := c.Client.R().
		SetFileReader(MultipartFilePart, name, bytes.NewReader(csar)).
		SetHeader("Content-Type", UploadMultipartContentType).
		SetContentLength(true).
		Put(c.BaseURL + "/telco/api/vnfpkgm/v2/vnf_packages/" + pkgId + "/package_content")

	if err != nil {
		glog.Error(err)
		return false, err
	}

	if c.dumpRespond && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return false, c.checkError(resp)

	}

	return true, nil
}
