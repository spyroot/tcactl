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
	"github.com/golang/glog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spyroot/tcactl/lib/api"
	"github.com/spyroot/tcactl/lib/client/response"
	"log"
	"strconv"
	"time"
)

func dataSourceMaster() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mem": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replica": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"clone_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func dataSourceWorker() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mem": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replica": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"clone_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"cpu_manager_policy": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrderRead,
		Schema: map[string]*schema.Schema{
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataSourceMaster(),
						},
						"worker": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataSourceWorker(),
						},
					},
				},
			},
		},
	}
}

// dataSourceOrderRead
func dataSourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type

	log.Println("[DEBUG] dataSourceOrderRead reading")

	var diags diag.Diagnostics

	tca := m.(*api.TcaApi)
	if tca == nil {
		return diag.FromErr(fmt.Errorf("nil client"))
	}

	_, err := tca.GetAuthorization()
	if err != nil {
		return diag.FromErr(err)
	}

	glog.Infof("dataSourceOrderRead")
	//u, err := uuid.Parse(d.Get("id").(string))
	//if err != nil {
	//	panic(err)
	//}

	//d.Get
	//d.Get("id")
	//
	//orderID, ok := d.Get("id").(string)
	//
	templates, err := tca.GetClusterTemplates()
	if err != nil {
		return diag.FromErr(err)
	}

	t := flattenTemplatesData(templates)
	if err := d.Set("templates", t); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	//d.SetId(orderID)

	return diags
}

// flattenTemplatesData - read cluster template
func flattenTemplatesData(templates *response.ClusterTemplates) []interface{} {

	if templates != nil {
		tis := make([]interface{}, len(templates.ClusterTemplates), len(templates.ClusterTemplates))

		for i, t := range templates.ClusterTemplates {
			ti := make(map[string]interface{})

			ti["id"] = t.Id
			ti["name"] = t.Name
			ti["cluster_type"] = t.ClusterType
			ti["description"] = t.Description

			// master node spec
			mis := make([]interface{}, len(t.MasterNodes), len(t.MasterNodes))
			for j, m := range t.MasterNodes {
				mi := make(map[string]interface{})
				mi["name"] = m.Name
				mi["cpu"] = m.Cpu
				mi["mem"] = m.Memory
				mi["storage"] = m.Storage
				mi["clone_mode"] = m.CloneMode
				mi["replica"] = m.Replica
				mi["labels"] = m.Labels
				var nets []string
				for _, n := range m.Networks {
					nets = append(nets, n.Label)
				}
				mi["network"] = nets
				mis[j] = mi
			}
			ti["master"] = mis

			// worker node
			wis := make([]interface{}, len(t.WorkerNodes), len(t.WorkerNodes))
			for j, m := range t.WorkerNodes {
				wi := make(map[string]interface{})
				wi["name"] = m.Name
				wi["cpu"] = m.Cpu
				wi["mem"] = m.Memory
				wi["storage"] = m.Storage
				wi["clone_mode"] = m.CloneMode
				wi["replica"] = m.Replica
				wi["labels"] = m.Labels
				var nets []string
				for _, n := range m.Networks {
					nets = append(nets, n.Label)
				}
				wi["network"] = nets

				if m.Config.CpuManagerPolicy != nil {

					var config []map[string]interface{}
					l := map[string]interface{}{
						"type":   m.Config.CpuManagerPolicy.Type,
						"policy": m.Config.CpuManagerPolicy.Policy,
					}
					config = append(config, l)
					wi["cpu_manager_policy"] = config

				}
				wis[j] = wi

			}
			ti["worker"] = wis

			tis[i] = ti
		}

		return tis
	}

	return make([]interface{}, 0)
}
