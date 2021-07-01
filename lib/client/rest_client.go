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
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/pkg/io"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Code    string `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
	Error   string `json:"error" yaml:"error"`
	Details string `json:"detail" yaml:"details"`
	Path    string `json:"path" yaml:"path"`
}

type PackageCreateError struct {
	Type     string `json:"type" yaml:"type"`
	Title    string `json:"title" yaml:"title"`
	Status   int    `json:"status" yaml:"status"`
	Detail   string `json:"detail" yaml:"detail"`
	Instance string `json:"instance" yaml:"instance"`
}

type ErrorsResponse struct {
	ErrorResponses []ErrorResponse `json:"errors" yaml:"errors"`
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
	Code int         `json:"code" yaml:"code"`
	Data interface{} `json:"data" yaml:"data"`
}

// AuthorizationReq - json body for authorization
type AuthorizationReq struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type RestClient struct {

	// BaseURL TCA Base API endpoint
	BaseURL string

	// apiKey key returned by TCA
	apiKey string

	// SkipSsl in case client need skip self sign cert check
	SkipSsl bool

	// Client rest client
	Client *resty.Client

	// Username tca username
	Username string

	// Password tca password
	Password string

	// CertFile path to cert "certs/client.pem"
	CertFile string

	// CertKey path to cert key "certs/client.key"
	CertKey string

	// dump server respond (debug server output)
	isTrace bool

	// authentication type
	isBasicAuthentication bool

	// set debug mode
	IsDebug bool
}

func NewRestClient(baseURL string, skipSsl bool, username string, password string) (*RestClient, error) {
	if len(baseURL) == 0 {
		return nil, errors.New("base url is empty string")
	}
	if len(username) == 0 {
		return nil, errors.New("username is empty string")
	}
	if len(password) == 0 {
		return nil, errors.New("password is empty string")
	}
	return &RestClient{BaseURL: baseURL, SkipSsl: skipSsl, Username: username, Password: password}, nil
}

// makeDefaultHeaders default headers
func (c *RestClient) makeDefaultHeaders() {
	c.Client.SetHeader("Content-Type", defaultContentType)
	c.Client.SetHeader("Version", defaultVersion)
	c.Client.SetHeader("Accept", defaultAccept)
	c.Client.SetHeader("Authorization", c.apiKey)
	c.Client.SetHeader("x-hm-authorization", c.apiKey)
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

// SetDumpRespond sets client in output mode.
func (c *RestClient) SetDumpRespond(dumpRespond bool) {
	c.isTrace = dumpRespond
}

// GetAuthorization retrieve API key from TCA
// and update internal state.
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

	if c.isTrace {

		glog.Infof("Response status code: %v", resp.StatusCode())
		glog.Infof("Response status: %v", resp.Status())
		glog.Infof("Response time: %v", resp.Time())

		fmt.Println(string(resp.Body()))
	}

	if !resp.IsSuccess() {
		var errRes ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errRes); err == nil {
			if strings.ToLower(errRes.Error) == "unauthorized" {
				glog.Errorf("Server return error %s", errRes.Message)
				return false, fmt.Errorf("authentication failed for username %s", c.Username)
			}
			glog.Errorf("Server return error %s", errRes.Message)
			return false, fmt.Errorf("error %s", errRes.Message)
		} else {
			glog.Errorf("Failed parse server respond.")
		}
	}

	if resp.StatusCode() == http.StatusOK {
		c.apiKey = resp.Header().Get(authorizationHeader)
		return len(c.apiKey) > 0, nil
	}

	return false, fmt.Errorf("server return %v", resp.StatusCode())
}

// checkError - check error, log it
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

// checkErrors - check list of errors.
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

func (c *RestClient) SetTrace(trace bool) {
	c.isTrace = trace
}

func (c *RestClient) GetApiKey() string {
	return c.apiKey
}
