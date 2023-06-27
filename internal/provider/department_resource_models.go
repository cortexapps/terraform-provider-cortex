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
func (o *DepartmentResourceModel) ToCreateRequest() cortex.CreateDepartmentRequest {
	var members []cortex.DepartmentMember
	for _, member := range o.Members {
		members = append(members, member.ToApiModel())
	}

	return cortex.CreateDepartmentRequest{
		Tag:         o.Tag.ValueString(),
		Name:        o.Name.ValueString(),
		Description: o.Description.ValueString(),
		Members:     members,
	}
}

// ToUpdateRequest https://docs.cortex.io/docs/api/update-department
func (o *DepartmentResourceModel) ToUpdateRequest() cortex.UpdateDepartmentRequest {
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

func (o *DepartmentResourceModel) FromApiModel(department *cortex.Department) {
	o.Id = types.StringValue(department.Tag)
	o.Tag = types.StringValue(department.Tag)
	o.Name = types.StringValue(department.Name)
	o.Description = types.StringValue(department.Description)
	if department.Members != nil {
		o.Members = make([]DepartmentMemberResourceModel, len(department.Members))
		for i, member := range department.Members {
			m := DepartmentMemberResourceModel{}
			o.Members[i] = m.FromApiModel(&member)
		}
	}
}

type DepartmentMemberResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
}

func (o *DepartmentMemberResourceModel) ToApiModel() cortex.DepartmentMember {
	return cortex.DepartmentMember{
		Name:        o.Name.ValueString(),
		Email:       o.Email.ValueString(),
		Description: o.Description.ValueString(),
	}
}

func (o *DepartmentMemberResourceModel) FromApiModel(member *cortex.DepartmentMember) DepartmentMemberResourceModel {
	return DepartmentMemberResourceModel{
		Name:        types.StringValue(member.Name),
		Email:       types.StringValue(member.Email),
		Description: types.StringValue(member.Description),
	}
}

// DepartmentDataSourceModel describes the data source data model.
type DepartmentDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Tag         types.String `tfsdk:"tag"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (o *DepartmentDataSourceModel) FromApiModel(department *cortex.Department) {
	o.Id = types.StringValue(department.Tag)
	o.Tag = types.StringValue(department.Tag)
	o.Name = types.StringValue(department.Name)
	o.Description = types.StringValue(department.Description)
}
