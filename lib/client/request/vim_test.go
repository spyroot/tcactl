package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/spyroot/tcactl/lib/testutil"
	"io"
	"strings"
	"testing"
)

func TestReadProviderSpec(t *testing.T) {

	tests := []struct {
		name              string
		b                 io.Reader
		wantErrReader     bool
		wantErrValidation bool
	}{
		{
			name:              "spec",
			b:                 strings.NewReader(vimRegistrationJson),
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			b:                 strings.NewReader(vimRegistrationYaml),
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			b:                 strings.NewReader(vimRegistrationYamlIp),
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			b:                 strings.NewReader(vimRegistrationYamlIpOnly),
			wantErrReader:     false,
			wantErrValidation: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ReadProviderSpec(tt.b)
			if (err != nil) != tt.wantErrReader {
				t.Errorf("ReadProviderSpec() error = %v, wantErr %v", err, tt.wantErrReader)
				return
			}

			specValidator := validator.New()
			err = specValidator.Struct(got)
			if tt.wantErrValidation {
				if err == nil {
					t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
			if !tt.wantErrValidation {
				if err != nil {
					t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
		})
	}
}

func TestProviderSpecsFromFromString(t *testing.T) {

	tests := []struct {
		name              string
		spec              string
		wantErrReader     bool
		wantErrValidation bool
	}{
		{
			name:              "spec",
			spec:              vimRegistrationJson,
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              vimRegistrationYaml,
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              vimRegistrationYamlIp,
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              vimRegistrationYamlIpOnly,
			wantErrReader:     false,
			wantErrValidation: true,
		},
		{
			name:              "spec",
			spec:              vimRegistrationYamlInvalid,
			wantErrReader:     true,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              vimRegistrationInvalidJson,
			wantErrReader:     true,
			wantErrValidation: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ProviderSpecsFromFromString(tt.spec)
			if tt.wantErrReader {
				if err == nil {
					t.Errorf("ReadProviderSpec() error = %v, wantErr %v", err, tt.wantErrReader)
					return
				}
				return
			}
			if !tt.wantErrReader && err != nil {
				t.Errorf("ReadProviderSpec() error = %v, wantErr %v", err, tt.wantErrReader)
				return

			}

			specValidator := validator.New()
			err = specValidator.Struct(got)
			if tt.wantErrValidation {
				if err == nil {
					t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
			if !tt.wantErrValidation {
				if err != nil {
					t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
		})
	}
}

func TestProviderSpecsFromFromFile(t *testing.T) {

	tests := []struct {
		name              string
		spec              string
		wantErrReader     bool
		wantErrValidation bool
	}{
		{
			name:              "spec",
			spec:              testutil.SpecTempFileName(vimRegistrationJson),
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              testutil.SpecTempFileName(vimRegistrationYaml),
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              testutil.SpecTempFileName(vimRegistrationYamlIp),
			wantErrReader:     false,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              testutil.SpecTempFileName(vimRegistrationYamlIpOnly),
			wantErrReader:     false,
			wantErrValidation: true,
		},
		{
			name:              "spec",
			spec:              testutil.SpecTempFileName(vimRegistrationYamlInvalid),
			wantErrReader:     true,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              testutil.SpecTempFileName(vimRegistrationInvalidJson),
			wantErrReader:     true,
			wantErrValidation: false,
		},
		{
			name:              "spec",
			spec:              testutil.SpecTempFileName("wrong_file"),
			wantErrReader:     true,
			wantErrValidation: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ProviderSpecsFromFile(tt.spec)
			if tt.wantErrReader {
				if err == nil {
					t.Errorf("ReadProviderSpec() error = %v, wantErr %v", err, tt.wantErrReader)
					return
				}
				return
			}
			if !tt.wantErrReader && err != nil {
				t.Errorf("ReadProviderSpec() error = %v, wantErr %v", err, tt.wantErrReader)
				return

			}

			specValidator := validator.New()
			err = specValidator.Struct(got)
			if tt.wantErrValidation {
				if err == nil {
					t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
			}
			if !tt.wantErrValidation {
				if err != nil {
					t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErrValidation)
					return
				}
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
