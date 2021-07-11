// Package cmds
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
package cmds

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/api"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"strings"
)

// CmdUpdateInstance root command
func (ctl *TcaCtl) CmdUpdateInstance() *cobra.Command {

	var cmdUpdateCnf = &cobra.Command{
		Use:   "cnf",
		Short: "Command updates, reconfigure scale, existing cnf instance state.",
		Long: templates.LongDesc(`

Command creates a new cnf instance.  By default it uses
a configuration as default parameter for cloud provider, cluster name,
node pool.

`),
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmdUpdateCnf.AddCommand(ctl.CmdTerminateInstances())
	cmdUpdateCnf.AddCommand(ctl.CmdUpdateInstances())
	cmdUpdateCnf.AddCommand(ctl.CmdRollbackInstances())
	cmdUpdateCnf.AddCommand(ctl.CmdReconfigure())

	return cmdUpdateCnf
}

// CmdGetInstances Get CNF/VNF active instances
// instance might be in different state. active define
// package that instantiate.
func (ctl *TcaCtl) CmdGetInstances() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle

		_defaultFilter string
		_instanceID    string
		_outputFilter  string
	)

	var cmdCnfInstance = &cobra.Command{
		Use:   "cnfi",
		Short: "Command returns cnf instance or all instances",
		Long:  templates.LongDesc(`Command returns cnf instance or all instance.`),

		Example: "tcactl get cnfi -o json --filter \"{eq,id,5c11bd9c-085d-4913-a453-572457ddffe2}\"",
		Run: func(cmd *cobra.Command, args []string) {

			var (
				err            error
				genericRespond interface{}
			)

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			// rest call
			if len(args) > 0 {
				genericRespond, err = ctl.tca.GetVnflcm(_defaultFilter, args[0])
			} else {
				genericRespond, err = ctl.tca.GetVnflcm(_defaultFilter)
			}
			CheckErrLogError(err)

			// for extension request we route to correct printer
			cnfsExt, ok := genericRespond.(*response.CnfsExtended)
			if ok {
				if printer, ok := ctl.CnfInstanceExtendedPrinters[_defaultPrinter]; ok {
					printer(cnfsExt, _defaultStyler)
				}
				return
			}

			// for regular request we route to correct printer
			cnfsReg, ok := genericRespond.(*response.Cnfs)
			if ok {
				if printer, ok := ctl.CnfInstancePrinters[_defaultPrinter]; ok {
					printer(cnfsReg, _defaultStyler)
				}
			}
		},
	}

	//
	cmdCnfInstance.Flags().StringVarP(&_instanceID,
		"package_id", "i", "", "VNF package id")

	//
	cmdCnfInstance.Flags().StringVar(&_defaultFilter,
		"filter", "",
		"filter for query example, filter by id --filter \"{eq,id,5c11bd9c-085d-4913-a453-572457ddffe2}\"")

	// output filter , filter specific value from data structure
	cmdCnfInstance.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return cmdCnfInstance
}

// CmdCreateCnf create a new cnf instance
func (ctl *TcaCtl) CmdCreateCnf() *cobra.Command {

	var (
		vimType             = models.VimTypeKubernetes
		namespace           string
		disableGrantFlag    bool
		disableAutoRollback bool
		ignoreGrantFailure  bool
		isDryRun            bool
		doBlock             bool
		doAutoName          bool
	)

	var cmdCreate = &cobra.Command{
		Use:   "cnf [catalog name or catalog id, and instance name]",
		Short: "Command creates a new cnf or vnf instance.",
		Long: templates.LongDesc(`
The create cnf command creates a new CNF instance.  By default, it uses
a tcactl configuration as default parameter for a cloud provider, cluster name
and node pool.
`),
		Example: "\t - tca create cnf myapp myapp-instance1\n" +
			"\t - tca create cnf myapp myapp-instance2 --disable_grant\n " +
			"\t - tca create cnf myapp myapp-instance3 --disable_grant --dry",
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			glog.Infof("Using cloud provider name='%s' cluster name='%s' repo name='%s' node pool name='%s'",
				ctl.DefaultCloudName,
				ctl.DefaultClusterName,
				ctl.DefaultRepoName,
				ctl.DefaultNodePoolName)

			if len(ctl.DefaultClusterName) == 0 {
				fmt.Println("Please indicate cloud provider, default is empty.")
				return
			}
			if len(ctl.DefaultClusterName) == 0 {
				fmt.Println("Please indicate cluster name, default is empty.")
				return
			}
			if len(ctl.DefaultRepoName) == 0 {
				fmt.Println("Please indicate repository, default is empty.")
				return
			}
			if len(ctl.DefaultNodePoolName) == 0 {
				fmt.Println("Please indicate node pool, default is empty.")
				return
			}

			var newInstanceReq = api.NewInstanceRequestSpec(ctl.DefaultCloudName,
				ctl.DefaultClusterName, vimType, args[0], ctl.DefaultRepoName, args[1], ctl.DefaultNodePoolName)

			newInstanceReq.AdditionalParams.DisableAutoRollback = disableAutoRollback
			newInstanceReq.AdditionalParams.IgnoreGrantFailure = ignoreGrantFailure
			newInstanceReq.AdditionalParams.DisableGrant = disableGrantFlag
			newInstanceReq.SetAutoName(doAutoName)

			instance, err := ctl.tca.CreateCnfNewInstance(context.Background(), newInstanceReq, isDryRun, doBlock)
			CheckErrLogError(err)

			if isDryRun == false {
				fmt.Printf("Instance %s created.", instance.Id)
			}
		},
	}

	// disables helm validation
	cmdCreate.Flags().BoolVar(&disableGrantFlag,
		CliDisableGran, false,
		"disables grant validation.")
	//
	cmdCreate.Flags().BoolVar(&disableAutoRollback,
		CliDisableAutoRollback, false,
		"disables auto rollback.")
	//
	cmdCreate.Flags().BoolVar(&ignoreGrantFailure,
		CliIgnoreGrantFailure, false,
		"disable grant failure check.")
	// namespace
	cmdCreate.Flags().StringVarP(&namespace,
		CliNamespace, "n", "default",
		"cnf namespace.")
	// dry run
	cmdCreate.Flags().BoolVar(&isDryRun,
		CliDryRun, false, "Flag instructs to run command in dry run.")
	// auto name
	cmdCreate.Flags().BoolVar(&doAutoName,
		CliAutoName, false, "Flag instructs to generate new name if name is conflicts.")
	// block
	cmdCreate.Flags().BoolVarP(&doBlock, CliBlock, "b", false,
		"Flag instructs to blocks and pools the status of the operation.")

	return cmdCreate
}

//
func (ctl *TcaCtl) CmdReconfigure() *cobra.Command {

	var (
		namespace           string
		disableGrantFlag    bool
		disableAutoRollback bool
		ignoreGrantFailure  bool
		isDryRun            bool
	)

	var cmdCreate = &cobra.Command{
		Use:   "reconfigure [instance name, vdu name, values.yaml]",
		Short: "Command reconfigures existing cnf instance.",
		Long: templates.LongDesc(`
The command reconfigures the existing CNF active instance.
`),
		Args: cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {

			glog.Infof("Using cloud provider name='%s' cluster name='%s' repo name='%s' node pool name='%s'",
				ctl.DefaultCloudName,
				ctl.DefaultClusterName,
				ctl.DefaultRepoName,
				ctl.DefaultNodePoolName)

			if len(ctl.DefaultClusterName) == 0 {
				fmt.Println("Please indicate cloud provider, default is empty.")
				return
			}
			if len(ctl.DefaultClusterName) == 0 {
				fmt.Println("Please indicate cluster name, default is empty.")
				return
			}
			if len(ctl.DefaultRepoName) == 0 {
				fmt.Println("Please indicate repository, default is empty.")
				return
			}
			if len(ctl.DefaultNodePoolName) == 0 {
				fmt.Println("Please indicate node pool, default is empty.")
				return
			}

			err := ctl.tca.CnfReconfigure(context.Background(), args[0], args[1], args[2], isDryRun)
			CheckErrLogError(err)
		},
	}

	cmdCreate.Flags().BoolVar(&disableGrantFlag,
		"disable_grant", false,
		"disable grant validation.")

	cmdCreate.Flags().BoolVar(&disableAutoRollback,
		"rollback", false,
		"disables auto Rollback.")

	cmdCreate.Flags().BoolVar(&ignoreGrantFailure,
		"ignore_failure", false,
		"disable grant failure.")

	// namespace
	cmdCreate.Flags().StringVarP(&namespace,
		"namespace", "n", "default",
		"cnf namespace.")

	cmdCreate.Flags().BoolVar(&isDryRun,
		"dry", false, "Flag instructs to run command in dry run.")

	return cmdCreate
}

// CmdTerminateInstances command to update CNF state. i.e terminate
func (ctl *TcaCtl) CmdTerminateInstances() *cobra.Command {

	var (
		doBlock      bool
		showProgress bool
	)

	var cmdTerminate = &cobra.Command{
		Use:   "terminate [instance name or id]",
		Short: "Command terminates CNF or VNF instance",
		Long: templates.LongDesc(`
The command terminates active CNF or VNF instances.
`),
		Aliases: []string{"down"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			doBlock, err := cmd.Flags().GetBool(CliBlock)
			CheckErrLogError(err)

			ctl.checkDefaultsConfig()
			err = ctl.tca.TerminateCnfInstance(context.Background(),
				&api.TerminateInstanceApiReq{
					InstanceName: args[0],
					ClusterName:  ctl.DefaultClusterName,
					IsBlocking:   doBlock,
					IsVerbose:    showProgress,
				})
			CheckErrLogError(err)

			fmt.Println("Successfully updated state.")
		},
	}

	//
	cmdTerminate.Flags().BoolVarP(&doBlock, CliBlock, "b", false,
		"Flag instructs to blocks and pools the status of the operation.")

	//
	cmdTerminate.Flags().BoolVarP(&showProgress, CliProgress, "s", false,
		"Show task progress.")

	return cmdTerminate
}

// CmdUpdateInstances Update state of instance
// if instance terminated , instantiate in same environment or update environment
// for example change VIM or Node Pool.
func (ctl *TcaCtl) CmdUpdateInstances() *cobra.Command {

	var (
		disableGrant        bool
		ignoreGrantFailure  bool
		disableAutoRollback bool
		doBlock             bool
		showProgress        bool
		_targetPool         = ""
	)

	// cnf instances
	var updateInstance = &cobra.Command{
		Use:   "instantiate [instance name or id]",
		Short: "Updates CNF or VNF instance state.",
		Long: templates.LongDesc(`

Updates CNF or VNF instance state, need to provide id or name of
of the instance. --block provides option to block and 
wait when task will finished.

`),
		Example: "\t - tcactl update cnf up testapp\n" +
			"\t - tcactl update cnf up --stderrthreshold INFO --block --pool my_pool01\n",
		Aliases: []string{"up"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctl.checkDefaultsConfig()

			// block or not call
			doBlock, err := cmd.Flags().GetBool(CliBlock)
			CheckErrLogError(err)

			// resolve pool id, if client indicated target pool
			poolName := cmd.Flags().Lookup(CliPool)

			req := api.CreateInstanceApiReq{
				InstanceName: args[0],
				PoolName:     poolName.Value.String(),
				VimName:      ctl.DefaultClusterName,
				ClusterName:  ctl.DefaultClusterName,
				IsBlocking:   doBlock,
				IsVerbose:    showProgress,
				AdditionalParam: &specs.AdditionalParams{
					DisableGrant:        disableGrant,
					IgnoreGrantFailure:  ignoreGrantFailure,
					DisableAutoRollback: disableAutoRollback,
				},
			}
			err = ctl.tca.CreateCnfInstance(context.Background(), &req)
			CheckErrLogError(err)

			fmt.Println("Successfully updated state.")
		},
	}

	// blocking flag
	updateInstance.Flags().BoolVarP(&doBlock, CliBlock, "b", false,
		"Flag instructs to blocks and pools the status of the operation.")

	// node pool flag
	updateInstance.Flags().StringVar(&_targetPool, CliPool, "",
		"Updates kubernetes node pool, note it will use same kubernetes cluster, "+
			"in case you need swap cloud, you need indicate explicitly,")
	//
	updateInstance.Flags().BoolVar(&disableGrant, CliDisableGran, false,
		"Flag disable Helm Grant validation.")
	//
	updateInstance.Flags().BoolVar(&ignoreGrantFailure, CliIgnoreGrantFailure, false,
		"Flag ignore grant failure.")
	//
	updateInstance.Flags().BoolVar(&disableAutoRollback, CliDisableAutoRollback, false,
		"Flag disable auto rollback.")
	//
	updateInstance.Flags().BoolVarP(&showProgress, CliProgress, "s", false,
		"Show task progress.")

	return updateInstance
}

// CmdDeleteInstances command deletes existing instance
// force flag provide option to terminate and delete
func (ctl *TcaCtl) CmdDeleteInstances() *cobra.Command {

	var (
		_isForce bool
	)

	// cnf instances
	var updateInstance = &cobra.Command{
		Use:   "cnf [instance name or id]",
		Short: "Command deletes CNF or VNF instance.",
		Long: templates.LongDesc(`
The command deletes CNF or VNF instances. The client must provide the ID or Name of the instance.
The instance must be in a  active state.
`),
		Example: "\t - tcactl delete cnf testapp\n" +
			"\t - tcactl delete cnf testapp --force",
		Aliases: []string{"del"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctl.checkDefaultsConfig()
			err := ctl.tca.DeleteCnfInstance(context.Background(), args[0], ctl.DefaultClusterName, _isForce)
			CheckErrLogError(err)

			fmt.Printf("Instance '%s' delete\n", args[0])
		},
	}

	// grand disable flag
	updateInstance.Flags().BoolVar(&_isForce, CliForce, false,
		"Terminate and Deletes instance.")

	return updateInstance
}

// checkDefaultsConfig - checks all default variables set.
func (ctl *TcaCtl) checkDefaultsConfig() {
	if len(ctl.DefaultClusterName) == 0 {
		panic("Please indicate cloud provider, default is empty.")
	}
	if len(ctl.DefaultClusterName) == 0 {
		panic("Please indicate cluster name, default is empty.")
	}
	if len(ctl.DefaultRepoName) == 0 {
		panic("Please indicate repository, default is empty.")
	}
	if len(ctl.DefaultNodePoolName) == 0 {
		panic("Please indicate node pool, default is empty.")
	}
}

// CmdRollbackInstances command to update CNF state. i.e terminate
func (ctl *TcaCtl) CmdRollbackInstances() *cobra.Command {

	var (
		_doBlock = false
	)

	var cmdTerminate = &cobra.Command{
		Use:   "rollback [instance name or id]",
		Short: "Command rollback CNF or VNF instance",
		Long: templates.LongDesc(
			`Rollback CNF instance, caller need to provide valid instance id or a name.`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctl.checkDefaultsConfig()

			err := ctl.tca.RollbackCnf(context.Background(), args[0], _doBlock, true)
			CheckErrLogError(err)

			fmt.Println("Successfully rollback state.")
		},
	}

	// blocking flag
	cmdTerminate.Flags().BoolVarP(&_doBlock, CliBlock, "b", false,
		"Blocks and Pool the operations status.")

	return cmdTerminate
}
