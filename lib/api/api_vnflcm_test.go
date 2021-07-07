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
			name:    "Get all packages shouldn't fail",
			rest:    rest,
			wantErr: false,
			vduName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			got, err := a.GetAllPackages()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackages() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			assert.NotNil(t, got)
		})
	}
}

func TestCreateCatalogEntity(t *testing.T) {

	tests := []struct {
		rest    *client.RestClient
		name    string
		wantErr bool
		vduName string
	}{
		{
			name:    "Get all packages shouldn't fail",
			rest:    rest,
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
