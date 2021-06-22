package request

import (
	"encoding/json"
	"testing"
)

func TestNewClusterYaml(t *testing.T) {

	data := `{"name":"test01",
	"clusterType":"WORKLOAD",
	"clusterTemplateId":"06060d88-5b85-4ebb-b9f6-615096d294be",
	"hcxCloudUrl":"https://test.io",
	"vmTemplate":"/Datacenter/vm/templates/photon-3-kube-v1.20.4+vmware.1",
	"endpointIP":"1.1.1.1",
    "placementParams":[
			{"type":"Folder","name":"templates"},
			{"type":"Datastore","name":"vsanDatastore"},
			{"type":"ResourcePool","name":"pod03"},
			{"type":"IsValidClusterCompute","name":"core"}],
	"clusterConfig":
		{"csi":[
				{
				"name":"nfs_client",
				"properties":{
					"serverIP":"10.241.0.250",
					"mountPath":"w3-nfv-pst-01-mus"}
				},
				{
					"name":"vsphere-csi",
					"properties":{
					"datastoreUrl":"ds:///vmfs/volumes/vsan:525abe9561d93fa1-4add35a851137588/",
					"datastoreName":"vsanDatastore"}}],
			"tools":[
				{
				"name":"harbor",
				"properties":{"type":"extension",
				"extensionId":"9d0d4ff4-1963-4d89-ac15-2d856768deeb"}}]},

				"managementClusterId":"85b839fb-f7ef-436f-8832-1045c6875a01",
				"clusterPassword":"Vk13YXJlMSE=",
"masterNodes":[{"name":"master","placementParams":[{"type":"IsValidClusterCompute","name":"core"},{"type":"Datastore","name":"vsanDatastore"},{"type":"ResourcePool","name":"pod03"}],"networks":[{"label":"MANAGEMENT","networkName":"/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0"}]}],"workerNodes":[{"name":"default-pool01","placementParams":[{"type":"IsValidClusterCompute","name":"core"},{"type":"Datastore","name":"vsanDatastore"},{"type":"ResourcePool","name":"pod03"}],"networks":[{"label":"MANAGEMENT","networkName":"/Datacenter/network/tkg-dhcp-vlan1007-10.241.7.0"}]}]}`

	tests := []struct {
		name string
		args string
	}{
		{
			name: "yaml parse",
			args: data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var spec Cluster
			source := []byte(tt.args)
			err := json.Unmarshal(source, &spec)
			if err != nil {
				t.Errorf("Failed.")
			}
		})
	}
}
