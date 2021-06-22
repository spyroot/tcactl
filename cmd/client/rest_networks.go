package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/models"
	"net/http"
)

// GetInfrNetworks - return list of cluster templates
// TODO
func (c *RestClient) GetInfrNetworks(tenantId string) (*models.CloudNetworks, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + TcaVmwareNfvNetworks + "/" + tenantId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.dumpRespond {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var net models.CloudNetworks
	if err := json.Unmarshal(resp.Body(), &net); err != nil {
		return nil, err
	}

	return &net, nil
}
