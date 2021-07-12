package specs

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

//Read spec from file and test
func TestExtensionSpecsFromFile(t *testing.T) {

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name: "Read instance workload spec from yaml",
			file: "/extension/positive/harbor.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file
			spec, err := SpecExtension{}.SpecsFromFile(fileName)
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

			extensionSpec, ok := (*spec).(*SpecExtension)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = extensionSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() Test failed validator "+
					"return error for positive case err %v file %s", err, fileName)
				return
			}
		})
	}
}

// Read spec from string and test
func TestSpecExtension_SpecsFromString(t *testing.T) {

	tests := []struct {
		name          string
		file          string
		wantErr       bool
		wantValidaErr bool
	}{
		{
			name:    "Read instance spec from yaml",
			file:    "/extension/positive/harbor.yaml",
			wantErr: false,
		},
		{
			name:          "No kind",
			file:          "/extension/negative/no_kind.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
		{
			name:          "Wrong kind",
			file:          "/extension/negative/wrong_kind.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
		{
			name:          "No name",
			file:          "/extension/negative/no_name.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
		{
			name:          "No version",
			file:          "/extension/negative/no_ver.yaml",
			wantErr:       false,
			wantValidaErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file

			b, err := ioutil.ReadFile(fileName)
			assert.NoError(t, err)

			spec, err := SpecExtension{}.SpecsFromString(string(b))
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
			extensionSpec, ok := (*spec).(*SpecExtension)
			if !ok {
				t.Errorf("SpecsFromString() failed method return wrong type")
				return
			}
			err = extensionSpec.Validate()
			if tt.wantValidaErr && err == nil {
				t.Errorf("SpecsFromString() failed spec validator return no error for negative case %v", err)
				return
			}
			if tt.wantValidaErr && err != nil {
				assert.Equal(t, false, extensionSpec.IsValid())
				t.Log(err)
				return
			}

			assert.Equal(t, true, extensionSpec.IsValid())
		})
	}
}

var testValidHarborJson = `{
  "kind": "extensions",
  "name": "repo",
  "version": "v2.x",
  "type": "Repository",
  "extensionKey": "",
  "extensionSubtype": "Harbor",
  "products": [],
  "vimInfo": [],
  "interfaceInfo": {
    "url": "https://1.1.1.1",
    "description": "",
    "trustedCertificate": ""
  },
  "additionalParameters": {
    "trustAllCerts": true
  },
  "autoScaleEnabled": true,
  "autoHealEnabled": true,
  "accessInfo": {
    "username": "admin",
    "password": "Vk13YXJlMSE="
  }
}`

var testValidHarborYaml = `
kind: extensions
name: repo
version: v2.x
type: Repository
extensionKey: ""
extensionSubtype: Harbor
products: []
vim_info: []
interfaceInfo:
  url: https://1.1.1.1
  description: ""
  trustedCertificate: ""
additionalParameters:
  trustAllCerts: true
autoScaleEnabled: true
autoHealEnabled: true
accessInfo:
  username: admin
  password: Vk13YXJlMSE=`
