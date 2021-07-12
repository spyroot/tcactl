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
package specs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Read cloud provider spec from string
func TestSpecCloudProvider_SpecsFromString(t *testing.T) {

	tests := []struct {
		name              string
		spec              string
		wantErr           bool
		wantErrValidation bool
	}{
		{
			name:              "spec",
			spec:              vimRegistrationJson,
			wantErr:           false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              vimRegistrationYaml,
			wantErr:           false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              vimRegistrationYamlIp,
			wantErr:           false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              vimRegistrationYamlIpOnly,
			wantErr:           false,
			wantErrValidation: true,
		},
		{
			name:              "spec",
			spec:              vimRegistrationYamlInvalid,
			wantErr:           true,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              vimRegistrationInvalidJson,
			wantErr:           true,
			wantErrValidation: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			spec, err := SpecCloudProvider{}.SpecsFromString(tt.spec)
			if tt.wantErr && err == nil {
				t.Errorf("Test failed must not return error")
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if spec == nil {
				t.Errorf("SpecsFromFile() return nil spec")
				return
			}

			clusterSpec, ok := (*spec).(*SpecCloudProvider)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = clusterSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromString() Test failed validator return error for positive case err %v", err)
				return
			}
		})
	}
}

// Read cloud provider spec from file
func TestSpecCloudProvider_SpecsFromFile(t *testing.T) {

	tests := []struct {
		name          string
		file          string
		wantErr       bool
		wantValidaErr bool
	}{
		{
			name:    "Read cluster workload spec from yaml",
			file:    "/provider/positive/provider.yaml",
			wantErr: false,
		},
		{
			name:          "No kind",
			file:          "/provider/negative/no_kind.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
		{
			name:          "Wrong kind",
			file:          "/provider/negative/wrong_kind.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
		{
			name:          "no username",
			file:          "/provider/negative/no_username.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
		{
			name:          "no password",
			file:          "/provider/negative/no_password.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
		{
			name:          "no password",
			file:          "/provider/negative/empty_password.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file

			spec, err := SpecCloudProvider{}.SpecsFromFile(fileName)
			if tt.wantErr && err == nil {
				t.Errorf("Test failed must not return error")
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if spec == nil {
				t.Errorf("SpecsFromFile() return nil spec")
				return
			}
			providerSpec, ok := (*spec).(*SpecCloudProvider)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = providerSpec.Validate()
			if tt.wantValidaErr && err == nil {
				t.Errorf("SpecsFromString() failed spec validator return no error for negative case %v", err)
				return
			}
			if tt.wantValidaErr && err != nil {
				assert.Equal(t, false, providerSpec.IsValid())
				t.Log(err)
				return
			}

			assert.Equal(t, true, providerSpec.IsValid())
		})
	}
}

// Read cluster spec , creates reader and test
// validation
func TestSpecCloudProvider_SpecsFromReader(t *testing.T) {

	tests := []struct {
		name          string
		file          string
		wantErr       bool
		wantValidaErr bool
	}{
		{
			name:    "Read cluster workload spec from yaml",
			file:    "/provider/positive/provider.yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file
			f, err := os.Open(fileName)
			assert.NoError(t, err)

			spec, err := SpecCloudProvider{}.SpecsFromReader(f)
			if tt.wantErr && err == nil {
				t.Errorf("Test failed must not return error")
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if spec == nil {
				t.Errorf("SpecsFromFile() return nil spec")
				return
			}
			providerSpec, ok := (*spec).(*SpecCloudProvider)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = providerSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() Test failed validator "+
					"return error for positive case err %v file %s", err, fileName)
				return
			}
		})
	}
}

var vimRegistrationJson = `{
"kind": "provider",
"hcxCloudUrl": "https://tca-cp03.cnfdemo.io",
"vimName": "core",
"tenantName": "",
"username": "administrator@vsphere.local",
"password": "VMware1!"
}
`

var vimRegistrationYaml = `
kind: provider
hcxCloudUrl: https://tca-cp03.cnfdemo.io
vimName: core
tenantName: ""
username: administrator@vsphere.local
password: VMware1!`

var vimRegistrationYamlIp = `
kind: provider
hcxCloudUrl: https://1.1.1.1
vimName: core
tenantName: ""
username: administrator@vsphere.local
password: VMware1!`

var vimRegistrationYamlIpOnly = `
kind: provider
hcxCloudUrl: 1.1.1.1
vimName: core
tenantName: ""
username: administrator@vsphere.local
password: VMware1!`

var vimRegistrationYamlInvalid = `
kind: provider
     hcxCloudUrl: 1.1.1.1
vimName: core
tenantName: ""
username: administrator@vsphere.local
password: VMware1!`

var vimRegistrationInvalidJson = `{
"kind": "provider"
"hcxCloudUrl": "https://tca-cp03.cnfdemo.io",
"vimName": "core",
"tenantName": "",
"username": "administrator@vsphere.local",
"password": "VMware1!"
}
`
