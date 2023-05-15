package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDepartmentResource(t *testing.T) {
	stub := tFactoryBuildDepartmentResource()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDepartmentResourceConfig(stub),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_department.engineering", "tag", stub.Tag),
					resource.TestCheckResourceAttr("cortex_department.engineering", "name", stub.Name),
					resource.TestCheckResourceAttr("cortex_department.engineering", "description", stub.Description),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cortex_department.engineering",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDepartmentResourceConfig(stub),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cortex_department.engineering", "tag", stub.Tag),
					resource.TestCheckResourceAttr("cortex_department.engineering", "name", stub.Name),
					resource.TestCheckResourceAttr("cortex_department.engineering", "description", stub.Description),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDepartmentResourceConfig(stub TestDepartmentResource) string {
	return fmt.Sprintf(`
resource "cortex_department" "engineering" {
  tag = %[1]q
  name = %[2]q
  description = %[3]q
}
`, stub.Tag, stub.Name, stub.Description)
}

type TestDepartmentResource struct {
	Tag         string
	Name        string
	Description string
	Members     []TestDepartmentMemberResource
}
type TestDepartmentMemberResource struct {
	Name        string
	Email       string
	Description string
}

func tFactoryBuildDepartmentResource() TestDepartmentResource {
	return TestDepartmentResource{
		Tag:         "engineering",
		Name:        "Engineering",
		Description: "The Engineering Department",
		Members: []TestDepartmentMemberResource{
			{
				Name:        "John Doe",
				Email:       "test+member1@cortex.io",
				Description: "John Doe is a member of the Engineering Department",
			},
		},
	}
}
