package provider

import (
	"context"
	"encoding/json"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Models
 **********************************************************************************************************************/

// ResourceDefinitionResourceModel describes the department data model within Terraform.
type ResourceDefinitionResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Source      types.String `tfsdk:"source"`
	Schema      types.String `tfsdk:"schema"`
}

func (r *ResourceDefinitionResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.ResourceDefinition) {
	r.Id = types.StringValue(entity.Type)
	r.Type = types.StringValue(entity.Type)
	r.Name = types.StringValue(entity.Name)
	r.Source = types.StringValue(entity.Source)
	if entity.Description != "" {
		r.Description = types.StringValue(entity.Description)
	} else {
		r.Description = types.StringNull()
	}

	schema := make(map[string]interface{})
	if entity.Schema != nil && len(entity.Schema) > 0 {
		schema = entity.Schema
	}
	sv, err := json.Marshal(schema)
	if err != nil {
		diagnostics.AddError("Error parsing schema: %s", err.Error())
		return
	}
	r.Schema = types.StringValue(string(sv))
}

func (r *ResourceDefinitionResourceModel) ToApiModel(diagnostics *diag.Diagnostics) cortex.ResourceDefinition {
	entity := cortex.ResourceDefinition{
		Type:        r.Type.ValueString(),
		Name:        r.Name.ValueString(),
		Description: r.Description.ValueString(),
		Source:      r.Source.ValueString(),
	}

	schema := make(map[string]interface{})
	if !r.Schema.IsNull() && !r.Schema.IsUnknown() && r.Schema.ValueString() != "" {
		err := json.Unmarshal([]byte(r.Schema.ValueString()), &schema)
		if err != nil {
			diagnostics.AddError(
				"Unable to Convert Resource Definition Schema",
				"An unexpected result occurred when deserializing the schema from JSON. Please ensure your resource definition schema is valid JSON.",
			)
			schema = make(map[string]interface{})
		}
	}
	entity.Schema = schema

	return entity
}
