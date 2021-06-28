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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"net/http"
)

// makeDefaultHeaders default headers
func (c *RestClient) makeDefaultHeaders() {
	c.Client.SetHeader("Content-Type", defaultContentType)
	c.Client.SetHeader("Version", defaultVersion)
	c.Client.SetHeader("Accept", defaultAccept)
	c.Client.SetHeader("Authorization", c.ApiKey)
	c.Client.SetHeader("x-hm-authorization", c.ApiKey)
}

// GetClient return rest client
func (c *RestClient) GetClient() {

	if c.Client == nil {
		glog.Infof("Creating a new rest client")
		c.Client = resty.New()
		if c.SkipSsl {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: c.SkipSsl},
			}
			c.Client.SetTransport(tr)
		}
	}

	c.makeDefaultHeaders()
}

// GetVnflcm - Retrieves information about a CNF/VNF instance by reading an "Individual VNF instance" resource.
// This method shall follow the provisions specified in the tables 5.4.3.3.2-1 and 5.4.3.3.2-2
// for URI query parameters, request and response data structures, and response codes.
//
// Example of filter
// (eq,id,5c11bd9c-085d-4913-a453-572457ddffe2)
func (c *RestClient) GetVnflcm(req ...string) (interface{}, error) {

	var (
		err     error
		resp    *resty.Response
		isArray = true
	)

	c.GetClient()

	if len(req) == 0 {
		resp, err = c.Client.R().Get(c.BaseURL + TcaApiVnfLcmExtensionVnfInstance)
	}

	if len(req) == 1 {
		// attach filter and dispatch
		var queryFilter = req[0]
		glog.Infof("Attaching request filter %v", queryFilter)
		resp, err = c.Client.R().SetQueryParams(map[string]string{
			"filter": queryFilter,
		},
		).Get(c.BaseURL + TcaApiVnfLcmExtensionVnfInstance)
	}
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

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes response.CnfInstancesError
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			return nil, fmt.Errorf(errRes.Detail)
		}
		return nil, fmt.Errorf("unknown error, status code: %d", resp.StatusCode())
	}

	// for single cnf request, pack return result in array.
	if isArray == false {
		var cnflcm response.LcmInfo
		var cnfslcm response.Cnfs
		if err := json.Unmarshal(resp.Body(), &cnflcm); err != nil {
			return nil, err
		}
		cnfslcm.CnfLcms = append(cnfslcm.CnfLcms, cnflcm)
		return &cnfslcm, nil
	}

	// default case
	var cnfs response.CnfsExtended
	if err := json.Unmarshal(resp.Body(), &cnfs.CnfLcms); err != nil {
		return nil, err
	}
	return &cnfs, nil
}

// GetRunningVnflcm return state of CNF
func (c *RestClient) GetRunningVnflcm(r string) (*response.LcmInfo, error) {

	var (
		err  error
		resp *resty.Response
	)

	c.GetClient()
	resp, err = c.Client.R().Get(c.BaseURL + TcaApiVnfLcmVnfInstance + "/" + r)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes response.CnfInstancesError
		if err = json.Unmarshal(resp.Body(), &errRes); err == nil {
			glog.Errorf("Server return error %v", errRes.Detail)
			return nil, fmt.Errorf(errRes.Detail)
		}
		glog.Errorf("unknown error, status code: %v %v", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("unknown error, status code: %v", resp.StatusCode())
	}

	//
	var cnflcm response.LcmInfo
	if err := json.Unmarshal(resp.Body(), &cnflcm); err != nil {
		return nil, err
	}

	return &cnflcm, nil
}

// TerminateInstance action
func (c *RestClient) TerminateInstance(terminateUri string, terminateReq request.TerminateVnfRequest) error {

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

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(resp.RawResponse.Body).Decode(&errRes); err == nil {
			return fmt.Errorf(errRes.Message)
		}
		return fmt.Errorf("unknown error, status code: %d", resp.StatusCode())
	}

	return nil
}

// CnfRollback action
func (c *RestClient) CnfRollback(ctx context.Context, instanceId string) error {

	if c == nil {
		return fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().
		Post(c.BaseURL + "/telco/api/vnflcm/v2/vnf_lcm_op_occs/" + instanceId + "/rollback")

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
func (c *RestClient) CreateInstance(ctx context.Context, req *request.CreateVnfLcm) (*response.VNFInstantiate, error) {

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
func (c *RestClient) InstanceInstantiate(ctx context.Context, instanceId string, req request.InstantiateVnfRequest) error {

	if c == nil {
		return fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).
		SetBody(req).
		Post(c.BaseURL + TcaVmwareVnflcmInstance + instanceId + "/instantiate")

	if err != nil {
		glog.Error(err)
		return err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
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
func (c *RestClient) InstanceUpdateState(ctx context.Context, instanceId string, req request.InstantiateVnfRequest) error {

	if c == nil {
		return fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().SetContext(ctx).
		SetBody(req).
		Post(c.BaseURL + "/hybridity/api/vnflcm/v1/vnf_instances/" + instanceId + "/update_state")

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

// InstanceReconfigure - reconfigure cnf instance.
func (c *RestClient) InstanceReconfigure(ctx context.Context, r *request.CnfReconfigure, id string) error {

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
