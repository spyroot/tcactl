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
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/api"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/io"
	"strings"
)

// CmdGetClouds - return list of cloud provider
// attached to Telco Cloud Automation
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
		Short:   "Command retrieves a list of cloud providers. (tenants)",
		Long: templates.LongDesc(`
Command retrieve a list of cloud providers currently attached to Telco Cloud Automation. 
Workload and Tenant Kubernetes clusters must be active in the target cloud provider.
`),
		Example: "\t- tcactl get clouds \n" +
			"\t- tcactl get clouds edge",
		Args: cobra.RangeArgs(0, 1),
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

			ctx := context.Background()
			tenants, vimErr := ctl.tca.GetVimTenants(ctx)
			CheckErrLogError(vimErr)

			if len(args) > 0 {
				r, err := ctl.tca.TenantsCloudProvider(ctx, args[0])
				io.CheckErr(err)
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(r, _defaultStyler)
				}
				return
			}

			if len(vimType) > 0 {
				r := tenants.Filter(response.FilterVimType, func(q string) bool {
					return strings.HasPrefix(strings.ToLower(q), strings.ToLower(vimType))
				})
				if printer, ok := ctl.TenantsPrinter[_defaultPrinter]; ok {
					printer(&response.Tenants{TenantsList: r}, _defaultStyler)
				}
				return
			}

			if len(hcxUuid) > 0 {
				r := tenants.Filter(response.FilterHcxUUID, func(q string) bool {
					return strings.HasPrefix(q, hcxUuid)
				})
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
	vimTypeUsage := fmt.Sprintf("filter by VIM Type. %s\n",
		strings.Join(models.VimType(), ","))

	_cmd.Flags().StringVar(&vimType,
		"vim_type", "", vimTypeUsage)
	//
	_cmd.Flags().StringVar(&hcxUuid,
		"hcx_uuid", "", "filter by HCX UUID.")

	fields := strings.Join(api.TenantFields(), ",")
	chunks := Chunks(fields, 50, ',')
	outputUsage := fmt.Sprintf("Output filter. (%s", chunks[0])
	for i, chunk := range chunks {
		if i > 0 {
			outputUsage += chunk
		}
		if i == len(chunk) {
			outputUsage += "\n)"
		} else {
			outputUsage += "\n"
		}
	}

	//outputFilterUsage := fmt.Sprintf("filter by VIM Type. %s\n", strings.Join(models.,","))
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "", outputUsage)

	return _cmd
}
