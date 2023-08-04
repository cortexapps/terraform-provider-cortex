package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

type testScorecardDataSource struct {
	Tag  string
	Name string
}

/***********************************************************************************************************************
 * Helper methods
 **********************************************************************************************************************/

func (t *testScorecardDataSource) DataSourceFullName() string {
	return "data." + t.DataSourceType() + "." + t.Tag
}

func (t *testScorecardDataSource) DataSourceType() string {
	return "cortex_scorecard"
}

func (t *testScorecardDataSource) ToTerraform() string {
	return fmt.Sprintf(`
data %[1]q %[2]q {
	tag = %[2]q
}`, t.DataSourceType(), t.Tag)
}

/***********************************************************************************************************************
 * Tests
 **********************************************************************************************************************/

func TestAccScorecardDataSource(t *testing.T) {
	stub := testScorecardDataSource{
		Tag:  "onboarding-scorecard",
		Name: "Manual Onboarding Scorecard",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: stub.ToTerraform(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(stub.DataSourceFullName(), "tag", stub.Tag),
					resource.TestCheckResourceAttr(stub.DataSourceFullName(), "name", stub.Name),
				),
			},
		},
	})
}
