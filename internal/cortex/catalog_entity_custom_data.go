package cortex

import (
	"context"
	"errors"
	"github.com/dghubble/sling"
)

type CatalogEntityCustomDataClientInterface interface {
	Get(ctx context.Context, entityTag string, key string) (*CatalogEntityCustomData, error)
	List(ctx context.Context, entityTag string, params *CatalogEntityCustomDataListParams) ([]CatalogEntityCustomData, error)
}

type CatalogEntityCustomDataClient struct {
	client *HttpClient
}

var _ CatalogEntityCustomDataClientInterface = &CatalogEntityCustomDataClient{}

func (c *CatalogEntityCustomDataClient) Client() *sling.Sling {
	return c.client.client
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

type CatalogEntityCustomData struct {
	ID          string      `json:"id"`
	Key         string      `json:"key"`
	Description string      `json:"description,omitempty"`
	Source      string      `json:"source"`
	Value       interface{} `json:"value"`
	DateUpdated string      `json:"dateUpdated,omitempty"`
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/:tag/custom-data/:key
 **********************************************************************************************************************/

func (c *CatalogEntityCustomDataClient) Get(ctx context.Context, entityTag string, key string) (*CatalogEntityCustomData, error) {
	entityResponse := &CatalogEntityCustomData{}
	apiError := &ApiError{}
	response, err := c.Client().Get(BaseUris["catalog_entities"]+entityTag+"/custom-data/"+key).Receive(entityResponse, apiError)
	if err != nil {
		return entityResponse, errors.New("could not get catalog entity custom data: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return entityResponse, errors.Join(errors.New("Failed getting catalog entity custom data: "), err)
	}

	return entityResponse, nil
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/:tag/custom-data
 **********************************************************************************************************************/

// CatalogEntityCustomDataListParams are the query parameters for the GET /v1/catalog/:tag/custom-data endpoint.
type CatalogEntityCustomDataListParams struct{}

// List retrieves a list of scorecards based on a query.
func (c *CatalogEntityCustomDataClient) List(ctx context.Context, entityTag string, params *CatalogEntityCustomDataListParams) ([]CatalogEntityCustomData, error) {
	var entitiesResponse []CatalogEntityCustomData
	apiError := &ApiError{}

	response, err := c.Client().Get(BaseUris["catalog_entities"]+entityTag+"/custom-data").QueryStruct(params).Receive(entitiesResponse, apiError)
	if err != nil {
		return nil, errors.New("could not get catalog entity custom data: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return nil, err
	}

	return entitiesResponse, nil
}
