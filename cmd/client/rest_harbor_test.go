// Package client
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
package client

import (
	"github.com/spyroot/hestia/pkg/io"
	"testing"
)

func TestRestClient_GetChart(t *testing.T) {

	tests := []struct {
		name    string
		client  *RestClient
		wantErr bool
		arg     *PackageUpload
	}{
		{
			name:    "get list of helm repos",
			client:  harbor,
			wantErr: false,
			arg:     NewPackageUpload("test_upload24"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.client.GetCharts()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestRestClient_GetChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRestClient_GetRepos(t *testing.T) {

	tests := []struct {
		name    string
		client  *RestClient
		wantErr bool
		arg     *PackageUpload
	}{
		{
			name:    "get list of helm repos",
			client:  harbor,
			wantErr: false,
			arg:     NewPackageUpload("test_upload24"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := tt.client.GetRepos()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestRestClient_GetChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			io.PrettyPrint(r)
		})
	}
}
