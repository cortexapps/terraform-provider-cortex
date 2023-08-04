package provider

import (
	"context"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewResourceDefinitionDataSourceModel() ResourceDefinitionDataSourceModel {
	return ResourceDefinitionDataSourceModel{}
}

// ResourceDefinitionDataSourceModel describes the data source data model.
type ResourceDefinitionDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Source      types.String `tfsdk:"source"`
	Schema      types.String `tfsdk:"schema"`
}

func (r *ResourceDefinitionDataSourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.ResourceDefinition) {
	r.Id = types.StringValue(entity.Type)
	r.Type = types.StringValue(entity.Type)
	r.Name = types.StringValue(entity.Name)
	r.Description = types.StringValue(entity.Description)
	r.Source = types.StringValue(entity.Source)

	schema, err := entity.SchemaAsString()
	if err != nil {
		diagnostics.AddError("Error parsing schema: %s", err.Error())
		return
	}
	r.Schema = types.StringValue(schema)
}
