package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogResourceMinimal(t *testing.T) {
	slug := "test-catalog-minimal"
	resourceName := "cortex_catalog." + slug
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogResourceMinimalConfig(slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "slug", slug),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Catalog Minimal"),
					resource.TestCheckResourceAttr(resourceName, "icon_tag", "cortex"),
					resource.TestCheckResourceAttr(resourceName, "is_draft", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete occurs automatically
		},
	})
}

func testAccCatalogResourceMinimalConfig(slug string) string {
	return fmt.Sprintf(`
resource "cortex_catalog" %[1]q {
  slug     = %[1]q
  name     = "Test Catalog Minimal"
  icon_tag = "cortex"
  is_draft = false
}`, slug)
}
