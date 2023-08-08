package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntityDomainMinimal(t *testing.T) {
	tag := "domain-test-minimal"
	resourceName := "cortex_catalog_entity." + tag
	name := "Test Minimal Domain"
	description := "Minimal Domain configuration for testing"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityDomainMinimal(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "type", "domain"),
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
				Config: testAccCatalogEntityDomainMinimal(tag, name, description+" 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description+" 2"),
					resource.TestCheckResourceAttr(resourceName, "type", "domain"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityDomainMinimal(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" %[1]q {
 tag = %[1]q
 name = %[2]q
 description = %[3]q
 type = "domain"
}`, tag, name, description)
}

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

func TestAccCatalogEntityDomainWithParent(t *testing.T) {
	tag := "domain-test-with-parent"
	resourceName := "cortex_catalog_entity." + tag
	name := "Test Domain With Parent"
	description := "Domain with parents defined"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityDomainWithParent(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "type", "domain"),
					resource.TestCheckResourceAttr(resourceName, "domain_parents.0.tag", "manual-test-domain"),
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
				Config: testAccCatalogEntityDomainWithParent(tag, name, description+" 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description+" 2"),
					resource.TestCheckResourceAttr(resourceName, "type", "domain"),
					resource.TestCheckResourceAttr(resourceName, "domain_parents.0.tag", "manual-test-domain"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityDomainWithParent(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" %[1]q {
 tag = %[1]q
 name = %[2]q
 description = %[3]q
 type = "domain"
 domain_parents = [
   {
     tag = "manual-test-domain"
   }
 ]

}`, tag, name, description)
}
