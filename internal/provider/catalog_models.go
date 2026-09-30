package provider

import (
	"context"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Filter Sub-Models
 **********************************************************************************************************************/

type CatalogTypeFilterModel struct {
	Include types.List `tfsdk:"include"`
	Exclude types.List `tfsdk:"exclude"`
}

type CatalogGroupFilterModel struct {
	Include types.List `tfsdk:"include"`
	Exclude types.List `tfsdk:"exclude"`
}

type CatalogFilterModel struct {
	Query  types.String             `tfsdk:"query"`
	Types  *CatalogTypeFilterModel  `tfsdk:"types"`
	Groups *CatalogGroupFilterModel `tfsdk:"groups"`
}

/***********************************************************************************************************************
 * Resource Model
 **********************************************************************************************************************/

// CatalogResourceModel describes the catalog resource data model within Terraform.
type CatalogResourceModel struct {
	Id                  types.String        `tfsdk:"id"`
	Slug                types.String        `tfsdk:"slug"`
	Name                types.String        `tfsdk:"name"`
	Description         types.String        `tfsdk:"description"`
	IconTag             types.String        `tfsdk:"icon_tag"`
	IsDraft             types.Bool          `tfsdk:"is_draft"`
	Type                types.String        `tfsdk:"type"`
	RelationshipTypeTag types.String        `tfsdk:"relationship_type_tag"`
	Filter              *CatalogFilterModel `tfsdk:"filter"`
	IsCortexManaged     types.Bool          `tfsdk:"is_cortex_managed"`
}

// NewCatalogResourceModel returns a zero-value CatalogResourceModel.
func NewCatalogResourceModel() CatalogResourceModel {
	return CatalogResourceModel{}
}

func (r *CatalogResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.Catalog) {
	r.Id = types.StringValue(entity.Slug)
	r.Slug = types.StringValue(entity.Slug)
	r.Name = types.StringValue(entity.Name)
	r.IconTag = types.StringValue(entity.IconTag)
	r.IsDraft = types.BoolValue(entity.IsDraft)
	r.IsCortexManaged = types.BoolValue(entity.IsCortexManaged)

	if entity.Description != "" {
		r.Description = types.StringValue(entity.Description)
	} else {
		r.Description = types.StringNull()
	}

	if entity.Type != "" {
		r.Type = types.StringValue(entity.Type)
	} else {
		r.Type = types.StringNull()
	}

	if entity.RelationshipTypeTag != "" {
		r.RelationshipTypeTag = types.StringValue(entity.RelationshipTypeTag)
	} else {
		r.RelationshipTypeTag = types.StringNull()
	}

	if entity.Filter == nil {
		r.Filter = nil
	} else {
		r.Filter = catalogFilterFromApiModel(entity.Filter)
	}
}

func (r *CatalogResourceModel) ToApiModel(ctx context.Context, diagnostics *diag.Diagnostics) cortex.UpsertCatalogRequest {
	req := cortex.UpsertCatalogRequest{
		Slug:    r.Slug.ValueString(),
		Name:    r.Name.ValueString(),
		IconTag: r.IconTag.ValueString(),
		IsDraft: r.IsDraft.ValueBool(),
	}

	if !r.Description.IsNull() && !r.Description.IsUnknown() {
		req.Description = r.Description.ValueString()
	}

	if !r.Type.IsNull() && !r.Type.IsUnknown() {
		req.Type = r.Type.ValueString()
	}

	if !r.RelationshipTypeTag.IsNull() && !r.RelationshipTypeTag.IsUnknown() {
		req.RelationshipTypeTag = r.RelationshipTypeTag.ValueString()
	}

	if r.Filter != nil {
		req.Filter = catalogFilterToApiModel(ctx, diagnostics, r.Filter)
	}

	return req
}

/***********************************************************************************************************************
 * Data Source Model
 **********************************************************************************************************************/

// CatalogDataSourceModel describes the catalog data source data model within Terraform.
type CatalogDataSourceModel struct {
	Id                  types.String        `tfsdk:"id"`
	Slug                types.String        `tfsdk:"slug"`
	Name                types.String        `tfsdk:"name"`
	Description         types.String        `tfsdk:"description"`
	IconTag             types.String        `tfsdk:"icon_tag"`
	IsDraft             types.Bool          `tfsdk:"is_draft"`
	Type                types.String        `tfsdk:"type"`
	RelationshipTypeTag types.String        `tfsdk:"relationship_type_tag"`
	Filter              *CatalogFilterModel `tfsdk:"filter"`
	IsCortexManaged     types.Bool          `tfsdk:"is_cortex_managed"`
}

// NewCatalogDataSourceModel returns a zero-value CatalogDataSourceModel.
func NewCatalogDataSourceModel() CatalogDataSourceModel {
	return CatalogDataSourceModel{}
}

func (r *CatalogDataSourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.Catalog) {
	r.Id = types.StringValue(entity.Slug)
	r.Slug = types.StringValue(entity.Slug)
	r.Name = types.StringValue(entity.Name)
	r.IconTag = types.StringValue(entity.IconTag)
	r.IsDraft = types.BoolValue(entity.IsDraft)
	r.IsCortexManaged = types.BoolValue(entity.IsCortexManaged)

	if entity.Description != "" {
		r.Description = types.StringValue(entity.Description)
	} else {
		r.Description = types.StringNull()
	}

	if entity.Type != "" {
		r.Type = types.StringValue(entity.Type)
	} else {
		r.Type = types.StringNull()
	}

	if entity.RelationshipTypeTag != "" {
		r.RelationshipTypeTag = types.StringValue(entity.RelationshipTypeTag)
	} else {
		r.RelationshipTypeTag = types.StringNull()
	}

	if entity.Filter == nil {
		r.Filter = nil
	} else {
		r.Filter = catalogFilterFromApiModel(entity.Filter)
	}
}

/***********************************************************************************************************************
 * Filter Conversion Helpers
 **********************************************************************************************************************/

func stringsToAttrValues(ss []string) []attr.Value {
	vals := make([]attr.Value, len(ss))
	for i, s := range ss {
		vals[i] = types.StringValue(s)
	}
	return vals
}

func catalogFilterFromApiModel(f *cortex.CatalogFilter) *CatalogFilterModel {
	model := &CatalogFilterModel{}

	if f.Query != "" {
		model.Query = types.StringValue(f.Query)
	} else {
		model.Query = types.StringNull()
	}

	if f.Types != nil {
		tf := &CatalogTypeFilterModel{}
		if len(f.Types.Include) > 0 {
			tf.Include = types.ListValueMust(types.StringType, stringsToAttrValues(f.Types.Include))
		} else {
			tf.Include = types.ListNull(types.StringType)
		}
		if len(f.Types.Exclude) > 0 {
			tf.Exclude = types.ListValueMust(types.StringType, stringsToAttrValues(f.Types.Exclude))
		} else {
			tf.Exclude = types.ListNull(types.StringType)
		}
		model.Types = tf
	} else {
		model.Types = nil
	}

	if f.Groups != nil {
		gf := &CatalogGroupFilterModel{}
		if len(f.Groups.Include) > 0 {
			gf.Include = types.ListValueMust(types.StringType, stringsToAttrValues(f.Groups.Include))
		} else {
			gf.Include = types.ListNull(types.StringType)
		}
		if len(f.Groups.Exclude) > 0 {
			gf.Exclude = types.ListValueMust(types.StringType, stringsToAttrValues(f.Groups.Exclude))
		} else {
			gf.Exclude = types.ListNull(types.StringType)
		}
		model.Groups = gf
	} else {
		model.Groups = nil
	}

	return model
}

func catalogFilterToApiModel(ctx context.Context, diagnostics *diag.Diagnostics, f *CatalogFilterModel) *cortex.CatalogFilter {
	apiFilter := &cortex.CatalogFilter{}

	if !f.Query.IsNull() && !f.Query.IsUnknown() {
		apiFilter.Query = f.Query.ValueString()
	}

	if f.Types != nil {
		tf := &cortex.CatalogTypeFilter{}
		if !f.Types.Include.IsNull() && !f.Types.Include.IsUnknown() {
			var include []string
			d := f.Types.Include.ElementsAs(ctx, &include, false)
			diagnostics.Append(d...)
			tf.Include = include
		}
		if !f.Types.Exclude.IsNull() && !f.Types.Exclude.IsUnknown() {
			var exclude []string
			d := f.Types.Exclude.ElementsAs(ctx, &exclude, false)
			diagnostics.Append(d...)
			tf.Exclude = exclude
		}
		apiFilter.Types = tf
	}

	if f.Groups != nil {
		gf := &cortex.CatalogGroupFilter{}
		if !f.Groups.Include.IsNull() && !f.Groups.Include.IsUnknown() {
			var include []string
			d := f.Groups.Include.ElementsAs(ctx, &include, false)
			diagnostics.Append(d...)
			gf.Include = include
		}
		if !f.Groups.Exclude.IsNull() && !f.Groups.Exclude.IsUnknown() {
			var exclude []string
			d := f.Groups.Exclude.ElementsAs(ctx, &exclude, false)
			diagnostics.Append(d...)
			gf.Exclude = exclude
		}
		apiFilter.Groups = gf
	}

	return apiFilter
}
