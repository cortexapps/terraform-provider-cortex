package provider_test

//func TestAccResourceDefinitionResource(t *testing.T) {
//	stub := tFactoryBuildResourceDefinitionResource()
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Create and Read testing
//			{
//				Config: testAccResourceDefinitionResourceConfig(stub),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("cortex_resource_definition.squid_proxy", "type", stub.Type),
//					resource.TestCheckResourceAttr("cortex_resource_definition.squid_proxy", "name", stub.Name),
//					resource.TestCheckResourceAttr("cortex_resource_definition.squid_proxy", "description", stub.Description),
//				),
//			},
//			// ImportState testing
//			{
//				ResourceName:      "cortex_department.engineering",
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//			// Update and Read testing
//			{
//				Config: testAccResourceDefinitionResourceConfig(stub),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("cortex_resource_definition.squid_proxy", "tag", stub.Type),
//					resource.TestCheckResourceAttr("cortex_resource_definition.squid_proxy", "name", stub.Name),
//					resource.TestCheckResourceAttr("cortex_resource_definition.squid_proxy", "description", stub.Description),
//				),
//			},
//			// Delete testing automatically occurs in TestCase
//		},
//	})
//}
//
//func testAccResourceDefinitionResourceConfig(stub TestResourceDefinitionResource) string {
//	return fmt.Sprintf(`
//resource "cortex_resource_definition" "squid_proxy" {
//  type = %[1]q
//  name = %[2]q
//  description = %[3]q
//  schema = {
//	properties = {
//	  tag = {
//		type = "string"
//	  }
//	}
//  }
//}
//`, stub.Type, stub.Name, stub.Description)
//}
//
//type TestResourceDefinitionResource struct {
//	Type        string
//	Name        string
//	Description string
//	Schema      map[string]interface{}
//}
//
//func tFactoryBuildResourceDefinitionResource() TestResourceDefinitionResource {
//	return TestResourceDefinitionResource{
//		Type:        "squid-proxy",
//		Name:        "squid-proxy",
//		Description: "Squid Proxy",
//		Schema: map[string]interface{}{
//			"properties": map[string]interface{}{
//				"tag": map[string]interface{}{
//					"type": "string",
//				},
//			},
//		},
//	}
//}
