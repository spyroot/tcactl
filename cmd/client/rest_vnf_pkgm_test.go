package client

import (
	"github.com/spyroot/hestia/cmd/csar"
	testUtil "github.com/spyroot/hestia/pkg/testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestRestClient_CreateVnfPkgmVnfd(t *testing.T) {

	tests := []struct {
		name    string
		client  *RestClient
		wantErr bool
		arg     *PackageUpload
	}{
		{
			name:    "create empty and delete",
			client:  rest,
			wantErr: false,
			arg:     NewPackageUpload("test_upload24"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := tt.client.GetAuthorization()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfrNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			respond, err := tt.client.CreateVnfPkgmVnfd(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfrNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotEqualValuesf(t, respond, nil, "respond must not be nil")
			assert.Equal(t, respond.OnboardingState, "CREATED", "Status must created.")
			assert.NotEqualValuesf(t, len(respond.Id), 0, "Server must respond with id")

			ok, err := tt.client.DeleteVnfPkgmVnfd(respond.Id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfrNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, ok, true, "Delete must return true.")

		})
	}
}

func TestRestClient_UploadVnfPkgmVnfd(t *testing.T) {

	tests := []struct {
		name     string
		client   *RestClient
		wantErr  bool
		arg      *PackageUpload
		fileName string
	}{
		{
			name:     "create, upload and delete",
			client:   rest,
			wantErr:  false,
			arg:      NewPackageUpload(testUtil.RandomString(8)),
			fileName: "/Users/spyroot/go/src/hestia/tests/smokeping-cnf.csar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := tt.client.GetAuthorization()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfrNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var substitution = map[string]string{}
			substitution["descriptorId"] = "nfd_1234"
			newCsarFile, err := csar.ApplyTransformation(tt.fileName, csar.SpecNfd, csar.NfdYamlPropertyTransformer, substitution)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestRestClient_UploadVnfPkgmVnfd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			file, _ := os.Open(newCsarFile)
			fileBytes, _ := ioutil.ReadAll(file)
			fileName := filepath.Base(newCsarFile)

			respond, err := tt.client.CreateVnfPkgmVnfd(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestRestClient_UploadVnfPkgmVnfd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotEqualValuesf(t, respond, nil, "respond must not be nil")
			assert.Equal(t, respond.OnboardingState, "CREATED", "Status must created.")
			assert.NotEqualValuesf(t, len(respond.Id), 0, "Server must respond with id")

			ok, err := tt.client.UploadVnfPkgmVnfd(respond.Id, fileBytes, fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestRestClient_UploadVnfPkgmVnfd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, ok, true, "return must return true.")

			//ok, err := tt.client.DeleteVnfPkgmVnfd(respond.Id)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("GetInfrNetworks() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			//
			//assert.Equal(t, ok, true, "Delete must return true.")

		})
	}
}
