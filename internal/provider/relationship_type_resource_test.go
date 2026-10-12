package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRelationshipTypeResourceMinimal(t *testing.T) {
	tag := "test-reltype-minimal"
	resourceType := "cortex_relationship_type"
	resourceName := resourceType + "." + tag

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRelationshipTypeResourceMinimalConfig(tag, "Minimal Relationship Type"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", "Minimal Relationship Type"),
					resource.TestCheckResourceAttr(resourceName, "allow_cycles", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_single_source", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_single_destination", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceType + "." + tag,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing — change name, verify no phantom diff
			{
				Config: testAccRelationshipTypeResourceMinimalConfig(tag, "Minimal Relationship Type Updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", "Minimal Relationship Type Updated"),
					resource.TestCheckResourceAttr(resourceName, "allow_cycles", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRelationshipTypeResourceMinimalConfig(tag string, name string) string {
	return fmt.Sprintf(`
resource "cortex_relationship_type" %[1]q {
  tag                   = %[1]q
  name                  = %[2]q
  allow_cycles          = false
  is_single_source      = false
  is_single_destination = false
  definition_location   = "SOURCE"

  sources_filter = {
    include = true
    types   = []
  }

  destinations_filter = {
    include = true
    types   = []
  }
}
`, tag, name)
}

func TestAccRelationshipTypeResourceFull(t *testing.T) {
	tag := "test-reltype-full"
	resourceType := "cortex_relationship_type"
	resourceName := resourceType + "." + tag

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with all optional fields
			{
				Config: testAccRelationshipTypeResourceFullConfig(tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", "Full Relationship Type"),
					resource.TestCheckResourceAttr(resourceName, "description", "A full relationship type with all optional fields."),
					resource.TestCheckResourceAttr(resourceName, "allow_cycles", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_single_source", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_single_destination", "false"),
					resource.TestCheckResourceAttr(resourceName, "source_label_singular", "Service"),
					resource.TestCheckResourceAttr(resourceName, "source_label_plural", "Services"),
					resource.TestCheckResourceAttr(resourceName, "destination_label_singular", "Database"),
					resource.TestCheckResourceAttr(resourceName, "destination_label_plural", "Databases"),
					resource.TestCheckResourceAttr(resourceName, "inheritances.#", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceType + "." + tag,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify plan is clean after apply (no diff from unordered set fields)
			{
				Config:   testAccRelationshipTypeResourceFullConfig(tag),
				PlanOnly: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRelationshipTypeResourceFullConfig(tag string) string {
	return fmt.Sprintf(`
resource "cortex_relationship_type" %[1]q {
  tag                   = %[1]q
  name                  = "Full Relationship Type"
  description           = "A full relationship type with all optional fields."
  allow_cycles          = false
  is_single_source      = false
  is_single_destination = false
  definition_location   = "SOURCE"

  sources_filter = {
    include = true
    types   = ["service"]
  }

  destinations_filter = {
    include = true
    types   = ["resource"]
  }

  source_label_singular      = "Service"
  source_label_plural        = "Services"
  destination_label_singular = "Database"
  destination_label_plural   = "Databases"

  inheritances = [
    {
      parameter_tag    = "x-cortex-owners"
      inheritance_type = "FALLBACK"
    }
  ]
}
`, tag)
}
