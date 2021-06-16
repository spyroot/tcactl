package client

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/models"
	"net/http"
)

// GetInfraContainers - return list of cluster templates
func (c *RestClient) GetInfraContainers(clusterId string) (*models.VcInfra, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + " /api/service/inventory/containers")
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var vcInfra models.VcInfra
	if err := json.Unmarshal(resp.Body(), &vcInfra); err != nil {
		return nil, err
	}

	return &vcInfra, nil
}
