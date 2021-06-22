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
	"github.com/spf13/cobra"
	"github.com/spyroot/hestia/cmd/client/main/app/ui"
	"github.com/spyroot/hestia/cmd/client/response"
	"strings"
)

// CmdDescClusterNodePools - describe all node pool
// command will output all node pool currently
func (ctl *TcaCtl) CmdDescClusterNodePools() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "pools [name or id]",
		Short: "Command describes all node pool",
		Long: `Command describes all node pool, and it 
outputs all node pool currently`,
		Example: "tcactl describe pools",
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

			allSpecs, err := ctl.tca.GetAllNodePools()
			CheckErrLogError(err)

			if _printer, ok := ctl.NodePoolPrinter[_defaultPrinter]; ok {
				_printer(&response.NodePool{
					Pools: allSpecs,
				}, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}
