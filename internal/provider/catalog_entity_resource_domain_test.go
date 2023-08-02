package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntityDomainWithChildren(t *testing.T) {
	tag := "domain-test-with-children"
	resourceName := "cortex_catalog_entity." + tag
	name := "Test Domain With Children"
	description := "Domain with children defined"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityDomainWithChildren(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "type", "domain"),
					resource.TestCheckResourceAttr(resourceName, "children.0.tag", "manual-test"),
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
				Config: testAccCatalogEntityDomainWithChildren(tag, name, description+" 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description+" 2"),
					resource.TestCheckResourceAttr(resourceName, "type", "domain"),
					resource.TestCheckResourceAttr(resourceName, "children.0.tag", "manual-test"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityDomainWithChildren(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" %[1]q {
 tag = %[1]q
 name = %[2]q
 description = %[3]q
 type = "domain"
 children = [
  {
    tag = "manual-test"
   }
 ]

}`, tag, name, description)
}
