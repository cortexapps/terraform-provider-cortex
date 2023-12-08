package provider

import (
	"context"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCatalogEntityCustomDataDataSourceModel() CatalogEntityCustomDataDataSourceModel {
	return CatalogEntityCustomDataDataSourceModel{}
}

// CatalogEntityCustomDataDataSourceModel describes the data source data model.
type CatalogEntityCustomDataDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Tag         types.String `tfsdk:"tag"`
	Key         types.String `tfsdk:"key"`
	Description types.String `tfsdk:"description"`
	Source      types.String `tfsdk:"source"`
	Value       types.String `tfsdk:"value"`
}

func (o *CatalogEntityCustomDataDataSourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.CatalogEntityCustomData) {
	o.Id = types.StringValue(entity.Tag)
	o.Tag = types.StringValue(entity.Tag)
	o.Key = types.StringValue(entity.Key)
	o.Description = types.StringValue(entity.Description)
	o.Source = types.StringValue(entity.Source)
	value, err := entity.ValueAsString()
	if err != nil {
		diagnostics.AddError("Error parsing value: %s", err.Error())
		return
	}
	o.Value = types.StringValue(value)
}
