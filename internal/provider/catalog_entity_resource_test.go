package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogEntityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("catalog_entity.test", "name", "one"),
					resource.TestCheckResourceAttr("catalog_entity.test", "description", "example value when not configured"),
					resource.TestCheckResourceAttr("catalog_entity.test", "tag", "test"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "catalog_entity.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"name", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccCatalogEntityResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("catalog_entity.test", "name", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "catalog_entity" "test" {
  name = %[1]q
}
`, name)
}
