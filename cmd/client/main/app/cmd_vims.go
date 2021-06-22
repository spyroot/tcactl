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
	"github.com/spyroot/hestia/cmd/api"
	"github.com/spyroot/hestia/cmd/client/main/app/templates"
	"github.com/spyroot/hestia/cmd/client/main/app/ui"
	"github.com/spyroot/hestia/pkg/io"
	"strings"
)

// CmdGetVim - get VIM root command.
// each sub-command gets particular facts about attached VIM, Cloud provider.
func (ctl *TcaCtl) CmdGetVim() *cobra.Command {

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "vim",
		Aliases: []string{"vims"},
		Short:   "Command retrieves a vim-cloud provider information.",
		Long: templates.LongDesc(`
									Command retrieves a vim-cloud provider information.`),
		Example: " - tcactl describe vim compute my_cloud_provider",
		//Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s requires a subcommand", cmd.Name())
		},
	}

	_cmd.AddCommand(ctl.CmdGetVmwareInfra())
	_cmd.AddCommand(ctl.CmdGetVmwareDatastore())
	_cmd.AddCommand(ctl.CmdGetVmwareNetworks())
	_cmd.AddCommand(ctl.CmdGetVmTemplates())
	_cmd.AddCommand(ctl.CmdGetVmFolders())
	_cmd.AddCommand(ctl.CmdGetVimResourcePool())

	return _cmd
}

// CmdDescribeVim - command describe VIM , Cloud Provider
// attached to system.
func (ctl *TcaCtl) CmdDescribeVim() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		_templateType   = ""
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "vim [name or id]",
		Aliases: []string{"vims"},
		Short:   "Command retrieves a vim-cloud provider information.",
		Long: templates.LongDesc(`
									Command retrieves a list of vim templates.`),
		Example: " - tcactl describe vim vmware_FB40D3DE2967483FBF9033B451DC7571",
		Args:    cobra.MinimumNArgs(1),
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

			tenant, err := ctl.tca.GetVim(args[0])
			CheckErrLogError(err)
			if printer, ok := ctl.TenantsResponsePrinter[_defaultPrinter]; ok {
				printer(tenant, _defaultStyler)
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

// CmdGetVmwareInfra - describe VIM
func (ctl *TcaCtl) CmdGetVmwareInfra() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		_templateType   = ""
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "compute",
		Aliases: []string{"compute"},
		Short:   "Command retrieves a vim information.",
		Long: templates.LongDesc(`
									Command retrieves a vim information.`),
		Example: " - tcactl get vim my_cloud",
		Args:    cobra.MinimumNArgs(1),
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

			clusterInventory, err := ctl.tca.GetVimComputeClusters(args[0])
			CheckErrLogError(err)

			if printer, ok := ctl.VMwareClusterPrinter[_defaultPrinter]; ok {
				printer(clusterInventory, _defaultStyler)
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

// CmdGetVmwareDatastore - command to get retrieve cloud provider
// datastore,  in case VMware ( supported now ) it VSAN/Local and NFS
// data stores.
func (ctl *TcaCtl) CmdGetVmwareDatastore() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "datastore",
		Aliases: []string{"datastore"},
		Short:   "Command retrieves a vim information.",
		Long: templates.LongDesc(`
									Command retrieves a list of vim templates.`),
		Example: " - tcactl get vim datastore vmware_FB40D3DE2967483FBF9033B451DC7571",
		Args:    cobra.MinimumNArgs(1),
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

			clusterInventory, err := ctl.tca.GetVimComputeClusters(args[0])
			CheckErrLogError(err)

			if printer, ok := ctl.VMwareDatastorePrinter[_defaultPrinter]; ok {
				printer(clusterInventory, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}

// CmdGetVmwareNetworks - describe VIM networks.
func (ctl *TcaCtl) CmdGetVmwareNetworks() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "networks",
		Aliases: []string{"network"},
		Short:   "Command retrieves a vim networks.",
		Long: templates.LongDesc(`
									Command retrieves a vim networks.`),
		Example: " - tcactl get vim networks my_cloud_provider",
		Args:    cobra.MinimumNArgs(1),
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

			clusterInventory, err := ctl.tca.GetVimNetworks(args[0])
			CheckErrLogError(err)

			if printer, ok := ctl.VmwareNetworkPrinter[_defaultPrinter]; ok {
				printer(clusterInventory, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}

// CmdGetVmTemplates - describe VIM
func (ctl *TcaCtl) CmdGetVmTemplates() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_templateName   string
		_outputFilter   string
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "templates",
		Aliases: []string{"template"},
		Short:   "Command retrieves a template VM and path.",
		Long: templates.LongDesc(`

Command retrieves a VM template used to retrieve cluster creation. 
This command only works for VMware VC and describes spec for VM template.
By default it search for template v1.20.4+vmware.1.`),

		Example: " - tcactl get vim templates my_cloud_provider\n" +
			" - tcactl get vim templates my_cloud_provider -o yaml\n",
		Args: cobra.MinimumNArgs(1),
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

			vmTemplate, err := ctl.tca.GetVimVMTemplates(args[0], api.VmwareTemplateK8s, _templateName)
			CheckErrLogError(err)

			if printer, ok := ctl.VmwareVmTemplatePrinter[_defaultPrinter]; ok {
				printer(vmTemplate, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	_cmd.Flags().StringVar(&_templateName, "template", "v1.20.4+vmware.1",
		"Default template command filters.")

	return _cmd
}

// CmdGetVmFolders - describe VIM
func (ctl *TcaCtl) CmdGetVmFolders() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "folders",
		Aliases: []string{"folder"},
		Short:   "Command retrieves a VIM folder.",
		Long: templates.LongDesc(`

Command retrieves a folder structure in VMware VC.This command only works for VMware VC
and describes all VM and VM Template folders.
During cluster creation, spec must contains valid folder path 
it is a mandatory requirement for correct placement.`),

		Example: " - tcactl get vim folder my_cloud_provider",
		Args:    cobra.MinimumNArgs(1),
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

			folders, err := ctl.tca.GetVimFolders(args[0])
			CheckErrLogError(err)

			io.PrettyPrint(folders)
			//if printer, ok := ctl.VmwareVmTemplatePrinter[_defaultPrinter]; ok {
			//	printer(vmTemplate, _defaultStyler)
			//}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}

// CmdGetVimResourcePool - return vmware resource pools
// defined in VIM.
func (ctl *TcaCtl) CmdGetVimResourcePool() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "resources",
		Aliases: []string{"resource"},
		Short:   "Command retrieves a VIM's/Cloud Provider a resource pool.",
		Long: templates.LongDesc(`

Command retrieves a VIMs resource pool. This command only works for VMware VC
and describes resources pool defined in VMware VC.
During cluster creation, spec must contains valid resource pool
for correct placement.`),

		Example: " - tcactl get vim resource my_cloud_provider_name\n" +
			" - tcactl get vim resource my_cloud_provider -o yaml",
		Args: cobra.MinimumNArgs(1),
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

			rps, err := ctl.tca.GetVimResourcePool(args[0])
			CheckErrLogError(err)

			if printer, ok := ctl.VmwareResourcePrinter[_defaultPrinter]; ok {
				printer(rps, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter,
		"ofilter", "",
		"Output filter.")

	return _cmd
}
