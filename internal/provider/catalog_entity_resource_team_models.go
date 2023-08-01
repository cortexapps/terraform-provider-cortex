package provider

import (
	"context"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Team Catalog Entity Models
 **********************************************************************************************************************/

// Team

type CatalogEntityTeamResourceModel struct {
	Members []CatalogEntityTeamMemberResourceModel  `tfsdk:"members"`
	Groups  []CatalogEntityGroupMemberResourceModel `tfsdk:"groups"`
}

func (o *CatalogEntityTeamResourceModel) AttrTypes() map[string]attr.Type {
	tm := CatalogEntityTeamMemberResourceModel{}
	gm := CatalogEntityGroupMemberResourceModel{}
	return map[string]attr.Type{
		"members": types.ListType{ElemType: types.ObjectType{AttrTypes: tm.AttrTypes()}},
		"groups":  types.ListType{ElemType: types.ObjectType{AttrTypes: gm.AttrTypes()}},
	}
}

func (o *CatalogEntityTeamResourceModel) ToApiModel() cortex.CatalogEntityTeam {
	var members = make([]cortex.CatalogEntityTeamMember, len(o.Members))
	for i, e := range o.Members {
		members[i] = e.ToApiModel()
	}
	var groups = make([]cortex.CatalogEntityGroupMember, len(o.Groups))
	for i, e := range o.Groups {
		groups[i] = e.ToApiModel()
	}
	return cortex.CatalogEntityTeam{
		Members: members,
		Groups:  groups,
	}
}

func (o *CatalogEntityTeamResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityTeam) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	tm := CatalogEntityTeamResourceModel{
		Members: make([]CatalogEntityTeamMemberResourceModel, len(entity.Members)),
		Groups:  make([]CatalogEntityGroupMemberResourceModel, len(entity.Groups)),
	}

	if len(entity.Members) > 0 {
		for i, e := range entity.Members {
			ob := CatalogEntityTeamMemberResourceModel{}
			if e.Enabled() {
				tm.Members[i] = ob.FromApiModel(&e)
			}
		}
	} else {
		tm.Members = nil
	}
	if len(entity.Groups) > 0 {
		for i, e := range entity.Groups {
			ob := CatalogEntityGroupMemberResourceModel{}
			if e.Enabled() {
				tm.Groups[i] = ob.FromApiModel(&e)
			}
		}
	} else {
		tm.Groups = nil
	}
	objectValue, d := types.ObjectValueFrom(ctx, tm.AttrTypes(), &tm)
	diagnostics.Append(d...)
	return objectValue
}

type CatalogEntityTeamMemberResourceModel struct {
	Name                 types.String `tfsdk:"name"`
	Email                types.String `tfsdk:"email"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

func (o *CatalogEntityTeamMemberResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                  types.StringType,
		"email":                 types.StringType,
		"notifications_enabled": types.BoolType,
	}
}

func (o *CatalogEntityTeamMemberResourceModel) ToApiModel() cortex.CatalogEntityTeamMember {
	return cortex.CatalogEntityTeamMember{
		Name:                 o.Name.ValueString(),
		Email:                o.Email.ValueString(),
		NotificationsEnabled: o.NotificationsEnabled.ValueBool(),
	}
}

func (o *CatalogEntityTeamMemberResourceModel) FromApiModel(entity *cortex.CatalogEntityTeamMember) CatalogEntityTeamMemberResourceModel {
	notificationsEnabled := types.BoolNull()
	if entity.NotificationsEnabled {
		notificationsEnabled = types.BoolValue(true)
	}

	return CatalogEntityTeamMemberResourceModel{
		Name:                 types.StringValue(entity.Name),
		Email:                types.StringValue(entity.Email),
		NotificationsEnabled: notificationsEnabled,
	}
}

type CatalogEntityGroupMemberResourceModel struct {
	Name     types.String `tfsdk:"name"`
	Provider types.String `tfsdk:"provider"`
}

func (o *CatalogEntityGroupMemberResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":     types.StringType,
		"provider": types.StringType,
	}
}

func (o *CatalogEntityGroupMemberResourceModel) ToApiModel() cortex.CatalogEntityGroupMember {
	return cortex.CatalogEntityGroupMember{
		Name:     o.Name.ValueString(),
		Provider: o.Provider.ValueString(),
	}
}

func (o *CatalogEntityGroupMemberResourceModel) FromApiModel(entity *cortex.CatalogEntityGroupMember) CatalogEntityGroupMemberResourceModel {
	return CatalogEntityGroupMemberResourceModel{
		Name:     types.StringValue(entity.Name),
		Provider: types.StringValue(entity.Provider),
	}
}
