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
	"github.com/spyroot/tcactl/lib/client/response"
	"net/http"
)

// ExtensionQuery - query for all extension api
func (c *RestClient) ExtensionQuery() (*response.Extensions, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()
	resp, err := c.Client.R().Get(c.BaseURL + TcaVmwareExtensions)
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

	var e response.Extensions
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &e, nil
}
