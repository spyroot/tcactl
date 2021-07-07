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
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spyroot/tcactl/lib/client/response"
	errnos "github.com/spyroot/tcactl/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func TemplateFields() []string {
	f := response.TenantSpecs{}
	fields, _ := f.GetFields()

	var keys []string
	for s, _ := range fields {
		keys = append(keys, s)
	}

	return keys
}

func ReadTemplateSpecFromFile(fileName string) (*response.ClusterTemplateSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadTemplateSpec(file)
}

// ReadTemplateSpecFromString read specString from reader
func ReadTemplateSpecFromString(str string) (*response.ClusterTemplateSpec, error) {
	r := strings.NewReader(str)
	return ReadTemplateSpec(r)
}

// ReadTemplateSpec - Read Template Spec
func ReadTemplateSpec(b io.Reader) (*response.ClusterTemplateSpec, error) {

	var spec response.ClusterTemplateSpec

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	return nil, &InvalidSpec{"unknown format"}
}

// GetClusterTemplates - return list of cluster templates
func (a *TcaApi) GetClusterTemplates() (*response.ClusterTemplates, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return nil, errnos.RestNilError
	}

	return a.rest.GetClusterTemplates()
}

// CreateClusterTemplate - create cluster template from initialSpec
func (a *TcaApi) CreateClusterTemplate(spec *response.ClusterTemplateSpec) (string, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return "", errnos.RestNilError
	}

	_id, err := a.ResolveTemplateId(spec.Name)
	if err == nil && len(_id) > 0 {
		// generate initialSpec name
		spec.Name = spec.Name + "-" + uuid.New().String()
		spec.Name = string(spec.Name[0:25])
	}

	// adjust case sensitivity
	spec.ClusterType = strings.ToUpper(spec.ClusterType)

	if len(spec.MasterNodes) == 0 {
		return "", &InvalidSpec{" master node section not present."}
	}

	if len(spec.WorkerNodes) == 0 {
		return "", &InvalidSpec{" worker node section not present."}
	}

	err = a.specValidator.Struct(spec)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return "", validationErrors
	}

	return spec.Name, a.rest.CreateClusterTemplate(spec)
}

// GetClusterTemplate return cluster template
func (a *TcaApi) GetClusterTemplate(templateId string) (*response.ClusterTemplateSpec, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return nil, errnos.RestNilError
	}

	var (
		tid = templateId
		err error
	)

	if !IsValidUUID(templateId) {
		tid, err = a.ResolveTemplateId(templateId)
		if err != nil {
			return nil, err
		}
	}

	return a.rest.GetClusterTemplate(tid)
}

// UpdateClusterTemplate - updates cluster template based
// on input initialSpec
func (a *TcaApi) UpdateClusterTemplate(spec *response.ClusterTemplateSpec) error {

	glog.Infof("Updating cluster template. %v", spec)

	if a.rest == nil {
		return errnos.RestNilError
	}

	if IsValidUUID(spec.Id) {
		glog.Infof("Validating template id %s", spec.Id)
		_, err := a.rest.GetClusterTemplate(spec.Id)
		if err != nil {
			glog.Error(err)
			return err
		}
	} else {
		// if template name used , resolve id
		id, err := a.ResolveTemplateId(spec.Id)
		if err != nil {
			return err
		}
		spec.Id = id
	}

	err := a.specValidator.Struct(spec)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	return a.rest.UpdateClusterTemplate(spec)
}

// DeleteTemplate deletes cluster template from TCA
// template argument can be name or ID.
func (a *TcaApi) DeleteTemplate(template string) error {

	if a.rest == nil {
		return errnos.RestNilError
	}

	var templateId = ""

	if IsValidUUID(template) {
		tmpl, err := a.rest.GetClusterTemplate(template)
		if err != nil {
			return err
		}
		templateId = tmpl.Id
	} else {
		templates_, err := a.rest.GetClusterTemplates()
		if err != nil {
			return err
		}
		templateId, err = templates_.GetTemplateId(template)
		if err != nil {
			return err
		}
		glog.Infof("Resolved template id %s", templateId)
	}

	err := a.rest.DeleteClusterTemplate(templateId)
	if err != nil {
		return err
	}

	return nil
}

// ResolveTemplateId - resolves template name to id
func (a *TcaApi) ResolveTemplateId(templateId string) (string, error) {

	// resolve template id, in case client used name instead id
	clusterTemplates, err := a.rest.GetClusterTemplates()
	if err != nil {
		return "", err
	}

	template, err := clusterTemplates.GetTemplate(templateId)
	if err != nil {
		return "", err
	}

	return template.Id, nil
}
