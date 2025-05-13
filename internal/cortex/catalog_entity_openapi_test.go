package cortex_test

import (
	"context"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
)

var testOpenAPISpec = &cortex.CatalogEntityOpenAPI{
	Tag: "test-catalog-entity",
	Spec: `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: OK`,
}

func TestGetCatalogEntityOpenAPI(t *testing.T) {
	tag := testOpenAPISpec.Tag
	c, teardown, err := setupClient(
		cortex.Route("catalog_entities", tag+"/documentation/openapi"),
		testOpenAPISpec,
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.CatalogEntityOpenAPI().Get(context.Background(), tag)
	assert.Nil(t, err, "error retrieving catalog entity OpenAPI spec")
	assert.Equal(t, testOpenAPISpec.Tag, res.Tag)
	assert.Equal(t, testOpenAPISpec.Spec, res.Spec)
}

func TestUpsertCatalogEntityOpenAPI(t *testing.T) {
	tag := testOpenAPISpec.Tag
	req := cortex.UpsertCatalogEntityOpenAPIRequest{
		Spec:  testOpenAPISpec.Spec,
		Force: true,
	}
	c, teardown, err := setupClient(
		cortex.Route("catalog_entities", tag+"/documentation/openapi"),
		testOpenAPISpec,
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	updatedSpec, err := c.CatalogEntityOpenAPI().Upsert(context.Background(), tag, req)
	assert.Nil(t, err, "error upserting catalog entity OpenAPI spec")
	assert.Equal(t, updatedSpec.Tag, tag)
	assert.Equal(t, updatedSpec.Spec, req.Spec)
}

func TestDeleteCatalogEntityOpenAPI(t *testing.T) {
	tag := testOpenAPISpec.Tag

	c, teardown, err := setupClient(
		cortex.Route("catalog_entities", tag+"/documentation/openapi"),
		nil,
		AssertRequestMethod(t, "DELETE"),
		AssertRequestURI(t, "/api/v1/catalog/"+tag+"/documentation/openapi"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.CatalogEntityOpenAPI().Delete(context.Background(), tag)
	assert.Nil(t, err, "error deleting catalog entity OpenAPI spec")
}
