package api

import (
	"context"
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

// specNodePoolStringReaderHelper test helper return node specString
func specNodePoolStringReaderHelper(s string) *request.NewNodePoolSpec {
	r, err := request.ReadNodeSpecFromString(s)
	io.CheckErr(err)
	return r
}

// specNodePoolStringReaderHelper test helper return node specString
func specNodePoolFromFile(spec string) *request.NewNodePoolSpec {

	tmpFile, err := ioutil.TempFile("", "tcactltest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write to file,close it and read specString
	if _, err = tmpFile.Write([]byte(spec)); err != nil {
		io.CheckErr(err)
	}

	if err := tmpFile.Close(); err != nil {
		io.CheckErr(err)
	}

	// read from file
	r, err := request.ReadNodeSpecFromFile(tmpFile.Name())
	io.CheckErr(err)

	return r
}

//  Test create a new node pool
//
func TestTcaCreateNewNodePool(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		spec         *request.NewNodePoolSpec
		wantErr      bool
		reset        bool
		withoutlabel bool
	}{
		{
			name:    "Create node pool from json spec",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(jsonNodeSpec),
			wantErr: false,
		},
		{
			name:    "Create node pool from yaml spec",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYaml),
			wantErr: false,
		},
		{
			name:    "Wrong specString no CPU",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlNoCPU),
			wantErr: true,
		},
		{
			name:    "Wrong specString no replica",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlNoReplica),
			wantErr: true,
		},
		{
			name:    "Wrong specString no network",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlNoNetwork),
			wantErr: true,
		},
		{
			name:         "Wrong specString without label",
			rest:         rest,
			spec:         specNodePoolStringReaderHelper(newNodePoolYamlWithoutLabel),
			wantErr:      true,
			withoutlabel: true,
		},
		{
			name:    "Wrong specString without config",
			rest:    rest,
			spec:    specNodePoolStringReaderHelper(newNodePoolYamlWithoutConfig),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				ctx  = context.Background()
				err  error
				task *models.TcaTask
			)

			a := getTcaApi(t, rest, false)
			assert.NotNil(t, tt.spec)

			tt.spec.CloneMode = request.LinkedClone

			if tt.spec != nil && tt.withoutlabel == false {
				tt.spec.Name = generateName()
				tt.spec.Labels[0] = "type=" + tt.spec.Name
			}

			if task, err = a.CreateNewNodePool(ctx,
				&NodePoolCreateApiReq{
					Spec:       tt.spec,
					Cluster:    getTestWorkloadClusterName(),
					IsVerbose:  false,
					IsBlocking: false,
					IsDryRun:   false,
				}); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewNodePool() error = %v, vimErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Errorf("CreateNewNodePool() error is nil")
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if task == nil {
				t.Errorf("Task is nil")
				return
			}
			if len(task.OperationId) == 0 {
				t.Errorf("Task is not nil op id is empty")
				return
			}
		})
	}
}

// Update test
//
func TestTcaUpdateNodePool(t *testing.T) {

	tests := []struct {
		name      string
		rest      *client.RestClient
		spec      *request.NewNodePoolSpec
		wantErr   bool
		reset     bool
		doBlock   bool
		doVerbose bool
		doDryRun  bool
	}{
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

			ctx := context.Background()
			a := getTcaApi(t, rest, false)

			if tt.reset {
				a.rest = nil
			}

			tt.spec.CloneMode = request.LinkedClone
			SetLoggingFlags()

			req := NodePoolCreateApiReq{
				Spec:       tt.spec,
				Cluster:    getTestWorkloadClusterName(),
				IsVerbose:  tt.doVerbose,
				IsBlocking: tt.doBlock,
				IsDryRun:   tt.doDryRun,
			}
			if task, err = a.UpdateNodePool(ctx, &req); (err != nil) != tt.wantErr {
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

			if task == nil {
				t.Logf("Recieved correct error %v", err)
				return
			}

			if len(task.OperationId) == 0 {
				t.Logf("Recieved correct error %v", err)
				return
			}
		})
	}
}

//
//  "Only Update of labels, replicas, and machine health check for node pools supported."
var updateYamlPoolSpec = `
kind: node_pool
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
kind: node_pool
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
kind: node_pool
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
kind: node_pool
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

//  "Only Update of labels, replicas, and machine health check for node pools supported."
var updateMustWrongCPUCount = `
kind: node_pool
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
kind: node_pool
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
	"kind": "node_pool"
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
kind: node_pool
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
kind: node_pool
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
kind: node_pool
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
kind: node_pool
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

var jsonNodeSpec = `
{
	"kind": "node_pool",
    "name": "temp1234",
    "storage": 50,
    "cpu": 2,
    "memory": 16384,
    "replica": 1,
    "labels": [
        "type=hub"
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
            "type": "ClusterComputeResource",
            "name": "hubsite"
        },
        {
            "type": "Datastore",
            "name": "vsanDatastore"
        },
        {
            "type": "ResourcePool",
            "name": "k8s"
        }
    ],
    "config": {
        "cpuManagerPolicy": {
            "type": "kubernetes",
            "policy": "default"
        }
    }
}
`

var yamlNodeSpec = `
id: ""
clone_mode: ""
cpu: 2
labels:
    - type=hub
memory: 16384
name: temp
networks:
    - label: MANAGEMENT
      network_name: ""
      nameservers:
        - 10.246.2.9
placement_params: []
replica: 1
storage: 50
config:
    cpu_manager_policy:
        type: ""
        policy: ""
        properties:
            kube_reserved:
                cpu: 0
                memoryInGiB: 0
            system_reserved:
                cpu: 0
                memoryInGiB: 0
    health_check:
        nodeStartupTimeout: ""
        unhealthy_conditions: []
status: ""
active_tasks_count: 0
nodes: []
is_node_customization_deprecated: false
`

var newNodePoolYamlWithoutLabel = `
kind: node_pool
cpu: 2
labels:
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

var newNodePoolYamlWithoutConfig = `
id: ""
clone_mode: ""
cpu: 2
labels:
    - type=hub
memory: 16384
name: temp
networks:
    - label: MANAGEMENT
      network_name: ""
      nameservers:
        - 10.246.2.9
placement_params: []
replica: 1
storage: 50
status: ""
active_tasks_count: 0
nodes: []
is_node_customization_deprecated: false
`
