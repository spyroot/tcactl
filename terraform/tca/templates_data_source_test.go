package tca

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const testAccDataSourceScaffolding = `
data "tca_data_source" "foo" {
  sample_attribute = "bar"
}
`

func TestdataSourceTemplates(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceScaffolding,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.tca_templates.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}
