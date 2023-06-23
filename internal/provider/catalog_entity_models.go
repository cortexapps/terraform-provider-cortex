package provider

import (
	"encoding/json"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CatalogEntityResourceModel describes the resource data model.
type CatalogEntityResourceModel struct {
	Id          types.String                      `tfsdk:"id"`
	Tag         types.String                      `tfsdk:"tag"`
	Name        types.String                      `tfsdk:"name"`
	Description types.String                      `tfsdk:"description"`
	Owners      []CatalogEntityOwnerResourceModel `tfsdk:"owners"`
	Groups      []types.String                    `tfsdk:"groups"`
	Links       []CatalogEntityLinkResourceModel  `tfsdk:"links"`
	Metadata    types.String                      `tfsdk:"metadata"`
}

func (o CatalogEntityResourceModel) ToApiModel() cortex.CatalogEntityData {
	owners := make([]cortex.CatalogEntityOwner, len(o.Owners))
	for i, owner := range o.Owners {
		owners[i] = owner.ToApiModel()
	}
	groups := make([]string, len(o.Groups))
	for i, group := range o.Groups {
		groups[i] = group.ValueString()
	}
	links := make([]cortex.CatalogEntityLink, len(o.Links))
	for i, link := range o.Links {
		links[i] = link.ToApiModel()
	}
	metadata := make(map[string]interface{})
	err := json.Unmarshal([]byte(o.Metadata.ValueString()), &metadata)
	if err != nil {
		fmt.Println(err)
		metadata = make(map[string]interface{})
	}

	return cortex.CatalogEntityData{
		Tag:         o.Tag.ValueString(),
		Title:       o.Name.ValueString(),
		Description: o.Description.ValueString(),
		Owners:      owners,
		Groups:      groups,
		Links:       links,
		Metadata:    metadata,
	}
}

// CatalogEntityOwnerResourceModel describes owners of the catalog entity. This can be a user, Slack channel, or group.
type CatalogEntityOwnerResourceModel struct {
	Type                 types.String `tfsdk:"type"` // group, user, slack
	Name                 types.String `tfsdk:"name"` // Must be of form <org>/<team>
	Description          types.String `tfsdk:"description"`
	Provider             types.String `tfsdk:"provider"`
	Email                types.String `tfsdk:"email"`
	Channel              types.String `tfsdk:"channel"` // for slack, do not add # to beginning
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

func (o CatalogEntityOwnerResourceModel) ToApiModel() cortex.CatalogEntityOwner {
	return cortex.CatalogEntityOwner{
		Type:                 o.Type.ValueString(),
		Name:                 o.Name.ValueString(),
		Email:                o.Email.ValueString(),
		Description:          o.Description.ValueString(),
		Provider:             o.Provider.ValueString(),
		Channel:              o.Channel.ValueString(),
		NotificationsEnabled: o.NotificationsEnabled.ValueBool(),
	}
}

type CatalogEntityLinkResourceModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	Url  types.String `tfsdk:"url"`
}

func (o CatalogEntityLinkResourceModel) ToApiModel() cortex.CatalogEntityLink {
	return cortex.CatalogEntityLink{
		Type: o.Type.ValueString(),
		Name: o.Name.ValueString(),
		Url:  o.Url.ValueString(),
	}
}
