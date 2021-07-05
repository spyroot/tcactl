package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTcaApiCreateExtension(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		spec         string
		wantErr      bool
		wantDel      bool
		password     string
		verifyAttach bool
	}{
		{
			name:    "Create harbor extension from string",
			rest:    rest,
			spec:    testHarborExt,
			wantErr: false,
			wantDel: true,
		},
		{
			name:    "Create harbor extension from string no delete",
			rest:    rest,
			spec:    testHarborExt,
			wantErr: false,
			wantDel: false,
		},
		{
			name:    "Create duplicate extension must fail",
			rest:    rest,
			spec:    testHarborExt,
			wantErr: true,
		},
		{
			name:     "Create harbor extension with attach",
			rest:     rest,
			spec:     testHarborWithVim,
			wantErr:  false,
			wantDel:  false,
			password: os.Getenv("TCA_REPO_PASSWORD"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf(tt.password)

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			spec, err := request.ExtensionSpecFromFromString(tt.spec)
			assert.NoError(t, err)

			if len(tt.password) > 0 {
				spec.AccessInfo.Password = tt.password
			}

			got, err := a.CreateExtension(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestTcaApiCreateExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !IsValidUUID(got) {
					t.Errorf("TestTcaApiCreateExtension() failed create extension must return UUID")
					return
				}
			}

			extId, err := a.ResolveExtensionId(spec.Name)
			if err != nil {
				return
			}

			if !tt.wantErr {
				assert.Equal(t, got, extId)
			} else {
				assert.Equal(t, "", got)
			}

			t.Logf("extension id %s", got)

			if tt.wantDel {
				_, err := a.DeleteExtension(got)
				assert.NoError(t, err)
			}
		})
	}
}

// TestTcaApiExtensionQuery
func TestTcaApiExtensionQuery(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		wantErr bool
	}{
		{
			name:    "Basic Parser",
			rest:    rest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := &TcaApi{
				rest: tt.rest,
			}

			got, err := a.ExtensionQuery()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtensionQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) && tt.wantErr == false {
				t.Errorf("ExtensionQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) && tt.wantErr == true {
				return
			}

			if len(got.ExtensionsList) == 0 && tt.wantErr == false {
				t.Errorf("ExtensionQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTcaApi_ExtensionQueryFindRepo(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		wantErr bool
	}{
		{
			name:    "Basic Find",
			rest:    rest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			got, err := a.ExtensionQuery()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtensionQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) && tt.wantErr == false {
				t.Errorf("ExtensionQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) && tt.wantErr == true {
				return
			}

			assert.NotNil(t, got)

			valid, err := got.FindRepo("9d0d4ff4-1963-4d89-ac15-2d856768deeb")
			if tt.wantErr && err != nil {
				return
			}

			if tt.wantErr == false && err != nil {
				t.Errorf("ExtensionQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtensionQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(valid.Name) == 0 && tt.wantErr == false {
				t.Errorf("ExtensionQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

var testHarborExt = `{
  "kind": "extensions",
  "name": "repo",
  "version": "v2.x",
  "type": "Repository",
  "extensionKey": "",
  "extensionSubtype": "Harbor",
  "products": [],
  "vimInfo": [],
  "interfaceInfo": {
    "url": "https://1.1.1.1",
    "description": "",
    "trustedCertificate": ""
  },
  "additionalParameters": {
    "trustAllCerts": true
  },
  "autoScaleEnabled": true,
  "autoHealEnabled": true,
  "accessInfo": {
    "username": "admin",
    "password": "Vk13YXJlMSE="
  }
}`

var testHarborMinSpec = `{
  "kind": "extensions",
  "name": "min",
  "version": "v2.x",
  "type": "Repository",
  "extensionSubtype": "Harbor",
  "interfaceInfo": {
    "url": "https://4.53.151.180",
  },
  "additionalParameters": {
    "trustAllCerts": true
  },
  "autoScaleEnabled": true,
  "autoHealEnabled": true,
  "accessInfo": {
    "username": "admin",
    "password": "VMware1!"
  }
}`

var testHarborWithVim = `{
  "kind": "extensions",
  "name": "min",
  "version": "v2.x",
  "type": "Repository",
  "extensionKey": "",
  "extensionSubtype": "Harbor",
  "products": [],
  "vimInfo": [
    {
      "vimName": "edge-test01"
    }
  ],
  "interfaceInfo": {
    "url": "https://10.252.80.135",
    "description": "",
    "trustedCertificate": ""
  },
  "additionalParameters": {
    "trustAllCerts": true
  },
  "autoScaleEnabled": true,
  "autoHealEnabled": true,
  "accessInfo": {
    "username": "admin",
    "password": ""
  }
}`
