package api

import (
	"github.com/google/uuid"
	"os"
)

const (
	// test name, if you already have cluster adjust or pass TCA_TEST_CLUSTER env
	testWorkloadClusterName = "edge-test01"

	// testMgmtClusterName
	testMgmtClusterName = "edge-mgmt-test01"

	// test cloud provider name
	testCloudName = "edge"

	// test
	testPoolName = "default-pool01"

	// testCatalogName
	testCatalogName = "unit_test"

	// testInstanceName
	testInstanceName = "unit_test_instance"

	//testRepoName
	testRepoName = "Repo"

	//testRepoUrl
	testRepoUrl = "Repo.cnfdemo.io/chartrepo/library"

	//testRepoUsername
	testRepoUsername = "admin"

	//testRepoPassword
	testRepoPassword = ""

	// testTenantId
	testTenantId = "BDC07231F50A4536AA6DCF6B8C04BA5C"

	//
	testClusterID = "6df10113-3c76-48ce-a742-0869fadd60b4"

	//
	testNodePoolId = "f532cde9-e574-40b6-856d-78fdfc8be3b9"

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

func getTestWorkloadClusterName() string {
	testCluster := os.Getenv("TCA_TEST_WORKLOAD_CLUSTER")
	if len(testCluster) == 0 {
		return testWorkloadClusterName
	}
	return testCluster
}

func getTestMgmtClusterName() string {
	e := os.Getenv("TCA_TEST_MGMT_CLUSTER")
	if len(e) == 0 {
		return testMgmtClusterName
	}
	return e
}

func getTestNodePoolName() string {
	testCluster := os.Getenv("TCA_TEST_NODE_POOL")
	if len(testCluster) == 0 {
		return testPoolName
	}
	return testCluster
}

func getTestRepoName() string {
	v := os.Getenv("TCA_TEST_REPO_NAME")
	if len(v) == 0 {
		return testRepoName
	}
	return v
}

func getTestInstanceName() string {
	v := os.Getenv("TCA_TEST_INSTANCE_NAME")
	if len(v) == 0 {
		return testInstanceName
	}
	return v
}

func getTestCatalogName() string {
	v := os.Getenv("TCA_TEST_CATALOG_NAME")
	if len(v) == 0 {
		return testCatalogName
	}
	return v
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
		return testWorkloadClusterName
	}
	return e
}

// getTestRepoUsername return test Repo either from env or default
func getTestRepoUsername() string {
	e := os.Getenv("TCA_TEST_REPO_USERNAME")
	if len(e) == 0 {
		return testRepoUsername
	}
	return e
}

// getTestRepoPassword return test Repo password from env or default
func getTestRepoPassword() string {
	e := os.Getenv("TCA_REPO_PASSWORD")
	if len(e) == 0 {
		return testRepoPassword
	}
	return e
}

//getTestRepoUrl return test Repo url from env or default
func getTestRepoUrl() string {
	e := os.Getenv("TCA_TEST_REPO_URL")
	if len(e) == 0 {
		return testRepoUrl
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
Description: ""
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
