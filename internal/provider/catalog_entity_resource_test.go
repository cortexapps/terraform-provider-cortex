package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceConfig("test", "A Test Service", "A test service for the Terraform provider"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "tag", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "name", "A Test Service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "description", "A test service for the Terraform provider"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.#", "3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_catalog_entity.test",
				ImportState:       true,
				ImportStateVerify: false,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"tag", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccCatalogEntityResourceConfig("test", "A Test Service", "A test service for the Terraform provider 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "tag", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "name", "A Test Service"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "description", "A test service for the Terraform provider 2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceConfig(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "test" {
 tag = %[1]q
 name = %[2]q
 description = %[3]q
 
 owners = [
    {
      type  = "EMAIL"
      name  = "John Doe"
      email = "john.doe@cortex.io"
    },
    {
      type     = "GROUP"
      name     = "Engineering"
      provider = "CORTEX"
    },
    {
      type                  = "SLACK"
      channel               = "engineering"
      notifications_enabled = false
    }
 ]

  groups = [
   "test",
   "test2"
  ]
}
`, tag, name, description)
}
