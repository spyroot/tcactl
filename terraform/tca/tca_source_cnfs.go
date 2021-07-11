// Package tca
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
package tca

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spyroot/tcactl/lib/client/response"
	"log"
	"strconv"
	"time"
)

func dataSourceCnfs() *schema.Resource {
	return &schema.Resource{
		// read
		ReadContext: dataSourceCnfRead,

		Schema: map[string]*schema.Schema{
			"cnfs": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vnf_instance_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vnf_instance_description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vnfd_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vnf_pkg_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vnf_catalog_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCnfRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	tca, err := GetApi(m)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[INFO] Getting list of cnfs")

	cnfs, err := tca.GetAllInstances()
	if err != nil {
		return diag.FromErr(err)
	}

	orderItems := flattenItemsData(cnfs)

	log.Println("[INFO]", len(cnfs.CnfLcms))

	if err := d.Set("cnfs", orderItems); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenItemsData(orderItems *response.CnfsExtended) []interface{} {

	if orderItems != nil {
		cnfs := make([]interface{}, len(orderItems.CnfLcms))
		cnf := make(map[string]interface{})
		for i, orderItem := range orderItems.CnfLcms {

			cnf["id"] = orderItem.CID
			cnf["vnf_instance_name"] = orderItem.VnfInstanceName
			cnf["vnf_instance_description"] = orderItem.VnfInstanceDescription
			cnf["vnfd_id"] = orderItem.VnfdID
			cnf["vnf_pkg_id"] = orderItem.VnfdID
			cnf["vnf_catalog_name"] = orderItem.VnfCatalogName
			//cnf["vnf_provider"] = orderItem.VnfProvider

			cnfs[i] = cnf
		}
		return cnfs
	}

	return make([]interface{}, 0)
}
