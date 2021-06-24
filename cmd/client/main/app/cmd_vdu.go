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
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/cmd/client/main/app/ui"
	"strings"
)

//CmdGetVdu retrieves information
func (ctl *TcaCtl) CmdGetVdu() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var cmdCreate = &cobra.Command{
		Use:   "vdu [package name or id]",
		Short: "Command retrieves CNF/VNF VDU information.",
		Long: `

Command retrieves CNF/VNF VDU information. 
The default output format tabular for detail output -o json`,

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Command Requires a catalog cnf/vnf name or ID. " +
					" Use get cnfc to retrieve a catalog list.")
			}
			return nil
		},
		Example: "tcactl get vdu 917a67eb-dcf2-481f-ae36-732aec1ba093",
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			vnfd, err := ctl.tca.GetVdu(args[0])
			CheckErrLogError(err)

			if printer, ok := ctl.VduPrinter[ctl.Printer]; ok {
				printer(vnfd, ctl.DefaultStyle)
			}
		},
	}

	return cmdCreate
}
