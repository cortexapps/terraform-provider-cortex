package integrations_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// These tests run against the shared tenant with fake credentials. Cortex does not check credentials when it saves
// a configuration.

const accResourceName = "cortex_integration_configuration.test"

// accIntegrationTest creates, imports, and updates a configuration. updatePlanChecks run on the plan of the update.
func accIntegrationTest(t *testing.T, importId string, create, update string, updatePlanChecks ...plancheck.PlanCheck) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: create, Check: resource.TestCheckResourceAttr(accResourceName, "id", importId)},
			{
				ResourceName:            accResourceName,
				ImportState:             true,
				ImportStateId:           importId,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
			{Config: update, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: updatePlanChecks}},
		},
	})
}

// The update changes the keys, region, custom subdomain, and environments in place.
func TestAccIntegrationConfigurationDatadog(t *testing.T) {
	body := func(keySuffix, settings string) string {
		return fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = "tf-acc-test-datadog"
  credentials = { key_pair = { key = "tf-acc-fake-api-key-%[1]s", secret = "tf-acc-fake-app-key-%[1]s" } }
  datadog     = %[2]s
}`, keySuffix, settings)
	}
	accIntegrationTest(t, "datadog/tf-acc-test-datadog",
		body("a1b2", `{ region = "US1", environments = ["tf-acc"] }`),
		body("e5f6", `{ region = "EU1", custom_subdomain = "tf-acc", environments = ["tf-acc", "tf-acc-2"] }`),
		plancheck.ExpectResourceAction(accResourceName, plancheck.ResourceActionUpdate))
}
