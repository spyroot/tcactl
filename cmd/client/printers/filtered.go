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
	"fmt"
	"github.com/spyroot/tcactl/cmd/client/main/app/ui"
	"github.com/spyroot/tcactl/cmd/client/response"
)

// VnfPackageFilteredOutput output filter for VnfPackages
func VnfPackageFilteredOutput(r *response.VnfPackages, style ui.PrinterStyle) {
	fields := style.GetFields()
	for _, vnfPackage := range r.Packages {
		for _, f := range fields {
			f = vnfPackage.GetField(f)
			fmt.Println(f)
		}
	}
}

// TenantsFilteredOutput output filter for tenants
func TenantsFilteredOutput(r *response.Tenants, style ui.PrinterStyle) {
	fields := style.GetFields()
	for _, t := range r.TenantsList {
		for _, f := range fields {
			f = t.GetField(f)
			fmt.Println(f)
		}
	}
}
