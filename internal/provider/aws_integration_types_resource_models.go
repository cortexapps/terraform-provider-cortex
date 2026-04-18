package provider

import (
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AwsTypeResourceModel struct {
	Type       types.String `tfsdk:"type"`
	Configured types.Bool   `tfsdk:"configured"`
}

type AwsIntegrationTypesResourceModel struct {
	Id    types.String           `tfsdk:"id"`
	Types []AwsTypeResourceModel `tfsdk:"types"`
}

func (r *AwsIntegrationTypesResourceModel) FromApiModel(entity []cortex.AwsType) {
	r.Id = types.StringValue("aws_types")
	if entity != nil {
		r.Types = make([]AwsTypeResourceModel, len(entity))
		for i, t := range entity {
			r.Types[i] = AwsTypeResourceModel{
				Type:       types.StringValue(t.Type),
				Configured: types.BoolValue(t.Configured),
			}
		}
	} else {
		r.Types = nil
	}
}

func (r *AwsIntegrationTypesResourceModel) ToApiModel() []cortex.AwsType {
	var typesList []cortex.AwsType
	for _, t := range r.Types {
		typesList = append(typesList, cortex.AwsType{
			Type:       t.Type.ValueString(),
			Configured: t.Configured.ValueBool(),
		})
	}
	return typesList
}
