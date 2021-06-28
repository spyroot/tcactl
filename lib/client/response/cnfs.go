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
	"time"
)

// CnfFilterType - cnf filter types
type CnfFilterType int32

const (
	// FilterCnfCID by cnf id
	FilterCnfCID CnfFilterType = 0

	// FilerVnfInstanceName filer by vnf instance name
	FilerVnfInstanceName CnfFilterType = 1

	// FilterVnfdID by vnfd id
	FilterVnfdID CnfFilterType = 2

	// FilterVnfCatalogName filters by vnf catalog name
	FilterVnfCatalogName CnfFilterType = 3
)

type CnfInstancesError struct {
	Type     string `json:"type" yaml:"type"`
	Title    string `json:"title" yaml:"title"`
	Status   int    `json:"status" yaml:"status"`
	Detail   string `json:"detail" yaml:"detail"`
	Instance string `json:"instance" yaml:"instance"`
}

type RespondID struct {
	Timestamp         int `json:"timestamp" yaml:"timestamp"`
	MachineIdentifier int `json:"machineIdentifier" yaml:"machineIdentifier"`
	ProcessIdentifier int `json:"processIdentifier" yaml:"processIdentifier"`
	Counter           int `json:"counter" yaml:"counter"`
}

//CnfMetadata Metadata information attached to respond for cnflcm req
type CnfMetadata struct {
	VnfPkgID          string                   `json:"vnfPkgId" yaml:"vnfPkgId"`
	VnfCatalogName    string                   `json:"vnfCatalogName" yaml:"vnfCatalogName"`
	ManagedBy         models.InternalManagedBy `json:"managedBy" yaml:"managedBy"`
	NfType            string                   `json:"nfType" yaml:"nfType"`
	LcmOperation      string                   `json:"lcmOperation" yaml:"lcmOperation"`
	LcmOperationState string                   `json:"lcmOperationState" yaml:"lcmOperationState"`
	IsUsedByNS        string                   `json:"isUsedByNS" yaml:"isUsedByNS"`
	AttachedNSCount   string                   `json:"attachedNSCount" yaml:"attachedNSCount"`
}

// VimConnectionInfo - Contains extra information including vim id , type
type VimConnectionInfo struct {
	Id            string `json:"id" yaml:"id"`
	VimId         string `json:"vimId" yaml:"vimId"`
	VimType       string `json:"vimType" yaml:"vimType"`
	InterfaceInfo struct {
	} `json:"interfaceInfo" yaml:"interfaceInfo"`
	AccessInfo struct {
	} `json:"accessInfo" yaml:"accessInfo"`
	Extra struct {
		DeploymentProfileId string `json:"deploymentProfileId" yaml:"deploymentProfileId"`
		NodeProfileName     string `json:"nodeProfileName" yaml:"nodeProfileName"`
		NodePoolId          string `json:"nodePoolId" yaml:"nodePoolId"`
		NodePoolName        string `json:"nodePoolName" yaml:"nodePoolName"`
		VimName             string `json:"vimName" yaml:"vimName"`
	} `json:"extra" yaml:"extra"`
}

// CnfInstantiateEntry VNFD Charts detail
type CnfInstantiateEntry struct {
	DispatchType       string `json:"dispatchType" yaml:"dispatch_type"`
	Namespace          string `json:"namespace" yaml:"namespace"`
	ChartName          string `json:"chartName" yaml:"chartName"`
	ChartVersion       string `json:"chartVersion" yaml:"chartVersion"`
	RepoURL            string `json:"repoUrl" yaml:"repoUrl"`
	Username           string `json:"username" yaml:"username"`
	Password           string `json:"password" yaml:"password"`
	HelmVersion        string `json:"helmVersion" yaml:"helmVersion"`
	VduID              string `json:"vduId" yaml:"vduId"`
	EntityID           string `json:"entityId" yaml:"entityId"`
	DeploymentName     string `json:"deploymentName" yaml:"deploymentName"`
	InstantiationState string `json:"instantiationState" yaml:"instantiationState"`
}

// CnfPolicyUri CNF policy URI
type CnfPolicyUri struct {
	Href string `json:"href,omitempty"`
}

type PolicyLinks struct {
	Self           CnfPolicyUri `json:"self,omitempty" yaml:"self"`
	Indicators     CnfPolicyUri `json:"indicators,omitempty" yaml:"indicators"`
	Instantiate    CnfPolicyUri `json:"instantiate,omitempty" yaml:"instantiate"`
	Terminate      CnfPolicyUri `json:"terminate,omitempty" yaml:"terminate"`
	Scale          CnfPolicyUri `json:"scale,omitempty" yaml:"scale"`
	ScaleToLevel   CnfPolicyUri `json:"scaleToLevel,omitempty" yaml:"scaleToLevel"`
	Heal           CnfPolicyUri `json:"heal,omitempty" yaml:"heal"`
	Update         CnfPolicyUri `json:"update,omitempty" yaml:"update"`
	UpgradePackage CnfPolicyUri `json:"upgrade_package,omitempty" yaml:"upgrade_package"`
	Upgrade        CnfPolicyUri `json:"upgrade,omitempty" yaml:"upgrade"`
	Reconfigure    CnfPolicyUri `json:"reconfigure,omitempty" yaml:"reconfigure"`
	ChangeFlavour  CnfPolicyUri `json:"changeFlavour,omitempty" yaml:"changeFlavour"`
	Operate        CnfPolicyUri `json:"operate,omitempty" yaml:"operate"`
	ChangeExtConn  CnfPolicyUri `json:"changeExtConn,omitempty" yaml:"changeExtConn"`
}

type CnfLcmExtended struct {
	RespId                 RespondID                      `json:"_id" yaml:"resp_id"`
	CID                    string                         `json:"id" yaml:"cid"`
	VnfInstanceName        string                         `json:"vnfInstanceName" yaml:"vnfInstanceName"`
	VnfInstanceDescription string                         `json:"vnfInstanceDescription" yaml:"vnfInstanceDescription"`
	VnfdID                 string                         `json:"vnfdId" yaml:"vnfdId"`
	VnfPkgID               string                         `json:"vnfPkgId" yaml:"vnfPkgId"`
	VnfCatalogName         string                         `json:"vnfCatalogName" yaml:"vnfCatalogName"`
	VnfProvider            string                         `json:"vnfProvider" yaml:"vnfProvider"`
	VnfProductName         string                         `json:"vnfProductName" yaml:"vnfProductName"`
	VnfSoftwareVersion     string                         `json:"vnfSoftwareVersion" yaml:"vnfSoftwareVersion"`
	VnfdVersion            string                         `json:"vnfdVersion" yaml:"vnfdVersion"`
	OnboardedVnfPkgInfoID  string                         `json:"onboardedVnfPkgInfoId" yaml:"onboardedVnfPkgInfoId"`
	InstantiationState     string                         `json:"instantiationState" yaml:"instantiationState"`
	ManagedBy              *models.InternalManagedBy      `json:"managedBy,omitempty" yaml:"managedBy,omitempty"`
	NfType                 string                         `json:"nfType" yaml:"nf_type"`
	Links                  PolicyLinks                    `json:"_links" yaml:"links"`
	LastUpdated            time.Time                      `json:"lastUpdated" yaml:"lastUpdated"`
	LastUpdateEnterprise   string                         `json:"lastUpdateEnterprise" yaml:"lastUpdateEnterprise"`
	LastUpdateOrganization string                         `json:"lastUpdateOrganization" yaml:"lastUpdateOrganization"`
	LastUpdateUser         string                         `json:"lastUpdateUser" yaml:"lastUpdateUser"`
	CreationDate           time.Time                      `json:"creationDate" yaml:"creationDate"`
	CreationEnterprise     string                         `json:"creationEnterprise" yaml:"creationEnterprise"`
	CreationOrganization   string                         `json:"creationOrganization" yaml:"creationOrganization"`
	CreationUser           string                         `json:"creationUser" yaml:"creationUser"`
	IsDeleted              bool                           `json:"isDeleted" yaml:"isDeleted"`
	VimConnectionInfo      []VimConnectionInfo            `json:"vimConnectionInfo" yaml:"vimConnectionInfo"`
	LcmOperation           string                         `json:"lcmOperation" yaml:"lcmOperation"`
	LcmOperationState      string                         `json:"lcmOperationState" yaml:"lcmOperationState"`
	RowType                string                         `json:"rowType" yaml:"rowType"`
	InstantiatedNfInfo     map[string]CnfInstantiateEntry `json:"instantiatedNfInfo,omitempty" yaml:"instantiatedNfInfo"`
	InstantiatedVnfInfo    map[string]CnfInstantiateEntry `json:"instantiatedVnfInfo,omitempty" yaml:"instantiatedVnfInfo"`
	IsUsedByNS             bool                           `json:"isUsedByNS" yaml:"isUsedByNS"`
	AttachedNSCount        int                            `json:"attachedNSCount" yaml:"attachedNSCount"`
	Meta                   CnfMetadata                    `json:"metadata,omitempty" yaml:"meta"`
}

// GetField - return struct field value
func (t *CnfLcmExtended) GetField(field string) string {

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
func (t *CnfLcmExtended) GetFields() (map[string]interface{}, error) {

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

//IsInCluster return true if cnf in cluster indicated
//vimName.
func (e *CnfLcmExtended) IsInCluster(vimName string) bool {

	if e == nil {
		return false
	}

	if e.VimConnectionInfo == nil {
		return false
	}

	for _, info := range e.VimConnectionInfo {
		if strings.ToLower(info.Extra.VimName) == strings.ToLower(vimName) {
			return true
		}
		fmt.Println(info.Extra.VimName, vimName)
	}
	return false
}

// CnfsExtended - list of CNF LCM respond
type CnfsExtended struct {
	CnfLcms []CnfLcmExtended
}

// Filter filters respond based on filter type and pass to callback
func (c *CnfsExtended) Filter(q CnfFilterType, f func(string) bool) ([]CnfLcmExtended, error) {

	if c == nil {
		return nil, fmt.Errorf("cnfs instance is nil")
	}

	filtered := make([]CnfLcmExtended, 0)
	for _, cnf := range c.CnfLcms {
		if q == FilterCnfCID && f(cnf.CID) {
			filtered = append(filtered, cnf)
		}
		if q == FilerVnfInstanceName && f(cnf.VnfInstanceName) {
			filtered = append(filtered, cnf)
		}
		if q == FilterVnfdID && f(cnf.VnfdID) {
			filtered = append(filtered, cnf)
		}
		if q == FilterVnfCatalogName && f(cnf.VnfCatalogName) {
			filtered = append(filtered, cnf)
		}
	}

	return filtered, nil
}

// CnfNotFound error raised if cnf not found
type CnfNotFound struct {
	errMsg string
}

//
func (m *CnfNotFound) Error() string {
	return "cnf '" + m.errMsg + "' not found"
}

// FindByName - tries to find CNF by product name, id.
func (c *CnfsExtended) FindByName(s string) (*CnfLcmExtended, error) {

	if c == nil {
		return nil, fmt.Errorf("cnfs instance is nil")
	}

	for _, cnf := range c.CnfLcms {
		if strings.HasPrefix(cnf.VnfProductName, s) ||
			strings.HasPrefix(cnf.VnfProductName, s) ||
			strings.HasPrefix(cnf.CID, s) {
			return &cnf, nil
		}
	}

	return nil, &CnfNotFound{s}
}

// ResolveFromName - tries to find CNF by product name or id.
func (c *CnfsExtended) ResolveFromName(s string) (*CnfLcmExtended, error) {

	if c == nil {
		return nil, fmt.Errorf("cnfs instance is nil")
	}

	for _, cnf := range c.CnfLcms {
		if strings.HasPrefix(cnf.VnfInstanceName, s) ||
			strings.HasPrefix(cnf.CID, s) {
			return &cnf, nil
		}
	}

	return nil, &CnfNotFound{s}
}
