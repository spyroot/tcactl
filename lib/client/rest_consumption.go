package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/models"
	"net/http"
)

func (c *RestClient) GetConsumption(ctx context.Context) (*models.ConsumptionResp, error) {

	var (
		resp *resty.Response
		err  error
	)

	c.GetClient()
	resp, err = c.Client.R().SetContext(ctx).Get(c.BaseURL + TcaConsumption)

	fmt.Println("##### Sending Request to", c.BaseURL+TcaConsumption)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return nil, c.checkError(resp)
	}

	var apiResp models.ConsumptionResp
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}
	return &apiResp, nil
}
