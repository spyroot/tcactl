package client

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/client/response"
	"net/http"
)

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

	c.GetClient()

	var restReq string
	restReq = c.BaseURL + "/hybridity/api/repositories/query"

	resp, err := c.Client.R().
		SetBody(query).
		Post(restReq)

	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		fmt.Println(resp.Body())
		fmt.Println("Status code", resp.Status())
		var errRes response.VnfPackagesError
		if err = json.NewDecoder(resp.RawResponse.Body).Decode(&errRes); err == nil {
			return nil, fmt.Errorf(errRes.Detail)
		}
		return nil, fmt.Errorf("unknown error, status code: %v", resp.Status())
	}

	var repos response.ReposList
	if err := json.Unmarshal(resp.Body(), &repos); err != nil {
		return nil, err
	}

	return &repos, nil
}
