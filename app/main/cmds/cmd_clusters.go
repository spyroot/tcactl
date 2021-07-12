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
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
	"github.com/spyroot/tcactl/lib/api"
	"github.com/spyroot/tcactl/lib/api/kubernetes"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	ioutils "github.com/spyroot/tcactl/pkg/io"
	osutil "github.com/spyroot/tcactl/pkg/os"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CmdGetClusters - Sub-command to get cluster information.
func (ctl *TcaCtl) CmdGetClusters() *cobra.Command {

	var _cmdClusters = &cobra.Command{
		Use:   "clusters",
		Short: "Command retrieves cluster related information.",
		Long: templates.LongDesc(`

Command retrieves cluster-related information. Each sub-command 
require either cluster name or cluster id.

`),
		Example: "\t - tcactl get clusters info mycluster\n" +
			"\t - tcactl get clusters pool mycluster",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"cluster", "cl"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s requires a subcommand", cmd.Name())
		},
	}

	_cmdClusters.AddCommand(ctl.CmdGetClustersList())
	_cmdClusters.AddCommand(ctl.CmdGetClustersK8SConfig())
	_cmdClusters.AddCommand(ctl.CmdGetClustersPool())
	_cmdClusters.AddCommand(ctl.CmdGetClustersPoolNodes())
	_cmdClusters.AddCommand(ctl.CmdGetClusterTasks())

	return _cmdClusters
}

// CmdGetClustersPool - command return cluster pools list
func (ctl *TcaCtl) CmdGetClustersPool() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:   "pool [name or id of cluster]",
		Short: "Command returns kubernetes node pool for a given cluster.",
		Long: templates.LongDesc(
			`Command returns a list kubernetes node pool for a given cluster name.`),

		Example: "\t - tcactl get clusters pool 794a675c-777a-47f4-8edb-36a686ef4065\n " +
			"\t - tcactl get cluster mycluster",
		Run: func(cmd *cobra.Command, args []string) {

			var (
				ctx  = context.Background()
				pool *response.NodePool
				err  error
			)

			// global output type, and terminal wide or not
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_isWide, err = cmd.Flags().GetBool(FlagCliWide)
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			// for exact match for cluster
			if len(args) > 0 {
				pool, err = ctl.tca.GetNodePool(ctx, args[0])
			}

			pool, err = ctl.tca.GetAllNodePool(ctx)
			CheckErrLogError(err)

			if _printer, ok := ctl.NodePoolPrinter[_defaultPrinter]; ok {
				_printer(pool, _defaultStyler)
			}
		},
	}

	// wide output
	_cmd.Flags().BoolVarP(&_isWide,
		"wide", "w", true, "Wide output")

	return _cmd
}

// CmdGetCluster - command to get CNF Catalog entity
func (ctl *TcaCtl) CmdGetCluster() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:   "cluster [name or id]",
		Short: "Command describes kubernetes cluster or clusters information",
		Long: templates.LongDesc(
			`Command returns kubernetes cluster or cluster information.`),
		Example: "\t - tcactl describe clusters 794a675c-777a-47f4-8edb-36a686ef4065 -o json",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var (
				ctx = context.Background()
				cid string
				err error
			)

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()
			// set wide or not
			_isWide, err := cmd.Flags().GetBool("wide")
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			cid = args[0]

			cluster, err := ctl.tca.GetCluster(ctx, cid)
			if err != nil {
				clusters, err := ctl.tca.GetClusters(ctx)
				if err != nil {
					return
				}
				fmt.Println("Unknown cluster, current cluster list:")
				for _, spec := range clusters.Clusters {
					fmt.Println(" * ", spec.ClusterName)
				}
				return
			}

			if printer, ok := ctl.ClusterPrinter[_defaultPrinter]; ok {
				printer(cluster, _defaultStyler)
			}
		},
	}

	// wide output
	_cmd.Flags().BoolVarP(&_isWide,
		"wide", "w", true, "Wide output")

	return _cmd
}

// CmdGetClustersPoolNodes - command to get CNF Catalog entity
func (ctl *TcaCtl) CmdGetClustersPoolNodes() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:   "nodes",
		Short: "Command returns kubernetes nodes in pool",
		Long: templates.LongDesc(
			`Command returns a list kubernetes node pool for a given cluster name.`),

		Example: "\t - tcactl get clusters nodes 794a675c-777a-47f4-8edb-36a686ef4065\n " +
			"\t - tcactl get clusters nodes edge",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// set wide or not
			_isWide, err := cmd.Flags().GetBool("wide")
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

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

	// wide output
	_cmd.Flags().BoolVarP(&_isWide,
		"wide", "w", true, "Wide output")

	return _cmd
}

// CmdDescClusterNodePool - describe node pool
func (ctl *TcaCtl) CmdDescClusterNodePool() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:   "pool [name or id]",
		Short: "Command describes kubernetes node pool",
		Long: templates.LongDesc(
			`Command describes kubernetes node pool for a given or default cluster name.`),
		Example: "\t - tcactl describe pool 794a675c-777a-47f4-8edb-36a686ef4065",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// set wide or not
			_isWide, err := cmd.Flags().GetBool(FlagCliWide)
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			poolName := args[0]
			clusterEntry := ctl.RootCmd.PersistentFlags().Lookup(ConfigDefaultCluster)
			if clusterEntry == nil {
				glog.Error("Please indicate default cluster name or indicate --cluster")
				return
			}
			glog.Infof("Using cluster %s to retrieve node pool.", clusterEntry.Value.String())
			_targetPoolID, _clusterId, err := ctl.ResolvePoolName(poolName, clusterEntry.Value.String())
			pool, err := ctl.tca.GetClusterNodePool(ctx, _clusterId, _targetPoolID)
			if err != nil {
				glog.Errorf("Failed retrieve node pools err: '%v'", err)
				return
			}
			//
			if _printer, ok := ctl.PoolSpecPrinter[_defaultPrinter]; ok {
				_printer(pool, _defaultStyler)
			}
		},
	}

	// wide output
	_cmd.Flags().BoolVarP(&_isWide,
		"wide", "w", true, "Wide output")

	return _cmd
}

// CmdDescClusterNodes  - describe node pool
// for all cluster
func (ctl *TcaCtl) CmdDescClusterNodes() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
	)

	var _cmd = &cobra.Command{
		Use:     "nodes",
		Short:   "Command describes all kubernetes nodes",
		Long:    `Command describes all kubernetes nodes`,
		Example: "tcactl describe nodes",
		//Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			clusters, err := ctl.tca.GetClusters(ctx)
			CheckErrLogError(err)

			var allSpecs []response.NodesSpecs
			for _, c := range clusters.Clusters {
				pools, err := ctl.tca.GetClusterNodePools(c.Id)
				CheckErrLogError(err)
				for _, p := range pools.Pools {
					pool, err := ctl.tca.GetClusterNodePool(ctx, c.Id, p.Id)
					CheckErrLogError(err)
					allSpecs = append(allSpecs, *pool)
				}
			}

			if _printer, ok := ctl.NodesPrinter[_defaultPrinter]; ok {
				_printer(&response.NodePool{
					Pools: allSpecs,
				}, _defaultStyler)
			}
		},
	}

	return _cmd
}

// CmdGetClustersList - command to get CNF Catalog entity
func (ctl *TcaCtl) CmdGetClustersList() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:   "info [optional cluster name]",
		Short: "Command returns kubernetes cluster or cluster information.",
		Long: templates.LongDesc(
			`Command returns kubernetes clusters or cluster information.`),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			ctx := context.Background()
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			clusters, err := ctl.tca.GetClusters(ctx)
			CheckErrLogError(err)
			// either get all or lookup by name
			if len(args) > 0 {
				cluster, err := clusters.GetClusterSpec(args[0])
				CheckErrLogError(err)
				if printer, ok := ctl.ClusterPrinter[_defaultPrinter]; ok {
					printer(cluster, _defaultStyler)
				}
			} else {
				if printer, ok := ctl.ClustersPrinter[_defaultPrinter]; ok {
					printer(clusters, _defaultStyler)
				}
			}
		},
	}

	//
	_cmd.Flags().BoolVarP(&_isWide,
		"wide", "w", true, "Wide output")
	return _cmd
}

// CmdGetClustersK8SConfig retrieve kubeconfig
// if active flag passed, will serialize to kubeconfig file
// if file indicated will save to a file
func (ctl *TcaCtl) CmdGetClustersK8SConfig() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		fileName        string
		activate        bool
	)

	var _cmd = &cobra.Command{
		Use:   "kubeconfig [cluster name]",
		Short: "Command returns cluster kubeconfig",
		Long: templates.LongDesc(
			`Command returns cluster kubeconfig.`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			clusters, err := ctl.tca.GetClusters(ctx)
			CheckErrLogError(err)
			for _, c := range clusters.Clusters {
				if strings.Contains(c.ClusterName, args[0]) {
					kubeconfig, err := b64.StdEncoding.DecodeString(c.KubeConfig)
					if err != nil {
						fmt.Println("Failed decode kubeconfig")
					}

					if activate {
						home := osutil.HomeDir()
						if home == "" {
							CheckErrLogError(errors.New("can't determine user home dir"))
						}

						defaultKubeconfig := filepath.Join(home,
							kubernetes.KUBEDIR, kubernetes.KUBEFILE)
						f, err := os.OpenFile(defaultKubeconfig,
							os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						CheckErrLogError(err)

						defer f.Close()
						if _, err := f.Write(kubeconfig); err != nil {
							log.Println(err)
						}
						continue
					}

					if len(fileName) == 0 {
						fmt.Println(string(kubeconfig))
						return
					}

					err = ioutil.WriteFile(fileName, kubeconfig, 0644)
					CheckErrLogError(err)
				}
			}
		},
	}

	_cmd.Flags().StringVarP(&fileName,
		"file_name", "f", "", "file to save.")

	_cmd.Flags().BoolVarP(&activate,
		"activate", "a", false, "set at active context.")

	return _cmd
}

// CmdCreateCluster -  command for cluster creation
// Read cluster spec , validate each spec param and create cluster
// if specs are valid, in Dry run resolve all name, parse spec
// and output final yaml if no error.
func (ctl *TcaCtl) CmdCreateCluster() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		isDry           = false
		doBlock         bool
		showProgress    bool
		showSpec        bool
	)

	var _cmd = &cobra.Command{
		Use:   "cluster [path to json or yaml file]",
		Short: "Command creates a management or tenant cluster",
		Long: `

The command creates cluster management or tenant cluster from input spec.  
Spec can be in YAML or JSON format. By default, the task is created asynchronously 
and doesn't wait for the task to finish.  The client can pass a block flag 
that will wait for the task to finish.' 
`,

		Example: "\t - tcactl create cluster examples/edge_mgmt_cluster.yaml --stderrthreshold INFO\n" +
			"\t - tcactl create cluster examples/edge_mgmt_cluster.yaml --block --progress --stderrthreshold INFO\n",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"cluster", "cl"},
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			// spec read from file
			var spec specs.SpecCluster
			if ioutils.FileExists(args[0]) {
				buffer, err := ioutil.ReadFile(args[0])
				CheckErrLogError(err)
				// first we try yaml if failed try json
				err = yaml.Unmarshal(buffer, &spec)
				if err != nil {
					// try json
					err = json.Unmarshal(buffer, &spec)
					CheckErrLogError(err)
				}
			} else {
				CheckErrLogError(fmt.Errorf("%v not found", args[0]))
			}

			if showSpec {
				err := ioutils.PrettyPrint(spec)
				CheckErrLogError(err)
				return
			}

			// otherwise create
			task, err := ctl.tca.CreateClusters(ctx, &api.ClusterCreateApiReq{
				Spec:          &spec,
				IsBlocking:    doBlock,
				IsDryRun:      isDry,
				IsVerbose:     showProgress,
				IsFixConflict: true})
			CheckErrLogError(err)

			if task != nil {
				fmt.Printf("SpecCluster created task create, task id %s\n", task.Id)
			}

			// dry run will output spec to screen after all validation and substitution
			if isDry {
				if printer, ok := ctl.ClusterRequestPrinter[_defaultPrinter]; ok {
					printer(&spec, _defaultStyler)
				}
			}
		},
	}

	_cmd.Flags().BoolVar(&isDry,
		CliDryRun, false, "Parses input template spec, "+
			"validates, outputs spec to the terminal screen. Format based on -o flag.")

	//
	_cmd.Flags().BoolVarP(&doBlock, CliBlock, "b", false,
		"Blocks and wait task to finish.")

	//
	_cmd.Flags().BoolVarP(&showProgress, CliProgress, "s", false,
		"Show task progress.")

	//
	_cmd.Flags().BoolVar(&showSpec, CliShow, false,
		"Show spec only.")

	return _cmd
}

// CmdDeleteCluster - command delete cluster
func (ctl *TcaCtl) CmdDeleteCluster() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		doBlock         bool
		showProgress    bool
	)

	var _cmd = &cobra.Command{
		Use:   "cluster [name or id of cluster]",
		Short: "Command delete cluster.",
		Long: templates.LongDesc(
			`Command deletes cluster.`),
		Example: "\t - tcactl delete cluster 794a675c-777a-47f4-8edb-36a686ef4065\n " +
			"\t -tcactl delete cluster mycluster",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			// global output type, and terminal wide or not
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			// otherwise create

			task, err := ctl.tca.DeleteCluster(ctx,
				&api.ClusterDeleteApiReq{
					Cluster:    args[0],
					IsBlocking: doBlock,
					IsVerbose:  showProgress,
				})
			if err != nil {
				return
			}
			if task != nil {
				fmt.Println("SpecCluster deleted.")
			}
		},
	}

	//
	_cmd.Flags().BoolVarP(&doBlock, CliBlock, "b", false,
		"Blocks and wait task to finish.")

	//
	_cmd.Flags().BoolVarP(&showProgress, CliProgress, "s", true,
		"Show task progress.")

	return _cmd
}

// CmdGetClusterTasks - command return current list of task.
func (ctl *TcaCtl) CmdGetClusterTasks() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
	)

	var _cmd = &cobra.Command{
		Use:     "tasks [cluster name or id]",
		Aliases: []string{"task"},
		Short:   "Command returns currently running task on a particular cluster.",
		Long: templates.LongDesc(`

Command returns currently running task on a particular cluster.`),

		Example: "- tcactl get cluster tasks 9411f70f-d24d-4842-ab56-b7214d",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			// global output type, and terminal wide or not
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			ctl.tca.SetTrace(ctl.IsTrace)

			task, err := ctl.tca.GetClusterTask(ctx, args[0], true)
			CheckErrLogError(err)

			if _printer, ok := ctl.TaskClusterPrinter[_defaultPrinter]; ok {
				_printer(task, _defaultStyler)
			}
		},
	}

	return _cmd
}

// CmdDescribeTask - command return cluster pools list
func (ctl *TcaCtl) CmdDescribeTask() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
	)

	var _cmd = &cobra.Command{
		Use:     "task [task_id]",
		Short:   "Command return current running task.",
		Long:    `Command return current running task.`,
		Example: "- tcactl desc task 9411f70f-d24d-4842-ab56-b7214d",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			// global output type, and terminal wide or not
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			task, err := ctl.tca.GetCurrentClusterTask(ctx, args[0])
			CheckErrLogError(err)

			if _printer, ok := ctl.TaskClusterPrinter[_defaultPrinter]; ok {
				_printer(task, _defaultStyler)
			}
		},
	}

	return _cmd
}
