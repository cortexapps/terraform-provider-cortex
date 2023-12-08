package cortex_test

import (
	"context"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCatalogCustomDataEntity = &cortex.CatalogEntityCustomData{
	Tag:         "test-catalog-entity",
	Key:         "example-key",
	Description: "example description",
	Value: map[string]interface{}{
		"example": "value",
		"nested": map[string]interface{}{
			"one": 2,
		},
	},
}

func TestGetCatalogEntityCustomData(t *testing.T) {
	tag := testCatalogCustomDataEntity.Tag
	key := testCatalogCustomDataEntity.Key
	c, teardown, err := setupClient(
		cortex.Route("catalog_entities", tag+"/custom-data/"+key),
		testCatalogCustomDataEntity,
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.CatalogEntityCustomData().Get(context.Background(), tag, key)
	assert.Nil(t, err, "error retrieving catalog entity custom data")
	assert.Equal(t, testCatalogCustomDataEntity.Tag, res.Tag)
}

func TestUpsertCatalogEntityCustomData(t *testing.T) {
	tag := testCatalogCustomDataEntity.Tag
	req := cortex.UpsertCatalogEntityCustomDataRequest{
		Key:         testCatalogCustomDataEntity.Key,
		Description: testCatalogCustomDataEntity.Description,
		Value:       testCatalogCustomDataEntity.Value,
		Force:       true,
	}
	c, teardown, err := setupClient(
		cortex.Route("catalog_entities", tag+"/custom-data"),
		testCatalogCustomDataEntity,
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	updatedEntity, err := c.CatalogEntityCustomData().Upsert(context.Background(), tag, req)
	assert.Nil(t, err, "error upserting catalog entity custom data")
	assert.Equal(t, updatedEntity.Tag, tag)
	assert.Equal(t, updatedEntity.Key, req.Key)
	assert.Equal(t, updatedEntity.Description, req.Description)
}

func TestDeleteCatalogEntityCustomData(t *testing.T) {
	tag := testCatalogCustomDataEntity.Tag
	key := testCatalogCustomDataEntity.Key

	c, teardown, err := setupClient(
		cortex.Route("catalog_entities", tag+"/custom-data"),
		cortex.DeleteCatalogEntityCustomDataResponse{},
		AssertRequestMethod(t, "DELETE"),
		AssertRequestURI(t, "/api/v1/catalog/"+tag+"/custom-data?force=true&key="+key),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.CatalogEntityCustomData().Delete(context.Background(), tag, key)
	assert.Nil(t, err, "error deleting catalog entity custom data")
}
