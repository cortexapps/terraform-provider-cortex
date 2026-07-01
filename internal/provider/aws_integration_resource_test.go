package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccAwsIntegrationResource(t *testing.T) {
	accountId := "123456789012"
	iamRole := "arn:aws:iam::123456789012:role/cortex-role"
	resourceName := "cortex_aws_integration_configuration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAwsIntegrationResourceConfig(accountId, iamRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_id", accountId),
					resource.TestCheckResourceAttr(resourceName, "iam_role", iamRole),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// AccountAlias is computed from the API, we will just ignore checking computed fields if they are missing
			},
			// Update and Read testing
			{
				Config: testAccAwsIntegrationResourceConfig(accountId, iamRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_id", accountId),
					resource.TestCheckResourceAttr(resourceName, "iam_role", iamRole),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAwsIntegrationResourceConfig(accountId string, iamRole string) string {
	return fmt.Sprintf(`
resource "cortex_aws_integration_configuration" "test" {
  account_id = %[1]q
  iam_role   = %[2]q
}
`, accountId, iamRole)
}
