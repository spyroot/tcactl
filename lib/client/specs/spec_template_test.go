package specs

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestSpecClusterTemplate_SpecsFromFile(t *testing.T) {

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Read cluster template spec from yaml",
			file:    "/template/positive/template.mgmt.yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file
			spec, err := SpecClusterTemplate{}.SpecsFromFile(fileName)
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
			templateSpec, ok := (*spec).(*SpecClusterTemplate)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}
			err = templateSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() Test failed validator "+
					"return error for positive case err %v file %s", err, fileName)
				return
			}
		})
	}
}

func TestClusterTemplateSpec_SpecsFromReader(t *testing.T) {

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Read cluster template spec from yaml",
			file:    "/template/positive/template.mgmt.yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file
			f, err := os.Open(fileName)
			assert.NoError(t, err)

			spec, err := SpecClusterTemplate{}.SpecsFromReader(f)
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

			clusterSpec, ok := (*spec).(*SpecClusterTemplate)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = clusterSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() Test failed validator "+
					"return error for positive case err %v file %s", err, fileName)
				return
			}

		})
	}
}

func TestClusterTemplateSpec_SpecsFromString(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		wantErr       bool
		wantValidaErr bool
	}{
		{
			name:    "Read instance spec from yaml",
			file:    "/template/positive/template.mgmt.yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file

			b, err := ioutil.ReadFile(fileName)
			assert.NoError(t, err)

			spec, err := SpecClusterTemplate{}.SpecsFromString(string(b))
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
			instanceSpec, ok := (*spec).(*SpecClusterTemplate)
			if !ok {
				t.Errorf("SpecsFromString() failed method return wrong type")
				return
			}
			err = instanceSpec.Validate()
			if tt.wantValidaErr && err == nil {
				t.Errorf("SpecsFromString() failed spec validator return no error for negative case %v", err)
				return
			}
			if tt.wantValidaErr && err != nil {
				t.Log(err)
				return
			}
		})
	}
}
