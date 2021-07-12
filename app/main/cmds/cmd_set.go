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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
	"github.com/spyroot/tcactl/pkg/io"
)

// CmdSetTca - return list of cloud provider attached to TCA
func (ctl *TcaCtl) CmdSetTca() *cobra.Command {

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:   "api",
		Short: "Command sets tca end-point and saves config.",
		Long:  templates.LongDesc(`Command sets tca end-point and saves config.`),
		Example: templates.LongDesc(
			"tcactl set https://tca-vip03.cnfdemo.io"),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set(ConfigTcaEndpoint, args[0])
			err := viper.WriteConfig()
			io.CheckErr(err)
		},
	}
	return _cmd
}

// CmdSetNodePool - return list of cloud provider attached to TCA
func (ctl *TcaCtl) CmdSetNodePool() *cobra.Command {

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "nodepool",
		Short:   "Command sets default node pool end-point and saves config.",
		Long:    `Command sets default node pool end-point and saves config.`,
		Example: "tcactl set nodepool mypool",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set(ConfigNodePool, args[0])
			err := viper.WriteConfig()
			io.CheckErr(err)
		},
	}
	return _cmd
}

// CmdSetCluster - return list of cloud provider attached to TCA
func (ctl *TcaCtl) CmdSetCluster() *cobra.Command {

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:   "cluster",
		Short: "Command sets cluster and saves config.",
		Long: templates.LongDesc(
			`Command sets cluster and saves config.`),
		Example: "tcactl set cluster mycluster",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set(ConfigDefaultCluster, args[0])
			err := viper.WriteConfig()
			io.CheckErr(err)
		},
	}
	return _cmd
}

// CmdSetUsername - return list of cloud provider attached to TCA
func (ctl *TcaCtl) CmdSetUsername() *cobra.Command {

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:   "username",
		Short: "Command sets TCA username and saves config.",
		Long: templates.LongDesc(
			`Command sets TCA username and saves config.`),
		Example: "tcactl set username administrator@vsphere.local",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set(ConfigTcaUsername, args[0])
			err := viper.WriteConfig()
			io.CheckErr(err)
		},
	}
	return _cmd
}

// CmdSetPassword - return list of cloud provider attached to TCA
func (ctl *TcaCtl) CmdSetPassword() *cobra.Command {

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:   "password",
		Short: "Command sets cluster username and saves config.",
		Long: templates.LongDesc(
			`Command sets cluster username and saves config.`),
		Example: "tcactl set password mypass",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set(ConfigTcaPassword, args[0])
			err := viper.WriteConfig()
			io.CheckErr(err)
		},
	}
	return _cmd
}
