package integrations_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

const unitResourceName = "cortex_integration_configuration.test"

func datadogUnit(url, alias, key, extra string) string {
	return unitConfig(url, fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias = %q
  %s
  credentials = { key_pair = { key = %q, secret = "fake-app-key-c3d4" } }
  datadog     = { region = "US1" }
}`, alias, extra, key))
}

func checkFake(f *fakeCortexApi, seg, alias, field string, want any) resource.TestCheckFunc {
	return func(*terraform.State) error {
		cfg, ok := f.get(seg, alias)
		if !ok {
			return fmt.Errorf("%s configuration %q does not exist", seg, alias)
		}
		if cfg[field] != want {
			return fmt.Errorf("%s configuration %q: %s is %v, want %v", seg, alias, field, cfg[field], want)
		}
		return nil
	}
}

func checkCreates(f *fakeCortexApi, seg string, want int) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if got := f.createCount(seg); got != want {
			return fmt.Errorf("%s create count is %d, want %d", seg, got, want)
		}
		return nil
	}
}

func TestUnitIntegrationConfiguration_DatadogLifecycle(t *testing.T) {
	fake, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: datadogUnit(url, "dd", "fake-api-key-a1b2", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "id", "datadog/dd"),
					resource.TestCheckResourceAttr(unitResourceName, "integration", "datadog"),
					resource.TestCheckResourceAttr(unitResourceName, "is_default", "true"),
					resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.key", "a1b2"),
					resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.secret", "c3d4"),
					resource.TestCheckResourceAttr(unitResourceName, "datadog.environments.#", "0"),
					checkFake(fake, "datadog", "dd", "apiKey", "fake-api-key-a1b2"),
				),
			},
			// Import: the API never returns the credentials.
			{
				ResourceName:            unitResourceName,
				ImportState:             true,
				ImportStateId:           "datadog/dd",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
			// Rename in place.
			{
				Config: datadogUnit(url, "dd-renamed", "fake-api-key-a1b2", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "id", "datadog/dd-renamed"),
					checkCreates(fake, "datadog", 1),
				),
			},
			// A new key updates the configuration in place.
			{
				Config: datadogUnit(url, "dd-renamed", "fake-api-key-e5f6", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.key", "e5f6"),
					checkFake(fake, "datadog", "dd-renamed", "apiKey", "fake-api-key-e5f6"),
					checkCreates(fake, "datadog", 1),
				),
			},
			// A key rotated outside Terraform shows as a change, and the update restores the configured key.
			{
				PreConfig: func() { fake.setField("datadog", "dd-renamed", "apiKey", "rotated-outside-9999") },
				Config:    datadogUnit(url, "dd-renamed", "fake-api-key-e5f6", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFake(fake, "datadog", "dd-renamed", "apiKey", "fake-api-key-e5f6"),
					checkCreates(fake, "datadog", 1),
				),
			},
		},
	})
}

func TestUnitIntegrationConfiguration_BecomesDefaultWhenAnotherDefaultExists(t *testing.T) {
	fake, url := newFakeCortexApi(t)
	fake.seed("datadog", map[string]any{"alias": "existing", "isDefault": true, "region": "US1", "environments": []any{}, "apiKey": "aaaa", "appKey": "bbbb"})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: datadogUnit(url, "dd", "fake-api-key-a1b2", "is_default = true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "is_default", "true"),
					checkFake(fake, "datadog", "existing", "isDefault", false),
				),
			},
			// The default moves back outside Terraform. With is_default not set, Terraform keeps the value from
			// Cortex, and the destroy can delete the configuration because it is no longer the default.
			{
				PreConfig: func() { fake.setDefault("datadog", "existing") },
				Config:    datadogUnit(url, "dd", "fake-api-key-a1b2", ""),
				Check:     resource.TestCheckResourceAttr(unitResourceName, "is_default", "false"),
			},
		},
	})
}

func TestUnitIntegrationConfiguration_RejectsNonDefaultFirstConfiguration(t *testing.T) {
	fake, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      datadogUnit(url, "dd", "fake-api-key-a1b2", "is_default = false"),
			ExpectError: regexp.MustCompile(`must be the default`),
		}},
	})
	assert.Equal(t, 0, fake.createCount("datadog"))
}

func TestUnitIntegrationConfiguration_RejectsUnsettingDefault(t *testing.T) {
	_, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: datadogUnit(url, "dd", "fake-api-key-a1b2", "")},
			{
				Config:      datadogUnit(url, "dd", "fake-api-key-a1b2", "is_default = false"),
				ExpectError: regexp.MustCompile(`Cannot unset the default configuration`),
			},
		},
	})
}

func TestUnitIntegrationConfiguration_UpdateKeepsLiveDefault(t *testing.T) {
	fake, url := newFakeCortexApi(t)
	withEnvironments := func(envs string) string {
		return unitConfig(url, fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = { key_pair = { key = "fake-api-key-a1b2", secret = "fake-app-key-c3d4" } }
  datadog     = { region = "US1", environments = %s }
}`, envs))
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: withEnvironments(`[]`)},
			// With is_default not set, the update sends the live value. Sending false on the default would fail.
			{
				Config: withEnvironments(`["prod"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "is_default", "true"),
					resource.TestCheckResourceAttr(unitResourceName, "datadog.environments.0", "prod"),
					checkCreates(fake, "datadog", 1),
				),
			},
		},
	})
}

func TestUnitIntegrationConfiguration_RecreatesConfigurationDeletedOutsideTerraform(t *testing.T) {
	fake, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: datadogUnit(url, "dd", "fake-api-key-a1b2", "")},
			{
				PreConfig: func() { fake.remove("datadog", "dd") },
				Config:    datadogUnit(url, "dd", "fake-api-key-a1b2", ""),
				Check:     checkCreates(fake, "datadog", 2),
			},
		},
	})
}

func TestUnitIntegrationConfiguration_Validation(t *testing.T) {
	_, url := newFakeCortexApi(t)
	cases := []struct {
		name string
		body string
		want string
	}{
		{"wrong credential kind", `
resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = { token = { value = "x" } }
  datadog     = { region = "US1" }
}`, `datadog accepts: key_pair`},
		{"missing alias", `
resource "cortex_integration_configuration" "test" {
  credentials = { key_pair = { key = "k", secret = "s" } }
  datadog     = { region = "US1" }
}`, `alias is required`},
		{"no integration block", `
resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = { key_pair = { key = "k", secret = "s" } }
}`, `Exactly one of these attributes must be configured`},
		{"two credential kinds", `
resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = { key_pair = { key = "k", secret = "s" }, token = { value = "x" } }
  datadog     = { region = "US1" }
}`, `Invalid Attribute Combination`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps:                    []resource.TestStep{{Config: unitConfig(url, c.body), ExpectError: regexp.MustCompile(c.want)}},
			})
		})
	}
}

func TestUnitIntegrationConfiguration_ImportIdErrors(t *testing.T) {
	_, url := newFakeCortexApi(t)
	for id, want := range map[string]string{
		"datadog":   `datadog/<alias>`,
		"unknown/x": `Unknown integration`,
	} {
		t.Run(id, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config:        datadogUnit(url, "dd", "fake-api-key-a1b2", ""),
					ResourceName:  unitResourceName,
					ImportState:   true,
					ImportStateId: id,
					ExpectError:   regexp.MustCompile(regexp.QuoteMeta(want)),
				}},
			})
		})
	}
}

// A change of region or custom_subdomain updates the configuration in place. The API keeps the current custom
// subdomain when an update omits it, so only a new configuration can remove it.
func TestUnitIntegrationConfiguration_DatadogSettingsUpdate(t *testing.T) {
	fake, url := newFakeCortexApi(t)
	withSettings := func(region, subdomain string) string {
		settings := fmt.Sprintf("{ region = %q }", region)
		if subdomain != "" {
			settings = fmt.Sprintf("{ region = %q, custom_subdomain = %q }", region, subdomain)
		}
		return unitConfig(url, fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = { key_pair = { key = "fake-api-key-a1b2", secret = "fake-app-key-c3d4" } }
  datadog     = %s
}`, settings))
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: withSettings("US1", "acme"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "datadog.custom_subdomain", "acme"),
					checkFake(fake, "datadog", "dd", "customSubdomain", "acme"),
				),
			},
			{
				Config: withSettings("EU1", "acme"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFake(fake, "datadog", "dd", "region", "EU1"),
					checkCreates(fake, "datadog", 1),
				),
			},
			{
				Config: withSettings("EU1", "acme2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFake(fake, "datadog", "dd", "customSubdomain", "acme2"),
					checkCreates(fake, "datadog", 1),
				),
			},
			{
				Config: withSettings("EU1", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(unitResourceName, "datadog.custom_subdomain"),
					checkFake(fake, "datadog", "dd", "customSubdomain", nil),
					checkCreates(fake, "datadog", 2),
				),
			},
		},
	})
}

// A delete that finds no configuration counts as done: another client deleted it after the refresh.
func TestUnitIntegrationConfiguration_DeleteToleratesNotFound(t *testing.T) {
	fake, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: datadogUnit(url, "dd", "fake-api-key-a1b2", "")},
			{
				PreConfig: fake.answerDeletesWithNotFound,
				Config:    datadogUnit(url, "dd", "fake-api-key-a1b2", ""),
				Destroy:   true,
			},
		},
	})
}

// When the follow-up update that makes a new configuration the default fails, the apply fails, and Terraform still
// tracks the created configuration so that it does not become an orphan.
func TestUnitIntegrationConfiguration_CreateKeepsStateWhenDefaultUpdateFails(t *testing.T) {
	fake, url := newFakeCortexApi(t)
	fake.seed("datadog", map[string]any{"alias": "existing", "isDefault": true, "region": "US1", "environments": []any{}, "apiKey": "aaaa", "appKey": "bbbb"})
	fake.failUpdates()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      datadogUnit(url, "dd", "fake-api-key-a1b2", "is_default = true"),
			ExpectError: regexp.MustCompile(`unable to make it the default`),
		}},
	})
	assert.Equal(t, 1, fake.createCount("datadog"))
	// The destroy after the steps deletes the tracked configuration; without state it would stay in Cortex.
	_, orphan := fake.get("datadog", "dd")
	assert.False(t, orphan, "datadog configuration dd is an orphan")
}

// Credentials can come from another resource, so the credentials object or the object of one kind is unknown until
// apply. The plan must accept that, and a later upstream change of the secret must reach Cortex.
func TestUnitIntegrationConfiguration_UnknownCredentials(t *testing.T) {
	for name, credentials := range map[string]string{
		"kind object": `{ key_pair = terraform_data.creds.output }`,
		"credentials": `terraform_data.creds.output.all`,
	} {
		t.Run(name, func(t *testing.T) {
			fake, url := newFakeCortexApi(t)
			withKey := func(key string) string {
				return unitConfig(url, fmt.Sprintf(`
resource "terraform_data" "creds" {
  input = { key = %q, secret = "fake-app-key-c3d4", all = { key_pair = { key = %q, secret = "fake-app-key-c3d4" } } }
}

resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = %s
  datadog     = { region = "US1" }
}`, key, key, credentials))
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: withKey("fake-api-key-a1b2"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.key", "a1b2"),
							checkFake(fake, "datadog", "dd", "apiKey", "fake-api-key-a1b2"),
						),
					},
					{
						Config: withKey("fake-api-key-e5f6"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.key", "e5f6"),
							checkFake(fake, "datadog", "dd", "apiKey", "fake-api-key-e5f6"),
							checkCreates(fake, "datadog", 1),
						),
					},
				},
			})
		})
	}
}

// Settings can come from another resource, so the whole settings object is unknown until apply. For an existing
// configuration the plan must then assume that a field that needs a new configuration changes; else the apply finds
// the change and Terraform fails on an inconsistent final plan. For Datadog, only removing a set custom subdomain
// needs a new configuration.
func TestUnitIntegrationConfiguration_UnknownSettings(t *testing.T) {
	for name, tc := range map[string]struct {
		settings string
		creates  int
	}{
		"custom subdomain set":     {settings: `{ region = %q, environments = ["prod"], custom_subdomain = "acme" }`, creates: 2},
		"custom subdomain not set": {settings: `{ region = %q, environments = ["prod"] }`, creates: 1},
	} {
		t.Run(name, func(t *testing.T) {
			fake, url := newFakeCortexApi(t)
			withRegion := func(region string) string {
				return unitConfig(url, fmt.Sprintf(`
resource "terraform_data" "settings" {
  input = %s
}

resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = { key_pair = { key = "fake-api-key-a1b2", secret = "fake-app-key-c3d4" } }
  datadog     = terraform_data.settings.output
}`, fmt.Sprintf(tc.settings, region)))
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: withRegion("US1"),
						Check:  checkFake(fake, "datadog", "dd", "region", "US1"),
					},
					{
						Config: withRegion("EU1"),
						Check: resource.ComposeAggregateTestCheckFunc(
							checkFake(fake, "datadog", "dd", "region", "EU1"),
							checkCreates(fake, "datadog", tc.creates),
						),
					},
				},
			})
		})
	}
}

// Only removing a set custom subdomain needs a new configuration, so an unknown custom_subdomain replaces a
// configuration that has one, and updates a configuration without one in place.
func TestUnitIntegrationConfiguration_DatadogUnknownCustomSubdomain(t *testing.T) {
	for name, tc := range map[string]struct {
		before  string
		creates int
	}{
		"was set":     {before: `custom_subdomain = "acme"`, creates: 2},
		"was not set": {before: ``, creates: 1},
	} {
		t.Run(name, func(t *testing.T) {
			fake, url := newFakeCortexApi(t)
			// A new input makes the output of terraform_data unknown until apply.
			withSettings := func(input, settings string) string {
				return unitConfig(url, fmt.Sprintf(`
resource "terraform_data" "subdomain" {
  input = %q
}

resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = { key_pair = { key = "fake-api-key-a1b2", secret = "fake-app-key-c3d4" } }
  datadog     = { region = "US1", %s }
}`, input, settings))
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: withSettings("acme1", tc.before)},
					{
						Config: withSettings("acme2", "custom_subdomain = terraform_data.subdomain.output"),
						Check: resource.ComposeAggregateTestCheckFunc(
							checkFake(fake, "datadog", "dd", "customSubdomain", "acme2"),
							checkCreates(fake, "datadog", tc.creates),
						),
					},
				},
			})
		})
	}
}

func pagerDutyUnit(url, token, extra string) string {
	return unitConfig(url, fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  %s
  credentials = { token = { value = %q } }
  pagerduty   = { is_token_readonly = true }
}`, extra, token))
}

func TestUnitIntegrationConfiguration_PagerDutyLifecycle(t *testing.T) {
	fake, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: pagerDutyUnit(url, "fake-token-a1b2", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "id", "pagerduty"),
					resource.TestCheckNoResourceAttr(unitResourceName, "is_default"),
					resource.TestCheckNoResourceAttr(unitResourceName, "alias"),
					resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.value", "a1b2"),
					checkFake(fake, "pagerduty", "", "isTokenReadonly", true),
				),
			},
			{
				ResourceName:            unitResourceName,
				ImportState:             true,
				ImportStateId:           "pagerduty",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
			// A new token updates in place through the full-replace PUT.
			{
				Config: pagerDutyUnit(url, "fake-token-c3d4", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.value", "c3d4"),
					checkFake(fake, "pagerduty", "", "token", "fake-token-c3d4"),
					checkCreates(fake, "pagerduty", 1),
				),
			},
			// A token rotated outside Terraform shows as a change, and the apply restores the configured token.
			{
				PreConfig: func() { fake.setField("pagerduty", "", "token", "rotated-outside-9999") },
				Config:    pagerDutyUnit(url, "fake-token-c3d4", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFake(fake, "pagerduty", "", "token", "fake-token-c3d4"),
					checkCreates(fake, "pagerduty", 1),
				),
			},
			// Deleted outside Terraform: Read removes it, and the apply creates it again.
			{
				PreConfig: func() { fake.remove("pagerduty", "") },
				Config:    pagerDutyUnit(url, "fake-token-c3d4", ""),
				Check:     checkCreates(fake, "pagerduty", 2),
			},
		},
	})
}

func TestUnitIntegrationConfiguration_PagerDutyCreateFailsWhenConfigurationExists(t *testing.T) {
	fake, url := newFakeCortexApi(t)
	fake.seed("pagerduty", map[string]any{"token": "real-token-9999", "isTokenReadonly": false})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      pagerDutyUnit(url, "fake-token-a1b2", ""),
			ExpectError: regexp.MustCompile(`Import it instead`),
		}},
	})
	cfg, _ := fake.get("pagerduty", "")
	assert.Equal(t, "real-token-9999", cfg["token"])
	assert.Equal(t, 0, fake.createCount("pagerduty"))
}

func TestUnitIntegrationConfiguration_PagerDutyValidation(t *testing.T) {
	_, url := newFakeCortexApi(t)
	for name, c := range map[string]struct{ extra, want string }{
		"alias":      {`alias = "x"`, `alias is not allowed`},
		"is_default": {`is_default = true`, `is_default is not allowed`},
	} {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps:                    []resource.TestStep{{Config: pagerDutyUnit(url, "t", c.extra), ExpectError: regexp.MustCompile(c.want)}},
			})
		})
	}
}

func TestUnitIntegrationConfiguration_PagerDutyImportIdError(t *testing.T) {
	_, url := newFakeCortexApi(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:        pagerDutyUnit(url, "t", ""),
			ResourceName:  unitResourceName,
			ImportState:   true,
			ImportStateId: "pagerduty/x",
			ExpectError:   regexp.MustCompile(`The import ID for PagerDuty is pagerduty\.`),
		}},
	})
}

func incidentIoUnit(url, alias, key string) string {
	return unitConfig(url, fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = %q
  credentials = { token = { value = %q } }
  incident_io = {}
}`, alias, key))
}

func TestUnitIntegrationConfiguration_IncidentIoLifecycle(t *testing.T) {
	fake, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: incidentIoUnit(url, "inc", "fake-key-a1b2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "id", "incident_io/inc"),
					resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.value", "a1b2"),
				),
			},
			{
				ResourceName:            unitResourceName,
				ImportState:             true,
				ImportStateId:           "incident_io/inc",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
			// A new key updates in place.
			{
				Config: incidentIoUnit(url, "inc", "fake-key-c3d4"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFake(fake, "incidentio", "inc", "apiKey", "fake-key-c3d4"),
					checkCreates(fake, "incidentio", 1),
				),
			},
		},
	})
}

// Changing the settings block to another integration replaces the configuration.
func TestUnitIntegrationConfiguration_ChangeOfIntegrationReplaces(t *testing.T) {
	fake, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: incidentIoUnit(url, "shared", "fake-key-a1b2")},
			{
				Config: datadogUnit(url, "shared", "fake-api-key-a1b2", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "integration", "datadog"),
					checkCreates(fake, "datadog", 1),
					func(*terraform.State) error {
						if _, ok := fake.get("incidentio", "shared"); ok {
							return fmt.Errorf("incident.io configuration still exists")
						}
						return nil
					},
				),
			},
		},
	})
}

func gitlabUnit(url, token, host, groups string) string {
	return unitConfig(url, fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = "gl"
  credentials = { token = { value = %q } }
  gitlab      = { host = %q, group_names = %s }
}`, token, host, groups))
}

func TestUnitIntegrationConfiguration_GitlabLifecycle(t *testing.T) {
	fake, url := newFakeCortexApi(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: gitlabUnit(url, "fake-token-a1b2", "https://gitlab.invalid", `[]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "gitlab.host", "https://gitlab.invalid"),
					resource.TestCheckResourceAttr(unitResourceName, "gitlab.hide_personal_projects", "false"),
					checkFake(fake, "gitlab", "gl", "personalAccessToken", "fake-token-a1b2"),
				),
			},
			// A new token and new groups update in place.
			{
				Config: gitlabUnit(url, "fake-token-c3d4", "https://gitlab.invalid", `["platform"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "gitlab.group_names.0", "platform"),
					checkFake(fake, "gitlab", "gl", "personalAccessToken", "fake-token-c3d4"),
					checkCreates(fake, "gitlab", 1),
				),
			},
			// The API ignores host on update, so a new host replaces the configuration.
			{
				Config: gitlabUnit(url, "fake-token-c3d4", "https://gitlab2.invalid", `["platform"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFake(fake, "gitlab", "gl", "host", "https://gitlab2.invalid"),
					checkCreates(fake, "gitlab", 2),
				),
			},
		},
	})
}

// The API drops blank group names, which would make the apply inconsistent, so validation rejects them.
func TestUnitIntegrationConfiguration_GitlabRejectsBlankGroupNames(t *testing.T) {
	_, url := newFakeCortexApi(t)
	for name, groups := range map[string]string{"empty": `[""]`, "blank": `[" "]`} {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config:      gitlabUnit(url, "fake-token-a1b2", "https://gitlab.invalid", groups),
					ExpectError: regexp.MustCompile(`group_names`),
				}},
			})
		})
	}
}
