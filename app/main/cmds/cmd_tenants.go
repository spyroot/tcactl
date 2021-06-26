// Package app
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
package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
)

// CmdDeleteTenantCluster - Deletes cluster template.
func (ctl *TcaCtl) CmdDeleteTenantCluster() *cobra.Command {

	var _cmd = &cobra.Command{
		Use:     "template [id or name of tenant cluster]",
		Aliases: []string{"templates"},
		Short:   "Command deletes a tenant cluster.",
		Long: templates.LongDesc(`

Command deletes a tenant cluster. Note in order to delete cluster
all instance must be removed`),

		Example: " - tcactl delete cluster cluster",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			ok, err := ctl.tca.DeleteTenantCluster(args[0])
			CheckErrLogError(err)
			if ok {
				fmt.Printf("Template %v deleted.", args[0])
			}
		},
	}

	return _cmd
}
