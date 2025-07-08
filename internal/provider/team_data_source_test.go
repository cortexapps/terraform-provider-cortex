package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccTeamDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cortex_team.engineering", "tag", "platform_engineering"),
					resource.TestCheckResourceAttrSet("data.cortex_team.engineering", "members.0.name"),
					resource.TestCheckResourceAttrSet("data.cortex_team.engineering", "members.0.email"),
					resource.TestCheckResourceAttrSet("data.cortex_team.engineering", "slack_channels.0.name"),
					resource.TestCheckResourceAttrSet("data.cortex_team.engineering", "slack_channels.0.notifications_enabled"),
				),
			},
		},
	})
}

const testAccTeamDataSourceConfig = `
data "cortex_team" "engineering" {
  tag = "platform_engineering"
}
`
