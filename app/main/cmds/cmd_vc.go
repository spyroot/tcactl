// Package cmds
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Mustafa mbayramo@vmware.com
package cmds

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"strings"
)

// CmdGetVc - a root sub-command for all vc sub-commands.
func (ctl *TcaCtl) CmdGetVc() *cobra.Command {

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "vc",
		Aliases: []string{"vc"},
		Short:   "Command retrieves a vc inventory object directly.",
		Long: templates.LongDesc(`

Command retrieves a vc inventory object details. Note tcactl config file
must contain vc fqdn, username and password.

`),
		Example: "\t - tcactl get vc datastore vsanDatastore\n " +
			"\t - tcactl get vc nic driver\n" +
			"\t - tcactl get vc get nic my_cloud_provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s requires a subcommand", cmd.Name())
		},
	}

	_cmd.AddCommand(ctl.CmdGetDatastore())
	return _cmd
}

// CmdGetDatastore - get datastores from vcenter
func (ctl *TcaCtl) CmdGetDatastore() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		_templateType   = ""
		defaultVc       = "default"
	)

	// datastore
	var _cmd = &cobra.Command{
		Use:     "datastores",
		Aliases: []string{"ds", "datastore"},
		Short:   "Command retrieve list of datastores.",
		Long: templates.LongDesc(`

Command retrieves list of all datastores or datastores for particular path or datacenter.

`),
		Example: " - tcactl get vc ds",
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			ctx := context.Background()
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			// swap filter if output filter is required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			err := ctl.VcConnect(ctx, defaultVc)
			CheckErrLogError(err)
			vcdss, err := ctl.vcRest.GetDatastores(ctx, "")
			if err != nil {
				CheckErrLogError(err)
				return
			}
			if printer, ok := ctl.VsphereDatastores[_defaultPrinter]; ok {
				printer(vcdss, _defaultStyler)
			}
		},
	}

	//
	_cmd.Flags().StringVar(&_templateType,
		"type", "", "filter by template type.")

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}
