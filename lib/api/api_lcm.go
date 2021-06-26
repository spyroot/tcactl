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
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
)

// CreateCnfInstance - create cnf instance that already
// in running or termination state.
// instanceName is instance that state must change
// poolName is target pool
// clusterName a cluster that will be used to query and search instance.
func (a *TcaApi) CreateCnfInstance(instanceName string, poolName string,
	clusterName string, vimName string, doBlock bool, _disableGrant bool, verbose bool) error {

	// resolve pool id, if client indicated target pool
	_targetPool := ""
	var err error

	if len(poolName) > 0 {
		_targetPool, _, err = a.ResolvePoolName(poolName, clusterName)
		if err != nil {
			return err
		}
	}

	_instances, err := a.rest.GetVnflcm()
	if err != nil {
		return err
	}

	// for extension request we route to correct printer
	instances, ok := _instances.(*response.CnfsExtended)
	if !ok {
		return errors.New("wrong instance type")
	}

	instance, err := instances.ResolveFromName(instanceName)
	if err != nil {
		return err
	}

	if instance.IsInCluster(vimName) == false {
		return fmt.Errorf("instance not found in %v cluster", vimName)
	}

	// Check the state
	glog.Infof("Name %v Instance ID %v State %v", instanceName, instance.CID, instance.LcmOperation)
	if IsInState(instance.Meta.LcmOperationState, StateStarting) {
		return fmt.Errorf("instance '%v', uuid '%v' need finish task", instanceName, instance.CID)
	}

	if IsInState(instance.Meta.LcmOperation, StateInstantiate) &&
		IsInState(instance.Meta.LcmOperationState, StateCompleted) {
		return fmt.Errorf("instance '%v', uuid '%v' already instantiated", instanceName, instance.CID)
	}

	var additionalVduParams = request.AdditionalParams{
		DisableGrant:        _disableGrant,
		IgnoreGrantFailure:  false,
		DisableAutoRollback: false,
	}

	for _, entry := range instance.InstantiatedNfInfo {
		additionalVduParams.VduParams = append(additionalVduParams.VduParams, request.VduParam{
			Namespace: entry.Namespace,
			RepoURL:   entry.RepoURL,
			Username:  entry.Username,
			Password:  entry.Password,
			VduName:   entry.VduID,
		})
	}

	currentNodePoolId := instance.VimConnectionInfo[0].Extra.NodePoolId
	if len(_targetPool) > 0 {
		glog.Infof("Updating target kubernetes node pool.")
		currentNodePoolId = _targetPool
	}

	var req = request.InstantiateVnfRequest{
		FlavourID:           "default",
		AdditionalVduParams: additionalVduParams,
		VimConnectionInfo: []request.VimConInfo{
			{
				ID:      instance.VimConnectionInfo[0].Id,
				VimType: "",
				Extra:   request.PoolExtra{NodePoolId: currentNodePoolId},
			},
		},
	}

	err = a.rest.CnfInstantiate(instance.CID, req)
	if err != nil {
		return err
	}

	if doBlock {
		err := a.BlockWaitStateChange(instance.CID, StateInstantiate, DefaultMaxRetry, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

// CnfRollback rollbacks instance
// if flag delete provide will also delete.
// TODO add blocking
func (a *TcaApi) CnfRollback(instanceName string) error {

	cnfs, err := a.GetAllPackages()
	if err != nil {
		return err
	}

	var cnfName = instanceName
	instance, err := cnfs.ResolveFromName(cnfName)
	if err != nil {
		return err
	}

	err = a.rest.CnfRollback(instance.CID)
	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}
