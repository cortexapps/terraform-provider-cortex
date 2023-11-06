package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntityTeamResourceMinimal(t *testing.T) {
	tag := "test-team-minimal"
	resourceName := "cortex_catalog_entity." + tag
	name := "Minimal team"
	description := "Minimal team"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityTeamResourceMinimal(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "team.members.0.name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "team.members.0.email", "test@cortex.io"),
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
				Config: testAccCatalogEntityTeamResourceMinimal(tag, name, description+" 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description+" 2"),
					resource.TestCheckResourceAttr(resourceName, "team.members.0.name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "team.members.0.email", "test@cortex.io"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityTeamResourceMinimal(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "%s" {
  type = "team"
  tag = "%s"
  name = "%s"
  description = "%s"
  team = {
    members = [
      {
        name = "Test"
        email = "test@cortex.io"
      }
    ]
  }
}
`, tag, tag, name, description)
}

func TestAccCatalogEntityTeamResourceComplete(t *testing.T) {
	tag := "test-team-complete"
	resourceName := "cortex_catalog_entity." + tag
	name := "Complete team"
	description := "Complete team"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityTeamResourceComplete(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "team.members.0.name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "team.members.0.email", "test@cortex.io"),
					resource.TestCheckResourceAttr(resourceName, "slack.channels.0.name", "test-channel"),
					resource.TestCheckResourceAttr(resourceName, "slack.channels.0.notifications_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "children.0.tag", tag+"-child"),
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
				Config: testAccCatalogEntityTeamResourceComplete(tag, name, description+" 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description+" 2"),
					resource.TestCheckResourceAttr(resourceName, "team.members.0.name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "team.members.0.email", "test@cortex.io"),
					resource.TestCheckResourceAttr(resourceName, "slack.channels.0.name", "test-channel"),
					resource.TestCheckResourceAttr(resourceName, "slack.channels.0.notifications_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "children.0.tag", tag+"-child"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityTeamResourceComplete(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "%[1]s-child" {
  type = "team"
  tag = "%[1]s-child"
  name = "%[2]s Child"
  description = "%[3]s Child"
  team = {
    members = [{ name = "Test", email = "test@cortex.io" }]
  }
}
resource "cortex_catalog_entity" %[1]q {
  type = "team"
  tag = %[1]q
  name = %[2]q
  description = %[3]q
  team = {
    members = [
      {
        name = "Test"
        email = "test@cortex.io"
      }
    ]
  }
  slack = {
   channels = [
     { 
        name = "test-channel"
        notifications_enabled = false
     }
   ]
  }
  children = [
	{
	  type = "team"
	  tag = "%[1]s-child"
    }
  ]
}
`, tag, name, description)
}

func TestAccCatalogEntityTeamResourceWithChildren(t *testing.T) {
	tag := "test-team-with-children"
	resourceName := "cortex_catalog_entity." + tag
	name := "Test Children Team"
	description := "A test suite team with children"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCatalogEntityTeamResourceWithChildren(tag, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "children.0.tag", tag+"-child"),
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
				Config: testAccCatalogEntityTeamResourceWithChildren(tag, name, description+" 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", tag),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description+" 2"),
					resource.TestCheckResourceAttr(resourceName, "children.0.tag", tag+"-child"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCatalogEntityTeamResourceWithChildren(tag string, name string, description string) string {
	return fmt.Sprintf(`
resource "cortex_catalog_entity" "%[1]s-child" {
  type = "team"
  tag = "%[1]s-child"
  name = "%[2]s Child"
  description = "%[3]s Child"
  team = {
    members = [{ name = "Test", email = "test@cortex.io" }]
  }
}
resource "cortex_catalog_entity" %[1]q {
  type = "team"
  tag = %[1]q
  name = %[2]q
  description = %[3]q
  team = {
    members = [
      {
        name  = "Test"
        email = "test@cortex.io"
        role  = "engineering-manager"
      }
    ]
  }
  children = [
	{
	  type = "team"
	  tag = "%[1]s-child"
    }
  ]
}
`, tag, name, description)
}
