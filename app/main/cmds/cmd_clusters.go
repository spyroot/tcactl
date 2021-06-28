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
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/lib/api/kubernetes"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/pkg/io"
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
		Long: `
Command retrieves cluster-related information. Each sub-command 
require either cluster name or cluster id.
`,
		Example: "- tcactl get clusters info\n - tcactl get clusters pool edge",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"cluster", "cl"},
		Run: func(cmd *cobra.Command, args []string) {
			return
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
		Long:  `Command returns a list kubernetes node pool for a given cluster name.`,
		Example: "\t - tcactl get clusters pool 794a675c-777a-47f4-8edb-36a686ef4065\n " +
			"\t - tcactl get cluster mycluster",
		Run: func(cmd *cobra.Command, args []string) {

			var (
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
				pool, err = ctl.tca.GetNodePool(args[0])
			}
			pool, err = ctl.tca.GetAllNodePool()
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
		Use:     "cluster [name or id]",
		Short:   "Command describes kubernetes cluster or clusters information",
		Long:    `Command returns kubernetes cluster or cluster information.`,
		Example: "- tcactl describe clusters 794a675c-777a-47f4-8edb-36a686ef4065 -o json",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var (
				_clusterId string
				err        error
			)

			_clusterId = args[0]
			if !IsValidUUID(args[0]) {
				_clusterId, err = ctl.ResolveClusterName(args[0])
				CheckErrLogError(err)
			}

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()
			// set wide or not
			_isWide, err := cmd.Flags().GetBool("wide")
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			cluster, err := ctl.tca.GetCluster(_clusterId)
			CheckErrLogError(err)

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
		Use:     "nodes",
		Short:   "Command returns kubernetes node pool",
		Long:    `Command returns a list kubernetes node pool for a given cluster name.`,
		Example: "- tcactl get clusters nodes 794a675c-777a-47f4-8edb-36a686ef4065\n - tcactl get clusters nodes edge",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// set wide or not
			_isWide, err := cmd.Flags().GetBool("wide")
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			clusters, err := ctl.tca.GetClusters()
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
		Use:     "pool [name or id]",
		Short:   "Command describes kubernetes node pool",
		Long:    `Command describes kubernetes node pool for a given or default cluster name.`,
		Example: "tcactl describe pool 794a675c-777a-47f4-8edb-36a686ef4065",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

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
			pool, err := ctl.tca.GetClusterNodePool(_clusterId, _targetPoolID)
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

			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			clusters, err := ctl.tca.GetClusters()
			CheckErrLogError(err)

			var allSpecs []response.NodesSpecs
			for _, c := range clusters.Clusters {
				pools, err := ctl.tca.GetClusterNodePools(c.Id)
				CheckErrLogError(err)
				for _, p := range pools.Pools {
					pool, err := ctl.tca.GetClusterNodePool(c.Id, p.Id)
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
		Long:  `Command returns kubernetes cluster or cluster information.`,
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			clusters, err := ctl.tca.GetClusters()
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
		Long:  `Command returns cluster kubeconfig.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			clusters, err := ctl.tca.GetClusters()
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
	)

	var _cmd = &cobra.Command{
		Use:   "clusters [path to file]",
		Short: "Command create cluster",
		Long: `

Command create cluster from input spec. Spec can be in yaml or json format.`,

		Example: "- ./tcactl create cluster examples/edge_mgmt_cluster.yaml --stderrthreshold INFO",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"cluster", "cl"},
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			// spec read from file
			var spec request.Cluster
			if io.FileExists(args[0]) {
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

			// otherwise create
			ok, err := ctl.tca.CreateClusters(&spec, isDry)
			CheckErrLogError(err)
			if ok {
				fmt.Println("Cluster created.")
			}

			// dry run will output template to screen after parser
			// resolved all name to id.
			if isDry {
				if printer, ok := ctl.ClusterRequestPrinter[_defaultPrinter]; ok {
					printer(&spec, _defaultStyler)
				}
			}
		},
	}

	_cmd.Flags().BoolVar(&isDry,
		"dry", false, "Parses input template spec, "+
			"validates, outputs spec to the terminal screen. Format based on -o flag.")

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

			// global output type, and terminal wide or not
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			task, err := ctl.tca.GetCurrentClusterTask(args[0])
			CheckErrLogError(err)

			if _printer, ok := ctl.TaskClusterPrinter[_defaultPrinter]; ok {
				_printer(task, _defaultStyler)
			}
		},
	}

	return _cmd
}

// CmdDeleteCluster - command delete cluster
func (ctl *TcaCtl) CmdDeleteCluster() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
	)

	var _cmd = &cobra.Command{
		Use:     "cluster [name or id of cluster]",
		Short:   "Command delete cluster.",
		Long:    `Command Command delete cluster.`,
		Example: "- tcactl delete cluster 794a675c-777a-47f4-8edb-36a686ef4065\n -tcactl delete cluster mycluster",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type, and terminal wide or not
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			ok, err := ctl.tca.DeleteCluster(args[0])
			if err != nil {
				return
			}
			if ok {
				fmt.Println("Cluster deleted.")
			}
		},
	}

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
		Long: `

Command returns currently running task on a particular cluster.`,

		Example: "- tcactl get cluster tasks 9411f70f-d24d-4842-ab56-b7214d",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type, and terminal wide or not
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			ctl.tca.SetTrace(ctl.IsTrace)

			task, err := ctl.tca.GetClusterTask(args[0], true)
			CheckErrLogError(err)

			if _printer, ok := ctl.TaskClusterPrinter[_defaultPrinter]; ok {
				_printer(task, _defaultStyler)
			}
		},
	}

	return _cmd
}
