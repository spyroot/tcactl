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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spyroot/tcactl/cmd/api"
	"github.com/spyroot/tcactl/cmd/client"
	"log"
	"strconv"
	"time"
)

func dataSourceCnfs() *schema.Resource {
	return &schema.Resource{
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
						"vnf_provider": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instantiated_nf_info": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

	tca, err := api.NewTcaApi(&client.RestClient{
		BaseURL:  "https://tca-vip03.cnfdemo.io",
		ApiKey:   "",
		SkipSsl:  true,
		Client:   nil,
		IsDebug:  true,
		Username: "administrator@vsphere.local",
		Password: "VMware1!",
	},
	)

	log.Println("Authorizing")

	_, err = tca.GetAuthorization()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("Getting list of cnfs")

	cnfs, err := tca.GetCnfs()
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Println("LEN CNFS", len(cnfs.CnfLcms))
	//cnfs := make([]map[string]interface{}, 0)
	//err = json.NewDecoder(cnfs).Decode(&cnfs)
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	if err := d.Set("cnfs", cnfs.CnfLcms); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
