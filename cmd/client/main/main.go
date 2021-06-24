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
	"github.com/spyroot/tcactl/pkg/io"
	_ "github.com/spyroot/tcactl/pkg/io"
	"net/url"
	"strings"

	//	pflag "github.com/spf13/pflag"
	//flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/spyroot/tcactl/cmd/client/main/app"

	_ "github.com/spyroot/tcactl/cmd/client/request"
	"os"
)

var (
	tcaCtl      = app.NewTcaCtl()
	userLicense string
)

// Inits logger
func init() {

	// builds all command tree
	tcaCtl.BuildCmd()

	// default values in case we don't have config
	viper.SetDefault(app.ConfigTcaEndpoint, "https://localhost")
	viper.SetDefault(app.ConfigTcaUsername, "administrator@vsphere.local")
	viper.SetDefault(app.ConfigTcaPassword, "VMware1!")
	viper.SetDefault(app.ConfigDefaultCluster, "default")
	viper.SetDefault(app.ConfigDefaultCloud, "default")
	viper.SetDefault(app.ConfigNodePool, "default")
	viper.SetDefault(app.ConfigStderrThreshold, "INFO")
	viper.SetDefault(app.ConfigHarborEndpoint, "repo.vmware.com")
	viper.SetDefault(app.ConfigHarborUsername, "admin")
	viper.SetDefault(app.ConfigHarborPassword, "VMware1!")
	viper.SetDefault(app.ConfigRepoName, "repo.vmware.com")

	cobra.OnInitialize(initConfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	//	pflag.Parse()
	err := flag.CommandLine.Parse([]string{})
	io.CheckErr(err)

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.Printer,
		app.FlagOutput, "o", "default",
		"Output format json, yaml. (default console)")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.CfgFile,
		app.FlagConfig, "c", "",
		"config file (default is $HOME/.tcacli/config.yaml)")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.DefaultCloudName,
		app.ConfigDefaultCloud, "p", "",
		"Default Cloud Provider used by tcactl.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.DefaultClusterName,
		app.ConfigDefaultCluster, "m", "",
		"Default Tenant Cluster Name.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.DefaultRepoName,
		app.ConfigRepoName, "z", "",
		"Default Cloud Provider used by tcactl.")

	tcaCtl.RootCmd.PersistentFlags().StringVar(&tcaCtl.DefaultNodePoolName,
		app.ConfigNodePool, "",
		"Default node pool to use.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.Harbor,
		app.ConfigHarborEndpoint, "r", "",
		"Harbor API end-point url.")

	tcaCtl.RootCmd.PersistentFlags().StringVar(&tcaCtl.HarborUsername,
		app.ConfigHarborUsername, "",
		"Harbor username.")

	tcaCtl.RootCmd.PersistentFlags().StringVar(&tcaCtl.HarborPassword,
		app.ConfigHarborPassword, "",
		"Harbor password.")

	tcaCtl.RootCmd.PersistentFlags().BoolVarP(&tcaCtl.IsColorTerm,
		app.FlagCliTerm, "t", false, "Disables color output.")

	tcaCtl.RootCmd.PersistentFlags().BoolVarP(&tcaCtl.IsWideTerm,
		app.FlagCliWide, "w", false, "Wide terminal output.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&userLicense,
		"license", "l", "", "name of license for the project")

	tcaCtl.RootCmd.PersistentFlags().Bool(
		"viper", true, "use Viper for configuration")

	err = viper.BindPFlag("useViper", tcaCtl.RootCmd.PersistentFlags().Lookup("viper"))
	io.CheckErr(err)

	err = viper.BindPFlag(app.FlagOutput, tcaCtl.RootCmd.PersistentFlags().Lookup(app.FlagOutput))
	io.CheckErr(err)

	err = viper.BindPFlag(app.ConfigDefaultCloud, tcaCtl.RootCmd.PersistentFlags().Lookup(app.ConfigDefaultCloud))
	io.CheckErr(err)

	err = viper.BindPFlag(app.ConfigDefaultCluster, tcaCtl.RootCmd.PersistentFlags().Lookup(app.ConfigDefaultCluster))
	io.CheckErr(err)

	err = viper.BindPFlag(app.ConfigNodePool, tcaCtl.RootCmd.PersistentFlags().Lookup(app.ConfigNodePool))
	io.CheckErr(err)

	err = viper.BindPFlag(app.ConfigRepoName, tcaCtl.RootCmd.PersistentFlags().Lookup(app.ConfigRepoName))
	io.CheckErr(err)

	err = viper.BindPFlag(app.ConfigHarborEndpoint, tcaCtl.RootCmd.PersistentFlags().Lookup(app.ConfigHarborEndpoint))
	io.CheckErr(err)

	err = viper.BindPFlag(app.ConfigHarborUsername, tcaCtl.RootCmd.PersistentFlags().Lookup(app.ConfigHarborUsername))
	io.CheckErr(err)

	err = viper.BindPFlag(app.ConfigHarborPassword, tcaCtl.RootCmd.PersistentFlags().Lookup(app.ConfigHarborPassword))
	io.CheckErr(err)

	viper.SetDefault("author", "spyroot@gmail.com")
	viper.SetDefault("license", "apache")
}

// initConfig - read tcactl configs
func initConfig() {

	// default values in case we don't have config
	viper.AutomaticEnv()

	if tcaCtl.CfgFile != "" {
		viper.SetConfigFile(tcaCtl.CfgFile)
	} else {
		// or search in default location.
		viper.SetConfigName(app.ConfigFile)
		viper.SetConfigType(app.ConfigFormat)
		viper.AddConfigPath("$HOME/.tcactl")
		viper.AddConfigPath("/etc/tcactl/")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// update default after we read
	viper.GetString(app.ConfigTcaEndpoint)
	_, err := url.ParseRequestURI(viper.GetString(app.ConfigTcaEndpoint))
	if err != nil {
		io.CheckErr(err)
	}

	if !strings.HasPrefix(viper.GetString(app.ConfigTcaEndpoint), "https") {
		io.CheckErr(fmt.Errorf("please indicate https protocol type"))
	}

	//_, err = url.ParseRequestURI(viper.GetString(app.ConfigHarborEndpoint))
	//if err != nil {
	//	io.CheckErr(err)
	//}

	//
	//_, err = url.ParseRequestURI(viper.GetString(app.ConfigRepoName))
	//if err != nil {
	//	io.CheckErr(err)
	//}

	tcaCtl.TcaClient.BaseURL = viper.GetString(app.ConfigTcaEndpoint)
	tcaCtl.TcaClient.Username = viper.GetString(app.ConfigTcaUsername)
	tcaCtl.TcaClient.Password = viper.GetString(app.ConfigTcaPassword)

	// default Cloud in TCA,  Cluster and node pool
	tcaCtl.DefaultCloudName = viper.GetString(app.ConfigDefaultCloud)
	tcaCtl.DefaultClusterName = viper.GetString(app.ConfigDefaultCluster)
	tcaCtl.DefaultNodePoolName = viper.GetString(app.ConfigNodePool)
	tcaCtl.DefaultRepoName = viper.GetString(app.ConfigRepoName)

	tcaCtl.Harbor = viper.GetString(app.ConfigHarborEndpoint)
	tcaCtl.HarborUsername = viper.GetString(app.ConfigHarborUsername)
	tcaCtl.HarborPassword = viper.GetString(app.ConfigHarborPassword)

	tcaCtl.Printer = viper.GetString("output")
	glog.Infof("TCA Base set to %v", viper.GetString(app.ConfigTcaEndpoint))

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
