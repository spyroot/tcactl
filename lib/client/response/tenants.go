// Package respons
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
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/models"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type TenantCloudFilter int32

const (
	// FilterVimType filter by VIM type i.e CNF/VC
	FilterVimType TenantCloudFilter = 0

	// VimLocationCity filter by City
	VimLocationCity TenantCloudFilter = 1

	// VimLocationCountry filter by Country
	VimLocationCountry TenantCloudFilter = 2

	// FilterHcxUUID filter by Hcx UUID
	FilterHcxUUID TenantCloudFilter = 3

	// FilterName filter by name
	FilterName TenantCloudFilter = 4

	// FilterId filter by id
	FilterId TenantCloudFilter = 5

	// FilterVimId filter by vim id
	FilterVimId TenantCloudFilter = 6
)

// InvalidTenantsSpec error if specs invalid
type InvalidTenantsSpec struct {
	errMsg string
}

//
func (m *InvalidTenantsSpec) Error() string {
	return m.errMsg
}

// InvalidClusterSpec error if specs invalid
type InvalidClusterSpec struct {
	errMsg string
}

//
func (m *InvalidClusterSpec) Error() string {
	return m.errMsg
}

// TenantCloudNotFound error raised if tenant cloud not found
type TenantCloudNotFound struct {
	errMsg string
}

//
func (m *TenantCloudNotFound) Error() string {
	return "tenant '" + m.errMsg + "' not found"
}

// Tenants list of Tenants
type Tenants struct {
	TenantsList []TenantsDetails `json:"items" yaml:"items"`
}

// IsVMware - return if cloud provider is Vmware
func (t *TenantsDetails) IsVMware() bool {
	return strings.ToLower(t.VimType) == strings.ToLower(models.VimTypeVmware)
}

// GetTenantClouds return list of tenant clouds
func (t *Tenants) GetTenantClouds(s string, vimType string) (*TenantsDetails, error) {

	glog.Infof("Acquiring cloud provider %s details, type %s", s, vimType)

	if t == nil {
		return nil, fmt.Errorf("uninitialized object")
	}

	for _, t := range t.TenantsList {
		if strings.Contains(strings.ToLower(t.VimType), strings.ToLower(vimType)) {
			if strings.Contains(t.VimName, s) || strings.Contains(t.ID, s) {
				return &t, nil
			}
		}
	}

	return nil, &TenantCloudNotFound{errMsg: s}
}

// FindCloudProvider search for a cloud provider
func (t *Tenants) FindCloudProvider(s string) (*TenantsDetails, error) {

	if t == nil {
		return nil, fmt.Errorf("uninitialized object")
	}

	for _, t := range t.TenantsList {
		if strings.ToLower(t.VimName) == strings.ToLower(s) || t.ID == s || t.VimID == s {
			return &t, nil
		}
	}

	return nil, &TenantCloudNotFound{errMsg: s}
}

// Contains search for a cloud provider that contains
// name in either ID , or Vim Name or Vim ID.
// it partial match
func (t *Tenants) Contains(s string) (*TenantsDetails, error) {

	if t == nil {
		return nil, fmt.Errorf("uninitialized object")
	}

	for _, t := range t.TenantsList {
		if strings.Contains(t.VimName, s) ||
			strings.Contains(t.ID, s) ||
			strings.Contains(t.VimID, s) {
			return &t, nil
		}
	}

	return nil, &TenantCloudNotFound{errMsg: s}
}

// FindCloudLink search for a cloud provider cloud link
//
func (t *Tenants) FindCloudLink(s string) (*TenantsDetails, error) {

	if t == nil {
		return nil, fmt.Errorf("uninitialized object")
	}

	for _, t := range t.TenantsList {
		if strings.Contains(t.HcxCloudURL, s) {
			return &t, nil
		}
	}

	return nil, &TenantCloudNotFound{errMsg: s}
}

// Filter filter on specific filed
func (t *Tenants) Filter(q TenantCloudFilter, f func(string) bool) []TenantsDetails {

	filtered := make([]TenantsDetails, 0)
	for _, v := range t.TenantsList {
		if q == FilterVimType && f(v.VimType) {
			filtered = append(filtered, v)
		}
		if q == FilterId && f(v.ID) {
			filtered = append(filtered, v)
		}
		if q == FilterName && f(v.Name) {
			filtered = append(filtered, v)
		}
		if q == FilterHcxUUID && f(v.HcxUUID) {
			filtered = append(filtered, v)
		}
		if q == VimLocationCity && f(v.Location.City) {
			filtered = append(filtered, v)
		}
		if q == VimLocationCountry && f(v.Location.Country) {
			filtered = append(filtered, v)
		}
		if q == FilterVimId && f(v.VimID) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// TenantsSpecsFromFile - reads tenant spec from file
// and return TenantSpecs instance
func TenantsSpecsFromFile(fileName string) (*Tenants, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadTenantsSpec(file)
}

// TenantsSpecsFromString take string that hold entire spec
// passed to reader and return TenantSpecs instance
func TenantsSpecsFromString(str string) (*Tenants, error) {
	r := strings.NewReader(str)
	return ReadTenantsSpec(r)
}

// ReadTenantsSpec - Read tenants spec from io interface
// detects format and use either yaml or json parse
func ReadTenantsSpec(b io.Reader) (*Tenants, error) {

	var spec Tenants

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	} else {
		fmt.Println(reflect.TypeOf(err).String())
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	} else {
		fmt.Println(err)
	}

	return nil, &InvalidTenantsSpec{"unknown format"}
}

//InstanceSpecsFromString method return instance form string
func (t Tenants) InstanceSpecsFromString(s string) (interface{}, error) {
	return TenantsSpecsFromString(s)
}
