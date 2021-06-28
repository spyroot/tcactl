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
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/pkg/io"
	"io/ioutil"
	"strings"
)

// CreateCnfInstance - create cnf instance that already
// in running or termination state.
// instanceName is instance that state must change
// poolName is target pool
// clusterName a cluster that will be used to query and search instance.
func (a *TcaApi) CreateCnfInstance(ctx context.Context, instanceName string, poolName string,
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

	err = a.rest.InstanceInstantiate(ctx, instance.CID, req)
	if err != nil {
		return err
	}

	if doBlock {
		err := a.BlockWaitStateChange(ctx, instance.CID, StateInstantiate, DefaultMaxRetry, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteCnfInstance - deletes instance
func (a *TcaApi) DeleteCnfInstance(ctx context.Context, instanceName string, vimName string, isForce bool) error {

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

	if isForce && !strings.Contains(instance.Meta.LcmOperation, StateTerminate) {

		fmt.Printf("Terminating cnf instance %s state %s status %s\n",
			instance.CID, instance.Meta.LcmOperation, instance.Meta.LcmOperationState)
		err := a.TerminateCnfInstance(ctx, instanceName, vimName, true, false)
		if err != nil {
			return err
		}

		terminated, err := a.rest.GetRunningVnflcm(instance.CID)
		if err != nil {
			return err
		}

		instance.Meta.LcmOperationState = terminated.Metadata.LcmOperationState
		instance.Meta.LcmOperation = terminated.Metadata.LcmOperation

		fmt.Printf("Instance state %s and operation state, %s\n",
			terminated.Metadata.LcmOperation,
			terminated.Metadata.LcmOperationState)
	}

	if strings.Contains(instance.Meta.LcmOperation, StateTerminate) &&
		strings.Contains(instance.Meta.LcmOperationState, StateCompleted) {
		// for force case we terminate and block.
		return a.rest.DeleteInstance(ctx, instance.CID)
	}

	return errors.New("Instance must be terminated before delete operation")
}

// CnfReconfigure - reconfigure existing instance
func (a *TcaApi) CnfReconfigure(ctx context.Context, instanceName string, valueFile string,
	vduName string, isDry bool) error {

	if a.rest == nil {
		return fmt.Errorf("rest interface is nil")
	}

	if len(instanceName) == 0 {
		return fmt.Errorf("instance name empty string")
	}

	if len(vduName) == 0 {
		return fmt.Errorf("vdu name empty string")
	}

	if len(valueFile) == 0 {
		return fmt.Errorf("value file is empty string")
	}

	if !io.FileExists(valueFile) {
		return fmt.Errorf("specify valid path to value file")
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

	vduId := instance.CID
	chartName := ""
	for _, vdu := range instance.InstantiatedVnfInfo {
		if vdu.ChartName == vduName {
			chartName = vdu.ChartName
		}
	}

	if len(chartName) == 0 {
		return fmt.Errorf("chart name %s not found", vduName)
	}

	b, err := ioutil.ReadFile(valueFile)
	if err != nil {
		glog.Errorf("Failed to read value file.")
		return err
	}

	override := b64.StdEncoding.EncodeToString(b)
	p := request.VduParams{
		Overrides: override,
		ChartName: chartName,
	}

	var newVduParams []request.VduParams
	newVduParams = append(newVduParams, p)

	req := request.CnfReconfigure{}
	req.AdditionalParams.VduParams = newVduParams
	req.AspectId = request.AspectId
	req.NumberOfSteps = 2
	req.Type = request.LcmTypeScaleOut

	if isDry {
		return nil
	}

	return a.rest.InstanceReconfigure(ctx, &req, vduId)
}

// TerminateCnfInstance method terminate cnf instance
// caller need provider either name or uuid and vimName
// doBlock block and wait task to finish
// verbose output status on screen after each pool timer.
func (a *TcaApi) TerminateCnfInstance(ctx context.Context, instanceName string,
	vimName string, doBlock bool, verbose bool) error {

	respond, err := a.rest.GetVnflcm()
	if err != nil {
		return err
	}

	cnfs, ok := respond.(*response.CnfsExtended)
	if !ok {
		return errors.New("Received wrong object type")
	}

	instance, err := cnfs.ResolveFromName(instanceName)
	if err != nil {
		return err
	}

	if instance.IsInCluster(vimName) == false {
		return fmt.Errorf("instance not found in %v cluster", vimName)
	}

	glog.Infof("terminating cnfName %v instance ID %v.", instanceName, instance.CID)
	if strings.Contains(instance.Meta.LcmOperationState, "STARTING") {
		return fmt.Errorf("'%v' instance ID %v "+
			"need finish current action", instanceName, instance.CID)
	}

	if strings.Contains(instance.Meta.LcmOperation, "TERMINATE") &&
		strings.Contains(instance.Meta.LcmOperationState, "COMPLETED") {
		return fmt.Errorf("'%v' instance ID %v "+
			"already terminated", instanceName, instance.CID)
	}

	if err = a.rest.TerminateInstance(
		instance.Links.Terminate.Href,
		request.TerminateVnfRequest{
			TerminationType:            "GRACEFUL",
			GracefulTerminationTimeout: 120,
		}); err != nil {
		return err
	}

	if doBlock {
		err := a.BlockWaitStateChange(ctx, instance.CID, "NOT_INSTANTIATED", DefaultMaxRetry, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

// RollbackCnf rollbacks instance
// if flag delete provide will also delete.
// TODO add blocking
func (a *TcaApi) RollbackCnf(ctx context.Context, instanceName string, doBlock bool, verbose bool) error {

	cnfs, err := a.GetAllPackages()
	if err != nil {
		return err
	}

	var cnfName = instanceName
	instance, err := cnfs.ResolveFromName(cnfName)
	if err != nil {
		return err
	}

	err = a.rest.CnfRollback(ctx, instance.CID)
	if err != nil {
		glog.Error(err)
		return err
	}

	if doBlock {
		err := a.BlockWaitStateChange(ctx, instance.CID, StateInstantiate, DefaultMaxRetry, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteCnf delete CNF or VNF instance
func (a *TcaApi) DeleteCnf(ctx context.Context, instanceName string) error {

	cnfs, err := a.GetAllPackages()
	if err != nil {
		return err
	}

	var cnfName = instanceName
	instance, err := cnfs.ResolveFromName(cnfName)
	if err != nil {
		return err
	}

	err = a.rest.CnfRollback(ctx, instance.CID)
	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}
