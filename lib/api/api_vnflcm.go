// Package app
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
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/client/response"
)

func (a *TcaApi) GetAllPackages() (*response.CnfsExtended, error) {

	respond, err := a.GetVnflcm()
	if err != nil {
		return nil, err
	}

	pkgs, ok := respond.(response.CnfsExtended)
	if !ok {
		return nil, errors.New("rest vnflcm return wrong type")
	}

	return &pkgs, nil
}

func (a *TcaApi) GetVnflcm(f ...string) (interface{}, error) {
	return a.rest.GetVnflcm(f...)
}
