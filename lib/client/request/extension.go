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
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/api_errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type SpecFormat int

const (
	YamlFile SpecFormat = iota
	JsonFile
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

type VimInfo struct {
	VimName       string `json:"vimName" yaml:"vimName"`
	VimId         string `json:"vimId" yaml:"vimId"`
	VimSystemUUID string `json:"vimSystemUUID" yaml:"vimSystemUuid"`
}

// ExtensionSpec extension such as Harbor registered in TCA
type ExtensionSpec struct {
	// SpecType indicate a spec type and meet Spec interface requirements.
	SpecType *SpecKind `json:"kind,omitempty" yaml:"kind,omitempty" validate:"required"`
	Name     string    `json:"name" yaml:"name" validate:"required"`
	// Version for harbor it 1.x 2.x
	Version          string        `json:"version" yaml:"version" validate:"required"`
	Type             string        `json:"type" yaml:"type" validate:"required"`
	ExtensionKey     string        `json:"extensionKey,omitempty" yaml:"extensionKey,omitempty"`
	ExtensionSubtype string        `json:"extensionSubtype" yaml:"extensionSubtype" validate:"required"`
	Products         []interface{} `json:"products" yaml:"products"`
	VimInfo          []VimInfo     `json:"vimInfo" yaml:"vimInfo"`
	InterfaceInfo    struct {
		Url                string `json:"url" yaml:"url" validate:"required,url"`
		Description        string `json:"description,omitempty" yaml:"description,omitempty"`
		TrustedCertificate string `json:"trustedCertificate,omitempty" yaml:"trustedCertificate,omitempty"`
	} `json:"interfaceInfo" yaml:"interfaceInfo"`
	AdditionalParameters struct {
		TrustAllCerts bool `json:"trustAllCerts" yaml:"trustAllCerts"`
	} `json:"additionalParameters" yaml:"additionalParameters"`
	AutoScaleEnabled bool `json:"autoScaleEnabled" yaml:"autoScaleEnabled"`
	AutoHealEnabled  bool `json:"autoHealEnabled" yaml:"autoHealEnabled"`
	//
	AccessInfo struct {
		// Username for Harbor it username that has admin access
		Username string `json:"username" yaml:"username" validate:"required"`
		// Password is base64 encoded string
		Password string `json:"password" yaml:"password" validate:"required"`
	} `json:"accessInfo" yaml:"accessInfo"`
}

// AddVim add target vim
func (s *ExtensionSpec) AddVim(name string) {
	s.VimInfo = append(s.VimInfo, VimInfo{VimName: name})
}

// GetVim return vim spec
func (s *ExtensionSpec) GetVim(name string) (*VimInfo, error) {
	n := strings.ToLower(name)
	for _, info := range s.VimInfo {
		if info.VimName == n {
			return &info, nil
		}
	}

	return nil, api_errors.NewVimNotFound("name")
}

// ExtensionSpecsFromFile - reads spec spec from file
// and return ReadExtensionSpec instance
func ExtensionSpecsFromFile(fileName string) (*ExtensionSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	fileName = path.Base(fileName)
	glog.Infof("Parsing file %s", file)
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		spec, err := ReadExtensionSpec(file, YamlFile)
		if err != nil {
			return nil, err
		}
		return spec, nil
	}
	if strings.HasSuffix(fileName, "json") {
		spec, err := ReadExtensionSpec(file, YamlFile)
		if err != nil {
			return nil, err
		}
		return spec, nil
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
func ReadExtensionSpec(b io.Reader, f ...SpecFormat) (*ExtensionSpec, error) {

	var spec ExtensionSpec

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	if len(f) == 0 {
		err = json.Unmarshal(buffer, &spec)
		if err == nil {
			return &spec, nil
		} else {
			glog.Errorf("Tried json and got error: %v", err)
		}

		err = yaml.Unmarshal(buffer, &spec)
		if err == nil {
			return &spec, nil
		} else {
			glog.Errorf("Tried yaml and got error: %v", err)
		}
	} else {
		glog.Infof("File type provided")
		if f[0] == YamlFile {
			err = yaml.Unmarshal(buffer, &spec)
			if err != nil {
				glog.Errorf("Error: %v", err)
			}
			return &spec, nil
		}
		if f[0] == YamlFile {
			err = yaml.Unmarshal(buffer, &spec)
			if err != nil {
				glog.Errorf("Error: %v", err)
			}
			return &spec, nil
		}
	}

	return nil, &InvalidClusterSpec{"unknown format"}
}
