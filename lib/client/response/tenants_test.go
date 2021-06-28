package response

import (
	"github.com/spyroot/tcactl/lib/models"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var testTenantJson = `
{
  "items" : [ {
    "tenantId" : "4F5F6CAF841D4D1DA829104D32B5AA97",
    "vimName" : "core",
    "tenantName" : "DEFAULT",
    "hcxCloudUrl" : "https://test-test.test.io",
    "username" : "administrator@vsphere.local",
    "password" : "xxx",
    "tags" : [ ],
    "vimType" : "VC",
    "vimUrl" : "https://xxx.xxx.io",
    "hcxUUID" : "20210212053126765-51947d48-91b1-447e-ba40-668eb411f545",
    "hcxTenantId" : "20210212053126765-51947d48-91b1-447e-ba40-668eb411f545",
    "location" : {
      "city" : "Palo Alto",
      "country" : "United States of America",
      "cityAscii" : "Palo Alto",
      "latitude" : 37.3913,
      "longitude" : -122.1467
    },
    "vimId" : "vmware_FB40D3DE2967483FBF9033B451DC7571",
    "audit" : {
      "creationUser" : "Administrator@VSPHERE.LOCAL",
      "creationTimestamp" : "Mon Feb 15 03:47:42 GMT 2021"
    },
    "connection" : {
      "status" : "ok",
      "remoteStatus" : "ok"
    },
    "compatible" : true,
    "id" : "4F5F6CAF841D4D1DA829104D32B5AA97",
    "name" : "DEFAULT"
  }, {
    "tenantId" : "4E010FF9918D438EAFCE88B5EE7D81F8",
    "vimName" : "edge",
    "tenantName" : "DEFAULT",
    "hcxCloudUrl" : "https://xxxx.io",
    "username" : "Administrator@vsphere.local",
    "password" : "aaaa",
    "tags" : [ ],
    "vimType" : "VC",
    "vimUrl" : "https://xxx.xxxx.io",
    "hcxUUID" : "20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6",
    "hcxTenantId" : "20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6",
    "location" : {
      "city" : "Palo Alto",
      "country" : "United States of America",
      "cityAscii" : "Palo Alto",
      "latitude" : 37.3913,
      "longitude" : -122.1467
    },
    "vimId" : "vmware_8581A92D4C224708BF2A5C79E4ED5A3D",
    "audit" : {
      "creationUser" : "Administrator@VSPHERE.LOCAL",
      "creationTimestamp" : "Thu Jun 03 12:10:33 GMT 2021"
    },
    "connection" : {
      "status" : "ok",
      "remoteStatus" : "ok"
    },
    "compatible" : true,
    "id" : "4E010FF9918D438EAFCE88B5EE7D81F8",
    "name" : "DEFAULT"
  }, {
    "tenantId" : "BDC07231F50A4536AA6DCF6B8C04BA5C",
    "vimName" : "edge-test01",
    "tenantName" : "edge-test01",
    "hcxCloudUrl" : "https://xxxx.xxxx.io",
    "authType" : "kubeconfig",
    "clusterName" : "edge-test01",
    "username" : "edge-test01",
    "vimType" : "KUBERNETES",
    "tags" : [ ],
    "vimUrl" : "https://10.241.7.228:6443",
    "hcxTenantId" : "0fd92808-4300-4eb5-a2d1-b72ea186727d",
    "hcxUUID" : "20210602134140183-1ddd8717-de09-4143-acb5-e51fb372ebf6",
    "location" : {
      "city" : "Palo Alto",
      "country" : "United States of America",
      "cityAscii" : "Palo Alto",
      "latitude" : 37.3913,
      "longitude" : -122.1467
    },
    "vimId" : "vmware_651E9CA99EDE4B7AA1A64B0BCAA69A0B",
    "audit" : {
      "creationUser" : "Administrator@VSPHERE.LOCAL",
      "creationTimestamp" : "Mon Jun 21 14:05:16 GMT 2021"
    },
    "connection" : {
      "status" : "ok",
      "remoteStatus" : "ok",
      "vimConnectionStatus" : "ok"
    },
    "compatible" : true,
    "id" : "BDC07231F50A4536AA6DCF6B8C04BA5C",
    "name" : "edge-test01"
  } ]
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

type TeantSpecCreator func(s string) (*TenantSpecs, error)

func TestTenantsCreator(t1 *testing.T) {

	tests := []struct {
		name       string
		specString string
		wantErr    bool
		spec       Tenants
	}{
		{
			name:       "basic creator from json string",
			specString: testTenantJson,
			wantErr:    false,
		},
		{
			name:       "basic creator from json string",
			specString: testTenantYaml,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			instance, err := tt.spec.InstanceSpecsFromString(tt.specString)
			if !tt.wantErr {
				assert.Nil(t1, err)
			}
			if tt.wantErr {
				assert.NotNil(t1, err)
			}

			i, ok := instance.(*Tenants)
			assert.Equal(t1, ok, true)
			assert.Equal(t1, len(i.TenantsList), 1)
		})
	}
}

type filterCallback func(string) bool

func TestTenants_Filter(t1 *testing.T) {

	tests := []struct {
		name       string
		specString string
		wantErr    bool
		expected   string
		spec       Tenants
		filterType TenantCloudFilter
		filterCallback
	}{
		{
			name:       "basic positive case for vim type",
			specString: testTenantJson,
			wantErr:    false,
			expected:   "edge-test01",
			filterType: FilterVimType,
			filterCallback: func(s string) bool {
				return strings.ToLower(s) == models.VimTypeKubernetes
			},
		},
		{
			name:       "basic negative case for vim type",
			specString: testTenantJson,
			wantErr:    false,
			expected:   "",
			filterType: FilterVimType,
			filterCallback: func(s string) bool {
				return strings.ToLower(s) == models.VimTypeVmware
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			instance, err := tt.spec.InstanceSpecsFromString(tt.specString)
			assert.NoError(t1, err)

			i, ok := instance.(*Tenants)
			assert.Equal(t1, ok, true)
			assert.Equal(t1, len(i.TenantsList), 1)

			tenants := i.Filter(tt.filterType, tt.filterCallback)
			if (err != nil) != tt.wantErr {
				t1.Errorf("TestTenants_Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			found := false
			for _, tenant := range tenants {
				if tenant.VimName == tt.expected {
					found = true
				}
			}

			if len(tt.expected) == 0 && found {
				t1.Errorf("TestTenants_Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !found {
					t1.Errorf("TestTenants_Filter() error = %v, wantErr %v", err, tt.wantErr)
					return
				} else {
					t1.Logf("Found %s", tt.expected)
				}
			}
		})
	}
}

//func TestTenants_Creator(t1 *testing.T) {
//
//	tests := []struct {
//		name    	string
//		specString  string
//		creator     TeantSpecCreator
//		wantErr 	bool
//		spec        TenantSpecs
//	}{
//		{
//			name: "basic string",
//			specString: testTenantJson,
//			creator: 	TenantSpecsFromString,
//		},
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t1 *testing.T) {
//
//			instance, err := tt.creator(tt.specString)
//			if err != nil {
//				return
//			}
//
//		})
//	}
//}
