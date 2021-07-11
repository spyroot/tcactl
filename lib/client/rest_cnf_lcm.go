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
	"github.com/go-resty/resty/v2"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
)

// GetVnflcm - Retrieves information about a CNF/VNF instance by reading
// an "Individual VNF instance" resource.
//
// Example of filter
// (eq,id,5c11bd9c-085d-4913-a453-572457ddffe2)
func (c *RestClient) GetVnflcm(req ...string) (interface{}, error) {

	var (
		resp    *resty.Response
		err     error
		isArray = true
	)

	c.GetClient()

	// no args will return entire list
	if len(req) == 0 {
		resp, err = c.Client.R().Get(c.BaseURL + TcaApiVnfLcmExtensionVnfInstance)
	}
	// attach filter and dispatch
	if len(req) == 1 {
		var queryFilter = req[0]
		resp, err = c.Client.R().
			SetQueryParams(map[string]string{"filter": queryFilter}).
			Get(c.BaseURL + TcaApiVnfLcmExtensionVnfInstance)
	}
	// if two argument will retrieve particular catalog entity
	if len(req) == 2 {
		isArray = false
		resp, err = c.Client.R().Get(c.BaseURL + TcaApiVnfLcmVnfInstance + "/" + req[1])
	}

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	//
	if !resp.IsSuccess() {
		var errRes response.CnfInstancesError
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			return nil, fmt.Errorf(errRes.Detail)
		}
		return nil, fmt.Errorf("unknown error, status code: %d", resp.StatusCode())
	}

	// for single cnf request, pack return result in array.
	if !isArray {
		var (
			cnflcm  response.LcmInfo
			cnfslcm response.Cnfs
		)

		if err := json.Unmarshal(resp.Body(), &cnflcm); err != nil {
			glog.Errorf("Failed parse servers respond. %v", err)
			return nil, err
		}

		cnfslcm.CnfLcms = append(cnfslcm.CnfLcms, cnflcm)
		return &cnfslcm, nil
	}

	// default case
	var extended response.CnfsExtended
	if err := json.Unmarshal(resp.Body(), &extended.CnfLcms); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	return &extended, nil
}

// GetRunningVnflcm rest call return state of CNF or VNF
// state described as response.LcmInfo
func (c *RestClient) GetRunningVnflcm(instanceId string) (*response.LcmInfo, error) {

	var (
		resp *resty.Response
		err  error
	)

	glog.Infof("Retrieving running instancing %v", instanceId)

	c.GetClient()
	resp, err = c.Client.R().
		Get(c.BaseURL + TcaApiVnfLcmVnfInstance + "/" + instanceId)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes response.CnfInstancesError
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error %v", errRes.Detail)
			return nil, fmt.Errorf(errRes.Detail)
		}
		glog.Errorf("unknown error, status code: %v %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	//
	var lcmInfo response.LcmInfo
	if err := json.Unmarshal(resp.Body(), &lcmInfo); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	return &lcmInfo, nil
}

// TerminateInstance rest call, terminates CNF/VNF
// terminateReq *specs.LcmTerminateRequest describes
// a request.
func (c *RestClient) TerminateInstance(terminateUri string, terminateReq *specs.LcmTerminateRequest) error {

	if terminateReq == nil {
		return fmt.Errorf("nil request")
	}

	glog.Infof("Terminating instancing %v", terminateUri)

	c.GetClient()
	resp, err := c.Client.R().
		SetBody(terminateReq).
		Post(terminateUri)

	if err != nil {
		glog.Error(err)
		return err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err = json.NewDecoder(resp.RawResponse.Body).Decode(&errRes); err == nil {
			return fmt.Errorf(errRes.Message)
		}
		return fmt.Errorf("unknown error, status code: %d", resp.StatusCode())
	}

	return nil
}

// CnfRollback rest api action, rollback
func (c *RestClient) CnfRollback(ctx context.Context, href string) error {

	if c == nil {
		return fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).Post(href)

	if err != nil {
		glog.Error(err)
		return err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error %v", errRes.Details)
			return fmt.Errorf(errRes.Message)
		}
		glog.Errorf("Server return unknown error %v", string(resp.Body()))
		return fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	return nil
}

// CnfResetState action
func (c *RestClient) CnfResetState(ctx context.Context, href string) error {

	if c == nil {
		return fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).SetBody("{}").Post(href)

	if err != nil {
		glog.Error(err)
		return err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error %v", errRes.Details)
			return fmt.Errorf(errRes.Message)
		}
		glog.Errorf("Server return unknown error %v", string(resp.Body()))
		return fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	return nil
}

// CreateInstance vnf instances
func (c *RestClient) CreateInstance(ctx context.Context, req *specs.LcmCreateRequest) (*response.VNFInstantiate, error) {

	if c == nil {
		return nil, fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).
		SetBody(req).
		Post(c.BaseURL + TcaVmwareVnflcmInstances)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error %v", errRes.Details)
			return nil, fmt.Errorf(errRes.Message)
		}
		glog.Errorf("Server return unknown error %v", string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %d", resp.StatusCode())
	}

	var vnfCreateResp response.VNFInstantiate
	if err := json.Unmarshal(resp.Body(), &vnfCreateResp); err != nil {
		glog.Errorf("Failed parse server respond.")
		return nil, err
	}

	return &vnfCreateResp, nil
}

// InstanceInstantiate - instantiate CNF or VNF
// Note instance state must be terminated.
func (c *RestClient) InstanceInstantiate(ctx context.Context, instanceId string, req specs.LcmInstantiateRequest) error {

	if c == nil {
		return fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).
		SetBody(req).
		Post(c.BaseURL + fmt.Sprintf(TcaVmwareVnflcmInstantiate, instanceId))

	if err != nil {
		glog.Error(err)
		return err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error %v", errRes.Details)
			return fmt.Errorf(errRes.Message)
		}
		glog.Errorf("Server return unknown error %v", string(resp.Body()))
		return fmt.Errorf("unknown error, status code: %v %v", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// InstanceUpdateState current state of running instance.
func (c *RestClient) InstanceUpdateState(ctx context.Context,
	instanceId string, req *specs.LcmInstantiateRequest) (*response.InstanceUpdate, error) {

	if c == nil {
		return nil, fmt.Errorf("unutilized client")
	}

	var (
		resp *resty.Response
		err  error
	)

	c.GetClient()
	if req != nil {
		resp, err = c.Client.R().SetContext(ctx).
			SetBody(req).Post(c.BaseURL + fmt.Sprintf(TcaVmwareVnflcmUpdate, instanceId))
	} else {
		resp, err = c.Client.R().SetContext(ctx).SetBody("{}").Post(c.BaseURL + fmt.Sprintf(TcaVmwareVnflcmUpdate, instanceId))
	}

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		return nil, c.checkErrors(resp)
	}

	var updateReplay response.InstanceUpdate
	if err := json.Unmarshal(resp.Body(), &updateReplay); err != nil {
		glog.Errorf("Failed parse servers respond. %v", err)
		return nil, err
	}

	return &updateReplay, nil
}

// InstanceReconfigure - reconfigure cnf instance.
func (c *RestClient) InstanceReconfigure(ctx context.Context, r *specs.LcmReconfigureRequest, id string) error {

	c.GetClient()

	resp, err := c.Client.R().SetContext(ctx).SetBody(r).Post(c.BaseURL + TcaVmwareVnflcmInstance + id + "/scale")

	if err != nil {
		glog.Error(err)
		return err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error %v", errRes.Details)
			return fmt.Errorf(errRes.Message)
		}
		glog.Errorf("Server return unknown error %v", string(resp.Body()))
		return fmt.Errorf("unknown error, status code: %v %v", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// DeleteInstance - delete cnf instance.
func (c *RestClient) DeleteInstance(ctx context.Context, id string) error {

	c.GetClient()
	req := c.BaseURL + fmt.Sprintf(TcaVmwareVnflcmInstance, id)

	glog.Infof("Sending Delete %v", req)

	resp, err := c.Client.R().SetContext(ctx).Delete(req)

	if err != nil {
		glog.Error(err)
		return err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error %v", errRes.Details)
			return fmt.Errorf(errRes.Message)
		}
		glog.Errorf("Server return unknown error %v", string(resp.Body()))
		return fmt.Errorf("unknown error, status code: %v %v", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
