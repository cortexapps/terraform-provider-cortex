package provider_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntitiesDataSource(t *testing.T) {
	recordName := "data.cortex_catalog_entities.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing - search by type
			{
				Config: testAccCatalogEntitiesDataSourceByType(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(recordName, "entities.#"),
				),
			},
		},
	})
}

func TestAccCatalogEntitiesDataSourceWithQuery(t *testing.T) {
	recordName := "data.cortex_catalog_entities.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing - search with query
			{
				Config: testAccCatalogEntitiesDataSourceWithQuery(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(recordName, "entities.#"),
				),
			},
		},
	})
}

func testAccCatalogEntitiesDataSourceByType() string {
	return `
data "cortex_catalog_entities" "test" {
	types = ["service"]
}`
}

func testAccCatalogEntitiesDataSourceWithQuery() string {
	return `
data "cortex_catalog_entities" "test" {
	query = "test"
	types = ["service"]
}`
}
