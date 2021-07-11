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
	"github.com/spyroot/tcactl/pkg/io"
	"strings"
)

// CmdDescClusterNodePools - describe all node pool
// command will output all node pool currently
func (ctl *TcaCtl) CmdDescClusterNodePools() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:     "pools [name or id]",
		Short:   "Command describes all node pool",
		Long:    templates.LongDesc(`Command describes all node pool, it outputs all node pool currently in a system`),
		Example: "tcactl describe pools",
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

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

			allSpecs, err := ctl.tca.GetAllNodePools(ctx)
			CheckErrLogError(err)

			if _printer, ok := ctl.NodePoolPrinter[_defaultPrinter]; ok {
				_printer(&response.NodePool{
					Pools: allSpecs,
				}, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}

// CmdGetPoolNodes - command to get CNF Catalog entity
func (ctl *TcaCtl) CmdGetPoolNodes() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "nodes",
		Short: "Command returns kubernetes node pool",
		Long: templates.LongDesc(
			`Command returns a list kubernetes node pool for a given cluster name.`),
		Example: "tcactl get clusters pool 794a675c-777a-47f4-8edb-36a686ef4065",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

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

			clusters, err := ctl.tca.GetClusters(ctx)
			if err != nil || clusters == nil {
				glog.Errorf("Failed retrieve cluster list %v", err)
				return
			}

			clusterId, err := clusters.GetClusterId(args[0])
			CheckErrLogError(err)

			pool, err := ctl.tca.GetClusterNodePools(clusterId)
			if err != nil {
				glog.Errorf("Failed retrieve node pools %v", err)
				return
			}
			if _printer, ok := ctl.NodesPrinter[_defaultPrinter]; ok {
				_printer(pool, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}

// CmdDeletePoolNodes - command delete a node pool from a cluster.
// Note worker node must not have any active instances.
func (ctl *TcaCtl) CmdDeletePoolNodes() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:     "pool [cluster name or id,  pool name or id]",
		Short:   "Command deletes kubernetes node pool.",
		Long:    `Command deletes kubernetes node pool.`,
		Example: "tcactl delete pool my_cluster my_pool",
		Aliases: []string{"pools", "node_pool"},
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

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

			task, err := ctl.tca.DeleteNodePool(ctx, args[0], args[1])
			CheckErrLogError(err)

			fmt.Printf("Node pool deleted, task id %v\n", task.OperationId)
		},
	}

	return _cmd
}

// CmdCreatePoolNodes - command create a node pool
// from a node pool spec.
func (ctl *TcaCtl) CmdCreatePoolNodes() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		isDry           bool
		doBlock         bool
		showProgress    bool
	)

	var _cmd = &cobra.Command{
		Use:   "pool [cluster name or id,  spec file]",
		Short: "Command create additional node pool on target kubernetes cluster.",
		Long: templates.LongDesc(`
Command create additional node pool on target kubernetes cluster.
`),
		Example: "tcactl create node-pool my_cluster example/node-pool.yaml",
		Aliases: []string{"pools", "pool"},
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

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

			nodePoolSpec, err := specs.ReadNodeSpecFromFile(args[1])
			CheckErrLogError(err)

			if isDry && nodePoolSpec != nil {
				err := io.YamlPrinter(nodePoolSpec, false)
				CheckErrLogError(err)
			}

			task, err := ctl.tca.CreateNewNodePool(ctx, &api.NodePoolCreateApiReq{
				Spec:       nodePoolSpec,
				Cluster:    args[0],
				IsDryRun:   isDry,
				IsVerbose:  showProgress,
				IsBlocking: doBlock,
			})

			CheckErrLogError(err)
			fmt.Printf("Node Pool task %v created.\n", task.OperationId)
		},
	}

	_cmd.Flags().BoolVar(&isDry,
		"dry", false,
		"Parses input spec, validates and outputs spec to the terminal screen.")
	//
	_cmd.Flags().BoolVarP(&doBlock, CliBlock, "b", false,
		"Blocks and wait task to finish.")

	//
	_cmd.Flags().BoolVarP(&showProgress, CliProgress, "p", true,
		"Show task progress.")

	return _cmd
}

// CmdUpdatePoolNodes - command create a node pool
// from a node pool spec.
// TODO
func (ctl *TcaCtl) CmdUpdatePoolNodes() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		isDry           bool
		doBlock         bool
		showProgress    bool
	)

	var _cmd = &cobra.Command{
		Use:   "pool [cluster name or id,  spec file]",
		Short: "Command update node pool for target kubernetes cluster.",
		Long: templates.LongDesc(`
Command update node pool for target kubernetes cluster.`),
		Example: "tcactl update node-pool my_cluster example/node-pool.yaml",
		Aliases: []string{"pools", "pool"},
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

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

			nodePoolSpec, err := specs.ReadNodeSpecFromFile(args[1])
			CheckErrLogError(err)

			if isDry && nodePoolSpec != nil {
				err := io.YamlPrinter(nodePoolSpec, false)
				CheckErrLogError(err)
			}

			task, err := ctl.tca.UpdateNodePool(ctx, &api.NodePoolCreateApiReq{
				Spec:       nodePoolSpec,
				Cluster:    args[0],
				IsDryRun:   isDry,
				IsVerbose:  showProgress,
				IsBlocking: doBlock,
			})
			CheckErrLogError(err)
			fmt.Printf("Node Pool task %v created.\n", task.OperationId)
		},
	}

	_cmd.Flags().BoolVar(&isDry,
		"dry", false, "Parses input template spec, "+
			"validates, outputs spec to the terminal screen. Format based on -o flag.")

	//
	_cmd.Flags().BoolVarP(&doBlock, CliBlock, "b", false,
		"Blocks and wait task to finish.")

	//
	_cmd.Flags().BoolVarP(&showProgress, CliProgress, "p", true,
		"Show task progress.")

	return _cmd
}
