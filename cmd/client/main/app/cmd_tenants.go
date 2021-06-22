package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spyroot/hestia/cmd/client/main/app/templates"
)

// CmdDeleteTenantCluster - Deletes cluster template.
func (ctl *TcaCtl) CmdDeleteTenantCluster() *cobra.Command {

	// delete template
	var _cmd = &cobra.Command{
		Use:     "tenant [id or name of tenant cluster]",
		Aliases: []string{"templates"},
		Short:   "Command deletes a tenant cluster.",
		Long: templates.LongDesc(`
									Command deletes a tenant cluster.`),
		Example: " - tcactl delete template my_template",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ok, err := ctl.tca.DeleteTenantCluster(args[0])
			CheckErrLogError(err)
			if ok {
				fmt.Printf("Template %v deleted.")
			}
		},
	}

	return _cmd
}
