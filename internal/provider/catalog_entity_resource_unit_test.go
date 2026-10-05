package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func catalogEntityWithStaticAnalysis(url, staticAnalysis string) string {
	return catalogUnitConfig(url, `
resource "cortex_catalog_entity" "test" {
  name = "Unit Test Entity"
  tag  = "unit-test-entity"

  static_analysis = `+staticAnalysis+`
}`)
}

// The API does not store an empty static analysis block, so these blocks are absent on read.
func TestUnitCatalogEntity_EmptyStaticAnalysis(t *testing.T) {
	for name, staticAnalysis := range map[string]string{
		"no integrations": `{}`,
		"empty mend":      `{ mend = {} }`,
		"empty mend ids":  `{ mend = { application_ids = [], project_ids = [] } }`,
	} {
		t.Run(name, func(t *testing.T) {
			_, url := newFakeCatalogApi(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: catalogEntityWithStaticAnalysis(url, staticAnalysis)},
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
