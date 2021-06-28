package main

import (
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/pkg/io"
	"log"
)

var testTenant = `
{
  "items": [
    {
      "tenantId": "BDC07231F50A4536AA6DCF6B8C04BA5C",
      "vimName": "edge-test01",
      "tenantName": "edge-test01",
      "hcxCloudUrl": "https://test.io",
      "username": "edge-test01",
      "vimType": "KUBERNETES",
      "vimUrl": "https://10.241.7.228:6443",
      "hcxUUID": "20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6",
      "hcxTenantId": "0fd92808-4300-4eb5-a2d1-b72ea186727d",
      "location": {
        "city": "Palo Alto",
        "country": "United States of America",
        "cityAscii": "Palo Alto",
        "latitude": 37.3913,
        "longitude": -122.1467
      },
      "vimId": "vmware_651E9CA99EDE4B7AA1A64B0BCAA69A0B",
      "audit": {
        "creationUser": "Administrator@VSPHERE.LOCAL",
        "creationTimestamp": "Mon Jun 21 14:05:16 GMT 2021"
      },
      "connection": {
        "status": "ok",
        "remoteStatus": "ok",
        "vimConnectionStatus": "ok"
      },
      "compatible": true,
      "id": "BDC07231F50A4536AA6DCF6B8C04BA5C",
      "name": "edge-test01",
      "authType": "kubeconfig",
      "clusterName": "edge-test01",
      "clusterNodeConfigList": null,
      "hasSupportedKubernetesVersion": false,
      "clusterStatus": "",
      "isCustomizable": false
    }
  ]
}`

var testTenantYaml = `
items:
  - tenantId: BDC07231F50A4536AA6DCF6B8C04BA5C
    vimName: edge-test01
    tenantName: edge-test01
    hcxCloudUrl: https://test.io
    username: edge-test01
    password: ""
    vimType: KUBERNETES
    vimUrl: https://10.241.7.228:6443
    hcxUUID: 20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6
    hcxTenantId: 0fd92808-4300-4eb5-a2d1-b72ea186727d
    location:
      city: Palo Alto
      country: United States of America
      cityAscii: Palo Alto
      latitude: 37.3913
      longitude: -122.1467
    vimId: vmware_651E9CA99EDE4B7AA1A64B0BCAA69A0B
    audit:
      creationUser: Administrator@VSPHERE.LOCAL
      creationTimestamp: Mon Jun 21 14:05:16 GMT 2021
    connection:
      status: ok
      remoteStatus: ok
      vimConnectionStatus: ok
    compatible: true
    id: BDC07231F50A4536AA6DCF6B8C04BA5C
    name: edge-test01
    authType: kubeconfig
    clusterName: edge-test01
    clusterNodeConfigList: []
    hasSupportedKubernetesVersion: false
    clusterStatus: ""
    isCustomizable: false
`

func main() {

	var tenant response.Tenants

	v, err := tenant.InstanceSpecsFromString(testTenantYaml)
	if err != nil {
		log.Fatal(err)
	}

	io.YamlPrinter(v)
}
