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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spyroot/tcactl/pkg/io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	CliBlock       = "block"
	CliPool        = "pool"
	CliDisableGran = "grant"
	CliForce       = "force"
)

// CmdInitConfig - initialize configuration file, for initial
// setup TCA and other defaults
func (ctl *TcaCtl) CmdInitConfig() *cobra.Command {

	var _cmd = &cobra.Command{
		Use:   "init",
		Short: "Command initializes default config file.",
		Long:  `Command Initializes default config file.`,

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

			fmt.Println("Default config file generated: ", configPath)
			fmt.Println("Now run tcactl set and set username, password and TCA Cluster Endpoint.")
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

func (ctl *TcaCtl) CmdCreate() *cobra.Command {
	// cnf instances

	var cmdCreate = &cobra.Command{
		Use:   "Create",
		Short: "Terminate CNF instance",
		Long:  `Terminate CNF instance, caller need to provide CNF Identifier.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmdCreate.AddCommand(ctl.CmdCreateCnf())
	return cmdCreate
}

// BuildCmd build all commands and attaches to root cmd
// in case you need add sub-command you can, add to plugin dir.
//(TODO) add dynamic loading plugin
func (ctl *TcaCtl) BuildCmd() {

	var describe = &cobra.Command{
		Use:     "describe [cloud or cluster or nodes or pool or template]",
		Long:    `Command describes object TCA. CNFI is CNFI in the inventory, CNFC Catalog entities.`,
		Short:   "Describe TCA object details",
		Aliases: []string{"desc"},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				CheckErrLogError(err)
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
				CheckErrLogError(err)
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
				CheckErrLogError(err)
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
				CheckErrLogError(err)
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// set
	var cmdSet = &cobra.Command{
		Use:   "set",
		Short: "Command sets config variables (Username, Password etc) for tcactl.",
		Long:  `Command sets config variables (Username, Password etc) for tcactl.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// delete root command
	var cmdDelete = &cobra.Command{
		Use:     "delete",
		Short:   "Command deletes object (template,cluster,cnf etc).",
		Long:    `Command deletes object (template,cluster,cnf etc).`,
		Aliases: []string{"del"},
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

	// Set root command
	cmdSet.AddCommand(
		ctl.CmdSetTca(),
		ctl.CmdSetCluster(),
		ctl.CmdSetNodePool(),
		ctl.CmdSetUsername(),
		ctl.CmdSetPassword())

	// Describe sub command
	describe.AddCommand(ctl.CmdDescribeVim())
	describe.AddCommand(ctl.CmdGetCluster())
	describe.AddCommand(ctl.CmdDescClusterNodePool())
	describe.AddCommand(ctl.CmdDescClusterNodePools())
	describe.AddCommand(ctl.CmdDescClusterNodes())
	describe.AddCommand(ctl.CmdDescribeTemplate())
	describe.AddCommand(ctl.CmdDescribeTask())

	// Update sub-commands
	cmdUpdate.AddCommand(ctl.CmdUpdateClusterTemplates())
	cmdUpdate.AddCommand(ctl.CmdUpdateInstance())

	// root command
	ctl.RootCmd.AddCommand(describe,
		cmdGet,
		cmdUpdate,
		cmdCreate,
		cmdDelete,
		cmdSet,
		ctl.CmdSaveConfig(),
		ctl.CmdInitConfig())

	// Get command
	cmdGet.AddCommand(ctl.CmdGetPackages())
	cmdGet.AddCommand(ctl.CmdGetInstances())
	cmdGet.AddCommand(ctl.CmdGetRepos())
	cmdGet.AddCommand(ctl.CmdGetClouds())
	cmdGet.AddCommand(ctl.CmdGetPool())
	cmdGet.AddCommand(ctl.CmdGetClusters())
	cmdGet.AddCommand(ctl.CmdGetVdu())
	cmdGet.AddCommand(ctl.CmdGetExtensions())
	cmdGet.AddCommand(ctl.CmdGetClusterTemplates())
	cmdGet.AddCommand(ctl.CmdGetVim())

	// Create root command
	cmdCreate.AddCommand(
		ctl.CmdCreateCluster(),
		ctl.CmdCreateCnf(),
		ctl.CmdCreateClusterTemplates(),
		ctl.CmdCreatePackage())

	// Delete
	cmdDelete.AddCommand(
		ctl.CmdDeleteClusterTemplates(),
		ctl.CmdDeleteCluster(),
		ctl.CmdDeleteTenantCluster(),
		ctl.CmdDeleteInstances())
}
