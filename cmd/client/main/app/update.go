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
package app

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/client/response"
	"github.com/spyroot/hestia/pkg/io"
	"strings"
	"time"
)

const (
	DefaultMaxRetry = 32
	CliBlock        = "block"
	CliPool         = "pool"
	CliDisableGran  = "grant"
)

// BlockWaitStateChange blocs main thread of execution and wait specific state to change a state,
// it busy waiting.
func (ctl *TcaCtl) BlockWaitStateChange(instanceId string, waitFor string, maxRetry int) error {

	for i := 1; i < maxRetry; i++ {
		instance, err := ctl.TcaClient.GetRunningVnflcm(instanceId)
		if err != nil {
			glog.Error(err)
			return err
		}
		fmt.Printf("Current state %s waiting for %s\n", instance.InstantiationState, waitFor)
		fmt.Printf("Current LCM Operation status %s target state %s\n\n",
			instance.Metadata.LcmOperationState, instance.Metadata.LcmOperation)

		if strings.HasPrefix(instance.InstantiationState, waitFor) {
			break
		}

		time.Sleep(30 * time.Second)
	}

	return nil
}

// CmdTerminateInstances command to update CNF state. i.e terminate
func (ctl *TcaCtl) CmdTerminateInstances() *cobra.Command {
	// cnf instances

	var CnfId string
	var CnfName string
	var doBlock = false

	var cmdTerminate = &cobra.Command{
		Use:     "terminate",
		Short:   "Terminate CNF instance",
		Long:    `Terminate CNF instance, caller need to provide CNF Identifier.`,
		Aliases: []string{"down"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			doBlock, err := cmd.Flags().GetBool(CliBlock)
			io.CheckErr(err)

			respond, err := ctl.TcaClient.GetVnflcm()
			if err != nil {
				glog.Error(err)
				return
			}
			cnfs, ok := respond.(*response.CnfsExtended)
			if !ok {
				glog.Error("Received wrong object type.")
				return
			}

			var cnfName = args[0]
			instance, err := cnfs.ResolveFromName(cnfName)
			if err != nil {
				glog.Error(err)
				return
			}

			glog.Infof("terminating cnfName %v instance ID %v.", cnfName, instance.CID)
			if strings.Contains(instance.Meta.LcmOperationState, "STARTING") {
				glog.Errorf("cnfName %v instance ID %v need finish current action.", cnfName, instance.CID)
				return
			}

			if strings.Contains(instance.Meta.LcmOperation, "TERMINATE") &&
				strings.Contains(instance.Meta.LcmOperationState, "COMPLETED") {
				glog.Errorf("cnfName %v instance ID %v already terminated.", cnfName, instance.CID)
				return
			}

			if err = ctl.TcaClient.CnfTerminate(
				instance.Links.Terminate.Href,
				request.TerminateVnfRequest{
					TerminationType:            "GRACEFUL",
					GracefulTerminationTimeout: 120,
				}); err == nil {
				glog.Infof("Instance terminated.")
			}

			if doBlock {
				err := ctl.BlockWaitStateChange(instance.CID, "NOT_INSTANTIATED", DefaultMaxRetry)
				if err != nil {
					glog.Error(err)
					return
				}
			}
		},
	}

	//
	cmdTerminate.Flags().StringVarP(&CnfId, "cnf_id", "i", "",
		"cnf running instance id.")
	//
	cmdTerminate.Flags().StringVarP(&CnfName, "cnf_name", "n", "",
		"cnf running instance name.")
	//
	cmdTerminate.Flags().BoolVarP(&doBlock, CliBlock, "b", true,
		"Blocks and re-check after interval operations status.")

	return cmdTerminate
}

// CmdRollbackInstances command to update CNF state. i.e terminate
func (ctl *TcaCtl) CmdRollbackInstances() *cobra.Command {
	// cnf instances

	var CnfId string
	var CnfName string

	var cmdTerminate = &cobra.Command{
		Use:   "rollback",
		Short: "Terminate CNF instance",
		Long:  `Terminate CNF instance, caller need to provide CNF Identifier.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			//cnfs1, ok := resp2.(respons.CnfsExtension)
			//if ok {
			//	if printer, ok := ctl.CnfInstancePrinters[_defaultPrinter]; ok {
			//		printer(&cnfs1, _defaultStyler)
			//	}
			//	return
			//}
			//cnfs2, ok := resp2.(respons.Cnfs)
			//if ok {
			//	if printer, ok := ctl.CnfInstanceExtendedPrinters[_defaultPrinter]; ok {
			//		printer(&cnfs2, _defaultStyler)
			//	}
			//}
			//

			respond, err := ctl.TcaClient.GetVnflcm()
			if err != nil {
				glog.Error(err)
			}

			cnfs, ok := respond.(response.CnfsExtended)
			if !ok {
				return
			}

			if len(args) == 0 {
				return
			}

			var cnfName = args[0]
			for _, instance := range cnfs.CnfLcms {
				if strings.Contains(instance.VnfInstanceName, cnfName) {
					glog.Infof("Rolling back cnfName %v instance ID %v.", cnfName, instance.CID)
					if strings.Contains(instance.Meta.LcmOperationState, "FAILED_TEMP") {
						err = ctl.TcaClient.CnfRollback(instance.CID)
						if err == nil {
							glog.Infof("Instance roll backed.")
						} else {
							glog.Error("Failed rollback %v", err)
						}
					}
				}
			}
		},
	}

	cmdTerminate.Flags().StringVarP(&CnfId, "cnf_id", "i", "",
		"cnf running instance id")
	cmdTerminate.Flags().StringVarP(&CnfName, "cnf_name", "n", "",
		"cnf running instance name")
	return cmdTerminate
}

// CmdUpdateInstances Update state of instance
// if instance terminated , instantiate in same environment or update environment
// for example change VIM or Node Pool.
func (ctl *TcaCtl) CmdUpdateInstances() *cobra.Command {

	var (
		_disableGrant = true
		doBlock       = false
		_targetPool   = ""
	)

	// cnf instances
	var updateInstance = &cobra.Command{
		Use:     "instantiate",
		Short:   "Updates CNF instance",
		Long:    `Updates CNF instance, caller need to provide CNF Identifier.`,
		Example: "tcactl update instantiate testapp --stderrthreshold INFO --block",
		Aliases: []string{"up"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var (
				cnfName string
			)

			cnfName = args[0]

			// block or not call
			_doBlock, err := cmd.Flags().GetBool(CliBlock)
			io.CheckErr(err)
			// disable grant validation or not
			_disableGrant, err := cmd.Flags().GetBool(CliDisableGran)
			io.CheckErr(err)
			// resolve pool id, if client indicated target pool
			poolName := cmd.Flags().Lookup(CliPool)
			if poolName.Value != nil {
				clusterEntry := ctl.RootCmd.PersistentFlags().Lookup(ConfigDefaultCloud)
				if clusterEntry == nil {
					glog.Errorf("Please indicate default cluster name or indicate --cluster")
					return
				}
				glog.Infof("Using cluster %s to retrieve node pool.", clusterEntry.Value.String())
				_targetPool, _, err = ctl.ResolvePoolName(poolName.Value.String(), clusterEntry.Value.String())
				io.CheckErr(err)
			}

			// get instance data
			respond, err := ctl.TcaClient.GetVnflcm()
			CheckErrLogError(err)
			cnfs, ok := respond.(*response.CnfsExtended)
			CheckNotOkLogError(ok, "Received wrong object type.")

			instance, err := cnfs.ResolveFromName(cnfName)
			io.CheckErr(err)

			// Check the state
			glog.Infof("Name %v Instance ID %v State %v", cnfName, instance.CID, instance.LcmOperation)
			if IsInState(instance.Meta.LcmOperation, StateInstantiate) &&
				IsInState(instance.Meta.LcmOperationState, StateCompleted) {
				glog.Errorf("cnfName %v instance ID %v already instantiated.", cnfName, instance.CID)
				return
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
				fmt.Println("Chaning target node pool.")
				currentNodePoolId = _targetPool
			}

			var req = request.InstantiateVnfRequest{
				FlavourID:           "default",
				AdditionalVduParams: additionalVduParams,
				VimConnectionInfo: []request.VimConInfo{
					request.VimConInfo{
						ID:      instance.VimConnectionInfo[0].Id,
						VimType: "",
						Extra:   request.PoolExtra{NodePoolId: currentNodePoolId},
					},
				},
			}

			err = ctl.TcaClient.CnfInstantiate(instance.CID, req)
			CheckErrLogError(err)

			if _doBlock {
				err := ctl.BlockWaitStateChange(instance.CID, StateInstantiate, DefaultMaxRetry)
				CheckErrLogError(err)
			}
		},
	}

	//
	updateInstance.Flags().BoolVarP(&doBlock, CliBlock, "b", true,
		"Blocks and re-check after interval operations status")

	updateInstance.Flags().StringVarP(&_targetPool, CliPool, "p", "",
		"Update node pool, note it will use same VIM, "+
			"in case you need swap cloud, you need indicate explicitly")

	updateInstance.Flags().BoolVar(&_disableGrant, CliDisableGran, false,
		"Disable Helm Grant validation")

	return updateInstance
}
