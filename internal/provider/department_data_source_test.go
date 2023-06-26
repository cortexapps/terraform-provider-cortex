package provider_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDepartmentDataSource(t *testing.T) {
	recordName := "data.cortex_department.test-manual-department-root"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDepartmentDataSourceBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(recordName, "tag", "test-manual-department-root"),
					resource.TestCheckResourceAttr(recordName, "name", "Manual Test Department (Root)"),
					resource.TestCheckResourceAttr(recordName, "description", "Department for testing data sources. DO NOT DELETE."),
				),
			},
		},
	})
}

func testAccDepartmentDataSourceBasic() string {
	return `
data "cortex_department" "test-manual-department-root" {
 tag = "test-manual-department-root"
}`
}
