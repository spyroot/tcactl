package api

import (
	"context"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/testutil"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestCreateExtension(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		spec         string
		wantErr      bool
		wantDel      bool
		verifyAttach bool
		password     string
		randRepoName bool
		randUser     bool
	}{
		{
			name:         "Create harbor extension from string",
			rest:         rest,
			spec:         testHarborExt,
			wantErr:      false,
			wantDel:      true,
			randRepoName: true,
			randUser:     true,
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

			a := getTcaApi(t, rest, false)
			spec, err := specs.ExtensionSpecFromFromString(tt.spec)
			assert.NoError(t, err)
			ctx := context.Background()

			if tt.randUser {
				spec.AccessInfo.Username = testutil.RandString(8)
			}
			if tt.randRepoName {
				spec.Name = testutil.RandString(8)
			}
			if len(tt.password) > 0 {
				spec.AccessInfo.Password = tt.password
			}

			got, err := a.CreateExtension(ctx, spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestCreateExtension() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !IsValidUUID(got) {
					t.Errorf("TestCreateExtension() failed create extension must return UUID")
					return
				}
			}

			extId, err := a.ResolveExtensionId(ctx, spec.Name)
			if err != nil {
				return
			}

			if !tt.wantErr {
				assert.Equal(t, got, extId)
			} else {
				assert.Equal(t, "", got)
			}

			if tt.verifyAttach {
				assert.NotNil(t, spec.VimInfo)
				for _, info := range spec.VimInfo {
					assert.NotEmpty(t, info.VimId)
					assert.NotEmpty(t, info.VimSystemUUID)
					assert.NotEmpty(t, info.VimName)
				}
			}

			time.Sleep(1 * time.Second)

			if tt.wantDel {
				_, err := a.DeleteExtension(ctx, got)
				assert.NoError(t, err)
			}

		})
	}
}

func TestGetExtension(t *testing.T) {

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
		randRepoName bool
		randUser     bool
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
			name:         "Get extension by name",
			rest:         rest,
			spec:         testGetExtention,
			eid:          "gettest",
			wantErr:      false,
			delAfter:     true,
			addBefore:    true,
			randUser:     true,
			randRepoName: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			api := getTcaApi(t, rest, false)
			spec, err := specs.ExtensionSpecFromFromString(tt.spec)
			assert.NoError(t, err)

			if tt.addBefore {
				// set password if needed
				if len(tt.password) > 0 {
					spec.AccessInfo.Password = tt.password
				}
				if tt.randUser {
					spec.AccessInfo.Username = testutil.RandString(8)
				}
				if tt.randRepoName {
					spec.Name = testutil.RandString(8)
				}

				// create extension if needed
				eid, err := api.CreateExtension(ctx, spec)
				if (err != nil) != tt.wantErr {
					t.Errorf("TestGetExtension() error = %v, vimErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && !IsValidUUID(eid) {
					t.Errorf("TestGetExtension() failed create extension must return UUID")
					return
				}

				tt.eid = eid
			}

			extensions, err := api.GetExtension(ctx, tt.eid)
			if tt.wantErr && err != nil {
				return
			}

			if tt.wantErr && err == nil {
				t.Errorf("TestGetExtension() must return not nil error")
				return
			}

			if extensions == nil {
				t.Errorf("TestGetExtension() failed create must not be nil")
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
				_, err := api.DeleteExtension(ctx, tt.eid)
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
			name:       "Create and attach harbor extension to cluster",
			rest:       rest,
			specString: testHarborCreateUpdate,
			vimName:    getTestWorkloadClusterName(),
			extName:    getTestRepoName(),
			password:   os.Getenv("TCA_REPO_PASSWORD"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				ctx        = context.Background()
				extensions *response.Extensions
				err        error
			)

			ctx = context.Background()
			api := getTcaApi(t, rest, false)
			spec, err := specs.ExtensionSpecFromFromString(tt.specString)
			assert.NoError(t, err)

			if len(tt.password) > 0 {
				spec.AccessInfo.Password = tt.password
			}

			extensionId := tt.extName
			if tt.addBefore {
				gotEid, err := api.CreateExtension(ctx, spec)
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
			extensions, err = api.GetExtension(ctx, extensionId)
			assert.NoError(t, err)
			extension, err := extensions.FindExtension(extensionId)
			assert.NoError(t, err)
			assert.NotNil(t, extension)

			// add vim details
			spec.AddVim(tt.vimName)
			_, err = spec.GetVim(tt.vimName)
			assert.NoError(t, err)

			_, err = api.UpdateExtension(ctx, spec)
			if !tt.errOnUpdate {
				assert.NoError(t, err)
			}

			if tt.delAfter {
				_, err := api.DeleteExtension(ctx, spec.Name)
				assert.NoError(t, err)
			}
		})
	}
}

// TestTcaApiExtensionQuery
func TestExtensionQuery(t *testing.T) {

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

			ctx := context.Background()
			a := getTcaApi(t, rest, false)
			got, err := a.ExtensionQuery(ctx)

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

func TestExtensionQueryFindRepo(t *testing.T) {

	tests := []struct {
		name     string
		rest     *client.RestClient
		wantErr  bool
		repoName string
	}{
		{
			name:     "Basic Find must return ok",
			rest:     rest,
			wantErr:  false,
			repoName: getTestRepoName(),
		},
		{
			name:     "Basic Find must return error",
			rest:     rest,
			wantErr:  true,
			repoName: "notfound",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			a := getTcaApi(t, rest, false)

			got, err := a.ExtensionQuery(ctx)
			assert.NoError(t, err)

			valid, err := got.FindRepo(tt.repoName)
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
  "name": "test_repo",
  "version": "v2.x",
  "type": "Repository",
  "extensionKey": "",
  "extensionSubtype": "Harbor",
  "products": [],
  "vimInfo": [],
  "interfaceInfo": {
    "url": "https://1.1.1.1",
    "Description": "",
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
  "name": "unit_test",
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
    "Description": "",
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
    "Description": "",
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
  "name": "Repo",
  "version": "v2.x",
  "type": "Repository",
  "extensionKey": "",
  "extensionSubtype": "Harbor",
  "products": [],
  "interfaceInfo": {
    "url": "https://10.252.80.135",
    "Description": "",
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
    "Description": "",
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
