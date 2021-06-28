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
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/io"
	"strings"
)

// CmdGetPackages - Command retrieves a list of CNFs
// or VNFs catalog entities
func (ctl *TcaCtl) CmdGetPackages() *cobra.Command {

	var (
		filter             string
		packageId          string
		vnfProductNameFlag string
		vnfdIdFlag         string
		_outputFilter      string

		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
	)

	var _cmd = &cobra.Command{
		Use:     "cnfc",
		Aliases: []string{"catalog"},
		Short:   "Command retrieves CNF or VNF catalogs entity.",
		Long: `

Command retrieves a list of CNFs or VNFs catalog entities or single element if -i id provide.`,

		Example: "\t - tcactl get catalog df5f3ba2-62f1-4c47-9498-6f7e1acc35cc -o json\n" +
			"\t - tcactl get catalog --vnfd_id nfd_1b6bed2e-6c93-4fd7-83a9-4a8d060fe728 --ofilter PID",

		Run: func(cmd *cobra.Command, args []string) {

			if len(args) > 0 {
				packageId = args[0]
			}

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

			p, err := ctl.tca.GetVnfPkgm(filter, packageId)
			CheckErrLogError(err)

			// filter by name
			if len(vnfProductNameFlag) > 0 {
				r, err := p.Filter(response.VnfProductName, func(q string) bool {
					return strings.HasPrefix(q, vnfProductNameFlag)
				})
				io.CheckErr(err)
				if _printer, ok := ctl.CnfPackagePrinters[_defaultPrinter]; ok {
					_printer(&response.VnfPackages{
						Packages: r,
					}, _defaultStyler)
				}
				return
			}

			// filter by id
			if len(vnfdIdFlag) > 0 {
				r, err := p.Filter(response.VnfdId, func(q string) bool {
					return strings.HasPrefix(q, vnfdIdFlag)
				})
				io.CheckErr(err)
				if _printer, ok := ctl.CnfPackagePrinters[_defaultPrinter]; ok {
					_printer(&response.VnfPackages{
						Packages: r,
					}, _defaultStyler)
				}
				return
			}

			if printer, ok := ctl.CnfPackagePrinters[_defaultPrinter]; ok {
				printer(p, _defaultStyler)
			}
		},
	}

	//
	_cmd.Flags().StringVar(&filter, "filter", "",
		"Adds filter for query to limit a scope of the query.")
	//
	_cmd.Flags().StringVar(&vnfProductNameFlag, "vnf_name", "",
		"Filters by product name.")
	//
	_cmd.Flags().StringVar(&vnfdIdFlag, "vnfd_id", "",
		"Filters by vnfd id.")

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}

// CmdCreatePackage create package in TCA catalog
// by default API interface will regenerate descriptor id
// based on substitution map. All other parameter can overwritten
// by passing respected key.  For example if caller need overwrite
// chart name or chart version.
func (ctl *TcaCtl) CmdCreatePackage() *cobra.Command {

	var (
		filter    string
		packageId string

		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle

		substitution                    = map[string]string{}
		_propertyDescriptorId           = ""
		_PropertyProvider               = ""
		_PropertyDescriptorVersion      = ""
		_PropertyFlavourId              = ""
		_PropertyFlavourDescription     = ""
		_PropertyProductName            = ""
		_PropertyVersion                = ""
		_PropertyId                     = ""
		_PropertySoftwareVersion        = ""
		_PropertyChartName              = ""
		_PropertyChartVersion           = ""
		_PropertyHelmVersion            = ""
		_PropertyName                   = ""
		_PropertyDescription            = ""
		_PropertyConfigurableProperties = ""
		_PropertyVnfmInfo               = ""
	)

	var _cmd = &cobra.Command{
		Use:     "catalog [csar file name, catalog name",
		Aliases: []string{"catalog", "cnfc"},
		Short:   "Command creates CNF or VNF package in TCA catalog.",
		Long: `

Command creates a CNF or VNF package.  It take CSAR file and uploads to a system,
if nfd id already exists in catalog, it will generate a new ID.

Command allow to overwrite some of CSAR Tosca values. Check flags.
`,

		Example: "\ttcactl create catalog my_cnf.csar my_cnf\n\t" +
			"tcactl create catalog my_cnf.csar my_cnf --chart_name my_chart_name\n\t" +
			"tcactl create catalog my_cnf.csar my_cnf --chart_name my_chart_name ----chart_version 1.0",
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) > 0 {
				packageId = args[0]
			}

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			pkgUUID, err := uuid.NewUUID()
			if err != nil {
				return
			}
			substitution[models.PropertyDescriptorId] = "nfd_" + pkgUUID.String()
			if len(_propertyDescriptorId) > 0 {
				substitution[models.PropertyDescriptorId] = _propertyDescriptorId
			}
			if len(_PropertyProvider) > 0 {
				substitution[models.PropertyProvider] = _PropertyProvider
			}
			if len(_PropertyDescriptorVersion) > 0 {
				substitution[models.PropertyDescriptorVersion] = _PropertyDescriptorVersion
			}
			if len(_PropertyFlavourId) > 0 {
				substitution[models.PropertyFlavourId] = _PropertyFlavourId
			}
			if len(_PropertyFlavourDescription) > 0 {
				substitution[models.PropertyFlavourDescription] = _PropertyFlavourDescription
			}
			if len(_PropertyProductName) > 0 {
				substitution[models.PropertyProductName] = _PropertyProductName
			}
			if len(_PropertyVersion) > 0 {
				substitution[models.PropertyVersion] = _PropertyVersion
			}
			if len(_PropertyId) > 0 {
				substitution[models.PropertyId] = _PropertyId
			}
			if len(_PropertySoftwareVersion) > 0 {
				substitution[models.PropertySoftwareVersion] = _PropertySoftwareVersion
			}
			if len(_PropertyChartName) > 0 {
				substitution[models.PropertyChartName] = _PropertyChartName
			}
			if len(_PropertyChartVersion) > 0 {
				substitution[models.PropertyChartVersion] = _PropertyChartVersion
			}
			if len(_PropertyHelmVersion) > 0 {
				substitution[models.PropertyHelmVersion] = _PropertyHelmVersion
			}
			if len(_PropertyName) > 0 {
				substitution[models.PropertyName] = _PropertyName
			}
			if len(_PropertyDescription) > 0 {
				substitution[models.PropertyDescription] = _PropertyDescription
			}
			if len(_PropertyConfigurableProperties) > 0 {
				substitution[models.PropertyConfigurableProperties] = _PropertyConfigurableProperties
			}
			if len(_PropertyVnfmInfo) > 0 {
				substitution[models.PropertyVnfmInfo] = _PropertyVnfmInfo
			}

			ok, err := ctl.tca.CreateNewPackage(args[0], args[1], substitution)
			if err != nil {
				glog.Errorf("Failed create new package. Error: %v", err)
				return
			}

			if ok {
				fmt.Println("Package created.")
			}
		},
	}

	//
	_cmd.Flags().StringVar(&filter, "filter", "",
		"Adds filter for query to limit a scope of the query.")
	//
	_cmd.Flags().StringVar(&_propertyDescriptorId, "nfd_id", "",
		"manually overwrite nfd_id, normally tcactl will "+
			"generate random string.")
	//
	_cmd.Flags().StringVar(&_PropertyProvider, "nfd_provider", "",
		"Overwrite provider field.")
	//
	_cmd.Flags().StringVar(&_PropertyDescriptorVersion, "descriptor_version", "",
		"Overwrite descriptor version.")
	//
	_cmd.Flags().StringVar(&_PropertyFlavourId, "flavour_id", "",
		"Overwrite flavor id.")
	//
	_cmd.Flags().StringVar(&_PropertyFlavourDescription, "flavor_description", "",
		"Overwrite flavor descriptor.")
	//
	_cmd.Flags().StringVar(&_PropertyProductName, "product_name", "",
		"Overwrite product name.")
	//
	_cmd.Flags().StringVar(&_PropertyVersion, "version", "",
		"Overwrite version.")
	//
	_cmd.Flags().StringVar(&_PropertyId, "id", "",
		"Overwrite id field.")
	//
	_cmd.Flags().StringVar(&_PropertySoftwareVersion, "software_version", "",
		"Overwrite software version.")
	//
	_cmd.Flags().StringVar(&_PropertyChartName, "chart_name", "",
		"Overwrite chart name.")
	//
	_cmd.Flags().StringVar(&_PropertyChartVersion, "chart_version", "",
		"Overwrite chart version.")
	//
	_cmd.Flags().StringVar(&_PropertyHelmVersion, "helm_version", "",
		"Overwrite helm version.")
	//
	_cmd.Flags().StringVar(&_PropertyName, "name", "",
		"Overwrite name.")
	//
	_cmd.Flags().StringVar(&_PropertyDescription, "description", "",
		"Overwrite description.")
	//
	_cmd.Flags().StringVar(&_PropertyConfigurableProperties, "conf_properties", "",
		"Overwrite configurable properties.")
	//
	_cmd.Flags().StringVar(&_PropertyVnfmInfo, "vnfm_info", "",
		"Overwrite vnfm info.")

	return _cmd
}
