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
	"encoding/json"
	"github.com/nsf/jsondiff"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/lib/testutil"
	iotuils "github.com/spyroot/tcactl/pkg/io"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

// specNodePoolStringReaderHelper helper return node specString
func specClusterStringReaderHelper(s string) *request.Cluster {
	r, err := request.ClusterSpecsFromString(s)
	iotuils.CheckErr(err)
	return r
}

// specNodePoolStringReaderHelper helper return node specString
func specClusterFromFile(spec string) *request.NewNodePoolSpec {

	tmpFile, err := ioutil.TempFile("", "tcactltest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write to file,close it and read specString
	if _, err = tmpFile.Write([]byte(spec)); err != nil {
		iotuils.CheckErr(err)
	}

	if err := tmpFile.Close(); err != nil {
		iotuils.CheckErr(err)
	}

	// read from file
	r, err := ReadNodeSpecFromFile(tmpFile.Name())
	iotuils.CheckErr(err)

	return r
}

// TestClusterSpecFromString read specString and validate parser
func TestClusterSpecFromString(t *testing.T) {

	tests := []struct {
		name    string
		spec    string
		wantErr bool
	}{
		{
			name:    "Read Yaml management cluster specString",
			spec:    NewManagementCluster,
			wantErr: false,
		},
		{
			name:    "Read Yaml workload cluster specString",
			spec:    WorkloadCluster,
			wantErr: false,
		},
		{
			name:    "Read Yaml broken workload cluster specString",
			spec:    YamlBrokenWorkload,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ReadNodeSpecFromString(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNodeSpecFromString() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			t.Log(err)

			if err != nil && tt.wantErr == false {
				t.Errorf("ReadNodeSpecFromString() error = %v, vimErr %v", err, tt.wantErr)
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

			opt := jsondiff.DefaultJSONOptions()
			diff, _ := jsondiff.Compare(newJson, oldJson, &opt)

			if tt.wantErr != true && diff > 0 {
				t.Errorf("ReadNodeSpecFromString() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			t.Logf("diff %d", diff)
		})
	}
}

// TestClusterSpecFromFile read specString and validate parser
func TestClusterSpecFromFile(t *testing.T) {

	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "Read  yaml management specString file",
			fileName: testutil.SpecTempFileName(NewManagementCluster),
			wantErr:  false,
		},
		{
			name:     "Read yaml workload specString file",
			fileName: testutil.SpecTempFileName(WorkloadCluster),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := request.ClusterSpecsFromFile(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestClusterSpecFromFile() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr == false {
				t.Errorf("TestClusterSpecFromFile() error = %v, vimErr %v", err, tt.wantErr)
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

			opt := jsondiff.DefaultJSONOptions()
			diff, _ := jsondiff.Compare(newJson, oldJson, &opt)

			if tt.wantErr != true && diff > 0 {
				t.Errorf("ReadNodeSpecFromString() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			t.Logf("diff %d", diff)
		})
	}
}

// Reads specString and validate parser
func TestClusterSpecFromReader(t *testing.T) {

	tests := []struct {
		name    string
		reader  io.Reader
		wantErr bool
	}{
		{
			name:    "Read yaml management specString from io.reader",
			reader:  testutil.SpecTempReader(NewManagementCluster),
			wantErr: false,
		},
		{
			name:    "Read Basic yaml workload specString from io.reader",
			reader:  testutil.SpecTempReader(WorkloadCluster),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := request.ReadClusterSpec(tt.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestClusterSpecFromFile() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr == false {
				t.Errorf("TestClusterSpecFromFile() error = %v, vimErr %v", err, tt.wantErr)
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

			opt := jsondiff.DefaultJSONOptions()
			diff, _ := jsondiff.Compare(newJson, oldJson, &opt)

			if tt.wantErr != true && diff > 0 {
				t.Errorf("ReadNodeSpecFromString() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			t.Logf("diff %d", diff)
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
// Test uses existing cluster ID and name
// make sure they adjusted or passed via env variable
func TestGetCluster(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		wantErr bool
		cluster string
		expect  string
	}{
		{
			name:    "Get cluster from name.",
			rest:    rest,
			cluster: getTestClusterName(),
			expect:  getTestClusterId(),
			wantErr: false,
		},
		{
			name:    "Get cluster from id.",
			rest:    rest,
			cluster: getTestClusterId(),
			expect:  getTestClusterId(),
			wantErr: false,
		},
		{
			name:    "Get cluster wrong id.",
			rest:    rest,
			cluster: "868636c9-868f-49fb-a6df-6a0d2d137141",
			expect:  getTestClusterId(),
			wantErr: true,
		},
		{
			name:    "Get cluster wrong name.",
			rest:    rest,
			cluster: "test",
			expect:  getTestClusterId(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			a.SetTrace(true)
			SetLoggingFlags()

			actual, err := a.GetCluster(tt.cluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCluster() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, actual)
				assert.Equal(t, tt.expect, actual.Id)
			}
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
func TestGetClusterNodePool(t *testing.T) {

	tests := []struct {
		name     string
		rest     *client.RestClient
		wantErr  bool
		cluster  string
		nodepool string
		expect   string
		trace    bool
	}{
		{
			name:     "Must resolve from cluster name and pool name.",
			rest:     rest,
			cluster:  getTestClusterName(),
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
			nodepool: getTestNodePoolId(),
			expect:   getTestNodePoolId(),
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

			a.SetTrace(tt.trace)
			SetLoggingFlags()

			actual, err := a.GetClusterNodePool(tt.cluster, tt.nodepool)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterNodePool() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, actual)
				assert.Equal(t, tt.expect, actual.Id)
			}

			if tt.wantErr {
				t.Log(err)
			}
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
func TestTcaApi_GetCurrentClusterTask(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		reader  io.Reader
		wantErr bool
		reset   bool
		trace   bool
		sec     time.Duration
	}{
		{
			name:    "Create cluster and check task list.",
			rest:    rest,
			reader:  testutil.SpecTempReader(NewManagementCluster),
			wantErr: false,
			trace:   false,
			sec:     3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				err  error
				task *models.TcaTask
			)

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			a.SetTrace(tt.trace)
			SetLoggingFlags()

			spec, err := request.ReadClusterSpec(tt.reader)

			assert.NoError(t, err)
			assert.NotNil(t, spec)

			if spec.IsManagement() {
				spec.ClusterTemplateId = getTestMgmtTemplateId()
			} else {
				spec.ClusterTemplateId = getTestWorkloadTemplateId()
			}

			assert.NoError(t, err)
			if task, err = a.CreateClusters(spec, false, false, true); (err != nil) != tt.wantErr {
				t.Errorf("CreateClusters() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			time.Sleep(tt.sec * time.Second)

			if tt.wantErr && err == nil {
				t.Errorf("CreateClusters() error is nil, vimErr %v", tt.wantErr)
				return
			}

			if err == nil {
				assert.NotNil(t, task)
				clusterCreateTask, err := a.GetCurrentClusterTask(task.Id)
				assert.NoError(t, err)
				if _, err := clusterCreateTask.FindEntityByName(spec.Name); (err != nil) != tt.wantErr {
					t.Errorf("GetCurrentClusterTask() error = %v, vimErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
func TestTcaApi_CreateClusters(t *testing.T) {

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
	}{
		{
			name:            "Create mgmt dryRun run.",
			rest:            rest,
			reader:          testutil.SpecTempReader(NewManagementCluster),
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
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong datastore",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidDatastoreMasterNode),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidFolderGlobal),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong master node folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidMasterNodeFolder),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
		{
			name:            "Create cluster with wrong master node folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidWorkerNodeFolder),
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
			name:            "Create cluster with wrong master node folder",
			rest:            rest,
			reader:          testutil.SpecTempReader(InvalidDatastoreUrl),
			wantErr:         true,
			dryRun:          true,
			trace:           false,
			useTestTemplate: true,
			sec:             3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				err  error
				task *models.TcaTask
			)

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			a.SetTrace(tt.trace)
			SetLoggingFlags()

			spec, err := request.ReadClusterSpec(tt.reader)

			assert.NoError(t, err)
			assert.NotNil(t, spec)

			if tt.useTestTemplate {
				if spec.IsManagement() {
					spec.ClusterTemplateId = getTestMgmtTemplateId()
				} else {
					spec.ClusterTemplateId = getTestWorkloadTemplateId()
				}
			}

			assert.NoError(t, err)
			if task, err = a.CreateClusters(spec, tt.dryRun, false, true); (err != nil) != tt.wantErr {
				t.Errorf("CreateClusters() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				t.Log(err)
			}

			if !tt.wantErr {
				if tt.dryRun {
					// fail case in dry run
					if len(task.Id) > 0 {
						t.Errorf("CreateClusters() in dry run task must be empty")
						return
					}
					t.Logf("Dry run passed.")
					return
				}

				if !tt.dryRun {
					t.Logf("task created %s", task.Id)
					return
				}
			}
		})
	}
}

// TestTcaApi_GetClusterTask validate running task list
func TestTcaApi_CreateBlockDeleteCluster(t *testing.T) {

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
	}{
		{
			name:            "Create and Delete a new tenant cluster.",
			rest:            rest,
			reader:          testutil.SpecTempReader(CreateDeleteTest01),
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

			a, err := NewTcaApi(tt.rest)
			assert.NoError(t, err)

			a.SetTrace(tt.trace)
			SetLoggingFlags()

			spec, err := request.ReadClusterSpec(tt.reader)
			assert.NoError(t, err)
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

			createTask, err = a.CreateClusters(spec, tt.dryRun, true, true)

			// if task failed we still should be able delete cluster.
			if _, ok := err.(*TcaTaskFailed); ok {
				t.Log("Create cluster task failed but test will try to delete it")
				if deleteTask, err = a.DeleteCluster(spec.Name, true, true); (err != nil) != tt.wantErr {
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

					if deleteTask, err = a.DeleteCluster(spec.Name, true, true); (err != nil) != tt.wantErr {
						t.Errorf("CreateBlockDeleteCluster() error = %v, vimErr %v", err, tt.wantErr)
						return
					}

					assert.NotNil(t, deleteTask)
				}
			}
		})
	}
}

var NewManagementCluster = `---
name: unit_test
clusterPassword: VMware1!
clusterTemplateId: "55e69a3c-d92b-40ca-be51-9c6585b89ad7"
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
          type: ClusterComputeResourcel
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
