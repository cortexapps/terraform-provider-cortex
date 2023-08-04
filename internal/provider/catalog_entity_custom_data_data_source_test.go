package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

type testCatalogEntityCustomDataManual struct {
	Tag         string
	Key         string
	Description string
	Source      string
	Value       interface{}
}

/***********************************************************************************************************************
 * Helper methods
 **********************************************************************************************************************/

func (t *testCatalogEntityCustomDataManual) ResourceTag() string {
	return t.Tag + "-" + t.Key
}

func (t *testCatalogEntityCustomDataManual) ResourceName() string {
	return "data." + t.ResourceType() + "." + t.ResourceTag()
}

func (t *testCatalogEntityCustomDataManual) ResourceType() string {
	return "cortex_catalog_entity_custom_data"
}

func (t *testCatalogEntityCustomDataManual) ValueAsString() string {
	return ""
}

func (t *testCatalogEntityCustomDataManual) ToTerraform() string {
	return fmt.Sprintf(`
data %[1]q %[2]q {
	tag = %[3]q
    key = %[4]q
}`, t.ResourceType(), t.ResourceTag(), t.Tag, t.Key)
}

/***********************************************************************************************************************
 * Tests
 **********************************************************************************************************************/

/* manual-test - manual attribute */

func TestAccCatalogEntityCustomDataDataSource(t *testing.T) {
	stub := testCatalogEntityCustomDataManual{
		Tag:         "manual-test",
		Key:         "manual",
		Description: "",
		Source:      "YAML",
		Value: map[string]interface{}{
			"test": "one",
			"things": []string{
				"two",
				"three",
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.ResourceName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.ResourceName(), "key", stub.Key),
					resource.TestCheckResourceAttr(stub.ResourceName(), "description", stub.Description),
					resource.TestCheckResourceAttr(stub.ResourceName(), "source", stub.Source),
				),
			},
		},
	})
}
