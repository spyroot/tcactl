package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/models"
	"net/http"
)

const (
	// TcaVmwareClusters VMware clusters
	TcaVmwareClusters = "/hybridity/api/infra/inventory/vc/clusters"
	// TcaVmwareVmTemplates VMware VM templates
	TcaVmwareVmTemplates = "/https://tca-vip03.cnfdemo.io/hybridity/api/infra/inventory/vc/templates"
	// TcaVmwareNetworks VMware virtual networks
	TcaVmwareNetworks = "/hybridity/api/nfv/networks"
)

//
//// GetVcClusters return list of cloud provider attached to TCA
//func (c *RestClient) GetVcClusters() (*response.Tenants, error) {
//
//	c.GetClient()
//	resp, err := c.Client.R().Get(c.BaseURL + apiTenants)
//	if err != nil {
//		glog.Error(err)
//		return nil, err
//	}
//
//	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
//		return nil, fmt.Errorf("template not found")
//	}
//
//	var tenants response.Tenants
//	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
//		return nil, err
//	}
//
//	return &tenants, nil
//}

func (c *RestClient) GetVmwareCluster(f *request.VmwareClusterQuery) (*models.VMwareClusters, error) {

	c.GetClient()
	resp, err := c.Client.R().SetBody(f).Post(c.BaseURL + TcaVmwareClusters)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	fmt.Println(string(resp.Body()))

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var tenants models.VMwareClusters
	if err := json.Unmarshal(resp.Body(), &tenants); err != nil {
		return nil, err
	}

	return &tenants, nil
}

// GetVmwareNetworks - return query for vmware network list
func (c *RestClient) GetVmwareNetworks(f *request.VMwareNetworkQuery) (*models.CloudNetworks, error) {

	c.GetClient()

	resp, err := c.Client.R().SetBody(f).Post(c.BaseURL + TcaVmwareNetworks)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var networks models.CloudNetworks
	if err := json.Unmarshal(resp.Body(), &networks); err != nil {
		return nil, err
	}

	return &networks, nil
}

// GetVMwareTemplates - return VMware VM templates
// Typically Query filters based on tenant id.
func (c *RestClient) GetVMwareTemplates(f *request.VMwareTemplateQuery) (*models.VcInventory, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + TcaVmwareVmTemplates)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var inventory models.VcInventory
	if err := json.Unmarshal(resp.Body(), &inventory); err != nil {
		return nil, err
	}

	return &inventory, nil
}
