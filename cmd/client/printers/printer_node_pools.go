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
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spyroot/tcactl/cmd/client/main/app/ui"
	"github.com/spyroot/tcactl/cmd/client/response"
	"github.com/spyroot/tcactl/pkg/io"
	"github.com/tidwall/pretty"
	"os"
)

//NodesTablePrinter - tabular format printer for node list in node pool
func NodesTablePrinter(nodePool *response.NodePool, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name"})
	for i, pool := range nodePool.Pools {
		for _, node := range pool.Nodes {
			t.AppendRows([]table.Row{
				{i, node.Ip, node.VmName},
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

// NodesJsonPrinter - json printer
func NodesJsonPrinter(nodePool *response.NodePool, style ui.PrinterStyle) {

	for _, pool := range nodePool.Pools {
		if style.IsColor() {
			b, err := json.MarshalIndent(pool.Nodes, "", "  ")
			io.CheckErr(err)
			fmt.Println(string(pretty.Color(b, nil)))
		} else {
			err := io.PrettyPrint(pool.Nodes)
			io.CheckErr(err)
		}
	}
}

// PoolSpecTablePrinter - tabular format printer for node list in node pool
func PoolSpecTablePrinter(spec *response.NodesSpecs, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Lables", "Name", "CPU", "Memory", "Num Nodes", "Status"})
	t.AppendRows([]table.Row{
		{spec.Id, spec.Labels, spec.Name, spec.Cpu, spec.Memory, len(spec.Nodes), spec.Status},
	})

	t.AppendSeparator()
	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// PoolSpecJsonPrinter - json printer
func PoolSpecJsonPrinter(spec *response.NodesSpecs, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}
