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

package cmds

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
	"github.com/spyroot/tcactl/pkg/io"
	"os"
	"path/filepath"
)

const (
	// CliBlock command line flag to block thread of execution
	CliBlock = "block"

	// CliPool node pool flag
	CliPool = "pool"

	// CliDisableGran disable grand validation flag
	CliDisableGran = "grant"

	//CliIgnoreGrantFailure flag sets ignore grant failure
	CliIgnoreGrantFailure = "ignoreGrantFailure"

	//CliDisableAutoRollback flag disables auto rollback during instantiation
	CliDisableAutoRollback = "disableAutoRollback"

	//CliAutoName flag generate new name upon conflict.
	CliAutoName = "auto_name"

	// CliForce force delete flag
	CliForce = "force"

	// CliNamespace change default name space
	CliNamespace = "namespace"

	// CliProgress show task progress
	CliProgress = "progress"

	// CliDryRun dry run flag
	CliDryRun = "dry"

	// CliShow output spec to stdio
	CliShow = "show"
)

// Chunks splits string to chunks,
// it uses sep to split near chunkSize limit.
// Each chunk is variable size. Method used to partition
// flags usage.
func Chunks(s string, chunkSize int, sep byte) []string {

	if len(s) == 0 {
		return nil
	}

	if chunkSize >= len(s) {
		return []string{s}
	}

	var chunks = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen >= chunkSize && s[i-1] == sep {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

// CmdInitConfig - initialize configuration file, for initial
// setup TCA and other defaults
func (ctl *TcaCtl) CmdInitConfig() *cobra.Command {

	var _cmd = &cobra.Command{
		Use:   "init",
		Short: "Command initializes default tcactl config file.",
		Long: templates.LongDesc(
			`Command Initializes default config file.`),

		Run: func(cmd *cobra.Command, args []string) {

			home, err := homedir.Dir()
			configPrefix := ".tcactl"
			configName := "config"
			configType := "yaml"
			configHome := filepath.Join(home, "/", configPrefix)
			configPath := filepath.Join(configHome, configName+"."+configType)

			_, err = os.ReadDir(home)
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
			fmt.Println("Now run tcactl set and set username, " +
				"password and TCA Cluster endpoint and other configuration settings.")
		},
	}

	return _cmd
}

// CmdSaveConfig - save config file
func (ctl *TcaCtl) CmdSaveConfig() *cobra.Command {

	var _cmd = &cobra.Command{
		Use:   "save",
		Short: "Saves config variables to .tcactl config file.",
		Long:  `Saves config variables to .tcactl config file.`,
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

	var cmdCreate = &cobra.Command{
		Use:   "Create",
		Short: "Command creates a new CNF instance.",
		Long: templates.LongDesc(
			`Command creates a new CNF instance, caller need to provide CNF Identifier.
`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmdCreate.AddCommand(ctl.CmdCreateCnf())
	return cmdCreate
}

// BuildCmd build all commands and attaches to root cmd
// in case you need add sub-command you can, add to plugin dir.
func (ctl *TcaCtl) BuildCmd() {

	var describe = &cobra.Command{
		Use:   "describe [cloud or cluster or nodes or pool or template]",
		Short: "Command describes in details a TCA object.",
		Long: templates.LongDesc(`

Command describes TCA entity. CNFI is CNFI in the inventory, CNFC Catalog entities.`),

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
			if ctl.IsTrace {
				ctl.GetApi().SetTrace(ctl.IsTrace)
			}
		},
	}

	// root cmd for all get
	var cmdGet = &cobra.Command{
		Use:   "get [cnfi, cnfc, clusters, pools]",
		Short: "Command retrieves TCA entity (cnf, catalog, cluster) etc",
		Long: templates.LongDesc(`

Command retrieves tca entities. Tenant, Cluster, Cnfs, Cluster Catalog etc.

`),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// get authentication token
			err := ctl.Authorize()
			if err != nil {
				CheckErrLogError(err)
			}
			// set trace
			if ctl.IsTrace {
				ctl.GetApi().SetTrace(ctl.IsTrace)
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// root cmd for all update commands
	var cmdUpdate = &cobra.Command{
		Use:   "update [cnfi or cnfc]",
		Short: "Command updates or apply changes tca entity cnf, cnf catalog , cluster or node pool.",
		Long: templates.LongDesc(`

Command updates or apply changes to tca entity (cnf, cnf catalog , cluster or node pool.)

`),
		Aliases: []string{"apply"},

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				CheckErrLogError(err)
			}
			if ctl.IsTrace {
				ctl.GetApi().SetTrace(ctl.IsTrace)
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// create root command
	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Command creates a new object in TCA.",
		Long: templates.LongDesc(`

Command creates a new object in TCA. For example new CNF instance, cluster , cluster template etc.`),

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				CheckErrLogError(err)
			}
			if ctl.IsTrace {
				ctl.GetApi().SetTrace(ctl.IsTrace)
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// set root command
	var cmdSet = &cobra.Command{
		Use:   "set",
		Short: "Command sets config variables (Username, Password etc) for tcactl.",
		Long: templates.LongDesc(`
Command sets config variables (Username, Password etc) for tcactl.`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// delete root command
	var cmdDelete = &cobra.Command{
		Use:   "delete",
		Short: "Command deletes object (template,cluster,cnf etc).",
		Long: templates.LongDesc(
			`Command deletes object in tca, (template,cluster,cnf etc).`),
		Aliases: []string{"del"},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := ctl.Authorize()
			if err != nil {
				return
			}
			if ctl.IsTrace {
				ctl.GetApi().SetTrace(ctl.IsTrace)
			}
		},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// Set root command. ( set tca api endpoint , cluster etc)
	cmdSet.AddCommand(
		ctl.CmdSetTca(),
		ctl.CmdSetCluster(),
		ctl.CmdSetNodePool(),
		ctl.CmdSetUsername(),
		ctl.CmdSetPassword())

	// List of all describe sub commands
	describe.AddCommand(
		ctl.CmdDescribeVim(),
		ctl.CmdGetCluster(),
		ctl.CmdDescClusterNodePool(),
		ctl.CmdDescClusterNodePools(),
		ctl.CmdDescClusterNodes(),
		ctl.CmdDescribeTemplate(),
		ctl.CmdDescribeTask())

	// List of all update sub-commands
	cmdUpdate.AddCommand(
		ctl.CmdUpdateExtension(),
		ctl.CmdUpdatePoolNodes(),
		ctl.CmdUpdateTenant(),
		ctl.CmdUpdateClusterTemplates(),
		ctl.CmdUpdateInstance())

	// TCA root command menu
	ctl.RootCmd.AddCommand(
		describe,
		cmdGet,
		cmdUpdate,
		cmdCreate,
		cmdDelete,
		cmdSet,
		ctl.CmdSaveConfig(),
		ctl.CmdInitConfig())

	// Get command
	cmdGet.AddCommand(
		ctl.CmdGetPackages(),
		ctl.CmdGetInstances(),
		ctl.CmdGetRepos(),
		ctl.CmdGetClouds(),
		ctl.CmdVims(),
		ctl.CmdGetClusters(),
		ctl.CmdGetVdu(),
		ctl.CmdGetExtensions(),
		ctl.CmdGetClusterTemplates(),
		ctl.CmdGetVim(),
		ctl.CmdGetTcaManager(),
		ctl.CmdGetVc())

	// Create root command
	cmdCreate.AddCommand(
		ctl.CmdCreateTenant(),
		ctl.CmdCreateCluster(),
		ctl.CmdCreateCnf(),
		ctl.CmdCreateClusterTemplates(),
		ctl.CmdCreatePackage(),
		ctl.CmdCreatePoolNodes(),
		ctl.CmdCreateExtension())

	// Delete
	cmdDelete.AddCommand(
		ctl.CmdDeleteTenant(),
		ctl.CmdDeleteClusterTemplates(),
		ctl.CmdDeleteCluster(),
		ctl.CmdDeleteCatalog(),
		ctl.CmdDeleteTenantCluster(),
		ctl.CmdDeleteInstances(),
		ctl.CmdDeletePoolNodes(),
		ctl.CmdDeleteExtension())

	var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate tcactl completion script",
		Long: `To load completions:

Bash:

$ source <(tcactl completion bash)

# To load completions for each session, execute once:
Linux:
  $ tcactl completion bash > /etc/bash_completion.d/tcactl
MacOS:
  $ tcactl completion bash > /usr/local/etc/bash_completion.d/tcactl

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ tcactl completion zsh > "${fpath[1]}/_tcactl"

# You will need to start a new shell for this setup to take effect.

Fish:

$ tcactl completion fish | source

# To load completions for each session, execute once:
$ tcactl completion fish > ~/.config/fish/completions/tcactl.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}

	ctl.RootCmd.AddCommand(completionCmd)
}
