package provider_test

import (
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type TestTeamDataSource struct {
	Tag           string
	Members       []cortex.TeamMember
	SlackChannels []cortex.TeamSlackChannel
}

func (t *TestTeamDataSource) ToStubTeam() *cortex.Team {
	return &cortex.Team{
		TeamTag: t.Tag,
		CortexTeam: cortex.TeamCortexManaged{
			Members: t.Members,
		},
		SlackChannels: t.SlackChannels,
	}
}

func testAccTeamDataSourceConfig(tag string) string {
	return `
data "cortex_team" "test-team" {
  tag = "` + tag + `"
}
`
}

func TestAccTeamDataSource(t *testing.T) {
	stub := TestTeamDataSource{
		Tag: "test-team",
		Members: []cortex.TeamMember{
			{Name: "Test User", Email: "test@example.com"},
		},
		SlackChannels: []cortex.TeamSlackChannel{
			{Name: "test-channel", NotificationsEnabled: true},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceConfig(stub.Tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cortex_team.test-team", "tag", stub.Tag),
					resource.TestCheckResourceAttr("data.cortex_team.test-team", "members.0.name", stub.Members[0].Name),
					resource.TestCheckResourceAttr("data.cortex_team.test-team", "members.0.email", stub.Members[0].Email),
					resource.TestCheckResourceAttr("data.cortex_team.test-team", "slack_channels.0.name", stub.SlackChannels[0].Name),
					resource.TestCheckResourceAttr("data.cortex_team.test-team", "slack_channels.0.notifications_enabled", "true"),
				),
			},
		},
	})
}
