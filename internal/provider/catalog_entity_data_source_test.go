package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogEntityDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCatalogEntityDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cortex_catalog_entity.test", "tag", "test"),
				),
			},
		},
	})
}

const testAccCatalogEntityDataSourceConfig = `
data "cortex_catalog_entity" "test" {
  tag = "test"
}
`
