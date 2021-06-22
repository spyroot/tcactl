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
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

// VnfPackage - TCA VNF Package TCA format
type VnfPackage struct {
	PID                string        `json:"id"`
	VnfdID             string        `json:"vnfdId"`
	VnfProvider        string        `json:"vnfProvider"`
	VnfProductName     string        `json:"vnfProductName"`
	VnfSoftwareVersion string        `json:"vnfSoftwareVersion"`
	VnfdVersion        string        `json:"vnfdVersion"`
	OnboardingState    string        `json:"onboardingState"`
	OperationalState   string        `json:"operationalState"`
	UsageState         string        `json:"usageState"`
	VnfmInfo           []interface{} `json:"vnfmInfo"`
	UserDefinedData    struct {
		Name                   string            `json:"name"`
		Tags                   []interface{}     `json:"tags"`
		NfType                 string            `json:"nfType"`
		ManagedBy              InternalManagedBy `json:"managedBy"`
		LocalFilePath          string            `json:"localFilePath"`
		LastUpdated            time.Time         `json:"lastUpdated"`
		LastUpdateEnterprise   string            `json:"lastUpdateEnterprise"`
		LastUpdateOrganization string            `json:"lastUpdateOrganization"`
		LastUpdateUser         string            `json:"lastUpdateUser"`
		CreationDate           time.Time         `json:"creationDate"`
		CreationEnterprise     string            `json:"creationEnterprise"`
		CreationOrganization   string            `json:"creationOrganization"`
		CreationUser           string            `json:"creationUser"`
		IsDeleted              bool              `json:"isDeleted"`
	}
}

// VnfPackages - array of VNF Packages.
type VnfPackages struct {
	Packages []VnfPackage
}

// GetField - return field from VNfPackage struct
func (p *VnfPackage) GetField(field string) string {
	r := reflect.ValueOf(p)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	return f.String()
}

// GetVnfdID - find by either VnfdID or ProductName, ID
func (v *VnfPackages) GetVnfdID(q string) (*VnfPackage, error) {

	if v == nil {
		return nil, fmt.Errorf("nil object")
	}

	for _, p := range v.Packages {
		if p.VnfProductName == q || p.VnfdID == q || p.PID == q {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("package not found")
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
