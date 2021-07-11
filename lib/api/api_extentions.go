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
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/golang/glog"
	_ "github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/api_errors"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/pkg/errors"
)

// GetRepos retrieve repos by tenant id
func (a *TcaApi) GetRepos(ctx context.Context) (*response.ReposList, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		glog.Error(err)
	}

	var allRepos response.ReposList
	for _, r := range tenants.TenantsList {
		repos, err := a.rest.RepositoriesQuery(&specs.RepoQuery{
			QueryFilter: specs.Filter{
				ExtraFilter: specs.AdditionalFilters{
					VimID: r.TenantID,
				},
			},
		})

		if err != nil {
			return nil, err
		}

		allRepos.Items = append(allRepos.Items, repos.Items...)
	}

	return &allRepos, nil
}

// GetExtensions returns all extension registered in TCA
// This method useful if we need find Harbor or any other extension
// for example attached to particular VIM
func (a *TcaApi) GetExtensions(ctx context.Context) (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	return a.rest.GetExtensions(ctx)
}

// GetExtension api call return extension register in TCA
// NameOrId is name of extension or id.  If name provided
// method need resolve name to id first.
func (a *TcaApi) GetExtension(ctx context.Context, NameOrId string) (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	if IsValidUUID(NameOrId) {
		return a.rest.GetExtension(ctx, NameOrId)
	}

	eid, err := a.ResolveExtensionId(ctx, NameOrId)
	if err != nil {
		return nil, err
	}

	if len(eid) == 0 {
		return nil, api_errors.NewExtensionsNotFound(NameOrId)
	}

	extension, err := a.rest.GetExtension(ctx, eid)
	if err != nil {
		return nil, err
	}

	return extension, err
}

func (a *TcaApi) resolveVimInfo(ctx context.Context, spec *specs.SpecExtension) error {
	if spec.VimInfo != nil {
		for i, vimInfo := range spec.VimInfo {
			// if vim name preset resolve all vim ids if needed
			if len(vimInfo.VimName) > 0 {
				vim, err := a.GetVim(ctx, vimInfo.VimName)
				if err != nil {
					return err
				}

				tenant, err := vim.GetTenant(vimInfo.VimName)
				if err != nil {
					return err
				}

				if len(vimInfo.VimId) == 0 {
					spec.VimInfo[i].VimId = tenant.ID
				}

				if len(vimInfo.VimSystemUUID) == 0 {
					spec.VimInfo[i].VimSystemUUID = vim.VimId
				}

				glog.Infof("Resolved vim attributes name %s,  vim id %s, vim uuid %s",
					vimInfo.VimName, vimInfo.VimId, vimInfo.VimSystemUUID)

			} else {
				return fmt.Errorf("vim name is empty")
			}
		}
	}

	return nil
}

// CreateExtension api call create extension in TCA
func (a *TcaApi) CreateExtension(ctx context.Context, spec *specs.SpecExtension) (string, error) {

	if a == nil {
		return "", errors.NilError
	}

	if spec == nil {
		return "", api_errors.NewInvalidSpec("Spec is nil")
	}

	err := spec.Validate()
	if err != nil {
		return "", err
	}

	// remove kind and encode password as base64
	specCopy := spec
	specCopy.SpecType = ""

	if len(spec.AccessInfo.Password) > 0 {
		spec.AccessInfo.Password = b64.StdEncoding.EncodeToString([]byte(spec.AccessInfo.Password))
	}

	err = a.resolveVimInfo(ctx, spec)
	if err != nil {
		return "", err
	}

	return a.rest.CreateExtension(ctx, spec)
}

// ResolveExtensionId method resolve Name or Id to extension
// if extension already valid id it no op
func (a *TcaApi) ResolveExtensionId(ctx context.Context, NameOrId string) (string, error) {

	extension, err := a.GetExtensions(ctx)
	if err != nil {
		return "", err
	}

	ext, err := extension.FindExtension(NameOrId)
	if err != nil {
		return "", err
	}

	return ext.ExtensionId, nil
}

// DeleteExtension api call delete extension from tca
// return true
func (a *TcaApi) DeleteExtension(ctx context.Context, NameOrId string) (bool, error) {

	if a == nil {
		return false, errors.NilError
	}

	eid, err := a.ResolveExtensionId(ctx, NameOrId)
	if err != nil {
		return false, err
	}

	// remove kind and encode password as base64
	return a.rest.DeleteExtension(ctx, eid)
}

// UpdateExtension api call delete extension from TCA
func (a *TcaApi) UpdateExtension(ctx context.Context, spec *specs.SpecExtension) (bool, error) {

	if a == nil {
		return false, errors.NilError
	}

	if spec == nil {
		return false, api_errors.NewInvalidSpec("Spec is nil")
	}

	err := spec.Validate()
	if err != nil {
		return false, err
	}

	// remove kind and encode password as base64
	if len(spec.AccessInfo.Password) > 0 {
		spec.AccessInfo.Password = b64.StdEncoding.EncodeToString([]byte(spec.AccessInfo.Password))
	}

	extension, err := a.GetExtension(ctx, spec.Name)
	if err != nil {
		return false, err
	}

	e, err := extension.FindExtension(spec.Name)
	if err != nil {
		return false, err
	}

	specCopy := spec
	specCopy.SpecType = ""

	err = a.resolveVimInfo(ctx, spec)
	if err != nil {
		return false, err
	}

	return a.rest.UpdateExtension(ctx, spec, e.ExtensionId)
}

// ExtensionQuery - query for all extension api
func (a *TcaApi) ExtensionQuery(ctx context.Context) (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	return a.rest.GetExtensions(ctx)
}

// GetFilteredExtension retrieve repos by tenant id
func (a *TcaApi) GetFilteredExtension(ctx context.Context) (*response.ReposList, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		glog.Error(err)
	}

	var allRepos response.ReposList
	for _, r := range tenants.TenantsList {
		repos, err := a.rest.RepositoriesQuery(&specs.RepoQuery{
			QueryFilter: specs.Filter{
				ExtraFilter: specs.AdditionalFilters{
					VimID: r.TenantID,
				},
			},
		})

		if err != nil {
			return nil, err
		}

		allRepos.Items = append(allRepos.Items, repos.Items...)
	}

	return &allRepos, nil
}
