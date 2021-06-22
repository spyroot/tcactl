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
	"github.com/spf13/viper"
	"github.com/spyroot/hestia/cmd/client/main/app"
	_ "github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/csar"
	"github.com/spyroot/hestia/pkg/io"
	_ "github.com/spyroot/hestia/pkg/io"
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

	tcaCtl.RootCmd.PersistentFlags().StringVar(&tcaCtl.DefaultClusterName,
		app.ConfigDefaultCloud, "",
		"Default cloud provider.")

	tcaCtl.RootCmd.PersistentFlags().StringVarP(&tcaCtl.DefaultClusterName,
		app.ConfigDefaultCluster, "k", "",
		"Default Cluster.")

	tcaCtl.RootCmd.PersistentFlags().StringVar(&tcaCtl.DefaultClusterName,
		app.ConfigNodePool, "",
		"Default node pool to use.")

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
	tcaCtl.TcaClient.BaseURL = viper.GetString(app.ConfigTcaEndpoint)
	tcaCtl.TcaClient.Username = viper.GetString(app.ConfigTcaUsername)
	tcaCtl.TcaClient.Password = viper.GetString(app.ConfigTcaPassword)

	// default Cloud in TCA,  Cluster and node pool
	tcaCtl.DefaultCloudName = viper.GetString(app.ConfigDefaultCloud)
	tcaCtl.DefaultClusterName = viper.GetString(app.ConfigDefaultCluster)
	tcaCtl.DefaultNodePoolName = viper.GetString(app.ConfigNodePool)

	tcaCtl.Printer = viper.GetString("output")
	glog.Infof("TCA Base set to %v", viper.GetString(app.ConfigTcaEndpoint))
}

// @title VMware TCA CTL proto API
// @version 1.0
// @description MWC proto API server.
// @termsOfService http://swagger.io/terms/
func main() {

	var substitution = map[string]string{}
	substitution["descriptorId"] = "nfd_1234"

	csar.ApplyTransformation(
		"/Users/spyroot/go/src/hestia/tests/smokeping-cnf.csar",
		"NFD.yaml", csar.NfdYamlPropertyTransformer, substitution)
}
