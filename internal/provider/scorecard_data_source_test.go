package provider_test

//func TestAccScorecardDataSource(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Read testing
//			{
//				Config: testAccScorecardDataSourceConfig,
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("data.cortex_scorecard.dora-metrics", "tag", "dora"),
//				),
//			},
//		},
//	})
//}
//
//const testAccScorecardDataSourceConfig = `
//data "cortex_scorecard" "dora-metrics" {
//  tag = "dora"
//}
//`
