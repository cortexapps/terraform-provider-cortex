package provider_test

import (
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCatalogEntityCustomDataMinimalResource(t *testing.T) {
	entityTag := "manual-test"
	key := "test-custom-data-key"
	resourceTag := entityTag + "-" + key
	resourceType := "cortex_catalog_entity_custom_data"
	resourceName := resourceType + "." + resourceTag
	stub := tFactoryBuildCatalogEntityCustomDataResource(
		entityTag,
		key,
		"A test custom data key",
		"test-custom-data-value",
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", stub.Tag),
					resource.TestCheckResourceAttr(resourceName, "key", stub.Key),
					resource.TestCheckResourceAttr(resourceName, "description", stub.Description),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", stub.Tag),
					resource.TestCheckResourceAttr(resourceName, "key", stub.Key),
					resource.TestCheckResourceAttr(resourceName, "description", stub.Description),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

type TestCatalogEntityCustomDataResource struct {
	Tag         string
	Key         string
	Description string
	Value       interface{}
}

func (t *TestCatalogEntityCustomDataResource) ToTerraform() string {
	value, _ := t.ValueAsString()
	return fmt.Sprintf(`
resource "cortex_catalog_entity_custom_data" "%[1]s-%[2]s" {
 	tag = %[1]q
 	key = %[2]q
 	description = %[3]q
    value = %[4]q
}
`, t.Tag, t.Key, t.Description, value)
}

func (t *TestCatalogEntityCustomDataResource) ValueAsString() (string, error) {
	value := ""
	if t.Value != nil {
		err := error(nil)
		value, err = cortex.InterfaceToString(t.Value)
		if err != nil {
			return "", err
		}
	}
	return value, nil
}

func tFactoryBuildCatalogEntityCustomDataResource(tag string, key string, description string, value interface{}) TestCatalogEntityCustomDataResource {
	r := TestCatalogEntityCustomDataResource{
		Tag:         tag,
		Key:         key,
		Description: description,
		Value:       value,
	}
	return r
}

func TestAccCatalogEntityCustomDataCompleteResource(t *testing.T) {
	entityTag := "manual-test"
	key := "test-custom-data-complete"
	resourceTag := entityTag + "-" + key
	resourceType := "cortex_catalog_entity_custom_data"
	resourceName := resourceType + "." + resourceTag
	stub := tFactoryBuildCatalogEntityCustomDataResource(
		entityTag,
		key,
		"A test custom data key",
		map[string]interface{}{
			"test": "test",
			"a": map[string]interface{}{
				"b": "c",
			},
		},
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", stub.Tag),
					resource.TestCheckResourceAttr(resourceName, "key", stub.Key),
					resource.TestCheckResourceAttr(resourceName, "description", stub.Description),
					resource.TestCheckResourceAttr(resourceName, "value", `{"a":{"b":"c"},"test":"test"}`),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", stub.Tag),
					resource.TestCheckResourceAttr(resourceName, "key", stub.Key),
					resource.TestCheckResourceAttr(resourceName, "description", stub.Description),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
