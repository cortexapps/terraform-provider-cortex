package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogEntityOpenAPIResource(t *testing.T) {
	entityTag := "test-service"
	spec := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: OK`

	updatedSpec := `openapi: 3.0.0
info:
  title: Updated Test API
  version: 1.0.1
paths:
  /test:
    get:
      responses:
        '200':
          description: OK`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityOpenAPIResourceConfig(entityTag, spec),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity_openapi.test", "entity_tag", entityTag),
					resource.TestCheckResourceAttr("cortex_catalog_entity_openapi.test", "spec", spec),
					resource.TestCheckResourceAttr("cortex_catalog_entity_openapi.test", "id", entityTag),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_catalog_entity_openapi.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCatalogEntityOpenAPIResourceConfig(entityTag, updatedSpec),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity_openapi.test", "entity_tag", entityTag),
					resource.TestCheckResourceAttr("cortex_catalog_entity_openapi.test", "spec", updatedSpec),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityOpenAPIResourceConfig(entityTag string, spec string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity_openapi" "test" {
  entity_tag = %[1]q
  spec       = %[2]q
}
`, entityTag, spec)
}
