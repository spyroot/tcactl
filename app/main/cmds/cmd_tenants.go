// Package cmds
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
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"strings"
)

// CmdVims - Command retrieves K8s tenant cluster.
func (ctl *TcaCtl) CmdVims() *cobra.Command {

	var (
		//		_defaultNfType  = "CNF"
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "tenant [tenant vim or nothing for list]",
		Short: "Command retrieves particular cloud provider or list of all providers.",
		Long: `

Command retrieves particular cloud provider or list of all providers.`,

		//Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var (
				ctx = context.Background()
				t   *response.Tenants
				err error
			)

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
			ctl.tca.SetTrace(ctl.IsTrace)

			if len(args) > 0 {
				t, err = ctl.tca.GetTenant(ctx, args[0])
				CheckErrLogError(err)
				if t != nil && len(t.TenantsList) == 0 {
					fmt.Printf("Tenant %s not found\n", args[0])
					return
				}
			} else {
				t, err = ctl.tca.GetVims(ctx)
				CheckErrLogError(err)
			}
			if t != nil {
				if printer, ok := ctl.TenantQueryPrinter[_defaultPrinter]; ok {
					printer(t, _defaultStyler)
				}
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter,
		"ofilter", "",
		"Output filter.")

	_cmd.Flags().StringVar(&_outputFilter,
		"--", "",
		"Output filter.")
	return _cmd
}

// CmdDeleteTenant - Command deletes tenant
// TODO do recursive force i.e remove all
func (ctl *TcaCtl) CmdDeleteTenant() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "provider [provider id, or name]",
		Short: "Command deletes cloud provider.",
		Long: templates.LongDesc(`

Command delete cloud provider. Note all entity must be removed.`),

		//Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var (
				ctx = context.Background()
				t   *response.Tenants
				err error
			)

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
			ctl.tca.SetTrace(ctl.IsTrace)

			if len(args) > 0 {
				_, err := ctl.tca.DeleteCloudProvider(ctx, args[0])
				CheckErrLogError(err)

				fmt.Printf("cloud provider %s delete\n", args[0])

			} else {
				// if no args
				fmt.Println("Please provide cloud provider name or id.")
				t, err = ctl.tca.GetVims(ctx)
				CheckErrLogError(err)
			}
			if t != nil {
				if printer, ok := ctl.TenantQueryPrinter[_defaultPrinter]; ok {
					printer(t, _defaultStyler)
				}
			}
		},
	}

	return _cmd
}

// CmdCreateTenant - Command create-registers a new tenant
func (ctl *TcaCtl) CmdCreateTenant() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "tenant [spec file]",
		Short: "Command attaches cloud provider to TCA.",
		Long: templates.LongDesc(`

Command attaches cloud provider to TCA.`),

		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			spec, err := request.ProviderSpecsFromFile(args[0])
			CheckErrLogError(err)

			_, err = ctl.tca.CreateTenantProvider(spec)
			CheckErrLogError(err)

			//fmt.Printf("Cloud provider %s registered\n", args[0])
		},
	}

	return _cmd
}

// CmdUpdateTenant - Command update tenant
// TODO
func (ctl *TcaCtl) CmdUpdateTenant() *cobra.Command {

	var (
		//		_defaultNfType  = "CNF"
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "tenant [spec file]",
		Short: "Command update cloud provider details and trigger re-attach.",
		Long: `

Command update cloud provider details and trigger re-attach.`,

		//Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var (
				ctx = context.Background()
				t   *response.Tenants
				err error
			)

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
			ctl.tca.SetTrace(ctl.IsTrace)

			if len(args) > 0 {
				t, err = ctl.tca.GetTenant(ctx, args[0])
				CheckErrLogError(err)
				if t != nil && len(t.TenantsList) == 0 {
					fmt.Printf("Tenant %s not found\n", args[0])
					return
				}
			} else {
				t, err = ctl.tca.GetVims(ctx)
				CheckErrLogError(err)
			}
			if t != nil {
				if printer, ok := ctl.TenantQueryPrinter[_defaultPrinter]; ok {
					printer(t, _defaultStyler)
				}
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter,
		"ofilter", "",
		"Output filter.")

	_cmd.Flags().StringVar(&_outputFilter,
		"--", "",
		"Output filter.")
	return _cmd
}

// CmdDeleteTenantCluster - Deletes tenant cluster
func (ctl *TcaCtl) CmdDeleteTenantCluster() *cobra.Command {

	var _cmd = &cobra.Command{
		Use:     "tenant [id or name of tenant cluster]",
		Aliases: []string{"templates"},
		Short:   "Command deletes a tenant cluster.",
		Long: templates.LongDesc(`

The command deletes a tenant cluster. Note to delete tenant cluster all 
active CNF instance must be removed first.

`),

		Example: " - tcactl delete tenant my_tenant_name",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()

			task, err := ctl.tca.DeleteTenantCluster(ctx, args[0])
			CheckErrLogError(err)
			fmt.Printf("Tenant cluster %v deleted. Task id %s\n", args[0], task.OperationId)
		},
	}

	return _cmd
}
