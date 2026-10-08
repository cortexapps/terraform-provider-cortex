package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccResourceDefinitionResourceMinimal(t *testing.T) {
	tag := "test-resource-definition-minimal"
	resourceType := "cortex_resource_definition"
	description := "A minimal resource definition."
	resourceName := resourceType + "." + tag
	stub := tFactoryBuildResourceDefinitionResource(tag, description)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceDefinitionResourceMinimalConfig(tag, stub),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", stub.Type),
					resource.TestCheckResourceAttr(resourceName, "name", stub.Name),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceType + "." + tag,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccResourceDefinitionResourceMinimalConfig(tag, stub),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", stub.Type),
					resource.TestCheckResourceAttr(resourceName, "name", stub.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccResourceDefinitionResourceMinimalConfig(resourceTag string, stub TestResourceDefinitionResource) string {
	return fmt.Sprintf(`
resource "cortex_resource_definition" %1[1]q {
 type = %[1]q
 name = %[2]q
 description = %[3]q
 schema = jsonencode({
	"properties": {
	  "region": {
		"type": "string"
	  }
	}
    "type": "object"
 })
}
`, resourceTag, stub.Type, stub.Name, stub.Description)
}

func TestAccResourceDefinitionResourceWithIconTag(t *testing.T) {
	tag := "test-resourcedef-icontag"
	resourceType := "cortex_resource_definition"
	resourceName := resourceType + "." + tag

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with icon_tag set
			{
				Config: fmt.Sprintf(`
resource "cortex_resource_definition" %[1]q {
  type     = %[1]q
  name     = %[1]q
  icon_tag = "cloud"
  schema   = jsonencode({
    "properties": {
      "region": {
        "type": "string"
      }
    }
    "type": "object"
  })
}
`, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", tag),
					resource.TestCheckResourceAttr(resourceName, "icon_tag", "cloud"),
				),
			},
			// Remove icon_tag — should produce no error and no phantom diff
			{
				Config: fmt.Sprintf(`
resource "cortex_resource_definition" %[1]q {
  type   = %[1]q
  name   = %[1]q
  schema = jsonencode({
    "properties": {
      "region": {
        "type": "string"
      }
    }
    "type": "object"
  })
}
`, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", tag),
					resource.TestCheckNoResourceAttr(resourceName, "icon_tag"),
				),
			},
		},
	})
}

type TestResourceDefinitionResource struct {
	Type        string
	Name        string
	Description string
	Schema      map[string]interface{}
}

func tFactoryBuildResourceDefinitionResource(tag string, description string) TestResourceDefinitionResource {
	return TestResourceDefinitionResource{
		Type:        tag,
		Name:        tag,
		Description: description,
		Schema: map[string]interface{}{
			"properties": map[string]interface{}{
				"region": map[string]interface{}{
					"type": "string",
				},
			},
			"type": "object",
		},
	}
}
