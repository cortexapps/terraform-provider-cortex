package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// These tests run against the shared tenant with fake credentials. Cortex does not check credentials when it saves
// a configuration.

const accResourceName = "cortex_integration_configuration.test"

func accIntegrationTest(t *testing.T, importId string, create, update string, extraIgnore ...string) {
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
				ImportStateVerifyIgnore: append([]string{"credentials"}, extraIgnore...),
			},
			{Config: update},
		},
	})
}

func TestAccIntegrationConfigurationDatadog(t *testing.T) {
	body := func(environments string) string {
		return fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = "tf-acc-test-datadog"
  credentials = { key_pair = { key = "tf-acc-fake-api-key-a1b2", secret = "tf-acc-fake-app-key-c3d4" } }
  datadog     = { region = "US1", environments = %s }
}`, environments)
	}
	accIntegrationTest(t, "datadog/tf-acc-test-datadog", body(`["tf-acc"]`), body(`["tf-acc", "tf-acc-2"]`))
}
