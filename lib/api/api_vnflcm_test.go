// Package app
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
package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

//TestTcaGetAllPackages Basic catalog getter test
func TestGetAllPackages(t *testing.T) {

	tests := []struct {
		rest    *client.RestClient
		name    string
		wantErr bool
		vduName string
	}{
		{
			name:    "Get all packages",
			rest:    getAuthenticatedClient(),
			wantErr: false,
			vduName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			entireCatalog, err := a.GetAllPackages()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackages() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			assert.NotEqual(t, 0, len(entireCatalog.CnfLcms))
		})
	}
}

//TestGetByCatalogName get entire catalog and search for unit_test catalog entities.
func TestGetByCatalogName(t *testing.T) {

	tests := []struct {
		rest    *client.RestClient
		name    string
		wantErr bool
		vduName string
	}{
		{
			name:    "Get by catalog name",
			rest:    getAuthenticatedClient(),
			wantErr: false,
			vduName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			entireCatalog, err := a.GetAllPackages()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackages() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			if entireCatalog == nil {
				t.Errorf("GetAllPackages() shouldn't return nil")
				return
			}

			catalogEntities, err := entireCatalog.GetByCatalogName(getTestCatalogName())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackages() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			assert.NotEqual(t, 0, len(catalogEntities))
		})
	}
}

//TestTcaGetAllPackages Basic catalog getter test
func TestResolveFromName(t *testing.T) {

	tests := []struct {
		rest    *client.RestClient
		name    string
		wantErr bool
		vduName string
	}{
		{
			name:    "Resolve by name unit test catalog",
			rest:    getAuthenticatedClient(),
			wantErr: false,
			vduName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			entireCatalog, err := a.GetAllPackages()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackages() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			if entireCatalog == nil {
				t.Errorf("GetAllPackages() shouldn't return nil")
				return
			}

			catalogEntities, err := entireCatalog.GetByCatalogName(getTestCatalogName())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackages() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			for _, entity := range catalogEntities {
				name, err := entireCatalog.ResolveFromName(entity.CID)
				if (err != nil) != tt.wantErr {
					t.Errorf("ResolveFromName() error = %v, wantOnGetErr %v", err, tt.wantErr)
					return
				}
				if name == nil {
					t.Errorf("ResolveFromName() should not return nil")
					return
				}
				assert.NotEqual(t, 0, len(name.VnfCatalogName))
			}
			assert.NotEqual(t, 0, len(catalogEntities))
		})
	}
}

// Create catalog entity
func TestCreateCatalogEntity(t *testing.T) {

	tests := []struct {
		rest    *client.RestClient
		name    string
		wantErr bool
		vduName string
	}{
		{
			name:    "Get all packages shouldn't fail",
			rest:    getAuthenticatedClient(),
			wantErr: false,
			vduName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dir := GetTestAssetsDir()
			t.Log(dir)

			//a := getTcaApi(t, rest, false)
			//got, err := a.CreateCatalogEntity()
			//
			//if (err != nil) != tt.wantOnGetErr {
			//	t.Errorf("GetAllPackages() error = %v, wantOnGetErr %v", err, tt.wantOnGetErr)
			//	return
			//}

			//assert.NotNil(t, got)
		})
	}
}
