package cortex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDepartmentResponse = &Department{
	Tag:         "test-department",
	Name:        "Test Department",
	Description: "A test department",
	Members: []DepartmentMember{
		{
			Name:        "Test User",
			Email:       "test-department-user@cortex.io",
			Description: "A test user",
		},
	},
}

var testGetDepartmentResponse = &DepartmentsResponse{
	Departments: []Department{
		*testDepartmentResponse,
	},
}

func TestGetDepartment(t *testing.T) {
	testDepartmentTag := "test-department"
	c, teardown, err := setupClient(BaseUris["departments"], testGetDepartmentResponse, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Departments().Get(context.Background(), testDepartmentTag)
	assert.Nil(t, err, "error retrieving a department")
	assert.Equal(t, testDepartmentResponse, res)
}

func TestCreateDepartment(t *testing.T) {
	tag := "test-department"
	req := CreateDepartmentRequest{
		Tag:         tag,
		Name:        "Test Department",
		Description: "A test department",
		Members: []DepartmentMember{
			{
				Name:        "Test User",
				Email:       "test-department-user@cortex.io",
				Description: "A test user",
			},
		},
	}
	c, teardown, err := setupClient(
		BaseUris["departments"],
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
	req := UpdateDepartmentRequest{
		Name:        "Test Department",
		Description: "A test department",
		Members: []DepartmentMember{
			{
				Name:        "Test User",
				Email:       "test-department-user@cortex.io",
				Description: "A test user",
			},
		},
	}
	tag := "test-department"

	c, teardown, err := setupClient(
		BaseUris["departments"],
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
		BaseUris["departments"],
		DeleteDepartmentResponse{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Departments().Delete(context.Background(), tag)
	assert.Nil(t, err, "error deleting a department")
}
