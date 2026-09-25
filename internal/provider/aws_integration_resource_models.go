package provider

import (
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AwsIntegrationResourceModel describes the AWS integration configuration data model within Terraform.
type AwsIntegrationResourceModel struct {
	Id           types.String `tfsdk:"id"`
	AccountId    types.String `tfsdk:"account_id"`
	IamRole      types.String `tfsdk:"iam_role"`
	AccountAlias types.String `tfsdk:"account_alias"`
}

func (r *AwsIntegrationResourceModel) FromApiModel(entity *cortex.AwsConfiguration) {
	r.Id = types.StringValue(entity.AccountID)
	r.AccountId = types.StringValue(entity.AccountID)
	r.IamRole = types.StringValue(entity.IAMRole)
	if entity.AccountAlias != "" {
		r.AccountAlias = types.StringValue(entity.AccountAlias)
	} else {
		r.AccountAlias = types.StringNull()
	}
}

func (r *AwsIntegrationResourceModel) ToApiModel() cortex.AwsConfiguration {
	return cortex.AwsConfiguration{
		AccountID: r.AccountId.ValueString(),
		IAMRole:   r.IamRole.ValueString(),
	}
}
