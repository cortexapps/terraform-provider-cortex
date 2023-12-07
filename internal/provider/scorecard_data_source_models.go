package provider

import (
	"context"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Models
 **********************************************************************************************************************/

// ScorecardDataSourceModel describes the data source data model.
type ScorecardDataSourceModel struct {
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

func (o *ScorecardDataSourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.Scorecard) {
	o.Id = types.StringValue(entity.Tag)
	o.Tag = types.StringValue(entity.Tag)
	o.Name = types.StringValue(entity.Name)
	o.Description = types.StringValue(entity.Description)
	o.Draft = types.BoolValue(entity.Draft)

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
