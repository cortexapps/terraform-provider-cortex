package integrations_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	api "github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/cortexapps/terraform-provider-cortex/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// These tests run against the shared tenant with fake credentials. Cortex does not check credentials when it saves
// a configuration. Hosts use the reserved .invalid domain, so background jobs cannot reach a real server.

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

func TestAccIntegrationConfigurationGitlab(t *testing.T) {
	body := func(groups string) string {
		return fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = "tf-acc-test-gitlab"
  credentials = { token = { value = "tf-acc-fake-token-a1b2" } }
  gitlab      = { host = "https://gitlab.invalid", group_names = %s }
}`, groups)
	}
	accIntegrationTest(t, "gitlab/tf-acc-test-gitlab", body(`[]`), body(`["tf-acc"]`))
}

func TestAccIntegrationConfigurationIncidentIo(t *testing.T) {
	body := func(key string) string {
		return fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = "tf-acc-test-incident-io"
  credentials = { token = { value = %q } }
  incident_io = {}
}`, key)
	}
	accIntegrationTest(t, "incident_io/tf-acc-test-incident-io", body("tf-acc-fake-key-a1b2"), body("tf-acc-fake-key-c3d4"))
}

// A tenant has one PagerDuty configuration. The test skips when the tenant already has one, so it never touches a
// real configuration.
func TestAccIntegrationConfigurationPagerDuty(t *testing.T) {
	body := func(readonly bool) string {
		return fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  credentials = { token = { value = "tf-acc-fake-token-a1b2" } }
  pagerduty   = { is_token_readonly = %t }
}`, readonly)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipWhenPagerDutyConfigured(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: body(true), Check: resource.TestCheckResourceAttr(accResourceName, "id", "pagerduty")},
			{ResourceName: accResourceName, ImportState: true, ImportStateId: "pagerduty", ImportStateVerify: true, ImportStateVerifyIgnore: []string{"credentials"}},
			{Config: body(false)},
		},
	})
}

func skipWhenPagerDutyConfigured(t *testing.T) {
	url := os.Getenv("CORTEX_API_URL")
	if url == "" {
		url = provider.DefaultBaseApiUrl
	}
	client, err := cortex.NewClient(cortex.WithURL(url), cortex.WithToken(os.Getenv("CORTEX_API_TOKEN")), cortex.WithVersion("acctest"))
	if err != nil {
		t.Fatalf("could not build client: %s", err)
	}
	_, err = api.PagerDuty(client).Get(context.Background())
	switch {
	case err == nil:
		t.Skip("the tenant already has a PagerDuty configuration")
	case !errors.Is(err, cortex.ApiErrorNotFound):
		t.Fatalf("could not read the PagerDuty configuration: %s", err)
	}
}
