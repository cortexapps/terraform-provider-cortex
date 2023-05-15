package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScorecardResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccScorecardResourceConfig("dora"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_scorecard.dora", "tag", "dora"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_scorecard.dora",
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
				Config: testAccScorecardResourceConfig("dora"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_scorecard.dora", "tag", "dora"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccScorecardResourceConfig(tag string) string {
	return fmt.Sprintf(`
resource "cortex_scorecard" "dora" {
  tag = %[1]q
}
`, tag)
}
