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

func dataMasterNodeNetworks() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nameservers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func dataMasterNodes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
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
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMasterNodeNetworks(),
			},
		},
	}
}

// dataSourceClusters - TCA cluster data source
func dataSourceClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterRead,
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vim_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vsphere_cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"management_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hcx_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active_tasks_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kube_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataMasterNodes(),
						},
					},
				},
			},
		},
	}
}

// dataSourceOrderRead
func dataSourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	glog.Infof("dataSourceClusterRead")
	//u, err := uuid.Parse(d.Get("id").(string))
	//if err != nil {
	//	panic(err)
	//}

	//d.Get
	//d.Get("id")
	//
	//orderID, ok := d.Get("id").(string)
	//
	clusters, err := tca.GetClusters()
	if err != nil {
		return diag.FromErr(err)
	}

	t := flattenClusterData(clusters)
	if err := d.Set("clusters", t); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	//d.SetId(orderID)

	return diags
}

// flattenClusterData - read clusters
func flattenClusterData(clusters *response.Clusters) []interface{} {

	if clusters != nil {
		ois := make([]interface{}, len(clusters.Clusters), len(clusters.Clusters))

		for i, t := range clusters.Clusters {
			oi := make(map[string]interface{})

			oi["id"] = t.Id
			oi["cluster_name"] = t.VimId
			oi["cluster_name"] = t.ClusterName
			oi["cluster_type"] = t.ClusterType
			oi["vsphere_cluster_name"] = t.VsphereClusterName
			oi["management_cluster_id"] = t.ManagementClusterId
			oi["hcx_uuid"] = t.HcxUUID
			oi["status"] = t.Status
			oi["active_tasks_count"] = t.ActiveTasksCount
			oi["cluster_id"] = t.ClusterId
			oi["cluster_url"] = t.ClusterUrl
			oi["kube_config"] = t.KubeConfig
			oi["endpoint_ip"] = t.EndpointIP
			oi["master_nodes"] = flattenMasterNodeData(&t.MasterNodes)
			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}

func flattenMasterNodeData(nodes *[]response.MasterNodesDetails) []interface{} {

	if nodes != nil {
		ois := make([]interface{}, len(*nodes), len(*nodes))

		for i, t := range *nodes {
			oi := make(map[string]interface{})

			oi["cpu"] = t.Cpu
			oi["memory"] = t.Memory
			oi["name"] = t.Name
			oi["storage"] = t.Storage
			oi["replica"] = t.Replica
			oi["clone_mode"] = t.CloneMode
			oi["networks"] = flattenMasterNodeNetworks(&t.Networks)

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}

func flattenMasterNodeNetworks(networks *[]response.ClusterNetwork) []interface{} {

	if networks != nil {
		ois := make([]interface{}, len(*networks), len(*networks))

		for i, t := range *networks {
			oi := make(map[string]interface{})

			oi["label"] = t.Label
			oi["network_name"] = t.NetworkName
			oi["nameservers"] = t.Nameservers
			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
