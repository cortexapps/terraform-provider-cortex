package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccAwsIntegrationTypesResource(t *testing.T) {
	resourceName := "cortex_aws_integration_types.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAwsIntegrationTypesResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "types.0.type", "ec2"),
					resource.TestCheckResourceAttr(resourceName, "types.0.configured", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing
			{
				Config: testAccAwsIntegrationTypesResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "types.0.type", "ec2"),
					resource.TestCheckResourceAttr(resourceName, "types.0.configured", "true"),
				),
			},
		},
	})
}

func testAccAwsIntegrationTypesResourceConfig() string {
	return fmt.Sprintf(`
resource "cortex_aws_integration_types" "test" {
  types = [
    {
      type       = "ec2"
      configured = true
    }
  ]
}
`)
}
