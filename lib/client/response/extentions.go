package response

import (
	"github.com/spyroot/tcactl/lib/api_errors"
	errorsconst "github.com/spyroot/tcactl/pkg/errors"
	"strings"
)

const (
	// ExtensionTypeRepository Type Repository extension
	ExtensionTypeRepository = "Repository"

	// ExtensionEnabled extension state enabled
	ExtensionEnabled = "ENABLED"
)

// ExtensionInterfaceInfo contains url for harbor , cert
type ExtensionInterfaceInfo struct {
	Url                string `json:"url" yaml:"url"`
	Description        string `json:"description" yaml:"description"`
	TrustedCertificate string `json:"trustedCertificate" yaml:"trustedCertificate"`
}

// ExtensionAccessInfo for Harbor it username and password
type ExtensionAccessInfo struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

// ExtensionAdditionalParameters contains data where it attached to.
type ExtensionAdditionalParameters struct {
	TrustAllCerts bool `json:"trustAllCerts" yaml:"trustAllCerts"`
	RepoSyncVim   struct {
		VimName       string `json:"vimName" yaml:"vimName"`
		VimId         string `json:"vimId" yaml:"vim_id"`
		VimSystemUUID string `json:"vimSystemUUID" yaml:"vim_system_uuid"`
	} `json:"repoSyncVim" yaml:"repo_sync_vim"`
}

type ExtensionVimInfo struct {
	VimName       string `json:"vimName" yaml:"vimName"`
	VimId         string `json:"vimId" yaml:"vimId"`
	VimSystemUUID string `json:"vimSystemUUID" yaml:"vimSystemUUID"`
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
	VimInfo              []ExtensionVimInfo             `json:"vimInfo" yaml:"vim_info"`
	Version              string                         `json:"version" yaml:"version"`
	VnfCount             int                            `json:"vnfCount" yaml:"vnf_count"`
	VnfCatalogCount      int                            `json:"vnfCatalogCount" yaml:"vnfCatalogCount"`
	Error                string                         `json:"error" yaml:"error"`
	AutoScaleEnabled     bool                           `json:"autoScaleEnabled" yaml:"autoScaleEnabled"`
	AutoHealEnabled      bool                           `json:"autoHealEnabled" yaml:"autoHealEnabled"`
}

func (e *Extension) IsEnabled() bool {
	return strings.ToLower(e.State) == "enabled"
}

type Extensions struct {
	ExtensionsList []Extension `json:"extensions"`
}

// FindRepo find repository Extension
func (e *Extensions) FindRepo(q string) (*Extension, error) {

	if e == nil {
		return nil, errorsconst.NilError
	}

	for _, l := range e.ExtensionsList {
		if strings.Contains(l.Name, q) || strings.Contains(l.ExtensionId, q) {
			return &l, nil
		}
	}

	return nil, api_errors.NewExtensionsNotFound(q)
}

// FindExtension find repository Extension by name or extension id
func (e *Extensions) FindExtension(NameOrId string) (*Extension, error) {

	if e == nil {
		return nil, errorsconst.NilError
	}

	q := strings.ToLower(NameOrId)

	for _, l := range e.ExtensionsList {
		if strings.ToLower(l.Name) == q || l.ExtensionId == q {
			return &l, nil
		}
	}

	return nil, api_errors.NewExtensionsNotFound(q)
}

// GetAllRepositories return all repository Extension
func (e *Extensions) GetAllRepositories() (*Extension, error) {

	if e == nil {
		return nil, errorsconst.NilError
	}

	for _, l := range e.ExtensionsList {
		if strings.Contains(l.Type, ExtensionTypeRepository) && strings.Contains(l.State, ExtensionEnabled) {
			return &l, nil
		}
	}

	return nil, api_errors.NewExtensionsNotFound(ExtensionTypeRepository)
}

// GetRepositoryByUrl return all repository by url
func (e *Extensions) GetRepositoryByUrl(url string) (*Extension, error) {

	if e == nil {
		return nil, errorsconst.NilError
	}

	for _, l := range e.ExtensionsList {
		if l.InterfaceInfo != nil && strings.Contains(l.InterfaceInfo.Url, url) {
			return &l, nil
		}
	}

	return nil, api_errors.NewExtensionsNotFound(url)
}

// GetVimAttached return list all vim extension ext
// attached to.  ext can be a name or id.
func (e *Extensions) GetVimAttached(ext string, vs ...string) ([]ExtensionVimInfo, error) {

	if e == nil {
		return nil, errorsconst.NilError
	}

	v := ""
	if len(vs) == 1 {
		v = vs[0]
	}

	var vimInfo []ExtensionVimInfo
	for _, l := range e.ExtensionsList {
		if strings.Contains(l.Name, ext) || strings.Contains(l.ExtensionId, ext) {
			for _, info := range l.VimInfo {
				if len(v) > 0 {
					if strings.Contains(info.VimId, v) ||
						strings.Contains(info.VimSystemUUID, v) ||
						strings.Contains(info.VimId, v) {
						vimInfo = append(vimInfo, info)
					}
				} else {
					vimInfo = append(vimInfo, info)
				}
			}
		}
	}

	return vimInfo, nil
}
