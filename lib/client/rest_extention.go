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
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	"net/http"
)

// GetExtensions - api call for all extension.
func (c *RestClient) GetExtensions() (*response.Extensions, error) {

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
		glog.Error("Failed parse server respond %v.", err)
		return nil, err
	}

	return &e, nil
}

//CreateExtension - add new extension
func (c *RestClient) CreateExtension(spec *request.ExtensionSpec) (string, error) {

	if c == nil {
		return "", fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()

	resp, err := c.Client.R().
		SetBody(spec).
		Post(c.BaseURL + TcaVmwareExtensions)

	if err != nil {
		glog.Error(err)
		return "", err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return "", c.checkError(resp)
	}

	var ext response.Extension
	if err := json.Unmarshal(resp.Body(), &ext); err != nil {
		glog.Error("Failed parse server respond. %v", err)
		return "", err
	}

	return ext.ExtensionId, nil
}

//DeleteExtension - query repositories linked to vim
func (c *RestClient) DeleteExtension(extensionId string) (*models.TcaTask, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized rest client")
	}

	glog.Infof("Deleting extension %v", extensionId)

	c.GetClient()
	resp, err := c.Client.R().Delete(c.BaseURL + fmt.Sprintf(TcaVmwareDeleteExtensions, extensionId))

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkError(resp)
	}

	var task models.TcaTask
	if err := json.Unmarshal(resp.Body(), &task); err != nil {
		glog.Error("Failed parse server respond. %v", err)
		return nil, err
	}

	return &task, nil
}

//UpdateExtension - update extension.
func (c *RestClient) UpdateExtension(spec *request.ExtensionSpec) (*response.ReposList, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()

	resp, err := c.Client.R().
		SetBody(spec).
		Put(c.BaseURL + TcaVmwareExtensions)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkError(resp)
	}

	var repos response.ReposList
	if err := json.Unmarshal(resp.Body(), &repos); err != nil {
		glog.Error("Failed parse server respond. %v", err)
		return nil, err
	}

	return &repos, nil
}
