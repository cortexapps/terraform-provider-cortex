package cortex

import (
	"context"
	"errors"
	"github.com/dghubble/sling"
)

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

type CatalogTypeFilter struct {
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}

type CatalogGroupFilter struct {
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}

type CatalogFilter struct {
	Query  string              `json:"query,omitempty"`
	Types  *CatalogTypeFilter  `json:"types,omitempty"`
	Groups *CatalogGroupFilter `json:"groups,omitempty"`
}

type Catalog struct {
	Slug                string         `json:"slug"`
	Name                string         `json:"name"`
	Description         string         `json:"description,omitempty"`
	IconTag             string         `json:"iconTag"`
	IsDraft             bool           `json:"isDraft"`
	Type                string         `json:"type,omitempty"`
	RelationshipTypeTag string         `json:"relationshipTypeTag,omitempty"`
	Filter              *CatalogFilter `json:"filter,omitempty"`
	IsCortexManaged     bool           `json:"isCortexManaged"`
}

type UpsertCatalogRequest struct {
	Slug                string         `json:"slug"`
	Name                string         `json:"name"`
	Description         string         `json:"description,omitempty"`
	IconTag             string         `json:"iconTag"`
	IsDraft             bool           `json:"isDraft"`
	Type                string         `json:"type,omitempty"`
	RelationshipTypeTag string         `json:"relationshipTypeTag,omitempty"`
	Filter              *CatalogFilter `json:"filter,omitempty"`
}

/***********************************************************************************************************************
 * Interface & Client
 **********************************************************************************************************************/

type CatalogsClientInterface interface {
	Get(ctx context.Context, slug string) (Catalog, error)
	Upsert(ctx context.Context, req UpsertCatalogRequest) (Catalog, error)
	Delete(ctx context.Context, slug string) error
}

type CatalogsClient struct {
	client *HttpClient
}

var _ CatalogsClientInterface = &CatalogsClient{}

func (c *CatalogsClient) Client() *sling.Sling {
	return c.client.Client()
}

/***********************************************************************************************************************
 * GET /api/v1/catalog-pages/:slug
 **********************************************************************************************************************/

func (c *CatalogsClient) Get(ctx context.Context, slug string) (Catalog, error) {
	data := Catalog{}
	apiError := ApiError{}

	response, err := c.Client().Get(Route("catalog_pages", slug)).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not get catalog: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, errors.Join(errors.New("Failed getting catalog: "), err)
	}

	return data, nil
}

/***********************************************************************************************************************
 * POST /api/v1/catalog-pages/
 **********************************************************************************************************************/

func (c *CatalogsClient) Upsert(ctx context.Context, req UpsertCatalogRequest) (Catalog, error) {
	data := Catalog{}
	apiError := ApiError{}

	response, err := c.Client().Post(Route("catalog_pages", "")).BodyJSON(&req).Receive(&data, &apiError)
	if err != nil {
		return data, errors.New("could not upsert catalog: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return data, err
	}

	return data, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/catalog-pages/:slug
 **********************************************************************************************************************/

type deleteCatalogResponse struct{}

func (c *CatalogsClient) Delete(ctx context.Context, slug string) error {
	deleteResponse := deleteCatalogResponse{}
	apiError := ApiError{}

	response, err := c.Client().Delete(Route("catalog_pages", slug)).Receive(&deleteResponse, &apiError)
	if err != nil {
		return errors.New("could not delete catalog: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return err
	}

	return nil
}
