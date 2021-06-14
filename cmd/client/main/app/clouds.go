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
	"github.com/spyroot/hestia/cmd/client/respons"
	"github.com/spyroot/hestia/pkg/io"
	"strings"
)

// CmdGetClouds - return list of cloud provider attached to TCA
func (ctl *TcaCtl) CmdGetClouds() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		_isWide         = false
		vimType         = ""
		hcxUuid         = ""
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "clouds",
		Aliases: []string{"cloud"},
		Short:   "Return list of cloud providers.",
		Long:    `Return list of cloud providers.`,
		Example: "tcactl get clouds edge",
		Args:    cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			// set wide or not
			_isWide, err := cmd.Flags().GetBool(CliWide)
			CheckErrLogError(err)
			_defaultStyler.SetWide(_isWide)

			tenants, err := ctl.TcaClient.GetTenants()
			CheckErrLogError(err)

			if len(args) > 0 {
				r, err := tenants.FindCloudProvider(args[0])
				io.CheckErr(err)
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(&respons.Tenants{
						TenantsList: []respons.TenantsDetails{*r},
					}, _defaultStyler)
				}
				return
			}

			if len(vimType) > 0 {
				r, err := tenants.Filter(respons.FilterVimType, func(q string) bool {
					return strings.HasPrefix(q, vimType)
				})
				CheckErrLogError(err)
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(&respons.Tenants{TenantsList: r}, _defaultStyler)
				}
				return
			}

			if len(hcxUuid) > 0 {
				r, err := tenants.Filter(respons.FilterHcxUUID, func(q string) bool {
					return strings.HasPrefix(q, hcxUuid)
				})
				CheckErrLogError(err)
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(&respons.Tenants{
						TenantsList: r,
					}, _defaultStyler)
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
		"vim_type", "KUBERNETES", "filter by VIM Type. KUBERNETES|VC")
	//
	_cmd.Flags().StringVar(&hcxUuid,
		"hcx_uuid", "", "filter by HCX UUID.")
	// wide output
	_cmd.Flags().BoolVarP(&_isWide,
		"wide", "w", true, "Wide output")
	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}
