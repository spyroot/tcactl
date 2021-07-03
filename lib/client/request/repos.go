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

// AdditionalFilters Filter for repo query
type AdditionalFilters struct {
	VimID string `json:"vimId" yaml:"vim_id"`
}

type Filter struct {
	ExtraFilter AdditionalFilters `json:"additionalFilters"`
}

type RepoQuery struct {
	QueryFilter Filter `json:"filter"`
}

type ExtensionSpec struct {
	// SpecType indicate a spec type and meet Spec interface requirements.
	SpecType         *SpecKind     `json:"kind,omitempty" yaml:"kind,omitempty" validate:"required"`
	Name             string        `json:"name" yaml:"name"`
	Version          string        `json:"version" yaml:"version"`
	Type             string        `json:"type" yaml:"type"`
	ExtensionKey     string        `json:"extensionKey" yaml:"extensionKey"`
	ExtensionSubtype string        `json:"extensionSubtype" yaml:"extensionSubtype"`
	Products         []interface{} `json:"products" yaml:"products"`
	VimInfo          []interface{} `json:"vimInfo" yaml:"vim_info"`
	InterfaceInfo    struct {
		Url                string `json:"url" yaml:"url"`
		Description        string `json:"description" yaml:"description"`
		TrustedCertificate string `json:"trustedCertificate" yaml:"trustedCertificate"`
	} `json:"interfaceInfo" yaml:"interfaceInfo"`
	AdditionalParameters struct {
		TrustAllCerts bool `json:"trustAllCerts" yaml:"trustAllCerts"`
	} `json:"additionalParameters" yaml:"additionalParameters"`
	AutoScaleEnabled bool `json:"autoScaleEnabled" yaml:"autoScaleEnabled"`
	AutoHealEnabled  bool `json:"autoHealEnabled" yaml:"autoHealEnabled"`
	AccessInfo       struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	} `json:"accessInfo" yaml:"accessInfo"`
}

// ExtensionSpecsFromFile - reads spec spec from file
// and return ReadExtensionSpec instance
func ExtensionSpecsFromFile(fileName string) (*ExtensionSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadExtensionSpec(file)
}

// ExtensionSpecFromFromString take string that holds entire spec
// passed to reader and return ReadExtensionSpec instance
func ExtensionSpecFromFromString(str string) (*ExtensionSpec, error) {
	r := strings.NewReader(str)
	return ReadExtensionSpec(r)
}

// ReadExtensionSpec - Read ReadExtensionSpec from io reader
func ReadExtensionSpec(b io.Reader) (*ExtensionSpec, error) {

	var spec ExtensionSpec

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
