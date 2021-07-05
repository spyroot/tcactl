package api

import (
	"github.com/google/uuid"
	"os"
)

const (
	// test name, if you already have cluster adjust or pass TCA_TEST_CLUSTER env
	testClusterName = "edge-test01"

	// test cloud provider name
	testCloudName = "edge"

	// test
	testPoolName = "test-cluster01"

	// testTenantId
	testTenantId = "BDC07231F50A4536AA6DCF6B8C04BA5C"

	//
	testClusterID = "868636c9-868f-49fb-a6df-6a0d2d137146"

	//
	testNodePoolId = "3acf9b79-f8e5-41d6-997b-58792d3955bb"

	// testMgmtTemplateId test mgmt template
	testMgmtTemplateId = "55e69a3c-d92b-40ca-be51-9c6585b89ad7"

	// testWorkloadTemplateId test template
	testWorkloadTemplateId = "c3e006c1-e6aa-4591-950b-6f3bedd944d3"
)

func getTestClusterId() string {
	testCluster := os.Getenv("TCA_TEST_CLUSTER_ID")
	if len(testCluster) == 0 {
		return testClusterID
	}
	return testCluster
}

func getTenantId() string {
	testTenantId := os.Getenv("TCA_TEST_TENANT_ID")
	if len(testTenantId) == 0 {
		return testTenantId
	}
	return testTenantId
}

func getTestClusterName() string {
	testCluster := os.Getenv("TCA_TEST_CLUSTER")
	if len(testCluster) == 0 {
		return testClusterName
	}
	return testCluster
}

func getTestNodePoolName() string {
	testCluster := os.Getenv("TCA_TEST_NODE_POOL")
	if len(testCluster) == 0 {
		return testPoolName
	}
	return testCluster
}

func getTestNodePoolId() string {
	testCluster := os.Getenv("TCA_TEST_NODE_POOL_ID")
	if len(testCluster) == 0 {
		return testNodePoolId
	}
	return testCluster
}

func getTestMgmtTemplateId() string {
	env := os.Getenv("TCA_TEST_MGMT_TEMPLATE_ID")
	if len(env) == 0 {
		return testMgmtTemplateId
	}
	return env
}

func getTestWorkloadTemplateId() string {
	env := os.Getenv("TCA_TEST_WORKLOAD_TEMPLATE_ID")
	if len(env) == 0 {
		return testWorkloadTemplateId
	}
	return env
}

// getTestCloudProvider should return a cloud provider
// that attached in TCA
func getTestCloudProvider() string {
	e := os.Getenv("TCA_TEST_CLOUD")
	if len(e) == 0 {
		return testCloudName
	}
	return e
}

// getTenantCluster - should return tenant that
// already present in TCA
func getTenantCluster() string {
	e := os.Getenv("TCA_TEST_TENANT")
	if len(e) == 0 {
		return testClusterName
	}
	return e
}

// generate name string
func generateName() string {
	n := uuid.New().String()
	return n[0:12]
}

var yamlMgmtTemplate = `
clusterType: MANAGEMENT
clusterConfig:
    kubernetesVersion: v1.20.4+vmware.1
masterNodes:
    - cpu: 4
      memory: 16384
      name: master
      networks:
        - label: MANAGEMENT
      storage: 50
      replica: 1
      labels: []
      cloneMode: linkedClone
name: min
workerNodes:
    - cpu: 4
      memory: 131072
      name: default-pool01
      networks:
        - label: MANAGEMENT
      storage: 80
      replica: 1
      labels:
        - type=pool01
      cloneMode: linkedClone
      config:
        cpuManagerPolicy:
            type: kubernetes
            policy: default
`

var yamlInvalidMgmtTemplate = `
clusterType: MANAGEMENT
clusterConfig:
    kubernetesVersion: v1.20.4+vmware.1
name: min
workerNodes:
    - cpu: 4
      memory: 131072
      name: default-pool01
      networks:
        - label: MANAGEMENT
      storage: 80
      replica: 1
      labels:
        - type=pool01
      cloneMode: linkedClone
      config:
        cpuManagerPolicy:
            type: kubernetes
            policy: default
`

var yamlInvalidMgmtTemplate2 = `
clusterType: MANAGEMENT
clusterConfig:
    kubernetesVersion: v1.20.4+vmware.1
masterNodes:
    - cpu: 4
      memory: 16384
      name: master
      networks:
        - label: MANAGEMENT
      storage: 50
      replica: 1
      labels: []
      cloneMode: linkedClone
name: min
`

var yamlInvalidMgmtTemplate3 = `
clusterType: 
clusterConfig:
    kubernetesVersion: v1.20.4+vmware.1
masterNodes:
    - cpu: 4
      memory: 16384
      name: master
      networks:
        - label: MANAGEMENT
      storage: 50
      replica: 1
      labels: []
      cloneMode: linkedClone
name: min
workerNodes:
    - cpu: 4
      memory: 131072
      name: default-pool01
      networks:
        - label: MANAGEMENT
      storage: 80
      replica: 1
      labels:
        - type=pool01
      cloneMode: linkedClone
      config:
        cpuManagerPolicy:
            type: kubernetes
            policy: default
`

var yamlInvalidMgmtTemplate4 = `
clusterType: MANAGEMENT
clusterConfig:
    kubernetesVersion: v1.20.4+vmware.1
masterNodes:
    - cpu: 4
      memory: 16384
      name: master
      networks:
        - label: MANAGEMENT
      storage: 50
      replica: 1
      labels: []
      cloneMode: linkedClone
name: min
workerNodes:
    - cpu: 4
      memory: 131072
      name: default-pool01
      networks:
        - label: 
      storage: 80
      replica: 1
      labels:
        - type=pool01
      cloneMode: linkedClone
      config:
        cpuManagerPolicy:
            type: kubernetes
            policy: default
`

var yamlWorkloadTemplate = `
clusterType: WORKLOAD
clusterConfig:
    cni:
        - name: multus
        - name: calico
    csi:
        - name: vsphere-csi
          properties:
            name: vsphere-sc
            isDefault: true
            timeout: "300"
        - name: nfs_client
          properties:
            name: nfs-client
            isDefault: false
    kubernetesVersion: v1.20.4+vmware.1
    tools:
        - name: helm
          version: 2.17.0
description: ""
masterNodes:
    - cpu: 4
      memory: 16384
      name: master
      networks:
        - label: MANAGEMENT
      storage: 50
      replica: 1
      labels: []
      cloneMode: linkedClone
name: myworkload
workerNodes:
    - cpu: 4
      memory: 131072
      name: default-pool01
      networks:
        - label: MANAGEMENT
      storage: 80
      replica: 1
      labels:
        - type=pool01
      cloneMode: linkedClone
      config:
        cpuManagerPolicy:
            type: kubernetes
            policy: default
`

var yamlWorkloadEmpty = `

`

var jsonNodeSpec = `
{
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
