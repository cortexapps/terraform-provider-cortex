package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
      expression = "entity.description() != null"
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
    query = "owners_is_set"
  }
  evaluation = {
    window = 24
  }
}`, t.ResourceType(), t.Tag, t.Name, t.Description, t.Draft)
}

func (t *testScorecardResource) ToTerraformWithoutDescriptionOrFilter() string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  tag = %[2]q
  name = %[3]q
  draft = %[4]t
  rules = [
    {
      title = "Has a Description"
      expression = "entity.description() != null"
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
  evaluation = {
    window = 24
  }

  lifecycle {
    ignore_changes = [description, filter]
  }
}`, t.ResourceType(), t.Tag, t.Name, t.Draft)
}

func (t *testScorecardResource) ToTerraformWithoutDraft() string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  tag = %[2]q
  name = %[3]q
  description = %[4]q
  rules = [
    {
      title = "Has a Description"
      expression = "entity.description() != null"
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
      },
      {
         name = "Silver"
         rank = 2
		 color = "#c3c3c3"
	  }
    ]
  }
  filter = {
    types = {
      include = ["service", "app"]
    }
    query = "owners_is_set"
  }
  evaluation = {
    window = 24
  }
}`, t.ResourceType(), t.Tag, t.Name, t.Description)
}

func (t *testScorecardResource) ToTerraformWithFilterGroups() string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  tag = %[2]q
  name = %[3]q
  description = %[4]q
  draft = %[5]t
  rules = [
    {
      title = "Has a Description"
      expression = "entity.description() != null"
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
    types = {
      include = ["service"]
    }
    groups = {
      include = ["team-a", "team-b"]
    }
    query = "owners_is_set"
  }
  evaluation = {
    window = 24
  }
}`, t.ResourceType(), t.Tag, t.Name, t.Description, t.Draft)
}

func (t *testScorecardResource) ToTerraformWithFilterTypesExclude() string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  tag = %[2]q
  name = %[3]q
  description = %[4]q
  draft = %[5]t
  rules = [
    {
      title = "Has a Description"
      expression = "entity.description() != null"
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
    types = {
      exclude = ["deprecated"]
    }
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
			// Read testing without description or filter
			{
				Config: stub.ToTerraformWithoutDescriptionOrFilter(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "name", stub.Name),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "draft", fmt.Sprintf("%t", stub.Draft)),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.title", "Has a Description"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.expression", "entity.description() != null"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.weight", "1"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.level", "Bronze"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.failure_message", "The description is required"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.description", "The service has a description"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.name", "Bronze"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.rank", "1"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.color", "#c38b5f"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "evaluation.window", "24"),
				),
			},
			// Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "name", stub.Name),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "description", stub.Description),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "draft", fmt.Sprintf("%t", stub.Draft)),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.title", "Has a Description"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.expression", "entity.description() != null"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.weight", "1"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.level", "Bronze"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.failure_message", "The description is required"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.description", "The service has a description"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.name", "Bronze"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.rank", "1"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.color", "#c38b5f"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "filter.query", "owners_is_set"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "evaluation.window", "24"),
				),
			},
			// Read testing with types filter
			{
				Config: stub.ToTerraformWithoutDraft(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "name", stub.Name),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "description", stub.Description),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "draft", "false"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.title", "Has a Description"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.expression", "entity.description() != null"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.weight", "1"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.level", "Bronze"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.failure_message", "The description is required"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "rules.0.description", "The service has a description"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.name", "Bronze"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.rank", "1"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.0.color", "#c38b5f"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.1.name", "Silver"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.1.rank", "2"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "ladder.levels.1.color", "#c3c3c3"),

					resource.TestCheckTypeSetElemAttr(stub.ResourceFullName(), "filter.types.include.*", "service"),
					resource.TestCheckTypeSetElemAttr(stub.ResourceFullName(), "filter.types.include.*", "app"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "filter.query", "owners_is_set"),

					resource.TestCheckResourceAttr(stub.ResourceFullName(), "evaluation.window", "24"),
				),
			},
			// Read testing with types and groups filter
			{
				Config: stub.ToTerraformWithFilterGroups(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "name", stub.Name),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "description", stub.Description),

					resource.TestCheckTypeSetElemAttr(stub.ResourceFullName(), "filter.types.include.*", "service"),
					resource.TestCheckTypeSetElemAttr(stub.ResourceFullName(), "filter.groups.include.*", "team-a"),
					resource.TestCheckTypeSetElemAttr(stub.ResourceFullName(), "filter.groups.include.*", "team-b"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "filter.query", "owners_is_set"),
				),
			},
			// Read testing with types exclude filter
			{
				Config: stub.ToTerraformWithFilterTypesExclude(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "name", stub.Name),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "description", stub.Description),

					resource.TestCheckTypeSetElemAttr(stub.ResourceFullName(), "filter.types.exclude.*", "deprecated"),
					resource.TestCheckResourceAttr(stub.ResourceFullName(), "filter.query", "owners_is_set"),
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

func TestAccScorecardResourceFilterTypesValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "cortex_scorecard" "test" {
  tag  = "test-validation"
  name = "Validation Test"

  ladder = {
    levels = [
      {
        name  = "Bronze"
        rank  = 1
        color = "#c38b5f"
      }
    ]
  }

  rules = [
    {
      title      = "Has Description"
      expression = "entity.description() != null"
      weight     = 1
      level      = "Bronze"
    }
  ]

  filter = {
    types = {
      include = ["service"]
      exclude = ["deprecated"]
    }
  }
}`,
				ExpectError: regexp.MustCompile(`Attribute "filter\.types\.(include|exclude)" cannot be specified when`),
			},
		},
	})
}
