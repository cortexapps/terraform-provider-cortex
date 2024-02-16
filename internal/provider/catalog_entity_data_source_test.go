package provider_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntityDataSource(t *testing.T) {
	recordName := "data.cortex_catalog_entity.manual-test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogEntityDataSourceBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(recordName, "tag", "manual-test"),
					resource.TestCheckResourceAttr(recordName, "name", "Manual Test Service"),
					resource.TestCheckResourceAttr(recordName, "description", "A manual service for data source testing. DO NOT DELETE."),
				),
			},
		},
	})
}

func TestAccCatalogDomainDataSource(t *testing.T) {
	recordName := "data.cortex_catalog_entity.domain-manual-test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogDomainEntityDataSourceBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(recordName, "tag", "domain-manual-test"),
					resource.TestCheckResourceAttr(recordName, "name", "Manual Test Domain"),
					resource.TestCheckResourceAttr(recordName, "description", "A manual domain for data source testing. DO NOT DELETE."),
					resource.TestCheckResourceAttr(recordName, "children.#", "1"),
					resource.TestCheckResourceAttr(recordName, "children.0.tag", "manual-test"),
				),
			},
		},
	})
}

func testAccCatalogEntityDataSourceBasic() string {
	return `
data "cortex_catalog_entity" "manual-test" {
	tag = "manual-test"
}`
}

func testAccCatalogDomainEntityDataSourceBasic() string {
	return `
data "cortex_catalog_entity" "domain-manual-test" {
	tag = "domain-manual-test"
}`
}
