package client

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/models"
	"net/http"
)

// GetInfrNetworks - return list of cluster templates
func (c *RestClient) GetInfrNetworks(clusterId string) (*models.CloudNetworks, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + "hybridity/api/nfv/networks")
	if err != nil {
		glog.Error(err)
		return nil, err
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
