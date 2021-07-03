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
	"github.com/go-playground/validator/v10"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
)

// TenantFields method return all struct
// fields name
func TenantFields() []string {
	f := response.TenantsDetails{}
	fields, _ := f.GetFields()

	var keys []string
	for s, _ := range fields {
		keys = append(keys, s)
	}

	return keys
}

// TenantsCloudProvider return a tenant attached to cloud provide for lookup query string
func (a *TcaApi) TenantsCloudProvider(query string) (*response.Tenants, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(query) == 0 {
		return nil, fmt.Errorf("empty query string")
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

// DeleteTenantsProvider delete a tenant attached to cloud provide.
func (a *TcaApi) DeleteTenantsProvider(tenantCluster string) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	r, err := tenants.FindCloudProvider(tenantCluster)
	if err != nil {
		return nil, err
	}

	return a.rest.DeleteTenant(r.TenantID)
}

// CreateTenantProvider method create, registers new target cloud provider
// as tenant infrastructure in TCA.
func (a *TcaApi) CreateTenantProvider(spec *request.RegisterVim) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	err := a.specValidator.Struct(spec)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	// remove kind and encode password as base64
	specCopy := spec
	specCopy.SpecType = nil
	//	spec.Password = b64.StdEncoding.EncodeToString([]byte(spec.Password))

	return a.rest.RegisterCloudProvider(specCopy)
}
