package tca

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spyroot/tcactl/cmd/api"
	"github.com/spyroot/tcactl/cmd/client"
	"log"
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
	//d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
