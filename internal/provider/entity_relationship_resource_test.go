package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	erSourceTag  = "test-er-source"
	erDestTag1   = "test-er-dest-1"
	erDestTag2   = "test-er-dest-2"
	erDestTag3   = "test-er-dest-3"
	erRelTypeTag = "test-er-reltype"
)

func TestAccEntityRelationshipResource(t *testing.T) {
	resourceName := "cortex_entity_relationship.example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with 2 destinations
			{
				Config: testAccEntityRelationshipConfig([]string{erDestTag1, erDestTag2}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_entity_tag", erSourceTag),
					resource.TestCheckResourceAttr(resourceName, "relationship_type_tag", erRelTypeTag),
					resource.TestCheckResourceAttr(resourceName, "destination_entity_tags.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "destination_entity_tags.*", erDestTag1),
					resource.TestCheckTypeSetElemAttr(resourceName, "destination_entity_tags.*", erDestTag2),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Idempotency check after create
			{
				Config:   testAccEntityRelationshipConfig([]string{erDestTag1, erDestTag2}),
				PlanOnly: true,
			},
			// Update to 3 destinations
			{
				Config: testAccEntityRelationshipConfig([]string{erDestTag1, erDestTag2, erDestTag3}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_entity_tag", erSourceTag),
					resource.TestCheckResourceAttr(resourceName, "relationship_type_tag", erRelTypeTag),
					resource.TestCheckResourceAttr(resourceName, "destination_entity_tags.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "destination_entity_tags.*", erDestTag1),
					resource.TestCheckTypeSetElemAttr(resourceName, "destination_entity_tags.*", erDestTag2),
					resource.TestCheckTypeSetElemAttr(resourceName, "destination_entity_tags.*", erDestTag3),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to empty (clear all destinations)
			{
				Config: testAccEntityRelationshipConfig([]string{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_entity_tag", erSourceTag),
					resource.TestCheckResourceAttr(resourceName, "relationship_type_tag", erRelTypeTag),
					resource.TestCheckResourceAttr(resourceName, "destination_entity_tags.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccEntityRelationshipConfig builds a Terraform config that:
//  1. Creates the source entity, destination entities, and relationship type as dependencies.
//  2. Creates the cortex_entity_relationship with the given destination tags.
func testAccEntityRelationshipConfig(activeDests []string) string {
	destList := ""
	for _, d := range activeDests {
		destList += fmt.Sprintf("%q,\n    ", d)
	}

	return fmt.Sprintf(`
resource "cortex_catalog_entity" "source" {
  tag  = %[1]q
  name = "Entity Relationship Test Source"
  type = "service"
}

resource "cortex_catalog_entity" "dest1" {
  tag  = %[2]q
  name = "Entity Relationship Test Dest 1"
  type = "service"
}

resource "cortex_catalog_entity" "dest2" {
  tag  = %[3]q
  name = "Entity Relationship Test Dest 2"
  type = "service"
}

resource "cortex_catalog_entity" "dest3" {
  tag  = %[4]q
  name = "Entity Relationship Test Dest 3"
  type = "service"
}

resource "cortex_relationship_type" "reltype" {
  tag                   = %[5]q
  name                  = "Test Entity Relationship Type"
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

resource "cortex_entity_relationship" "example" {
  source_entity_tag       = cortex_catalog_entity.source.tag
  relationship_type_tag   = cortex_relationship_type.reltype.tag
  destination_entity_tags = [%[6]s]

  depends_on = [
    cortex_catalog_entity.source,
    cortex_catalog_entity.dest1,
    cortex_catalog_entity.dest2,
    cortex_catalog_entity.dest3,
    cortex_relationship_type.reltype,
  ]
}
`, erSourceTag, erDestTag1, erDestTag2, erDestTag3, erRelTypeTag, destList)
}
