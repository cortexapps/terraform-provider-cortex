package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

type testResourceDefinition struct {
	Type string
}

/***********************************************************************************************************************
 * Helper methods
 **********************************************************************************************************************/

func (t *testResourceDefinition) DataSourceFullName() string {
	return "data." + t.DataSourceType() + "." + t.Type
}

func (t *testResourceDefinition) DataSourceType() string {
	return "cortex_resource_definition"
}

func (t *testResourceDefinition) ToTerraform() string {
	return fmt.Sprintf(`
data %[1]q %[2]q {
	type = %[2]q
}`, t.DataSourceType(), t.Type)
}

/***********************************************************************************************************************
 * Tests
 **********************************************************************************************************************/

func TestAccResourceDefinitionDataSource(t *testing.T) {
	stub := testResourceDefinition{
		Type: "test-resource-definition",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.DataSourceFullName(), "type", stub.Type),
				),
			},
		},
	})
}
