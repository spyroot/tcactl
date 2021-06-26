package printer

import (
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/client/response"
)

// ExtensionsJsonPrinter - json printer
func ExtensionsJsonPrinter(t *response.Clusters, style ui.PrinterStyle) {
	DefaultJsonPrinter(t, style)
}

func ExtensionsYamlPrinter(t *response.Clusters, style ui.PrinterStyle) {
	DefaultYamlPrinter(t, style)
}
