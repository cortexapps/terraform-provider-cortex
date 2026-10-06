package cortex_test

import (
	"context"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Compile-time interface compliance check.
var _ cortex.CatalogsClientInterface = &cortex.CatalogsClient{}

var testCatalogResponse = &cortex.Catalog{
	Slug:            "my-catalog",
	Name:            "My Catalog",
	IconTag:         "cortex",
	IsDraft:         false,
	IsCortexManaged: false,
}

func TestCatalogsClient_Get(t *testing.T) {
	slug := "my-catalog"
	c, teardown, err := setupClient(cortex.Route("catalog_pages", slug), testCatalogResponse, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Catalogs().Get(context.Background(), slug)
	assert.Nil(t, err, "error retrieving a catalog")
	assert.Equal(t, testCatalogResponse.Slug, res.Slug)
	assert.Equal(t, testCatalogResponse.Name, res.Name)
	assert.Equal(t, testCatalogResponse.IconTag, res.IconTag)
	assert.Equal(t, testCatalogResponse.IsDraft, res.IsDraft)
	assert.Equal(t, testCatalogResponse.IsCortexManaged, res.IsCortexManaged)
}

func TestCatalogsClient_Upsert(t *testing.T) {
	req := cortex.UpsertCatalogRequest{
		Slug:    testCatalogResponse.Slug,
		Name:    testCatalogResponse.Name,
		IconTag: testCatalogResponse.IconTag,
		IsDraft: testCatalogResponse.IsDraft,
	}
	c, teardown, err := setupClient(
		cortex.Route("catalog_pages", ""),
		testCatalogResponse,
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Catalogs().Upsert(context.Background(), req)
	assert.Nil(t, err, "error upserting a catalog")
	assert.Equal(t, res.Slug, req.Slug)
	assert.Equal(t, res.Name, req.Name)
	assert.Equal(t, res.IconTag, req.IconTag)
}

func TestCatalogsClient_Delete(t *testing.T) {
	slug := testCatalogResponse.Slug
	c, teardown, err := setupClient(
		cortex.Route("catalog_pages", slug),
		struct{}{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Catalogs().Delete(context.Background(), slug)
	assert.Nil(t, err, "error deleting a catalog")
}
