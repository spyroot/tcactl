package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTcaApiGetEntireCatalog(t *testing.T) {
	//
	tests := []struct {
		rest    *client.RestClient
		name    string
		wantErr bool
		vduName string
	}{
		{
			name:    "Get entire catalog shouldn't fail",
			rest:    rest,
			wantErr: false,
			vduName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			got, err := api.GetEntireCatalog()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVnfPkgm() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			assert.NotNil(t, got)
		})
	}
}

// TestTcaApiGetVnfPkgm
func TestTcaApiGetVnfPkgm(t *testing.T) {
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

			api, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			catalog, err := api.GetEntireCatalog()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVnfPkgm() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			if catalog == nil {
				t.Errorf("GetVnfPkgm() shouldn't return nil")
				return
			}

			for _, p := range catalog.Packages {
				pkgm, err := api.GetVnfPkgm("", p.PID)
				if err != nil {
					return
				}

				assert.NotNil(t, pkgm.Packages)
			}
		})
	}
}

// TestTcaApiGetVnfPkgm
func TestTcaApiGetCatalogAndVdu(t *testing.T) {
	tests := []struct {
		rest          *client.RestClient
		name          string
		wantErr       bool
		isUserdefined bool
		isPackageid   bool
	}{
		{
			name:          "Get all packages shouldn't fail",
			rest:          rest,
			wantErr:       false,
			isUserdefined: true,
			isPackageid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			catalog, err := api.GetEntireCatalog()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVnfPkgm() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}

			if catalog == nil {
				t.Errorf("GetVnfPkgm() shouldn't return nil")
				return
			}

			for _, p := range catalog.Packages {
				// search by catalog user defined name
				if tt.isUserdefined {
					pkg, vdu, err2 := api.GetCatalogAndVdu(p.UserDefinedData.Name)
					if !tt.wantErr {
						assert.NoError(t, err2)
						assert.NotNil(t, pkg)
						assert.NotNil(t, vdu)
					}
				}
			}
		})
	}
}
