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
	"net/http"
)

// LinkedRepositories - return linked repo to tenant's vim
// note repo must be linked
func (c *RestClient) LinkedRepositories(tenantId string, repo string) (string, error) {

	if c == nil {
		return "", fmt.Errorf("uninitialized object")
	}

	repos, err := c.RepositoriesQuery(&request.RepoQuery{
		QueryFilter: request.Filter{
			ExtraFilter: request.AdditionalFilters{
				VimID: tenantId,
			},
		},
	})

	if err != nil {
		return "", err
	}

	return repos.GetRepoId(repo)
}

//RepositoriesQuery - query repositories linked to vim
func (c *RestClient) RepositoriesQuery(query *request.RepoQuery) (*response.ReposList, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()

	resp, err := c.Client.R().
		SetBody(query).
		Post(c.BaseURL + TcaVmwareRepositories)

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

	var repos response.ReposList
	if err := json.Unmarshal(resp.Body(), &repos); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &repos, nil
}

//GetRepositoriesQuery - query repositories linked to vim
func (c *RestClient) GetRepositoriesQuery(query *request.RepoQuery) (*response.ReposList, error) {

	if c == nil {
		return nil, fmt.Errorf("uninitialized rest client")
	}

	c.GetClient()

	resp, err := c.Client.R().
		SetBody(query).
		Post(c.BaseURL + TcaVmwareRepos)

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

	var repos response.ReposList
	if err := json.Unmarshal(resp.Body(), &repos); err != nil {
		glog.Error("Failed parse server respond.")
		return nil, err
	}

	return &repos, nil
}
