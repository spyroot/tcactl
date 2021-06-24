package printer

import (
	"github.com/spyroot/tcactl/cmd/client/main/app/ui"
	"github.com/spyroot/tcactl/cmd/client/response"
)

// ExtensionsJsonPrinter - json printer
func ExtensionsJsonPrinter(t *response.Clusters, style ui.PrinterStyle) {
	DefaultJsonPrinter(t, style)
}

func ExtensionsYamlPrinter(t *response.Clusters, style ui.PrinterStyle) {
	DefaultYamlPrinter(t, style)
}
