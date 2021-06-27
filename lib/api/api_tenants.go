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
	"github.com/spyroot/tcactl/lib/client/response"
)

// TenantsCloudProvider return a tenant attached to cloud provide for lookup query string
func (a *TcaApi) TenantsCloudProvider(query string) (*response.Tenants, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	r, err := tenants.FindCloudProvider(query)
	if err != nil {
		return nil, err
	}

	return &response.Tenants{
		TenantsList: []response.TenantsDetails{*r},
	}, nil
}
