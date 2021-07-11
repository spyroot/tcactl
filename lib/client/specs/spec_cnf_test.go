package specs

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestInstanceRequestSpec_SpecsFromFile(t *testing.T) {

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Read instance workload spec from yaml",
			file:    "/cnf/positive/cnf.spec.yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file
			spec, err := InstanceRequestSpec{}.SpecsFromFile(fileName)
			if tt.wantErr && err == nil {
				t.Errorf("SpecsFromFile() failed must not return error")
				return
			}
			if tt.wantErr && err != nil {
				return
			}

			if spec == nil {
				t.Errorf("SpecsFromFile() return nil spec")
				return
			}

			instanceSpec, ok := (*spec).(*InstanceRequestSpec)
			if !ok {
				t.Errorf("SpecsFromFile() failed method return wrong type")
				return
			}

			err = instanceSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() Test failed validator "+
					"return error for positive case err %v file %s", err, fileName)
				return
			}
		})
	}
}

func TestInstanceRequestSpec_SpecsFromString(t *testing.T) {

	tests := []struct {
		name          string
		file          string
		wantErr       bool
		wantValideErr bool
	}{
		{
			name:    "Read instance spec from yaml",
			file:    "/cnf/positive/cnf.spec.yaml",
			wantErr: false,
		},
		{
			name:          "Read instance spec no kind details",
			file:          "/cnf/negative/no_kind_cnf.spec.yaml",
			wantErr:       false,
			wantValideErr: true,
		},
		{
			name:          "Read instance spec no cluster details",
			file:          "/cnf/negative/no_cluster_name.yaml",
			wantErr:       false,
			wantValideErr: true,
		},
		{
			name:          "Read instance spec no instance name",
			file:          "/cnf/negative/no_instance_name.yaml",
			wantErr:       false,
			wantValideErr: true,
		},
		{
			name:          "Read instance spec no instance name",
			file:          "/cnf/negative/no_catalog_name.yaml",
			wantErr:       false,
			wantValideErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file

			b, err := ioutil.ReadFile(fileName)
			assert.NoError(t, err)

			spec, err := InstanceRequestSpec{}.SpecsFromString(string(b))
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
			instanceSpec, ok := (*spec).(*InstanceRequestSpec)
			if !ok {
				t.Errorf("SpecsFromString() failed method return wrong type")
				return
			}
			err = instanceSpec.Validate()
			if tt.wantValideErr && err == nil {
				t.Errorf("SpecsFromString() failed spec validator return no error for negative case %v", err)
				return
			}
			if tt.wantValideErr && err != nil {
				t.Log(err)
				return
			}
		})
	}
}

var testSpec = `
kind: instance
cloud_name: edge
cluster_name: edge-test01
vim_type: k8s
catalog_name: default
repo_url: https://test
repo_username: admin
repo_password: pass
instance_name: test_instance01
node_pool: default
namespace: default
flavor_name: default
description: default
additionalParams:
    disableGrant: true
    ignoreGrantFailure: true
    disableAutoRollback: true`
