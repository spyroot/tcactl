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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/pkg/io"
	"net/http"
)

const (
	HarborChartRepo = "/api/chartrepo/library/charts"
	HarborRepos     = "/api/v2.0/projects/library/repositories"
)

// Harbor
func (c *RestClient) HarborAuthenticate() (bool, error) {

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

	resp, err := c.Client.R().SetBasicAuth(c.Username, c.Password).
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

// UploadHelm - Uploads Harbor chart
func (c *RestClient) UploadHelm(csar []byte, fileName string) (bool, error) {

	c.GetClient()

	resp, err := c.Client.R().SetBasicAuth(c.Username, c.Password).
		SetFileReader("file", fileName, bytes.NewReader(csar)).
		SetHeader("Content-Type", "application/zip").
		SetContentLength(true).
		Post(c.BaseURL + HarborChartRepo)

	if err != nil {
		glog.Error(err)
		return false, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return false, c.checkError(resp)

	}

	return true, nil
}

// GetCharts - return all charts
func (c *RestClient) GetCharts() ([]response.HelmChart, error) {

	c.GetClient()
	var (
		helmCharts []response.HelmChart
		resp       *resty.Response
		err        error
	)

	if c.isBasicAuthentication {
		resp, err = c.Client.R().SetBasicAuth(c.Username, c.Password).Get(c.BaseURL + HarborChartRepo)
	} else {
		resp, err = c.Client.R().Get(c.BaseURL + HarborChartRepo)
	}

	if err != nil {
		glog.Error(err)
		return helmCharts, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return helmCharts, c.checkError(resp)
	}

	if err := json.Unmarshal(resp.Body(), &helmCharts); err != nil {
		glog.Error("Failed parse server respond.")
		return helmCharts, err
	}

	return helmCharts, nil
}

// GetRepos - return all charts
func (c *RestClient) GetRepos() ([]response.Repos, error) {

	c.GetClient()
	var (
		repos []response.Repos
		resp  *resty.Response
		err   error
	)

	if c.isBasicAuthentication {
		resp, err = c.Client.R().SetBasicAuth(c.Username, c.Password).Get(c.BaseURL + HarborChartRepo)
	} else {
		resp, err = c.Client.R().Get(c.BaseURL + HarborRepos)
	}

	if err != nil {
		glog.Error(err)
		return repos, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return repos, c.checkError(resp)
	}

	if err := json.Unmarshal(resp.Body(), &repos); err != nil {
		glog.Error("Failed parse server respond.")
		return repos, err
	}

	return repos, nil
}

// GetChart -
func (c *RestClient) GetChart(chartName string) (bool, error) {

	c.GetClient()

	resp, err := c.Client.R().Get(c.BaseURL + HarborChartRepo + "/" + chartName)

	if err != nil {
		glog.Error(err)
		return false, err
	}

	if c.isTrace && resp != nil {
		fmt.Println(string(resp.Body()))
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return false, c.checkError(resp)

	}

	return true, nil
}
