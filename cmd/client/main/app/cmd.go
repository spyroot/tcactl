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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spyroot/hestia/pkg/io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CmdInitConfig - initialize configuration file, for initial
// setup TCA and other defaults
func (ctl *TcaCtl) CmdInitConfig() *cobra.Command {

	var _cmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes config file.",
		Long:  `Initializes config file.`,

		Run: func(cmd *cobra.Command, args []string) {

			home, err := homedir.Dir()
			configPrefix := ".tcactl"
			configName := "config"
			configType := "yaml"
			configHome := filepath.Join(home, "/", configPrefix)
			configPath := filepath.Join(configHome, configName+"."+configType)

			_, err = ioutil.ReadDir(home)
			io.CheckErr(err)
			err = os.MkdirAll(configHome, 0755)
			io.CheckErr(err)

			_, err = os.Stat(configPath)
			if !os.IsExist(err) {
				if _, err := os.Create(configPath); err != nil {
					io.CheckErr(err)
				}
			}

			err = viper.WriteConfig()
			io.CheckErr(err)
		},
	}

	return _cmd
}

// CmdSaveConfig - save config file
func (ctl *TcaCtl) CmdSaveConfig() *cobra.Command {

	var _cmd = &cobra.Command{
		Use:   "save",
		Short: "Saves config variables to config file.",
		Long:  `Saves config variables to config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := viper.WriteConfig()
			if err != nil {
				io.CheckErr(err)
				return
			}
			io.CheckErr(err)
		},
	}

	return _cmd
}

// BuildCmd build all commands and attach to root cmd
func (ctl *TcaCtl) BuildCmd() {

	var describe = &cobra.Command{
		Use:   "describe",
		Short: "Describe TCA object details",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				return
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// root cmd for all get
	var cmdGet = &cobra.Command{
		Use:   "get [cnfi, cnfc, clusters, pools]",
		Short: "Gets object from TCA, cnfi, cnfc etc",
		Long:  `Gets object from TCA. CNFI is CNFI in the inventory, CNFC Catalog entities.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				return
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// root cmd for all update commands
	var cmdUpdate = &cobra.Command{
		Use:   "update [cnfi or cnfc]",
		Short: "Updates cnf, cnf catalog etc",
		Long:  `Updates cnf state, (terminate, scale etc), or CNF catalog`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				return
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// create
	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a new object in TCA.",
		Long:  `Create a new object in TCA. For example new CNF instance.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				return
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// create
	var CmdsetCreate = &cobra.Command{
		Use:   "set",
		Short: "Create a new object in TCA.",
		Long:  `Create a new object in TCA. For example new CNF instance.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				return
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// root describe
	describe.AddCommand(ctl.CmdGetCluster())
	describe.AddCommand(ctl.CmdDescClusterNodePool())
	describe.AddCommand(ctl.CmdDescClusterNodePools())
	describe.AddCommand(ctl.CmdDescClusterNodes())

	// update sub-commands
	cmdUpdate.AddCommand(ctl.CmdTerminateInstances())
	cmdUpdate.AddCommand(ctl.CmdUpdateInstances())
	cmdUpdate.AddCommand(ctl.CmdRollbackInstances())
	ctl.RootCmd.AddCommand(describe, cmdGet, cmdUpdate, cmdCreate, CmdsetCreate,
		ctl.CmdSaveConfig(), ctl.CmdInitConfig())

	// get command
	cmdGet.AddCommand(ctl.CmdGetPackages())
	cmdGet.AddCommand(ctl.CmdGetInstances())
	cmdGet.AddCommand(ctl.CmdGetRepos())
	cmdGet.AddCommand(ctl.CmdGetClouds())
	cmdGet.AddCommand(ctl.CmdGetPool())
	cmdGet.AddCommand(ctl.CmdGetClusters())
	cmdGet.AddCommand(ctl.CmdGetVdu())
	cmdGet.AddCommand(ctl.CmdGetExtensions())

	// add root command
	cmdCreate.AddCommand(ctl.CmdCreateCnf())
}
