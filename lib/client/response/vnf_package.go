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
	"github.com/spyroot/tcactl/lib/api_errors"
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
	PID                string                  `json:"id" yaml:"pid"`
	VnfdID             string                  `json:"vnfdId" yaml:"vnfdId"`
	VnfProvider        string                  `json:"vnfProvider" yaml:"vnfProvider"`
	VnfProductName     string                  `json:"vnfProductName" yaml:"vnfProductName"`
	VnfSoftwareVersion string                  `json:"vnfSoftwareVersion" yaml:"vnfSoftwareVersion"`
	VnfdVersion        string                  `json:"vnfdVersion" yaml:"vnfdVersion"`
	OnboardingState    string                  `json:"onboardingState" yaml:"onboardingState"`
	OperationalState   string                  `json:"operationalState" yaml:"operationalState"`
	UsageState         string                  `json:"usageState" yaml:"usage_state"`
	VnfmInfo           []interface{}           `json:"vnfmInfo" yaml:"vnfmInfo"`
	UserDefinedData    *models.UserDefinedData `json:"userDefinedData" yaml:"userDefinedData"`
}

// IsEnabled return true if catalog entity enabled
func (p *VnfPackage) IsEnabled() bool {
	return p != nil && p.OperationalState == "ENABLED"
}

// IsOnboarded return true if successfully onboarded
func (p *VnfPackage) IsOnboarded() bool {
	return p != nil && p.OnboardingState == "ONBOARDED"
}

// IsCnf return true if catalog entity is CNF
func (p *VnfPackage) IsCnf() bool {
	return p != nil && p.UserDefinedData != nil && p.UserDefinedData.NfType == "CNF"
}

// GetField - return struct field value
func (p *VnfPackage) GetField(field string) string {

	r := reflect.ValueOf(p)
	fields, _ := p.GetFields()
	if _, ok := fields[field]; ok {
		f := reflect.Indirect(r).FieldByName(strings.Title(field))
		return f.String()
	}

	return ""
}

// GetFields return VduPackage fields name as
// map[string], each key is field name
func (p *VnfPackage) GetFields() (map[string]interface{}, error) {

	var m map[string]interface{}

	b, err := json.Marshal(p)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}

	return m, nil
}

// VnfPackages - list of VnfPackage
type VnfPackages struct {
	Entity []VnfPackage
}

// GetVnfdID - find by either VnfdID or ProductName, id
// or user defined data, the main objective to resolve some name
// to catalog entity.
func (v *VnfPackages) GetVnfdID(NameOrId string) (*VnfPackage, error) {

	if v == nil {
		return nil, fmt.Errorf("nil object")
	}

	q := strings.ToLower(NameOrId)

	for _, p := range v.Entity {
		// by product name
		if p.VnfProductName == q || p.VnfdID == q || p.PID == q {
			return &p, nil
		}

		// by User define data
		if p.UserDefinedData != nil && strings.ToLower(p.UserDefinedData.Name) == q {
			return &p, nil
		}
	}

	return nil, api_errors.NewCatalogNotFound(q)
}

// FindByCatalogName - find by either VnfdID or ProductName, id
// the main objective to resolve some name to catalog
// entity.
func (v *VnfPackages) FindByCatalogName(q string) (*VnfPackage, error) {

	if v == nil {
		return nil, fmt.Errorf("nil object")
	}

	for _, p := range v.Entity {
		if p.UserDefinedData != nil {
			if strings.ToLower(p.UserDefinedData.Name) == strings.ToLower(q) {
				return &p, nil
			}
		}
	}

	return nil, api_errors.NewCatalogNotFound(q)
}

// Filter filters respond based on filter type and pass to callback
func (v *VnfPackages) Filter(q VnfdFilterType, f func(string) bool) ([]VnfPackage, error) {

	filtered := make([]VnfPackage, 0)
	for _, v := range v.Entity {
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
