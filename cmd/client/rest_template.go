package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/client/response"
	"net/http"
	"strings"
)

const (
	apiClusterTemplates = "/hybridity/api/infra/cluster-templates"
	apiClusterTemplate  = "/hybridity/api/infra/cluster-template"
)

// GetClusterTemplates - return list of cluster templates
func (c *RestClient) GetClusterTemplates() (*response.ClusterTemplates, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + apiClusterTemplates)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			return nil, fmt.Errorf(errRes.Message)
		}
		return nil, fmt.Errorf("unknown error, status code: %v ", resp.StatusCode())
	}

	var template response.ClusterTemplates
	if err := json.Unmarshal(resp.Body(), &template.ClusterTemplates); err != nil {
		return nil, err
	}

	return &template, nil
}

// CreateClusterTemplate - creates cluster template from specs
func (c *RestClient) CreateClusterTemplate(spec *response.ClusterTemplate) error {

	c.GetClient()
	resp, err := c.Client.R().SetBody(spec).Post(c.BaseURL + apiClusterTemplates)
	if err != nil {
		glog.Error(err)
		return err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errResp ErrorsResponse
		if err = json.Unmarshal(resp.Body(), &errResp); err == nil {
			errs := errResp.GetErrors()
			return fmt.Errorf(strings.Join(errs, "\n"))
		} else {
			glog.Errorf("Failed to parse error respond. %v", err)
		}
		return fmt.Errorf("unknown error, status code: %v ", resp.StatusCode())
	}

	if resp.StatusCode() == http.StatusOK {
		glog.Infof("Template created.")
	}

	return nil
}

// UpdateClusterTemplate - updates existing cluster template
func (c *RestClient) UpdateClusterTemplate(spec *response.ClusterTemplate) error {

	c.GetClient()
	resp, err := c.Client.R().SetBody(spec).Put(c.BaseURL + apiClusterTemplates + "/" + spec.Id)
	if err != nil {
		glog.Error(err)
		return err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errResp ErrorsResponse
		if err = json.Unmarshal(resp.Body(), &errResp); err == nil {
			errs := errResp.GetErrors()
			return fmt.Errorf(strings.Join(errs, "\n"))
		} else {
			glog.Errorf("Failed to parse error respond. %v", err)
		}
		return fmt.Errorf("unknown error, status code: %v ", resp.StatusCode())
	}

	if resp.StatusCode() == http.StatusOK {
		glog.Infof("Template updated.")
	}

	return nil
}

// GetClusterTemplate - return list of cluster templates
func (c *RestClient) GetClusterTemplate(clusterId string) (*response.ClusterTemplate, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + apiClusterTemplates + "/" + clusterId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	// it doesn't responds with not found or proper payload.
	if resp.StatusCode() == http.StatusInternalServerError {
		return nil, fmt.Errorf("template not found")
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var template response.ClusterTemplate
	if err := json.Unmarshal(resp.Body(), &template); err != nil {
		return nil, err
	}

	return &template, nil
}

// DeleteClusterTemplate - deletes cluster template
func (c *RestClient) DeleteClusterTemplate(clusterId string) error {

	c.GetClient()
	resp, err := c.Client.R().Delete(c.BaseURL + apiClusterTemplates + "/" + clusterId)
	if err != nil {
		glog.Error(err)
		return err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			return fmt.Errorf(errRes.Message)
		}
		return fmt.Errorf("unknown error, status code: %v ", resp.StatusCode())
	}

	if resp.StatusCode() == http.StatusOK {
		glog.Infof("Template deleted.")
	}

	return nil
}
