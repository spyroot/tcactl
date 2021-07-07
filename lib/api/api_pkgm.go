package api

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/response"
)

// GetEntireCatalog - return all CNF/VNF packages
// onboarded in TCA.
func (a *TcaApi) GetEntireCatalog() (*response.VnfPackages, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetAllCatalog()
}

// GetVnfPkgm - return packages
// pkgId is id of package in TCA
// filter
func (a *TcaApi) GetVnfPkgm(filter string, pkgId string) (*response.VnfPackages, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetVnfPkgm(filter, pkgId)
}

// GetCatalogId return vnf Package ID and VNFD ID
func (a *TcaApi) GetCatalogId(catalogId string) (string, string, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return "", "", fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetPackageCatalogId(catalogId)
}

// GetCatalogAndVdu API method returns
// catalog entity and vdu package.
func (a *TcaApi) GetCatalogAndVdu(nfdName string) (*response.VnfPackage, *response.VduPackage, error) {

	vnfCatalog, err := a.rest.GetVnfPkgm("", "")
	if err != nil || vnfCatalog == nil {
		glog.Errorf("Failed acquire vnf package information. Error %v", err)
		return nil, nil, err
	}

	catalogEntity, err := vnfCatalog.GetVnfdID(nfdName)
	if err != nil || catalogEntity == nil {
		glog.Errorf("Failed acquire catalog information for catalog name %v", nfdName)
		return nil, nil, err
	}

	v, err := a.rest.GetVnfPkgmVnfd(catalogEntity.PID)

	return catalogEntity, v, err
}
