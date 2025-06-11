package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/google/go-cmp/cmp"
)

func TestAccCatalogEntityOpenAPIResource(t *testing.T) {
	entityTag := "test-service"
	spec := `{"info":{"title":"Test API","version":"1.0.0"},"paths":{"/test":{"get":{"responses":{"200":{"description":"OK"}}}}},"openapi":"3.0.0","servers":[{"url":"/"}]}`

	updatedSpec := `{"info":{"title":"UpdatedTest API","version":"1.0.1"},"paths":{"/test":{"get":{"responses":{"200":{"description":"OK"}}}}},"openapi":"3.0.0","servers":[{"url":"/"}]}`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityOpenAPIResourceConfig(entityTag, spec),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							rs := s.RootModule().Resources["cortex_catalog_entity_openapi.test"]
							actual := rs.Primary.Attributes["spec"]
							expected := spec

							if diff := cmp.Diff(expected, actual); diff != "" {
								t.Errorf("attribute mismatch (-expected +actual):\n%s", diff)
							}

							return nil
						},
					),
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
