package provider_test

//func TestAccResourceDefinitionDataSource(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Read testing
//			{
//				Config: testAccResourceDefinitionDataSourceConfig,
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("data.cortex_resource_definition.squid_proxy", "type", "squid-proxy"),
//				),
//			},
//		},
//	})
//}
//
//const testAccResourceDefinitionDataSourceConfig = `
//data "cortex_department" "engineering" {
//  type = "squid-proxy"
//}
//`
