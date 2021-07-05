// Package request
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
package request

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type SpecType string

// RegisterVimSpec main spec for cloud provider registration
type RegisterVimSpec struct {
	SpecType *SpecKind `json:"kind,omitempty" yaml:"kind,omitempty" validate:"required"`
	// HcxCloudUrl HCX CP full url
	HcxCloudUrl string `json:"hcxCloudUrl" yaml:"hcxCloudUrl" validate:"required,url"`
	// VimName Name that TCA show
	VimName string `json:"vimName" yaml:"vimName" validate:"required"`
	// TenantName for VC it "DEFAULT"
	TenantName string `json:"tenantName,omitempty" yaml:"tenantName,omitempty"`
	// Username authenticated to SSO
	Username string `json:"username" yaml:"username" validate:"required"`
	// Password password for SSO
	Password string `json:"password" yaml:"password" validate:"required"`
}

//GetKind return spec kind
func (p *RegisterVimSpec) GetKind() *SpecKind {
	return p.SpecType
}

// ProviderSpecsFromFile - reads tenant spec from file
// and return TenantSpecs instance
func ProviderSpecsFromFile(fileName string) (*RegisterVimSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadProviderSpec(file)
}

// ProviderSpecsFromFromString take string that holdw entire spec
// passed to reader and return TenantSpecs instance
func ProviderSpecsFromFromString(str string) (*RegisterVimSpec, error) {
	r := strings.NewReader(str)
	return ReadProviderSpec(r)
}

// ReadProviderSpec - Read spec from io reader
// detects format and use either yaml or json parse
func ReadProviderSpec(b io.Reader) (*RegisterVimSpec, error) {

	var spec RegisterVimSpec

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	return nil, &InvalidClusterSpec{"unknown format"}
}
