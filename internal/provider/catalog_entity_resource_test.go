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
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.name", "John Doe"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.0.email", "john.doe@cortex.io"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.type", "GROUP"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.name", "Engineering"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.1.provider", "CORTEX"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.type", "SLACK"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.channel", "engineering"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.2.notifications_enabled", "false"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.#", "2"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.0", "test"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.1", "test2"),

					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.#", "1"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.name", "Internal Docs"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.type", "documentation"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.0.url", "https://internal-docs.cortex.io/products-service"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_catalog_entity.test",
				ImportState:       true,
				ImportStateVerify: false, // TODO: Fix this
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
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "owners.#", "3"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "groups.#", "2"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "links.#", "1"),
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

  links = [
    {
      name = "Internal Docs"
      type = "documentation"
      url  = "https://internal-docs.cortex.io/products-service"
    }
  ]

  metadata = jsonencode({
	"my-key": "the value",
	"another-key": {
		"this": "is",
		"an": "object"
	},
	"final-key": [
		"also",
		"use",
		"lists!"
	]
  })
}
`, tag, name, description)
}
