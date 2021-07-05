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
	b64 "encoding/base64"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang/glog"
	_ "github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/api_errors"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/errors"
)

// GetExtensions returns all extension registered in TCA
// This method useful if we need find Harbor or any other extension
// for example attached to particular VIM
func (a *TcaApi) GetExtensions() (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	return a.rest.GetExtensions()
}

// GetExtension api call return extension register in TCA
// NameOrId is name of extension or id.  If name provided
// method need resolve name to id first.
func (a *TcaApi) GetExtension(NameOrId string) (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	if IsValidUUID(NameOrId) {
		return a.rest.GetExtension(NameOrId)
	}

	eid, err := a.ResolveExtensionId(NameOrId)
	if err != nil {
		return nil, err
	}

	if len(eid) == 0 {
		return nil, api_errors.NewExtensionsNotFound(NameOrId)
	}

	extension, err := a.rest.GetExtension(eid)
	if err != nil {
		return nil, err
	}

	return extension, err
}

func (a *TcaApi) resolveVimInfo(spec *request.ExtensionSpec) error {
	if spec.VimInfo != nil {
		for i, vimInfo := range spec.VimInfo {
			// if vim name preset resolve all vim ids if needed
			if len(vimInfo.VimName) > 0 {
				vim, err := a.GetVim(vimInfo.VimName)
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
func (a *TcaApi) CreateExtension(spec *request.ExtensionSpec) (string, error) {

	if a == nil {
		return "", errors.NilError
	}

	err := a.specValidator.Struct(spec)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return "", validationErrors
	}

	// remove kind and encode password as base64
	specCopy := spec
	specCopy.SpecType = nil
	if len(spec.AccessInfo.Password) > 0 {
		spec.AccessInfo.Password = b64.StdEncoding.EncodeToString([]byte(spec.AccessInfo.Password))
	}

	err = a.resolveVimInfo(spec)
	if err != nil {
		return "", err
	}

	return a.rest.CreateExtension(spec)
}

// ResolveExtensionId method resolve Name or Id to extension
// if extension already valid id it no op
func (a *TcaApi) ResolveExtensionId(NameOrId string) (string, error) {

	extension, err := a.GetExtensions()
	if err != nil {
		return "", err
	}

	ext, err := extension.FindExtension(NameOrId)
	if err != nil {
		return "", err
	}

	return ext.ExtensionId, nil
}

// DeleteExtension api call delete extension from TCA
func (a *TcaApi) DeleteExtension(NameOrId string) (*models.TcaTask, error) {

	if a == nil {
		return nil, errors.NilError
	}

	eid, err := a.ResolveExtensionId(NameOrId)
	if err != nil {
		return nil, err
	}

	// remove kind and encode password as base64
	return a.rest.DeleteExtension(eid)
}

// UpdateExtension api call delete extension from TCA
func (a *TcaApi) UpdateExtension(spec *request.ExtensionSpec, eid string) (interface{}, error) {

	if a == nil {
		return nil, errors.NilError
	}

	fmt.Println("SPec type", spec.SpecType)

	err := a.specValidator.Struct(spec)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	// remove kind and encode password as base64
	if len(spec.AccessInfo.Password) > 0 {
		spec.AccessInfo.Password = b64.StdEncoding.EncodeToString([]byte(spec.AccessInfo.Password))
	}

	specCopy := spec
	specCopy.SpecType = nil

	err = a.resolveVimInfo(spec)
	if err != nil {
		return nil, err
	}

	return a.rest.UpdateExtension(spec, eid)
}

// ExtensionQuery - query for all extension api
func (a *TcaApi) ExtensionQuery() (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	return a.rest.GetExtensions()
}
