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
	"fmt"
	tca2 "github.com/spyroot/hestia/cmd/models"
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

	//
	VimTypeVmware = "VC"

	//
	VimTypeKubernetes = "kubernetes"
)

type AuditField struct {
	CreationUser      string `json:"creationUser" yaml:"creationUser"`
	CreationTimestamp string `json:"creationTimestamp" yaml:"creationTimestamp"`
}

type ClusterNodeConifgList struct {
	Labels           []string `json:"labels" yaml:"labels"`
	Id               string   `json:"id" yaml:"id"`
	Name             string   `json:"name" yaml:"name"`
	Status           string   `json:"status" yaml:"status"`
	ActiveTasksCount int      `json:"activeTasksCount" yaml:"activeTasksCount"`
	Compatible       bool     `json:"compatible" yaml:"compatible"`
}

// TenantsDetails Tenant Cloud Details
type TenantsDetails struct {
	TenantID                      string                  `json:"tenantId" yaml:"tenantId"`
	VimName                       string                  `json:"vimName" yaml:"vimName"`
	TenantName                    string                  `json:"tenantName" yaml:"tenantName"`
	HcxCloudURL                   string                  `json:"hcxCloudUrl" yaml:"hcxCloudUrl"`
	Username                      string                  `json:"username" yaml:"username"`
	Password                      string                  `json:"password,omitempty" yaml:"password"`
	VimType                       string                  `json:"vimType" yaml:"vimType"`
	VimURL                        string                  `json:"vimUrl" yaml:"vimUrl"`
	HcxUUID                       string                  `json:"hcxUUID" yaml:"hcxUUID"`
	HcxTenantID                   string                  `json:"hcxTenantId" yaml:"hcxTenantId"`
	Location                      tca2.Location           `json:"location" yaml:"location"`
	VimID                         string                  `json:"vimId" yaml:"vimId"`
	Audit                         AuditField              `json:"audit" yaml:"audit"`
	VimConn                       tca2.VimConnection      `json:"connection,omitempty" yaml:"connection"`
	Compatible                    bool                    `json:"compatible" yaml:"compatible"`
	ID                            string                  `json:"id" yaml:"id"`
	Name                          string                  `json:"name" yaml:"name"`
	AuthType                      string                  `json:"authType,omitempty" yaml:"authType"`
	ClusterName                   string                  `json:"clusterName,omitempty" yaml:"clusterName"`
	ClusterList                   []ClusterNodeConifgList `json:"clusterNodeConfigList" yaml:"clusterNodeConfigList"`
	HasSupportedKubernetesVersion bool                    `json:"hasSupportedKubernetesVersion" yaml:"hasSupportedKubernetesVersion"`
	ClusterStatus                 string                  `json:"clusterStatus" yaml:"clusterStatus"`
	IsCustomizable                bool                    `json:"isCustomizable" yaml:"isCustomizable"`
}

type TenantSpecs struct {
	CloudOwner    string           `json:"cloud_owner" yaml:"cloud_owner"`
	CloudRegionId string           `json:"cloud_region_id" yaml:"cloud_region_id"`
	VimId         string           `json:"vimId" yaml:"vimId"`
	VimName       string           `json:"vimName" yaml:"vimName"`
	Tenants       []TenantsDetails `json:"tenants" yaml:"tenants"`
}

// TenantCloudNotFound error raised if tenant cloud not found
type TenantCloudNotFound struct {
	errMsg string
}

//
func (m *TenantCloudNotFound) Error() string {
	return m.errMsg + " tenant cloud not found"
}

// Tenants list of Tenants
type Tenants struct {
	TenantsList []TenantsDetails `json:"items"`
}

// GetField - return field from VNfPackage struct
func (t *TenantsDetails) GetField(field string) string {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	return f.String()
}

// IsVMware - return if cloud provider is Vmware
func (t *TenantsDetails) IsVMware() bool {
	return t.VimType == VimTypeVmware
}

// GetTenantClouds return list of tenant clouds
func (t *Tenants) GetTenantClouds(s string, vimType string) (*TenantsDetails, error) {

	if t == nil {
		return nil, fmt.Errorf("uninitialized object")
	}

	for _, t := range t.TenantsList {
		if strings.Contains(t.VimType, vimType) {
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
		if strings.Contains(t.VimName, s) || strings.Contains(t.ID, s) || strings.Contains(t.VimID, s) {
			return &t, nil
		}
	}
	return nil, &TenantCloudNotFound{errMsg: s}
}

// FindCloudLink search for a cloud provider cloud link
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
func (t *Tenants) Filter(q TenantCloudFilter, f func(string) bool) ([]TenantsDetails, error) {

	filtered := make([]TenantsDetails, 0)
	for _, v := range t.TenantsList {
		if q == FilterVimType {
			if f(v.VimType) {
				filtered = append(filtered, v)
			}
		}
		if q == FilterHcxUUID {
			if f(v.HcxUUID) {
				filtered = append(filtered, v)
			}
		}
		if q == VimLocationCity {
			if f(v.Location.City) {
				filtered = append(filtered, v)
			}
		}
		if q == VimLocationCountry {
			if f(v.Location.Country) {
				filtered = append(filtered, v)
			}
		}
	}

	return filtered, nil
}
