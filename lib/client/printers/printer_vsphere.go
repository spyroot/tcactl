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
	t.AppendHeader(table.Row{"#", "Name", "Owner", "Parent type", "VC ID"})

	//for i, d := range specs.Data.Items {
	//	if d.EntityType == models.EntityTypeResourcePool {
	//		t.AppendRows([]table.Row{{i, d.Name, d.Owner.Value, d.Parent.Type, d.EntityId}})
	//
	//	}
	//	t.AppendSeparator()
	//}

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
