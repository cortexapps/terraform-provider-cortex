package provider

import (
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DepartmentResource defines the resource implementation.
type DepartmentResource struct {
	client *cortex.HttpClient
}

// DepartmentResourceModel describes the department data model within Terraform.
type DepartmentResourceModel struct {
	Id          types.String                    `tfsdk:"id"`
	Tag         types.String                    `tfsdk:"tag"`
	Name        types.String                    `tfsdk:"name"`
	Description types.String                    `tfsdk:"description"`
	Members     []DepartmentMemberResourceModel `tfsdk:"members"`
}

// ToCreateRequest https://docs.cortex.io/docs/api/create-department
func (o DepartmentResourceModel) ToCreateRequest() cortex.CreateDepartmentRequest {
	var members []cortex.DepartmentMember
	for _, member := range o.Members {
		members = append(members, member.ToApiModel())
	}

	return cortex.CreateDepartmentRequest{
		Tag:         o.Tag.ValueString(),
		Name:        o.Name.ValueString(),
		Description: o.Description.String(),
		Members:     members,
	}
}

// ToUpdateRequest https://docs.cortex.io/docs/api/update-department
func (o DepartmentResourceModel) ToUpdateRequest() cortex.UpdateDepartmentRequest {
	var members []cortex.DepartmentMember
	for _, member := range o.Members {
		members = append(members, member.ToApiModel())
	}
	return cortex.UpdateDepartmentRequest{
		Name:        o.Name.ValueString(),
		Description: o.Description.ValueString(),
		Members:     members,
	}
}

type DepartmentMemberResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
}

func (o DepartmentMemberResourceModel) ToApiModel() cortex.DepartmentMember {
	return cortex.DepartmentMember{
		Name:        o.Name.ValueString(),
		Email:       o.Email.ValueString(),
		Description: o.Description.ValueString(),
	}
}
