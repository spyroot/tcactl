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
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
)

type ExtensionDeleteReplay struct {
	ExtensionId string `json:"extensionId" yaml:"extensionId"`
	Deleted     bool   `json:"deleted" yaml:"deleted"`
}

type ExtensionUpdateReplay struct {
	ExtensionId string `json:"extensionId" yaml:"extensionId"`
	Updated     bool   `json:"updated" yaml:"updated"`
}

// GetExtensions - api call retrieves all extension.
// client cna filter on particular field, VIM Id , type etc.
func (c *RestClient) GetExtensions(ctx context.Context) (*response.Extensions, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).Get(c.BaseURL + TcaVmwareExtensions)
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

	var e response.Extensions
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		glog.Error("Failed parse server respond %v.", err)
		return nil, err
	}

	return &e, nil
}

// GetExtension - api call retrieves extension,
// eid is internal extension id.
func (c *RestClient) GetExtension(ctx context.Context, eid string) (*response.Extensions, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()

	resp, err := c.Client.R().SetContext(ctx).Get(c.BaseURL + fmt.Sprintf(TcaVmwareGetExtensions, eid))
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

	var e response.Extensions
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		glog.Error("Failed parse server respond %v.", err)
		return nil, err
	}

	return &e, nil
}

//CreateExtension - method creates new extension
// spec can contain optional VimInfo that indicates
// cluster or cluster where extension will be attach.
func (c *RestClient) CreateExtension(ctx context.Context, spec *specs.SpecExtension) (string, error) {

	if c == nil {
		return "", fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()

	resp, err := c.Client.R().SetContext(ctx).
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

//DeleteExtension - api call delete extension
//method return true if extension deleted.
func (c *RestClient) DeleteExtension(ctx context.Context, extensionId string) (bool, error) {

	if c == nil {
		return false, fmt.Errorf("uninitialized rest client")
	}

	req := fmt.Sprintf(TcaVmwareDeleteExtensions, extensionId)
	glog.Infof("Deleting extension %v, req", extensionId, req)

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).Delete(c.BaseURL + req)

	if err != nil {
		glog.Error(err)
		return false, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return false, c.checkError(resp)
	}

	var r ExtensionDeleteReplay
	if err := json.Unmarshal(resp.Body(), &r); err != nil {
		glog.Error("Failed parse server respond. %v", err)
		return false, err
	}

	return r.Deleted, nil
}

//UpdateExtension - update extension
//Returns true if extension updated
func (c *RestClient) UpdateExtension(ctx context.Context, spec *specs.SpecExtension, eid string) (bool, error) {

	if c == nil {
		return false, fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()
	resp, err := c.Client.R().
		SetContext(ctx).SetBody(spec).Post(c.BaseURL + fmt.Sprintf(TcaVmwareUpdateExtensions, eid))

	if err != nil {
		glog.Error(err)
		return false, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return false, c.checkError(resp)
	}

	var replay ExtensionUpdateReplay
	if err := json.Unmarshal(resp.Body(), &replay); err != nil {
		glog.Error("Failed parse server respond. %v", err)
		return false, err
	}

	return replay.Updated, nil
}
