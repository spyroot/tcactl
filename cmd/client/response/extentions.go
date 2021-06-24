package response

import (
	"fmt"
	"strings"
)

type ExtensionInterfaceInfo struct {
	Url                string `json:"url" yaml:"url"`
	Description        string `json:"description" yaml:"description"`
	TrustedCertificate string `json:"trustedCertificate" yaml:"trustedCertificate"`
}

type ExtensionAccessInfo struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type ExtensionAdditionalParameters struct {
	TrustAllCerts bool `json:"trustAllCerts" yaml:"trustAllCerts"`
	RepoSyncVim   struct {
		VimName       string `json:"vimName" yaml:"vimName"`
		VimId         string `json:"vimId" yaml:"vim_id"`
		VimSystemUUID string `json:"vimSystemUUID" yaml:"vim_system_uuid"`
	} `json:"repoSyncVim" yaml:"repo_sync_vim"`
}

type Extension struct {
	ExtensionId          string                         `json:"extensionId" yaml:"extensionId"`
	Name                 string                         `json:"name" yaml:"name"`
	Type                 string                         `json:"type" yaml:"type"`
	ExtensionKey         string                         `json:"extensionKey" yaml:"extensionKey"`
	Description          string                         `json:"description" yaml:"description"`
	InterfaceInfo        *ExtensionInterfaceInfo        `json:"interfaceInfo" yaml:"interfaceInfo"`
	AccessInfo           *ExtensionAccessInfo           `json:"accessInfo" yaml:"accessInfo"`
	AdditionalParameters *ExtensionAdditionalParameters `json:"additionalParameters" yaml:"additional_parameters"`
	State                string                         `json:"state" yaml:"state"`
	ExtensionSubtype     string                         `json:"extensionSubtype" yaml:"extension_subtype"`
	Products             []interface{}                  `json:"products" yaml:"products"`
	VimInfo              []struct {
		VimName       string `json:"vimName" yaml:"vimName"`
		VimId         string `json:"vimId" yaml:"vimId"`
		VimSystemUUID string `json:"vimSystemUUID" yaml:"vimSystemUUID"`
	} `json:"vimInfo" yaml:"vim_info"`
	Version          string `json:"version" yaml:"version"`
	VnfCount         int    `json:"vnfCount" yaml:"vnf_count"`
	VnfCatalogCount  int    `json:"vnfCatalogCount" yaml:"vnfCatalogCount"`
	Error            string `json:"error" yaml:"error"`
	AutoScaleEnabled bool   `json:"autoScaleEnabled" yaml:"autoScaleEnabled"`
	AutoHealEnabled  bool   `json:"autoHealEnabled" yaml:"autoHealEnabled"`
}

func (e *Extension) IsEnabled() bool {
	return strings.ToLower(e.State) == "enabled"
}

type Extensions struct {
	ExtensionsList []Extension `json:"extensions"`
}

// FindRepo
func (e *Extensions) FindRepo(q string) (*Extension, error) {

	if e == nil {
		return nil, fmt.Errorf("uninitialized object")
	}

	for _, l := range e.ExtensionsList {
		if strings.Contains(l.Name, q) || strings.Contains(l.ExtensionId, q) {
			return &l, nil
		}
	}

	return nil, fmt.Errorf("repository not found")
}
