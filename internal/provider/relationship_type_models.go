package provider

import (
	"context"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Models
 **********************************************************************************************************************/

// inheritanceObjectType is the types.ObjectType for a single inheritance entry.
var inheritanceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"parameter_tag":    types.StringType,
		"inheritance_type": types.StringType,
	},
}

type RelationshipTypeFilterModel struct {
	Include   types.Bool `tfsdk:"include"`
	Types     types.Set  `tfsdk:"types"`
	Providers types.Set  `tfsdk:"providers"`
}

type RelationshipTypeResourceModel struct {
	ID                       types.String                 `tfsdk:"id"`
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

// stringsToSet converts a string slice to a types.Set of strings.
// Returns null for empty slices so that Optional attributes remain null
// when the user omits them, avoiding plan/state inconsistency.
// Callers should pair this with a SizeAtLeast(1) validator so that
// an explicit empty set (types = []) is rejected at plan time rather
// than silently converting to null in state.
func stringsToSet(ss []string) types.Set {
	if len(ss) == 0 {
		return types.SetNull(types.StringType)
	}
	vals := make([]attr.Value, len(ss))
	for i, s := range ss {
		vals[i] = types.StringValue(s)
	}
	return types.SetValueMust(types.StringType, vals)
}

// setToStrings converts a types.Set of strings to a string slice.
func setToStrings(s types.Set) []string {
	if s.IsNull() || s.IsUnknown() {
		return []string{}
	}
	elements := s.Elements()
	result := make([]string, len(elements))
	for i, elem := range elements {
		result[i] = elem.(types.String).ValueString()
	}
	return result
}

func filterFromApi(f cortex.EntityTypeFilter) RelationshipTypeFilterModel {
	return RelationshipTypeFilterModel{
		Include:   types.BoolValue(f.Include),
		Types:     stringsToSet(f.Types),
		Providers: stringsToSet(f.Providers),
	}
}

func filterToApi(f *RelationshipTypeFilterModel) cortex.EntityTypeFilter {
	if f == nil {
		return cortex.EntityTypeFilter{}
	}
	return cortex.EntityTypeFilter{
		Include:   f.Include.ValueBool(),
		Types:     setToStrings(f.Types),
		Providers: setToStrings(f.Providers),
	}
}

// inheritancesToSet converts API inheritances to a types.Set.
// Returns null for empty slices because inheritances is Optional (not Computed):
// when the user omits the block, state must be null to match the null plan value.
func inheritancesToSet(inheritances []cortex.RelationshipInheritance) types.Set {
	if len(inheritances) == 0 {
		return types.SetNull(inheritanceObjectType)
	}
	vals := make([]attr.Value, len(inheritances))
	for i, inh := range inheritances {
		obj, _ := types.ObjectValue(inheritanceObjectType.AttrTypes, map[string]attr.Value{
			"parameter_tag":    types.StringValue(inh.ParameterTag),
			"inheritance_type": types.StringValue(inh.InheritanceType),
		})
		vals[i] = obj
	}
	return types.SetValueMust(inheritanceObjectType, vals)
}

// inheritancesFromSet converts a types.Set back to API inheritance slice.
func inheritancesFromSet(s types.Set) []cortex.RelationshipInheritance {
	if s.IsNull() || s.IsUnknown() {
		return []cortex.RelationshipInheritance{}
	}
	result := make([]cortex.RelationshipInheritance, 0, len(s.Elements()))
	for _, elem := range s.Elements() {
		obj := elem.(types.Object)
		attrs := obj.Attributes()
		result = append(result, cortex.RelationshipInheritance{
			ParameterTag:    attrs["parameter_tag"].(types.String).ValueString(),
			InheritanceType: attrs["inheritance_type"].(types.String).ValueString(),
		})
	}
	return result
}

func (r *RelationshipTypeResourceModel) FromApiModel(_ context.Context, _ *diag.Diagnostics, entity cortex.EntityRelationshipType) {
	r.ID = types.StringValue(entity.Tag)
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

func (r *RelationshipTypeResourceModel) ToCreateRequest() cortex.CreateRelationshipTypeRequest {
	req := cortex.CreateRelationshipTypeRequest{
		Tag:                      r.Tag.ValueString(),
		Name:                     r.Name.ValueString(),
		Description:              r.Description.ValueString(),
		AllowCycles:              r.AllowCycles.ValueBool(),
		IsSingleSource:           r.IsSingleSource.ValueBool(),
		IsSingleDestination:      r.IsSingleDestination.ValueBool(),
		DefinitionLocation:       r.DefinitionLocation.ValueString(),
		SourcesFilter:            filterToApi(r.SourcesFilter),
		DestinationsFilter:       filterToApi(r.DestinationsFilter),
		SourceLabelSingular:      r.SourceLabelSingular.ValueString(),
		SourceLabelPlural:        r.SourceLabelPlural.ValueString(),
		DestinationLabelSingular: r.DestinationLabelSingular.ValueString(),
		DestinationLabelPlural:   r.DestinationLabelPlural.ValueString(),
		CreateCatalog:            r.CreateCatalog.ValueBool(),
		Inheritances:             inheritancesFromSet(r.Inheritances),
	}
	if req.Inheritances == nil {
		req.Inheritances = []cortex.RelationshipInheritance{}
	}
	return req
}

func (r *RelationshipTypeResourceModel) ToUpdateRequest() cortex.UpdateRelationshipTypeRequest {
	req := cortex.UpdateRelationshipTypeRequest{
		Tag:                      r.Tag.ValueString(),
		Name:                     r.Name.ValueString(),
		Description:              r.Description.ValueString(),
		AllowCycles:              r.AllowCycles.ValueBool(),
		IsSingleSource:           r.IsSingleSource.ValueBool(),
		IsSingleDestination:      r.IsSingleDestination.ValueBool(),
		DefinitionLocation:       r.DefinitionLocation.ValueString(),
		SourcesFilter:            filterToApi(r.SourcesFilter),
		DestinationsFilter:       filterToApi(r.DestinationsFilter),
		SourceLabelSingular:      r.SourceLabelSingular.ValueString(),
		SourceLabelPlural:        r.SourceLabelPlural.ValueString(),
		DestinationLabelSingular: r.DestinationLabelSingular.ValueString(),
		DestinationLabelPlural:   r.DestinationLabelPlural.ValueString(),
		CreateCatalog:            r.CreateCatalog.ValueBool(),
		Inheritances:             inheritancesFromSet(r.Inheritances),
	}
	if req.Inheritances == nil {
		req.Inheritances = []cortex.RelationshipInheritance{}
	}
	return req
}
