package vc

import (
	"context"
	"os"
	"reflect"
	"testing"
)

func TestUpload(t *testing.T) {

	if VcEnvChecker(t, false) == false {
		t.Skipf("Skipping VC environment unset.")
	}

	type uploadArgs struct {
		datastore string
		src       string
		dst       string
	}

	tests := []struct {
		name       string
		vcArgs     vcArgs
		actionArgs uploadArgs
		wantIsVc   bool
		wantErr    bool
	}{
		{
			name:     "nil args",
			wantIsVc: false,
			wantErr:  true,
		},
		{
			name:     "empty vc details",
			vcArgs:   generateEmptyCredential(),
			wantIsVc: false,
			wantErr:  true,
		},
		{
			name: "basic upload",
			vcArgs: vcArgs{
				os.Getenv("VC_HOST"),
				os.Getenv("VC_USERNAME"),
				os.Getenv("VC_PASSWORD"),
			},
			actionArgs: uploadArgs{
				datastore: "vsanDatastore",
				src:       "core-13.0.iso",
				dst:       "/ISO/core-13.0.iso",
			},
			wantIsVc: true,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Connect(context.TODO(), tt.vcArgs.hostname, tt.vcArgs.username, tt.vcArgs.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				return
			}

			if !reflect.DeepEqual(c.IsVC(), tt.wantIsVc) {
				t.Errorf("Connect() is vc got = %v, want %v", c.IsVC(), tt.wantIsVc)
			}

			vcRest := VSphereRest{Ctl: c.Client}
			err = vcRest.Upload(context.TODO(), tt.actionArgs.datastore, tt.actionArgs.src, tt.actionArgs.dst)
		})
	}
}
