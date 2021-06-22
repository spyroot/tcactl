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
	"github.com/spyroot/hestia/cmd/client/main/app/templates"
	"github.com/spyroot/hestia/cmd/client/main/app/ui"
	"github.com/spyroot/hestia/cmd/client/response"
	"github.com/spyroot/hestia/pkg/io"
	"strings"
)

// CmdGetClouds - return list of cloud provider attached to TCA
func (ctl *TcaCtl) CmdGetClouds() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		vimType         = ""
		hcxUuid         = ""
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "clouds [name or id]",
		Aliases: []string{"cloud"},
		Short:   "Command retrieves a list of cloud providers.",
		Long: templates.LongDesc(`
									Command retrieve a list of cloud providers currently attached to TCA.`),
		Example: " - tcactl get clouds \n - tcactl get clouds edge",
		Args:    cobra.RangeArgs(0, 1),
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

			tenants, vimErr := ctl.tca.GetVimTenants()
			CheckErrLogError(vimErr)

			if len(args) > 0 {
				r, err := ctl.tca.TenantsCloudProvider(args[0])
				io.CheckErr(err)
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(r, _defaultStyler)
				}
				return
			}

			if len(vimType) > 0 {
				r, err := tenants.Filter(response.FilterVimType, func(q string) bool {
					return strings.HasPrefix(q, vimType)
				})
				CheckErrLogError(err)
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(&response.Tenants{TenantsList: r}, _defaultStyler)
				}
				return
			}

			if len(hcxUuid) > 0 {
				r, err := tenants.Filter(response.FilterHcxUUID, func(q string) bool {
					return strings.HasPrefix(q, hcxUuid)
				})
				CheckErrLogError(err)
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(&response.Tenants{TenantsList: r}, _defaultStyler)
				}
				return
			}

			if tenants != nil {
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(tenants, _defaultStyler)
				}
			}
		},
	}

	//
	_cmd.Flags().StringVar(&vimType,
		"vim_type", "", "filter by VIM Type. KUBERNETES|VC")
	//
	_cmd.Flags().StringVar(&hcxUuid,
		"hcx_uuid", "", "filter by HCX UUID.")

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}
