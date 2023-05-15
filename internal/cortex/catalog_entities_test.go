package cortex_test

import (
	"context"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCatalogEntity = &cortex.CatalogEntity{
	Tag: "test-catalog-entity",
}

func TestGetCatalogEntity(t *testing.T) {
	testTag := "test-catalog-entity"
	c, teardown, err := setupClient(cortex.BaseUris["catalog_entities"]+testTag, testCatalogEntity, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.CatalogEntities().Get(context.Background(), testTag)
	assert.Nil(t, err, "error retrieving a catalog entity")
	assert.Equal(t, testCatalogEntity, res)
}

func TestListCatalogEntities(t *testing.T) {
	firstTag := "test-catalog-entity"
	resp := &cortex.CatalogEntitiesResponse{
		Entities: []cortex.CatalogEntity{
			*testCatalogEntity,
		},
	}
	c, teardown, err := setupClient(cortex.BaseUris["catalog_entities"], resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	var queryParams cortex.CatalogEntityListParams
	res, err := c.CatalogEntities().List(context.Background(), &queryParams)
	assert.Nil(t, err, "error retrieving entities")
	assert.NotEmpty(t, res.Entities, "returned no entities")
	assert.Equal(t, res.Entities[0].Tag, firstTag)
}
