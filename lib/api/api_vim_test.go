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
	"context"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/pkg/io"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetVim tests basic retrieval vim
// and vims tenant details
func TestGetVim(t *testing.T) {

	tests := []struct {
		name        string
		rest        *client.RestClient
		vimErr      bool
		errTenant   bool // invalid tenant
		reset       bool
		dumpJson    bool
		provideName string
		tenantName  string
	}{
		{
			name:        "Get vim and tenant by id",
			rest:        getAuthenticatedClient(),
			vimErr:      false,
			dumpJson:    true,
			provideName: getTestWorkloadClusterName(),
			tenantName:  getTenantId(),
		},
		{
			name:        "Get vim and tenant by name provider k8s",
			rest:        getAuthenticatedClient(),
			vimErr:      false,
			dumpJson:    false,
			provideName: getTestWorkloadClusterName(),
			tenantName:  getTestWorkloadClusterName(),
		},
		{
			name:        "Get vim and tenant by name provider vc",
			rest:        getAuthenticatedClient(),
			vimErr:      false,
			dumpJson:    false,
			provideName: getTestCloudProvider(),
		},
		{
			name:        "Get invalid vim",
			rest:        getAuthenticatedClient(),
			vimErr:      true,
			dumpJson:    false,
			provideName: "invalid",
		},
		{
			name:        "Get invalid test",
			rest:        getAuthenticatedClient(),
			vimErr:      false,
			errTenant:   true,
			dumpJson:    false,
			provideName: getTestCloudProvider(),
			tenantName:  "invalid",
		},
		{
			name:        "Get empty name test",
			rest:        getAuthenticatedClient(),
			vimErr:      true,
			errTenant:   false,
			dumpJson:    false,
			provideName: "",
			tenantName:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			a := getTcaApi(t, tt.rest, false)

			provider, err := a.GetVim(ctx, tt.provideName)
			if err != nil != tt.vimErr {
				t.Errorf("GetVim() error = %v, vimErr %v", err, tt.vimErr)
				return
			}
			if tt.vimErr && err != nil {
				return
			}

			if !tt.vimErr {
				if provider == nil {
					t.Errorf("provide must not be nil")
					return
				}
				if len(tt.provideName) > 0 {
					tenant, err := provider.GetTenant(tt.tenantName)
					if err != nil {
						return
					}

					if tt.errTenant {
						assert.Error(t, err)
						return
					}

					assert.NoError(t, err)
					if tt.dumpJson {
						err := io.PrettyPrint(tenant)
						if err != nil {
							return
						}
					}
				}
			}
		})
	}
}

// test vim computer clusters
func TestTcaApi_GetVimComputeClusters(t *testing.T) {

	tests := []struct {
		name              string
		rest              *client.RestClient
		wantErr           bool
		reset             bool
		dumpJson          bool
		cloudProviderName string
	}{
		{
			name:              "Basic positive get list vim details",
			rest:              getAuthenticatedClient(),
			wantErr:           false,
			dumpJson:          true,
			cloudProviderName: getTestCloudProvider(),
		},
		{
			name:              "Basic error invalid cloud provider",
			rest:              getAuthenticatedClient(),
			wantErr:           true,
			dumpJson:          false,
			cloudProviderName: "test",
		},
		{
			name:              "Empty cloud provider",
			rest:              getAuthenticatedClient(),
			wantErr:           true,
			dumpJson:          false,
			cloudProviderName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			a := getTcaApi(t, tt.rest, false)
			if a == nil {
				t.Errorf("clinet must not be nil")
				return
			}

			if tt.reset {
				a.rest = nil
			}

			vimCompute, err := a.GetVimComputeClusters(ctx, tt.cloudProviderName)
			if err != nil != tt.wantErr {
				t.Errorf("GetVimComputeClusters() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && vimCompute == nil {
				t.Errorf("GetVimComputeClusters() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			// we want err and got not nil instance
			if tt.wantErr && vimCompute != nil {
				t.Errorf("GetVimComputeClusters() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && vimCompute == nil {
				return
			}

			// if we got at least something in list we are ok
			if !tt.wantErr && vimCompute != nil && len(vimCompute.Items) > 0 {
				return
			}

			if vimCompute != nil {
				t.Log(vimCompute.Items)
			}

			if tt.dumpJson {
				err := io.PrettyPrint(vimCompute)
				if err != nil {
					return
				}
			}

			if tt.wantErr && err == nil {
				t.Errorf("GetVimComputeClusters() error is nil, vimErr %v", tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				t.Logf("Recieved correct error %v", err)
				return
			}
		})
	}
}

// Test vim network fetching
func TestTcaApi_GetVimNetworks(t *testing.T) {

	tests := []struct {
		name              string
		rest              *client.RestClient
		wantErr           bool
		reset             bool
		dumpJson          bool
		cloudProviderName string
	}{
		{
			name:              "Basic get list network details",
			rest:              getAuthenticatedClient(),
			wantErr:           false,
			dumpJson:          true,
			cloudProviderName: getTestCloudProvider(),
		},
		{
			name:              "Basic error case bogus cloud provider",
			rest:              getAuthenticatedClient(),
			wantErr:           true,
			dumpJson:          false,
			cloudProviderName: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			a := getTcaApi(t, tt.rest, false)
			if tt.reset {
				a.rest = nil
			}

			got, err := a.GetVimNetworks(ctx, tt.cloudProviderName)
			if err != nil != tt.wantErr {
				t.Errorf("GetVimNetworks() error = %v, vimErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("GetVimNetworks() error = %v, vimErr %v", err, tt.wantErr)
				return
			}
			// we want err and got not nil instance
			if tt.wantErr && got != nil {
				t.Errorf("GetVimNetworks() error = %v, vimErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				return
			}

			if got == nil {
				t.Errorf("GetVimNetworks() return nil with no error")
				return
			}

			// if we got at least something in list we are ok
			if len(got.Network) > 0 {
				return
			}

			if tt.dumpJson {
				err := io.PrettyPrint(got)
				if err != nil {
					return
				}
			}
		})
	}
}

// test vim network
func TestTcaApi_GetVimNetworksAdv(t *testing.T) {

	tests := []struct {
		name              string
		rest              *client.RestClient
		wantErr           bool
		reset             bool
		dumpJson          bool
		cloudProviderName string
	}{
		{
			name:              "Basic get list network details",
			rest:              getAuthenticatedClient(),
			wantErr:           false,
			dumpJson:          true,
			cloudProviderName: getTestCloudProvider(),
		},
		{
			name:              "Basic error case bogus cloud provider",
			rest:              getAuthenticatedClient(),
			wantErr:           true,
			dumpJson:          false,
			cloudProviderName: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			a := getTcaApi(t, tt.rest, false)
			if tt.reset {
				a.rest = nil
			}

			got, err := a.GetVimNetworks(ctx, tt.cloudProviderName)

			if err != nil != tt.wantErr {
				t.Errorf("GetVimNetworks() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("GetVimNetworks() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			// we want err and got not nil instance
			if tt.wantErr && got != nil {
				t.Errorf("GetVimNetworks() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				return
			}

			// if we got at least something in list we are ok
			if !tt.wantErr && len(got.Network) > 0 {
				return
			}

			if got != nil {
				t.Log(got.Network)
			}

			if tt.dumpJson {
				err := io.PrettyPrint(got)
				if err != nil {
					return
				}
			}

			if tt.wantErr && err == nil {
				t.Errorf("GetVimNetworks() error is nil, vimErr %v", tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				t.Logf("Recieved correct error %v", err)
				return
			}
		})
	}
}
