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

//TestTcaApiGetVdu Basic catalog getter test
func TestGetVdu(t *testing.T) {

	tests := []struct {
		rest             *client.RestClient
		name             string
		wantErr          bool
		vduName          string
		simulateFailure  bool
		failureCondition string
	}{
		{
			name:    "Wrong name must fail",
			rest:    getAuthenticatedClient(),
			wantErr: true,
			vduName: "",
		},
		{
			name:    "Wrong name must fail",
			rest:    getAuthenticatedClient(),
			wantErr: true,
			vduName: "abc",
		},
		{
			name:    "Valid user defined name must be resolved",
			rest:    getAuthenticatedClient(),
			wantErr: false,
			vduName: "app",
		},
		{
			name:             "Valid user defined name and simulated failure",
			rest:             getAuthenticatedClient(),
			wantErr:          true,
			vduName:          "app",
			simulateFailure:  true,
			failureCondition: client.TcaVmwareTelcoPackages,
		},
		{
			name:             "Valid user defined name and simulated failure",
			rest:             getAuthenticatedClient(),
			wantErr:          true,
			vduName:          "app",
			simulateFailure:  true,
			failureCondition: client.TcaVmwareTelcoPackages + "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := getTcaApi(t, tt.rest, false)
			if tt.simulateFailure {
				a.rest.SetSimulateFailure(true)
				a.rest.AddFailureCondition(tt.failureCondition)
			}
			defer a.rest.SetSimulateFailure(false)

			got, err := a.GetVdu(tt.vduName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVdu() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.NotEqual(t, 0, len(got.Vdus))
			}
		})
	}
}
