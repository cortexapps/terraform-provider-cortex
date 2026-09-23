package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The public API does not check the keys against Datadog on create, so fake keys work here.
const (
	testDatadogApiKey = "tf-acc-fake-api-key-a1b2"
	testDatadogAppKey = "tf-acc-fake-app-key-c3d4"
)

func TestAccDatadogConfigurationResource(t *testing.T) {
	resourceName := "cortex_datadog_configuration.test"
	alias := "tf-acc-test-datadog"
	renamedAlias := "tf-acc-test-datadog-renamed"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatadogConfigurationResourceConfig(alias, `["tf-acc"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", alias),
					resource.TestCheckResourceAttr(resourceName, "alias", alias),
					resource.TestCheckResourceAttr(resourceName, "region", "US1"),
					resource.TestCheckResourceAttr(resourceName, "environments.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "environments.0", "tf-acc"),
					resource.TestCheckResourceAttr(resourceName, "api_key_last_four", "a1b2"),
					resource.TestCheckResourceAttr(resourceName, "app_key_last_four", "c3d4"),
					resource.TestCheckResourceAttrSet(resourceName, "is_default"),
				),
			},
			// ImportState testing. The API never returns the keys, so import cannot verify them.
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateId:                        alias,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "alias",
				ImportStateVerifyIgnore:              []string{"api_key", "app_key"},
			},
			// Update and Read testing: rename in place and change environments.
			{
				Config: testAccDatadogConfigurationResourceConfig(renamedAlias, `["tf-acc", "tf-acc-2"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", renamedAlias),
					resource.TestCheckResourceAttr(resourceName, "alias", renamedAlias),
					resource.TestCheckResourceAttr(resourceName, "environments.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "environments.1", "tf-acc-2"),
					resource.TestCheckResourceAttr(resourceName, "api_key_last_four", "a1b2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDatadogConfigurationResourceConfig(alias string, environments string) string {
	return fmt.Sprintf(`
resource "cortex_datadog_configuration" "test" {
  alias        = %[1]q
  api_key      = %[2]q
  app_key      = %[3]q
  region       = "US1"
  environments = %[4]s
}
`, alias, testDatadogApiKey, testDatadogAppKey, environments)
}
