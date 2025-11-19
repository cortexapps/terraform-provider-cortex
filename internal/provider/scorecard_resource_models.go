package provider

import (
	"context"
	"fmt"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// ScorecardResourceModel describes the scorecard data model within Terraform.
type ScorecardResourceModel struct {
	Id          types.String                 `tfsdk:"id"`
	Tag         types.String                 `tfsdk:"tag"`
	Name        types.String                 `tfsdk:"name"`
	Description types.String                 `tfsdk:"description"`
	Draft       types.Bool                   `tfsdk:"draft"`
	Ladder      types.Object                 `tfsdk:"ladder"`
	Rules       []ScorecardRuleResourceModel `tfsdk:"rules"`
	Filter      types.Object                 `tfsdk:"filter"`
	Evaluation  types.Object                 `tfsdk:"evaluation"`
}

type ScorecardLadderResourceModel struct {
	Levels []ScorecardLevelResourceModel `tfsdk:"levels"`
}

type ScorecardLevelResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Rank        types.Int64  `tfsdk:"rank"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
}

type ScorecardRuleResourceModel struct {
	Title          types.String `tfsdk:"title"`
	Expression     types.String `tfsdk:"expression"`
	Weight         types.Int64  `tfsdk:"weight"`
	Level          types.String `tfsdk:"level"`
	Description    types.String `tfsdk:"description"`
	FailureMessage types.String `tfsdk:"failure_message"`
}

type ScorecardFilterResourceModel struct {
	Types  types.Object `tfsdk:"types"`
	Groups types.Object `tfsdk:"groups"`
	Query  types.String `tfsdk:"query"`
}

type ScorecardFilterTypesResourceModel struct {
	Include []types.String `tfsdk:"include"`
	Exclude []types.String `tfsdk:"exclude"`
}

type ScorecardFilterGroupsResourceModel struct {
	Include []types.String `tfsdk:"include"`
	Exclude []types.String `tfsdk:"exclude"`
}

type ScorecardEvaluationResourceModel struct {
	Window types.Int64 `tfsdk:"window"`
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func NewScorecardResourceModel() ScorecardResourceModel {
	return ScorecardResourceModel{}
}

/***********************************************************************************************************************
 * ScorecardResourceModel
 **********************************************************************************************************************/

func (o *ScorecardResourceModel) ToApiModel(ctx context.Context, diagnostics *diag.Diagnostics) cortex.Scorecard {
	defaultObjOptions := getDefaultObjectOptions()

	rules := make([]cortex.ScorecardRule, len(o.Rules))
	for i, rule := range o.Rules {
		rules[i] = rule.ToApiModel()
	}

	ladder := ScorecardLadderResourceModel{}
	err := o.Ladder.As(ctx, &ladder, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing scorecard ladder", fmt.Sprintf("%+v", err))
	}

	filter := ScorecardFilterResourceModel{}
	err = o.Filter.As(ctx, &filter, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing scorecard filter", fmt.Sprintf("%+v", err))
	}

	evaluation := ScorecardEvaluationResourceModel{}
	err = o.Evaluation.As(ctx, &evaluation, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing scorecard evaluation", fmt.Sprintf("%+v", err))
	}

	return cortex.Scorecard{
		Tag:         o.Tag.ValueString(),
		Name:        o.Name.ValueString(),
		Description: o.Description.ValueString(),
		Draft:       o.Draft.ValueBool(),
		Ladder:      ladder.ToApiModel(),
		Rules:       rules,
		Filter:      filter.ToApiModel(ctx, diagnostics),
		Evaluation:  evaluation.ToApiModel(),
	}
}

func (o *ScorecardResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.Scorecard) {
	o.Id = types.StringValue(entity.Tag)
	o.Tag = types.StringValue(entity.Tag)
	o.Name = types.StringValue(entity.Name)
	o.Draft = types.BoolValue(entity.Draft)
	o.Description = types.StringValue(entity.Description)

	ladder := ScorecardLadderResourceModel{}
	o.Ladder = ladder.FromApiModel(ctx, diagnostics, &entity.Ladder)

	rules := make([]ScorecardRuleResourceModel, len(entity.Rules))
	for i, e := range entity.Rules {
		rrm := ScorecardRuleResourceModel{}
		rules[i] = rrm.FromApiModel(&e)
	}
	o.Rules = rules

	filter := ScorecardFilterResourceModel{}
	o.Filter = filter.FromApiModel(ctx, diagnostics, &entity.Filter)

	evaluation := ScorecardEvaluationResourceModel{}
	o.Evaluation = evaluation.FromApiModel(ctx, diagnostics, &entity.Evaluation)
}

/***********************************************************************************************************************
 * Ladder/Levels
 **********************************************************************************************************************/

func (o *ScorecardLadderResourceModel) AttrTypes() map[string]attr.Type {
	sl := ScorecardLevelResourceModel{}
	return map[string]attr.Type{
		"levels": types.ListType{ElemType: types.ObjectType{AttrTypes: sl.AttrTypes()}},
	}
}

func (o *ScorecardLadderResourceModel) ToApiModel() cortex.ScorecardLadder {
	levels := make([]cortex.ScorecardLevel, len(o.Levels))
	for i, level := range o.Levels {
		levels[i] = level.ToApiModel()
	}
	return cortex.ScorecardLadder{
		Levels: levels,
	}
}

func (o *ScorecardLadderResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.ScorecardLadder) types.Object {
	obj := ScorecardLadderResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(obj.AttrTypes())
	}

	levels := make([]ScorecardLevelResourceModel, len(entity.Levels))
	for i, e := range entity.Levels {
		lrm := ScorecardLevelResourceModel{}
		levels[i] = lrm.FromApiModel(&e)
	}
	obj.Levels = levels

	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

func (o *ScorecardLevelResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"rank":        types.Int64Type,
		"description": types.StringType,
		"color":       types.StringType,
	}
}

func (o *ScorecardLevelResourceModel) ToApiModel() cortex.ScorecardLevel {
	return cortex.ScorecardLevel{
		Name:        o.Name.ValueString(),
		Rank:        o.Rank.ValueInt64(),
		Description: o.Description.ValueString(),
		Color:       o.Color.ValueString(),
	}
}

func (o *ScorecardLevelResourceModel) FromApiModel(entity *cortex.ScorecardLevel) ScorecardLevelResourceModel {
	lrm := ScorecardLevelResourceModel{
		Name:  types.StringValue(entity.Name),
		Rank:  types.Int64Value(entity.Rank),
		Color: types.StringValue(entity.Color),
	}
	if entity.Description != "" {
		lrm.Description = types.StringValue(entity.Description)
	} else {
		lrm.Description = types.StringNull()
	}
	return lrm
}

/***********************************************************************************************************************
 * Rules
 **********************************************************************************************************************/

func (o *ScorecardRuleResourceModel) ToApiModel() cortex.ScorecardRule {
	return cortex.ScorecardRule{
		Title:          o.Title.ValueString(),
		Expression:     o.Expression.ValueString(),
		Description:    o.Description.ValueString(),
		Weight:         o.Weight.ValueInt64(),
		FailureMessage: o.FailureMessage.ValueString(),
		Level:          o.Level.ValueString(),
	}
}

func (o *ScorecardRuleResourceModel) FromApiModel(entity *cortex.ScorecardRule) ScorecardRuleResourceModel {
	rm := ScorecardRuleResourceModel{
		Title:      types.StringValue(entity.Title),
		Expression: types.StringValue(entity.Expression),
		Weight:     types.Int64Value(entity.Weight),
		Level:      types.StringValue(entity.Level),
	}
	if entity.Description != "" {
		rm.Description = types.StringValue(entity.Description)
	} else {
		rm.Description = types.StringNull()
	}
	if entity.FailureMessage != "" {
		rm.FailureMessage = types.StringValue(entity.FailureMessage)
	} else {
		rm.FailureMessage = types.StringNull()
	}
	return rm
}

/***********************************************************************************************************************
 * Filter
 **********************************************************************************************************************/

func (o *ScorecardFilterResourceModel) AttrTypes() map[string]attr.Type {
	typesModel := ScorecardFilterTypesResourceModel{}
	groupsModel := ScorecardFilterGroupsResourceModel{}
	return map[string]attr.Type{
		"types":  types.ObjectType{AttrTypes: typesModel.AttrTypes()},
		"groups": types.ObjectType{AttrTypes: groupsModel.AttrTypes()},
		"query":  types.StringType,
	}
}

func (o *ScorecardFilterTypesResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"include": types.SetType{ElemType: types.StringType},
		"exclude": types.SetType{ElemType: types.StringType},
	}
}

func (o *ScorecardFilterGroupsResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"include": types.SetType{ElemType: types.StringType},
		"exclude": types.SetType{ElemType: types.StringType},
	}
}

func (o *ScorecardFilterResourceModel) ToApiModel(ctx context.Context, diagnostics *diag.Diagnostics) cortex.ScorecardFilter {
	defaultObjOptions := getDefaultObjectOptions()
	filter := cortex.ScorecardFilter{
		Kind:  "GENERIC",
		Query: o.Query.ValueString(),
	}

	// Convert Types if present
	if !o.Types.IsNull() {
		typesModel := ScorecardFilterTypesResourceModel{}
		err := o.Types.As(ctx, &typesModel, defaultObjOptions)
		if err != nil {
			diagnostics.AddError("error parsing filter types", fmt.Sprintf("%+v", err))
		} else {
			include := make([]string, len(typesModel.Include))
			for i, v := range typesModel.Include {
				include[i] = v.ValueString()
			}
			exclude := make([]string, len(typesModel.Exclude))
			for i, v := range typesModel.Exclude {
				exclude[i] = v.ValueString()
			}
			if len(include) > 0 || len(exclude) > 0 {
				filter.Types = &cortex.ScorecardFilterTypes{
					Include: include,
					Exclude: exclude,
				}
			}
		}
	}

	// Convert Groups if present
	if !o.Groups.IsNull() {
		groupsModel := ScorecardFilterGroupsResourceModel{}
		err := o.Groups.As(ctx, &groupsModel, defaultObjOptions)
		if err != nil {
			diagnostics.AddError("error parsing filter groups", fmt.Sprintf("%+v", err))
		} else {
			include := make([]string, len(groupsModel.Include))
			for i, v := range groupsModel.Include {
				include[i] = v.ValueString()
			}
			exclude := make([]string, len(groupsModel.Exclude))
			for i, v := range groupsModel.Exclude {
				exclude[i] = v.ValueString()
			}
			if len(include) > 0 || len(exclude) > 0 {
				filter.Groups = &cortex.ScorecardFilterGroups{
					Include: include,
					Exclude: exclude,
				}
			}
		}
	}

	return filter
}

func (o *ScorecardFilterResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.ScorecardFilter) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	obj := ScorecardFilterResourceModel{}

	// Convert Query
	if entity.Query != "" {
		obj.Query = types.StringValue(entity.Query)
	} else {
		obj.Query = types.StringNull()
	}

	// Convert Types
	if entity.Types != nil {
		typesModel := ScorecardFilterTypesResourceModel{}

		// Convert Include
		if len(entity.Types.Include) > 0 {
			typesModel.Include = make([]types.String, len(entity.Types.Include))
			for i, v := range entity.Types.Include {
				typesModel.Include[i] = types.StringValue(v)
			}
		}

		// Convert Exclude
		if len(entity.Types.Exclude) > 0 {
			typesModel.Exclude = make([]types.String, len(entity.Types.Exclude))
			for i, v := range entity.Types.Exclude {
				typesModel.Exclude[i] = types.StringValue(v)
			}
		}

		typesObj, d := types.ObjectValueFrom(ctx, typesModel.AttrTypes(), &typesModel)
		diagnostics.Append(d...)
		obj.Types = typesObj
	} else {
		typesModel := ScorecardFilterTypesResourceModel{}
		obj.Types = types.ObjectNull(typesModel.AttrTypes())
	}

	// Convert Groups
	if entity.Groups != nil {
		groupsModel := ScorecardFilterGroupsResourceModel{}

		// Convert Include
		if len(entity.Groups.Include) > 0 {
			groupsModel.Include = make([]types.String, len(entity.Groups.Include))
			for i, v := range entity.Groups.Include {
				groupsModel.Include[i] = types.StringValue(v)
			}
		}

		// Convert Exclude
		if len(entity.Groups.Exclude) > 0 {
			groupsModel.Exclude = make([]types.String, len(entity.Groups.Exclude))
			for i, v := range entity.Groups.Exclude {
				groupsModel.Exclude[i] = types.StringValue(v)
			}
		}

		groupsObj, d := types.ObjectValueFrom(ctx, groupsModel.AttrTypes(), &groupsModel)
		diagnostics.Append(d...)
		obj.Groups = groupsObj
	} else {
		groupsModel := ScorecardFilterGroupsResourceModel{}
		obj.Groups = types.ObjectNull(groupsModel.AttrTypes())
	}

	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

/***********************************************************************************************************************
 * Evaluation
 **********************************************************************************************************************/

func (o *ScorecardEvaluationResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"window": types.Int64Type,
	}
}

func (o *ScorecardEvaluationResourceModel) ToApiModel() cortex.ScorecardEvaluation {
	return cortex.ScorecardEvaluation{
		Window: o.Window.ValueInt64(),
	}
}

func (o *ScorecardEvaluationResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.ScorecardEvaluation) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	obj := ScorecardEvaluationResourceModel{}
	if entity.Window > 0 {
		obj.Window = types.Int64Value(entity.Window)
	} else {
		obj.Window = types.Int64Null()
	}

	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}
