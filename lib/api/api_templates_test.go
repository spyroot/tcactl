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
	"github.com/google/uuid"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/pkg/io"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// return initialSpec from file
func specFileReaderHelper(f string) *response.ClusterTemplate {
	s, err := ReadTemplateSpecFromFile(f)
	io.CheckErr(err)
	return s
}

// return initialSpec from string
func specStringReaderHelper(s string) *response.ClusterTemplate {
	r, err := ReadTemplateSpecFromString(s)
	io.CheckErr(err)
	return r
}

// Create new cluster template
//
func TestTcaApi_CreateClusterTemplate(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		spec    *response.ClusterTemplate
		wantErr bool
		reset   bool
	}{
		{
			name:    "Create From File mgmt template",
			rest:    rest,
			spec:    specStringReaderHelper(yamlMgmtTemplate),
			wantErr: false,
		},
		{
			name:    "Invalid template no master",
			rest:    rest,
			spec:    specStringReaderHelper(yamlInvalidMgmtTemplate),
			wantErr: true,
		},
		{
			name:    "Invalid template no worker",
			rest:    rest,
			spec:    specStringReaderHelper(yamlInvalidMgmtTemplate2),
			wantErr: true,
		},
		{
			name:    "Invalid template no template type",
			rest:    rest,
			spec:    specStringReaderHelper(yamlInvalidMgmtTemplate3),
			wantErr: true,
		},
		{
			name:    "Invalid template no template no network label",
			rest:    rest,
			spec:    specStringReaderHelper(yamlInvalidMgmtTemplate4),
			wantErr: true,
		},
		{
			name:    "Empty",
			rest:    rest,
			spec:    specStringReaderHelper(yamlWorkloadEmpty),
			wantErr: true,
		},
		{
			name:    "Nil rest",
			rest:    rest,
			spec:    specStringReaderHelper(yamlWorkloadEmpty),
			wantErr: true,
			reset:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				err      error
				err2     error
				name     string
				template *response.ClusterTemplate
			)

			a, err := NewTcaApi(tt.rest)

			if tt.reset {
				a.rest = nil
			}

			if name, err = a.CreateClusterTemplate(tt.spec); (err != nil) != tt.wantErr {
				t.Errorf("CreateClusterTemplate() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Errorf("CreateClusterTemplate() error is nil, vimErr %v", tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				t.Logf("Recieved correct error %v", err)
				return
			}

			time.Sleep(3 * time.Second)

			if template, err2 = a.GetClusterTemplate(name); (err2 != nil) != tt.wantErr {
				t.Errorf("CreateClusterTemplate() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Template created %s", name)

			if template == nil {
				if !tt.wantErr {
					t.Errorf("CreateClusterTemplate() error = %v, vimErr %v", err, tt.wantErr)
					return
				}
			} else {
				t.Logf("Template retrieved %s", template.Id)
				if err = a.DeleteTemplate(name); (err != nil) != tt.wantErr {
					t.Errorf("DeleteTemplate() error = %v, vimErr %v", err, tt.wantErr)
					return
				}
				t.Logf("Template %s deleted", template.Id)
			}
		})
	}
}

// TestTcaApi_GetClusterTemplate
func TestTcaApi_GetClusterTemplate(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		specs   *response.ClusterTemplate
		repeat  int
		wantErr bool
	}{
		{
			name:    "Create cluster and check list",
			rest:    rest,
			specs:   specStringReaderHelper(yamlWorkloadTemplate),
			repeat:  3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, _ := NewTcaApi(tt.rest)

			var names []string
			for i := 0; i < tt.repeat; i++ {
				if tt.specs != nil {
					u := uuid.New().String()
					tt.specs.Name = u[0:8]
					_, err := a.CreateClusterTemplate(tt.specs)
					assert.NoError(t, err)
					time.Sleep(1 * time.Second)
					names = append(names, u[0:8])
				}
			}

			for _, n := range names {
				_, err := a.GetClusterTemplate(n)
				if tt.wantErr == false && err != nil {
					t.Errorf("GetClusterTemplates() error = %v, vimErr %v", err, tt.wantErr)
					return
				}
			}

			// delete all
			for _, n := range names {
				err := a.DeleteTemplate(n)
				if tt.wantErr == false && err != nil {
					t.Errorf("GetClusterTemplates() error = %v, vimErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestTcaApi_GetClusterTemplates(t *testing.T) {
	tests := []struct {
		name    string
		rest    *client.RestClient
		specs   *response.ClusterTemplate
		repeat  int
		wantErr bool
	}{
		{
			name:    "Create cluster and check list",
			rest:    rest,
			specs:   specStringReaderHelper(yamlWorkloadTemplate),
			repeat:  3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, _ := NewTcaApi(tt.rest)

			var names []string
			for i := 0; i < tt.repeat; i++ {
				if tt.specs != nil {
					u := uuid.New().String()
					tt.specs.Name = u[0:8]
					_, err := a.CreateClusterTemplate(tt.specs)
					assert.NoError(t, err)
					time.Sleep(1 * time.Second)
					names = append(names, u[0:8])
				}
			}

			got, err := a.GetClusterTemplates()
			assert.NotNil(t, got)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterTemplates() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			for _, n := range names {
				_, err := got.GetTemplateId(n)
				if tt.wantErr == false && err != nil {
					t.Errorf("GetClusterTemplates() error = %v, vimErr %v", err, tt.wantErr)
				}
			}

			// delete all
			for _, n := range names {
				err = a.DeleteTemplate(n)
				if tt.wantErr == false && err != nil {
					t.Errorf("GetClusterTemplates() error = %v, vimErr %v", err, tt.wantErr)
				}
			}

		})
	}
}

func changeDescription(t *response.ClusterTemplate) *response.ClusterTemplate {
	t.Description = "test"
	return t
}

func checkerDescription(t *response.ClusterTemplate) bool {
	return t.Description == "test"
}

func changeReplica(t *response.ClusterTemplate) *response.ClusterTemplate {
	t.WorkerNodes[0].Replica = 5
	return t
}

func checkerReplica(t *response.ClusterTemplate) bool {
	return t.WorkerNodes[0].Replica == 5
}

func changeTemplateId(t *response.ClusterTemplate) *response.ClusterTemplate {
	t.Id = t.Name
	return t
}

func checkerTemplateId(t *response.ClusterTemplate) bool {
	return IsValidUUID(t.Id)
}

func changeInvalidTemplateId(t *response.ClusterTemplate) *response.ClusterTemplate {
	t.Id = "01ac6fc7-435e-4eb6-9444-c6ed07999999"
	return t
}

func TestTcaApi_UpdateClusterTemplate(t *testing.T) {

	tests := []struct {
		name        string
		rest        *client.RestClient
		initialSpec *response.ClusterTemplate
		wantErr     bool
		transformer func(*response.ClusterTemplate) *response.ClusterTemplate
		checker     func(*response.ClusterTemplate) bool
		recheck     bool
	}{
		{
			name:        "Update description",
			rest:        rest,
			initialSpec: specStringReaderHelper(yamlMgmtTemplate),
			wantErr:     false,
			transformer: changeDescription,
			checker:     checkerDescription,
			recheck:     true,
		},
		{
			name:        "Update replica",
			rest:        rest,
			initialSpec: specStringReaderHelper(yamlMgmtTemplate),
			wantErr:     false,
			transformer: changeReplica,
			checker:     checkerReplica,
			recheck:     true,
		},
		{
			name:        "Use name instead if in update",
			rest:        rest,
			initialSpec: specStringReaderHelper(yamlMgmtTemplate),
			wantErr:     false,
			transformer: changeTemplateId,
			checker:     checkerTemplateId,
			recheck:     false,
		},
		{
			name:        "Use invalid id",
			rest:        rest,
			initialSpec: specStringReaderHelper(yamlMgmtTemplate),
			wantErr:     true,
			transformer: changeInvalidTemplateId,
			checker:     nil,
			recheck:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, _ := NewTcaApi(tt.rest)

			var (
				err      error
				name     string
				template *response.ClusterTemplate
			)

			if tt.initialSpec != nil {
				name, err = a.CreateClusterTemplate(tt.initialSpec)
				assert.NoError(t, err)
				time.Sleep(3 * time.Second)
				template, err = a.GetClusterTemplate(name)
				assert.NoError(t, err)
			}

			t.Logf("Template created %s", name)

			assert.NotNil(t, template)

			transformed := tt.transformer(template)

			if tt.recheck {
				assert.NotEqual(t, tt.checker(transformed), false)
			}

			var errUpdate error
			if errUpdate = a.UpdateClusterTemplate(transformed); (errUpdate != nil) != tt.wantErr {
				t.Errorf("UpdateClusterTemplate() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			t.Log(errUpdate)

			if errUpdate != nil && tt.wantErr {
				t.Logf("Wanted error and got error ted %s", errUpdate)
				return
			}

			t.Logf("Template updated %s", template.Id)

			var updateTemplate *response.ClusterTemplate
			// get result and compare
			if updateTemplate, errUpdate = a.GetClusterTemplate(name); (errUpdate != nil) != tt.wantErr {
				t.Errorf("GetClusterTemplate() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			io.PrettyPrint(updateTemplate)

			if tt.checker(updateTemplate) == true && tt.wantErr == false {
				if err = a.DeleteTemplate(name); (err != nil) != tt.wantErr {
					t.Errorf("DeleteTemplate() error = %v, vimErr %v", err, tt.wantErr)
					return
				}
				t.Logf("Template %s deleted", template.Id)
			}

		})
	}
}

func TestTcaApi_DeleteClusterTemplate(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		spec    *response.ClusterTemplate
		tid     string
		wantErr bool
	}{
		{
			name:    "Should produce error.",
			rest:    rest,
			spec:    nil,
			tid:     "",
			wantErr: true,
		},
		{
			name:    "Should produce error.",
			rest:    rest,
			spec:    nil,
			tid:     "abc",
			wantErr: true,
		},
		{
			name:    "Should produce error nil rest.",
			rest:    nil,
			spec:    nil,
			tid:     "abc",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := &TcaApi{
				rest: tt.rest,
			}

			var (
				err  error
				err2 error
				name string
			)

			// if need create initialSpec
			if tt.spec != nil {
				if name, err = a.CreateClusterTemplate(tt.spec); (err != nil) != tt.wantErr {
					t.Errorf("CreateClusterTemplate() error = %v, vimErr %v", err, tt.wantErr)
					return
				}

				if tt.wantErr && err == nil {
					t.Errorf("CreateClusterTemplate() error is nil, vimErr %v", tt.wantErr)
					return
				}

				if tt.wantErr && err != nil {
					t.Logf("Recieved correct error %v", err)
					return
				}

				time.Sleep(3 * time.Second)

				if _, err2 = a.GetClusterTemplate(name); (err2 != nil) != tt.wantErr {
					t.Errorf("CreateClusterTemplate() error = %v, vimErr %v", err, tt.wantErr)
					return
				}

				t.Logf("Template created %s", name)
			}

			if err = a.DeleteTemplate(name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTemplate() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Errorf("DeleteTemplate() error is nil, vimErr %v", tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				t.Logf("Recieved correct error %v", err)
				return
			}
		})
	}
}
