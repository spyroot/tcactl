package request

import (
	"encoding/json"
	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Reads specString and validate parser
func TestReadNodeSpecFromString(t *testing.T) {

	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{
			name:    "Read from json node spec",
			args:    jsonNodeSpec,
			wantErr: false,
		},
		{
			name:    "Read from yaml node spec",
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
