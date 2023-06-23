package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// CatalogEntityResourceModel describes the resource data model.
type CatalogEntityResourceModel struct {
	Id           types.String                           `tfsdk:"id"`
	Tag          types.String                           `tfsdk:"tag"`
	Name         types.String                           `tfsdk:"name"`
	Description  types.String                           `tfsdk:"description"`
	Owners       []CatalogEntityOwnerResourceModel      `tfsdk:"owners"`
	Groups       []types.String                         `tfsdk:"groups"`
	Links        []CatalogEntityLinkResourceModel       `tfsdk:"links"`
	Metadata     types.String                           `tfsdk:"metadata"`
	Dependencies []CatalogEntityDependencyResourceModel `tfsdk:"dependencies"`
	Alerts       []CatalogEntityAlertResourceModel      `tfsdk:"alerts"`
	Git          types.Object                           `tfsdk:"git"`
}

func (o CatalogEntityResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityData {
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
	if !o.Metadata.IsNull() && !o.Metadata.IsUnknown() && o.Metadata.ValueString() != "" {
		err := json.Unmarshal([]byte(o.Metadata.ValueString()), &metadata)
		if err != nil {
			fmt.Println(err)
			metadata = make(map[string]interface{})
		}
	}
	dependencies := make([]cortex.CatalogEntityDependency, len(o.Dependencies))
	for i, dependency := range o.Dependencies {
		dependencies[i] = dependency.ToApiModel()
	}

	alerts := make([]cortex.CatalogEntityAlert, len(o.Alerts))
	for i, alert := range o.Alerts {
		alerts[i] = alert.ToApiModel()
	}

	git := &CatalogEntityGitResourceModel{}
	err := o.Git.As(ctx, git, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	if err != nil {
		fmt.Println(err)
	}

	return cortex.CatalogEntityData{
		Tag:          o.Tag.ValueString(),
		Title:        o.Name.ValueString(),
		Description:  o.Description.ValueString(),
		Owners:       owners,
		Groups:       groups,
		Links:        links,
		Metadata:     metadata,
		Dependencies: dependencies,
		Alerts:       alerts,
		Git:          git.ToApiModel(),
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

type CatalogEntityDependencyResourceModel struct {
	Tag         types.String `tfsdk:"tag"`
	Method      types.String `tfsdk:"method"`
	Path        types.String `tfsdk:"path"`
	Description types.String `tfsdk:"description"`
	Metadata    types.String `tfsdk:"metadata"`
}

func (o CatalogEntityDependencyResourceModel) ToApiModel() cortex.CatalogEntityDependency {
	metadata := make(map[string]interface{})
	if !o.Metadata.IsNull() && !o.Metadata.IsUnknown() && o.Metadata.ValueString() != "" {
		err := json.Unmarshal([]byte(o.Metadata.ValueString()), &metadata)
		if err != nil {
			fmt.Println(err)
			metadata = make(map[string]interface{})
		}
	}

	return cortex.CatalogEntityDependency{
		Tag:         o.Tag.ValueString(),
		Method:      o.Method.ValueString(),
		Path:        o.Path.ValueString(),
		Description: o.Description.ValueString(),
		Metadata:    metadata,
	}
}

type CatalogEntityAlertResourceModel struct {
	Type  types.String `tfsdk:"type"`
	Tag   types.String `tfsdk:"tag"`
	Value types.String `tfsdk:"value"`
}

func (o CatalogEntityAlertResourceModel) ToApiModel() cortex.CatalogEntityAlert {
	return cortex.CatalogEntityAlert{
		Type:  o.Type.ValueString(),
		Tag:   o.Tag.ValueString(),
		Value: o.Value.ValueString(),
	}
}

/***********************************************************************************************************************
 * Git
 ***********************************************************************************************************************/

type CatalogEntityGitResourceModel struct {
	Github    CatalogEntityGithubResourceModel    `tfsdk:"github"`
	Gitlab    CatalogEntityGitlabResourceModel    `tfsdk:"gitlab"`
	Azure     CatalogEntityAzureResourceModel     `tfsdk:"azure"`
	Bitbucket CatalogEntityBitbucketResourceModel `tfsdk:"bitbucket"`
}

func (o CatalogEntityGitResourceModel) ToApiModel() cortex.CatalogEntityGit {
	git := cortex.CatalogEntityGit{}
	if o.Github.Repository.ValueString() != "" {
		git.Github = o.Github.ToApiModel()
	}
	if o.Gitlab.Repository.ValueString() != "" {
		git.Gitlab = o.Gitlab.ToApiModel()
	}
	if o.Azure.Repository.ValueString() != "" {
		git.Azure = o.Azure.ToApiModel()
	}
	if o.Bitbucket.Repository.ValueString() != "" {
		git.BitBucket = o.Bitbucket.ToApiModel()
	}
	return git
}

type CatalogEntityGithubResourceModel struct {
	Repository types.String `tfsdk:"repository"`
	BasePath   types.String `tfsdk:"base_path"`
}

func (o CatalogEntityGithubResourceModel) ToApiModel() cortex.CatalogEntityGitGithub {
	return cortex.CatalogEntityGitGithub{
		Repository: o.Repository.ValueString(),
		BasePath:   o.BasePath.ValueString(),
	}
}

type CatalogEntityGitlabResourceModel struct {
	Repository types.String `tfsdk:"repository"`
	BasePath   types.String `tfsdk:"base_path"`
}

func (o CatalogEntityGitlabResourceModel) ToApiModel() cortex.CatalogEntityGitGitlab {
	return cortex.CatalogEntityGitGitlab{
		Repository: o.Repository.ValueString(),
		BasePath:   o.BasePath.ValueString(),
	}
}

type CatalogEntityAzureResourceModel struct {
	Project    types.String `tfsdk:"project"`
	Repository types.String `tfsdk:"repository"`
	BasePath   types.String `tfsdk:"base_path"`
}

func (o CatalogEntityAzureResourceModel) ToApiModel() cortex.CatalogEntityGitAzureDevOps {
	return cortex.CatalogEntityGitAzureDevOps{
		Project:    o.Project.ValueString(),
		Repository: o.Repository.ValueString(),
		BasePath:   o.BasePath.ValueString(),
	}
}

type CatalogEntityBitbucketResourceModel struct {
	Repository types.String `tfsdk:"repository"`
}

func (o CatalogEntityBitbucketResourceModel) ToApiModel() cortex.CatalogEntityGitBitBucket {
	return cortex.CatalogEntityGitBitBucket{
		Repository: o.Repository.ValueString(),
	}
}
