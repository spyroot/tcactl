package specs

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

//Read spec from file and test
func TestNewNodePoolSpec_SpecsFromFile(t *testing.T) {

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name: "Read instance workload spec from yaml",
			file: "/pool/positive/pool.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file
			spec, err := SpecNodePool{}.SpecsFromFile(fileName)
			if tt.wantErr && err == nil {
				t.Errorf("Test failed must not return error")
				return
			}
			if tt.wantErr && err != nil {
				return
			}

			if spec == nil {
				t.Errorf("SpecsFromFile() return nil spec")
				return
			}

			poolSpec, ok := (*spec).(*SpecNodePool)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = poolSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() Test failed validator "+
					"return error for positive case err %v file %s", err, fileName)
				return
			}
		})
	}
}

// Reads nodes spec from string and test
func TestNewNodePoolSpec_SpecsFromString(t *testing.T) {

	tests := []struct {
		name          string
		file          string
		wantErr       bool
		wantValidaErr bool
	}{
		{
			name:    "Read instance spec from yaml",
			file:    "/pool/positive/pool.yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file

			b, err := ioutil.ReadFile(fileName)
			assert.NoError(t, err)

			spec, err := SpecNodePool{}.SpecsFromString(string(b))
			if tt.wantErr && err == nil {
				t.Errorf("Test failed must not return error")
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if spec == nil {
				t.Errorf("SpecsFromFile() return nil spec")
				return
			}
			instanceSpec, ok := (*spec).(*SpecNodePool)
			if !ok {
				t.Errorf("SpecsFromString() failed method return wrong type")
				return
			}
			err = instanceSpec.Validate()
			if tt.wantValidaErr && err == nil {
				t.Errorf("SpecsFromString() failed spec validator return no error for negative case %v", err)
				return
			}
			if tt.wantValidaErr && err != nil {
				t.Log(err)
				return
			}
		})
	}
}

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
