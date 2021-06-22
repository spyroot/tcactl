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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/client/response"
	"net/http"
)

// makeDefaultHeaders default headers
func (c *RestClient) makeDefaultHeaders() {
	c.Client.SetHeader("Content-Type", "application/json")
	c.Client.SetHeader("Version", "2")
	c.Client.SetHeader("Accept", "application/json")
	c.Client.SetHeader("Authorization", c.ApiKey)
	c.Client.SetHeader("x-hm-authorization", c.ApiKey)
}

// GetClient return rest client
func (c *RestClient) GetClient() {

	if c.Client == nil {
		glog.Infof("Creating rest client")
		c.Client = resty.New()
		if c.SkipSsl {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
// API Docs
// /hybridity/docs/apis/index.html?component=api-nfvm-v2&schema=swagger/1.0/SOL002-SOL003/SOL003/VNFLifecycleManagement/VNFLifecycleManagement.yaml
func (c *RestClient) GetVnflcm(req ...string) (interface{}, error) {

	var (
		err     error
		resp    *resty.Response
		isArray bool = true
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
		glog.Infof("ID %v", req[1])
		resp, err = c.Client.R().Get(c.BaseURL + TcaApiVnfLcmVnfInstance + "/" + req[1])
	}

	if err != nil {
		glog.Error(err)
		return nil, err
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
func (c *RestClient) GetRunningVnflcm(req string) (*response.LcmInfo, error) {

	var (
		err  error
		resp *resty.Response
	)

	c.GetClient()
	resp, err = c.Client.R().Get(c.BaseURL + TcaApiVnfLcmVnfInstance + "/" + req)
	if err != nil {
		glog.Error(err)
		return nil, err
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

// CnfScale action
func (c *RestClient) CnfScale(scaleUri string) error {

	c.GetClient()

	resp, err := c.Client.R().Get(scaleUri)

	if err != nil {
		if c.IsDebug {
			glog.Error(err)
		}
		return err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(resp.RawResponse.Body).Decode(&errRes); err == nil {
			return fmt.Errorf(errRes.Message)
		}
		return fmt.Errorf("unknown error, status code: %d", resp.StatusCode())
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Body())
	return nil

	//var fullResponse CnfLcmResp
	//if err := json.Unmarshal(resp.Body(), &fullResponse.CnfLcms); err != nil {
	//	return err
	//}
	//
	//return &fullResponse, nil
}

// CnfTerminate action
func (c *RestClient) CnfTerminate(terminateUri string, terminateReq request.TerminateVnfRequest) error {

	c.GetClient()

	resp, err := c.Client.R().
		SetBody(terminateReq).
		Post(terminateUri)

	if err != nil {
		if c.IsDebug {
			glog.Error(err)
		}
		return err
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
func (c *RestClient) CnfRollback(instanceId string) error {

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

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
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

// CnfVnfInstantiate vnf instances
func (c *RestClient) CnfVnfInstantiate(req *request.CreateVnfLcm) (*response.VNFInstantiate, error) {

	if c == nil {
		return nil, fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().
		SetBody(req).
		Post(c.BaseURL + "/telco/api/vnflcm/v2/vnf_instances")

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
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
		return nil, err
	}

	return &vnfCreateResp, nil
}

// CnfInstantiate instantiate CNF/VNF
// The CNF that status must be terminated.
func (c *RestClient) CnfInstantiate(instanceId string, req request.InstantiateVnfRequest) error {

	if c == nil {
		return fmt.Errorf("unutilized client")
	}

	c.GetClient()
	resp, err := c.Client.R().
		SetBody(req).
		Post(c.BaseURL + "/telco/api/vnflcm/v2/vnf_instances/" + instanceId + "/instantiate")

	if err != nil {
		glog.Error(err)
		return err
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
