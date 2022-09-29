// Package printer
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Mustafa mbayramo@vmware.com
package printer

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"os"
)

// ClusterSpecTablePrinter - tabular format printer for repos
func ClusterSpecTablePrinter(c *response.ClusterSpec, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Type", "VC Name", "Endpoint", "Status"})
	t.AppendRows([]table.Row{
		{c.Id, c.ClusterName, c.ClusterType, c.VsphereClusterName, c.ClusterUrl, c.Status},
	})
	t.AppendSeparator()

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// ClusterSpecJsonPrinter - json printer existing cluster details
func ClusterSpecJsonPrinter(spec *response.ClusterSpec, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}

// ClusterSpecYamlPrinter - json printer for existing cluster details
func ClusterSpecYamlPrinter(spec *response.ClusterSpec, style ui.PrinterStyle) {
	DefaultYamlPrinter(spec, style)

}

// ClusterRequestJsonPrinter - json printer for new cluster creation request
func ClusterRequestJsonPrinter(spec *specs.SpecCluster, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}

// ClusterRequestYamlPrinter - json printer for new cluster creation request
func ClusterRequestYamlPrinter(specs *specs.SpecCluster, style ui.PrinterStyle) {
	DefaultYamlPrinter(specs, style)
}

func ClusterTaskTablePrinter(specs *models.ClusterTask, style ui.PrinterStyle) {
	if specs == nil {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Type", "Sub-type", "Task Step",
		"Step Status", "Step Status", "Task Status", "Progress"})
	for i, task := range specs.Items {
		for _, step := range task.Steps {
			t.AppendRows([]table.Row{
				{i, task.EntityDetails.Id, task.EntityDetails.Type, task.Type,
					step.Title, step.Status, step.Status, task.Status, task.Progress},
			})
		}
		t.AppendSeparator()
	}
	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// ClusterTaskJsonPrinter - json printer for new cluster creation request
func ClusterTaskJsonPrinter(spec *models.ClusterTask, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}

// ClusterTaskYamlPrinter - json printer for new cluster creation request
func ClusterTaskYamlPrinter(specs *models.ClusterTask, style ui.PrinterStyle) {
	DefaultYamlPrinter(specs, style)
}

// ConsumptionSpecYamlPrinter - json printer for tca lic consumption
func ConsumptionSpecYamlPrinter(spec *models.ConsumptionResp, style ui.PrinterStyle) {
	DefaultYamlPrinter(spec, style)

}

// ConsumptionJsonPrinter - json printer for new cluster creation request
func ConsumptionJsonPrinter(spec *models.ConsumptionResp, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}

// ConsumptionTablePrinter - tabular format printer for lic consumption
func ConsumptionTablePrinter(specs *models.ConsumptionResp, style ui.PrinterStyle) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"License Qt", "Consumed Qt", "License Ut", "Display Ut", "RawUsage Ut"})

	t.AppendRows([]table.Row{{
		specs.LicenseQuantity,
		specs.ConsumedQuantity,
		specs.LicenseUnit,
		specs.LicenseUnit,
		specs.LicenseDisplayUnit,
		specs.RawUsageUnit},
	})
	t.AppendSeparator()
	fmt.Println("")

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()

	detailTab := table.NewWriter()
	detailTab.SetOutputMirror(os.Stdout)
	detailTab.AppendHeader(table.Row{"vim", "vim name", "vim url", "vim type", "tenant", "consumed qt"})

	for _, c := range specs.Details {
		detailTab.AppendRows([]table.Row{{
			c.VimID,
			c.VimName,
			c.VimURL,
			c.VimType,
			c.TenantName,
			c.ConsumedQuantity},
		})
		t.AppendSeparator()
	}

	detailTabStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		detailTab.SetStyle(detailTabStyle)
	}
	detailTab.Render()
}
