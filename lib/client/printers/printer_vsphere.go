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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/pkg/vmware/vc"
	"os"
)

// VsphereDatastoresTablePrinters - a tabular format printer for vSphere datastores list cmd
func VsphereDatastoresTablePrinters(specs *vc.VsphereDatastores, style ui.PrinterStyle) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Inventory Path", "DC Path", "Type"})

	for s := range specs.Datastores {
		t.AppendRows([]table.Row{{specs.Datastores[s].Name, specs.Datastores[s].InventoryPath, specs.Datastores[s].DatacenterPath, specs.Datastores[s].Type}})
	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// VsphereDatastoresJsonPrinters - json printer for vSphere datastores list cmd
func VsphereDatastoresJsonPrinters(specs *vc.VsphereDatastores, style ui.PrinterStyle) {
	DefaultJsonPrinter(specs, style)
}

// VsphereDatastoresYamlPrinters - json printer for vSphere datastores list cmd
func VsphereDatastoresYamlPrinters(specs *vc.VsphereDatastores, style ui.PrinterStyle) {
	DefaultYamlPrinter(specs, style)
}
