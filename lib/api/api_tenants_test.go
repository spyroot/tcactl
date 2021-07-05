package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/pkg/io"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTenantsCloudProvider(t *testing.T) {

	tests := []struct {
		name              string
		rest              *client.RestClient
		wantErr           bool
		reset             bool
		dumpJson          bool
		cloudProviderName string
		expectedErrorMsg  string
	}{
		{
			name:              "Basic get tenant positive case",
			rest:              rest,
			wantErr:           false,
			dumpJson:          true,
			cloudProviderName: getTenantCluster(),
			expectedErrorMsg:  "",
		},
		{
			name:              "Basic get tenant negative case",
			rest:              rest,
			wantErr:           true,
			dumpJson:          true,
			cloudProviderName: "test123",
			expectedErrorMsg:  "tenant 'test123' not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			if tt.reset {
				a.rest = nil
			}

			got, err := a.TenantsCloudProvider(tt.cloudProviderName)
			if got != nil && tt.dumpJson {
				err := io.PrettyPrint(got)
				if err != nil {
					return
				}
			}

			// for error case we check what we got
			if tt.wantErr && err != nil {
				assert.EqualErrorf(t, err, tt.expectedErrorMsg, "Error should be: %v, got: %v", tt.expectedErrorMsg, err)
			}

			if err != nil != tt.wantErr {
				t.Errorf("TenantsCloudProvider() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("TenantsCloudProvider() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			// we want err and got not nil instance
			if tt.wantErr && got != nil {
				t.Errorf("TenantsCloudProvider() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got == nil {
				return
			}

			// if we got at least something in list we are ok
			if !tt.wantErr && len(got.TenantsList) > 0 {
				return
			}

			if got != nil {
				t.Log(got)
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
