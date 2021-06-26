package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/io"
	"testing"
	"time"
)

func TestTcaApi_GetClusterTask(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		spec    *request.NewNodePoolSpec
		wantErr bool
		reset   bool
	}{
		{
			name:    "Create cluster and check task list",
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

			if tt.spec != nil {
				tt.spec.Name = generateName()
				tt.spec.Labels[0] = "type=" + tt.spec.Name
			}

			tt.spec.CloneMode = request.LinkedClone
			// note cluster statically defined
			// TODO it take time to create cluster so only way for now create test cluster and use for unit testing
			if task, err = a.CreateNewNodePool(tt.spec, getTestClusterName(), false); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewNodePool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			time.Sleep(3 * time.Second)

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

			clusterTask, err := a.GetClusterTask(getTestClusterName())
			if err != nil {
				return
			}

			io.PrettyPrint(clusterTask)
		})
	}
}
