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

// VduTablePrinter ClusterTablePrinter - tabular format printer for repos
func VduTablePrinter(vdus *response.VduPackage, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name", "Type", "VC Name", "Endpoint", "Status"})
	for i, v := range vdus.Vdus {
		t.AppendRows([]table.Row{
			{i, v.VduId, v.Type, v.Properties.Name, v.Properties.ChartName, v.Properties.ChartVersion},
		})
		t.AppendSeparator()
	}
	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// VduJsonPrinter - json printer
func VduJsonPrinter(t *response.VduPackage, style ui.PrinterStyle) {
	DefaultJsonPrinter(t, style)
}

// VduYamlPrinter  - json printer
func VduYamlPrinter(t *response.VduPackage, style ui.PrinterStyle) {
	DefaultYamlPrinter(t, style)
}
