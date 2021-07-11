// Package api
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
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/response"
)

// GetEntireCatalog - api call return entire TCA catalog.
func (a *TcaApi) GetEntireCatalog() (*response.VnfPackages, error) {
	glog.Infof("Fetching entire tca catalog")
	return a.rest.GetAllCatalog()
}

// GetVnfPkgm - return catalog entity
// pkgId is id of package in TCA
// filter
func (a *TcaApi) GetVnfPkgm(filter string, pkgId string) (*response.VnfPackages, error) {
	glog.Infof("Fetching catalog entity id %s", pkgId)
	return a.rest.GetVnfPkgm(filter, pkgId)
}

// GetCatalogId return vnf Package ID and VNFD ID
func (a *TcaApi) GetCatalogId(catalogId string) (string, string, error) {
	glog.Infof("Retrieving vnf packages for catalog entity %s.", catalogId)
	return a.rest.GetPackageCatalogId(catalogId)
}

// GetCatalogAndVdu API method returns
// catalog entity and vdu package.
func (a *TcaApi) GetCatalogAndVdu(nfdName string) (*response.VnfPackage, *response.VduPackage, error) {

	glog.Infof("Fetching catalog entity %s", nfdName)

	vnfCatalog, err := a.rest.GetVnfPkgm("", "")
	if err != nil || vnfCatalog == nil {
		glog.Errorf("Failed acquire vnf package information. Error %v", err)
		return nil, nil, err
	}

	catalogEntity, err := vnfCatalog.GetVnfdID(nfdName)
	if err != nil || catalogEntity == nil {
		glog.Errorf("Failed acquire catalog information for catalog name %v", nfdName)
		return nil, nil, err
	}

	v, err := a.rest.GetVnfPkgmVnfd(catalogEntity.PID)
	return catalogEntity, v, err
}
