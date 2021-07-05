package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/pkg/io"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetVim tests basic retrieval vim
//and vims tenant details
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
			rest:        rest,
			vimErr:      false,
			dumpJson:    true,
			provideName: getTestClusterName(),
			tenantName:  getTenantId(),
		},
		{
			name:        "Get vim and tenant by name provider k8s",
			rest:        rest,
			vimErr:      false,
			dumpJson:    false,
			provideName: getTestClusterName(),
			tenantName:  getTestClusterName(),
		},
		{
			name:        "Get vim and tenant by name provider vc",
			rest:        rest,
			vimErr:      false,
			dumpJson:    false,
			provideName: getTestCloudProvider(),
		},
		{
			name:        "Get invalid vim",
			rest:        rest,
			vimErr:      true,
			dumpJson:    false,
			provideName: "invalid",
		},
		{
			name:        "Get invalid test",
			rest:        rest,
			vimErr:      false,
			errTenant:   true,
			dumpJson:    false,
			provideName: getTestCloudProvider(),
			tenantName:  "invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			if a == nil {
				t.Errorf("clinet must not be nil")
				return
			}

			if tt.reset {
				a.rest = nil
			}

			provider, err := a.GetVim(tt.provideName)
			if err != nil != tt.vimErr {
				t.Errorf("GetVimComputeClusters() error = %v, vimErr %v", err, tt.vimErr)
				return
			}

			io.PrettyPrint(provider)

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

			//			io.PrettyPrint(provider)
		})
	}
}

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
			rest:              rest,
			wantErr:           false,
			dumpJson:          true,
			cloudProviderName: getTestCloudProvider(),
		},
		{
			name:              "Basic error invalid cloud provider",
			rest:              rest,
			wantErr:           true,
			dumpJson:          false,
			cloudProviderName: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			if a == nil {
				t.Errorf("clinet must not be nil")
				return
			}

			if tt.reset {
				a.rest = nil
			}

			vimCompute, err := a.GetVimComputeClusters(tt.cloudProviderName)
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
				t.Errorf("CreateClusterTemplate() error is nil, vimErr %v", tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				t.Logf("Recieved correct error %v", err)
				return
			}
		})
	}
}

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
			rest:              rest,
			wantErr:           false,
			dumpJson:          true,
			cloudProviderName: getTestCloudProvider(),
		},
		{
			name:              "Basic error case bogus cloud provider",
			rest:              rest,
			wantErr:           true,
			dumpJson:          false,
			cloudProviderName: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)

			if tt.reset {
				a.rest = nil
			}

			got, err := a.GetVimNetworks(tt.cloudProviderName)
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

			if tt.wantErr && got == nil {
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
			rest:              rest,
			wantErr:           false,
			dumpJson:          true,
			cloudProviderName: getTestCloudProvider(),
		},
		{
			name:              "Basic error case bogus cloud provider",
			rest:              rest,
			wantErr:           true,
			dumpJson:          false,
			cloudProviderName: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			if tt.reset {
				a.rest = nil
			}

			got, err := a.GetVimNetworks(tt.cloudProviderName)

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

			if tt.wantErr && got == nil {
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
