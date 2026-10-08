package provider_test

import (
	"fmt"
	"regexp"
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
  type     = "DOMAIN"
}`, slug)
}

func TestAccCatalogDataSource(t *testing.T) {
	slug := "test-catalog-ds"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCatalogDataSourceConfig(slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cortex_catalog.test", "slug", slug),
					resource.TestCheckResourceAttr("data.cortex_catalog.test", "name", "Test Catalog DS"),
				),
			},
		},
	})
}

func TestAccCatalogResourceFilterConflict(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "cortex_catalog" "conflict" {
  slug     = "test-catalog-conflict"
  name     = "Conflict Test"
  icon_tag = "cortex"
  is_draft = false
  type     = "FILTER"
  filter = {
    types = {
      include = ["service"]
      exclude = ["domain"]
    }
  }
}`,
				ExpectError: regexp.MustCompile(`cannot be specified when`),
			},
		},
	})
}

func TestAccCatalogResourceWithFilter(t *testing.T) {
	slug := "test-catalog-filter"
	resourceName := "cortex_catalog." + slug
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCatalogResourceWithFilterConfig(slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "slug", slug),
					resource.TestCheckResourceAttr(resourceName, "type", "FILTER"),
					resource.TestCheckResourceAttr(resourceName, "filter.types.include.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "filter.types.include.*", "service"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCatalogResourceWithFilterConfig(slug string) string {
	return fmt.Sprintf(`
resource "cortex_catalog" %[1]q {
  slug     = %[1]q
  name     = "Test Catalog Filter"
  icon_tag = "cortex"
  is_draft = false
  type     = "FILTER"
  filter = {
    types = {
      include = ["service"]
    }
  }
}`, slug)
}

func testAccCatalogDataSourceConfig(slug string) string {
	return fmt.Sprintf(`
resource "cortex_catalog" %[1]q {
  slug     = %[1]q
  name     = "Test Catalog DS"
  icon_tag = "cortex"
  is_draft = false
  type     = "DOMAIN"
}

data "cortex_catalog" "test" {
  slug = cortex_catalog.%[1]s.slug
}`, slug)
}
