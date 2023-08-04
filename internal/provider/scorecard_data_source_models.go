package provider

import (
	"context"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Models
 **********************************************************************************************************************/

// ScorecardDataSourceModel describes the data source data model.
type ScorecardDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Tag         types.String `tfsdk:"tag"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Draft       types.Bool   `tfsdk:"draft"`
}

func (r *ScorecardDataSourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.Scorecard) {
	r.Id = types.StringValue(entity.Tag)
	r.Tag = types.StringValue(entity.Tag)
	r.Name = types.StringValue(entity.Name)
	r.Description = types.StringValue(entity.Description)
	r.Draft = types.BoolValue(entity.Draft)
}
