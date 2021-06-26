// Package response
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Mustafa mbayramo@vmware.com
package response

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// VnfdFilterType - vnfd filter types
type VnfdFilterType int32

const (

	// VnfProductName filter by VNF/CNF product name
	VnfProductName VnfdFilterType = 0

	// VnfdId filter by VNFD id
	VnfdId VnfdFilterType = 1

	// OperationalState filter type Operation Status in Catalog
	OperationalState VnfdFilterType = 2
)

// VnfPackagesError - TCA Error rest API error response format
type VnfPackagesError struct {
	Type     string `json:"type" yaml:"type"`
	Title    string `json:"title" yaml:"title"`
	Status   int    `json:"status" yaml:"status"`
	Detail   string `json:"detail" yaml:"detail"`
	Instance string `json:"instance" yaml:"instance"`
}

// VnfPackage - TCA VNF Package TCA format
type VnfPackage struct {
	PID                string        `json:"id" yaml:"pid"`
	VnfdID             string        `json:"vnfdId" yaml:"vnfd_id"`
	VnfProvider        string        `json:"vnfProvider" yaml:"vnf_provider"`
	VnfProductName     string        `json:"vnfProductName" yaml:"vnf_product_name"`
	VnfSoftwareVersion string        `json:"vnfSoftwareVersion" yaml:"vnf_software_version"`
	VnfdVersion        string        `json:"vnfdVersion" yaml:"vnfd_version"`
	OnboardingState    string        `json:"onboardingState" yaml:"onboarding_state"`
	OperationalState   string        `json:"operationalState" yaml:"operational_state"`
	UsageState         string        `json:"usageState" yaml:"usage_state"`
	VnfmInfo           []interface{} `json:"vnfmInfo" yaml:"vnfm_info"`
	UserDefinedData    struct {
		Name                   string            `json:"name" yaml:"name"`
		Tags                   []interface{}     `json:"tags" yaml:"tags"`
		NfType                 string            `json:"nfType" yaml:"nf_type"`
		ManagedBy              InternalManagedBy `json:"managedBy" yaml:"managed_by"`
		LocalFilePath          string            `json:"localFilePath" yaml:"local_file_path"`
		LastUpdated            time.Time         `json:"lastUpdated" yaml:"last_updated"`
		LastUpdateEnterprise   string            `json:"lastUpdateEnterprise" yaml:"last_update_enterprise"`
		LastUpdateOrganization string            `json:"lastUpdateOrganization" yaml:"last_update_organization"`
		LastUpdateUser         string            `json:"lastUpdateUser" yaml:"last_update_user"`
		CreationDate           time.Time         `json:"creationDate" yaml:"creation_date"`
		CreationEnterprise     string            `json:"creationEnterprise" yaml:"creation_enterprise"`
		CreationOrganization   string            `json:"creationOrganization" yaml:"creation_organization"`
		CreationUser           string            `json:"creationUser" yaml:"creation_user"`
		IsDeleted              bool              `json:"isDeleted" yaml:"is_deleted"`
	} `yaml:"user_defined_data"`
}

// VnfPackages - array of VNF Packages.
type VnfPackages struct {
	Packages []VnfPackage
}

// CatalogNotFound error raised if tenant cloud not found
type CatalogNotFound struct {
	errMsg string
}

//
func (m *CatalogNotFound) Error() string {
	return "Catalog entity '" + m.errMsg + "' not found"
}

// GetField - return field from VNfPackage struct
func (p *VnfPackage) GetField(field string) string {
	r := reflect.ValueOf(p)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	return f.String()
}

// GetVnfdID - find by either VnfdID or ProductName, id
// the main objective to resolve some name to catalog
// entity.
func (v *VnfPackages) GetVnfdID(q string) (*VnfPackage, error) {

	if v == nil {
		return nil, fmt.Errorf("nil object")
	}

	for _, p := range v.Packages {
		if p.VnfProductName == q || p.VnfdID == q || p.PID == q {
			return &p, nil
		}
	}

	return nil, &CatalogNotFound{q}
}

// Filter filters respond based on filter type and pass to callback
func (v *VnfPackages) Filter(q VnfdFilterType, f func(string) bool) ([]VnfPackage, error) {

	filtered := make([]VnfPackage, 0)
	for _, v := range v.Packages {
		if q == VnfProductName {
			if f(v.VnfProductName) {
				filtered = append(filtered, v)
			}
		}
		if q == VnfdId {
			if f(v.VnfdID) {
				filtered = append(filtered, v)
			}
		}
		if q == OperationalState {
			if f(v.VnfdID) {
				filtered = append(filtered, v)
			}
		}
	}

	return filtered, nil
}
