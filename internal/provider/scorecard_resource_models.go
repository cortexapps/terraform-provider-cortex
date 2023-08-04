package provider

import (
	"context"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
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
	Category types.String `tfsdk:"category"`
	Query    types.String `tfsdk:"query"`
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

func (o *ScorecardResourceModel) ToApiModel(ctx context.Context) cortex.Scorecard {
	defaultObjOptions := getDefaultObjectOptions()

	rules := make([]cortex.ScorecardRule, len(o.Rules))
	for i, rule := range o.Rules {
		rules[i] = rule.ToApiModel()
	}

	ladder := ScorecardLadderResourceModel{}
	err := o.Ladder.As(ctx, &ladder, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing scorecard ladder: ", err)
	}

	filter := ScorecardFilterResourceModel{}
	err = o.Filter.As(ctx, &filter, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing scorecard filter: ", err)
	}

	evaluation := ScorecardEvaluationResourceModel{}
	err = o.Evaluation.As(ctx, &evaluation, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing scorecard evaluation: ", err)
	}

	return cortex.Scorecard{
		Tag:         o.Tag.ValueString(),
		Name:        o.Name.ValueString(),
		Description: o.Description.ValueString(),
		Draft:       o.Draft.ValueBool(),
		Ladder:      ladder.ToApiModel(),
		Rules:       rules,
		Filter:      filter.ToApiModel(),
		Evaluation:  evaluation.ToApiModel(),
	}
}

func (o *ScorecardResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.Scorecard) {
	o.Id = types.StringValue(entity.Tag)
	o.Tag = types.StringValue(entity.Tag)
}

/***********************************************************************************************************************
 * Ladder/Levels
 **********************************************************************************************************************/

func (o *ScorecardLadderResourceModel) ToApiModel() cortex.ScorecardLadder {
	levels := make([]cortex.ScorecardLevel, len(o.Levels))
	for i, level := range o.Levels {
		levels[i] = level.ToApiModel()
	}
	return cortex.ScorecardLadder{
		Levels: levels,
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

/***********************************************************************************************************************
 * Filter
 **********************************************************************************************************************/

func (o *ScorecardFilterResourceModel) ToApiModel() cortex.ScorecardFilter {
	return cortex.ScorecardFilter{
		Category: o.Category.ValueString(),
		Query:    o.Query.ValueString(),
	}
}

/***********************************************************************************************************************
 * Evaluation
 **********************************************************************************************************************/

func (o *ScorecardEvaluationResourceModel) ToApiModel() cortex.ScorecardEvaluation {
	return cortex.ScorecardEvaluation{
		Window: o.Window.ValueInt64(),
	}
}
