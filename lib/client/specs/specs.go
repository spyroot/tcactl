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
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type SpecFormatType int

type SpecType string

const (
	// SpecKindProvider provider registration spec
	SpecKindProvider SpecType = "provider"

	// SpecKindExtension spec type extension registration
	SpecKindExtension SpecType = "extensions"

	// SpecKindNodePool spec type node pool
	SpecKindNodePool SpecType = "node_pool"

	//SpecKindTemplate spec type cluster template
	SpecKindTemplate SpecType = "template"

	//SpecKindInstance spec type instance
	SpecKindInstance SpecType = "instance"

	//SpecKindCluster spec type cluster
	SpecKindCluster SpecType = "cluster"

	Unknown SpecFormatType = iota
	Yaml
	Json
	Xml
)

func (t SpecFormatType) String() string {
	return [...]string{"yaml", "json", "xml"}[t]
}

type RequestSpec interface {
	Kind() SpecType
	Default() error
	Validate() error
	IsValid() bool
}

// ReadSpecFromFromFile - reads instance spec from file
// and return TenantSpecs instance
func ReadSpecFromFromFile(fileName string, spec RequestSpec, f ...SpecFormatType) (*RequestSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	if len(f) == 0 {
		fileName = path.Base(fileName)
		if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
			spec, err := ReadSpec(file, spec, Yaml)
			if err != nil {
				return nil, err
			}
			return spec, nil
		} else {
			if strings.HasSuffix(fileName, ".json") {
				spec, err := ReadSpec(file, spec, Json)
				if err != nil {
					return nil, err
				}
				return spec, nil
			}
		}
	}

	return ReadSpec(file, spec, f...)
}

// ReadSpecFromFromString take string that holds entire spec
// passed to reader and returns InstanceRequestSpec instance
func ReadSpecFromFromString(str string, spec RequestSpec, f ...SpecFormatType) (*RequestSpec, error) {
	r := strings.NewReader(str)
	return ReadSpec(r, spec, f...)
}

// ReadSpec - Read spec from io reader
// detects format and uses either yaml or json parser
func ReadSpec(b io.Reader, spec RequestSpec, f ...SpecFormatType) (*RequestSpec, error) {

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	sType := Unknown
	if len(f) > 0 {
		sType = f[0]
	}

	isYamlSpec := false
	isJsonSpec := false

	switch sType {
	case Yaml:
		isYamlSpec = true
	case Json:
		isJsonSpec = true
	default:
		isYamlSpec = true
		isJsonSpec = true
	}

	if isJsonSpec {
		if err = json.Unmarshal(buffer, spec); err == nil {
			err = spec.Default()
			if err != nil {
				return nil, err
			}
			return &spec, nil
		}
	}

	if isYamlSpec {
		if err = yaml.Unmarshal(buffer, spec); err == nil {
			err = spec.Default()
			if err != nil {
				return nil, err
			}
			return &spec, nil
		}
	}

	return nil, &InvalidClusterSpec{"unknown format"}
}
