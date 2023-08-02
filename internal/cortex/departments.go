package cortex

import (
	"context"
	"errors"
	"fmt"
	"github.com/dghubble/sling"
)

type DepartmentsClientInterface interface {
	Get(ctx context.Context, tag string) (*Department, error)
	Create(ctx context.Context, req CreateDepartmentRequest) (*Department, error)
	Update(ctx context.Context, tag string, req UpdateDepartmentRequest) (*Department, error)
	Delete(ctx context.Context, tag string) error
}

type DepartmentsClient struct {
	client *HttpClient
}

var _ DepartmentsClientInterface = &DepartmentsClient{}

func (c *DepartmentsClient) Client() *sling.Sling {
	return c.client.client
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// Department is the response from the GET /v1/teams/departments/ endpoint.
type Department struct {
	Tag         string             `json:"departmentTag" url:"departmentTag"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Members     []DepartmentMember `json:"members"`
}

type DepartmentMember struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}

/***********************************************************************************************************************
 * GET /api/v1/teams/departments/?departmentTag
 **********************************************************************************************************************/

type DepartmentGetParams struct {
	DepartmentTag string `json:"departmentTag" url:"departmentTag"`
}

type DepartmentsResponse struct {
	Departments []Department `json:"department"`
}

func (c *DepartmentsClient) Get(ctx context.Context, tag string) (*Department, error) {
	department := &Department{}
	apiError := &ApiError{}

	params := DepartmentGetParams{
		DepartmentTag: tag,
	}
	body, err := c.Client().Get(Route("departments", "")).QueryStruct(params).Receive(department, apiError)
	if err != nil {
		return department, fmt.Errorf("failed getting department: %+v", err)
	}

	err = c.client.handleResponseStatus(body, apiError)
	if err != nil {
		return department, fmt.Errorf("failed getting department: %+v", err)
	}
	return department, nil
}

/***********************************************************************************************************************
 * POST /api/v1/teams/departments
 **********************************************************************************************************************/

type CreateDepartmentRequest struct {
	Tag         string             `json:"departmentTag" url:"departmentTag"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Members     []DepartmentMember `json:"members"`
}

func (c *DepartmentsClient) Create(ctx context.Context, req CreateDepartmentRequest) (*Department, error) {
	department := &Department{}
	apiError := &ApiError{}

	body, err := c.Client().Post(Route("departments", "")).BodyJSON(&req).Receive(department, apiError)
	if err != nil {
		return department, fmt.Errorf("failed creating department: %+v", err)
	}

	err = c.client.handleResponseStatus(body, apiError)
	if err != nil {
		return department, err
	}

	return department, nil
}

/***********************************************************************************************************************
 * PUT /api/v1/teams/departments/:tag
 **********************************************************************************************************************/

type UpdateDepartmentRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Members     []DepartmentMember `json:"members"`
}

func (c *DepartmentsClient) Update(ctx context.Context, tag string, req UpdateDepartmentRequest) (*Department, error) {
	department := &Department{}
	apiError := &ApiError{}

	body, err := c.Client().Put(Route("departments", tag)).BodyJSON(&req).Receive(department, apiError)
	if err != nil {
		return department, errors.New("could not update department: " + err.Error())
	}

	err = c.client.handleResponseStatus(body, apiError)
	if err != nil {
		return department, err
	}

	return department, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/teams/departments/:tag - Delete a department
 **********************************************************************************************************************/

type DeleteDepartmentRequest struct {
	DepartmentTag string `json:"departmentTag" url:"departmentTag"`
}
type DeleteDepartmentResponse struct{}

func (c *DepartmentsClient) Delete(ctx context.Context, tag string) error {
	response := &DeleteDepartmentResponse{}
	apiError := &ApiError{}
	params := DeleteDepartmentRequest{
		DepartmentTag: tag,
	}

	body, err := c.Client().Delete(Route("departments", "")).QueryStruct(params).Receive(response, apiError)
	if err != nil {
		return errors.New("could not delete department: " + err.Error())
	}

	err = c.client.handleResponseStatus(body, apiError)
	if err != nil {
		return err
	}

	return nil
}
