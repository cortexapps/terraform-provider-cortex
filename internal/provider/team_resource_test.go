package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamResource(t *testing.T) {
	r := tFactoryBuildTeamResource("Platform Engineering Team")
	//r2 := tFactoryBuildTeamResource("Platform Engineering Team with changed description")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: r.String(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// core attributes
					resource.TestCheckResourceAttr("cortex_team.platform-engineering", "tag", "platform-engineering"),
					resource.TestCheckResourceAttr("cortex_team.platform-engineering", "name", "Platform Engineering"),
					resource.TestCheckResourceAttr("cortex_team.platform-engineering", "description", "Platform Engineering Team"),
					resource.TestCheckResourceAttr("cortex_team.platform-engineering", "summary", "The Cortex Platform Engineering Team"),
					// links
					//resource.TestCheckResourceAttr("cortex_team.platform-engineering", "link.0.name", "GitHub"),
					//resource.TestCheckResourceAttr("cortex_team.platform-engineering", "link.0.url", "https://github.com/cortexapp/cortex"),
					//resource.TestCheckResourceAttr("cortex_team.platform-engineering", "link.0.type", "documentation"),
					//resource.TestCheckResourceAttr("cortex_team.platform-engineering", "link.0.description", "GitHub repository for Cortex"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_team.platform-engineering",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			// TODO: Fix this, Cortex API 400s here, may have to check with Cortex
			//{
			//	Config: r2.String(),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		// core attributes
			//		resource.TestCheckResourceAttr("cortex_team.platform-engineering", "tag", "platform-engineering"),
			//		resource.TestCheckResourceAttr("cortex_team.platform-engineering", "name", "Platform Engineering"),
			//		resource.TestCheckResourceAttr("cortex_team.platform-engineering", "description", "Platform Engineering Team with changed description"),
			//		resource.TestCheckResourceAttr("cortex_team.platform-engineering", "summary", "The Cortex Platform Engineering Team"),
			//		// links
			//		//resource.TestCheckResourceAttr("cortex_team.platform-engineering", "link.0.name", "GitHub"),
			//		//resource.TestCheckResourceAttr("cortex_team.platform-engineering", "link.0.url", "https://github.com/cortexapp/cortex"),
			//		//resource.TestCheckResourceAttr("cortex_team.platform-engineering", "link.0.type", "documentation"),
			//		//resource.TestCheckResourceAttr("cortex_team.platform-engineering", "link.0.description", "GitHub repository for Cortex"),
			//	),
			//},
			// Delete testing automatically occurs in TestCase
		},
	})
}

type TestTeamResource struct {
	Tag         string
	Name        string
	Description string
	Summary     string
}

func (r TestTeamResource) String() string {
	return fmt.Sprintf(`
resource "cortex_team" "platform-engineering" {
  	tag = %[1]q
  	name = %[2]q
  	description = %[3]q
  	summary = %[4]q
  	slack_channels = [
		{
      		name = "#platform-engineering"
	  		notifications_enabled = false
		}
  	]
  	links = [
		{
			name = "Internal Docs"
			url = "https://internal-docs.cortex.io"
			type = "documentation"
			description = "Internal documentation for Cortex"
		}
	]
	additional_members = [
		{
			name = "John Doe"
			description = "Sir John"
			email = "john.doe@cortex.io"
		}
	]
}
`, r.Tag, r.Name, r.Description, r.Summary)
}

func tFactoryBuildTeamResource(description string) TestTeamResource {
	return TestTeamResource{
		Tag:         "platform-engineering",
		Name:        "Platform Engineering",
		Description: description,
		Summary:     "The Cortex Platform Engineering Team",
	}
}
