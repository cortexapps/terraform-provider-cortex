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
	Issues       types.Object                           `tfsdk:"issues"`
	OnCall       types.Object                           `tfsdk:"on_call"`
	Sentry       types.Object                           `tfsdk:"sentry"`
	Snyk         types.Object                           `tfsdk:"snyk"`
}

func (o CatalogEntityResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityData {
	defaultObjOptions := basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}

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
	err := o.Git.As(ctx, git, defaultObjOptions)
	if err != nil {
		fmt.Println(err)
	}
	issues := &CatalogEntityIssuesResourceModel{}
	err = o.Issues.As(ctx, issues, defaultObjOptions)
	if err != nil {
		fmt.Println(err)
	}
	onCall := &CatalogEntityOnCallResourceModel{}
	err = o.OnCall.As(ctx, onCall, defaultObjOptions)
	if err != nil {
		fmt.Println(err)
	}

	sentry := &CatalogEntitySentryResourceModel{}
	err = o.Sentry.As(ctx, sentry, defaultObjOptions)
	if err != nil {
		fmt.Println(err)
	}
	snyk := &CatalogEntitySnykResourceModel{}
	err = o.Snyk.As(ctx, snyk, defaultObjOptions)
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
		Issues:       issues.ToApiModel(ctx),
		OnCall:       onCall.ToApiModel(),
		Sentry:       sentry.ToApiModel(),
		Snyk:         snyk.ToApiModel(),
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

/***********************************************************************************************************************
 * Issues
 ***********************************************************************************************************************/

type CatalogEntityIssuesResourceModel struct {
	Jira CatalogEntityIssuesJiraResourceModel `tfsdk:"jira"`
}

func (o CatalogEntityIssuesResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityIssues {
	return cortex.CatalogEntityIssues{
		Jira: o.Jira.ToApiModel(ctx),
	}
}

type CatalogEntityIssuesJiraResourceModel struct {
	DefaultJQL types.String `tfsdk:"default_jql"`
	Projects   types.Set    `tfsdk:"projects"`
	Labels     types.Set    `tfsdk:"labels"`
	Components types.Set    `tfsdk:"components"`
}

func (o CatalogEntityIssuesJiraResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityIssuesJira {
	var projects = make([]string, len(o.Projects.Elements()))
	err := o.Labels.ElementsAs(ctx, projects, true)
	if err != nil {
		fmt.Println(err)
	}
	var labels = make([]string, len(o.Labels.Elements()))
	err = o.Labels.ElementsAs(ctx, labels, true)
	if err != nil {
		fmt.Println(err)
	}
	var components = make([]string, len(o.Components.Elements()))
	err = o.Components.ElementsAs(ctx, components, true)
	if err != nil {
		fmt.Println(err)
	}
	return cortex.CatalogEntityIssuesJira{
		DefaultJQL: o.DefaultJQL.ValueString(),
		Projects:   projects,
		Labels:     labels,
		Components: components,
	}
}

/***********************************************************************************************************************
 * On-Call
 ***********************************************************************************************************************/

type CatalogEntityOnCallResourceModel struct {
	PagerDuty CatalogEntityOnCallPagerDutyResourceModel `tfsdk:"pager_duty"`
	OpsGenie  CatalogEntityOnCallOpsGenieResourceModel  `tfsdk:"ops_genie"`
	VictorOps CatalogEntityOnCallVictorOpsResourceModel `tfsdk:"victor_ops"`
}

func (o CatalogEntityOnCallResourceModel) ToApiModel() cortex.CatalogEntityOnCall {
	return cortex.CatalogEntityOnCall{
		PagerDuty: o.PagerDuty.ToApiModel(),
		OpsGenie:  o.OpsGenie.ToApiModel(),
		VictorOps: o.VictorOps.ToApiModel(),
	}
}

type CatalogEntityOnCallPagerDutyResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (o CatalogEntityOnCallPagerDutyResourceModel) ToApiModel() cortex.CatalogEntityOnCallPagerDuty {
	return cortex.CatalogEntityOnCallPagerDuty{
		ID:   o.ID.ValueString(),
		Type: o.Type.ValueString(),
	}
}

type CatalogEntityOnCallOpsGenieResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (o CatalogEntityOnCallOpsGenieResourceModel) ToApiModel() cortex.CatalogEntityOnCallOpsGenie {
	return cortex.CatalogEntityOnCallOpsGenie{
		ID:   o.ID.ValueString(),
		Type: o.Type.ValueString(),
	}
}

type CatalogEntityOnCallVictorOpsResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (o CatalogEntityOnCallVictorOpsResourceModel) ToApiModel() cortex.CatalogEntityOnCallVictorOps {
	return cortex.CatalogEntityOnCallVictorOps{
		ID:   o.ID.ValueString(),
		Type: o.Type.ValueString(),
	}
}

/***********************************************************************************************************************
 * Sentry
 **********************************************************************************************************************/

type CatalogEntitySentryResourceModel struct {
	Project types.String `tfsdk:"project"`
}

func (o CatalogEntitySentryResourceModel) ToApiModel() cortex.CatalogEntitySentry {
	return cortex.CatalogEntitySentry{
		Project: o.Project.ValueString(),
	}
}

/***********************************************************************************************************************
 * Snyk
 **********************************************************************************************************************/

type CatalogEntitySnykResourceModel struct {
	Projects []CatalogEntitySnykProjectResourceModel `tfsdk:"projects"`
}

func (o CatalogEntitySnykResourceModel) ToApiModel() cortex.CatalogEntitySnyk {
	var projects = make([]cortex.CatalogEntitySnykProject, len(o.Projects))
	for i, e := range o.Projects {
		projects[i] = e.ToApiModel()
	}
	return cortex.CatalogEntitySnyk{
		Projects: projects,
	}
}

type CatalogEntitySnykProjectResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	ProjectID    types.String `tfsdk:"project_id"`
	Source       types.String `tfsdk:"source"`
}

func (o CatalogEntitySnykProjectResourceModel) ToApiModel() cortex.CatalogEntitySnykProject {
	return cortex.CatalogEntitySnykProject{
		Organization: o.Organization.ValueString(),
		ProjectID:    o.ProjectID.ValueString(),
		Source:       o.Source.ValueString(),
	}
}
