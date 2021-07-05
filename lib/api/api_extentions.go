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
	_ "github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/errors"
)

// GetExtension return all extension
func (a *TcaApi) GetExtension() (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	return a.rest.GetExtensions()
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

	if spec.VimInfo != nil {
		for _, vimInfo := range spec.VimInfo {
			// if vim name preset resolve all vim ids if needed
			if len(vimInfo.VimName) > 0 {
				vim, err := a.GetVim(vimInfo.VimName)
				if err != nil {
					return "", err
				}

				//				glog.Infof("resolved vim id %s %s", vim.VimId, vim.Tenants[0].)
				if len(vimInfo.VimId) == 0 {
					vimInfo.VimId = vim.VimId
				}
				if len(vimInfo.VimSystemUUID) == 0 {

				}
			} else {
				return "", fmt.Errorf("Vim name is empty")
			}
		}
	}
	return a.rest.CreateExtension(spec)
}

func (a *TcaApi) ResolveExtensionId(NameOrId string) (string, error) {

	extension, err := a.GetExtension()
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
func (a *TcaApi) UpdateExtension(spec *request.ExtensionSpec) (interface{}, error) {

	if a == nil {
		return nil, errors.NilError
	}

	err := a.specValidator.Struct(spec)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	// remove kind and encode password as base64
	specCopy := spec
	specCopy.SpecType = nil

	return a.rest.UpdateExtension(spec)
}

// ExtensionQuery - query for all extension api
func (a *TcaApi) ExtensionQuery() (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	return a.rest.GetExtensions()
}
