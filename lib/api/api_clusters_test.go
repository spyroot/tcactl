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
	"context"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/lib/testutil"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
	"time"
)

func getTcaApi(t *testing.T, rest *client.RestClient, isLogEnabled bool) *TcaApi {
	a, err := NewTcaApi(rest)
	assert.NoError(t, err)
	if isLogEnabled {
		a.SetTrace(true)
		SetLoggingFlags()
	}
	return a
}

func getTestWorkloadCluster(t *testing.T, api *TcaApi) *response.ClusterSpec {
	assert.NotNil(t, api)
	c, err := api.GetCluster(context.Background(), getTestWorkloadClusterName())
	assert.NoError(t, err)
	return c
}

func getTestMgmtCluster(t *testing.T, api *TcaApi) *response.ClusterSpec {
	assert.NotNil(t, api)
	c, err := api.GetCluster(context.Background(), getTestMgmtClusterName())
	assert.NoError(t, err)
	return c
}

func NormalizeClusterIP(s string) string {
	u := s

	if strings.HasPrefix(u, "https://") {
		u = strings.TrimPrefix(u, "https://")
	}

	if strings.HasSuffix(u, ":6443") {
		u = strings.TrimSuffix(u, ":6443")
	}

	return u
}

func TestGetClusterIPs(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		reader       io.Reader
		wantErr      bool
		isLogEnabled bool
	}{
		{
			name:    "Test cluster IP must be resolved",
			rest:    rest,
			wantErr: false,
			reader:  testutil.SpecTempReader(newManagementCluster),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			a := getTcaApi(t, rest, tt.isLogEnabled)
			wCluster := getTestWorkloadCluster(t, a)
			mCluster := getTestMgmtCluster(t, a)

			clusters, err := a.GetClusters(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCluster() error = %v, vimErr %v", err, tt.wantErr)
				return
			}
			if clusters == nil {
				t.Errorf("GetClusters() must not return nil")
				return

			}

			if !tt.wantErr {
				m := clusters.GetClusterIPs()
				_, okWorkloadCluster := m[NormalizeClusterIP(wCluster.ClusterUrl)]
				assert.Equal(t, true, okWorkloadCluster)
				_, okMgmtCluster := m[NormalizeClusterIP(mCluster.ClusterUrl)]
				assert.Equal(t, true, okMgmtCluster)

			}
		})
	}
}

func TestAllocateNewClusterIp(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		reader       io.Reader
		wantErr      bool
		isLogEnabled bool
		isConflict   bool
		ipAddr       string
	}{
		{
			name:       "Check IP overlap mgmt cluster must generate new ip ",
			rest:       rest,
			wantErr:    false,
			reader:     testutil.SpecTempReader(newManagementCluster),
			isConflict: true,
		},
		{
			name:       "Check IP overlap workload cluster must generate new ip ",
			rest:       rest,
			wantErr:    false,
			reader:     testutil.SpecTempReader(newTestWorkloadCluster),
			isConflict: true,
		},
		{
			name:       "Check IP overlap workload cluster must not generate new ip ",
			rest:       rest,
			wantErr:    false,
			reader:     testutil.SpecTempReader(newTestWorkloadCluster),
			isConflict: false,
			ipAddr:     "222.22.2.2",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			a := getTcaApi(t, rest, tt.isLogEnabled)
			wCluster := getTestWorkloadCluster(t, a)
			_spec, err := specs.SpecCluster{}.SpecsFromReader(tt.reader)
			assert.NoError(t, err)
			assert.NotNil(t, _spec)
			spec, ok := (*_spec).(*specs.SpecCluster)
			assert.Equal(t, true, ok)

			// create conflict
			if tt.isConflict {
				conflictedAddr := NormalizeClusterIP(wCluster.ClusterUrl)
				spec.EndpointIP = conflictedAddr
				err := a.allocateNewClusterIp(context.Background(), spec)
				if !tt.wantErr && err != nil {
					t.Errorf("TestAllocateNewClusterIp() Test must not fail")
					return
				}

				// make sure it resolved
				if spec.EndpointIP == conflictedAddr {
					t.Errorf("TestAllocateNewClusterIp() test faield confict not resolved")
					return
				}
			} else {
				spec.EndpointIP = tt.ipAddr
				err := a.allocateNewClusterIp(context.Background(), spec)
				if !tt.wantErr && err != nil {
					t.Errorf("TestAllocateNewClusterIp() Test must not fail")
					return
				}

				// make sure it resolved
				if spec.EndpointIP != tt.ipAddr {
					t.Errorf("TestAllocateNewClusterIp() test faield confict not resolved")
					return
				}
			}
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
// Test uses existing cluster ID and name
// make sure they adjusted or passed via env variable
func TestGetCluster(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		wantErr      bool
		cluster      string
		expect       string
		isLogEnabled bool
	}{
		{
			name:    "Get cluster by name must not fail",
			rest:    rest,
			cluster: getTestWorkloadClusterName(),
			expect:  getTestClusterId(),
			wantErr: false,
		},
		{
			name:    "Get cluster by cluster id must not fail",
			rest:    rest,
			cluster: getTestClusterId(),
			expect:  getTestClusterId(),
			wantErr: false,
		},
		{
			name:    "Get cluster wrong id must fail.",
			rest:    rest,
			cluster: "868636c9-868f-49fb-a6df-6a0d2d137141",
			expect:  getTestClusterId(),
			wantErr: true,
		},
		{
			name:    "Get cluster wrong cluster name must fail.",
			rest:    rest,
			cluster: "test",
			expect:  getTestClusterId(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := getTcaApi(t, rest, tt.isLogEnabled)

			actual, err := a.GetCluster(context.Background(), tt.cluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCluster() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, actual)
				if actual != nil {
					assert.Equal(t, tt.expect, actual.Id)
				} else {
					t.Errorf("TestGetCluster() cluster must not be nil")
					return
				}
			}
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
func TestGetClusterNodePool(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		wantErr      bool
		cluster      string
		nodepool     string
		expect       string
		trace        bool
		isLogEnabled bool
	}{
		{
			name:     "Must resolve from cluster name and pool name.",
			rest:     rest,
			cluster:  getTestWorkloadClusterName(),
			nodepool: getTestNodePoolName(),
			expect:   getTestNodePoolId(),
			wantErr:  false,
			trace:    false,
		},
		{
			name:     "Must resolve from cluster id and pool name.",
			rest:     rest,
			cluster:  getTestClusterId(),
			nodepool: getTestNodePoolName(),
			expect:   getTestNodePoolId(),
			wantErr:  false,
			trace:    false,
		},
		{
			name:     "Must resolve from cluster id and pool id.",
			rest:     rest,
			cluster:  getTestClusterId(),
			nodepool: getTestNodePoolId(), // this must be adjust based on deployment
			expect:   getTestNodePoolId(), // this must be adjust based on deployment
			wantErr:  false,
			trace:    false,
		},
		{
			name:     "Wrong cluster id valid pool id must fail.",
			rest:     rest,
			cluster:  "test123213",
			nodepool: getTestNodePoolName(),
			expect:   getTestNodePoolId(),
			wantErr:  true,
			trace:    false,
		},
		{
			name:     "Correct cluster id wrong pool name must fail.",
			rest:     rest,
			cluster:  getTestClusterId(),
			nodepool: "asdasd",
			expect:   getTestNodePoolId(),
			wantErr:  true,
			trace:    false,
		},
		{
			name:     "Correct cluster id wrong pool id must fail.",
			rest:     rest,
			cluster:  getTestClusterId(),
			nodepool: "3acf9b79-f8e5-4155-997b-58792d395555",
			expect:   getTestNodePoolId(),
			wantErr:  true,
			trace:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			if tt.isLogEnabled {
				a.SetTrace(tt.trace)
				SetLoggingFlags()
			}

			actual, err := a.GetClusterNodePool(context.Background(), tt.cluster, tt.nodepool)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterNodePool() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, actual)
				if actual != nil {
					t.Log(actual.Id)
				} else {
					t.Errorf("GetClusterNodePool() error = %v, vimErr %v", err, tt.wantErr)
					return
				}
			}

			if tt.wantErr {
				t.Log(err)
			}
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
func TestGetCurrentClusterTask(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		reader       io.Reader
		wantErr      bool
		reset        bool
		trace        bool
		sec          time.Duration
		isLogEnabled bool
	}{
		{
			name:    "Create mgmt cluster must not fail.",
			rest:    rest,
			reader:  testutil.SpecTempReader(newManagementCluster),
			wantErr: false,
			trace:   false,
			sec:     1,
		},
		{
			name:    "Create workload cluster must not fail.",
			rest:    rest,
			reader:  testutil.SpecTempReader(newTestWorkloadCluster),
			wantErr: false,
			trace:   false,
			sec:     1,
		},
		{
			name:    "Create workload cluster must fail.",
			rest:    rest,
			reader:  testutil.SpecTempReader(InvalidWorkerNodeFolder),
			wantErr: true,
			trace:   false,
			sec:     1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				err  error
				task *models.TcaTask
			)

			a := getTcaApi(t, rest, tt.isLogEnabled)
			_spec, err := specs.SpecCluster{}.SpecsFromReader(tt.reader)
			assert.NoError(t, err)
			assert.NotNil(t, _spec)
			spec, ok := (*_spec).(*specs.SpecCluster)
			assert.Equal(t, true, ok)
			assert.NotNil(t, spec)

			if spec.IsManagement() {
				spec.ClusterTemplateId = getTestMgmtTemplateId()
			} else {
				spec.ClusterTemplateId = getTestWorkloadTemplateId()
			}

			assert.NoError(t, err)
			if task, err = a.CreateClusters(context.Background(), &ClusterCreateApiReq{
				Spec:          spec,
				IsBlocking:    false,
				IsDryRun:      false,
				IsVerbose:     false,
				IsFixConflict: true}); (err != nil) != tt.wantErr {
				t.Errorf("CreateClusters() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			// sleep
			time.Sleep(tt.sec * time.Second)

			if tt.wantErr && err == nil {
				t.Errorf("CreateClusters() error is nil, vimErr %v", tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				return
			}

			if task == nil {
				t.Errorf("CreateClusters() error task must not be nil")
				return
			}

			assert.NotNil(t, task)
			t.Logf("SpecCluster create task created %s", task.Id)
			clusterCreateTask, err := a.GetCurrentClusterTask(context.Background(), task.Id)
			assert.NoError(t, err)
			if _, err := clusterCreateTask.FindEntityByName(spec.Name); (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentClusterTask() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Task created %v", clusterCreateTask.Items[0].TaskId)
			cluster, err := a.GetCluster(context.Background(), spec.Name)
			assert.NoError(t, err)
			assert.NotEmpty(t, cluster.Id)

			t.Logf("New cluster id %s", cluster.Id)
			delTask, err := a.DeleteCluster(context.Background(),
				&ClusterDeleteApiReq{
					Cluster:    cluster.Id,
					IsBlocking: false,
					IsVerbose:  false,
				})

			if err != nil {
				t.Log("Failed delete cluster.", err)
				return
			}

			t.Logf("Delete cluster task create %v", delTask.Id)
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
func TestCreateClusters(t *testing.T) {

	tests := []struct {
		name            string
		rest            *client.RestClient
		reader          io.Reader
		wantErr         bool
		dryRun          bool
		reset           bool
		trace           bool
		useTestTemplate bool
		isLogEnabled    bool
		wantDelete      bool
		sec             time.Duration
	}{
		{
			name:            "Create mgmt dryRun run.",
			rest:            rest,
			reader:          testutil.SpecTempReader(newManagementCluster),
			wantErr:         false,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create mgmt dryRun run - wrong template01 no template id.",
			rest:            rest,
			reader:          testutil.SpecTempReader(NewManagementClusterFailCase01),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: false,
			sec:             3,
		},
		{
			name:            "Create mgmt dryRun run - wrong template02 no cloud id.",
			rest:            rest,
			reader:          testutil.SpecTempReader(NewManagementClusterFailCase02),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong mgmt cluster should fail",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidMgmtCluster),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong template id should fail",
			rest:            rest,
			reader:          testutil.SpecTempReader(WrongMgmtTemplateId),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: false,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong datastore",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidDatastoreGlobal),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong datastore",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidDatastoreWorker),
			wantErr:         true,
			dryRun:          true,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong datastore",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidDatastoreMasterNode),
			wantErr:         true,
			dryRun:          true,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidFolderGlobal),
			wantErr:         true,
			dryRun:          true,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong master node folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidMasterNodeFolder),
			wantErr:         true,
			dryRun:          true,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong master node folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidWorkerNodeFolder),
			wantErr:         true,
			dryRun:          true,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong master node folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidDatastoreUrl),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong master node folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidDatastoreUrl),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create mgmt cluster must not fail",
			rest:            rest,
			reader:          testutil.SpecTempReader(newManagementCluster),
			useTestTemplate: true,
			wantDelete:      true,
			sec:             3,
		},
		{
			name:            "Create workload cluster must not fail",
			rest:            rest,
			reader:          testutil.SpecTempReader(newTestWorkloadCluster),
			useTestTemplate: true,
			wantDelete:      true,
			sec:             3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				err  error
				task *models.TcaTask
			)

			a := getTcaApi(t, rest, tt.isLogEnabled)
			_spec, err := specs.SpecCluster{}.SpecsFromReader(tt.reader)
			assert.NoError(t, err)
			assert.NotNil(t, _spec)
			spec, ok := (*_spec).(*specs.SpecCluster)
			assert.Equal(t, true, ok)
			assert.NotNil(t, spec)

			if tt.useTestTemplate {
				if spec.IsManagement() {
					spec.ClusterTemplateId = getTestMgmtTemplateId()
				} else {
					spec.ClusterTemplateId = getTestWorkloadTemplateId()
				}
			}

			if task, err = a.CreateClusters(context.Background(),
				&ClusterCreateApiReq{
					Spec:          spec,
					IsBlocking:    false,
					IsDryRun:      tt.dryRun,
					IsVerbose:     false,
					IsFixConflict: true}); (err != nil) != tt.wantErr {
				t.Errorf("CreateClusters() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Errorf("CreateClusters() didn't fail")
				return
			}

			if tt.wantErr && err != nil {
				return
			}

			if task == nil {
				t.Errorf("CreateClusters() create task must return task id")
				return
			}

			if tt.dryRun {
				if len(task.Id) > 0 {
					t.Errorf("CreateClusters() in dry run task must be empty")
					return
				}
				return
			}

			if !tt.dryRun {
				t.Logf("Task created %s", task.Id)
				return
			}

			if tt.wantDelete {
				cluster, err := a.GetCluster(context.Background(), spec.Name)
				assert.NoError(t, err)
				assert.NotEmpty(t, cluster.Id)

				_, _ = a.DeleteCluster(context.Background(),
					&ClusterDeleteApiReq{
						Cluster:    cluster.Id,
						IsBlocking: false,
						IsVerbose:  false,
					})
			}
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
func TestCreateBlockDeleteCluster(t *testing.T) {

	tests := []struct {
		name            string
		rest            *client.RestClient
		reader          io.Reader
		wantErr         bool
		dryRun          bool
		reset           bool
		trace           bool
		useTestTemplate bool
		sec             time.Duration
		isLogEnabled    bool
	}{
		{
			name:            "Create block and Delete a new tenant cluster.",
			rest:            rest,
			reader:          testutil.SpecTempReader(newManagementCluster),
			wantErr:         false,
			dryRun:          false,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				err        error
				createTask *models.TcaTask
				deleteTask *models.TcaTask
			)

			a := getTcaApi(t, rest, tt.isLogEnabled)
			_spec, err := specs.SpecCluster{}.SpecsFromReader(tt.reader)
			assert.NoError(t, err)
			assert.NotNil(t, _spec)
			spec, ok := (*_spec).(*specs.SpecCluster)
			assert.Equal(t, true, ok)
			assert.NotNil(t, spec)

			// adjust template id
			if tt.useTestTemplate {
				if spec.IsManagement() {
					spec.ClusterTemplateId = getTestMgmtTemplateId()
				} else {
					spec.ClusterTemplateId = getTestWorkloadTemplateId()
				}
			}

			assert.NoError(t, err)
			createTask, err = a.CreateClusters(context.Background(),
				&ClusterCreateApiReq{
					Spec:          spec,
					IsBlocking:    true,
					IsDryRun:      tt.dryRun,
					IsVerbose:     false,
					IsFixConflict: true})

			// if task failed we still should be able delete cluster.
			if _, ok := err.(*TcaTaskFailed); ok {
				t.Log("Create cluster task failed but test will try to delete it")
				if deleteTask, err = a.DeleteCluster(context.Background(),
					&ClusterDeleteApiReq{
						Cluster:    spec.Name,
						IsBlocking: false,
						IsVerbose:  false}); (err != nil) != tt.wantErr {
					t.Errorf("CreateBlockDeleteCluster() error = %v, vimErr %v", err, tt.wantErr)
					return
				}
				return
			} else {
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateBlockDeleteCluster() error = %v, vimErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr {
					assert.NotNil(t, createTask)
					// wait 2 sec and send delete
					time.Sleep(2 * time.Second)

					if deleteTask, err = a.DeleteCluster(context.Background(),
						&ClusterDeleteApiReq{
							Cluster:    spec.Name,
							IsBlocking: false,
							IsVerbose:  false,
						}); (err != nil) != tt.wantErr {
						t.Errorf("CreateBlockDeleteCluster() error = %v, vimErr %v", err, tt.wantErr)
						return
					}

					assert.NotNil(t, deleteTask)
				}
			}
		})
	}
}

var newManagementCluster = `---
name: unittest
clusterPassword: VMware1!
clusterTemplateId: "55e69a3c-d92b-40ca-be51-9c6585b89ad7"
clusterType: MANAGEMENT
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.201
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01 
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
    - name: tkg
      type: Folder
    - name: vsanDatastore
      type: Datastore
    - name: k8s
      type: ResourcePool
    - name: hubsite
      type: ClusterComputeResource`

var WorkloadCluster = `
---
name: edge-workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var YamlBrokenWorkload = `
---
name: edge-workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
      - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`
var CreateDeleteTest01 = `
---
name: edge-workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

// no template id
var NewManagementClusterFailCase01 = `---
name: edge-mgmt-test01
clusterPassword: VMware1!
clusterType: MANAGEMENT
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01 
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
    - name: tkg
      type: Folder
    - name: vsanDatastore
      type: Datastore
    - name: k8s
      type: ResourcePool
    - name: hubsite
      type: ClusterComputeResource`

var NewManagementClusterFailCase02 = `---
name: unit_test
clusterPassword: VMware1!
clusterTemplateId: "55e69a3c-d92b-40ca-be51-9c6585b89ad7"
clusterType: MANAGEMENT
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01 
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
    - name: tkg
      type: Folder
    - name: vsanDatastore
      type: Datastore
    - name: k8s
      type: ResourcePool
    - name: hubsite
      type: ClusterComputeResource`

var InvalidMgmtCluster = `
---
name: edge-workload-test01
managementClusterId: wrong
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var WrongMgmtTemplateId = `
---
name: edge-workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
clusterTemplateId: wrongtemplate
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var InvalidDatastoreGlobal = `
---
name: workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastoreWrong
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var InvalidDatastoreWorker = `
---
name: workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastoreWrong
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastoreWrong
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var InvalidDatastoreMasterNode = `
---
name: workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastoreWrong
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var InvalidFolderGlobal = `
---
name: workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: WrongFolder
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var InvalidMasterNodeFolder = `
---
name: workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: Wrong
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var InvalidWorkerNodeFolder = `
---
name: workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: WrongFolder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: Wrong
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var InvalidDatastoreUrl = `
---
name: workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: WrongFolder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkgIn
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var InvalidNetworkPath = `
---
name: workload-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
# we can use name or id c3e006c1-e6aa-4591-950b-6f3bedd944d3
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            datastoreUrl: ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/
    tools:
        - name: harbor
          properties:
            extensionId: 9d0d4ff4-1963-4d89-ac15-2d856768deeb
            type: extension
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.189
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: WrongFolder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`

var newTestWorkloadCluster = `
name: edge-test01
managementClusterId: edge-mgmt-test01
clusterPassword: VMware1!
clusterTemplateId: myworkload
clusterType: workload
clusterConfig:
    csi:
        - name: nfs_client
          properties:
            serverIP: 10.241.0.250
            mountPath: w3-nfv-pst-01-mus
        - name: vsphere-csi
          properties:
            # you will need adjust that
            datastoreUrl: ds:///vmfs/volumes/vsan:528724284ea01639-d098d64191b96c2a/
            datastoreName: "vsanDatastore"
hcxCloudUrl: https://tca-pod03-cp.cnfdemo.io
endpointIP: 10.241.7.190
vmTemplate: photon-3-kube-v1.20.4+vmware.1
masterNodes:
    - name: master
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
workerNodes:
    - name: default-pool01
      networks:
        - label: MANAGEMENT
          networkName: tkg-dhcp-vlan1007-10.241.7.0
          nameservers:
            - 10.246.2.9
      placementParams:
        - name: tkg
          type: Folder
        - name: vsanDatastore
          type: Datastore
        - name: k8s
          type: ResourcePool
        - name: hubsite
          type: ClusterComputeResource
placementParams:
  - name: tkg
    type: Folder
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
  - name: hubsite
    type: ClusterComputeResource
`
