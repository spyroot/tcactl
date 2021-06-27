package api

import (
	"encoding/json"
	"github.com/nsf/jsondiff"
	_ "github.com/nsf/jsondiff"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/io"
	"testing"
)

// Reads spec and validate parser
func TestReadNodeSpecFromString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{
			name:    "Basic JSON Template",
			args:    jsonNodeSpec,
			wantErr: false,
		},
		{
			name:    "Basic Yaml Template",
			args:    yamlNodeSpec,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadNodeSpecFromString(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr == false {
				t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			newJson, err := json.Marshal(got)
			if err != nil {
				return
			}
			oldJson, err := json.Marshal(got)
			if err != nil {
				return
			}

			//fmt.Println(string(b))

			opt := jsondiff.DefaultJSONOptions()
			diff, _ := jsondiff.Compare(newJson, oldJson, &opt)

			if tt.wantErr != true && diff > 0 {
				t.Errorf("ReadNodeSpecFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("diff %d", diff)
		})
	}
}

// specNodePoolStringReaderHelper helper return node spec
func specNodePoolStringReaderHelper(s string) *request.NewNodePoolSpec {
	r, err := ReadNodeSpecFromString(s)
	io.CheckErr(err)
	return r
}

func TestTcaCreateNewNodePool(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		spec    *request.NewNodePoolSpec
		wantErr bool
		reset   bool
	}{
		{
			name:    "Create From File mgmt template",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(jsonNodeSpec),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				err  error
				task *models.TcaTask
			)

			a, err := NewTcaApi(tt.rest)

			if tt.reset {
				a.rest = nil
			}

			tt.spec.CloneMode = request.LinkedClone
			// note cluster statically defined
			// TODO it take time to create cluster so only way for now create test cluster and use for unit testing

			if tt.spec != nil {
				tt.spec.Name = generateName()
				tt.spec.Labels[0] = "type=" + tt.spec.Name
			}

			if task, err = a.CreateNewNodePool(tt.spec, getTestClusterName(), false); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewNodePool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Errorf("CreateClusterTemplate() error is nil, wantErr %v", tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				t.Logf("Recieved correct error %v", err)
				return
			}

			if tt.wantErr == false && task == nil {
				t.Logf("Recieved correct error %v", err)
				return
			}

			if tt.wantErr == false && len(task.OperationId) == 0 {
				t.Logf("Recieved correct error %v", err)
				return
			}
		})
	}
}
