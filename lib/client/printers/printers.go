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
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/pkg/io"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// DefaultJsonPrinter Default Json printer
func DefaultJsonPrinter(t interface{}, style ui.PrinterStyle) {
	if style.IsColor() {
		b, err := json.MarshalIndent(t, "", "  ")
		io.CheckErr(err)
		fmt.Println(string(pretty.Color(b, nil)))
	} else {
		err := io.PrettyPrint(t)
		io.CheckErr(err)
	}
}

// DefaultYamlPrinter Default Json printer
func DefaultYamlPrinter(t interface{}, style ui.PrinterStyle) {
	b, err := yaml.Marshal(t)
	io.CheckErr(err)
	if style.IsColor() {
		fmt.Println(string(pretty.Color(b, nil)))
	} else {
		fmt.Println(string(b))
	}
}

// CnfPackageTablePrinter table printer
func CnfPackageTablePrinter(cnfs *response.VnfPackages, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Catalog Name", "VnfID", "Type", "State", "Operation"})

	for i, p := range cnfs.Entity {
		t.AppendRows([]table.Row{
			{i, p.PID, p.VnfProductName, p.VnfdID, p.OnboardingState, p.UsageState, p.VnfSoftwareVersion},
		})

		t.AppendSeparator()
	}
	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// CnfPackageJsonPrinter json pretty printer
func CnfPackageJsonPrinter(cnfs *response.VnfPackages, style ui.PrinterStyle) {

	if cnfs == nil {
		return
	}

	DefaultJsonPrinter(cnfs, style)
}

// CnfPackageYamlPrinter json pretty printer
func CnfPackageYamlPrinter(cnfs *response.VnfPackages, style ui.PrinterStyle) {

	if cnfs == nil {
		return
	}

	DefaultYamlPrinter(cnfs, style)
}

// CnfInstanceTablePrinter table printer
func CnfInstanceTablePrinter(cnfs *response.Cnfs, style ui.PrinterStyle) {

	if cnfs == nil {
		return
	}

	if len(cnfs.CnfLcms) == 0 {
		fmt.Println("No active instance.")
		return
	}

	// default stdout
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name", "Catalog Name", "Vnf pkg ID", "Type", "State", "Operation"})
	for i, s := range cnfs.CnfLcms {
		t.AppendRows([]table.Row{
			{i, s.Id, s.VnfInstanceName, s.VnfProductName, s.Metadata.VnfPkgId,
				s.Metadata.NfType, s.InstantiationState, s.Metadata.LcmOperationState},
		})
		t.AppendSeparator()
	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// CnfInstanceJsonPrinter json pretty printer
func CnfInstanceJsonPrinter(cnfs *response.Cnfs, style ui.PrinterStyle) {

	if cnfs == nil {
		return
	}

	DefaultJsonPrinter(cnfs, style)
}

// CnfInstanceYamlPrinter json pretty printer
func CnfInstanceYamlPrinter(cnfs *response.Cnfs, style ui.PrinterStyle) {

	if cnfs == nil {
		return
	}

	DefaultYamlPrinter(cnfs, style)
}

// CnfInstanceExtendedTablePrinter table printer
func CnfInstanceExtendedTablePrinter(cnfs *response.CnfsExtended, style ui.PrinterStyle) {

	if cnfs == nil {
		return
	}

	// default stdout
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if style.IsWide() {
		t.AppendHeader(table.Row{"#", "ID", "Name", "Catalog Name", "VnfID", "Type", "State", "Operation", "Vim", "Pool"})
		for i, s := range cnfs.CnfLcms {
			var _vimName = ""
			var _nodePool = ""
			if len(s.VimConnectionInfo) > 0 {
				_vimName = s.VimConnectionInfo[0].Extra.VimName
			}
			if len(s.VimConnectionInfo) > 0 {
				_nodePool = s.VimConnectionInfo[0].Extra.NodePoolName
			}
			t.AppendRows([]table.Row{
				{i, s.CID, s.VnfInstanceName, s.VnfProductName,
					s.VnfdID, s.NfType, s.LcmOperationState, s.InstantiationState, _vimName, _nodePool},
			})
			t.AppendSeparator()
		}
	} else {
		t.AppendHeader(table.Row{"#", "ID", "Name", "Catalog Name", "State", "Operation"})
		for i, s := range cnfs.CnfLcms {
			t.AppendRows([]table.Row{
				{i, s.CID, s.VnfInstanceName, s.VnfProductName,
					s.LcmOperationState, s.InstantiationState},
			})
			t.AppendSeparator()
		}
	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// CnfInstanceExtendedJsonPrinter json pretty printer
func CnfInstanceExtendedJsonPrinter(cnfs *response.CnfsExtended, style ui.PrinterStyle) {

	if cnfs == nil {
		return
	}
	DefaultJsonPrinter(cnfs, style)
}

func CnfInstanceExtendedYamlPrinter(cnfs *response.CnfsExtended, style ui.PrinterStyle) {

	if cnfs == nil {
		return
	}
	DefaultYamlPrinter(cnfs, style)
}

//RepoTablePrinter - tabular format printer for repos
func RepoTablePrinter(repo *response.ReposList, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name", "Type", "State"})

	for i, r := range repo.Items {
		for _, e := range r.Repos {
			t.AppendRows([]table.Row{
				{i, r.ID, e.Name, e.Type, r.State},
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

//RepoJsonPrinter - json printer
func RepoJsonPrinter(r *response.ReposList, style ui.PrinterStyle) {

	if r == nil {
		return
	}

	DefaultJsonPrinter(r, style)
}

func RepoYamlPrinter(r *response.ReposList, style ui.PrinterStyle) {

	if r == nil {
		return
	}

	DefaultYamlPrinter(r, style)
}

//TenantsTablePrinter - tabular format printer for repos
func TenantsTablePrinter(tenants *response.Tenants, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "" +
		"Tenant ID",
		"VIM Name", "HCX Cloud",
		"Vim Type", "City", "Latitude", "Longitude", "Remote Status", "Local Status"})

	for i, r := range tenants.TenantsList {
		var (
			la           float64
			lo           float64
			remoteStatus = "Disconnected"
			localStatus  = "Disconnected"
			city         = ""
		)
		if r.Location != nil {
			la = r.Location.Latitude
			lo = r.Location.Longitude
			city = r.Location.City
		}

		if r.VimConn != nil {
			remoteStatus = r.VimConn.RemoteStatus
			localStatus = r.VimConn.Status
		}

		t.AppendRows([]table.Row{
			{i, r.TenantID, r.VimName, r.HcxCloudURL,
				r.VimType, city, la, lo, remoteStatus, localStatus},
		})
		t.AppendSeparator()
	}
	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

//TenantsJsonPrinter - json printer
func TenantsJsonPrinter(t *response.Tenants, style ui.PrinterStyle) {
	DefaultJsonPrinter(t, style)
}

//TenantsYamlPrinter - json printer
func TenantsYamlPrinter(t *response.Tenants, style ui.PrinterStyle) {
	DefaultYamlPrinter(t, style)
}

//NodePoolTablePrinter - tabular format printer for node pool
func NodePoolTablePrinter(p *response.NodePool, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Pool ID", "Pool Name", "Pool Label", "Mem", "CPU", "Compute", "DS", "Status"})

	if p == nil {
		t.Render()
		return
	}

	for i, r := range p.Pools {

		var compute []string
		var datastore []string

		for _, t := range r.PlacementParams {
			if t.Type == "IsValidClusterCompute" {
				compute = append(compute, t.Name)
			}
			if t.Type == "Datastore" {
				datastore = append(datastore, t.Name)
			}
			if t.Type == "ResourcePool" {
			}
		}

		t.AppendRows([]table.Row{
			{i, r.Id, r.Name, r.Labels, r.Memory, r.Cpu,
				strings.Join(compute[:], ","),
				strings.Join(datastore[:], ","),
				r.Status},
		})
		t.AppendSeparator()
	}
	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

//NodePoolJsonPrinter - json printer
func NodePoolJsonPrinter(t *response.NodePool, style ui.PrinterStyle) {
	if style.IsColor() {
		b, err := json.MarshalIndent(t, "", "  ")
		io.CheckErr(err)
		fmt.Println(string(pretty.Color(b, nil)))
	} else {
		err := io.PrettyPrint(t)
		io.CheckErr(err)
	}
}

//NodePoolYamlPrinter - json printer
func NodePoolYamlPrinter(t *response.NodePool, style ui.PrinterStyle) {
	if t == nil {
		return
	}
	DefaultYamlPrinter(t, style)
}

// TenantTabularPinter - print tenant data in tabular format
func TenantTabularPinter(tenants *response.Tenants, style ui.PrinterStyle) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"#", "ID", "Name", "Vim Name", "Tenant Name", "VIM ID", "Type"})
	for i, it := range tenants.TenantsList {
		t.AppendRows([]table.Row{
			{i, it.ID, it.Name, it.VimName, it.TenantName, it.VimID, it.VimType},
		})

		t.AppendSeparator()
	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// TenantJsonPrinter ClusterJsonPrinter - json printer
func TenantJsonPrinter(t *response.Tenants, style ui.PrinterStyle) {
	if style.IsColor() {
		b, err := json.MarshalIndent(t, "", "  ")
		io.CheckErr(err)
		fmt.Println(string(pretty.Color(b, nil)))
	} else {
		err := io.PrettyPrint(t)
		io.CheckErr(err)
	}
}

// TenantYamlPrinter  - yaml printer for tenant printer
func TenantYamlPrinter(t *response.Tenants, style ui.PrinterStyle) {
	if t == nil {
		return
	}

	DefaultYamlPrinter(t, style)
}
