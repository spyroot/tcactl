package vc

import (
	"context"
	"os"
	"reflect"
	"testing"
)

type vcArgs struct {
	hostname string
	username string
	password string
}

func generateEmptyCredential() vcArgs {
	return vcArgs{"", "", ""}
}

func VcEnvChecker(t *testing.T, strict bool) bool {
	if len(os.Getenv("VC_HOST")) == 0 {
		if strict {
			t.Fatalf("Test requires VC_HOST environment.")
		}
		return false
	}
	if len(os.Getenv("VC_USERNAME")) == 0 {
		if strict {
			t.Fatalf("Test requires VC_USERNAME environment.")
		}
		return false
	}
	if len(os.Getenv("VC_PASSWORD")) == 0 {
		if strict {
			t.Fatalf("Test requires VC_PASSWORD environment.")
		}
		return false
	}
	return true
}

func TestConnect(t *testing.T) {

	if VcEnvChecker(t, false) == false {
		t.Skipf("Skipping VC environment unset.")
	}

	tests := []struct {
		name     string
		args     vcArgs
		wantIsVc bool
		wantErr  bool
	}{
		{
			name:     "nil args",
			wantIsVc: false,
			wantErr:  true,
		},
		{
			name:     "nil args value",
			args:     generateEmptyCredential(),
			wantIsVc: false,
			wantErr:  true,
		},
		{
			name: "basic connect",
			args: vcArgs{
				os.Getenv("VC_HOST"),
				os.Getenv("VC_USERNAME"),
				os.Getenv("VC_PASSWORD"),
			},
			wantIsVc: true,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Connect(context.TODO(), tt.args.hostname, tt.args.username, tt.args.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				return
			}

			if !reflect.DeepEqual(got.IsVC(), tt.wantIsVc) {
				t.Errorf("Connect() is vc got = %v, want %v", got, tt.wantIsVc)
			}
		})
	}
}
