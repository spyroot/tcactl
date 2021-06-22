// Package client
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Mustafa mbayramo@vmware.com
package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/models"
	"net/http"
)

// GetInfraNetworks - return list of cluster templates
// TODO
func (c *RestClient) GetInfraNetworks(tenantId string) (*models.CloudNetworks, error) {

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
