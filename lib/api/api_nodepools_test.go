package api

import (
	"encoding/json"
	"github.com/nsf/jsondiff"
	_ "github.com/nsf/jsondiff"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/io"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
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
				t.Errorf("ReadNodeSpecFromString() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr == false {
				t.Errorf("ReadNodeSpecFromString() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			newJson, err := json.Marshal(got)
			assert.NoError(t, err)

			oldJson, err := json.Marshal(got)
			assert.NoError(t, err)

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

// specNodePoolStringReaderHelper test helper return node spec
func specNodePoolStringReaderHelper(s string) *request.NewNodePoolSpec {
	r, err := ReadNodeSpecFromString(s)
	io.CheckErr(err)
	return r
}

// specNodePoolStringReaderHelper test helper return node spec
func specNodePoolFromFile(spec string) *request.NewNodePoolSpec {

	tmpFile, err := ioutil.TempFile("", "tcactltest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write to file,close it and read spec
	if _, err = tmpFile.Write([]byte(spec)); err != nil {
		io.CheckErr(err)
	}

	if err := tmpFile.Close(); err != nil {
		io.CheckErr(err)
	}

	// read from file
	r, err := ReadNodeSpecFromFile(tmpFile.Name())
	io.CheckErr(err)

	return r
}

//  Test create a new node pool
//
func TestTcaCreateNewNodePool(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		spec    *request.NewNodePoolSpec
		wantErr bool
		reset   bool
	}{
		{
			name:    "Create node pool from json",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(jsonNodeSpec),
			wantErr: false,
		},
		{
			name:    "Create node pool from yaml",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYaml),
			wantErr: false,
		},
		{
			name:    "Wrong spec no CPU",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlNoCPU),
			wantErr: true,
		},
		{
			name:    "Wrong spec no replica",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlNoCPU),
			wantErr: false,
		},
		{
			name:    "Wrong spec no replica",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlNoReplica),
			wantErr: false,
		},
		{
			name:    "Wrong spec no replica",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlNoReplica),
			wantErr: false,
		},
		{
			name:    "Wrong spec no replica",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlNoNetwork),
			wantErr: false,
		},
	}

	SetLoggingFlags()

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

			if tt.spec != nil {
				tt.spec.Name = generateName()
				tt.spec.Labels[0] = "type=" + tt.spec.Name
			}

			if task, err = a.CreateNewNodePool(tt.spec, getTestClusterName(),
				false, true, true); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewNodePool() error = %v, vimErr %v", err, tt.wantErr)
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

// Update test
//
func TestTcaUpdateNodePool(t *testing.T) {

	tests := []struct {
		name    string
		rest    *client.RestClient
		spec    *request.NewNodePoolSpec
		wantErr bool
		reset   bool
	}{
		{
			name:    "Yaml wrong file.",
			rest:    rest,
			spec:    specNodePoolFromFile("test"),
			wantErr: true,
		},
		{
			name:    "Yaml template wrong cpu count.",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(updateMustWrongCPUCount),
			wantErr: true,
		},
		{
			name:    "Yaml template no cpu value.",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(updateMustFailNoCPU),
			wantErr: true,
		},
		{
			name:    "Yaml without pool id.",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(updateMinYamlNoID),
			wantErr: false,
		},
		{
			name:    "Read from yaml file without node pool id.",
			rest:    rest,
			spec:    specNodePoolFromFile(updateMinYamlNoID),
			wantErr: false,
		},
		{
			name:    "Read from yaml file and add label",
			rest:    rest,
			spec:    specNodePoolFromFile(updatePoolAddLabel),
			wantErr: false,
		},
		{
			name:    "Yaml updates node pool replica",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(updateYamlPoolSpec),
			wantErr: false,
		},
		{
			name:    "Yaml min update replica",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(updateMinYamlPoolSpec),
			wantErr: false,
		},
		{
			name:    "Json updates node pool replica",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(updateJsonPoolSpec),
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
			SetLoggingFlags()

			if task, err = a.UpdateNodePool(tt.spec,
				getTestClusterName(),
				getTestNodePoolName(), false, true, true); (err != nil) != tt.wantErr {
				t.Errorf("TestTcaUpdateNodePool() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Errorf("TestTcaUpdateNodePool() error is nil, vimErr %v", tt.wantErr)
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

//
//  "Only Update of labels, replicas, and machine health check for node pools supported."
var updateYamlPoolSpec = `
# in this example we just update replica for existing pool
name: test-cluster01
cpu: 2
id: 3acf9b79-f8e5-41d6-997b-58792d3955bb
labels:
  - type=test_cluster01
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
status: ACTIVE,
activeTasksCount: 0
isNodeCustomizationDeprecated: false
`

//
//  "Only Update of labels, replicas, and machine health check for node pools supported."
var updateMinYamlPoolSpec = `
# in this example we just update replica for existing pool
name: test-cluster01
cpu: 2
id: 3acf9b79-f8e5-41d6-997b-58792d3955bb
labels:
  - type=test_cluster01
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
status: ACTIVE,
isNodeCustomizationDeprecated: false
`

//
//  "Only Update of labels, replicas, and machine health check for node pools supported."
var updateMinYamlNoID = `
# in this example we just update replica for existing pool
name: test-cluster01
cpu: 2
labels:
  - type=test_cluster01
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
status: ACTIVE,
isNodeCustomizationDeprecated: false
`

//
//  "Only Update of labels, replicas, and machine health check for node pools supported."
var updatePoolAddLabel = `
# in this example we just update replica for existing pool
name: test-cluster01
cpu: 2
labels:
  - type=test_cluster01
  - loc=paloalot
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
status: ACTIVE,
isNodeCustomizationDeprecated: false
`

//
//  "Only Update of labels, replicas, and machine health check for node pools supported."
var updateMustWrongCPUCount = `
# in this example we just update replica for existing pool
name: test-cluster01
cpu: 3
id: 3acf9b79-f8e5-41d6-997b-58792d3955bb
labels:
  - type=test_cluster01
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
status: ACTIVE,
isNodeCustomizationDeprecated: false
`

//
//  "Only Update of labels, replicas, and machine health check for node pools supported."
var updateMustFailNoCPU = `
# in this example we just update replica for existing pool
name: test-cluster01
id: 3acf9b79-f8e5-41d6-997b-58792d3955bb
labels:
  - type=test_cluster01
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
status: ACTIVE,
isNodeCustomizationDeprecated: false
`

var updateJsonPoolSpec = `
{
    "name": "test-cluster01",
    "id": "3acf9b79-f8e5-41d6-997b-58792d3955bb",
    "cpu": 2,
    "memory": 16384,
    "labels": [
        "type=test_cluster"
    ],
    "networks": [
        {
            "label": "MANAGEMENT",
            "networkName": "/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0",
            "nameservers": [
                "10.246.2.9"
            ]
        }
    ],
    "placementParams": [
        {
            "name": "hubsite",
            "type": "ClusterComputeResource"
        },
        {
            "name": "vsanDatastore",
            "type": "Datastore"
        },
        {
            "name": "/Datacenter/host/hubsite/Resources/k8s",
            "type": "ResourcePool"
        }
    ],
    "replica": 2,
    "storage": 50,
    "config": {},
    "status": "ACTIVE",
    "activeTasksCount": 0,
    "nodes": [],
    "isNodeCustomizationDeprecated": false
}
`

var newNodePoolYaml = `
cpu: 2
labels:
  - type=test_cluster02
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
`

var newNodePoolYamlNoCPU = `
labels:
  - type=test_cluster02
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
`

var newNodePoolYamlNoReplica = `
cpu: 2
labels:
  - type=test_cluster02
memory: 16384
networks:
  - label: MANAGEMENT
    networkName: /Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0
    nameservers:
      - 10.246.2.9
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
storage: 50
`

var newNodePoolYamlNoNetwork = `
labels:
  - type=test_cluster02
memory: 16384
networks:
  - label: MANAGEMENT
placementParams:
  - name: hubsite
    type: ClusterComputeResource
  - name: vsanDatastore
    type: Datastore
  - name: k8s
    type: ResourcePool
replica: 1
storage: 50
`
