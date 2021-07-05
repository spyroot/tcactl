// Package api
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
package api

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
)

// GetRepos retrieve repos by tenant id
func (a *TcaApi) GetRepos() (*response.ReposList, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		glog.Error(err)
	}

	var allRepos response.ReposList
	for _, r := range tenants.TenantsList {
		repos, err := a.rest.RepositoriesQuery(&request.RepoQuery{
			QueryFilter: request.Filter{
				ExtraFilter: request.AdditionalFilters{
					VimID: r.TenantID,
				},
			},
		})

		if err != nil {
			return nil, err
		}

		allRepos.Items = append(allRepos.Items, repos.Items...)
	}

	return &allRepos, nil
}

// GetFilteredExtension retrieve repos by tenant id
func (a *TcaApi) GetFilteredExtension() (*response.ReposList, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		glog.Error(err)
	}

	var allRepos response.ReposList
	for _, r := range tenants.TenantsList {
		repos, err := a.rest.RepositoriesQuery(&request.RepoQuery{
			QueryFilter: request.Filter{
				ExtraFilter: request.AdditionalFilters{
					VimID: r.TenantID,
				},
			},
		})

		if err != nil {
			return nil, err
		}

		allRepos.Items = append(allRepos.Items, repos.Items...)
	}

	return &allRepos, nil
}
