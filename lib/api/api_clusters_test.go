// Package api
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Mustafa mbayramo@vmware.com
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

			clusterTask, err := a.GetClusterTask(getTestClusterName(), true)
			if err != nil {
				return
			}

			io.PrettyPrint(clusterTask)
		})
	}
}
