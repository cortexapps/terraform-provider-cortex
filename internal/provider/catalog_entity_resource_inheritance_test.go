package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogEntityResourceWithInheritance(t *testing.T) {
	tag := "domain-test-with-inheritance"
	resourceName := "cortex_catalog_entity." + tag
	name := "Test Domain With Ownership Inheritance"
	description := "Domain with owners that have inheritance configured"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceWithInheritance(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "type", "domain"),

					// Check owners
					resource.TestCheckResourceAttr(resourceName, "owners.#", "3"),

					// Email owner with APPEND inheritance
					resource.TestCheckResourceAttr(resourceName, "owners.0.type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceName, "owners.0.name", "Domain Owner"),
					resource.TestCheckResourceAttr(resourceName, "owners.0.email", "owner@example.com"),
					resource.TestCheckResourceAttr(resourceName, "owners.0.inheritance", "APPEND"),

					// Group owner with FALLBACK inheritance
					resource.TestCheckResourceAttr(resourceName, "owners.1.type", "GROUP"),
					resource.TestCheckResourceAttr(resourceName, "owners.1.name", "Platform Team"),
					resource.TestCheckResourceAttr(resourceName, "owners.1.provider", "CORTEX"),
					resource.TestCheckResourceAttr(resourceName, "owners.1.inheritance", "FALLBACK"),

					// Slack owner with NONE inheritance
					resource.TestCheckResourceAttr(resourceName, "owners.2.type", "SLACK"),
					resource.TestCheckResourceAttr(resourceName, "owners.2.channel", "platform"),
					resource.TestCheckResourceAttr(resourceName, "owners.2.inheritance", "NONE"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing - change inheritance values
			{
				Config: testAccCatalogEntityResourceWithInheritanceUpdated(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),

					// Check updated inheritance values
					resource.TestCheckResourceAttr(resourceName, "owners.0.inheritance", "FALLBACK"),
					resource.TestCheckResourceAttr(resourceName, "owners.1.inheritance", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "owners.2.inheritance", "APPEND"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceWithInheritance(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" %[1]q {
  tag         = %[1]q
  name        = %[2]q
  description = %[3]q
  type        = "domain"

  owners = [
    {
      type        = "EMAIL"
      name        = "Domain Owner"
      email       = "owner@example.com"
      inheritance = "APPEND"
    },
    {
      type        = "GROUP"
      name        = "Platform Team"
      provider    = "CORTEX"
      inheritance = "FALLBACK"
    },
    {
      type        = "SLACK"
      channel     = "platform"
      inheritance = "NONE"
    }
  ]
}`, tag, name, description)
}

func testAccCatalogEntityResourceWithInheritanceUpdated(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" %[1]q {
  tag         = %[1]q
  name        = %[2]q
  description = %[3]q
  type        = "domain"

  owners = [
    {
      type        = "EMAIL"
      name        = "Domain Owner"
      email       = "owner@example.com"
      inheritance = "FALLBACK"
    },
    {
      type        = "GROUP"
      name        = "Platform Team"
      provider    = "CORTEX"
      inheritance = "NONE"
    },
    {
      type        = "SLACK"
      channel     = "platform"
      inheritance = "APPEND"
    }
  ]
}`, tag, name, description)
}

func TestAccCatalogEntityResourceWithoutInheritance(t *testing.T) {
	tag := "entity-test-without-inheritance"
	resourceName := "cortex_catalog_entity." + tag
	name := "Test Entity Without Inheritance"
	description := "Entity with owners that have no inheritance configured"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityResourceWithoutInheritance(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),

					// Check that owners exist but inheritance is not set (should be empty/null)
					resource.TestCheckResourceAttr(resourceName, "owners.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "owners.0.type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceName, "owners.1.type", "GROUP"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityResourceWithoutInheritance(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" %[1]q {
  tag         = %[1]q
  name        = %[2]q
  description = %[3]q

  owners = [
    {
      type  = "EMAIL"
      name  = "Owner"
      email = "owner@example.com"
    },
    {
      type     = "GROUP"
      name     = "Team"
      provider = "CORTEX"
    }
  ]
}`, tag, name, description)
}

func TestAccCatalogEntityResourceInheritanceValidation(t *testing.T) {
	tag := "entity-test-invalid-inheritance"
	name := "Test Entity Invalid Inheritance"
	description := "Entity with invalid inheritance value"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid inheritance value - should fail validation
			{
				Config:      testAccCatalogEntityResourceInvalidInheritance(tag, name, description),
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
		},
	})
}

func testAccCatalogEntityResourceInvalidInheritance(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" %[1]q {
  tag         = %[1]q
  name        = %[2]q
  description = %[3]q

  owners = [
    {
      type        = "EMAIL"
      name        = "Owner"
      email       = "owner@example.com"
      inheritance = "INVALID_VALUE"
    }
  ]
}`, tag, name, description)
}
