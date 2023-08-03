package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccDepartmentMinimalResource(t *testing.T) {
	tag := "test-department-minimal"
	resourceType := "cortex_department"
	resourceName := resourceType + "." + tag
	stub := tFactoryBuildDepartmentMinimalResource(tag)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDepartmentMinimalResourceConfig(resourceType, stub),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", stub.Tag),
					resource.TestCheckResourceAttr(resourceName, "name", stub.Name),
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
				Config: testAccDepartmentMinimalResourceConfig(resourceType, stub),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag", stub.Tag),
					resource.TestCheckResourceAttr(resourceName, "name", stub.Name),
					resource.TestCheckResourceAttr(resourceName, "description", stub.Description),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDepartmentMinimalResourceConfig(resourceType string, stub TestDepartmentResource) string {
	return fmt.Sprintf(`
resource %[1]q %[2]q {
 	tag = %[3]q
 	name = %[4]q
 	description = %[5]q
 	members = [
		{
			name = "John Doe"
			email = "john.doe@cortex.io"
			description = "John Doe is a member of the Engineering Department"
		}
	]
}
`, resourceType, stub.Tag, stub.Tag, stub.Name, stub.Description)
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

func tFactoryBuildDepartmentMinimalResource(tag string) TestDepartmentResource {
	return TestDepartmentResource{
		Tag:         tag,
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
