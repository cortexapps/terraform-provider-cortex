package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

type testScorecardDataSource struct {
	Tag  string
	Name string
}

/***********************************************************************************************************************
 * Helper methods
 **********************************************************************************************************************/

func (t *testScorecardDataSource) ResourceFullName() string {
	return t.ResourceType() + "." + t.Tag
}

func (t *testScorecardDataSource) ResourceType() string {
	return "cortex_scorecard"
}

func (t *testScorecardDataSource) ToTerraform() string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
  tag = %[2]q
  name = %[3]q
  rules = [
    {
      title = "Has a Description"
      expression = "description != null"
      weight = 1
      level = "Bronze"
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
}`, t.ResourceType(), t.Tag, t.Name)
}

/***********************************************************************************************************************
 * Tests
 **********************************************************************************************************************/

func TestAccScorecardResourceMinimal(t *testing.T) {
	stub := testScorecardDataSource{
		Tag:  "test-minimal-scorecard",
		Name: "Test Scorecard - Minimal",
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
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
