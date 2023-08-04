package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

type testScorecardResource struct {
	Tag         string
	Name        string
	Description string
	Draft       bool
}

/***********************************************************************************************************************
 * Helper methods
 **********************************************************************************************************************/

func (t *testScorecardResource) ResourceFullName() string {
	return t.ResourceType() + "." + t.Tag
}

func (t *testScorecardResource) ResourceType() string {
	return "cortex_scorecard"
}

func (t *testScorecardResource) ToTerraform() string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  tag = %[2]q
  name = %[3]q
  description = %[4]q
  draft = %[5]t
  rules = [
    {
      title = "Has a Description"
      expression = "description != null"
      weight = 1
      level = "Bronze"
      failure_message = "The description is required"
      description = "The service has a description"
    }
  ]
  ladder = {
    levels = [
      {
         name = "Bronze"
         rank = 1
         color = "#c38b5f"
      }
    ]
  }
  filter = {
    category = "SERVICE"
    query = "owners_is_set"
  }
  evaluation = {
    window = 24
  }
}`, t.ResourceType(), t.Tag, t.Name, t.Description, t.Draft)
}

/***********************************************************************************************************************
 * Tests
 **********************************************************************************************************************/

func TestAccScorecardResourceComplete(t *testing.T) {
	stub := testScorecardResource{
		Tag:  "test-complete-scorecard",
		Name: "Test Scorecard - Complete",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "name", stub.Name),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "description", stub.Description),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "draft", fmt.Sprintf("%t", stub.Draft)),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.title", "Has a Description"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.expression", "description != null"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.weight", "1"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.level", "Bronze"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.failure_message", "The description is required"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.description", "The service has a description"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.name", "Bronze"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.rank", "1"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.color", "#c38b5f"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "filter.category", "SERVICE"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "filter.query", "owners_is_set"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "evaluation.window", "24"),
				),
			},
			// ImportState testing
			{
				ResourceName:      stub.ResourceFullName(),
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "name", stub.Name),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "description", stub.Description),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "draft", fmt.Sprintf("%t", stub.Draft)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
