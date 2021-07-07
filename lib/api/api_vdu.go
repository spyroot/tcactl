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
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/api_errors"
	"github.com/spyroot/tcactl/lib/client/response"
)

// GetVdu retrieve Vdu
func (a *TcaApi) GetVdu(nfdName string) (*response.VduPackage, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(nfdName) == 0 {
		return nil, api_errors.NewCatalogNotFound(nfdName)
	}

	vnfCatalog, err := a.rest.GetVnfPkgm("", "")
	if err != nil || vnfCatalog == nil {
		glog.Errorf("Failed retrieve vnf package information, error %v", err)
		return nil, err
	}

	pkgCnf, err := vnfCatalog.GetVnfdID(nfdName)
	if err != nil || pkgCnf == nil {
		glog.Errorf("Failed retrieve vnfd information for %v.", nfdName)
		return nil, err
	}

	vnfd, err := a.rest.GetVnfPkgmVnfd(pkgCnf.PID)
	if err != nil || vnfd == nil {
		glog.Errorf("Failed acquire VDU information for %v.", pkgCnf.PID)
		return nil, err
	}

	return vnfd, nil
}
