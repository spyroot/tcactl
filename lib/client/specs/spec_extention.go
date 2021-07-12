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
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spyroot/tcactl/lib/api_errors"
	"io"
	"strings"
)

// AdditionalFilters Filter for repo query
type AdditionalFilters struct {
	VimID string `json:"vimId" yaml:"vimId"`
}

type Filter struct {
	ExtraFilter AdditionalFilters `json:"additionalFilters"`
}

type RepoQuery struct {
	QueryFilter Filter `json:"filter"`
}

// InvalidExtensionSpec error if specs invalid
type InvalidExtensionSpec struct {
	errMsg string
}

//
func (m *InvalidExtensionSpec) Error() string {
	return m.errMsg
}

type VimInfo struct {
	VimName       string `json:"vimName" yaml:"vimName"`
	VimId         string `json:"vimId" yaml:"vimId"`
	VimSystemUUID string `json:"vimSystemUUID" yaml:"vimSystemUuid"`
}

type SpecInterfaceInfo struct {
	Url                string `json:"url" yaml:"url" validate:"required,url"`
	Description        string `json:"description,omitempty" yaml:"description,omitempty"`
	TrustedCertificate string `json:"trustedCertificate,omitempty" yaml:"trustedCertificate,omitempty"`
}

type SpecAdditionalParameters struct {
	TrustAllCerts bool `json:"trustAllCerts" yaml:"trustAllCerts"`
}

type SpecAccessInfo struct {
	// Username for Harbor it username that has admin access
	Username string `json:"username" yaml:"username" validate:"required"`
	// Password is base64 encoded string
	Password string `json:"password" yaml:"password" validate:"required"`
}

// SpecExtension extension such as Harbor registered in TCA
type SpecExtension struct {
	// SpecType indicate a spec type and meet Spec interface requirements.
	SpecType SpecType `json:"kind,omitempty" yaml:"kind,omitempty" valid:"required~kind is mandatory spec field"`
	//Name is extension name
	Name string `json:"name" yaml:"name" valid:"required~name is mandatory spec field"`
	// Version for harbor it 1.x 2.x
	Version              string                    `json:"version" yaml:"version" valid:"required~version is mandatory spec field"`
	Type                 string                    `json:"type" yaml:"type"  valid:"required~kind is mandatory spec field"`
	ExtensionKey         string                    `json:"extensionKey,omitempty" yaml:"extensionKey,omitempty"`
	ExtensionSubtype     string                    `json:"extensionSubtype" yaml:"extensionSubtype"`
	Products             []interface{}             `json:"products,omitempty" yaml:"products,omitempty"`
	VimInfo              []VimInfo                 `json:"vimInfo,omitempty" yaml:"vimInfo,omitempty"`
	InterfaceInfo        *SpecInterfaceInfo        `json:"interfaceInfo,omitempty" yaml:"interfaceInfo,omitempty"`
	AdditionalParameters *SpecAdditionalParameters `json:"additionalParameters,omitempty" yaml:"additionalParameters,omitempty"`
	AutoScaleEnabled     bool                      `json:"autoScaleEnabled" yaml:"autoScaleEnabled"`
	AutoHealEnabled      bool                      `json:"autoHealEnabled" yaml:"autoHealEnabled"`
	AccessInfo           *SpecAccessInfo           `json:"accessInfo,omitempty" yaml:"accessInfo,omitempty"`
	// hold specError
	specError error
}

// AddVim add target vim
func (t *SpecExtension) AddVim(name string) {
	t.VimInfo = append(t.VimInfo, VimInfo{VimName: name})
}

// GetVim return vim spec
func (t *SpecExtension) GetVim(name string) (*VimInfo, error) {
	n := strings.ToLower(name)
	for _, info := range t.VimInfo {
		if info.VimName == n {
			return &info, nil
		}
	}

	return nil, api_errors.NewVimNotFound(name)
}

// Kind return spec kind
func (t *SpecExtension) Kind() SpecType {
	return t.SpecType
}

//SpecsFromString method read cluster spec from string
//and return instance
func (t SpecExtension) SpecsFromString(s string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromString(s, new(SpecExtension), f...)
}

//SpecsFromFile method return instance form string
func (t SpecExtension) SpecsFromFile(fileName string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromFile(fileName, new(SpecExtension), f...)
}

// SpecsFromReader create spec from reader
func (t SpecExtension) SpecsFromReader(r io.Reader, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpec(r, new(SpecExtension), f...)
}

//Validate method validate specs
func (t *SpecExtension) Validate() error {

	if t == nil {
		return &InvalidExtensionSpec{errMsg: "nil instance"}
	}

	_, err := govalidator.ValidateStruct(t)
	if err != nil {
		t.specError = err
		return err
	}

	if t.Kind() != SpecKindExtension {
		t.specError = &InvalidExtensionSpec{errMsg: fmt.Sprintf(
			"Invalid spec kind. Extension spec must use kind %s", SpecKindExtension)}
		return t.specError
	}

	return nil
}

// Default TODO
func (t *SpecExtension) Default() error {
	return nil
}

// IsValid return false if validator set error
func (t *SpecExtension) IsValid() bool {

	if t.specError != nil {
		return false
	}

	return true
}
