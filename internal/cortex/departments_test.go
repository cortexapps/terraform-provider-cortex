package cortex_test

import (
	"context"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDepartmentResponse = &cortex.Department{
	Tag:         "test-department",
	Name:        "Test Department",
	Description: "A test department",
	Members: []cortex.DepartmentMember{
		{
			Name:        "Test User",
			Email:       "test-department-user@cortex.io",
			Description: "A test user",
		},
	},
}

func TestGetDepartment(t *testing.T) {
	testDepartmentTag := "test-department"
	c, teardown, err := setupClient(cortex.Route("departments", ""), testDepartmentResponse, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Departments().Get(context.Background(), testDepartmentTag)
	assert.Nil(t, err, "error retrieving a department")
	assert.Equal(t, testDepartmentResponse, res)
}

func TestCreateDepartment(t *testing.T) {
	tag := "test-department"
	req := cortex.CreateDepartmentRequest{
		Tag:         tag,
		Name:        "Test Department",
		Description: "A test department",
		Members: []cortex.DepartmentMember{
			{
				Name:        "Test User",
				Email:       "test-department-user@cortex.io",
				Description: "A test user",
			},
		},
	}
	c, teardown, err := setupClient(
		cortex.Route("departments", ""),
		testDepartmentResponse,
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Departments().Create(context.Background(), req)
	assert.Nil(t, err, "error creating a department")
	assert.Equal(t, res.Tag, tag)
}

func TestUpdateDepartment(t *testing.T) {
	req := cortex.UpdateDepartmentRequest{
		Name:        "Test Department",
		Description: "A test department",
		Members: []cortex.DepartmentMember{
			{
				Name:        "Test User",
				Email:       "test-department-user@cortex.io",
				Description: "A test user",
			},
		},
	}
	tag := "test-department"

	c, teardown, err := setupClient(
		cortex.Route("departments", tag),
		testDepartmentResponse,
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Departments().Update(context.Background(), tag, req)
	assert.Nil(t, err, "error updating a department")
	assert.Equal(t, res.Tag, tag)
}

func TestDeleteDepartment(t *testing.T) {
	tag := "test-department"

	c, teardown, err := setupClient(
		cortex.Route("departments", ""),
		cortex.DeleteDepartmentResponse{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Departments().Delete(context.Background(), tag)
	assert.Nil(t, err, "error deleting a department")
}
