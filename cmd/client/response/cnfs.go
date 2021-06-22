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
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

type RespondID struct {
	Timestamp         int `json:"timestamp"`
	MachineIdentifier int `json:"machineIdentifier"`
	ProcessIdentifier int `json:"processIdentifier"`
	Counter           int `json:"counter"`
}

type InternalManagedBy struct {
	ExtensionSubtype string `json:"extensionSubtype"`
	ExtensionName    string `json:"extensionName"`
}

//CnfMetadata Metadata information attached to respond for cnflcm req
type CnfMetadata struct {
	VnfPkgID          string            `json:"vnfPkgId"`
	VnfCatalogName    string            `json:"vnfCatalogName"`
	ManagedBy         InternalManagedBy `json:"managedBy"`
	NfType            string            `json:"nfType"`
	LcmOperation      string            `json:"lcmOperation"`
	LcmOperationState string            `json:"lcmOperationState"`
	IsUsedByNS        string            `json:"isUsedByNS"`
	AttachedNSCount   string            `json:"attachedNSCount"`
}

// VimConnectionInfo - Contains extra information including vim id , type
type VimConnectionInfo struct {
	Id            string `json:"id"`
	VimId         string `json:"vimId"`
	VimType       string `json:"vimType"`
	InterfaceInfo struct {
	} `json:"interfaceInfo"`
	AccessInfo struct {
	} `json:"accessInfo"`
	Extra struct {
		DeploymentProfileId string `json:"deploymentProfileId"`
		NodeProfileName     string `json:"nodeProfileName"`
		NodePoolId          string `json:"nodePoolId"`
		NodePoolName        string `json:"nodePoolName"`
		VimName             string `json:"vimName"`
	} `json:"extra"`
}

// CnfInstantiateEntry VNFD Charts detail
type CnfInstantiateEntry struct {
	DispatchType       string `json:"dispatchType"`
	Namespace          string `json:"namespace"`
	ChartName          string `json:"chartName"`
	ChartVersion       string `json:"chartVersion"`
	RepoURL            string `json:"repoUrl"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	HelmVersion        string `json:"helmVersion"`
	VduID              string `json:"vduId"`
	EntityID           string `json:"entityId"`
	DeploymentName     string `json:"deploymentName"`
	InstantiationState string `json:"instantiationState"`
}

// CnfPolicyUri CNF policy URI
type CnfPolicyUri struct {
	Href string `json:"href,omitempty"`
}

type PolicyLinks struct {
	Self           CnfPolicyUri `json:"self,omitempty"`
	Indicators     CnfPolicyUri `json:"indicators,omitempty"`
	Instantiate    CnfPolicyUri `json:"instantiate,omitempty"`
	Terminate      CnfPolicyUri `json:"terminate,omitempty"`
	Scale          CnfPolicyUri `json:"scale,omitempty"`
	ScaleToLevel   CnfPolicyUri `json:"scaleToLevel,omitempty"`
	Heal           CnfPolicyUri `json:"heal,omitempty"`
	Update         CnfPolicyUri `json:"update,omitempty"`
	UpgradePackage CnfPolicyUri `json:"upgrade_package,omitempty"`
	Upgrade        CnfPolicyUri `json:"upgrade,omitempty"`
	Reconfigure    CnfPolicyUri `json:"reconfigure,omitempty"`
	ChangeFlavour  CnfPolicyUri `json:"changeFlavour,omitempty"`
	Operate        CnfPolicyUri `json:"operate,omitempty"`
	ChangeExtConn  CnfPolicyUri `json:"changeExtConn,omitempty"`
}

type CnfLcmExtended struct {
	RespId                 RespondID                      `json:"_id"`
	CID                    string                         `json:"id"`
	VnfInstanceName        string                         `json:"vnfInstanceName"`
	VnfInstanceDescription string                         `json:"vnfInstanceDescription"`
	VnfdID                 string                         `json:"vnfdId"`
	VnfPkgID               string                         `json:"vnfPkgId"`
	VnfCatalogName         string                         `json:"vnfCatalogName"`
	VnfProvider            string                         `json:"vnfProvider"`
	VnfProductName         string                         `json:"vnfProductName"`
	VnfSoftwareVersion     string                         `json:"vnfSoftwareVersion"`
	VnfdVersion            string                         `json:"vnfdVersion"`
	OnboardedVnfPkgInfoID  string                         `json:"onboardedVnfPkgInfoId"`
	InstantiationState     string                         `json:"instantiationState"`
	ManagedBy              InternalManagedBy              `json:"managedBy"`
	NfType                 string                         `json:"nfType"`
	Links                  PolicyLinks                    `json:"_links"`
	LastUpdated            time.Time                      `json:"lastUpdated"`
	LastUpdateEnterprise   string                         `json:"lastUpdateEnterprise"`
	LastUpdateOrganization string                         `json:"lastUpdateOrganization"`
	LastUpdateUser         string                         `json:"lastUpdateUser"`
	CreationDate           time.Time                      `json:"creationDate"`
	CreationEnterprise     string                         `json:"creationEnterprise"`
	CreationOrganization   string                         `json:"creationOrganization"`
	CreationUser           string                         `json:"creationUser"`
	IsDeleted              bool                           `json:"isDeleted"`
	VimConnectionInfo      []VimConnectionInfo            `json:"vimConnectionInfo"`
	LcmOperation           string                         `json:"lcmOperation"`
	LcmOperationState      string                         `json:"lcmOperationState"`
	RowType                string                         `json:"rowType"`
	InstantiatedNfInfo     map[string]CnfInstantiateEntry `json:"instantiatedNfInfo,omitempty"`
	InstantiatedVnfInfo    map[string]CnfInstantiateEntry `json:"instantiatedVnfInfo,omitempty"`
	IsUsedByNS             bool                           `json:"isUsedByNS"`
	AttachedNSCount        int                            `json:"attachedNSCount"`
	Meta                   CnfMetadata                    `json:"metadata,omitempty"`
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

	return nil, fmt.Errorf("CNF not found")
}

// ResolveFromName - tries to find CNF by product name, id.
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

	return nil, fmt.Errorf("CNF not found")
}
