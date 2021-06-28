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
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/api"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/pkg/io"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

// CmdGetClusterTemplates - return list of cloud provider attached to TCA
// output filter allow to filter by specific key
// filter allow to filter on template type
func (ctl *TcaCtl) CmdGetClusterTemplates() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
		_templateType   = ""
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "templates [name or id]",
		Aliases: []string{"template"},
		Short:   "Command retrieves a list of cluster templates.",
		Long: templates.LongDesc(`
									Command retrieves a list of cluster templates.`),
		Example: " - tcactl get templates --type WORKLOAD\n -tcactl get templates --type WORKLOAD -o json -t",
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

			tmpl, err := ctl.tca.GetClusterTemplates()
			CheckErrLogError(err)
			if len(_templateType) > 0 {
				_templateType = strings.ToUpper(_templateType)
				if isValidTemplateType(_templateType) == false {
					CheckErrLogError("template must be workload or management.")
				}
				tmpl, err = tmpl.Filter(response.FilterTemplateType, func(q string) bool {
					return strings.HasPrefix(q, _templateType)
				})
				io.CheckErr(err)
			}
			if len(tmpl.ClusterTemplates) == 1 {
				if printer, ok := ctl.TemplatePrinter[_defaultPrinter]; ok {
					printer(&tmpl.ClusterTemplates[0], _defaultStyler)
				}
			} else {
				if printer, ok := ctl.TemplatesPrinter[_defaultPrinter]; ok {
					printer(tmpl.ClusterTemplates, _defaultStyler)
				}
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

// CmdCreateClusterTemplates - Create new cluster template
func (ctl *TcaCtl) CmdCreateClusterTemplates() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		isDry           = false
	)
	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "template [file]",
		Aliases: []string{"template"},
		Short:   "Command creates a cluster template.",
		Long: templates.LongDesc(`
Command creates a cluster template from input spec.
`),
		Example: " - tcactl create template template_spec.yaml -o json --dry",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			var spec response.ClusterTemplate
			if io.FileExists(args[0]) {
				buffer, err := ioutil.ReadFile(args[0])
				CheckErrLogError(err)
				err = yaml.Unmarshal(buffer, &spec)
				CheckErrLogError(err)
			} else {
				CheckErrLogError(fmt.Errorf("%v not found", args[0]))
			}

			if isDry {
				if printer, ok := ctl.TemplatePrinter[_defaultPrinter]; ok {
					printer(&spec, _defaultStyler)
				}
				return
			}

			name, err := ctl.tca.CreateClusterTemplate(&spec)
			CheckErrLogError(err)
			fmt.Printf("Template %v created.\n", name)
		},
	}

	_cmd.Flags().BoolVar(&isDry,
		"dry", false, "Parses template spec and validate, dry run outputs spec "+
			"to terminal screen and format based based on -o.")

	return _cmd
}

// CmdDescribeTemplate - describe single template
func (ctl *TcaCtl) CmdDescribeTemplate() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	// cloud - tenants
	var _cmd = &cobra.Command{
		Use:     "templates [name or id]",
		Aliases: []string{"template"},
		Short:   "Command describes a cluster templates.",
		Long: templates.LongDesc(`
									Command describes a cluster templates.`),
		Example: " - tcactl describe templates ",
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

			var templateId = args[0]
			t, err := ctl.tca.GetNamedClusterTemplate(templateId)
			CheckErrLogError(err)

			if printer, ok := ctl.TemplatePrinter[_defaultPrinter]; ok {
				printer(t, _defaultStyler)
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}

// CmdDeleteClusterTemplates - Deletes cluster template.
func (ctl *TcaCtl) CmdDeleteClusterTemplates() *cobra.Command {

	// delete template
	var _cmd = &cobra.Command{
		Use:     "tenant [id or name of tenant cluster]",
		Aliases: []string{"templates"},
		Short:   "Command deletes a cluster template.",
		Long: templates.LongDesc(`
Command deletes a cluster template.  
`),
		Example: " - tcactl delete template my_template",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			err := ctl.tca.DeleteTemplate(args[0])
			CheckErrLogError(err)

			fmt.Printf("Template %v deleted.", args[0])
		},
	}

	return _cmd
}

// CmdUpdateClusterTemplates - Updates cluster template
// from a file spec.
func (ctl *TcaCtl) CmdUpdateClusterTemplates() *cobra.Command {

	var templateId = ""

	// delete template
	var _cmd = &cobra.Command{
		Use:     "template [id or name of template]",
		Aliases: []string{"templates"},
		Short:   "Command updates a cluster template.",
		Long: templates.LongDesc(`

Command update a cluster template. 
It take input spec and apply to existing cluster template.`),

		Example: " - tcactl update template template_spec_min.yaml",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// read template.
			var (
				spec *response.ClusterTemplate
				err  error
			)

			if io.FileExists(args[0]) {
				spec, err = api.ReadTemplateSpecFromFile(args[0])
				CheckErrLogError(err)
			} else {
				CheckErrLogError(fmt.Errorf("%v not found", args[0]))
			}

			// validate id
			if len(templateId) == 0 {
				templateId = spec.Id
				if len(templateId) == 0 {
					CheckErrLogError(fmt.Errorf(
						"you must indicate tempalte id in spec or provide template_id argument"))
				}
			} else {
				glog.Infof("Using template id indicate by client %s", templateId)
				spec.Id = templateId
			}

			err = ctl.tca.UpdateClusterTemplate(spec)
			CheckErrLogError(err)

			fmt.Printf("Template %v Updated.", templateId)
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&templateId, "template_id", "",
		"template id.")

	return _cmd
}

// CmdDeleteTenantCluster - Deletes cluster template.
func (ctl *TcaCtl) CmdDeleteTenantCluster() *cobra.Command {

	var _cmd = &cobra.Command{
		Use:     "template [id or name of tenant cluster]",
		Aliases: []string{"templates"},
		Short:   "Command deletes a tenant cluster.",
		Long: templates.LongDesc(`

Command deletes a tenant cluster. Note in order to delete cluster
all instance must be removed`),

		Example: " - tcactl delete cluster cluster",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			task, err := ctl.tca.DeleteTenantCluster(args[0])
			CheckErrLogError(err)
			fmt.Printf("Template %v deleted. Task id %s\n", args[0], task.OperationId)
		},
	}

	return _cmd
}
