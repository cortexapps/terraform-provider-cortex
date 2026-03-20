package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTeamsDataSource(t *testing.T) {
	recordName := "data.cortex_teams.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing - list all teams
			{
				Config: testAccTeamsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(recordName, "id"),
					resource.TestCheckResourceAttrSet(recordName, "teams.#"),
				),
			},
		},
	})
}

func testAccTeamsDataSourceConfig() string {
	return `
data "cortex_teams" "test" {
}
`
}
