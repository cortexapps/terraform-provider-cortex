package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamResource(t *testing.T) {
	r := tFactoryBuildTeamResource("Platform Engineering")
	r2 := tFactoryBuildTeamResource("Platform Engineering With Changed Name")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: r.String(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// core attributes
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "tag", "platform-engineering"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "name", "Platform Engineering"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "description", "Platform Engineering Team"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "summary", "The Cortex Platform Engineering Team"),
					// links
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "link.0.name", "GitHub"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "link.0.url", "https://github.com/cortexapp/cortex"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "link.0.type", "documentation"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "link.0.description", "GitHub repository for Cortex"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_team.platform_engineering",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"tag", "defaulted"},
			},
			// Update and Read testing
			{
				Config: r2.String(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// core attributes
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "tag", "platform-engineering"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "name", "Platform Engineering"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "description", "Platform Engineering Team"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "summary", "The Cortex Platform Engineering Team"),
					// links
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "link.0.name", "GitHub"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "link.0.url", "https://github.com/cortexapp/cortex"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "link.0.type", "documentation"),
					resource.TestCheckResourceAttr("cortex_team.platform_engineering", "link.0.description", "GitHub repository for Cortex"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

type TestTeamResource struct {
	Tag         string
	Name        string
	Description string
	Summary     string
	Links       []TestTeamLinkResource
}

type TestTeamLinkResource struct {
	Name        string
	URL         string
	Type        string
	Description string
}

func (r TestTeamResource) String() string {
	return fmt.Sprintf(`
resource "cortex_team" "platform_engineering" {
  tag = %[1]q
  name = %[2]q
  description = %[3]q
  summary = %[4]q
  link {
	name = %[5]q
	url = %[6]q
	type = %[7]q
	description = %[8]q
  }
}
`, r.Tag, r.Name, r.Description, r.Summary, r.Links[0].Name, r.Links[0].URL, r.Links[0].Type, r.Links[0].Description)
}

func tFactoryBuildTeamResource(name string) TestTeamResource {
	return TestTeamResource{
		Tag:         "platform_engineering",
		Name:        name,
		Description: "Platform Engineering Team",
		Summary:     "Platform Engineering Team",
		Links: []TestTeamLinkResource{
			{
				Name:        "GitHub",
				URL:         "https://github.com/cortexapp/cortex",
				Type:        "documentation",
				Description: "GitHub repository for Cortex",
			},
		},
	}
}
