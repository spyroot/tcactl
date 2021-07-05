package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	ioutils "github.com/spyroot/tcactl/pkg/io"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestTcaApiCreateExtension(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		spec         string
		wantErr      bool
		wantDel      bool
		verifyAttach bool
		password     string
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
			name:         "Create harbor extension with vim attach",
			rest:         rest,
			spec:         testHarborWithVim,
			wantErr:      false,
			wantDel:      true,
			verifyAttach: true,
			password:     os.Getenv("TCA_REPO_PASSWORD"),
		},
		{
			name:         "Create harbor extension with wrong vim",
			rest:         rest,
			spec:         testHarborWithInvalidVim,
			wantErr:      true,
			wantDel:      true,
			verifyAttach: true,
			password:     os.Getenv("TCA_REPO_PASSWORD"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			spec, err := request.ExtensionSpecFromFromString(tt.spec)
			assert.NoError(t, err)

			if len(tt.password) > 0 {
				spec.AccessInfo.Password = tt.password
			}

			got, err := a.CreateExtension(spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestTcaApiCreateExtension() error = %v, vimErr %v", err, tt.wantErr)
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

			if tt.verifyAttach {
				assert.NotNil(t, spec.VimInfo)
				for _, info := range spec.VimInfo {
					t.Log(info.VimId)
					assert.NotEmpty(t, info.VimId)
					assert.NotEmpty(t, info.VimSystemUUID)
					assert.NotEmpty(t, info.VimName)
				}
			}

			time.Sleep(1 * time.Second)

			if tt.wantDel {
				_, err := a.DeleteExtension(got)
				assert.NoError(t, err)
			}

		})
	}
}

func TestTcaApiGetExtension(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		spec         string
		eid          string
		wantErr      bool
		delAfter     bool
		addBefore    bool
		verifyAttach bool
		password     string
	}{
		{
			name:      "Get extension by wrong name",
			rest:      rest,
			eid:       "not found",
			wantErr:   true,
			delAfter:  false,
			addBefore: false,
		},
		{
			name:      "Get extension by name",
			rest:      rest,
			spec:      testGetExtention,
			eid:       "gettest",
			wantErr:   false,
			delAfter:  true,
			addBefore: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			spec, err := request.ExtensionSpecFromFromString(tt.spec)
			assert.NoError(t, err)

			if tt.addBefore {
				// set password if needed
				if len(tt.password) > 0 {
					spec.AccessInfo.Password = tt.password
				}

				// create extension if needed
				eid, err := api.CreateExtension(spec)
				if (err != nil) != tt.wantErr {
					t.Errorf("TestTcaApiGetExtension() error = %v, vimErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr && !IsValidUUID(eid) {
					t.Errorf("TestTcaApiGetExtension() failed create extension must return UUID")
					return
				}
			}

			extensions, err := api.GetExtension(tt.eid)
			if tt.wantErr && err != nil {
				return
			}

			extension, err := extensions.FindExtension(tt.eid)
			if !tt.wantErr {
				assert.NoError(t, err)
				assert.NotNil(t, extension)
				assert.Equal(t, spec.Name, extension.Name)
			} else {
				assert.Error(t, err)
			}

			if tt.addBefore && tt.delAfter {
				_, err := api.DeleteExtension(tt.eid)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTcaApiCreateUpdate(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		specString   string
		vimName      string
		extName      string
		errOnCreate  bool
		errOnUpdate  bool
		addBefore    bool
		delAfter     bool
		verifyAttach bool
		password     string
	}{
		{
			name:        "Create, Attach harbor extension to cluster",
			rest:        rest,
			specString:  testHarborCreateUpdate,
			vimName:     getTestClusterName(),
			extName:     "min",
			addBefore:   false,
			errOnCreate: false,
			errOnUpdate: false,
			delAfter:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				extensions *response.Extensions
				err        error
			)

			api, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			spec, err := request.ExtensionSpecFromFromString(tt.specString)
			assert.NoError(t, err)

			if len(tt.password) > 0 {
				spec.AccessInfo.Password = tt.password
			}

			gotEid := tt.extName
			if tt.addBefore {
				gotEid, err := api.CreateExtension(spec)
				if (err != nil) != tt.errOnCreate {
					t.Errorf("TestTcaApiCreateUpdate() error = %v, vimErr %v", err, tt.errOnCreate)
					return
				}
				if !tt.errOnCreate && !IsValidUUID(gotEid) {
					t.Errorf("TestTcaApiCreateUpdate() failed create extension must return UUID")
					return
				}
			}

			// lookup ext id if needed
			if !tt.errOnCreate {
				extensions, err = api.GetExtension(gotEid)
				assert.NoError(t, err)
				extension, err := extensions.FindExtension(gotEid)
				assert.NoError(t, err)
				assert.NotNil(t, extension)
			}

			// Update and add vim specString and update
			time.Sleep(1 * time.Second)
			spec.AddVim(tt.vimName)
			_, err = spec.GetVim(tt.vimName)
			assert.NoError(t, err)
			ioutils.YamlPrinter(spec)

			//t.Log(spec.VimInfo[0].VimId)
			//t.Log(spec.VimInfo[0].VimName)
			//t.Log(spec.VimInfo[0].VimSystemUUID)
			//
			//_, err = api.UpdateExtension(spec, gotEid)
			//if !tt.errOnUpdate {
			//	assert.NoError(t, err)
			//}
			//
			//if tt.delAfter {
			//	_, err := api.DeleteExtension(got)
			//	assert.NoError(t, err)
			//}

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
				t.Errorf("ExtensionQuery() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) && tt.wantErr == false {
				t.Errorf("ExtensionQuery() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) && tt.wantErr == true {
				return
			}

			if len(got.ExtensionsList) == 0 && tt.wantErr == false {
				t.Errorf("ExtensionQuery() error = %v, vimErr %v", err, tt.wantErr)
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
				t.Errorf("ExtensionQuery() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) && tt.wantErr == false {
				t.Errorf("ExtensionQuery() error = %v, vimErr %v", err, tt.wantErr)
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
				t.Errorf("ExtensionQuery() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtensionQuery() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if len(valid.Name) == 0 && tt.wantErr == false {
				t.Errorf("ExtensionQuery() error = %v, vimErr %v", err, tt.wantErr)
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

var testHarborWithInvalidVim = `{
  "kind": "extensions",
  "name": "min",
  "version": "v2.x",
  "type": "Repository",
  "extensionKey": "",
  "extensionSubtype": "Harbor",
  "products": [],
  "vimInfo": [
    {
      "vimName": "fail"
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

var testHarborCreateUpdate = `{
  "kind": "extensions",
  "name": "min",
  "version": "v2.x",
  "type": "Repository",
  "extensionKey": "",
  "extensionSubtype": "Harbor",
  "products": [],
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
    "password": "test"
  }
}`

var testGetExtention = `{
  "kind": "extensions",
  "name": "gettest",
  "version": "v2.x",
  "type": "Repository",
  "extensionKey": "",
  "extensionSubtype": "Harbor",
  "products": [],
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
    "password": "test"
  }
}`
