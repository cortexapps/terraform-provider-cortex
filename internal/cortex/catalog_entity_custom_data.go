package cortex

import (
	"context"
	"errors"
	"fmt"
	"github.com/dghubble/sling"
)

type CatalogEntityCustomDataClientInterface interface {
	Get(ctx context.Context, entityTag string, key string) (CatalogEntityCustomData, error)
	List(ctx context.Context, entityTag string, params CatalogEntityCustomDataListParams) ([]CatalogEntityCustomData, error)
	Upsert(ctx context.Context, entityTag string, req UpsertCatalogEntityCustomDataRequest) (CatalogEntityCustomData, error)
	Delete(ctx context.Context, entityTag string, key string) error
}

type CatalogEntityCustomDataClient struct {
	client *HttpClient
}

var _ CatalogEntityCustomDataClientInterface = &CatalogEntityCustomDataClient{}

func (c *CatalogEntityCustomDataClient) Client() *sling.Sling {
	return c.client.Client()
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

type CatalogEntityCustomData struct {
	Tag         string      `json:"tag"` // tag of catalog entity
	Key         string      `json:"key"` // key of custom data
	Description string      `json:"description,omitempty"`
	Source      string      `json:"source"`
	Value       interface{} `json:"value"`
	DateUpdated string      `json:"dateUpdated,omitempty"`
}

func (c *CatalogEntityCustomData) ID() string {
	return c.Tag + ":" + c.Key
}

func (c *CatalogEntityCustomData) ValueAsString() (string, error) {
	value := ""
	if c.Value != nil {
		err := error(nil)
		value, err = InterfaceToString(c.Value)
		if err != nil {
			return "", err
		}
	}
	return value, nil
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/:tag/custom-data/:key
 **********************************************************************************************************************/

func (c *CatalogEntityCustomDataClient) Get(ctx context.Context, entityTag string, key string) (CatalogEntityCustomData, error) {
	entity := CatalogEntityCustomData{}
	apiError := ApiError{}
	response, err := c.Client().Get(Route("catalog_entities", entityTag+"/custom-data/"+key)).Receive(&entity, &apiError)
	if err != nil {
		return entity, errors.New("could not get catalog entity custom data: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return entity, errors.Join(errors.New("Failed getting catalog entity custom data: "), err)
	}

	entity.Tag = entityTag
	return entity, nil
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/:tag/custom-data
 **********************************************************************************************************************/

// CatalogEntityCustomDataListParams are the query parameters for the GET /v1/catalog/:tag/custom-data endpoint.
type CatalogEntityCustomDataListParams struct{}

// List retrieves a list of scorecards based on a query.
func (c *CatalogEntityCustomDataClient) List(ctx context.Context, entityTag string, params CatalogEntityCustomDataListParams) ([]CatalogEntityCustomData, error) {
	var entities []CatalogEntityCustomData
	apiError := ApiError{}

	response, err := c.Client().Get(Route("catalog_entities", entityTag+"/custom-data")).QueryStruct(&params).Receive(entities, &apiError)
	if err != nil {
		return nil, errors.New("could not get catalog entity custom data: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return nil, err
	}

	for entitiesIndex := range entities {
		entities[entitiesIndex].Tag = entityTag
	}
	return entities, nil
}

/***********************************************************************************************************************
 * POST /api/v1/catalog/:tag/custom-data
 **********************************************************************************************************************/

type UpsertCatalogEntityCustomDataRequest struct {
	Key         string      `json:"key"`
	Description string      `json:"description,omitempty"`
	Value       interface{} `json:"value"`
	Force       bool        `json:"force" url:"force,omitempty"`
}

// ToUpsertRequest https://docs.cortex.io/docs/api/add-custom-data-for-entity
func (c *CatalogEntityCustomData) ToUpsertRequest() UpsertCatalogEntityCustomDataRequest {
	return UpsertCatalogEntityCustomDataRequest{
		Key:         c.Key,
		Description: c.Description,
		Value:       c.Value,
	}
}

func (c *CatalogEntityCustomDataClient) Upsert(ctx context.Context, entityTag string, req UpsertCatalogEntityCustomDataRequest) (CatalogEntityCustomData, error) {
	entity := CatalogEntityCustomData{}
	apiError := ApiError{}

	req.Force = true

	body, err := c.Client().Post(Route("catalog_entities", entityTag+"/custom-data")).BodyJSON(&req).Receive(&entity, &apiError)
	if err != nil {
		return entity, fmt.Errorf("failed upserting custom data for entity: %+v", err)
	}

	err = c.client.handleResponseStatus(body, &apiError)
	if err != nil {
		return entity, err
	}

	entity.Tag = entityTag
	return entity, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/catalog/:tag/custom-data - Delete custom data for a catalog entity by key
 **********************************************************************************************************************/

type DeleteCatalogEntityCustomDataRequest struct {
	Key   string `json:"key" url:"key"`
	Force bool   `json:"force" url:"force,omitempty"`
}
type DeleteCatalogEntityCustomDataResponse struct{}

func (c *CatalogEntityCustomDataClient) Delete(ctx context.Context, entityTag string, key string) error {
	response := DeleteCatalogEntityCustomDataResponse{}
	apiError := ApiError{}
	params := DeleteCatalogEntityCustomDataRequest{
		Key:   key,
		Force: true,
	}

	body, err := c.Client().Delete(Route("catalog_entities", entityTag+"/custom-data")).QueryStruct(&params).Receive(&response, &apiError)
	if err != nil {
		return errors.New("could not delete custom data for catalog entity: " + err.Error())
	}

	err = c.client.handleResponseStatus(body, &apiError)
	if err != nil {
		return err
	}

	return nil
}
