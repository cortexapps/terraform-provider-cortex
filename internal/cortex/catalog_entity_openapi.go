package cortex

import (
	"context"
	"errors"
	"fmt"

	"github.com/dghubble/sling"
)

type CatalogEntityOpenAPIClientInterface interface {
	Get(ctx context.Context, entityTag string) (CatalogEntityOpenAPI, error)
	Upsert(ctx context.Context, entityTag string, req UpsertCatalogEntityOpenAPIRequest) (CatalogEntityOpenAPI, error)
	Delete(ctx context.Context, entityTag string) error
}

type CatalogEntityOpenAPIClient struct {
	client *HttpClient
}

var _ CatalogEntityOpenAPIClientInterface = &CatalogEntityOpenAPIClient{}

func (c *CatalogEntityOpenAPIClient) Client() *sling.Sling {
	return c.client.Client()
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

type CatalogEntityOpenAPI struct {
	Tag  string `json:"tag"`  // tag of catalog entity
	Spec string `json:"spec"` // OpenAPI specification
}

func (c *CatalogEntityOpenAPI) ID() string {
	return c.Tag
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/:tag/documentation/openapi
 **********************************************************************************************************************/

func (c *CatalogEntityOpenAPIClient) Get(ctx context.Context, entityTag string) (CatalogEntityOpenAPI, error) {
	entity := CatalogEntityOpenAPI{}
	apiError := ApiError{}
	response, err := c.Client().Get(Route("catalog_entities", entityTag+"/documentation/openapi")).Receive(&entity, &apiError)
	if err != nil {
		return entity, errors.New("could not get catalog entity OpenAPI spec: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return entity, errors.Join(errors.New("failed getting catalog entity OpenAPI spec: "), err)
	}

	entity.Tag = entityTag
	return entity, nil
}

/***********************************************************************************************************************
 * POST /api/v1/catalog/:tag/documentation/openapi
 **********************************************************************************************************************/

type UpsertCatalogEntityOpenAPIRequest struct {
	Spec  string `json:"spec"`
	Force bool   `json:"force" url:"force,omitempty"`
}

func (c *CatalogEntityOpenAPI) ToUpsertRequest() UpsertCatalogEntityOpenAPIRequest {
	return UpsertCatalogEntityOpenAPIRequest{
		Spec: c.Spec,
	}
}

func (c *CatalogEntityOpenAPIClient) Upsert(ctx context.Context, entityTag string, req UpsertCatalogEntityOpenAPIRequest) (CatalogEntityOpenAPI, error) {
	entity := CatalogEntityOpenAPI{}
	apiError := ApiError{}

	req.Force = true

	body, err := c.Client().Put(Route("catalog_entities", entityTag+"/documentation/openapi")).BodyJSON(&req).Receive(&entity, &apiError)
	if err != nil {
		return entity, fmt.Errorf("failed upserting OpenAPI spec for entity: %+v", err)
	}

	err = c.client.handleResponseStatus(body, &apiError)
	if err != nil {
		return entity, err
	}

	entity.Tag = entityTag
	return entity, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/catalog/:tag/documentation/openapi
 **********************************************************************************************************************/

func (c *CatalogEntityOpenAPIClient) Delete(ctx context.Context, entityTag string) error {
	apiError := ApiError{}

	response, err := c.Client().Delete(Route("catalog_entities", entityTag+"/documentation/openapi")).Receive(nil, &apiError)
	if err != nil {
		return errors.New("could not delete OpenAPI spec: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return errors.Join(errors.New("failed deleting OpenAPI spec: "), err)
	}

	return nil
}
