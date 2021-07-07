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
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/csar"
	"io/ioutil"
	"os"
	"path/filepath"
)

// GetAllPackages return all packages
func (a *TcaApi) GetAllPackages() (*response.CnfsExtended, error) {

	respond, err := a.GetVnflcm()
	if err != nil {
		return nil, err
	}

	pkgs, ok := respond.(response.CnfsExtended)
	if !ok {
		return nil, errors.New("rest vnflcm return wrong type")
	}

	return &pkgs, nil
}

func (a *TcaApi) GetVnflcm(f ...string) (interface{}, error) {
	return a.rest.GetVnflcm(f...)
}

// GetCnfs method return list of cnf instances in
// response.CnfsExtended that encapsulate in collection
func (a *TcaApi) GetCnfs() (*response.CnfsExtended, error) {

	genericRespond, err := a.rest.GetVnflcm()
	if err != nil {
		return nil, err
	}

	// for extension request we route to correct printer
	cnfs, ok := genericRespond.(*response.CnfsExtended)
	if ok {
		return cnfs, nil
	}

	return nil, err
}

// CreateCatalogEntity method create a new package
// it take file name that must compressed zip file
// package catalog name and a substitution map.
// substitution map used to replace CSAR values.
// a key of map is key in CSAR and value a new value
// that used to replace value in actual CSAR.
// i.e  existing CSAR used as template and substitution
// map applied a transformation.
func (a *TcaApi) CreateCatalogEntity(
	fileName string,
	catalogName string,
	substitution map[string]string) (bool, error) {

	glog.Infof("Create new package. Received substitution %v.", substitution)

	if a.rest == nil {
		return false, fmt.Errorf("rest interface is nil")
	}

	// Apply transformation to a CSAR file
	newCsarFile, err := csar.ApplyTransformation(
		fileName,
		csar.SpecNfd,                    // a file inside a CSAR that we need apply transformation
		csar.NfdYamlPropertyTransformer, // a callback that apply transformation
		substitution)
	if err != nil {
		glog.Errorf("Failed apply transformation %v", err)
		return false, err
	}

	file, err := os.Open(newCsarFile)
	if err != nil {
		glog.Errorf("Failed read , newly generated csar %v", err)
		return false, err
	}

	// Read new CSAR file, to buffer
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		glog.Errorf("Failed read generated csar %v", err)
		return false, err
	}

	newFileName := filepath.Base(newCsarFile)
	uploadReq := client.NewPackageUpload(catalogName)
	respond, err := a.rest.CreateVnfPkgmVnfd(uploadReq)
	if err != nil {
		glog.Errorf("Failed create cnf package entity generated csar %v", err)
		return false, err
	}

	if len(respond.Id) == 0 {
		glog.Error("Something is wrong, server must contain package id in respond")
		return false, fmt.Errorf("respond doesn't contain package id")
	}

	// upload csar to a catalog
	ok, err := a.rest.UploadVnfPkgmVnfd(respond.Id, fileBytes, newFileName)
	if err != nil {
		return false, err
	}

	// TODO do GET to cross check and respond with ok if package is created.
	return ok, nil
}

// DeleteCatalogEntity method delete catalog entity
func (a *TcaApi) DeleteCatalogEntity(
	catalogName string) (bool, error) {

	glog.Infof("Delete catalog entity %v.", catalogName)

	if a.rest == nil {
		return false, fmt.Errorf("rest interface is nil")
	}

	pid, _, err := a.rest.GetPackageCatalogId(catalogName)
	if err != nil {
		return false, err
	}

	// upload csar to a catalog
	ok, err := a.rest.DeleteVnfPkgmVnfd(pid)
	if err != nil {
		return false, err
	}

	return ok, nil
}
