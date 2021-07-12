package tca

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spyroot/tcactl/lib/api"
	"github.com/spyroot/tcactl/lib/client/response"
	"time"
)

func resourceClusterTemplate() *schema.Resource {
	return &schema.Resource{
		//
		CreateContext: ClusterTemplateCreate,
		//
		ReadContext: ClusterTemplateRead,
		//
		UpdateContext: ClusterTemplateUpdate,
		//
		DeleteContext: ClusterTemplateDelete,
		//
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ClusterTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	_, err := GetApi(m)
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//
	//items := d.Get("items").([]interface{})
	//ois := []hc.OrderItem{}
	//
	//for _, item := range items {
	//	i := item.(map[string]interface{})
	//	co := i["coffee"].([]interface{})[0]
	//	coffee := co.(map[string]interface{})
	//
	//	oi := hc.OrderItem{
	//		Coffee: hc.Coffee{
	//			ID: coffee["id"].(int),
	//		},
	//		Quantity: i["quantity"].(int),
	//	}
	//
	//	ois = append(ois, oi)
	//}
	//
	//o, err := c.CreateOrder(ois)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//d.SetId(strconv.Itoa(o.ID))

	ClusterTemplateRead(ctx, d, m)
	return diags
}

// ClusterTemplateRead
func ClusterTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tca, err := GetApi(m)
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//	orderID := d.Id()

	templates, err := tca.GetClusterTemplates()
	if err != nil {
		return diag.FromErr(err)
	}

	templateItems := flattenTemplateItems(templates)
	if err := d.Set("templates", templateItems); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

//
func ClusterTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := GetApi(m)
	if err != nil {
		return diag.FromErr(err)
	}

	//orderID := d.Id()

	if d.HasChange("templates") {
		//items := d.Get("templates").([]interface{})
		//ois := []hc.OrderItem{}
		//
		//for _, item := range items {
		//	i := item.(map[string]interface{})
		//
		//	co := i["coffee"].([]interface{})[0]
		//	coffee := co.(map[string]interface{})
		//
		//	oi := hc.OrderItem{
		//		Coffee: hc.Coffee{
		//			ID: coffee["id"].(int),
		//		},
		//		Quantity: i["quantity"].(int),
		//	}
		//	ois = append(ois, oi)
		//}

		//_, err := c.UpdateOrder(orderID, ois)
		//if err != nil {
		//	return diag.FromErr(err)
		//}

		err := d.Set("last_updated", time.Now().Format(time.RFC850))
		if err != nil {
			return nil
		}
	}

	return ClusterTemplateRead(ctx, d, m)
}

func ClusterTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tca := m.(*api.TcaApi)
	if tca == nil {
		return diag.FromErr(fmt.Errorf("nil client"))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	templateId := d.Id()
	err := tca.DeleteTemplate(templateId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

//
func flattenTemplateItems(templates *response.ClusterTemplates) []interface{} {
	if templates != nil {
		ois := make([]interface{}, len(templates.ClusterTemplates), len(templates.ClusterTemplates))

		for i, orderItem := range templates.ClusterTemplates {
			oi := make(map[string]interface{})

			oi["id"] = orderItem.Id
			oi["name"] = orderItem.Name
			oi["cluster_type"] = orderItem.ClusterType
			//	oi["clusterConfig"] = orderItem.ClusterConfig
			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}

//flattenTemplate
func flattenTemplate(template *response.ClusterTemplateSpec) []interface{} {
	c := make(map[string]interface{})
	c["id"] = template.Id
	c["name"] = template.Name
	c["cluster_type"] = template.ClusterType

	return []interface{}{c}
}
