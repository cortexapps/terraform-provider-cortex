package cortex

import (
	"context"
	"errors"
	"github.com/dghubble/sling"
)

type CatalogEntitiesClientInterface interface {
	Get(ctx context.Context, tag string) (*CatalogEntity, error)
	List(ctx context.Context, params *CatalogEntityListParams) (*CatalogEntitiesResponse, error)
	Delete(ctx context.Context, tag string) error
}

type CatalogEntitiesClient struct {
	client *HttpClient
}

var _ CatalogEntitiesClientInterface = &CatalogEntitiesClient{}

func (c *CatalogEntitiesClient) Client() *sling.Sling {
	return c.client.client
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

type CatalogEntity struct {
	Tag         string                 `json:"tag" yaml:"x-cortex-tag"`
	Description string                 `json:"description" yaml:"description"`
	Groups      []string               `json:"groups" yaml:"x-cortex-groups"`
	Links       []CatalogEntityLink    `json:"links" yaml:"x-cortex-link"`
	Metadata    map[string]interface{} `json:"metadata" yaml:"x-cortex-custom-metadata"`
	Ownership   CatalogEntityOwnership `json:"ownership" yaml:"x-cortex-owners"`
}

type CatalogEntityOwnership struct {
	Emails        []CatalogEntityEmail      `json:"emails"`
	Groups        []CatalogEntityGroup      `json:"groups"`
	SlackChannels CatalogEntitySlackChannel `json:"slackChannels"`
}

type CatalogEntityEmail struct {
	Email       string `json:"email" yaml:"email"`
	Description string `json:"description" yaml:"description"`
}

type CatalogEntityGroup struct {
	GroupName   string `json:"group" yaml:"groupName"`
	Description string `json:"description" yaml:"description"`
	Provider    string `json:"provider" yaml:"provider"`
}

type CatalogEntitySlackChannel struct {
	Channel              string `json:"channel" yaml:"channel"`
	Description          string `json:"description" yaml:"description"`
	NotificationsEnabled bool   `json:"notificationsEnabled" yaml:"notificationsEnabled"`
}

type CatalogEntityLink struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
	Url  string `json:"url" yaml:"url"`
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/:tag
 **********************************************************************************************************************/

func (c *CatalogEntitiesClient) Get(ctx context.Context, tag string) (*CatalogEntity, error) {
	catalogEntityResponse := &CatalogEntity{}
	apiError := &ApiError{}
	response, err := c.Client().Get(BaseUris["catalog_entities"]+tag).Receive(catalogEntityResponse, apiError)
	if err != nil {
		return catalogEntityResponse, errors.New("could not get catalog entity: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return catalogEntityResponse, errors.Join(errors.New("Failed getting catalog entity: "), err)
	}

	return catalogEntityResponse, nil
}

/***********************************************************************************************************************
 * GET /api/v1/catalog
 **********************************************************************************************************************/

// CatalogEntityListParams are the query parameters for the GET /v1/catalog endpoint.
type CatalogEntityListParams struct {
	Groups          []string `url:"groups,omitempty"`
	Types           []string `url:"types,omitempty"`
	GitRepositories []string `url:"gitRepositories,omitempty"`
}

// CatalogEntitiesResponse is the response from the GET /v1/scorecards endpoint.
type CatalogEntitiesResponse struct {
	Entities []CatalogEntity `json:"entities" yaml:"entities"`
}

// List retrieves a list of scorecards based on a query.
func (c *CatalogEntitiesClient) List(ctx context.Context, params *CatalogEntityListParams) (*CatalogEntitiesResponse, error) {
	entitiesResponse := &CatalogEntitiesResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Get(BaseUris["catalog_entities"]).QueryStruct(params).Receive(entitiesResponse, apiError)
	if err != nil {
		return nil, errors.New("could not get entities: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return nil, err
	}

	return entitiesResponse, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/catalog/:tag - Delete a catalog entity
 **********************************************************************************************************************/

type DeleteCatalogEntityResponse struct{}

func (c *CatalogEntitiesClient) Delete(ctx context.Context, tag string) error {
	entityResponse := &DeleteCatalogEntityResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Delete(BaseUris["catalog_entities"]+tag).Receive(entityResponse, apiError)
	if err != nil {
		return errors.New("could not delete catalog entity: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
