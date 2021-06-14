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
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spyroot/hestia/cmd/client/respons"
	"github.com/spyroot/hestia/pkg/io"
)

// CmdGetPoolNodes - command to get CNF Catalog entity
func (ctl *TcaCtl) CmdGetPoolNodes() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:     "nodes",
		Short:   "Command returns kubernetes node pool",
		Long:    `Command returns a list kubernetes node pool for a given cluster name.`,
		Example: "tcactl get clusters pool 794a675c-777a-47f4-8edb-36a686ef4065",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// set wide or not
			_isWide, err := cmd.Flags().GetBool("wide")
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			clusters, err := ctl.TcaClient.GetClusters()
			if err != nil || clusters == nil {
				glog.Errorf("Failed retrieve cluster list %v", err)
				return
			}

			clusterId, err := clusters.GetClusterId(args[0])
			CheckErrLogError(err)

			pool, err := ctl.TcaClient.GetClusterNodePools(clusterId)
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

// CmdGetClustersPool - command return cluster pools list
func (ctl *TcaCtl) CmdGetClustersPool() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:     "pool [name or id of cluster]",
		Short:   "Command returns kubernetes node pool for a given cluster",
		Long:    `Command returns a list kubernetes node pool for a given cluster name.`,
		Example: "- tcactl get clusters pool 794a675c-777a-47f4-8edb-36a686ef4065\n -tcactl get cluster mycluster",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// set wide or not
			_isWide, err := cmd.Flags().GetBool("wide")
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			clusters, err := ctl.TcaClient.GetClusters()
			if err != nil || clusters == nil {
				glog.Errorf("Failed retrieve cluster list %v", err)
				return
			}

			clusterId, err := clusters.GetClusterId(args[0])
			CheckErrLogError(err)

			pool, err := ctl.TcaClient.GetClusterNodePools(clusterId)
			if err != nil {
				glog.Errorf("Failed retrieve node pools %v", err)
				return
			}
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

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
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

			cluster, err := ctl.TcaClient.GetCluster(_clusterId)
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

			clusters, err := ctl.TcaClient.GetClusters()
			if err != nil || clusters == nil {
				glog.Errorf("Failed retrieve cluster list %v", err)
				return
			}

			clusterId, err := clusters.GetClusterId(args[0])
			CheckErrLogError(err)

			pool, err := ctl.TcaClient.GetClusterNodePools(clusterId)
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
			_isWide, err := cmd.Flags().GetBool(CliWide)
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
			pool, err := ctl.TcaClient.GetClusterNodePool(_clusterId, _targetPoolID)
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

// CmdDescClusterNodePools - describe node pool
func (ctl *TcaCtl) CmdDescClusterNodePools() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:     "pools [name or id]",
		Short:   "Command describes all node pool",
		Long:    `Command describes all node pool.`,
		Example: "tcactl describe pools",
		//Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// set wide or not
			_isWide, err := cmd.Flags().GetBool(CliWide)
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			clusters, err := ctl.TcaClient.GetClusters()
			CheckErrLogError(err)

			var allSpecs []respons.NodesSpecs
			for _, c := range clusters.Clusters {
				pools, err := ctl.TcaClient.GetClusterNodePools(c.Id)
				CheckErrLogInfoMsg(err)
				for _, p := range pools.Pools {
					pool, err := ctl.TcaClient.GetClusterNodePool(c.Id, p.Id)
					CheckErrLogInfoMsg(err)
					allSpecs = append(allSpecs, *pool)
				}
			}

			if _printer, ok := ctl.NodePoolPrinter[_defaultPrinter]; ok {
				_printer(&respons.NodePool{
					Pools: allSpecs,
				}, _defaultStyler)
			}
		},
	}

	// wide output
	_cmd.Flags().BoolVarP(&_isWide,
		"wide", "w", true, "Wide output")

	return _cmd
}

// CmdDescClusterNodes  - describe node pool
func (ctl *TcaCtl) CmdDescClusterNodes() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_isWide         = false
	)

	var _cmd = &cobra.Command{
		Use:     "nodes",
		Short:   "Command describes all kubernetes nodes",
		Long:    `Command describes all kubernetes nodes`,
		Example: "tcactl describe nodes",
		//Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// set wide or not
			_isWide, err := cmd.Flags().GetBool(CliWide)
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			clusters, err := ctl.TcaClient.GetClusters()
			CheckErrLogError(err)

			var allSpecs []respons.NodesSpecs
			for _, c := range clusters.Clusters {
				pools, err := ctl.TcaClient.GetClusterNodePools(c.Id)
				CheckErrLogError(err)
				for _, p := range pools.Pools {
					pool, err := ctl.TcaClient.GetClusterNodePool(c.Id, p.Id)
					CheckErrLogError(err)
					allSpecs = append(allSpecs, *pool)
				}
			}

			if _printer, ok := ctl.NodesPrinter[_defaultPrinter]; ok {
				_printer(&respons.NodePool{
					Pools: allSpecs,
				}, _defaultStyler)
			}
		},
	}

	// wide output
	_cmd.Flags().BoolVarP(&_isWide,
		"wide", "w", true, "Wide output")

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
		Short: "Command returns kubernetes cluster or cluster information",
		Long:  `Command returns kubernetes cluster or cluster information.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()
			// set wide or not
			_isWide, err := cmd.Flags().GetBool("wide")
			io.CheckErr(err)
			_defaultStyler.SetWide(_isWide)

			clusters, err := ctl.TcaClient.GetClusters()
			CheckErrLogError(err)

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
