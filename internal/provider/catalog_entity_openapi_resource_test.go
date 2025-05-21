package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogEntityOpenAPIResource(t *testing.T) {
	entityTag := "manual-test"
	spec := `
info:
  title: Test API
  version: 1.0.0
openapi: 3.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: OK`

	updatedSpec := `
info:
  title: Updated Test API
  version: 1.0.1
openapi: 3.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: OK`

	resourceName := "cortex_catalog_entity_openapi." + entityTag

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			// TODO:  This is currently failing due to some perma-drift that shows up on refresh.
			{
				Config: testAccCatalogEntityOpenAPIResourceConfig(entityTag, spec),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "entity_tag", entityTag),
					resource.TestCheckResourceAttr(resourceName, "spec", spec),
					resource.TestCheckResourceAttr(resourceName, "id", entityTag),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCatalogEntityOpenAPIResourceConfig(entityTag, updatedSpec),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "entity_tag", entityTag),
					resource.TestCheckResourceAttr(resourceName, "spec", updatedSpec),
				),
			},
		},
	})
}

func testAccCatalogEntityOpenAPIResourceConfig(entityTag string, spec string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity_openapi" %[1]q {
  entity_tag = %[1]q
  spec       = %[2]q
}
`, entityTag, spec)
}
