package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDepartmentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDepartmentResourceConfig("engineering"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_department.engineering", "tag", "engineering"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_department.engineering",
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
				Config: testAccDepartmentResourceConfig("engineering"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_department.engineering", "tag", "engineering"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDepartmentResourceConfig(tag string) string {
	return fmt.Sprintf(`
resource "cortex_department" "engineering" {
  tag = %[1]q
}
`, tag)
}
