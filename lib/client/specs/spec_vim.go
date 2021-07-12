// Package specs
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
package specs

import (
	"github.com/asaskevich/govalidator"
	"io"
)

// InvalidCloudSpec error if cloud provider specs is invalid
type InvalidCloudSpec struct {
	errMsg string
}

// Error return error msg
func (m *InvalidCloudSpec) Error() string {
	return m.errMsg
}

// SpecCloudProvider main spec for cloud provider registration
type SpecCloudProvider struct {
	//
	SpecType SpecType `json:"kind,omitempty" yaml:"kind,omitempty" valid:"required~kind is mandatory spec field"`

	// HcxCloudUrl HCX CP full url
	HcxCloudUrl string `json:"hcxCloudUrl" yaml:"hcxCloudUrl" valid:"required~hcxCloudUrl is mandatory spec field"`

	// VimName Name that TCA show
	VimName string `json:"vimName" yaml:"vimName" valid:"required~kind is mandatory spec field"`

	// Username used used to SSO domain
	Username string `json:"username" yaml:"username" valid:"required~username is mandatory spec field"`

	// Password password used for for SSO domain authentication
	Password string `json:"password" yaml:"password" valid:"required~password is mandatory spec field"`

	// TenantName for VC it "DEFAULT"
	TenantName string `json:"tenantName,omitempty" yaml:"tenantName,omitempty"`

	// specError hold validator error
	specError error
}

// Kind must return "provider" SpecType
func (c *SpecCloudProvider) Kind() SpecType {
	return c.SpecType
}

// Default set default values
func (c *SpecCloudProvider) Default() error {
	c.TenantName = ""
	return nil
}

// IsValid return false if validator set error
func (c *SpecCloudProvider) IsValid() bool {
	if c.specError != nil {
		return false
	}
	return true
}

//Validate method validates cloud provider spec.
//and if mandatory field not set return error.
func (c *SpecCloudProvider) Validate() error {

	if c == nil {
		return &InvalidCloudSpec{errMsg: "nil instance"}
	}

	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		c.specError = err
		return err
	}

	if c.Kind() != SpecKindProvider {
		c.specError = &InvalidCloudSpec{errMsg: "spec must contain kind field"}
		return c.specError
	}

	return nil
}

//SpecsFromString method read cluster spec from string
//and return instance
func (c SpecCloudProvider) SpecsFromString(s string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromString(s, new(SpecCloudProvider), f...)
}

//SpecsFromFile method return instance form string
func (c SpecCloudProvider) SpecsFromFile(fileName string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromFile(fileName, new(SpecCloudProvider), f...)
}

// SpecsFromReader create spec from reader
func (c SpecCloudProvider) SpecsFromReader(r io.Reader, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpec(r, new(SpecCloudProvider), f...)
}
