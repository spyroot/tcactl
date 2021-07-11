package specs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//
func testFiles(dir string, fileType string) []string {

	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, fileType) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(file)
	}

	return files
}

// Read cluster spec from file.
func TestSpecCluster_SpecsFromFile(t *testing.T) {

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Read cluster workload spec from yaml",
			file:    "/cluster/positive/edge_workload_cluster.yaml",
			wantErr: false,
		},
		{
			name:    "Read cluster workload spec from json",
			file:    "/cluster/positive/edge_workload_cluster.json",
			wantErr: false,
		},
		{
			name:    "Read mgmt cluster spec from yaml",
			file:    "/cluster/positive/edge_mgmt_cluster.yaml",
			wantErr: false,
		},
		{
			name:    "Read mgmt cluster spec from json",
			file:    "/cluster/positive/edge_mgmt_cluster.json",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file
			spec, err := SpecCluster{}.SpecsFromFile(fileName)
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

			clusterSpec, ok := (*spec).(*SpecCluster)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = clusterSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() Test failed validator "+
					"return error for positive case err %v file %s", err, fileName)
				return
			}

		})
	}
}

func TestSpecCluster_SpecsFromString(t *testing.T) {

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Read cluster workload spec from yaml",
			file:    "/cluster/positive/edge_workload_cluster.yaml",
			wantErr: false,
		},
		{
			name:    "Read cluster workload spec from json",
			file:    "/cluster/positive/edge_workload_cluster.json",
			wantErr: false,
		},
		{
			name:    "Read mgmt cluster spec from yaml",
			file:    "/cluster/positive/edge_mgmt_cluster.yaml",
			wantErr: false,
		},
		{
			name:    "Read mgmt cluster spec from json",
			file:    "/cluster/positive/edge_mgmt_cluster.json",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file

			b, err := ioutil.ReadFile(fileName)
			assert.NoError(t, err)

			spec, err := SpecCluster{}.SpecsFromString(string(b))
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

			clusterSpec, ok := (*spec).(*SpecCluster)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = clusterSpec.Validate()
			if err != nil {
				t.Errorf("Test failed spec validator return error for positive case %v", err)
				return
			}

		})
	}
}

func TestSpecCluster_SpecsReader(t *testing.T) {

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Read cluster workload spec from yaml",
			file:    "/cluster/positive/edge_workload_cluster.yaml",
			wantErr: false,
		},
		{
			name:    "Read cluster workload spec from json",
			file:    "/cluster/positive/edge_workload_cluster.json",
			wantErr: false,
		},
		{
			name:    "Read mgmt cluster spec from yaml",
			file:    "/cluster/positive/edge_mgmt_cluster.yaml",
			wantErr: false,
		},
		{
			name:    "Read mgmt cluster spec from json",
			file:    "/cluster/positive/edge_mgmt_cluster.json",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file
			f, err := os.Open(fileName)
			assert.NoError(t, err)

			spec, err := SpecCluster{}.SpecsFromReader(f)
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

			clusterSpec, ok := (*spec).(*SpecCluster)
			if !ok {
				t.Errorf("Test failed method return wrong type")
				return
			}

			err = clusterSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() Test failed validator "+
					"return error for positive case err %v file %s", err, fileName)
				return
			}

		})
	}
}

func TestSpecCluster_Validate(t *testing.T) {

	tests := []struct {
		name       string
		testDir    string
		wantErr    bool
		wantValErr bool
	}{
		{
			name:       "Read invalid spec and do validation",
			testDir:    "/cluster/negative",
			wantErr:    false,
			wantValErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			testFileDir := assetsDir + tt.testDir

			files := testFiles(testFileDir, "yaml")
			for _, file := range files {
				spec, err := SpecCluster{}.SpecsFromFile(file)
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
				clusterSpec, ok := (*spec).(*SpecCluster)
				if !ok {
					t.Errorf("TestSpecCluster_Validate() failed method return wrong type")
					return
				}
				if clusterSpec == nil {
					t.Errorf("TestSpecCluster_Validate() failed SpecsFromFile return wrong type")
					return
				}
				err = clusterSpec.Validate()
				if tt.wantValErr && err == nil {
					t.Errorf("TestSpecCluster_Validate() "+
						"failed on file %s validator must return error", file)
					return
				}

				if clusterSpec.Kind() != "cluster" {
					t.Errorf("TestSpecCluster_Validate() wrong kind")
					return
				}

				t.Logf("validator return error %v\n", err)
			}
		})
	}
}

func TestSpecCluster_FindNodePoolByName(t *testing.T) {

	tests := []struct {
		name       string
		file       string
		masterPool string
		workerPool string
		wantErr    bool
		expect     bool
	}{
		{
			name:       "Read cluster workload spec from yaml",
			file:       "/cluster/positive/edge_workload_cluster.yaml",
			wantErr:    false,
			masterPool: "master",
			workerPool: "default-pool01",
			expect:     true,
		},
		{
			name:       "Read cluster workload spec from json",
			file:       "/cluster/positive/edge_workload_cluster.json",
			wantErr:    false,
			masterPool: "master",
			workerPool: "default-pool01",
			expect:     true,
		},
		{
			name:       "Read mgmt cluster spec from yaml",
			file:       "/cluster/positive/edge_mgmt_cluster.yaml",
			wantErr:    false,
			masterPool: "master",
			workerPool: "default-pool01",
			expect:     true,
		},
		{
			name:       "Read mgmt cluster spec from json",
			file:       "/cluster/positive/edge_mgmt_cluster.json",
			wantErr:    false,
			masterPool: "master",
			workerPool: "default-pool01",
			expect:     true,
		},
		{
			name:       "Read mgmt cluster spec from yaml",
			file:       "/cluster/positive/edge_mgmt_cluster.yaml",
			wantErr:    false,
			masterPool: "master1",
			workerPool: "default-pool",
			expect:     false,
		},
		{
			name:       "Read mgmt cluster spec from json",
			file:       "/cluster/positive/edge_mgmt_cluster.json",
			wantErr:    false,
			masterPool: "master1",
			workerPool: "default-pool",
			expect:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assetsDir := GetTestAssetsDir()
			fileName := assetsDir + tt.file

			b, err := ioutil.ReadFile(fileName)
			assert.NoError(t, err)

			spec, err := SpecCluster{}.SpecsFromString(string(b))
			if tt.wantErr && err == nil {
				t.Errorf("SpecsFromFile() failed must not return error")
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if spec == nil {
				t.Errorf("SpecsFromFile() return nil spec")
				return
			}
			clusterSpec, ok := (*spec).(*SpecCluster)
			if !ok {
				t.Errorf("SpecsFromFile() failed method return wrong type")
				return
			}
			err = clusterSpec.Validate()
			if err != nil {
				t.Errorf("SpecsFromFile() failed spec validator return error for positive case %v", err)
				return
			}

			if clusterSpec.FindNodePoolByName(tt.masterPool, false) != tt.expect {
				t.Errorf("SpecsFromFile() failed find a pool %v", tt.masterPool)
				return
			}

			if clusterSpec.FindNodePoolByName(tt.workerPool, true) != tt.expect {
				t.Errorf("SpecsFromFile() failed find a pool %v", tt.workerPool)
				return
			}
		})
	}
}
