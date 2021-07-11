package response

import (
	"encoding/json"
	"github.com/nsf/jsondiff"
	"github.com/spyroot/tcactl/lib/testutil"
	ioutils "github.com/spyroot/tcactl/pkg/io"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// specNodePoolStringReaderHelper test helper return node spec
func clusterSpecFromFile(spec string) (*ClusterSpec, error) {

	tmpFile, err := ioutil.TempFile("", "tcactltest")
	ioutils.CheckErr(err)

	defer os.Remove(tmpFile.Name())

	// write to file,close it and read spec
	if _, err = tmpFile.Write([]byte(spec)); err != nil {
		ioutils.CheckErr(err)
	}

	if err := tmpFile.Close(); err != nil {
		ioutils.CheckErr(err)
	}

	// read from file
	r, err := ClusterSpecsFromFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Test cluster spec from yaml or json string
// It must create new instance of cluster
func TestClusterInstanceSpecsFromString(t *testing.T) {

	tests := []struct {
		name    string
		spec    string
		wantErr bool
	}{
		{
			name:    "Basic Test",
			spec:    testJsonCluster,
			wantErr: false,
		},
		{
			name:    "Basic broken spec",
			spec:    YamlBrokenSpec,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClusterSpec{}.InstanceSpecsFromString(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpecsFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			newJson, err := json.Marshal(got)
			assert.NoError(t, err)

			assert.NotNil(t, got)
			opt := jsondiff.DefaultJSONOptions()
			diff, _ := jsondiff.Compare(newJson, []byte(tt.spec), &opt)

			if tt.wantErr != true && diff > 0 {
				t.Errorf("SpecsFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewClusterSpecs(t *testing.T) {

	tests := []struct {
		name        string
		reader      io.Reader
		wantErr     bool
		closeReader bool
	}{
		{
			name:        "Basic spec from json",
			reader:      strings.NewReader(testJsonCluster),
			wantErr:     false,
			closeReader: false,
		},
		{
			name:        "Basic spec from yaml",
			reader:      strings.NewReader(YamlPartialSpec),
			wantErr:     false,
			closeReader: false,
		},
		{
			name:        "Basic broken yaml spec",
			reader:      strings.NewReader(YamlBrokenSpec),
			wantErr:     true,
			closeReader: false,
		},
		{
			name:        "From a file ",
			reader:      testutil.SpecTempReader(testJsonCluster),
			wantErr:     false,
			closeReader: false,
		},
		{
			name:        "Closed file reader",
			reader:      strings.NewReader(testJsonCluster),
			wantErr:     true,
			closeReader: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := tt.reader
			if tt.closeReader {
				tmpFile, _ := ioutil.TempFile("", "tcactltest")
				tmpFile.Close()
				r = tmpFile
			}

			got, err := NewClusterSpecs(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClusterSpecs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				return
			}

			assert.NotNil(t, got)
		})
	}
}

func TestNilClusterSpecs(t *testing.T) {

	tests := []struct {
		name        string
		cluster     *Clusters
		wantErr     bool
		closeReader bool
	}{
		{
			name: "Basic nil test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := tt.cluster.GetClusterIds()
			t.Log(err)
			assert.Error(t, err)

			_, err = tt.cluster.GetClusterSpec("")
			t.Log(err)
			assert.Error(t, err)

			_, err = tt.cluster.GetClusterId("")
			t.Log(err)
			assert.Error(t, err)
		})
	}
}

// Test cluster spec from file.
//
func TestClusterSpecsFromFile(t *testing.T) {

	tests := []struct {
		name    string
		spec    string
		wantErr bool
	}{
		{
			name:    "Basic spec from json string",
			spec:    testJsonCluster,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := clusterSpecFromFile(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClusterSpecsFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

// Test Clusters list from yaml or json string
// it return Clusters object that encapsulate entire
// list
func TestClustersInstanceSpecsFromString(t *testing.T) {

	tests := []struct {
		name    string
		spec    string
		wantErr bool
	}{
		{
			name:    "Basic test clusters from string",
			spec:    testClusters,
			wantErr: false,
		},
		{
			name:    "Basic broken yaml spec. clusters from string",
			spec:    YamlBrokenSpec,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := Clusters{}.InstanceSpecsFromString(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpecsFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			newJson, err := json.Marshal(got)
			assert.NoError(t, err)

			assert.NotNil(t, got)
			opt := jsondiff.DefaultJSONOptions()
			diff, _ := jsondiff.Compare(newJson, []byte(tt.spec), &opt)

			if tt.wantErr != true && diff > 0 {
				t.Errorf("TestClustersInstanceSpecsFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// Test NewClustersSpec from io reader
//
func TestNewClustersSpecs(t *testing.T) {

	tests := []struct {
		name    string
		reader  io.Reader
		wantErr bool
	}{
		{
			name:    "Basic spec from json",
			reader:  strings.NewReader(testClusters),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClustersSpecs(tt.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClusterSpecs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

// Test GetClusterSpec
func TestGetClusterId(t *testing.T) {

	tests := []struct {
		name        string
		spec        string
		clusterName string
		wantErr     bool
		expect      string
	}{
		{
			name:        "Basic test get cluster from name",
			spec:        testClusters,
			wantErr:     false,
			clusterName: "edge-mgmt-test01",
			expect:      "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
		},
		{
			name:        "Basic test get cluster from id",
			spec:        testClusters,
			wantErr:     false,
			clusterName: "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
			expect:      "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
		},
		{
			name:        "Basic test get cluster not found",
			spec:        testClusters,
			wantErr:     true,
			clusterName: "edge-mgmt-test0",
			expect:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Clusters{}.InstanceSpecsFromString(tt.spec)
			assert.NoError(t, err)

			clusters, ok := got.(*Clusters)
			assert.Equal(t, true, ok)

			id, err := clusters.GetClusterId(tt.clusterName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, id, tt.expect)
				return
			}

		})
	}
}

// Test GetClusterSpec
func TestGetClusterSpec(t *testing.T) {

	tests := []struct {
		name        string
		spec        string
		clusterName string
		wantErr     bool
		expect      string
	}{
		{
			name:        "Basic test get, cluster from name",
			spec:        testClusters,
			wantErr:     false,
			clusterName: "edge-mgmt-test01",
			expect:      "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
		},
		{
			name:        "Basic test, get cluster from id",
			spec:        testClusters,
			wantErr:     false,
			clusterName: "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
			expect:      "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
		},
		{
			name:        "Basic test get, cluster not found",
			spec:        testClusters,
			wantErr:     true,
			clusterName: "edge-mgmt-test0",
			expect:      "",
		},
		{
			name:        "Basic test get, cluster on empty",
			spec:        emptyClustersSpec,
			wantErr:     true,
			clusterName: "edge-mgmt-test0",
			expect:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Clusters{}.InstanceSpecsFromString(tt.spec)
			assert.NoError(t, err)

			clusters, ok := got.(*Clusters)
			assert.Equal(t, true, ok)

			clusterSpec, err := clusters.GetClusterSpec(tt.clusterName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, tt.expect, clusterSpec.Id)
				return
			}
		})
	}
}

// Test TestGetClusterIds
func TestGetClusterIds(t *testing.T) {

	tests := []struct {
		name      string
		spec      string
		wantErr   bool
		expectLen int
	}{
		{
			name:      "Basic positive test get clusters id",
			spec:      testClusters,
			wantErr:   false,
			expectLen: 2,
		},
		{
			name:      "Basic negative test get clusters id",
			spec:      emptyClustersSpec,
			wantErr:   false,
			expectLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Clusters{}.InstanceSpecsFromString(tt.spec)
			assert.NoError(t, err)

			clusters, ok := got.(*Clusters)
			assert.Equal(t, true, ok)

			ids, err := clusters.GetClusterIds()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, tt.expectLen, len(ids))
				return
			}
		})
	}
}

// Test GetClusterSpec
func TestGetFields(t *testing.T) {

	tests := []struct {
		name        string
		spec        string
		clusterName string
		wantErr     bool
		expect      string
	}{
		{
			name:        "Basic test get, cluster from name",
			spec:        testClusters,
			wantErr:     false,
			clusterName: "edge-mgmt-test01",
			expect:      "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
		},
		{
			name:        "Basic test get, cluster empty",
			spec:        emptyClustersSpec,
			wantErr:     true,
			clusterName: "edge-mgmt-test01",
			expect:      "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Clusters{}.InstanceSpecsFromString(tt.spec)
			assert.NoError(t, err)

			clusters, ok := got.(*Clusters)
			assert.Equal(t, true, ok)
			if !ok {
				return
			}

			clusterSpec, err := clusters.GetClusterSpec(tt.clusterName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fields, err := clusterSpec.GetFields()
			assert.NoError(t, err)

			if !tt.wantErr {
				v, ok := fields["id"]
				assert.Equal(t, true, ok)
				assert.Equal(t, tt.expect, v)
				t.Logf("key id points to %s", v)
				return
			}
		})
	}
}

// Test GetClusterSpec
func TestGetField(t *testing.T) {

	tests := []struct {
		name          string
		spec          string
		clusterName   string
		wantErr       bool
		negativeField string
	}{
		{
			name:        "Basic positive test get struct field",
			spec:        testClusters,
			wantErr:     false,
			clusterName: "edge-mgmt-test01",
		},
		{
			name:          "Basic negative case",
			spec:          testClusters,
			wantErr:       false,
			clusterName:   "edge-mgmt-test01",
			negativeField: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Clusters{}.InstanceSpecsFromString(tt.spec)
			assert.NoError(t, err)

			clusters, ok := got.(*Clusters)
			assert.Equal(t, true, ok)
			if !ok {
				return
			}

			clusterSpec, err := clusters.GetClusterSpec(tt.clusterName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fields, err := clusterSpec.GetFields()
			assert.NoError(t, err)

			for fieldName, _ := range fields {
				val := clusterSpec.GetField(fieldName)
				t.Logf("Filed %s val %s", fieldName, val)
			}

			// shouldn't crash
			if len(tt.negativeField) > 0 {
				val := clusterSpec.GetField(tt.negativeField)
				t.Logf("Filed %s val %v", tt.negativeField, val)
			}
		})
	}
}

var testClusters = `
{
  "Clusters": [
    {
      "id": "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
      "clusterName": "edge-mgmt-test01",
      "clusterType": "MANAGEMENT",
      "vsphereClusterName": "hubsite",
      "managementClusterId": "",
      "hcxUUID": "20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6",
      "status": "ACTIVE",
      "activeTasksCount": 0,
      "clusterTemplate": {
        "name": "min",
        "version": "",
        "id": "55e69a3c-d92b-40ca-be51-9c6585b89ad7"
      },
      "clusterId": "",
      "clusterUrl": "https://10.241.7.224:6443",
      "kubeConfig": "YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VONWVrTkRRV0pQWjBGM1NVSkJaMGxDUVVSQlRrSm5hM0ZvYTJsSE9YY3dRa0ZSYzBaQlJFRldUVkpOZDBWUldVUldVVkZFUlhkd2NtUlhTbXdLWTIwMWJHUkhWbnBOUWpSWVJGUkplRTFFV1hoUFJFVXpUWHBWZVUweGIxaEVWRTE0VFVSWmVFNXFSVE5PUkVGNVRURnZkMFpVUlZSTlFrVkhRVEZWUlFwQmVFMUxZVE5XYVZwWVNuVmFXRkpzWTNwRFEwRlRTWGRFVVZsS1MyOWFTV2gyWTA1QlVVVkNRbEZCUkdkblJWQkJSRU5EUVZGdlEyZG5SVUpCVDFOUUNtUnRUWFJ6Ym5WRFIxQm5NRmR1ZUV4YVYwUm5Sek56YlhZeGJFMVdTM0puVldFMlYwcGtiRGhIY1hOclNXNXRlR0ZxZFZCRmMwY3ZhelZFVUhKUFJrUUthWGgyY0hWRk1YSXlaVXB5TUVRM2VEbGxjVEI1YURGaWQzcHZOMms0WWpCU1ppOUZTU3RaTW1vcmFGRTNaV3BJZVVKamExbDVMMjkwT1hCVVpGVTJkd3BZY0RGclpsZDZSRVJvWjBsdVkwZFFaelp1UzFkbEwwSlliRmhsTVhaSE5EbHBMMWxFU0dSeVZFTjRVVWcxWldsMlZrRmhWekZ0TkVsRlZHbzNXa3RTQ200MllUZHVUa3MxVTBjeFEydG1TR0p0ZVRNdk5WZFBORXhVVlV4NmRUZzRabmxCTDJSVk5YRlNNM3BsVWtaTVIyUmxSMUFyVEVsaFdqbFhSV0ptV0ZFS1ZuUldkMVppZDI1elZsaDVOWEJ6ZW1wYUt6SnlVa0Z2WlZNNWVGRjBlamhyUVVsTVp5OUlXVzB4ZDJoVVFrcGpjemR3S3pSRVRVeHJUemsyVEZVME1RcFZSVzlUTUVwcmIxRnJZVmRRVm1SS1ZFTTRRMEYzUlVGQllVMXRUVU5SZDBSbldVUldVakJRUVZGSUwwSkJVVVJCWjB0clRVSkpSMEV4VldSRmQwVkNDaTkzVVVsTlFWbENRV1k0UTBGUlFYZEVVVmxLUzI5YVNXaDJZMDVCVVVWTVFsRkJSR2RuUlVKQlRtdGFTMDl0TWxsVU9VZDRjRFY1ZDNWUVJYSk1OVmdLTVd0YVV6QmpZa3BMY2tsbWFVdENlV3BwUkhCWmNqVkROMDl5VkhSSVdURXJSMFZIU0Rrck5ITkRZVWRqVldaaFEyTnBWaTlVVTBoMFJtaHZNVTVDYlFwMWVVWlZjRlZNU0ZGR1VIYzJielpwVmpoR1VYTlFhbVkyUnpSaVYyaHZXUzlZVWpKM1NpdFlWRWRHVmpWeE1qbElNVTFaTldSTlVreFJNbTlWTmpaUENrZHJVMll5V2xsYU1uaDZSRWt6YjFaTllsQjBObFJOVFdscWMwbGlSREV3ZEVzMlYwZDRXSFYzUWtvNGNHcE1jMmN2UkdaUlJrNW9aMjh5Um5OM1VEQUthMk5hYzAxUFoydENkMU5pUlM5U01tdEthSEJYUXpGcVNWQktSVU0wV21oYWF5OVhhMHN3Y3psUldFNDFSVWN4TUhSbkt6TkxZVEpPVkZocE0yZzRkQXBYTTFGRk5tWkxTbWRCZW5wM1Z6VjVMMmxVY2l0YVpXSk5RblJDV21ZdmVIaFhka1ZrTXpKVmVIZGtTR3RIVjBwcE1XWjZiMk01Tnk5RmJHODVObGs5Q2kwdExTMHRSVTVFSUVORlVsUkpSa2xEUVZSRkxTMHRMUzBLCiAgICBzZXJ2ZXI6IGh0dHBzOi8vMTAuMjQxLjcuMjI0OjY0NDMKICBuYW1lOiBlZGdlLW1nbXQtdGVzdDAxCmNvbnRleHRzOgotIGNvbnRleHQ6CiAgICBjbHVzdGVyOiBlZGdlLW1nbXQtdGVzdDAxCiAgICB1c2VyOiBlZGdlLW1nbXQtdGVzdDAxLWFkbWluCiAgbmFtZTogZWRnZS1tZ210LXRlc3QwMS1hZG1pbkBlZGdlLW1nbXQtdGVzdDAxCmN1cnJlbnQtY29udGV4dDogZWRnZS1tZ210LXRlc3QwMS1hZG1pbkBlZGdlLW1nbXQtdGVzdDAxCmtpbmQ6IENvbmZpZwpwcmVmZXJlbmNlczoge30KdXNlcnM6Ci0gbmFtZTogZWRnZS1tZ210LXRlc3QwMS1hZG1pbgogIHVzZXI6CiAgICBjbGllbnQtY2VydGlmaWNhdGUtZGF0YTogTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVTTRha05EUVdSeFowRjNTVUpCWjBsSlMzQkpVbkJqYWxkcGNtTjNSRkZaU2t0dldrbG9kbU5PUVZGRlRFSlJRWGRHVkVWVVRVSkZSMEV4VlVVS1FYaE5TMkV6Vm1sYVdFcDFXbGhTYkdONlFXVkdkekI1VFZSQk1rMVVaM2hPZWsweFRXcE9ZVVozTUhsTmFrRXlUVlJuZUU1NlVYZE5hazVoVFVSUmVBcEdla0ZXUW1kT1ZrSkJiMVJFYms0MVl6TlNiR0pVY0hSWldFNHdXbGhLZWsxU2EzZEdkMWxFVmxGUlJFVjRRbkprVjBwc1kyMDFiR1JIVm5wTVYwWnJDbUpYYkhWTlNVbENTV3BCVGtKbmEzRm9hMmxIT1hjd1FrRlJSVVpCUVU5RFFWRTRRVTFKU1VKRFowdERRVkZGUVhsSU1rSnhXV2M1ZFVaUGIwazRTMU1LU2xKNVVWTjFla0p1YmxSUlNVaGxWSGREVmpGSlJEWnBkbVZaTW1kMWRFZEVZMGR2T1dacFYyVXlZWHBoVEhCVFkzY3ZVbXhIZG14NmJtSllVSFUxY2dvdmExTmxhMlJvU1VKbmJ5OHhRMUZCVVdwNVRYVkVkbUp2Vnk5c1V5ODFkekVyY201M2RpdFlOWGRMZVdOdlNEY3plRGt2UjBwcllXOHlPRkI0T1VOS0NuUm9abEZtZEZoaWNuTXplRkJIZDJvMmN5OUlaeXRpTXpaeWRtMUdOR1psTDFSa1FtdGpibUpIV2xaU1NWZEhWVzgwY0ZWRk1XaHlaMmhpTlZOd09VTUtLM0ZqVVVaUEszRmFVVlZYWWpsRk1GTkhOak4zVDNkVmFEaEJOVXRvYXlzeVMyMDJWbXd4Wm14RllWcENlVFphVEc5WFpWWTFSWE42TTBKU2VpOVdRZ3BpVEVwalNYVm1ZVGR4ZFZONk5rWmpVV3d3T0hFNGNDdFhibWRoV0N0b1ZrdHlkM0ZhT0RscmVFVkVTVFoyZFRSUlpGaHVMMHRNT0dWNlVtaEdTRVF6Q21ST1RHNW5VVWxFUVZGQlFtOTVZM2RLVkVGUFFtZE9Wa2hST0VKQlpqaEZRa0ZOUTBKaFFYZEZkMWxFVmxJd2JFSkJkM2REWjFsSlMzZFpRa0pSVlVnS1FYZEpkMFJSV1VwTGIxcEphSFpqVGtGUlJVeENVVUZFWjJkRlFrRk9SMHB6VVZoQk4waEJNVlF2T0VGVVdrdEtSMU5STDJaTVNFZHhXRWxCYUdkT2NRcExkVWxTZEUxR1VYWlpNMFJZWW5odU5FOTZhVlkzUjJkblZYRjZWbFV4UzFJeE0wNUJaMlpOT1c5dWIyRnplbFZtVm5WeU4ycEhTa0UzU2xkS2VubHVDakJsYm5jd2RYazFiSGxFTDNCemIyZGlVQ3RxTldsVGFWTlVlR2QxV2pCd01WQkRLek13WWtGc2JFaEZaRVJMVEdsQ2FUQTJlRmh5YWt0WGFrbFlabVVLVjNwaWQzaHFlbXAxT0ROMVowZGtSRXhxUjAwd05tUkVSSFZ4WVVFNU5GbzNOalZZZGxSTlN6Z3haa2RZUW10MlRFUjNWMlJyU0ZWWFkydDNSMVZ4YndweU1uWlVNREpzUjNSRGQwTklabXhJWjIxMVUzZDJkaTlwTUZaUWRrbHhWa1JpUW5OMU1DOWtkR2RqY1VWTU56TkhWR0ZGUjJab1FYSlNXSE5PZUhac0NucG1XakZrUW5GalRUbEVRbE5hZWtKSGIxY3hOa3Q2VWpWT2VXTnZla2ROVW5vMmNqQkhXVVpDYTBabVFYaG9SVTlJWnowS0xTMHRMUzFGVGtRZ1EwVlNWRWxHU1VOQlZFVXRMUzB0TFFvPQogICAgY2xpZW50LWtleS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJTVTBFZ1VGSkpWa0ZVUlNCTFJWa3RMUzB0TFFwTlNVbEZjRUZKUWtGQlMwTkJVVVZCZVVneVFuRlpaemwxUms5dlNUaExVMHBTZVZGVGRYcENibTVVVVVsSVpWUjNRMVl4U1VRMmFYWmxXVEpuZFhSSENrUmpSMjg1Wm1sWFpUSmhlbUZNY0ZOamR5OVNiRWQyYkhwdVlsaFFkVFZ5TDJ0VFpXdGthRWxDWjI4dk1VTlJRVkZxZVUxMVJIWmliMWN2YkZNdk5YY0tNU3R5Ym5kMksxZzFkMHQ1WTI5SU56TjRPUzlIU210aGJ6STRVSGc1UTBwMGFHWlJablJZWW5Kek0zaFFSM2RxTm5NdlNHY3JZak0yY25adFJqUm1aUW92VkdSQ2EyTnVZa2RhVmxKSlYwZFZielJ3VlVVeGFISm5hR0kxVTNBNVF5dHhZMUZHVHl0eFdsRlZWMkk1UlRCVFJ6WXpkMDkzVldnNFFUVkxhR3NyQ2pKTGJUWldiREZtYkVWaFdrSjVObHBNYjFkbFZqVkZjM296UWxKNkwxWkNZa3hLWTBsMVptRTNjWFZUZWpaR1kxRnNNRGh4T0hBclYyNW5ZVmdyYUZZS1MzSjNjVm80T1d0NFJVUkpObloxTkZGa1dHNHZTMHc0WlhwU2FFWklSRE5rVGt4dVoxRkpSRUZSUVVKQmIwbENRVWhUTlU1VFRVcGpZakU1TTNaamNRb3dWVVZRT1hsWVdETjNWM040Ym5OUlozaEhjaTlTZGxCMU5XNUthazVVT1hWMVVrcGhiVUZDTkZWS2J6QlhlbGxvVTFnclFVdDBhbWRuUTBwRFpHNVhDbTR3VFVOdGFuTlpObmd5ZUVjcldXSkxOV1JyWm5oNlNIQktMMEpqU1V0ck1sQmxjME5HUW5Jd0wydHVhR1JNUzNkSFVIWm9NbEV2ZUhSV09URklZMndLWm5aa2JuUXZaMUE0THpBMldYZGpiR2R0WlhwdWJHMDFkak0yUmpkSU5XaHdjME5FTldSVE1sWk9lbW8yUjNrMmRrMVVZMk56TkhwSE9USTFkRzh5VlFvMVFrSkxhMHhXTVhGdFZsRkVjamhGYWxCbmMwdFZLMWxhZFRCeU1XSmpkbVoxY1RWM1MydEdTa2MzUkdoTldXUjNkeTl1UnlzMFNWaENTVXRMTUhJMkNteEJORGhzVVhkclowaG1jM1Z6VEdWVlIwczBWVUZsVjNaaGVGaEVNV3BUU0RWcU5WcEdhVTQxWlUxSFZFaHFZVFJuUVhWSFMxQTBiSE55V1hrek9VZ0tLMWRIVEZSYVJVTm5XVVZCTm5nek1WaE1TM0JNWm01SGMzaHRRMGRMVUVkQk9YWktiRVpYUms4cllqRXdRbFJpTkhoT2NVcHNRakJRVGtOSVVWaGpPUXBWVEhOVFkwOHpZVmxGV1RKNGVsQnZTV1JEVlVNek1HbEJaMVZyUW1aR1JUSkNNbFpNYW1oVVNWZEZNbmRYY2xwc01uRTRha1kzU0hvd0swcGxTVTB2Q2l0NVdFVTRLMFprVVdSalRGVTVNa2dyVjJSWmJYTnRWVmxNYW0xQlRVMDBUVWxDWWtWNVVGTkhaWEpKTWprMk5FdHVVVzVKY0ZWRFoxbEZRVEpyZHpJS1VVeFplV2w1WTAwM01EazFMMjlzYUhWRldta3hjbGhqWjNZMFVWSllabXhRZUZsRk1XUnhSelYyWmxCRlEyeFlSa3gzU0RKMFUxZFpWMlJSUm1rd1N3cFdja3Q1UTBGTmJIUm5SVXBGV0dVM01FTXlZbXB5Y2s5aFYyRmFOMmhSVHl0RVRqaFBhUzlpUVRaUmF6RkVSMkZPVUdGUWJGaFliUzlDYUdSR1dWWlhDblpoTHpORVlVVk1MeXRMVUZST0wxcGhVelJ1UXpOTE9HbFhiVlUxVkdOWlIzSjNkMmRxTUVObldVSldhV0pIZUdoMVNtSTRjWEJwYTBZMGRsaHhNekFLZVRKd09FSTVNRkZJVVdaRVN6QlBVVk52YW5NMVdYQlZZMkl3ZW1sU1JVdHdVRk5vYkhwa1EyNVhhbWxqYWxRMk1FZHdOeTgwZERJeGNUTjZXalJMV0FwaFlsQlVaUzlTY1dsWVlXeFRXakZsU2tKMmVYRmhlbWc1Um5KalQwVjJNM05DWTBkaWVrSkNaWFEwU2twbGFVcFRVbk5KVUZaM05VVkRka1pxVlVOT0NrUk1WMVp6WjNWSVIwVTNOVUoyT1ZSdGNYZ3JVRkZMUW1kUlExSXZXV2RITjFKTFQxWXhXbWxOVkM5QlFqZGhkbTB6Y2poYWRrZE1UbFZvUTJaVVZqSUtXa1UxYkdsMGFsQllTV3hDU1V4VFdreEZjMEZuVlRNd05VUm1LMHhMT0U5WWEwZE1PVVJpYWt0dWQySkVNekIxWXpKSWJVdFhOV3AwVUdSaE4yeFZOQXBxYnpSV1FpczBjakpFYlRsSk1VMXBLMjlZTmtOaWFYQm1ZekkyTVRoeGVqTjRkVlpLVGtFeWJVSnJTRFl5YVd3eUwxUkhlR3AyY1d4bk5qWTBTV0ZvQ2xWd2QwdGlVVXRDWjFGRFYwdzJiMGxDVDNBMWRWQkNRVXBYV0c0eGRubDBjVGRDYVZKVlR6SmhSbGhqUlhJNGIyWjVjbmRoUjJGc1ptTk5kakp5VHpZS1FrTkxVRk54TjFVck5HTjJVeTlQTjJRNVNIcDJkbGhNVlM5TGRYZEViVU5ZU1U1d0sxTnliazVzUzFkNWQzaHhSWG80V0haSVRWRXpPRFpaYmpocGVRbzBOM1JFWldsQ2VGVjJjVmw1TVdkRmNHZ3lSVWt2V0hOVk5FOXBVVEU1UWxWc1ZXUjFOM0YyWW5CWWVrdDBUV2hYVVhWdVprRTlQUW90TFMwdExVVk9SQ0JTVTBFZ1VGSkpWa0ZVUlNCTFJWa3RMUzB0TFFvPQo=",
      "endpointIP": "",
      "masterNodes": [
        {
          "cpu": 4,
          "memory": 16384,
          "name": "master",
          "networks": [
            {
              "label": "MANAGEMENT",
              "networkName": "/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0",
              "nameservers": [
                "10.246.2.9"
              ]
            }
          ],
          "storage": 50,
          "replica": 1,
          "labels": [],
          "cloneMode": "linkedClone"
        }
      ],
      "workerNodes": null,
      "vimId": "",
      "error": ""
    },
    {
      "id": "868636c9-868f-49fb-a6df-6a0d2d137146",
      "clusterName": "edge-test01",
      "clusterType": "WORKLOAD",
      "vsphereClusterName": "hubsite",
      "managementClusterId": "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
      "hcxUUID": "20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6",
      "status": "ACTIVE",
      "activeTasksCount": 0,
      "clusterTemplate": {
        "name": "myworkload",
        "version": "",
        "id": "c3e006c1-e6aa-4591-950b-6f3bedd944d3"
      },
      "clusterId": "",
      "clusterUrl": "https://10.241.7.228:6443",
      "kubeConfig": "YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VONWVrTkRRV0pQWjBGM1NVSkJaMGxDUVVSQlRrSm5hM0ZvYTJsSE9YY3dRa0ZSYzBaQlJFRldUVkpOZDBWUldVUldVVkZFUlhkd2NtUlhTbXdLWTIwMWJHUkhWbnBOUWpSWVJGUkplRTFFV1hsTlZFVjZUbFJOZWsweGIxaEVWRTE0VFVSWmVFOVVSWHBPVkdkNlRURnZkMFpVUlZSTlFrVkhRVEZWUlFwQmVFMUxZVE5XYVZwWVNuVmFXRkpzWTNwRFEwRlRTWGRFVVZsS1MyOWFTV2gyWTA1QlVVVkNRbEZCUkdkblJWQkJSRU5EUVZGdlEyZG5SVUpCVERsWkNqUjBWamRVVjBGaVlXeEZURTg0YkN0d09WcElVMGhHVTJOUVJWcGxhRzFhVEZnNFZIaHFTa3N5WXpWWlVHNTBNVVU1UWsxSE4yazRVREEwV2tKWVN6UUtRazlKSzFKUlIwNUtWWEl2UVhKNE0wWkxOMDh4V1ZwaU9XdzBXREF3V2psM0wzWlFhVE5RTUVoWk5uTkNhV2xyUlROd01XRnBiRkV5UjA4MlJ6WjJWd3AxWVVOTmVGZ3lVWHBwWTJsTlNEazNkR3QyZDI1bUt6QjNWRlZLVjJkVFlrdDZkblJoYVdGaldYaHlXa3AxT1VNeWJEbFdNVU0zYUUxSmNVUlRhRUZGQ2tGWWF6TnFZV3RIV0VKaFoycHBkbXBMYmxBdmJ6bDJjR3BxTlZOeFVtVTRUa0pzVGs5emVqSnJaVUpXZEhGdFJuZ3hlbEpLZEZONU1tZEdVbUUxYWpnS04ySk1jSEpRWVZWemF5czRhM0JFVUc5aFdWQkllRUZqYlhkUmJIWXpXbGxYVFVkeVZqSlJLMVJVTUc5QlJtSmhaWGxtWlhOdFJEWmxaVEU1VjNneGN3cE1ObXRLYlRKQ1RVOWtjMUJhWkdaNGRDOXJRMEYzUlVGQllVMXRUVU5SZDBSbldVUldVakJRUVZGSUwwSkJVVVJCWjB0clRVSkpSMEV4VldSRmQwVkNDaTkzVVVsTlFWbENRV1k0UTBGUlFYZEVVVmxLUzI5YVNXaDJZMDVCVVVWTVFsRkJSR2RuUlVKQlMxVnRWRVpQT0Vwc1JUZEtlQ3RNZFdsUWVFOHdTSG9LVUU5eVdHZDBlbEIxVTFSdmJGbDJTVTVqUTJsQ0swbHlTVlowTjBReWJHUmpWMkpSY1ZwcVpXRjBObmg2VDBkMVoxYzRXWEYzU0d4Mk5uZHVTMWRPU2dwRGVVMDFhRGxOVDNWclJESjVkbmMxTlZCRk1XMTZaRVpyUVZSamJVRlZNQ3RzZFROUlNEQlpZa1psU1VkV1dsVkJPV04zVXlzMVlVbzFkWGMxYnlzMUNpdDJZM0kwYlRsNVJVVjZlWGhaZG1wNGNGTm5RMGx4ZEc1SFdpdHZPSFpNUW10UFFpdGlVRGhvVDFSalFuUk9VM1pETm1rd1ZWTkZXRTVMZHl0V2Vtb0tMMVpQU0VwMFVISk5kVW95VVc4M1QxaElVa0Z3U0cxRlVuTjZhazU2ZWxOaWFIaHBSV2xDVUZneFYyWjRUMHd4V2prMmNGVmFiMlJNUWxwbVJtTmtTZ3BhY0dZMFp6UmFjVFExYjJZck0xRjZjM2hoTnk5WlowbzFMMXBqUldwVWIwWkpZak5DUTFJelZ6SlNNalZ3VXpOV2JVNHlNR1ZvVkhSMWFrSnZXakE5Q2kwdExTMHRSVTVFSUVORlVsUkpSa2xEUVZSRkxTMHRMUzBLCiAgICBzZXJ2ZXI6IGh0dHBzOi8vMTAuMjQxLjcuMjI4OjY0NDMKICBuYW1lOiBlZGdlLXRlc3QwMQpjb250ZXh0czoKLSBjb250ZXh0OgogICAgY2x1c3RlcjogZWRnZS10ZXN0MDEKICAgIHVzZXI6IGVkZ2UtdGVzdDAxLWFkbWluCiAgbmFtZTogZWRnZS10ZXN0MDEtYWRtaW5AZWRnZS10ZXN0MDEKY3VycmVudC1jb250ZXh0OiBlZGdlLXRlc3QwMS1hZG1pbkBlZGdlLXRlc3QwMQpraW5kOiBDb25maWcKcHJlZmVyZW5jZXM6IHt9CnVzZXJzOgotIG5hbWU6IGVkZ2UtdGVzdDAxLWFkbWluCiAgdXNlcjoKICAgIGNsaWVudC1jZXJ0aWZpY2F0ZS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VNNGFrTkRRV1J4WjBGM1NVSkJaMGxKUkRoa1NqZEJVVVJpWWsxM1JGRlpTa3R2V2tsb2RtTk9RVkZGVEVKUlFYZEdWRVZVVFVKRlIwRXhWVVVLUVhoTlMyRXpWbWxhV0VwMVdsaFNiR042UVdWR2R6QjVUVlJCTWsxcVJYaE5lbFY2VFhwT1lVWjNNSGxOYWtFeVRXcEZlRTE2VlRSTmVsSmhUVVJSZUFwR2VrRldRbWRPVmtKQmIxUkViazQxWXpOU2JHSlVjSFJaV0U0d1dsaEtlazFTYTNkR2QxbEVWbEZSUkVWNFFuSmtWMHBzWTIwMWJHUkhWbnBNVjBackNtSlhiSFZOU1VsQ1NXcEJUa0puYTNGb2EybEhPWGN3UWtGUlJVWkJRVTlEUVZFNFFVMUpTVUpEWjB0RFFWRkZRVzVWWm1sTVlWQTFVV3AzYXpSb1FUWUtWR05HYlZJd01YaHpSVGRoT0V0M2JXRkVZM0UxSzJSYVoxaFNkeXREVGxOeFJIb3hOREY0ZEN0MmFqRlZjMnh0ZEZSaU1uSlFXRzVuVW5Wa1ZqRm9VZ3B6VWpZM1NXeHViVU5UT0RJd2FrTmFVM2gzTWt0RE5UaEZXVkpDWVhSMksza3dRV3BUU1dkc1RrY3ZjV00wYUVkYVNqTkVaMDg0TjFoVU0ycHJaMWhtQ20xdk5FNHlXVkpKYUhKR1QyWnRRbVJyYWpGQ01tMTBZWGxxVkdkMFVXODNibVJYYUdwUVNXMW1iVGhOYjFSSFlUVk1RVVZ1ZEZOSmRpOU1RVWhuYVhVS1UwSnhPVmt5UTI0MkwydEhTMHQwTmpSaVVYWkJOMmNyYlZvMFJEa3JRMjlPY2t4bmNVaG9abmxHTkd0cFpEQjJSVUp3YlV4Q2VtUTRjV0l6SzJock9RcFNlbmhzY1VGakswMHdSV2gyZWpGTU1uWlBOalExZFROSmMwbEVhSFFyZVdweVVYbGlTbFZJTVZKb1JIVTNRbUpWWTNsRlNIZHZkbGh1UVN0QlMzZERDa05LY0VaSVVVbEVRVkZCUW05NVkzZEtWRUZQUW1kT1ZraFJPRUpCWmpoRlFrRk5RMEpoUVhkRmQxbEVWbEl3YkVKQmQzZERaMWxKUzNkWlFrSlJWVWdLUVhkSmQwUlJXVXBMYjFwSmFIWmpUa0ZSUlV4Q1VVRkVaMmRGUWtGRlozUnZSa0pGYzBoNU9WaFlUVWRpVGtoSFNXaHdRVUU1UXlzNWRIWmpOMngwVUFwSU9VUnhOVkF6UVhad1dEVnlTRUV3YzNSblYzQk9VSHBJWkRGa1NtaG5jbll3VkZKVVdGaHlVSGRCTUdsYVNXVjFUakJaYnpaVmNYSXpNQzlNTUd0dENsaGhWa2RKUlVZd1lUaERlRFZCTURsaVYxWllhMWhNWmxST2MwRTNhbUZLZDBkUU1tcEllVVJGZEVOMlNXbExUamgzUzJkYWFGTjBjbkZRZDA5bVkwd0tlRmRWVDJKelIwNXNURlZ4T1ZrdlNYWjRlSG8xYTB4SFVHSTBRV0ZxU1VaRFRuTlFSV2xHYjNkQlNGWXJkMkZyZVdWek1XWXllbFZRTHpCR2NEWXZPUXBSVFhreVdGcHZWbEJyTTJSMWMweFBRa2xqUkcxS05ERTRRbmt2Wms5T1EyazRiREZpY25FeFRESmlTRk4wVVVObFZUbHJTbUoyT1c5d1Nub3pjM2w0Q2xKTFRDOHJWbVZHTTJRMk1tWkNWV2Q0VEZvcmNuQnJTMjAxU1ZaalVuTlVRak4yVWpkdk9YQk1Ta2NyUjIxRmRFUXlTVDBLTFMwdExTMUZUa1FnUTBWU1ZFbEdTVU5CVkVVdExTMHRMUW89CiAgICBjbGllbnQta2V5LWRhdGE6IExTMHRMUzFDUlVkSlRpQlNVMEVnVUZKSlZrRlVSU0JMUlZrdExTMHRMUXBOU1VsRmIyZEpRa0ZCUzBOQlVVVkJibFZtYVV4aFVEVlJhbmRyTkdoQk5sUmpSbTFTTURGNGMwVTNZVGhMZDIxaFJHTnhOU3RrV21kWVVuY3JRMDVUQ25GRWVqRTBNWGgwSzNacU1WVnpiRzEwVkdJeWNsQllibWRTZFdSV01XaFNjMUkyTjBsc2JtMURVemd5TUdwRFdsTjRkekpMUXpVNFJWbFNRbUYwZGlzS2VUQkJhbE5KWjJ4T1J5OXhZelJvUjFwS00wUm5UemczV0ZRemFtdG5XR1p0YnpST01sbFNTV2h5Ums5bWJVSmthMm94UWpKdGRHRjVhbFJuZEZGdk53cHVaRmRvYWxCSmJXWnRPRTF2VkVkaE5VeEJSVzUwVTBsMkwweEJTR2RwZFZOQ2NUbFpNa051Tmk5clIwdExkRFkwWWxGMlFUZG5LMjFhTkVRNUswTnZDazV5VEdkeFNHaG1lVVkwYTJsa01IWkZRbkJ0VEVKNlpEaHhZak1yYUdzNVVucDRiSEZCWXl0Tk1FVm9kbm94VERKMlR6WTBOWFV6U1hOSlJHaDBLM2tLYW5KUmVXSktWVWd4VW1oRWRUZENZbFZqZVVWSWQyOTJXRzVCSzBGTGQwTkRTbkJHU0ZGSlJFRlJRVUpCYjBsQ1FVUkVabTV2VVV0UFZFZEJNbEZGTmdwbk4wOVllbTR5U0hKcVZsbFBObUZuUTBGMFNWbFdabEY0ZW5BMFFWbGlTME41ZUcxVVVUVkxjMWsxVldkSU9IRlJTRlU1VlVObWN5OW1VbUp5VjJaMkNqbFBjWGxJZG5GTmNuWlFXbkpTU21oellrUnhXV3hZV1VwQmRWcGhiVFpZVlVSVEsybElhRWxtYVRoMVR6Z3JUVFJFY25nek1HVlphak42YnpreVpIWUtNV013UlV0WUsxaEpTMXBwVVVnM05qSlFhbEpsYWpSNVRDOU5NMUZJTUhGRFdVaEZOMnAwVGxWM2VUVlBhVWR0TldKRU5HNWtaMVZ3VDAxNGRYVkdXQXAxYm1sTmVFazVLMm80VEhaWVVqRjRNWFJZTDNoVGRGcFhRV0p4WlROUlNsUm1ZWGxZWjBwa1luZG1WRTFwYkZsMUwyaFdjazE1VkRGM1RsbGpVemd3Q2xad2NXSlZTekJwVFV4TlJqZzBhVnAwUzA0MVJXcFBkbEJxUWs1R2JUaE5SbEpYWXpRMWFrYzFVME5EYkZWYVZVYzFPSHBJTDNCQ2VtbzBiMmN5YmtNS0wxZE1SMlJEUlVObldVVkJkM05sWVVObmIzQTRiR0UwYWtKRlFUSjVZM2s0TkhkaWVFdEtTbWd5U0hGYU1uQTVOM0JZZW5CbU5HOXRSR053TVVFNU1RcENaWFl2WjJ4R1kxWlJVM1JRZFV4TVpUSlRaMVJEYm5SUEsyZGhOVU5VUjFKTmNrcHVkek53U0hOa1ZsY3JaMEZuTnpOUVVrMWtLemhxUm5WSloxWlJDa2xwVldSdWJIWnpZWEZXWXpOclNWQXJVbkl2V0d4SmVVUlRXbEZqVlZsVmJGcDBjMGhUUjNsQ1p6WlROMmcxV0hKSk5uZFFjMnREWjFsRlFYcHlZMG9LTmxSTGJFNWxZMXBTTmsxUmVuVmFTRUpMZGpCMVFURlBWbmhCYWsxclNucENPSGRpYnpKRVQwbEVVMHB0Y2s1clVFVlROblV3VHpSQ2JWQnlORUpEWXdwVmEzTnBkaXRFTTA1R2NGSXZRVkpEWkV0WFkxWnZTSEZ3UkVNclYyVlJXVWcySzJ0VmIwMHZZbGRTVlRZcmFEUTJSSEJXWjFsQlF6SkVSa2N4TVhKSUNsRkJlVkZ3WkV0SE5rUjVVbGh1YWtOU2RFbDJiSEZqU214elR6TXhVSEJtVGt0SVlWZGlWVU5uV1VJME1VUmxhbWxOUTFOMUsyUktlRWhoUTJaUFQyNEtiamN2VDBWR2NVdHhSbmQzUVZOcVZEWTJUekF3YlVNcldrZEpTRlZtZEc1VVFXODFMMng2U0hoT2JuZERLMHh6ZW1aU0szRjJaVzB4VVhCNFozSnhXUXBFUkVKa1FXUlZaWEZqYzNrNE1ETXlZVGRZTjFaalNsUllaMU5VUTBSaE5IQTRkV2sxUTFaYU1YWkZkMVV3ZW5wWVEwMVZhMDl5U25Fek16RmtkVXh2Q2tZdk1VOWhVWFV4Y0VNM1ZHeHVSVTg0Tkd0NmRWRkxRbWRJWmxoM1IxRXhWek5xWW5GV2JGSXpXVE5IUkZsb1ZURlJWR0p0WmpVMWRWbzBTM3BxZGtrS1FtOUlOblp4V1VoTU0zWkVOSGx1V25kMFUwMXNlVlZHVTNoMlZrcFlka1ZYWW05aGFEbGpVRkZvWVRCMFdEWkNRVkpzWVV0d2JsVk9TbVJUTWxCWFRnb3dPWFZQWjJabkszZGphMlpPZGtkVU1rUjZRelZYVGtoa2ExZHRUMGgzVjJwTmJrZFBlbHB4TTNaRFFsaEZha2t4U2xGdU9WUlJNVGxrT0VFMVN6bFVDalJ2TkRsQmIwZEJUMU56TVVGR1NtOXljWGs0UTFCeFIyMXFZbGxVV0RrM1NVaE1OMlZaUWpsMlEyTmllVmt6VHpWck5FaHdPVlF3UW5SRE0xcGpka01LYUVaUlNUQkNjek5DZGpSeFdFSlZlVVo1V2tsQldHaFdjVVJLUVM4MEwyNXNabXh1UjJscU9WaFFSVEZrWXpKdVdHSjBhM0E1Vm5Bek1VZ3haR2R4ZHdweFQybHZibkJsUnpWSWNEUXJVRWREV213MWVtdHJUVXBuVWxaRVNWZzVXWEJ2ZFZSWU1GcG9TMHhCZEVaUVVEQXZXRUU5Q2kwdExTMHRSVTVFSUZKVFFTQlFVa2xXUVZSRklFdEZXUzB0TFMwdENnPT0K",
      "endpointIP": "",
      "masterNodes": [
        {
          "cpu": 4,
          "memory": 16384,
          "name": "master",
          "networks": [
            {
              "label": "MANAGEMENT",
              "networkName": "/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0",
              "nameservers": [
                "10.246.2.9"
              ]
            }
          ],
          "storage": 50,
          "replica": 1,
          "labels": [],
          "cloneMode": "linkedClone"
        }
      ],
      "workerNodes": null,
      "vimId": "vmware_651E9CA99EDE4B7AA1A64B0BCAA69A0B",
      "error": ""
    }
  ]
}`

var emptyClustersSpec = `
{
  "Clusters": [

  ]
}`

var testJsonCluster = `{
  "id": "868636c9-868f-49fb-a6df-6a0d2d137146",
  "clusterName": "edge-test01",
  "clusterType": "WORKLOAD",
  "vsphereClusterName": "hubsite",
  "managementClusterId": "9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae",
  "hcxUUID": "20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6",
  "status": "ACTIVE",
  "activeTasksCount": 0,
  "clusterTemplate": {
    "name": "myworkload",
    "version": "",
    "id": "c3e006c1-e6aa-4591-950b-6f3bedd944d3"
  },
  "clusterId": "",
  "clusterUrl": "https://10.241.7.228:6443",
  "kubeConfig": "YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VONWVrTkRRV0pQWjBGM1NVSkJaMGxDUVVSQlRrSm5hM0ZvYTJsSE9YY3dRa0ZSYzBaQlJFRldUVkpOZDBWUldVUldVVkZFUlhkd2NtUlhTbXdLWTIwMWJHUkhWbnBOUWpSWVJGUkplRTFFV1hsTlZFVjZUbFJOZWsweGIxaEVWRTE0VFVSWmVFOVVSWHBPVkdkNlRURnZkMFpVUlZSTlFrVkhRVEZWUlFwQmVFMUxZVE5XYVZwWVNuVmFXRkpzWTNwRFEwRlRTWGRFVVZsS1MyOWFTV2gyWTA1QlVVVkNRbEZCUkdkblJWQkJSRU5EUVZGdlEyZG5SVUpCVERsWkNqUjBWamRVVjBGaVlXeEZURTg0YkN0d09WcElVMGhHVTJOUVJWcGxhRzFhVEZnNFZIaHFTa3N5WXpWWlVHNTBNVVU1UWsxSE4yazRVREEwV2tKWVN6UUtRazlKSzFKUlIwNUtWWEl2UVhKNE0wWkxOMDh4V1ZwaU9XdzBXREF3V2psM0wzWlFhVE5RTUVoWk5uTkNhV2xyUlROd01XRnBiRkV5UjA4MlJ6WjJWd3AxWVVOTmVGZ3lVWHBwWTJsTlNEazNkR3QyZDI1bUt6QjNWRlZLVjJkVFlrdDZkblJoYVdGaldYaHlXa3AxT1VNeWJEbFdNVU0zYUUxSmNVUlRhRUZGQ2tGWWF6TnFZV3RIV0VKaFoycHBkbXBMYmxBdmJ6bDJjR3BxTlZOeFVtVTRUa0pzVGs5emVqSnJaVUpXZEhGdFJuZ3hlbEpLZEZONU1tZEdVbUUxYWpnS04ySk1jSEpRWVZWemF5czRhM0JFVUc5aFdWQkllRUZqYlhkUmJIWXpXbGxYVFVkeVZqSlJLMVJVTUc5QlJtSmhaWGxtWlhOdFJEWmxaVEU1VjNneGN3cE1ObXRLYlRKQ1RVOWtjMUJhWkdaNGRDOXJRMEYzUlVGQllVMXRUVU5SZDBSbldVUldVakJRUVZGSUwwSkJVVVJCWjB0clRVSkpSMEV4VldSRmQwVkNDaTkzVVVsTlFWbENRV1k0UTBGUlFYZEVVVmxLUzI5YVNXaDJZMDVCVVVWTVFsRkJSR2RuUlVKQlMxVnRWRVpQT0Vwc1JUZEtlQ3RNZFdsUWVFOHdTSG9LVUU5eVdHZDBlbEIxVTFSdmJGbDJTVTVqUTJsQ0swbHlTVlowTjBReWJHUmpWMkpSY1ZwcVpXRjBObmg2VDBkMVoxYzRXWEYzU0d4Mk5uZHVTMWRPU2dwRGVVMDFhRGxOVDNWclJESjVkbmMxTlZCRk1XMTZaRVpyUVZSamJVRlZNQ3RzZFROUlNEQlpZa1psU1VkV1dsVkJPV04zVXlzMVlVbzFkWGMxYnlzMUNpdDJZM0kwYlRsNVJVVjZlWGhaZG1wNGNGTm5RMGx4ZEc1SFdpdHZPSFpNUW10UFFpdGlVRGhvVDFSalFuUk9VM1pETm1rd1ZWTkZXRTVMZHl0V2Vtb0tMMVpQU0VwMFVISk5kVW95VVc4M1QxaElVa0Z3U0cxRlVuTjZhazU2ZWxOaWFIaHBSV2xDVUZneFYyWjRUMHd4V2prMmNGVmFiMlJNUWxwbVJtTmtTZ3BhY0dZMFp6UmFjVFExYjJZck0xRjZjM2hoTnk5WlowbzFMMXBqUldwVWIwWkpZak5DUTFJelZ6SlNNalZ3VXpOV2JVNHlNR1ZvVkhSMWFrSnZXakE5Q2kwdExTMHRSVTVFSUVORlVsUkpSa2xEUVZSRkxTMHRMUzBLCiAgICBzZXJ2ZXI6IGh0dHBzOi8vMTAuMjQxLjcuMjI4OjY0NDMKICBuYW1lOiBlZGdlLXRlc3QwMQpjb250ZXh0czoKLSBjb250ZXh0OgogICAgY2x1c3RlcjogZWRnZS10ZXN0MDEKICAgIHVzZXI6IGVkZ2UtdGVzdDAxLWFkbWluCiAgbmFtZTogZWRnZS10ZXN0MDEtYWRtaW5AZWRnZS10ZXN0MDEKY3VycmVudC1jb250ZXh0OiBlZGdlLXRlc3QwMS1hZG1pbkBlZGdlLXRlc3QwMQpraW5kOiBDb25maWcKcHJlZmVyZW5jZXM6IHt9CnVzZXJzOgotIG5hbWU6IGVkZ2UtdGVzdDAxLWFkbWluCiAgdXNlcjoKICAgIGNsaWVudC1jZXJ0aWZpY2F0ZS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VNNGFrTkRRV1J4WjBGM1NVSkJaMGxKUkRoa1NqZEJVVVJpWWsxM1JGRlpTa3R2V2tsb2RtTk9RVkZGVEVKUlFYZEdWRVZVVFVKRlIwRXhWVVVLUVhoTlMyRXpWbWxhV0VwMVdsaFNiR042UVdWR2R6QjVUVlJCTWsxcVJYaE5lbFY2VFhwT1lVWjNNSGxOYWtFeVRXcEZlRTE2VlRSTmVsSmhUVVJSZUFwR2VrRldRbWRPVmtKQmIxUkViazQxWXpOU2JHSlVjSFJaV0U0d1dsaEtlazFTYTNkR2QxbEVWbEZSUkVWNFFuSmtWMHBzWTIwMWJHUkhWbnBNVjBackNtSlhiSFZOU1VsQ1NXcEJUa0puYTNGb2EybEhPWGN3UWtGUlJVWkJRVTlEUVZFNFFVMUpTVUpEWjB0RFFWRkZRVzVWWm1sTVlWQTFVV3AzYXpSb1FUWUtWR05HYlZJd01YaHpSVGRoT0V0M2JXRkVZM0UxSzJSYVoxaFNkeXREVGxOeFJIb3hOREY0ZEN0MmFqRlZjMnh0ZEZSaU1uSlFXRzVuVW5Wa1ZqRm9VZ3B6VWpZM1NXeHViVU5UT0RJd2FrTmFVM2gzTWt0RE5UaEZXVkpDWVhSMksza3dRV3BUU1dkc1RrY3ZjV00wYUVkYVNqTkVaMDg0TjFoVU0ycHJaMWhtQ20xdk5FNHlXVkpKYUhKR1QyWnRRbVJyYWpGQ01tMTBZWGxxVkdkMFVXODNibVJYYUdwUVNXMW1iVGhOYjFSSFlUVk1RVVZ1ZEZOSmRpOU1RVWhuYVhVS1UwSnhPVmt5UTI0MkwydEhTMHQwTmpSaVVYWkJOMmNyYlZvMFJEa3JRMjlPY2t4bmNVaG9abmxHTkd0cFpEQjJSVUp3YlV4Q2VtUTRjV0l6SzJock9RcFNlbmhzY1VGakswMHdSV2gyZWpGTU1uWlBOalExZFROSmMwbEVhSFFyZVdweVVYbGlTbFZJTVZKb1JIVTNRbUpWWTNsRlNIZHZkbGh1UVN0QlMzZERDa05LY0VaSVVVbEVRVkZCUW05NVkzZEtWRUZQUW1kT1ZraFJPRUpCWmpoRlFrRk5RMEpoUVhkRmQxbEVWbEl3YkVKQmQzZERaMWxKUzNkWlFrSlJWVWdLUVhkSmQwUlJXVXBMYjFwSmFIWmpUa0ZSUlV4Q1VVRkVaMmRGUWtGRlozUnZSa0pGYzBoNU9WaFlUVWRpVGtoSFNXaHdRVUU1UXlzNWRIWmpOMngwVUFwSU9VUnhOVkF6UVhad1dEVnlTRUV3YzNSblYzQk9VSHBJWkRGa1NtaG5jbll3VkZKVVdGaHlVSGRCTUdsYVNXVjFUakJaYnpaVmNYSXpNQzlNTUd0dENsaGhWa2RKUlVZd1lUaERlRFZCTURsaVYxWllhMWhNWmxST2MwRTNhbUZLZDBkUU1tcEllVVJGZEVOMlNXbExUamgzUzJkYWFGTjBjbkZRZDA5bVkwd0tlRmRWVDJKelIwNXNURlZ4T1ZrdlNYWjRlSG8xYTB4SFVHSTBRV0ZxU1VaRFRuTlFSV2xHYjNkQlNGWXJkMkZyZVdWek1XWXllbFZRTHpCR2NEWXZPUXBSVFhreVdGcHZWbEJyTTJSMWMweFBRa2xqUkcxS05ERTRRbmt2Wms5T1EyazRiREZpY25FeFRESmlTRk4wVVVObFZUbHJTbUoyT1c5d1Nub3pjM2w0Q2xKTFRDOHJWbVZHTTJRMk1tWkNWV2Q0VEZvcmNuQnJTMjAxU1ZaalVuTlVRak4yVWpkdk9YQk1Ta2NyUjIxRmRFUXlTVDBLTFMwdExTMUZUa1FnUTBWU1ZFbEdTVU5CVkVVdExTMHRMUW89CiAgICBjbGllbnQta2V5LWRhdGE6IExTMHRMUzFDUlVkSlRpQlNVMEVnVUZKSlZrRlVSU0JMUlZrdExTMHRMUXBOU1VsRmIyZEpRa0ZCUzBOQlVVVkJibFZtYVV4aFVEVlJhbmRyTkdoQk5sUmpSbTFTTURGNGMwVTNZVGhMZDIxaFJHTnhOU3RrV21kWVVuY3JRMDVUQ25GRWVqRTBNWGgwSzNacU1WVnpiRzEwVkdJeWNsQllibWRTZFdSV01XaFNjMUkyTjBsc2JtMURVemd5TUdwRFdsTjRkekpMUXpVNFJWbFNRbUYwZGlzS2VUQkJhbE5KWjJ4T1J5OXhZelJvUjFwS00wUm5UemczV0ZRemFtdG5XR1p0YnpST01sbFNTV2h5Ums5bWJVSmthMm94UWpKdGRHRjVhbFJuZEZGdk53cHVaRmRvYWxCSmJXWnRPRTF2VkVkaE5VeEJSVzUwVTBsMkwweEJTR2RwZFZOQ2NUbFpNa051Tmk5clIwdExkRFkwWWxGMlFUZG5LMjFhTkVRNUswTnZDazV5VEdkeFNHaG1lVVkwYTJsa01IWkZRbkJ0VEVKNlpEaHhZak1yYUdzNVVucDRiSEZCWXl0Tk1FVm9kbm94VERKMlR6WTBOWFV6U1hOSlJHaDBLM2tLYW5KUmVXSktWVWd4VW1oRWRUZENZbFZqZVVWSWQyOTJXRzVCSzBGTGQwTkRTbkJHU0ZGSlJFRlJRVUpCYjBsQ1FVUkVabTV2VVV0UFZFZEJNbEZGTmdwbk4wOVllbTR5U0hKcVZsbFBObUZuUTBGMFNWbFdabEY0ZW5BMFFWbGlTME41ZUcxVVVUVkxjMWsxVldkSU9IRlJTRlU1VlVObWN5OW1VbUp5VjJaMkNqbFBjWGxJZG5GTmNuWlFXbkpTU21oellrUnhXV3hZV1VwQmRWcGhiVFpZVlVSVEsybElhRWxtYVRoMVR6Z3JUVFJFY25nek1HVlphak42YnpreVpIWUtNV013UlV0WUsxaEpTMXBwVVVnM05qSlFhbEpsYWpSNVRDOU5NMUZJTUhGRFdVaEZOMnAwVGxWM2VUVlBhVWR0TldKRU5HNWtaMVZ3VDAxNGRYVkdXQXAxYm1sTmVFazVLMm80VEhaWVVqRjRNWFJZTDNoVGRGcFhRV0p4WlROUlNsUm1ZWGxZWjBwa1luZG1WRTFwYkZsMUwyaFdjazE1VkRGM1RsbGpVemd3Q2xad2NXSlZTekJwVFV4TlJqZzBhVnAwUzA0MVJXcFBkbEJxUWs1R2JUaE5SbEpYWXpRMWFrYzFVME5EYkZWYVZVYzFPSHBJTDNCQ2VtbzBiMmN5YmtNS0wxZE1SMlJEUlVObldVVkJkM05sWVVObmIzQTRiR0UwYWtKRlFUSjVZM2s0TkhkaWVFdEtTbWd5U0hGYU1uQTVOM0JZZW5CbU5HOXRSR053TVVFNU1RcENaWFl2WjJ4R1kxWlJVM1JRZFV4TVpUSlRaMVJEYm5SUEsyZGhOVU5VUjFKTmNrcHVkek53U0hOa1ZsY3JaMEZuTnpOUVVrMWtLemhxUm5WSloxWlJDa2xwVldSdWJIWnpZWEZXWXpOclNWQXJVbkl2V0d4SmVVUlRXbEZqVlZsVmJGcDBjMGhUUjNsQ1p6WlROMmcxV0hKSk5uZFFjMnREWjFsRlFYcHlZMG9LTmxSTGJFNWxZMXBTTmsxUmVuVmFTRUpMZGpCMVFURlBWbmhCYWsxclNucENPSGRpYnpKRVQwbEVVMHB0Y2s1clVFVlROblV3VHpSQ2JWQnlORUpEWXdwVmEzTnBkaXRFTTA1R2NGSXZRVkpEWkV0WFkxWnZTSEZ3UkVNclYyVlJXVWcySzJ0VmIwMHZZbGRTVlRZcmFEUTJSSEJXWjFsQlF6SkVSa2N4TVhKSUNsRkJlVkZ3WkV0SE5rUjVVbGh1YWtOU2RFbDJiSEZqU214elR6TXhVSEJtVGt0SVlWZGlWVU5uV1VJME1VUmxhbWxOUTFOMUsyUktlRWhoUTJaUFQyNEtiamN2VDBWR2NVdHhSbmQzUVZOcVZEWTJUekF3YlVNcldrZEpTRlZtZEc1VVFXODFMMng2U0hoT2JuZERLMHh6ZW1aU0szRjJaVzB4VVhCNFozSnhXUXBFUkVKa1FXUlZaWEZqYzNrNE1ETXlZVGRZTjFaalNsUllaMU5VUTBSaE5IQTRkV2sxUTFaYU1YWkZkMVV3ZW5wWVEwMVZhMDl5U25Fek16RmtkVXh2Q2tZdk1VOWhVWFV4Y0VNM1ZHeHVSVTg0Tkd0NmRWRkxRbWRJWmxoM1IxRXhWek5xWW5GV2JGSXpXVE5IUkZsb1ZURlJWR0p0WmpVMWRWbzBTM3BxZGtrS1FtOUlOblp4V1VoTU0zWkVOSGx1V25kMFUwMXNlVlZHVTNoMlZrcFlka1ZYWW05aGFEbGpVRkZvWVRCMFdEWkNRVkpzWVV0d2JsVk9TbVJUTWxCWFRnb3dPWFZQWjJabkszZGphMlpPZGtkVU1rUjZRelZYVGtoa2ExZHRUMGgzVjJwTmJrZFBlbHB4TTNaRFFsaEZha2t4U2xGdU9WUlJNVGxrT0VFMVN6bFVDalJ2TkRsQmIwZEJUMU56TVVGR1NtOXljWGs0UTFCeFIyMXFZbGxVV0RrM1NVaE1OMlZaUWpsMlEyTmllVmt6VHpWck5FaHdPVlF3UW5SRE0xcGpka01LYUVaUlNUQkNjek5DZGpSeFdFSlZlVVo1V2tsQldHaFdjVVJLUVM4MEwyNXNabXh1UjJscU9WaFFSVEZrWXpKdVdHSjBhM0E1Vm5Bek1VZ3haR2R4ZHdweFQybHZibkJsUnpWSWNEUXJVRWREV213MWVtdHJUVXBuVWxaRVNWZzVXWEJ2ZFZSWU1GcG9TMHhCZEVaUVVEQXZXRUU5Q2kwdExTMHRSVTVFSUZKVFFTQlFVa2xXUVZSRklFdEZXUzB0TFMwdENnPT0K",
  "endpointIP": "",
  "masterNodes": [
    {
      "cpu": 4,
      "memory": 16384,
      "name": "master",
      "networks": [
        {
          "label": "MANAGEMENT",
          "networkName": "/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0",
          "nameservers": [
            "10.246.2.9"
          ]
        }
      ],
      "storage": 50,
      "replica": 1,
      "labels": [],
      "cloneMode": "linkedClone"
    }
  ],
  "workerNodes": null,
  "vimId": "vmware_651E9CA99EDE4B7AA1A64B0BCAA69A0B",
  "error": ""
}`

var YamlPartialSpec = `
id: 868636c9-868f-49fb-a6df-6a0d2d137146
clustername: edge-test01
clustertype: WORKLOAD
vsphereclustername: hubsite
managementclusterid: 9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae
hcxuuid: 20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6
status: ACTIVE
activetaskscount: 0
clustertemplate:
    name: myworkload
    version: ""
    id: c3e006c1-e6aa-4591-950b-6f3bedd944d3
`

var YamlBrokenSpec = `
id: 868636c9-868f-49fb-a6df-6a0d2d137146
clustername: edge-test01
clustertype: WORKLOAD
vsphereclustername: hubsite
managementclusterid: 9ceb62ef-c48d-4504-86c5-cc9ce6ae1aae
hcxuuid: 20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6
status: ACTIVE
activetaskscount: 0
clustertemplate:
    name: myworkload
     version: ""
         id: c3e006c1-e6aa-4591-950b-6f3bedd944d3
`
