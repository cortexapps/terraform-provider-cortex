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

type CatalogEntityCustomDataResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Tag         types.String `tfsdk:"tag"`
	Key         types.String `tfsdk:"key"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
}

func (r *CatalogEntityCustomDataResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.CatalogEntityCustomData) {
	r.Id = types.StringValue(entity.ID())
	r.Tag = types.StringValue(entity.Tag)
	r.Key = types.StringValue(entity.Key)
	if entity.Description != "" {
		r.Description = types.StringValue(entity.Description)
	} else {
		r.Description = types.StringNull()
	}

	value, err := entity.ValueAsString()
	if err != nil {
		diagnostics.AddError("Error parsing value: %s", err.Error())
		return
	}
	r.Value = types.StringValue(value)
}

func (r *CatalogEntityCustomDataResourceModel) ToApiModel(diagnostics *diag.Diagnostics) cortex.CatalogEntityCustomData {
	entity := cortex.CatalogEntityCustomData{
		Tag:         r.Tag.ValueString(),
		Key:         r.Key.ValueString(),
		Description: r.Description.ValueString(),
	}

	var value interface{}
	if !r.Value.IsNull() && !r.Value.IsUnknown() && r.Value.ValueString() != "" {
		err := error(nil)
		value, err = cortex.StringToInterface(r.Value.ValueString())
		if err != nil {
			diagnostics.AddError(
				"Unable to Convert Custom Data Value",
				"An unexpected result occurred when deserializing the value from JSON. Please ensure your value is valid JSON.",
			)
			value = ""
		}
	} else {
		value = ""
	}
	entity.Value = value

	return entity
}
