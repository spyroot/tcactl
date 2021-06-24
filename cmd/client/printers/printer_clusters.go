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
	"github.com/spyroot/tcactl/cmd/client/main/app/ui"
	"github.com/spyroot/tcactl/cmd/client/response"
	"os"
)

//ClusterTablePrinter - tabular format printer for repos
func ClusterTablePrinter(clusters *response.Clusters, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name", "Type", "VC Name", "Endpoint", "Status"})
	for i, c := range clusters.Clusters {
		t.AppendRows([]table.Row{
			{i, c.Id, c.ClusterName, c.ClusterType, c.VsphereClusterName, c.ClusterUrl, c.Status},
		})
		t.AppendSeparator()
	}
	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// ClusterJsonPrinter - json printer
func ClusterJsonPrinter(t *response.Clusters, style ui.PrinterStyle) {
	DefaultJsonPrinter(t, style)
}

func ClusterYamlPrinter(t *response.Clusters, style ui.PrinterStyle) {
	DefaultYamlPrinter(t, style)
}
