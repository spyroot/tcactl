package api

import (
	"flag"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/testutil"
	"github.com/spyroot/tcactl/pkg/io"
	"os"
)

const (
	// test name, if you already have cluster adjust or pass TCA_TEST_CLUSTER env
	testWorkloadClusterName = "edge-test01"

	// testMgmtClusterName a cluster that must be defined int TCA
	testMgmtClusterName = "edge-mgmt-test01"

	// testCloudName a test cloud provider that must be defined in TCA
	testCloudName = "edge"

	// testPoolName a node pool that must be define and attached to test workload and mgmt cluster
	testPoolName = "default-pool01"

	// testCatalogName a test catalog entity
	testCatalogName = "unit_test"

	// testInstanceName a test instance
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

//
func getAuthenticatedClient() *client.RestClient {

	r := getClient()
	ok, err := r.GetAuthorization()
	if err != nil {
		io.CheckErr(err)
	}

	if !ok {
		io.PrintAndExit("failed authenticated")
	}

	r.SetTrace(false)

	return r
}

//harbor = &client.RestClient{
//	BaseURL:               os.Getenv("HARBOR_URL"),
//	apiKey:                "",
//	IsDebug:               true,
//	Username:              os.Getenv("HARBOR_USERNAME"),
//	Password:              os.Getenv("HARBOR_PASSWORD"),
//	SkipSsl:               true,
//}

//GetTestAssetsDir return dir where test assets are
//it call RunOnRootFolder that change dir
func GetTestAssetsDir() string {

	wd := testutil.RunOnRootFolder()
	wd = wd + "/test_assets"

	_, err := os.Stat(wd)
	if os.IsNotExist(err) {
		io.CheckErr(err)
	}

	return wd
}

// SetLoggingFlags Sets logging flag for logging tracer
func SetLoggingFlags() {

	err := flag.Set("alsologtostderr", "true")
	if err != nil {
		io.CheckErr(err)
	}

	err = flag.Set("v", "3")
	if err != nil {
		io.CheckErr(err)
	}

	flag.Parse()
	glog.Info("Logging configured")
}

// getClient() return tca client for unit testing
func getClient() *client.RestClient {

	tcaUrl := os.Getenv("TCA_URL")
	if len(tcaUrl) == 0 {
		io.PrintAndExit("TCA_URL not set")
	}

	tcaUsername := os.Getenv("TCA_USERNAME")
	if len(tcaUrl) == 0 {
		io.PrintAndExit("TCA_USERNAME not set")
	}

	tcaPassword := os.Getenv("TCA_PASSWORD")
	if len(tcaUrl) == 0 {
		io.PrintAndExit("TCA_PASSWORD not set")
	}

	r, err := client.NewRestClient(tcaUrl,
		true,
		tcaUsername,
		tcaPassword)

	if err != nil {
		io.CheckErr(err)
	}

	return r
}

// getTestClusterId() returns a
func getTestClusterId() string {
	e := os.Getenv("TCA_TEST_CLUSTER_ID")
	if len(e) == 0 {
		return testClusterID
	}
	return e
}

// return tenant id either default value or from env
func getTenantId() string {
	e := os.Getenv("TCA_TEST_TENANT_ID")
	if len(e) == 0 {
		return testTenantId
	}
	return e
}

func getTestWorkloadClusterName() string {
	e := os.Getenv("TCA_TEST_WORKLOAD_CLUSTER")
	if len(e) == 0 {
		return testWorkloadClusterName
	}
	return e
}

func getTestMgmtClusterName() string {
	e := os.Getenv("TCA_TEST_MGMT_CLUSTER")
	if len(e) == 0 {
		return testMgmtClusterName
	}
	return e
}

// return a default node pool name used for testing
func getTestNodePoolName() string {
	e := os.Getenv("TCA_TEST_NODE_POOL")
	if len(e) == 0 {
		return testPoolName
	}
	return e
}

// returns a default repo used for testing
func getTestRepoName() string {
	e := os.Getenv("TCA_TEST_REPO_NAME")
	if len(e) == 0 {
		return testRepoName
	}
	return e
}

// returns a default instance name used for testing
func getTestInstanceName() string {
	v := os.Getenv("TCA_TEST_INSTANCE_NAME")
	if len(v) == 0 {
		return testInstanceName
	}
	return v
}

// returns a default catalog name used for testing
func getTestCatalogName() string {
	e := os.Getenv("TCA_TEST_CATALOG_NAME")
	if len(e) == 0 {
		return testCatalogName
	}
	return e
}

// returns a default node pool id
func getTestNodePoolId() string {
	e := os.Getenv("TCA_TEST_NODE_POOL_ID")
	if len(e) == 0 {
		return testNodePoolId
	}
	return e
}

// returns a default mgmt cluster template id
func getTestMgmtTemplateId() string {
	e := os.Getenv("TCA_TEST_MGMT_TEMPLATE_ID")
	if len(e) == 0 {
		return testMgmtTemplateId
	}
	return e
}

// returns a default workload cluster template id
func getTestWorkloadTemplateId() string {
	e := os.Getenv("TCA_TEST_WORKLOAD_TEMPLATE_ID")
	if len(e) == 0 {
		return testWorkloadTemplateId
	}
	return e
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
kind: template
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
