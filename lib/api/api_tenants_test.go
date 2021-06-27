package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/pkg/io"
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
				t.Errorf("GetVimNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("GetVimNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// we want err and got not nil instance
			if tt.wantErr && got != nil {
				t.Errorf("GetVimNetworks() error = %v, wantErr %v", err, tt.wantErr)
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
				t.Errorf("GetVimNetworks() error is nil, wantErr %v", tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				t.Logf("Recieved correct error %v", err)
				return
			}
		})
	}
}
