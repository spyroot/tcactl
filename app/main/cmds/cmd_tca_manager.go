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

// CmdGetTcaManager - a root sub-command for all tca manager sub-commands.
func (ctl *TcaCtl) CmdGetTcaManager() *cobra.Command {

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "tca",
		Aliases: []string{"manager"},
		Short:   "Command retrieves a TCA manager state (license, consumption).",
		Long: templates.LongDesc(`

Command retrieves a TCA manager state (license, consumption).

`),
		Example: "\t - tcactl get vim compute my_cloud_provider\n " +
			"\t - tcactl get tca consumption\n" +
			"\t - tcactl get tca folders my_cloud_provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s requires a subcommand", cmd.Name())
		},
	}

	_cmd.AddCommand(ctl.CmdConsumption())
	return _cmd
}

// CmdConsumption - describe VIM
func (ctl *TcaCtl) CmdConsumption() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		_templateType   = ""
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "consumption",
		Aliases: []string{"con"},
		Short:   "Command license consumption.",
		Long: templates.LongDesc(`

Command retrieves a license consumption.

`),
		Example: " - tcactl get tca consumption",
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			// swap filter if output filter is required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			consumption, err := ctl.tca.GetConsumption(ctx)
			CheckErrLogError(err)

			if printer, ok := ctl.TcaConsumptionPrinter[_defaultPrinter]; ok {
				printer(consumption, _defaultStyler)
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
