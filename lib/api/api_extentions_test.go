package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTcaApi_ExtensionQuery(t *testing.T) {

	tests := []struct {
		name string
		rest *client.RestClient
		//want    *response.Extensions
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
