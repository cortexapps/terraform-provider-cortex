package provider

import (
	"context"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestCatalogResourceModel_FromApiModel_Minimal(t *testing.T) {
	entity := cortex.Catalog{
		Slug:            "my-catalog",
		Name:            "My Catalog",
		IconTag:         "cortex",
		IsDraft:         false,
		Type:            "FILTER",
		IsCortexManaged: false,
	}

	ctx := context.Background()
	var diagnostics diag.Diagnostics
	model := NewCatalogResourceModel()
	model.FromApiModel(ctx, &diagnostics, entity)

	assert.False(t, diagnostics.HasError())
	assert.Equal(t, types.StringValue("my-catalog"), model.Id)
	assert.Equal(t, types.StringValue("my-catalog"), model.Slug)
	assert.Equal(t, types.StringValue("My Catalog"), model.Name)
	assert.Equal(t, types.StringValue("cortex"), model.IconTag)
	assert.Equal(t, types.BoolValue(false), model.IsDraft)
	assert.Equal(t, types.StringValue("FILTER"), model.Type)
	assert.Equal(t, types.BoolValue(false), model.IsCortexManaged)
	assert.True(t, model.Description.IsNull())
	assert.Nil(t, model.Filter)
}

func TestCatalogResourceModel_FromApiModel_WithFilter(t *testing.T) {
	entity := cortex.Catalog{
		Slug:    "filtered-catalog",
		Name:    "Filtered Catalog",
		IconTag: "cortex",
		IsDraft: false,
		Filter: &cortex.CatalogFilter{
			Query: "tag != nil",
			Types: &cortex.CatalogTypeFilter{
				Include: []string{"service", "domain"},
			},
			Groups: &cortex.CatalogGroupFilter{
				Exclude: []string{"my-group"},
			},
		},
	}

	ctx := context.Background()
	var diagnostics diag.Diagnostics
	model := NewCatalogResourceModel()
	model.FromApiModel(ctx, &diagnostics, entity)

	assert.False(t, diagnostics.HasError())
	assert.NotNil(t, model.Filter)
	assert.Equal(t, types.StringValue("tag != nil"), model.Filter.Query)

	// Types filter
	assert.NotNil(t, model.Filter.Types)
	assert.False(t, model.Filter.Types.Include.IsNull())
	includeElems := model.Filter.Types.Include.Elements()
	assert.Len(t, includeElems, 2)
	assert.Contains(t, includeElems, types.StringValue("service"))
	assert.Contains(t, includeElems, types.StringValue("domain"))
	// Exclude not set — should be null
	assert.True(t, model.Filter.Types.Exclude.IsNull())

	// Groups filter
	assert.NotNil(t, model.Filter.Groups)
	assert.False(t, model.Filter.Groups.Exclude.IsNull())
	excludeElems := model.Filter.Groups.Exclude.Elements()
	assert.Len(t, excludeElems, 1)
	assert.Contains(t, excludeElems, types.StringValue("my-group"))
	assert.True(t, model.Filter.Groups.Include.IsNull())
}

func TestCatalogResourceModel_ToApiModel_Minimal(t *testing.T) {
	model := CatalogResourceModel{
		Id:                  types.StringValue("my-catalog"),
		Slug:                types.StringValue("my-catalog"),
		Name:                types.StringValue("My Catalog"),
		IconTag:             types.StringValue("cortex"),
		IsDraft:             types.BoolValue(false),
		Description:         types.StringNull(),
		Type:                types.StringNull(),
		RelationshipTypeTag: types.StringNull(),
		Filter:              nil,
		IsCortexManaged:     types.BoolValue(false),
	}

	ctx := context.Background()
	var diagnostics diag.Diagnostics
	req := model.ToApiModel(ctx, &diagnostics)

	assert.False(t, diagnostics.HasError())
	assert.Equal(t, "my-catalog", req.Slug)
	assert.Equal(t, "My Catalog", req.Name)
	assert.Equal(t, "cortex", req.IconTag)
	assert.Equal(t, false, req.IsDraft)
	assert.Equal(t, "", req.Description)
	assert.Nil(t, req.Filter)
}

func TestCatalogResourceModel_RoundTrip(t *testing.T) {
	entity := cortex.Catalog{
		Slug:            "round-trip-catalog",
		Name:            "Round Trip Catalog",
		IconTag:         "my-icon",
		IsDraft:         true,
		Description:     "some description",
		Type:            "FILTER",
		IsCortexManaged: false,
	}

	ctx := context.Background()
	var diagnostics diag.Diagnostics

	// FromApiModel
	model := NewCatalogResourceModel()
	model.FromApiModel(ctx, &diagnostics, entity)
	assert.False(t, diagnostics.HasError())

	// ToApiModel
	req := model.ToApiModel(ctx, &diagnostics)
	assert.False(t, diagnostics.HasError())

	assert.Equal(t, entity.Slug, req.Slug)
	assert.Equal(t, entity.Name, req.Name)
	assert.Equal(t, entity.IconTag, req.IconTag)
	assert.Equal(t, entity.IsDraft, req.IsDraft)
	assert.Equal(t, entity.Description, req.Description)
	assert.Equal(t, entity.Type, req.Type)
	assert.Nil(t, req.Filter)
}

func TestCatalogResourceModel_RoundTrip_WithFilter(t *testing.T) {
	entity := cortex.Catalog{
		Slug:    "filter-round-trip",
		Name:    "Filter Round Trip",
		IconTag: "cortex",
		IsDraft: false,
		Filter: &cortex.CatalogFilter{
			Query: "tag != nil",
			Types: &cortex.CatalogTypeFilter{
				Include: []string{"service"},
				Exclude: []string{"job"},
			},
			Groups: &cortex.CatalogGroupFilter{
				Include: []string{"group-a"},
			},
		},
	}

	ctx := context.Background()
	var diagnostics diag.Diagnostics

	model := NewCatalogResourceModel()
	model.FromApiModel(ctx, &diagnostics, entity)
	assert.False(t, diagnostics.HasError())

	req := model.ToApiModel(ctx, &diagnostics)
	assert.False(t, diagnostics.HasError())

	assert.Equal(t, entity.Slug, req.Slug)
	assert.NotNil(t, req.Filter)
	assert.Equal(t, "tag != nil", req.Filter.Query)
	assert.NotNil(t, req.Filter.Types)
	assert.Equal(t, []string{"service"}, req.Filter.Types.Include)
	assert.Equal(t, []string{"job"}, req.Filter.Types.Exclude)
	assert.NotNil(t, req.Filter.Groups)
	assert.Equal(t, []string{"group-a"}, req.Filter.Groups.Include)
}

func TestCatalogDataSourceModel_FromApiModel(t *testing.T) {
	entity := cortex.Catalog{
		Slug:            "ds-catalog",
		Name:            "DS Catalog",
		IconTag:         "cortex",
		IsDraft:         false,
		IsCortexManaged: true,
	}

	ctx := context.Background()
	var diagnostics diag.Diagnostics
	model := NewCatalogDataSourceModel()
	model.FromApiModel(ctx, &diagnostics, entity)

	assert.False(t, diagnostics.HasError())
	assert.Equal(t, types.StringValue("ds-catalog"), model.Id)
	assert.Equal(t, types.StringValue("ds-catalog"), model.Slug)
	assert.Equal(t, types.StringValue("DS Catalog"), model.Name)
	assert.Equal(t, types.BoolValue(true), model.IsCortexManaged)
	assert.True(t, model.Description.IsNull())
	assert.Nil(t, model.Filter)
}
