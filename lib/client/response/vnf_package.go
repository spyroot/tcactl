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
	"encoding/json"
	"fmt"
	"github.com/spyroot/tcactl/lib/models"
	"reflect"
	"strings"
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

// CatalogNotFound error raised if tenant cloud not found
type CatalogNotFound struct {
	errMsg string
}

//
func (m *CatalogNotFound) Error() string {
	return "Catalog entity '" + m.errMsg + "' not found"
}

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
	PID                string                 `json:"id" yaml:"pid"`
	VnfdID             string                 `json:"vnfdId" yaml:"vnfdId"`
	VnfProvider        string                 `json:"vnfProvider" yaml:"vnfProvider"`
	VnfProductName     string                 `json:"vnfProductName" yaml:"vnfProductName"`
	VnfSoftwareVersion string                 `json:"vnfSoftwareVersion" yaml:"vnfSoftwareVersion"`
	VnfdVersion        string                 `json:"vnfdVersion" yaml:"vnfdVersion"`
	OnboardingState    string                 `json:"onboardingState" yaml:"onboardingState"`
	OperationalState   string                 `json:"operationalState" yaml:"operationalState"`
	UsageState         string                 `json:"usageState" yaml:"usage_state"`
	VnfmInfo           []interface{}          `json:"vnfmInfo" yaml:"vnfmInfo"`
	UserDefinedData    models.UserDefinedData `yaml:"user_defined_data"`
}

// VnfPackages - array of VNF Packages.
type VnfPackages struct {
	Packages []VnfPackage
}

// GetField - return struct field value
func (t *VnfPackage) GetField(field string) string {

	r := reflect.ValueOf(t)
	fields, _ := t.GetFields()
	if _, ok := fields[field]; ok {
		f := reflect.Indirect(r).FieldByName(strings.Title(field))
		return f.String()
	}

	return ""
}

// GetFields return VduPackage fields name as
// map[string], each key is field name
func (t *VnfPackage) GetFields() (map[string]interface{}, error) {

	var m map[string]interface{}

	b, err := json.Marshal(t)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}

	return m, nil
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
		if q == VnfProductName && f(v.VnfProductName) {
			filtered = append(filtered, v)

		}
		if q == VnfdId && f(v.VnfdID) {
			filtered = append(filtered, v)
		}
		if q == OperationalState && f(v.VnfdID) {
			filtered = append(filtered, v)
		}
	}

	return filtered, nil
}
