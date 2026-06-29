package cortex_test

import (
	"context"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCatalogEntity = &cortex.CatalogEntity{
	Tag: "test-catalog-entity",
}

func TestGetCatalogEntity(t *testing.T) {
	testTag := "test-catalog-entity"
	c, teardown, err := setupClient(cortex.Route("catalog_entities", testTag), testCatalogEntity, AssertRequestMethod(t, "GET"))
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
	c, teardown, err := setupClient(cortex.Route("catalog_entities", ""), resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	var queryParams cortex.CatalogEntityListParams
	res, err := c.CatalogEntities().List(context.Background(), &queryParams)
	assert.Nil(t, err, "error retrieving entities")
	assert.NotEmpty(t, res.Entities, "returned no entities")
	assert.Equal(t, res.Entities[0].Tag, firstTag)
}

func TestListCatalogEntitiesWithOwners(t *testing.T) {
	resp := &cortex.CatalogEntitiesResponse{
		Entities: []cortex.CatalogEntity{
			*testCatalogEntity,
		},
	}
	c, teardown, err := setupClient(cortex.Route("catalog_entities", ""), resp,
		AssertRequestMethod(t, "GET"),
		AssertRequestURI(t, "/api/v1/catalog/?includeOwners=true"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	queryParams := cortex.CatalogEntityListParams{IncludeOwners: true}
	res, err := c.CatalogEntities().List(context.Background(), &queryParams)
	assert.Nil(t, err, "error retrieving entities")
	assert.NotEmpty(t, res.Entities, "returned no entities")
}

func TestListCatalogEntitiesDeserializesOwnersAndGit(t *testing.T) {
	entity := cortex.CatalogEntity{
		Tag: "test-catalog-entity",
		Owners: cortex.CatalogEntityOwnersResponse{
			Teams: []cortex.CatalogEntityOwnerTeam{
				{
					Name:        "My Team",
					Tag:         "my-team",
					Description: "Owns core services",
					Provider:    "GITHUB",
					Inheritance: "NONE",
				},
			},
			Individuals: []cortex.CatalogEntityOwnerIndividual{
				{
					Email:       "owner@example.com",
					Description: "Lead engineer",
				},
			},
		},
		Git: cortex.CatalogEntityGitSummary{
			Provider:      "GITHUB",
			Repository:    "my-org/my-repo",
			RepositoryUrl: "https://github.com/my-org/my-repo",
		},
	}
	resp := &cortex.CatalogEntitiesResponse{
		Entities: []cortex.CatalogEntity{entity},
	}
	c, teardown, err := setupClient(cortex.Route("catalog_entities", ""), resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	queryParams := cortex.CatalogEntityListParams{IncludeOwners: true}
	res, err := c.CatalogEntities().List(context.Background(), &queryParams)
	assert.Nil(t, err, "error retrieving entities")
	assert.Len(t, res.Entities, 1)

	got := res.Entities[0]

	assert.Len(t, got.Owners.Teams, 1)
	assert.Equal(t, "my-team", got.Owners.Teams[0].Tag)
	assert.Equal(t, "My Team", got.Owners.Teams[0].Name)
	assert.Equal(t, "GITHUB", got.Owners.Teams[0].Provider)
	assert.Equal(t, "NONE", got.Owners.Teams[0].Inheritance)

	assert.Len(t, got.Owners.Individuals, 1)
	assert.Equal(t, "owner@example.com", got.Owners.Individuals[0].Email)
	assert.Equal(t, "Lead engineer", got.Owners.Individuals[0].Description)

	assert.Equal(t, "GITHUB", got.Git.Provider)
	assert.Equal(t, "my-org/my-repo", got.Git.Repository)
	assert.Equal(t, "https://github.com/my-org/my-repo", got.Git.RepositoryUrl)
}
