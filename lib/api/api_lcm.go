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
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/io"
	"io/ioutil"
	"strings"
)

func (a *TcaApi) GetInstance(ctx context.Context, NameOrId string) (*response.LcmInfo, error) {

	var (
		iid = NameOrId
		err error
	)

	if !IsValidUUID(NameOrId) {
		iid, err = a.ResolveInstanceName(NameOrId)
		if err != nil {
			return nil, err
		}
	}

	return a.rest.GetRunningVnflcm(iid)
}

//ResolveInstanceName resolve instance name to id
func (a *TcaApi) ResolveInstanceName(name string) (string, error) {

	_instances, err := a.rest.GetVnflcm()
	if err != nil {
		return "", err
	}

	instances, ok := _instances.(*response.CnfsExtended)
	if ok {
		instance, err := instances.ResolveFromName(name)
		if err != nil {
			return "", err
		}

		return instance.CID, nil
	}

	return "", fmt.Errorf("not found %s", name)
}

// CreateCnfInstance - create cnf instance that already in running, not instantiated state
// or termination state.
// method take req *CreateInstanceApiReq
//  where   InstanceName is instance that state must change
// 		    poolName is target pool
// 			ClusterName a cluster that will be used to query and search instance.
func (a *TcaApi) CreateCnfInstance(ctx context.Context, req *CreateInstanceApiReq) error {

	// resolve pool id, if client indicated target pool
	_targetPool := ""
	var err error

	if len(req.PoolName) > 0 {
		_targetPool, _, err = a.ResolvePoolName(ctx, req.PoolName, req.ClusterName)
		if err != nil {
			return err
		}
	}

	_instances, err := a.rest.GetVnflcm()
	if err != nil {
		return err
	}

	instances, ok := _instances.(*response.CnfsExtended)
	if !ok {
		return errors.New("wrong instance type")
	}

	instance, err := instances.ResolveFromName(req.InstanceName)
	if err != nil {
		return err
	}

	if instance.IsInCluster(req.ClusterName) == false {
		return fmt.Errorf("instance not found in %v cluster", req.ClusterName)
	}

	if instance.IsStarting() {
		glog.Warning("Instance in starting state.")
		return nil
	}

	if !instance.IsInstantiated() {
		// Check the state
		glog.Infof("Name %v Instance ID %v State %v", req.InstanceName, instance.CID, instance.LcmOperation)
		if IsInState(instance.Meta.LcmOperationState, StateStarting) {
			return fmt.Errorf("instance '%v', uuid '%v' need finish task", req.InstanceName, instance.CID)
		}
		if IsInState(instance.Meta.LcmOperation, StateInstantiate) &&
			IsInState(instance.Meta.LcmOperationState, StateCompleted) {
			return fmt.Errorf("instance '%v', uuid '%v' already instantiated", req.InstanceName, instance.CID)
		}
	}

	var additionalVduParams specs.AdditionalParams
	if req.AdditionalParam != nil {
		additionalVduParams = *req.AdditionalParam
	} else {
		// default if not provided
		additionalVduParams = specs.AdditionalParams{
			DisableGrant:        false,
			IgnoreGrantFailure:  false,
			DisableAutoRollback: false,
		}
	}

	for _, entry := range instance.InstantiatedNfInfo {

		var (
			namespace = entry.Namespace
			repoUrl   = entry.RepoURL
			username  = entry.Username
			password  = entry.Password
			//	req.Spec.ClusterPassword = b64.StdEncoding.EncodeToString([]byte(req.Spec.ClusterPassword))
		)
		namespace = entry.Namespace
		if len(req.Namespace) > 0 {
			namespace = req.Namespace
		}
		if len(req.RepoUrl) > 0 {
			namespace = req.Namespace
		}
		if len(req.RepoUsername) > 0 {
			namespace = req.RepoUsername
		}
		if len(req.RepoPassword) > 0 {
			namespace = b64.StdEncoding.EncodeToString([]byte(req.RepoPassword))
		}
		additionalVduParams.VduParams = append(additionalVduParams.VduParams, specs.VduParam{
			Namespace: namespace,
			RepoURL:   repoUrl,
			Username:  username,
			Password:  password,
			VduName:   entry.VduID,
		})
	}

	currentNodePoolId := instance.VimConnectionInfo[0].Extra.NodePoolId
	if len(_targetPool) > 0 {
		glog.Infof("Updating target kubernetes node pool.")
		currentNodePoolId = _targetPool
	}

	var instantiateReq = specs.LcmInstantiateRequest{
		FlavourID:           "default",
		AdditionalVduParams: &additionalVduParams,
		// construct placement from request.
		VimConnectionInfo: []models.VimConnectionInfo{
			{
				Id:      instance.VimConnectionInfo[0].Id,
				VimType: "",
				Extra:   &models.VimExtra{NodePoolId: currentNodePoolId},
			},
		},
	}

	err = a.rest.InstanceInstantiate(ctx, instance.CID, instantiateReq)
	if err != nil {
		return err
	}

	if req.IsBlocking {
		err := a.BlockWaitStateChange(ctx, instance.CID, StateInstantiate, DefaultMaxRetry, req.IsVerbose)
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

	// if instance in roll back state just delete it
	if instance.IsStateRollback() {
		// for force case we terminate and block.
		return a.rest.DeleteInstance(ctx, instance.CID)
	}

	if isForce && !strings.Contains(instance.Meta.LcmOperation, StateTerminate) {

		fmt.Printf("Terminating cnf instance %s state %s status %s\n",
			instance.CID, instance.Meta.LcmOperation, instance.Meta.LcmOperationState)

		err := a.TerminateCnfInstance(ctx, &TerminateInstanceApiReq{
			InstanceName: instanceName,
			ClusterName:  vimName,
			IsBlocking:   true,
			IsVerbose:    false,
		})
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

	if (strings.Contains(instance.Meta.LcmOperation, StateTerminate) &&
		strings.Contains(instance.Meta.LcmOperationState, StateCompleted)) || (instance.IsStateRollback()) {
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
	p := specs.VduParams{
		Overrides: override,
		ChartName: chartName,
	}

	var newVduParams []specs.VduParams
	newVduParams = append(newVduParams, p)

	req := specs.LcmReconfigureRequest{}
	req.AdditionalParams.VduParams = newVduParams
	req.AspectId = specs.AspectId
	req.NumberOfSteps = 2
	req.Type = specs.LcmTypeScaleOut

	if isDry {
		return nil
	}

	return a.rest.InstanceReconfigure(ctx, &req, vduId)
}

// TerminateCnfInstance method terminate cnf instance
// caller need provider either name or uuid and vimName
// doBlock block and wait task to finish
// verbose output status on screen after each pool timer.
func (a *TcaApi) TerminateCnfInstance(ctx context.Context, req *TerminateInstanceApiReq) error {

	respond, err := a.rest.GetVnflcm()
	if err != nil {
		return err
	}

	cnfs, ok := respond.(*response.CnfsExtended)
	if !ok {
		return errors.New("Received wrong object type")
	}

	instance, err := cnfs.ResolveFromName(req.InstanceName)
	if err != nil {
		return err
	}

	if instance.IsInCluster(req.ClusterName) == false {
		return fmt.Errorf("instance not found in %v cluster", req.ClusterName)
	}

	glog.Infof("terminating cnfName %v instance ID %v.", req.InstanceName, instance.CID)
	if strings.Contains(instance.Meta.LcmOperationState, "STARTING") {
		return fmt.Errorf("'%v' instance ID %v "+
			"need finish current action", req.InstanceName, instance.CID)
	}

	if strings.Contains(instance.Meta.LcmOperation, "TERMINATE") &&
		strings.Contains(instance.Meta.LcmOperationState, "COMPLETED") {
		return fmt.Errorf("'%v' instance ID %v "+
			"already terminated", req.InstanceName, instance.CID)
	}

	if err = a.rest.TerminateInstance(
		instance.Links.Terminate.Href,
		&specs.LcmTerminateRequest{
			TerminationType:            "GRACEFUL",
			GracefulTerminationTimeout: 120,
		}); err != nil {
		return err
	}

	if req.IsBlocking {
		err := a.BlockWaitStateChange(ctx, instance.CID, "NOT_INSTANTIATED", DefaultMaxRetry, req.IsVerbose)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetLcmActions return list of available actions
// in current state
func (a *TcaApi) GetLcmActions(ctx context.Context, instanceName string) (*models.PolicyLinks, error) {

	lcmState, err := a.rest.GetVnflcm()
	if err != nil {
		return nil, err
	}

	cnfs, ok := lcmState.(*response.CnfsExtended)
	if !ok {
		return nil, errors.New("Received wrong object type")
	}

	instance, err := cnfs.ResolveFromName(instanceName)
	if err != nil {
		return nil, err
	}

	return &instance.Links, nil
}

// RollbackCnf rollbacks instance
// if flag delete provide will also delete.
func (a *TcaApi) RollbackCnf(ctx context.Context, instanceName string, doBlock bool, verbose bool) error {

	lcmState, err := a.rest.GetVnflcm()
	if err != nil {
		return err
	}

	cnfs, ok := lcmState.(*response.CnfsExtended)
	if !ok {
		return errors.New("Received wrong object type")
	}

	instance, err := cnfs.ResolveFromName(instanceName)
	if err != nil {
		return err
	}

	if len(instance.Links.Rollback.Href) == 0 {
		return fmt.Errorf("rollback action not avaliable in current state")
	}

	err = a.rest.CnfRollback(ctx, instance.Links.Rollback.Href)
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

// ResetState reset instance state.
func (a *TcaApi) ResetState(ctx context.Context, req *ResetInstanceApiReq) error {

	lcmState, err := a.rest.GetVnflcm()
	if err != nil {
		return err
	}

	cnfs, ok := lcmState.(*response.CnfsExtended)
	if !ok {
		return errors.New("Received wrong object type")
	}

	instance, err := cnfs.ResolveFromName(req.InstanceName)
	if err != nil {
		return err
	}

	if len(instance.Links.UpdateState.Href) == 0 {
		return fmt.Errorf("update action not avaliable in current state")
	}

	err = a.rest.CnfResetState(ctx, instance.Links.UpdateState.Href)
	if err != nil {
		glog.Error(err)
		return err
	}

	if req.IsBlocking {
		err := a.BlockWaitStateChange(ctx, instance.CID, StateInstantiate, DefaultMaxRetry, req.IsVerbose)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateCnfState update instance state.
func (a *TcaApi) UpdateCnfState(ctx context.Context, req *UpdateInstanceApiReq) (*response.InstanceUpdate, error) {
	var (
		instanceId = req.InstanceName
		err        error
	)

	if !IsValidUUID(req.InstanceName) {
		instanceId, err = a.ResolveInstanceName(req.InstanceName)
		if err != nil {
			return nil, err
		}
	}

	rep, err := a.rest.InstanceUpdateState(ctx, instanceId, req.UpdateReq)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if req.IsBlocking {
		err := a.BlockWaitStateChange(ctx, instanceId, StateInstantiate, DefaultMaxRetry, req.IsVerbose)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
	}

	return rep, nil
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
