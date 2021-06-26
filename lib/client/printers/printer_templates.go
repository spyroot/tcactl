// Package printer
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
package printer

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/client/response"
	"os"
)

// TemplateSpecJsonPrinter - json printer for cluster templates
func TemplateSpecJsonPrinter(spec *response.ClusterTemplate, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}

// TemplatesJsonPrinter - json printer for cluster templates
func TemplatesJsonPrinter(specs []response.ClusterTemplate, style ui.PrinterStyle) {
	DefaultJsonPrinter(specs, style)
}

// TemplateSpecYamlPrinter - yaml printer for cluster templates
func TemplateSpecYamlPrinter(spec *response.ClusterTemplate, style ui.PrinterStyle) {
	DefaultYamlPrinter(spec, style)
}

// TemplatesYamlPrinter - yaml printer for cluster templates
func TemplatesYamlPrinter(specs []response.ClusterTemplate, style ui.PrinterStyle) {
	DefaultYamlPrinter(specs, style)
}

// TemplateSpecTablePrinter - tabular format printer for cluster templates
func TemplateSpecTablePrinter(spec *response.ClusterTemplate, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Type", "CNI", "K8S Ver."})
	t.AppendRows([]table.Row{
		{spec.Id, spec.Name, spec.ClusterType, spec.ClusterConfig.Cni, spec.ClusterConfig.KubernetesVersion},
	})
	t.AppendSeparator()

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// TemplatesSpecTablePrinter - tabular format printer for
// TCA Cluster templates.
func TemplatesSpecTablePrinter(specs []response.ClusterTemplate, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name", "Type", "CNI", "K8S Ver."})
	for i, spec := range specs {
		t.AppendRows([]table.Row{
			{i, spec.Id, spec.Name, spec.ClusterType, spec.ClusterConfig.Cni, spec.ClusterConfig.KubernetesVersion},
		})
		t.AppendSeparator()
	}
	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}
