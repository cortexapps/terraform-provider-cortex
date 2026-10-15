package cortex

import (
	"context"
	"errors"

	"github.com/dghubble/sling"
)

type RelationshipTypesClientInterface interface {
	Get(ctx context.Context, tag string) (EntityRelationshipType, error)
	List(ctx context.Context) (RelationshipTypesResponse, error)
	Create(ctx context.Context, req CreateRelationshipTypeRequest) (EntityRelationshipType, error)
	Update(ctx context.Context, tag string, req UpdateRelationshipTypeRequest) (EntityRelationshipType, error)
	Delete(ctx context.Context, tag string) error
}

type RelationshipTypesClient struct {
	client *HttpClient
}

var _ RelationshipTypesClientInterface = &RelationshipTypesClient{}

func (c *RelationshipTypesClient) Client() *sling.Sling {
	return c.client.Client()
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// EntityTypeFilter constrains which entity types can appear at one end of a relationship.
type EntityTypeFilter struct {
	Include   bool     `json:"include"`
	Types     []string `json:"types,omitempty"`
	Providers []string `json:"providers,omitempty"`
}

// RelationshipInheritance defines parameter inheritance for a relationship type.
type RelationshipInheritance struct {
	ParameterTag    string `json:"parameterTag"`
	InheritanceType string `json:"inheritanceType"`
}

// EntityRelationshipType is the API response/model for a relationship type.
type EntityRelationshipType struct {
	Tag                      string                    `json:"tag"`
	Name                     string                    `json:"name"`
	Description              string                    `json:"description,omitempty"`
	AllowCycles              bool                      `json:"allowCycles"`
	IsCortexManaged          bool                      `json:"isCortexManaged"`
	IsSingleSource           bool                      `json:"isSingleSource"`
	IsSingleDestination      bool                      `json:"isSingleDestination"`
	DefinitionLocation       string                    `json:"definitionLocation"`
	SourcesFilter            EntityTypeFilter          `json:"sourcesFilter"`
	DestinationsFilter       EntityTypeFilter          `json:"destinationsFilter"`
	SourceLabelSingular      string                    `json:"sourceLabelSingular,omitempty"`
	SourceLabelPlural        string                    `json:"sourceLabelPlural,omitempty"`
	DestinationLabelSingular string                    `json:"destinationLabelSingular,omitempty"`
	DestinationLabelPlural   string                    `json:"destinationLabelPlural,omitempty"`
	Inheritances             []RelationshipInheritance `json:"inheritances"`
	CreateCatalog            bool                      `json:"createCatalog"`
}

// CreateRelationshipTypeRequest is the POST body for creating a relationship type.
type CreateRelationshipTypeRequest struct {
	Tag                      string                    `json:"tag"`
	Name                     string                    `json:"name"`
	Description              string                    `json:"description,omitempty"`
	AllowCycles              bool                      `json:"allowCycles"`
	IsSingleSource           bool                      `json:"isSingleSource"`
	IsSingleDestination      bool                      `json:"isSingleDestination"`
	DefinitionLocation       string                    `json:"definitionLocation"`
	SourcesFilter            EntityTypeFilter          `json:"sourcesFilter"`
	DestinationsFilter       EntityTypeFilter          `json:"destinationsFilter"`
	SourceLabelSingular      string                    `json:"sourceLabelSingular,omitempty"`
	SourceLabelPlural        string                    `json:"sourceLabelPlural,omitempty"`
	DestinationLabelSingular string                    `json:"destinationLabelSingular,omitempty"`
	DestinationLabelPlural   string                    `json:"destinationLabelPlural,omitempty"`
	Inheritances             []RelationshipInheritance `json:"inheritances"`
	CreateCatalog            bool                      `json:"createCatalog"`
}

// UpdateRelationshipTypeRequest is the PUT body. Tag is required in both path and body.
type UpdateRelationshipTypeRequest struct {
	Tag                      string                    `json:"tag"`
	Name                     string                    `json:"name"`
	Description              string                    `json:"description,omitempty"`
	AllowCycles              bool                      `json:"allowCycles"`
	IsSingleSource           bool                      `json:"isSingleSource"`
	IsSingleDestination      bool                      `json:"isSingleDestination"`
	DefinitionLocation       string                    `json:"definitionLocation"`
	SourcesFilter            EntityTypeFilter          `json:"sourcesFilter"`
	DestinationsFilter       EntityTypeFilter          `json:"destinationsFilter"`
	SourceLabelSingular      string                    `json:"sourceLabelSingular,omitempty"`
	SourceLabelPlural        string                    `json:"sourceLabelPlural,omitempty"`
	DestinationLabelSingular string                    `json:"destinationLabelSingular,omitempty"`
	DestinationLabelPlural   string                    `json:"destinationLabelPlural,omitempty"`
	Inheritances             []RelationshipInheritance `json:"inheritances"`
	CreateCatalog            bool                      `json:"createCatalog"`
}

// RelationshipTypesResponse is the list response.
type RelationshipTypesResponse struct {
	RelationshipTypes []EntityRelationshipType `json:"relationshipTypes"`
}

/***********************************************************************************************************************
 * GET /api/v1/relationship-types/:tag
 **********************************************************************************************************************/

func (c *RelationshipTypesClient) Get(ctx context.Context, tag string) (EntityRelationshipType, error) {
	data := EntityRelationshipType{}
	apiError := ApiError{}

	response, err := c.Client().Get(Route("relationship_types", tag)).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not get relationship type: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, errors.Join(errors.New("Failed getting relationship type: "), err)
	}

	return data, nil
}

/***********************************************************************************************************************
 * GET /api/v1/relationship-types
 **********************************************************************************************************************/

func (c *RelationshipTypesClient) List(ctx context.Context) (RelationshipTypesResponse, error) {
	data := RelationshipTypesResponse{}
	apiError := ApiError{}

	response, err := c.Client().Get(Route("relationship_types", "")).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not get relationship types: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, err
	}

	return data, nil
}

/***********************************************************************************************************************
 * POST /api/v1/relationship-types
 **********************************************************************************************************************/

func (c *RelationshipTypesClient) Create(ctx context.Context, req CreateRelationshipTypeRequest) (EntityRelationshipType, error) {
	data := EntityRelationshipType{}
	apiError := ApiError{}

	response, err := c.Client().Post(Route("relationship_types", "")).BodyJSON(&req).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not create a relationship type: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, err
	}

	return data, nil
}

/***********************************************************************************************************************
 * PUT /api/v1/relationship-types/:tag
 **********************************************************************************************************************/

func (c *RelationshipTypesClient) Update(ctx context.Context, tag string, req UpdateRelationshipTypeRequest) (EntityRelationshipType, error) {
	data := EntityRelationshipType{}
	apiError := ApiError{}

	response, err := c.Client().Put(Route("relationship_types", tag)).BodyJSON(&req).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not update a relationship type: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, err
	}

	return data, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/relationship-types/:tag
 **********************************************************************************************************************/

type DeleteRelationshipTypeResponse struct{}

func (c *RelationshipTypesClient) Delete(ctx context.Context, tag string) error {
	deleteResponse := DeleteRelationshipTypeResponse{}
	apiError := ApiError{}

	response, err := c.Client().Delete(Route("relationship_types", tag)).Receive(&deleteResponse, &apiError)
	if err != nil {
		return errors.New("could not delete relationship type: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return err
	}

	return nil
}
