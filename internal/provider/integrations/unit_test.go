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

//nolint:unparam // The next integrations in the stack use other segments.
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

//nolint:unparam // The next integrations in the stack use other segments.
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
			// A new key replaces the configuration, because the API cannot update Datadog keys.
			{
				Config: datadogUnit(url, "dd-renamed", "fake-api-key-e5f6", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(unitResourceName, "credentials_last_four.key", "e5f6"),
					checkCreates(fake, "datadog", 2),
				),
			},
			// A key rotated outside Terraform shows as a change and replaces the configuration again.
			{
				PreConfig: func() { fake.setField("datadog", "dd-renamed", "apiKey", "rotated-outside-9999") },
				Config:    datadogUnit(url, "dd-renamed", "fake-api-key-e5f6", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFake(fake, "datadog", "dd-renamed", "apiKey", "fake-api-key-e5f6"),
					checkCreates(fake, "datadog", 3),
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

// The API cannot change region or custom_subdomain in place, so a change to one of them replaces the configuration.
func TestUnitIntegrationConfiguration_DatadogSettingsReplace(t *testing.T) {
	fake, url := newFakeCortexApi(t)
	withSettings := func(region, subdomain string) string {
		return unitConfig(url, fmt.Sprintf(`
resource "cortex_integration_configuration" "test" {
  alias       = "dd"
  credentials = { key_pair = { key = "fake-api-key-a1b2", secret = "fake-app-key-c3d4" } }
  datadog     = { region = %q, custom_subdomain = %q }
}`, region, subdomain))
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
					checkCreates(fake, "datadog", 2),
				),
			},
			{
				Config: withSettings("EU1", "acme2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFake(fake, "datadog", "dd", "customSubdomain", "acme2"),
					checkCreates(fake, "datadog", 3),
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
