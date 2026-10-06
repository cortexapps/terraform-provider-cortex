package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func catalogEntityWithStaticAnalysis(url, staticAnalysis string) string {
	return catalogUnitConfig(url, `
resource "cortex_catalog_entity" "test" {
  name = "Unit Test Entity"
  tag  = "unit-test-entity"

  static_analysis = `+staticAnalysis+`
}`)
}

// The provider omits an empty static analysis block on write, so these blocks are absent on read.
func TestUnitCatalogEntity_EmptyStaticAnalysis(t *testing.T) {
	for name, staticAnalysis := range map[string]string{
		"no integrations": `{}`,
		"empty mend":      `{ mend = {} }`,
		"empty mend ids":  `{ mend = { application_ids = [], project_ids = [] } }`,
	} {
		t.Run(name, func(t *testing.T) {
			fake, url := newFakeCatalogApi(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: catalogEntityWithStaticAnalysis(url, staticAnalysis),
						Check:  checkStoredInfoAbsent(fake, "unit-test-entity", "x-cortex-static-analysis"),
					},
				},
			})
		})
	}
}

func TestUnitCatalogEntity_MendWithOneIdList(t *testing.T) {
	_, url := newFakeCatalogApi(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: catalogEntityWithStaticAnalysis(url, `{ mend = { application_ids = ["123"] } }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.application_ids.0", "123"),
					resource.TestCheckNoResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.project_ids"),
				),
			},
			{
				ResourceName:                         "cortex_catalog_entity.test",
				ImportState:                          true,
				ImportStateId:                        "unit-test-entity",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "tag",
			},
		},
	})
}

func TestUnitCatalogEntity_MendWithBothIdLists(t *testing.T) {
	_, url := newFakeCatalogApi(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: catalogEntityWithStaticAnalysis(url, `{ mend = { application_ids = ["123"], project_ids = ["456"] } }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.application_ids.0", "123"),
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.project_ids.0", "456"),
				),
			},
		},
	})
}

func checkStoredInfoAbsent(f *fakeCatalogApi, tag, key string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		descriptor, ok := f.descriptors[tag]
		if !ok {
			return fmt.Errorf("entity %q does not exist", tag)
		}
		info, _ := descriptor["info"].(map[string]any)
		if _, ok := info[key]; ok {
			return fmt.Errorf("entity %q stores %s: %v", tag, key, info[key])
		}
		return nil
	}
}

// Each integration reads as null when the descriptor has none, so a block with one integration imports cleanly.
func TestUnitCatalogEntity_StaticAnalysisWithOneIntegration(t *testing.T) {
	for name, staticAnalysis := range map[string]string{
		"sonar qube only":         `{ sonar_qube = { project = "p" } }`,
		"veracode names only":     `{ veracode = { application_names = ["app"] } }`,
		"veracode sandboxes only": `{ veracode = { sandboxes = [{ application_name = "app", sandbox_name = "staging" }] } }`,
	} {
		t.Run(name, func(t *testing.T) {
			_, url := newFakeCatalogApi(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: catalogEntityWithStaticAnalysis(url, staticAnalysis)},
					{
						ResourceName:                         "cortex_catalog_entity.test",
						ImportState:                          true,
						ImportStateId:                        "unit-test-entity",
						ImportStateVerify:                    true,
						ImportStateVerifyIdentifierAttribute: "tag",
					},
				},
			})
		})
	}
}

func TestUnitCatalogEntity_UpdateMendToEmptyAndBack(t *testing.T) {
	_, url := newFakeCatalogApi(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: catalogEntityWithStaticAnalysis(url, `{ mend = { application_ids = ["1"], project_ids = ["2"] } }`)},
			{
				Config: catalogEntityWithStaticAnalysis(url, `{ mend = {} }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("cortex_catalog_entity.test", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: catalogEntityWithStaticAnalysis(url, `{ mend = { project_ids = ["3"] } }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.project_ids.0", "3"),
					resource.TestCheckNoResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.application_ids"),
				),
			},
		},
	})
}

// Keeping a configured empty block must not hide a change made outside Terraform.
func TestUnitCatalogEntity_EmptyStaticAnalysisShowsRemoteChange(t *testing.T) {
	fake, url := newFakeCatalogApi(t)
	config := catalogEntityWithStaticAnalysis(url, `{ mend = {} }`)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					fake.setInfo("unit-test-entity", "x-cortex-static-analysis", map[string]any{"mend": map[string]any{"applicationIds": []any{"remote"}}})
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check:              resource.TestCheckResourceAttr("cortex_catalog_entity.test", "static_analysis.mend.application_ids.0", "remote"),
			},
		},
	})
}

// ImportStateVerify treats an empty list as absent, so plan after the import to catch an empty list in state.
func TestUnitCatalogEntity_ImportThenPlanIsEmpty(t *testing.T) {
	for name, tc := range map[string]struct{ descriptor, staticAnalysis string }{
		"mend one id list":        {`mend: {applicationIds: ["123"]}`, `{ mend = { application_ids = ["123"] } }`},
		"sonar qube only":         {`sonarqube: {project: p}`, `{ sonar_qube = { project = "p" } }`},
		"veracode names only":     {`veracode: {applicationNames: [app]}`, `{ veracode = { application_names = ["app"] } }`},
		"veracode sandboxes only": {`veracode: {sandboxes: [{applicationName: app, sandboxName: staging}]}`, `{ veracode = { sandboxes = [{ application_name = "app", sandbox_name = "staging" }] } }`},
	} {
		t.Run(name, func(t *testing.T) {
			fake, url := newFakeCatalogApi(t)
			fake.seed(t, `
info:
  title: Unit Test Entity
  x-cortex-tag: unit-test-entity
  x-cortex-type: service
  x-cortex-static-analysis: {`+tc.descriptor+`}
`)
			config := catalogEntityWithStaticAnalysis(url, tc.staticAnalysis)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:             config,
						ResourceName:       "cortex_catalog_entity.test",
						ImportState:        true,
						ImportStateId:      "unit-test-entity",
						ImportStatePersist: true,
					},
					{Config: config, PlanOnly: true},
				},
			})
		})
	}
}
