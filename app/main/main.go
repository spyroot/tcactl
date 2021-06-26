// Package main
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
package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spyroot/tcactl/app/main/cmds"
	"github.com/spyroot/tcactl/pkg/io"
	_ "github.com/spyroot/tcactl/pkg/io"
	"net/url"
	"strings"

	//	pflag "github.com/spf13/pflag"
	//flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spyroot/tcactl/lib/client/request"
	"os"
)

var (
	// tcactl main class
	tcaCtl      = cmds.NewTcaCtl()
	userLicense string
)

// Inits logger
func init() {

	// builds all command tree
	tcaCtl.BuildCmd()

	// default values in case we don't have config
	viper.SetDefault(cmds.ConfigTcaEndpoint, "https://localhost")
	viper.SetDefault(cmds.ConfigTcaUsername, "administrator@vsphere.local")
	viper.SetDefault(cmds.ConfigTcaPassword, "VMware1!")
	viper.SetDefault(cmds.ConfigDefaultCluster, "default")
	viper.SetDefault(cmds.ConfigDefaultCloud, "default")
	viper.SetDefault(cmds.ConfigNodePool, "default")
	viper.SetDefault(cmds.ConfigStderrThreshold, "INFO")
	viper.SetDefault(cmds.ConfigHarborEndpoint, "https://repo.vmware.com")
	viper.SetDefault(cmds.ConfigHarborUsername, "admin")
	viper.SetDefault(cmds.ConfigHarborPassword, "VMware1!")
	viper.SetDefault(cmds.ConfigRepoName, "repo.vmware.com")

	cobra.OnInitialize(initConfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	//	pflag.Parse()
	err := flag.CommandLine.Parse([]string{})
	io.CheckErr(err)

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.Printer,
		cmds.FlagOutput, "o", "default",
		"Output format json, yaml. (default console)")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.CfgFile,
		cmds.FlagConfig, "c", "",
		"config file (default is $HOME/.tcacli/config.yaml)")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.DefaultCloudName,
		cmds.ConfigDefaultCloud, "p", "",
		"Default Cloud Provider used by tcactl.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.DefaultClusterName,
		cmds.ConfigDefaultCluster, "m", "",
		"Default Tenant Cluster Name.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.DefaultRepoName,
		cmds.ConfigRepoName, "z", "",
		"Default Cloud Provider used by tcactl.")

	tcaCtl.RootCmd.PersistentFlags().StringVar(&tcaCtl.DefaultNodePoolName,
		cmds.ConfigNodePool, "",
		"Default node pool to use.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.Harbor,
		cmds.ConfigHarborEndpoint, "r", "",
		"Harbor API end-point url.")

	tcaCtl.RootCmd.PersistentFlags().StringVar(&tcaCtl.HarborUsername,
		cmds.ConfigHarborUsername, "",
		"Harbor username.")

	tcaCtl.RootCmd.PersistentFlags().StringVar(&tcaCtl.HarborPassword,
		cmds.ConfigHarborPassword, "",
		"Harbor password.")

	tcaCtl.RootCmd.PersistentFlags().BoolVarP(&tcaCtl.IsColorTerm,
		cmds.FlagCliTerm, "t", false, "Disables color output.")

	tcaCtl.RootCmd.PersistentFlags().BoolVarP(&tcaCtl.IsWideTerm,
		cmds.FlagCliWide, "w", false, "Wide terminal output.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&userLicense,
		"license", "l", "", "name of license for the project")

	tcaCtl.RootCmd.PersistentFlags().Bool(
		"viper", true, "use Viper for configuration")

	err = viper.BindPFlag("useViper",
		tcaCtl.RootCmd.PersistentFlags().Lookup("viper"))
	io.CheckErr(err)

	err = viper.BindPFlag(cmds.FlagOutput,
		tcaCtl.RootCmd.PersistentFlags().Lookup(cmds.FlagOutput))
	io.CheckErr(err)

	err = viper.BindPFlag(cmds.ConfigDefaultCloud,
		tcaCtl.RootCmd.PersistentFlags().Lookup(cmds.ConfigDefaultCloud))
	io.CheckErr(err)

	err = viper.BindPFlag(cmds.ConfigDefaultCluster,
		tcaCtl.RootCmd.PersistentFlags().Lookup(cmds.ConfigDefaultCluster))
	io.CheckErr(err)

	err = viper.BindPFlag(cmds.ConfigNodePool,
		tcaCtl.RootCmd.PersistentFlags().Lookup(cmds.ConfigNodePool))
	io.CheckErr(err)

	err = viper.BindPFlag(cmds.ConfigRepoName,
		tcaCtl.RootCmd.PersistentFlags().Lookup(cmds.ConfigRepoName))
	io.CheckErr(err)

	err = viper.BindPFlag(cmds.ConfigHarborEndpoint,
		tcaCtl.RootCmd.PersistentFlags().Lookup(cmds.ConfigHarborEndpoint))
	io.CheckErr(err)

	err = viper.BindPFlag(cmds.ConfigHarborUsername,
		tcaCtl.RootCmd.PersistentFlags().Lookup(cmds.ConfigHarborUsername))
	io.CheckErr(err)

	err = viper.BindPFlag(cmds.ConfigHarborPassword,
		tcaCtl.RootCmd.PersistentFlags().Lookup(cmds.ConfigHarborPassword))
	io.CheckErr(err)

	viper.SetDefault("author", "spyroot@gmail.com")
	viper.SetDefault("license", "apache")
}

// IsUrl return true if str is URL
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// initConfig - read tcactl configs
func initConfig() {

	// default values in case we don't have config
	viper.AutomaticEnv()

	if tcaCtl.CfgFile != "" {
		viper.SetConfigFile(tcaCtl.CfgFile)
	} else {
		// or search in default location.
		viper.SetConfigName(cmds.ConfigFile)
		viper.SetConfigType(cmds.ConfigFormat)
		viper.AddConfigPath("$HOME/.tcactl")
		viper.AddConfigPath("/etc/tcactl/")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// update default after we read
	glog.Infof("Using tca endpoint %v", viper.GetString(cmds.ConfigTcaEndpoint))
	viper.GetString(cmds.ConfigTcaEndpoint)
	ok := IsUrl(viper.GetString(cmds.ConfigTcaEndpoint))
	if !ok {
		io.CheckErr("Invalid tca url")
	}

	if !strings.HasPrefix(viper.GetString(cmds.ConfigTcaEndpoint), "https") {
		io.CheckErr(fmt.Errorf("please indicate https protocol type"))
	}

	glog.Infof("Using harbor endpoint %v", viper.GetString(cmds.ConfigHarborEndpoint))
	ok = IsUrl(viper.GetString(cmds.ConfigHarborEndpoint))
	if !ok {
		io.CheckErr("Invalid harbor url")
	}

	ok = IsUrl(viper.GetString(cmds.ConfigRepoName))
	if !ok {
		io.CheckErr("Invalid repo name")
	}

	tcaCtl.SetTcaBase(viper.GetString(cmds.ConfigTcaEndpoint))
	tcaCtl.SetTcaUsername(viper.GetString(cmds.ConfigTcaUsername))
	tcaCtl.SetPassword(viper.GetString(cmds.ConfigTcaPassword))

	//tcaCtl.TcaClient.BaseURL = viper.GetString(cmds.ConfigTcaEndpoint)
	//tcaCtl.TcaClient.Username = viper.GetString(cmds.ConfigTcaUsername)
	//tcaCtl.TcaClient.Password = viper.GetString(cmds.ConfigTcaPassword)

	// default Cloud in TCA,  Cluster and node pool
	tcaCtl.DefaultCloudName = viper.GetString(cmds.ConfigDefaultCloud)
	tcaCtl.DefaultClusterName = viper.GetString(cmds.ConfigDefaultCluster)
	tcaCtl.DefaultNodePoolName = viper.GetString(cmds.ConfigNodePool)
	tcaCtl.DefaultRepoName = viper.GetString(cmds.ConfigRepoName)

	tcaCtl.Harbor = viper.GetString(cmds.ConfigHarborEndpoint)
	tcaCtl.HarborUsername = viper.GetString(cmds.ConfigHarborUsername)
	tcaCtl.HarborPassword = viper.GetString(cmds.ConfigHarborPassword)

	tcaCtl.Printer = viper.GetString("output")
	glog.Infof("TCA Base set to %v", viper.GetString(cmds.ConfigTcaEndpoint))

}

// @title VMware TCA CTL proto API
// @version 1.0
// @description MWC proto API server.
// @termsOfService http://swagger.io/terms/
func main() {

	glog.Infof("Using config file %v", tcaCtl.CfgFile)
	if err := tcaCtl.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		glog.Error(err)
		os.Exit(1)
	}
}
