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
	"github.com/spyroot/hestia/pkg/io"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Details string `json:"detail"`
	Path    string `json:"path"`
}

type PackageCreateError struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

type ErrorsResponse struct {
	ErrorResponses []ErrorResponse `json:"errors"`
}

// GetErrors - return all errors in single array
func (e *ErrorsResponse) GetErrors() []string {
	var messages []string
	for _, er := range e.ErrorResponses {
		messages = append(messages, er.Message)
	}
	return messages
}

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// AuthorizationReq - json body for authorization
type AuthorizationReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RestClient struct {
	// BaseURL TCA Base API endpoint
	BaseURL string
	// ApiKey key returned by TCA
	ApiKey string
	// SkipSsl in case client need skip self sign cert check
	SkipSsl bool
	// Client rest client
	Client *resty.Client
	// Username tca username
	Username string
	// Password tca password
	Password string
	//
	IsDebug bool
	// CertFile path to cert "certs/client.pem"
	CertFile string
	// CertKey path to cert key "certs/client.key"
	CertKey string
	// dump server respond ( dubug )
	dumpRespond           bool
	isBasicAuthentication bool
}

const (
	defaultContentType  = "application/json"
	uriAuthorize        = "/hybridity/api/sessions"
	authorizationHeader = "x-hm-authorization"
)

// GetAuthorization retrieve API key from TCA
func (c *RestClient) GetAuthorization() (bool, error) {

	c.Client = resty.New()

	// loads cert or skip ssl
	if len(c.CertKey) > 0 && len(c.CertKey) > 0 {
		if io.FileExists(c.CertFile) && io.FileExists(c.CertKey) {
			cert, err := tls.LoadX509KeyPair(c.CertFile, c.CertKey)
			if err != nil {
				glog.Fatalf("ERROR client certificate: %s", err)
			}
			c.Client.SetCertificates(cert)
		}
	}

	if c.SkipSsl {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.Client.SetTransport(tr)
	}

	resp, err := c.Client.R().
		SetHeader("Content-Type", defaultContentType).
		SetBody(AuthorizationReq{
			Username: c.Username,
			Password: c.Password,
		}).
		Post(c.BaseURL + uriAuthorize)

	if err != nil {
		glog.Errorf("Failed authorize %v", err)
		return false, err
	}

	glog.Infof("Response status code: %v", resp.StatusCode())
	glog.Infof("Response status: %v", resp.Status())
	glog.Infof("Response time: %v", resp.Time())

	if resp.StatusCode() == http.StatusOK {
		c.ApiKey = resp.Header().Get(authorizationHeader)
		return len(c.ApiKey) > 0, nil
	}

	return false, fmt.Errorf("server return %v", resp.StatusCode())
}

// checkError - check error , log it
func (c *RestClient) checkError(r *resty.Response) error {

	if r.StatusCode() == http.StatusNotFound {
		glog.Error("API resource not found")
		return fmt.Errorf("api resource not found")
	}

	var errRes ErrorResponse
	if err := json.Unmarshal(r.Body(), &errRes); err == nil {
		glog.Errorf("Server return error %s path %s msg %s", errRes.Error, errRes.Path, errRes.Message)
		return fmt.Errorf("error %s path %s msg %s", errRes.Error, errRes.Path, errRes.Message)
	} else {
		glog.Errorf("Failed parse server respond.")
	}

	return fmt.Errorf("unknown error, status code: %v ", r.StatusCode())
}

func (c *RestClient) checkErrors(r *resty.Response) error {

	if r.StatusCode() == http.StatusNotFound {
		glog.Errorf("API resource not found")
		return fmt.Errorf("api resource not found")
	}

	var errRes ErrorsResponse
	if err := json.Unmarshal(r.Body(), &errRes); err == nil {
		glog.Errorf("Server return error %v", errRes.GetErrors())
		return fmt.Errorf("error %v", errRes.GetErrors())
	} else {
		glog.Errorf("Failed parse server respond.")
	}

	return fmt.Errorf("unknown error, status code: %v ", r.StatusCode())
}
