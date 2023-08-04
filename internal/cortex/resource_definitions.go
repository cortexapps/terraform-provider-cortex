package cortex

import (
	"context"
	"errors"
	"github.com/dghubble/sling"
)

type ResourceDefinitionsClientInterface interface {
	Get(ctx context.Context, typeName string) (ResourceDefinition, error)
	List(ctx context.Context, params *ResourceDefinitionListParams) (ResourceDefinitionsResponse, error)
	Create(ctx context.Context, req CreateResourceDefinitionRequest) (ResourceDefinition, error)
	Update(ctx context.Context, typeName string, req UpdateResourceDefinitionRequest) (ResourceDefinition, error)
	Delete(ctx context.Context, typeName string) error
}

type ResourceDefinitionsClient struct {
	client *HttpClient
}

var _ ResourceDefinitionsClientInterface = &ResourceDefinitionsClient{}

func (c *ResourceDefinitionsClient) Client() *sling.Sling {
	return c.client.Client()
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// ResourceDefinition is the response object that is typically returned from the resource definitions endpoints.
type ResourceDefinition struct {
	Type        string                 `json:"type" yaml:"type"`
	Name        string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty" yaml:"schema,omitempty"`
	Source      string                 `json:"source,omitempty" yaml:"source,omitempty"`
}

func (r *ResourceDefinition) SchemaAsString() (string, error) {
	value := ""
	if r.Schema != nil {
		err := error(nil)
		value, err = InterfaceToString(r.Schema)
		if err != nil {
			return "", err
		}
	}
	return value, nil
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/definitions/:typeName
 **********************************************************************************************************************/

func (c *ResourceDefinitionsClient) Get(ctx context.Context, typeName string) (ResourceDefinition, error) {
	data := ResourceDefinition{}
	apiError := ApiError{}
	response, err := c.Client().Get(Route("resource_definitions", typeName)).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not get resource definition: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, errors.Join(errors.New("Failed getting resource definition: "), err)
	}

	return data, nil
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/definitions
 **********************************************************************************************************************/

// ResourceDefinitionListParams are the query parameters for the GET /v1/catalog/definitions endpoint.
type ResourceDefinitionListParams struct {
	IncludeBuiltIn bool `url:"includeBuiltIn,omitempty"`
}

// ResourceDefinitionsResponse is the response from the GET /v1/catalog/definitions endpoint.
type ResourceDefinitionsResponse struct {
	ResourceDefinitions []ResourceDefinition `json:"definitions"`
}

// List retrieves a list of resource definitions based on a query.
func (c *ResourceDefinitionsClient) List(ctx context.Context, params *ResourceDefinitionListParams) (ResourceDefinitionsResponse, error) {
	data := ResourceDefinitionsResponse{}
	apiError := ApiError{}

	response, err := c.Client().Get(Route("resource_definitions", "")).QueryStruct(&params).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not get resource definitions: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, err
	}

	return data, nil
}

/***********************************************************************************************************************
 * POST /api/v1/catalog/definitions
 **********************************************************************************************************************/

type CreateResourceDefinitionRequest struct {
	Type        string                 `json:"type"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
	Source      string                 `json:"source,omitempty"`
}

func (r *ResourceDefinition) ToCreateRequest() CreateResourceDefinitionRequest {
	return CreateResourceDefinitionRequest{
		Type:        r.Type,
		Name:        r.Name,
		Description: r.Description,
		Schema:      r.Schema,
	}
}

type CreateResourceDefinitionResponse struct {
	ResourceDefinition *ResourceDefinition `json:"scorecard"`
}

func (c *ResourceDefinitionsClient) Create(ctx context.Context, req CreateResourceDefinitionRequest) (ResourceDefinition, error) {
	data := ResourceDefinition{}
	apiError := ApiError{}

	response, err := c.Client().Post(Route("resource_definitions", "")).BodyJSON(&req).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not create a resource definition: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, err
	}

	return data, nil
}

/***********************************************************************************************************************
 * PUT /api/v1/catalog/definitions/:typeName
 **********************************************************************************************************************/

type UpdateResourceDefinitionRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

func (r *ResourceDefinition) ToUpdateRequest() UpdateResourceDefinitionRequest {
	return UpdateResourceDefinitionRequest{
		Name:        r.Name,
		Description: r.Description,
		Schema:      r.Schema,
	}
}

func (c *ResourceDefinitionsClient) Update(ctx context.Context, typeName string, req UpdateResourceDefinitionRequest) (ResourceDefinition, error) {
	data := ResourceDefinition{}
	apiError := ApiError{}

	response, err := c.Client().Put(Route("resource_definitions", typeName)).BodyJSON(&req).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not update a resource definition: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, err
	}

	return data, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/catalog/definitions/:typeName - Delete a resource definition
 **********************************************************************************************************************/

type DeleteResourceDefinitionResponse struct{}

func (c *ResourceDefinitionsClient) Delete(ctx context.Context, typeName string) error {
	scorecardResponse := DeleteResourceDefinitionResponse{}
	apiError := ApiError{}

	response, err := c.Client().Delete(Route("resource_definitions", typeName)).Receive(&scorecardResponse, &apiError)
	if err != nil {
		return errors.New("could not delete scorecard: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return err
	}

	return nil
}
