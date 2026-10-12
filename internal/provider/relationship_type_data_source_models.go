package provider

import (
	"context"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewRelationshipTypeDataSourceModel() RelationshipTypeDataSourceModel {
	return RelationshipTypeDataSourceModel{}
}

// RelationshipTypeDataSourceModel describes the data source data model.
type RelationshipTypeDataSourceModel struct {
	Tag                      types.String                 `tfsdk:"tag"`
	Name                     types.String                 `tfsdk:"name"`
	Description              types.String                 `tfsdk:"description"`
	AllowCycles              types.Bool                   `tfsdk:"allow_cycles"`
	IsSingleSource           types.Bool                   `tfsdk:"is_single_source"`
	IsSingleDestination      types.Bool                   `tfsdk:"is_single_destination"`
	DefinitionLocation       types.String                 `tfsdk:"definition_location"`
	SourcesFilter            *RelationshipTypeFilterModel `tfsdk:"sources_filter"`
	DestinationsFilter       *RelationshipTypeFilterModel `tfsdk:"destinations_filter"`
	SourceLabelSingular      types.String                 `tfsdk:"source_label_singular"`
	SourceLabelPlural        types.String                 `tfsdk:"source_label_plural"`
	DestinationLabelSingular types.String                 `tfsdk:"destination_label_singular"`
	DestinationLabelPlural   types.String                 `tfsdk:"destination_label_plural"`
	Inheritances             types.Set                    `tfsdk:"inheritances"`
	IsCortexManaged          types.Bool                   `tfsdk:"is_cortex_managed"`
	CreateCatalog            types.Bool                   `tfsdk:"create_catalog"`
}

func (r *RelationshipTypeDataSourceModel) FromApiModel(_ context.Context, _ *diag.Diagnostics, entity cortex.EntityRelationshipType) {
	r.Tag = types.StringValue(entity.Tag)
	r.Name = types.StringValue(entity.Name)
	if entity.Description != "" {
		r.Description = types.StringValue(entity.Description)
	} else {
		r.Description = types.StringNull()
	}
	r.AllowCycles = types.BoolValue(entity.AllowCycles)
	r.IsSingleSource = types.BoolValue(entity.IsSingleSource)
	r.IsSingleDestination = types.BoolValue(entity.IsSingleDestination)
	r.DefinitionLocation = types.StringValue(entity.DefinitionLocation)
	r.IsCortexManaged = types.BoolValue(entity.IsCortexManaged)
	r.CreateCatalog = types.BoolValue(entity.CreateCatalog)

	sf := filterFromApi(entity.SourcesFilter)
	r.SourcesFilter = &sf
	df := filterFromApi(entity.DestinationsFilter)
	r.DestinationsFilter = &df

	if entity.SourceLabelSingular != "" {
		r.SourceLabelSingular = types.StringValue(entity.SourceLabelSingular)
	} else {
		r.SourceLabelSingular = types.StringNull()
	}
	if entity.SourceLabelPlural != "" {
		r.SourceLabelPlural = types.StringValue(entity.SourceLabelPlural)
	} else {
		r.SourceLabelPlural = types.StringNull()
	}
	if entity.DestinationLabelSingular != "" {
		r.DestinationLabelSingular = types.StringValue(entity.DestinationLabelSingular)
	} else {
		r.DestinationLabelSingular = types.StringNull()
	}
	if entity.DestinationLabelPlural != "" {
		r.DestinationLabelPlural = types.StringValue(entity.DestinationLabelPlural)
	} else {
		r.DestinationLabelPlural = types.StringNull()
	}

	r.Inheritances = inheritancesToSet(entity.Inheritances)
}
