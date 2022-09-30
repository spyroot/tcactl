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
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/models"
	"os"
)

// TenantsResponseYamlPrinter - json printer for cluster templates
func TenantsResponseYamlPrinter(spec *response.TenantSpecs, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}

// TenantsResponseJsonPrinter - json printer for cluster templates
func TenantsResponseJsonPrinter(specs *response.TenantSpecs, style ui.PrinterStyle) {
	DefaultJsonPrinter(specs, style)
}

// VimTablePrinter - tabular format printer for node list in node pool
func VimTablePrinter(specs *response.TenantSpecs, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Cloud Owner", "Cloud Name", "Vim ID", "Tenant Name", "Type", "VIM API"})
	for i, c := range specs.Tenants {
		t.AppendRows([]table.Row{{i, specs.CloudOwner, specs.VimName, specs.VimId, c.Name, c.VimType, c.VimURL}})
		t.AppendSeparator()
	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// VmwareInventoryYamlPrinter - json printer for cluster templates
func VmwareInventoryYamlPrinter(spec *models.VMwareClusters, style ui.PrinterStyle) {
	DefaultYamlPrinter(spec, style)
}

// VmwareInventoryJsonPrinter - json printer for cluster templates
func VmwareInventoryJsonPrinter(specs *models.VMwareClusters, style ui.PrinterStyle) {
	DefaultJsonPrinter(specs, style)
}

// VmwareInventoryTablePrinter - tabular format printer for compute nodes, clusters
// attached to cloud provider.
func VmwareInventoryTablePrinter(specs *models.VMwareClusters, style ui.PrinterStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"#", "Name", "Type", "Num hosts", "Memory",
		"CPU", "Deployed", "Num Mgmt", "Num Clusters"})

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
		{Number: 4, AutoMerge: true},
		{Number: 5, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 6, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 7, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 8, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 9, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
	})

	for i, c := range specs.Items {
		t.AppendRows([]table.Row{{i, c.Name, c.EntityType, c.NumOfHosts, c.Memory,
			c.Cpu, c.K8ClusterDeployed, c.NumK8SMgmtClusterDeployed, c.NumK8SWorkloadClusterDeployed},
		})
		t.AppendSeparator()
	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// VmwareDatastoreTablePrinter - tabular format printer for node
// datastores attached to compute.
func VmwareDatastoreTablePrinter(specs *models.VMwareClusters, style ui.PrinterStyle) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Cluster", "Name", "Type", "Capacity", "Free", "Accessible"})

	for _, c := range specs.Items {
		for j, s := range c.Datastore {
			t.AppendRows([]table.Row{{j,
				c.Name, s.Name, s.Summary.Type, s.Summary.Capacity, s.Summary.FreeSpace, s.Summary.Accessible},
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

// VmwareNetworkTablePrinter - tabular format printer for cloud networks
// attached to cloud provider.
func VmwareNetworkTablePrinter(specs *models.CloudNetworks, style ui.PrinterStyle) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "VC ID", "port-group", "Type", "DVS Name", "Path"})

	for i, c := range specs.Network {
		t.AppendRows([]table.Row{{i, c.Name, c.Id, c.Type, c.NetworkType, c.DvsName, c.FullNetworkPath}})
		t.AppendSeparator()

	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// VmwareNetworkJsonPrinter - json printer for cloud networks
// attached to cloud provider.
func VmwareNetworkJsonPrinter(spec *models.CloudNetworks, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}

// VmwareNetworkYamlPrinter - json printer for cloud networks
// attached to cloud provider.
func VmwareNetworkYamlPrinter(specs *models.CloudNetworks, style ui.PrinterStyle) {
	DefaultYamlPrinter(specs, style)
}

// VmwareTemplateTablePrinter - tabular format printer printer for VMware template
// attached to compute cluster
func VmwareTemplateTablePrinter(specs *models.VcInventory, style ui.PrinterStyle) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Cluster", "Name", "IS Template", "Path"})

	for i, c := range specs.Items {
		t.AppendRows([]table.Row{{i, c.Id, c.Name, c.IsTemplate, c.FullPath}})
		t.AppendSeparator()

	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// VmwareTemplateJsonPrinter - json printer for VMware template
// attached to compute cluster
func VmwareTemplateJsonPrinter(spec *models.VcInventory, style ui.PrinterStyle) {
	DefaultJsonPrinter(spec, style)
}

// VmwareTemplateYamlPrinter - json printer for VMware template
// attached to compute cluster
func VmwareTemplateYamlPrinter(specs *models.VcInventory, style ui.PrinterStyle) {
	DefaultYamlPrinter(specs, style)
}

// VmwareResourcePoolTablePrinter - a tabular format printer resource pools
// attached to compute cluster
func VmwareResourcePoolTablePrinter(specs *models.ResourcePool, style ui.PrinterStyle) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Owner", "Parent type", "VC ID"})

	for i, d := range specs.Data.Items {
		if d.EntityType == models.EntityTypeResourcePool {
			t.AppendRows([]table.Row{{i, d.Name, d.Owner.Value, d.Parent.Type, d.EntityId}})

		}
		t.AppendSeparator()
	}

	tableStyle, ok := style.GetTableStyle().(table.Style)
	if ok {
		t.SetStyle(tableStyle)
	}
	t.Render()
}

// VmwareResourcePoolJsonPrinter - json printer for VMware resource pools
// attached to compute cluster
func VmwareResourcePoolJsonPrinter(specs *models.ResourcePool, style ui.PrinterStyle) {
	DefaultJsonPrinter(specs, style)
}

// VmwareResourcePoolYamlPrinter - json printer for VMware resource pools
// attached to compute cluster
func VmwareResourcePoolYamlPrinter(specs *models.ResourcePool, style ui.PrinterStyle) {
	DefaultYamlPrinter(specs, style)
}
