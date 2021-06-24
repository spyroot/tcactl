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
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/cmd/client/main/app/ui"
	"strings"
)

// CmdGetPool - Command retrieves K8s tenant cluster.
func (ctl *TcaCtl) CmdGetPool() *cobra.Command {

	var (
		_defaultNfType  = "CNF"
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "tenant [tca catalog name or id]",
		Short: "Command retrieves K8s tenant cluster, -t option provided option retrieves VIM clusters.",
		Long:  `Command retrieves K8s tenant cluster, -t option provided option retrieves VIM clusters`,
		//Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			// global output type, and terminal wide or not
			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			q, err := ctl.tca.GetTenantsQuery(args[0], _defaultNfType)
			CheckErrLogError(err)
			if printer, ok := ctl.TenantQueryPrinter[_defaultPrinter]; ok {
				printer(q, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter,
		"ofilter", "",
		"Output filter.")

	return _cmd
}

//
func (ctl *TcaCtl) CmdGetExtensions() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var cmdCreate = &cobra.Command{
		Use:   "extensions",
		Short: "Command retrieves API extensions and respected object.",
		Long:  `Command retrieves CNF/VNF VDU information, The default output format tabular for detail output -o json`,
		Run: func(cmd *cobra.Command, args []string) {

			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			ext, err := ctl.TcaClient.ExtensionQuery()
			if err != nil || ext == nil {
				glog.Errorf("Failed retrieve extension information. %v", err)
				return
			}
		},
	}

	return cmdCreate
}
