package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJiraAssetsObjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "jiraassets_object" "test" {
					type_id = ""
					attributes = [
						{
							attr_type_id = ""
							attr_value = ""
						}
					]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jiraassets_object.test", "id", "1"),
					resource.TestCheckResourceAttr("jiraassets_object.test", "attribute", "1"),
				),
			},
		},
	})
}
