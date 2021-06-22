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

package response

import (
	"fmt"
	"strings"
)

const (
	ActiveRepo = "ENABLED"
)

type Repos struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
}

// RepoSpec - list entry in repos
type RepoSpec struct {
	ID    string  `json:"id" yaml:"id"`
	State string  `json:"state" yaml:"state"`
	Repos []Repos `json:"repos" yaml:"repos"`
	Error string  `json:"error" yaml:"error"`
}

// isActive() Return true if repo is active
func (r *RepoSpec) isActive() bool {
	return r.State == ActiveRepo
}

type ReposList struct {
	Items []RepoSpec `json:"items" yaml:"items"`
}

// GetRepoId return repo uuid and lookup performed by repo url
// http://repo.cnfdemo.io/chartrepo/library
func (r *ReposList) GetRepoId(repoUrl string) (string, error) {

	if r == nil {
		return "", fmt.Errorf("uninitialized object")
	}

	for _, it := range r.Items {
		for _, repos := range it.Repos {
			if strings.Contains(repos.Name, repoUrl) {
				return it.ID, nil
			}
		}
	}

	return "", fmt.Errorf("repository not found")
}
