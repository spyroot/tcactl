package response

import (
	"github.com/nsf/jsondiff"
	_ "github.com/spyroot/tcactl/lib/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
	"testing"
)

var TestNodePoolSpecJson = `
{
  "name": "test_cluster02",
  "storage": 50,
  "cpu": 2,
  "memory": 16384,
  "replica": 1,
  "labels": [
    "type=test_cluster02"
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

var TestNodePoolSpecYaml = `
name: test_cluster02
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

func TestNewNodePoolSpecs(t *testing.T) {

	tests := []struct {
		name       string
		args       io.Reader
		wantErr    bool
		difference int
		spec       string
		format     string
	}{

		{
			name:       "Json basic positive case",
			args:       strings.NewReader(TestNodePoolSpecJson),
			wantErr:    false,
			spec:       TestNodePoolSpecJson,
			difference: 0,
			format:     "json",
		},
		{
			name:       "Yaml basic positive case",
			args:       strings.NewReader(TestNodePoolSpecYaml),
			wantErr:    false,
			difference: 0,
			spec:       TestNodePoolSpecYaml,
			format:     "yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNodePoolSpecs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNodePoolSpecs() error = %v, wantErr %v got back %v", err, tt.wantErr, got)
				return
			}

			var b []byte
			a, err := got.AsJsonString()
			assert.NoError(t, err)
			if tt.format == "json" {
				b = []byte(tt.spec)
			}

			// if input yaml validate that both json and yaml produce same output
			if tt.format == "yaml" {
				var sourceSpec NodesSpecs
				err = yaml.Unmarshal([]byte(tt.spec), &sourceSpec)
				assert.NoError(t, err)
				js, err := sourceSpec.AsJsonString()
				assert.NoError(t, err)
				b = []byte(js)
			}

			opt := jsondiff.DefaultJSONOptions()
			diff, diffStr := jsondiff.Compare([]byte(a), b, &opt)

			if tt.wantErr == false {
				if int(diff) != tt.difference {
					t.Errorf("NewNodePoolSpecs() "+
						"accepted difference %d actual difference %d, diff str %s", tt.difference, int(diff), diffStr)
					return
				}
			}
		})
	}
}

// TODO add other filters
func TestNodePool_Filter(t *testing.T) {
	tests := []struct {
		name        string
		args        io.Reader
		wantErr     bool
		expect      string
		filter      NodePoolFilterType
		query       string
		expectCount int
	}{
		{
			name:        "Basic positive filter by label",
			args:        strings.NewReader(TestNodePoolSpecJson),
			wantErr:     false,
			expectCount: 1,
			expect:      "test_cluster02",
			filter:      FilterByLabel,
			query:       "test_cluster02",
		},
		{
			name:        "Basic negative filter by label",
			args:        strings.NewReader(TestNodePoolSpecJson),
			wantErr:     false,
			expectCount: 0,
			expect:      "test_cluster02",
			filter:      FilterByLabel,
			query:       "test_cluster01",
		},
		{
			name:        "Basic negative filter by label",
			args:        strings.NewReader(TestNodePoolSpecJson),
			wantErr:     false,
			expectCount: 0,
			expect:      "test_cluster02",
			filter:      FilterPoolByID,
			query:       "test_cluster01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			n, err := NewNodePoolSpecs(tt.args)
			assert.NoError(t, err)
			pool := NewNodePool(n)
			assert.NotNil(t, pool)
			assert.Equal(t, len(pool.Pools), 1)
			assert.Equal(t, pool.Pools[0].Name, tt.expect)

			filtered, err := pool.Filter(tt.filter, func(s string) bool {
				return s == tt.query
			})

			assert.NoError(t, err)
			assert.NotNil(t, filtered)
			assert.Equal(t, tt.expectCount, len(filtered.Pools))

			if tt.expectCount == 0 {
				return
			}

			if tt.wantErr == false {
				p, err := filtered.GetPoolByName(tt.expect)
				assert.NoError(t, err)
				assert.Contains(t, p.Name, tt.expect)
			}
		})
	}
}

//
//func TestNodePool_GetIds(t *testing.T) {
//	type fields struct {
//		Pools []NodesSpecs
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   []string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			n := &NodePool{
//				Pools: tt.fields.Pools,
//			}
//			if got := n.GetIds(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetIds() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//
//func TestNodePool_GetPool(t *testing.T) {
//	type fields struct {
//		Pools []NodesSpecs
//	}
//	type args struct {
//		q string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *NodesSpecs
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			n := &NodePool{
//				Pools: tt.fields.Pools,
//			}
//			got, err := n.GetPool(tt.args.q)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetPool() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetPool() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNodeSpecsFromFile(t *testing.T) {
//	type args struct {
//		fileName string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *NodesSpecs
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NodeSpecsFromFile(tt.args.fileName)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NodeSpecsFromFile() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NodeSpecsFromFile() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNodeSpecsFromString(t *testing.T) {
//	type args struct {
//		str string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *NodesSpecs
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NodeSpecsFromString(tt.args.str)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NodeSpecsFromString() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NodeSpecsFromString() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNodesSpecs_GetField(t *testing.T) {
//	type fields struct {
//		CloneMode       string
//		Cpu             int
//		Id              string
//		Labels          []string
//		Memory          int
//		Name            string
//		Networks        []models.Network
//		PlacementParams []struct {
//			Name string `json:"name" yaml:"name"`
//			Type string `json:"type" yaml:"type"`
//		}
//		Replica                       int
//		Storage                       int
//		Config                        *NodePoolConfig
//		Status                        string
//		ActiveTasksCount              int
//		Nodes                         []models.Nodes
//		IsNodeCustomizationDeprecated bool
//	}
//	type args struct {
//		field string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			n := &NodesSpecs{
//				CloneMode:                     tt.fields.CloneMode,
//				Cpu:                           tt.fields.Cpu,
//				Id:                            tt.fields.Id,
//				Labels:                        tt.fields.Labels,
//				Memory:                        tt.fields.Memory,
//				Name:                          tt.fields.Name,
//				Networks:                      tt.fields.Networks,
//				PlacementParams:               tt.fields.PlacementParams,
//				Replica:                       tt.fields.Replica,
//				Storage:                       tt.fields.Storage,
//				Config:                        tt.fields.Config,
//				Status:                        tt.fields.Status,
//				ActiveTasksCount:              tt.fields.ActiveTasksCount,
//				Nodes:                         tt.fields.Nodes,
//				IsNodeCustomizationDeprecated: tt.fields.IsNodeCustomizationDeprecated,
//			}
//			if got := n.GetField(tt.args.field); got != tt.want {
//				t.Errorf("GetField() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNodesSpecs_GetFields(t *testing.T) {
//	type fields struct {
//		CloneMode       string
//		Cpu             int
//		Id              string
//		Labels          []string
//		Memory          int
//		Name            string
//		Networks        []models.Network
//		PlacementParams []struct {
//			Name string `json:"name" yaml:"name"`
//			Type string `json:"type" yaml:"type"`
//		}
//		Replica                       int
//		Storage                       int
//		Config                        *NodePoolConfig
//		Status                        string
//		ActiveTasksCount              int
//		Nodes                         []models.Nodes
//		IsNodeCustomizationDeprecated bool
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    map[string]interface{}
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			n := &NodesSpecs{
//				CloneMode:                     tt.fields.CloneMode,
//				Cpu:                           tt.fields.Cpu,
//				Id:                            tt.fields.Id,
//				Labels:                        tt.fields.Labels,
//				Memory:                        tt.fields.Memory,
//				Name:                          tt.fields.Name,
//				Networks:                      tt.fields.Networks,
//				PlacementParams:               tt.fields.PlacementParams,
//				Replica:                       tt.fields.Replica,
//				Storage:                       tt.fields.Storage,
//				Config:                        tt.fields.Config,
//				Status:                        tt.fields.Status,
//				ActiveTasksCount:              tt.fields.ActiveTasksCount,
//				Nodes:                         tt.fields.Nodes,
//				IsNodeCustomizationDeprecated: tt.fields.IsNodeCustomizationDeprecated,
//			}
//			got, err := n.GetFields()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetFields() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetFields() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNodesSpecs_GetNodesSpecs(t *testing.T) {
//	type fields struct {
//		CloneMode       string
//		Cpu             int
//		Id              string
//		Labels          []string
//		Memory          int
//		Name            string
//		Networks        []models.Network
//		PlacementParams []struct {
//			Name string `json:"name" yaml:"name"`
//			Type string `json:"type" yaml:"type"`
//		}
//		Replica                       int
//		Storage                       int
//		Config                        *NodePoolConfig
//		Status                        string
//		ActiveTasksCount              int
//		Nodes                         []models.Nodes
//		IsNodeCustomizationDeprecated bool
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   []models.Nodes
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			n := &NodesSpecs{
//				CloneMode:                     tt.fields.CloneMode,
//				Cpu:                           tt.fields.Cpu,
//				Id:                            tt.fields.Id,
//				Labels:                        tt.fields.Labels,
//				Memory:                        tt.fields.Memory,
//				Name:                          tt.fields.Name,
//				Networks:                      tt.fields.Networks,
//				PlacementParams:               tt.fields.PlacementParams,
//				Replica:                       tt.fields.Replica,
//				Storage:                       tt.fields.Storage,
//				Config:                        tt.fields.Config,
//				Status:                        tt.fields.Status,
//				ActiveTasksCount:              tt.fields.ActiveTasksCount,
//				Nodes:                         tt.fields.Nodes,
//				IsNodeCustomizationDeprecated: tt.fields.IsNodeCustomizationDeprecated,
//			}
//			if got := n.GetNodesSpecs(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetNodesSpecs() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNodesSpecs_InstanceSpecsFromString(t *testing.T) {
//	type fields struct {
//		CloneMode       string
//		Cpu             int
//		Id              string
//		Labels          []string
//		Memory          int
//		Name            string
//		Networks        []models.Network
//		PlacementParams []struct {
//			Name string `json:"name" yaml:"name"`
//			Type string `json:"type" yaml:"type"`
//		}
//		Replica                       int
//		Storage                       int
//		Config                        *NodePoolConfig
//		Status                        string
//		ActiveTasksCount              int
//		Nodes                         []models.Nodes
//		IsNodeCustomizationDeprecated bool
//	}
//	type args struct {
//		s string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    interface{}
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			n := NodesSpecs{
//				CloneMode:                     tt.fields.CloneMode,
//				Cpu:                           tt.fields.Cpu,
//				Id:                            tt.fields.Id,
//				Labels:                        tt.fields.Labels,
//				Memory:                        tt.fields.Memory,
//				Name:                          tt.fields.Name,
//				Networks:                      tt.fields.Networks,
//				PlacementParams:               tt.fields.PlacementParams,
//				Replica:                       tt.fields.Replica,
//				Storage:                       tt.fields.Storage,
//				Config:                        tt.fields.Config,
//				Status:                        tt.fields.Status,
//				ActiveTasksCount:              tt.fields.ActiveTasksCount,
//				Nodes:                         tt.fields.Nodes,
//				IsNodeCustomizationDeprecated: tt.fields.IsNodeCustomizationDeprecated,
//			}
//			got, err := n.SpecsFromString(tt.args.s)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("SpecsFromString() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("SpecsFromString() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPoolNotFound_Error(t *testing.T) {
//	type fields struct {
//		ErrMsg string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := &PoolNotFound{
//				ErrMsg: tt.fields.ErrMsg,
//			}
//			if got := m.Error(); got != tt.want {
//				t.Errorf("Error() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestReadNodeSpecs(t *testing.T) {
//	type args struct {
//		b io.Reader
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *NodesSpecs
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ReadNodeSpecs(tt.args.b)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ReadNodeSpecs() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ReadNodeSpecs() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

var NodePoolList = `{
  "items": [
    {
      "cpu": 4,
      "id": "6a18d59e-bc44-4ff7-a7b9-7c2d298d73c3",
      "labels": [
        "type=pool01"
      ],
      "memory": 131072,
      "name": "default-pool01",
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
          "name": "k8s",
          "type": "ResourcePool"
        }
      ],
      "replica": 1,
      "storage": 80,
      "config": {
        "cpuManagerPolicy": {
          "type": "kubernetes",
          "policy": "default"
        }
      },
      "status": "ACTIVE",
      "nodes": [
        {
          "ip": "10.241.7.143",
          "vmName": "edge-mgmt-test01-default-pool01-5855bcbdb4-rxm57"
        }
      ]
    },
    {
      "cpu": 4,
      "id": "dfe87837-b0ff-4b2a-9524-21ede3f61b03",
      "labels": [
        "type=pool01"
      ],
      "memory": 131072,
      "name": "default-pool01",
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
          "name": "k8s",
          "type": "ResourcePool"
        }
      ],
      "replica": 1,
      "storage": 80,
      "config": {
        "cpuManagerPolicy": {
          "type": "kubernetes",
          "policy": "default"
        }
      },
      "status": "ACTIVE",
      "nodes": [
        {
          "ip": "10.241.7.221",
          "vmName": "edge-test01-default-pool01-7797bdc978-chgrc"
        }
      ]
    },
    {
      "cpu": 8,
      "id": "95596645-6c79-433f-adc4-cb6687c445af",
      "labels": [
        "site=atlanta"
      ],
      "memory": 16384,
      "name": "atlanta",
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
      "replica": 1,
      "storage": 80,
      "config": {
        "cpuManagerPolicy": {
          "type": "kubernetes",
          "policy": "default"
        }
      },
      "status": "ACTIVE"
    },
    {
      "cpu": 8,
      "id": "c27f3fdc-0871-49df-bbce-fb42138940a8",
      "labels": [
        "cite=dallas"
      ],
      "memory": 16384,
      "name": "dallas",
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
      "storage": 80,
      "config": {
        "cpuManagerPolicy": {
          "type": "kubernetes",
          "policy": "default"
        },
        "healthCheck": {
          "nodeStartupTimeout": "20m",
          "unhealthyConditions": [
            {
              "type": "Ready",
              "status": "True",
              "timeout": "15m"
            },
            {
              "type": "MemoryPressure",
              "status": "True",
              "timeout": "15m"
            },
            {
              "type": "DiskPressure",
              "status": "True",
              "timeout": "300s"
            },
            {
              "type": "PIDPressure",
              "status": "True",
              "timeout": "300s"
            },
            {
              "type": "NetworkUnavailable",
              "status": "True",
              "timeout": "300s"
            }
          ]
        }
      },
      "status": "ACTIVE"
    },
    {
      "cpu": 2,
      "id": "5806bb2c-a50e-4269-aa1b-d13b902d6349",
      "labels": [
        "type=hub"
      ],
      "memory": 16384,
      "name": "temp123",
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
      "replica": 1,
      "storage": 50,
      "config": {},
      "status": "ACTIVE"
    },
    {
      "cpu": 2,
      "id": "eaeab452-793f-44bb-8ba1-209f1708abd8",
      "labels": [
        "type=0c720ead-682"
      ],
      "memory": 16384,
      "name": "0c720ead-682",
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
      "replica": 1,
      "storage": 50,
      "config": {},
      "status": "ACTIVE"
    },
    {
      "cpu": 2,
      "id": "1b049d46-ea20-4b1b-957e-3e78d26f460b",
      "labels": [
        "type=275485d4-baf"
      ],
      "memory": 16384,
      "name": "275485d4-baf",
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
      "replica": 1,
      "storage": 50,
      "config": {},
      "status": "ACTIVE"
    },
    {
      "cpu": 2,
      "id": "3acf9b79-f8e5-41d6-997b-58792d3955bb",
      "labels": [
        "type=test_cluster"
      ],
      "memory": 16384,
      "name": "test-cluster01",
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
      "replica": 1,
      "storage": 50,
      "config": {},
      "status": "ACTIVE"
    }
  ]
}`
