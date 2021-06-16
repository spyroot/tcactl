package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/client/response"
	"net/http"
)

func (c *RestClient) ExtensionQuery() (*response.Extensions, error) {

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + "/hybridity/api/extensions")
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes response.VnfPackagesError
		if err = json.NewDecoder(resp.RawResponse.Body).Decode(&errRes); err == nil {
			return nil, fmt.Errorf(errRes.Detail)
		}
		return nil, fmt.Errorf("unknown error, status code: %v", resp.Status())
	}

	var e response.Extensions
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		return nil, err
	}

	return &e, nil
}
