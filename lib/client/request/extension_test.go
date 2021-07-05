package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/spyroot/tcactl/lib/testutil"
	"io"
	"testing"
)

// TestExtensionSpecFromFromString
func TestExtensionSpecFromFromString(t *testing.T) {

	tests := []struct {
		name              string
		specString        string
		wantErr           bool
		wantErrValidation bool
	}{
		{
			name:              "Read Json spec",
			specString:        testValidHarborJson,
			wantErr:           false,
			wantErrValidation: false,
		},
		{
			name:              "Read Json spec",
			specString:        testValidHarborYaml,
			wantErr:           false,
			wantErrValidation: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Read spec
			got, err := ExtensionSpecFromFromString(tt.specString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtensionSpecFromFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			specValidator := validator.New()
			err = specValidator.Struct(got)
			if tt.wantErrValidation {
				if err == nil {
					t.Errorf("ExtensionSpecFromFromString() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
			if !tt.wantErrValidation {
				if err != nil {
					t.Errorf("ExtensionSpecFromFromString() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
		})
	}
}

//TestExtensionSpecsFromFile
func TestExtensionSpecsFromFile(t *testing.T) {

	tests := []struct {
		name              string
		fileName          string
		wantErr           bool
		wantErrValidation bool
	}{
		{
			name:              "Read Json spec",
			fileName:          testutil.SpecTempFileName(testValidHarborJson),
			wantErr:           false,
			wantErrValidation: false,
		},
		{
			name:              "Read Json spec",
			fileName:          testutil.SpecTempFileName(testValidHarborYaml),
			wantErr:           false,
			wantErrValidation: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Read spec
			got, err := ExtensionSpecsFromFile(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtensionSpecsFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			specValidator := validator.New()
			err = specValidator.Struct(got)
			if tt.wantErrValidation {
				if err == nil {
					t.Errorf("ExtensionSpecsFromFile() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
			if !tt.wantErrValidation {
				if err != nil {
					t.Errorf("ExtensionSpecsFromFile() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
		})
	}
}

func TestReadExtensionSpec(t *testing.T) {
	tests := []struct {
		name              string
		r                 io.Reader
		wantErr           bool
		wantErrValidation bool
	}{
		{
			name:              "Read Json spec",
			r:                 testutil.SpecTempReader(testValidHarborJson),
			wantErr:           false,
			wantErrValidation: false,
		},
		{
			name:              "Read Yaml spec",
			r:                 testutil.SpecTempReader(testValidHarborYaml),
			wantErr:           false,
			wantErrValidation: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Read spec
			got, err := ReadExtensionSpec(tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadExtensionSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			specValidator := validator.New()
			err = specValidator.Struct(got)
			if tt.wantErrValidation {
				if err == nil {
					t.Errorf("ExtensionSpecsFromFile() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
			if !tt.wantErrValidation {
				if err != nil {
					t.Errorf("ExtensionSpecsFromFile() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
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
