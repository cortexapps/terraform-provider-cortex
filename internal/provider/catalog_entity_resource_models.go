package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/life4/genesis/slices"
	"strings"
)

// CatalogEntityResourceModel describes the resource data model.
type CatalogEntityResourceModel struct {
	Id             types.String                      `tfsdk:"id"`
	Tag            types.String                      `tfsdk:"tag"`
	Name           types.String                      `tfsdk:"name"`
	Description    types.String                      `tfsdk:"description"`
	Owners         []CatalogEntityOwnerResourceModel `tfsdk:"owners"`
	Groups         []types.String                    `tfsdk:"groups"`
	Links          []CatalogEntityLinkResourceModel  `tfsdk:"links"`
	Metadata       types.String                      `tfsdk:"metadata"`
	Dependencies   []types.Object                    `tfsdk:"dependencies"`
	Alerts         []CatalogEntityAlertResourceModel `tfsdk:"alerts"`
	Apm            types.Object                      `tfsdk:"apm"`
	Dashboards     types.Object                      `tfsdk:"dashboards"`
	Git            types.Object                      `tfsdk:"git"`
	Issues         types.Object                      `tfsdk:"issues"`
	OnCall         types.Object                      `tfsdk:"on_call"`
	SLOs           types.Object                      `tfsdk:"slos"`
	StaticAnalysis types.Object                      `tfsdk:"static_analysis"`
	BugSnag        types.Object                      `tfsdk:"bug_snag"`
	Checkmarx      types.Object                      `tfsdk:"checkmarx"`
	Rollbar        types.Object                      `tfsdk:"rollbar"`
	Sentry         types.Object                      `tfsdk:"sentry"`
	Snyk           types.Object                      `tfsdk:"snyk"`
}

func getDefaultObjectOptions() basetypes.ObjectAsOptions {
	return basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}
}

func (o *CatalogEntityResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityData {
	defaultObjOptions := getDefaultObjectOptions()

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
			fmt.Println("Error parsing custom metadata: ", err)
			metadata = make(map[string]interface{})
		}
	} else {
		metadata = make(map[string]interface{})
	}
	dependencies := make([]cortex.CatalogEntityDependency, len(o.Dependencies))
	for i, dependency := range o.Dependencies {
		dep := &CatalogEntityDependencyResourceModel{}
		err := dependency.As(ctx, dep, defaultObjOptions)
		if err != nil {
			fmt.Println("Error parsing dependency: ", err)
		}
		dependencies[i] = dep.ToApiModel()
	}
	alerts := make([]cortex.CatalogEntityAlert, len(o.Alerts))
	for i, alert := range o.Alerts {
		alerts[i] = alert.ToApiModel()
	}
	dashboards := &CatalogEntityDashboardResourceModel{}
	err := o.Dashboards.As(ctx, dashboards, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing Dashboards configuration: ", err)
	}
	apm := &CatalogEntityApmResourceModel{}
	err = o.Apm.As(ctx, apm, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing APM configuration: ", err)
	}
	git := &CatalogEntityGitResourceModel{}
	err = o.Git.As(ctx, git, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing Git configuration: ", err)
	}
	issues := &CatalogEntityIssuesResourceModel{}
	err = o.Issues.As(ctx, issues, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing Issues configuration: ", err)
	}
	onCall := &CatalogEntityOnCallResourceModel{}
	err = o.OnCall.As(ctx, onCall, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing On-Call configuration: ", err)
	}
	serviceLevelObjectives := &CatalogEntitySLOsResourceModel{}
	err = o.SLOs.As(ctx, serviceLevelObjectives, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing SLO configuration: ", err)
	}
	staticAnalysis := &CatalogEntityStaticAnalysisResourceModel{}
	err = o.StaticAnalysis.As(ctx, staticAnalysis, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing Static Analysis configuration: ", err)
	}
	bugSnag := &CatalogEntityBugSnagResourceModel{}
	err = o.BugSnag.As(ctx, bugSnag, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing BugSnag configuration: ", err)
	}
	checkmarx := &CatalogEntityCheckmarxResourceModel{}
	err = o.Checkmarx.As(ctx, checkmarx, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing Checkmarx configuration: ", err)
	}
	rollbar := &CatalogEntityRollbarResourceModel{}
	err = o.Rollbar.As(ctx, rollbar, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing Rollbar configuration: ", err)
	}
	sentry := &CatalogEntitySentryResourceModel{}
	err = o.Sentry.As(ctx, sentry, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing Sentry configuration: ", err)
	}
	snyk := &CatalogEntitySnykResourceModel{}
	err = o.Snyk.As(ctx, snyk, defaultObjOptions)
	if err != nil {
		fmt.Println("Error parsing Snyk configuration: ", err)
	}

	return cortex.CatalogEntityData{
		Tag:            o.Tag.ValueString(),
		Title:          o.Name.ValueString(),
		Description:    o.Description.ValueString(),
		Owners:         owners,
		Groups:         groups,
		Links:          links,
		Metadata:       metadata,
		Dependencies:   dependencies,
		Alerts:         alerts,
		Dashboards:     dashboards.ToApiModel(),
		Apm:            apm.ToApiModel(),
		Git:            git.ToApiModel(ctx),
		Issues:         issues.ToApiModel(),
		OnCall:         onCall.ToApiModel(),
		SLOs:           serviceLevelObjectives.ToApiModel(),
		StaticAnalysis: staticAnalysis.ToApiModel(),
		BugSnag:        bugSnag.ToApiModel(),
		Checkmarx:      checkmarx.ToApiModel(),
		Rollbar:        rollbar.ToApiModel(),
		Sentry:         sentry.ToApiModel(),
		Snyk:           snyk.ToApiModel(),
	}
}

func (o *CatalogEntityResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityData) {
	o.Id = types.StringValue(entity.Tag)
	o.Name = types.StringValue(entity.Title)
	o.Description = types.StringValue(entity.Description)

	if len(entity.Owners) > 0 {
		o.Owners = make([]CatalogEntityOwnerResourceModel, len(entity.Owners))
		for i, owner := range entity.Owners {
			m := CatalogEntityOwnerResourceModel{}
			o.Owners[i] = m.FromApiModel(&owner)
		}
	} else {
		o.Owners = nil
	}

	if len(entity.Groups) > 0 {
		o.Groups = make([]types.String, len(entity.Groups))
		for i, group := range entity.Groups {
			o.Groups[i] = types.StringValue(group)
		}
	} else {
		o.Groups = nil
	}

	if len(entity.Links) > 0 {
		o.Links = make([]CatalogEntityLinkResourceModel, len(entity.Links))
		for i, link := range entity.Links {
			m := CatalogEntityLinkResourceModel{}
			o.Links[i] = m.FromApiModel(&link)
		}
	} else {
		o.Links = nil
	}

	// coerce map of unknown types into string
	if entity.Metadata != nil && len(entity.Metadata) > 0 {
		metadata, err := json.Marshal(entity.Metadata)
		if err != nil {
			diagnostics.AddError("Error parsing metadata: %s", err.Error())
			return
		}
		o.Metadata = types.StringValue(string(metadata))
	} else {
		o.Metadata = types.StringNull()
	}

	if len(entity.Dependencies) > 0 {
		o.Dependencies = make([]types.Object, len(entity.Dependencies))
		for i, dependency := range entity.Dependencies {
			m := CatalogEntityDependencyResourceModel{}
			o.Dependencies[i] = m.FromApiModel(ctx, diagnostics, &dependency)
		}
	} else {
		o.Dependencies = nil
	}

	if len(entity.Alerts) > 0 {
		o.Alerts = make([]CatalogEntityAlertResourceModel, len(entity.Alerts))
		for i, alert := range entity.Alerts {
			m := CatalogEntityAlertResourceModel{}
			o.Alerts[i] = m.FromApiModel(&alert)
		}
	} else {
		o.Alerts = nil
	}

	dashboards := CatalogEntityDashboardResourceModel{}
	o.Dashboards = dashboards.FromApiModel(ctx, diagnostics, &entity.Dashboards)

	apm := CatalogEntityApmResourceModel{}
	o.Apm = apm.FromApiModel(ctx, diagnostics, &entity.Apm)

	git := CatalogEntityGitResourceModel{}
	o.Git = git.FromApiModel(ctx, diagnostics, &entity.Git)

	issues := CatalogEntityIssuesResourceModel{}
	o.Issues = issues.FromApiModel(ctx, diagnostics, &entity.Issues)

	onCall := CatalogEntityOnCallResourceModel{}
	o.OnCall = onCall.FromApiModel(ctx, diagnostics, &entity.OnCall)

	serviceLevelObjectives := CatalogEntitySLOsResourceModel{}
	o.SLOs = serviceLevelObjectives.FromApiModel(ctx, diagnostics, &entity.SLOs)

	staticAnalysis := CatalogEntityStaticAnalysisResourceModel{}
	o.StaticAnalysis = staticAnalysis.FromApiModel(ctx, diagnostics, &entity.StaticAnalysis)

	bugSnag := CatalogEntityBugSnagResourceModel{}
	o.BugSnag = bugSnag.FromApiModel(ctx, diagnostics, &entity.BugSnag)

	checkmarx := CatalogEntityCheckmarxResourceModel{}
	o.Checkmarx = checkmarx.FromApiModel(ctx, diagnostics, &entity.Checkmarx)

	rollbar := CatalogEntityRollbarResourceModel{}
	o.Rollbar = rollbar.FromApiModel(ctx, diagnostics, &entity.Rollbar)

	sentry := CatalogEntitySentryResourceModel{}
	o.Sentry = sentry.FromApiModel(ctx, diagnostics, &entity.Sentry)

	snyk := CatalogEntitySnykResourceModel{}
	o.Snyk = snyk.FromApiModel(ctx, diagnostics, &entity.Snyk)
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

func (o *CatalogEntityOwnerResourceModel) ToApiModel() cortex.CatalogEntityOwner {
	switch strings.ToLower(o.Type.ValueString()) {
	case "user", "email":
		return cortex.CatalogEntityOwner{
			Type:        o.Type.ValueString(),
			Name:        o.Name.ValueString(),
			Description: o.Description.ValueString(),
			Email:       o.Email.ValueString(),
		}
	case "slack":
		return cortex.CatalogEntityOwner{
			Type:                 o.Type.ValueString(),
			Description:          o.Description.ValueString(),
			Channel:              o.Channel.ValueString(),
			NotificationsEnabled: o.NotificationsEnabled.ValueBool(),
		}
	}
	return cortex.CatalogEntityOwner{
		Type:        o.Type.ValueString(),
		Name:        o.Name.ValueString(),
		Description: o.Description.ValueString(),
		Provider:    o.Provider.ValueString(),
	}
}

func (o *CatalogEntityOwnerResourceModel) FromApiModel(owner *cortex.CatalogEntityOwner) CatalogEntityOwnerResourceModel {
	obj := CatalogEntityOwnerResourceModel{
		Type: types.StringValue(owner.Type),
	}
	if owner.Description != "" {
		obj.Description = types.StringValue(owner.Description)
	} else {
		obj.Description = types.StringNull()
	}

	switch strings.ToLower(owner.Type) {
	case "user", "email":
		obj.Email = types.StringValue(owner.Email)
		obj.Channel = types.StringNull()
		obj.Name = types.StringValue(owner.Name)
		obj.NotificationsEnabled = types.BoolNull()
		obj.Provider = types.StringNull()
	case "slack":
		obj.Email = types.StringNull()
		obj.Channel = types.StringValue(owner.Channel)
		obj.Name = types.StringNull()
		obj.NotificationsEnabled = types.BoolValue(owner.NotificationsEnabled)
		obj.Provider = types.StringNull()
	default: // group
		obj.Email = types.StringNull()
		obj.Channel = types.StringNull()
		obj.Name = types.StringValue(owner.Name)
		obj.NotificationsEnabled = types.BoolNull()
		obj.Provider = types.StringValue(owner.Provider)
	}
	return obj
}

type CatalogEntityLinkResourceModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	Url  types.String `tfsdk:"url"`
}

func (o *CatalogEntityLinkResourceModel) ToApiModel() cortex.CatalogEntityLink {
	return cortex.CatalogEntityLink{
		Type: o.Type.ValueString(),
		Name: o.Name.ValueString(),
		Url:  o.Url.ValueString(),
	}
}

func (o *CatalogEntityLinkResourceModel) FromApiModel(link *cortex.CatalogEntityLink) CatalogEntityLinkResourceModel {
	return CatalogEntityLinkResourceModel{
		Type: types.StringValue(link.Type),
		Name: types.StringValue(link.Name),
		Url:  types.StringValue(link.Url),
	}
}

type CatalogEntityDependencyResourceModel struct {
	Tag         types.String `tfsdk:"tag"`
	Method      types.String `tfsdk:"method"`
	Path        types.String `tfsdk:"path"`
	Description types.String `tfsdk:"description"`
	Metadata    types.String `tfsdk:"metadata"`
}

func (o *CatalogEntityDependencyResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag":         types.StringType,
		"method":      types.StringType,
		"path":        types.StringType,
		"description": types.StringType,
		"metadata":    types.StringType,
	}
}

func (o *CatalogEntityDependencyResourceModel) ToApiModel() cortex.CatalogEntityDependency {
	metadata := make(map[string]interface{})
	if !o.Metadata.IsNull() && !o.Metadata.IsUnknown() && o.Metadata.ValueString() != "" {
		err := json.Unmarshal([]byte(o.Metadata.ValueString()), &metadata)
		if err != nil {
			fmt.Println("Error parsing Dependency configuration: ", err)
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

func (o *CatalogEntityDependencyResourceModel) FromApiModel(ctx context.Context, diag *diag.Diagnostics, dependency *cortex.CatalogEntityDependency) types.Object {
	obj := CatalogEntityDependencyResourceModel{
		Tag:         types.StringValue(dependency.Tag),
		Method:      types.StringValue(dependency.Method),
		Path:        types.StringValue(dependency.Path),
		Description: types.StringValue(dependency.Description),
	}
	if dependency.Metadata == nil {
		dependency.Metadata = map[string]interface{}{}
	}
	depMetadata, err := json.Marshal(dependency.Metadata)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error marshalling Dependency metadata: %+v", err))
		depMetadata = []byte{}
	}
	obj.Metadata = types.StringValue(string(depMetadata))

	retObj, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diag.Append(d...)
	return retObj
}

type CatalogEntityAlertResourceModel struct {
	Type  types.String `tfsdk:"type"`
	Tag   types.String `tfsdk:"tag"`
	Value types.String `tfsdk:"value"`
}

func (o *CatalogEntityAlertResourceModel) ToApiModel() cortex.CatalogEntityAlert {
	return cortex.CatalogEntityAlert{
		Type:  o.Type.ValueString(),
		Tag:   o.Tag.ValueString(),
		Value: o.Value.ValueString(),
	}
}
func (o *CatalogEntityAlertResourceModel) FromApiModel(alert *cortex.CatalogEntityAlert) CatalogEntityAlertResourceModel {
	return CatalogEntityAlertResourceModel{
		Type:  types.StringValue(alert.Type),
		Tag:   types.StringValue(alert.Tag),
		Value: types.StringValue(alert.Value),
	}
}

/***********************************************************************************************************************
 * Git
 ***********************************************************************************************************************/

type CatalogEntityGitResourceModel struct {
	Github    types.Object `tfsdk:"github"`
	Gitlab    types.Object `tfsdk:"gitlab"`
	Azure     types.Object `tfsdk:"azure"`
	Bitbucket types.Object `tfsdk:"bitbucket"`
}

func (o *CatalogEntityGitResourceModel) AttrTypes() map[string]attr.Type {
	gh := CatalogEntityGithubResourceModel{}
	gl := CatalogEntityGitlabResourceModel{}
	az := CatalogEntityAzureResourceModel{}
	bb := CatalogEntityBitbucketResourceModel{}
	return map[string]attr.Type{
		"github":    types.ObjectType{AttrTypes: gh.AttrTypes()},
		"gitlab":    types.ObjectType{AttrTypes: gl.AttrTypes()},
		"azure":     types.ObjectType{AttrTypes: az.AttrTypes()},
		"bitbucket": types.ObjectType{AttrTypes: bb.AttrTypes()},
	}
}

func (o *CatalogEntityGitResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityGit {
	git := cortex.CatalogEntityGit{}
	defaultObjOptions := getDefaultObjectOptions()

	if !o.Github.IsNull() {
		om := CatalogEntityGithubResourceModel{}
		o.Github.As(ctx, &om, defaultObjOptions)
		git.Github = om.ToApiModel()
	}
	if !o.Gitlab.IsNull() {
		om := CatalogEntityGitlabResourceModel{}
		o.Github.As(ctx, &om, defaultObjOptions)
		git.Gitlab = om.ToApiModel()
	}
	if !o.Azure.IsNull() {
		om := CatalogEntityAzureResourceModel{}
		o.Azure.As(ctx, &om, defaultObjOptions)
		git.Azure = om.ToApiModel()
	}
	if !o.Bitbucket.IsNull() {
		om := CatalogEntityBitbucketResourceModel{}
		o.Bitbucket.As(ctx, &om, defaultObjOptions)
		git.BitBucket = om.ToApiModel()
	}
	return git
}

func (o *CatalogEntityGitResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityGit) types.Object {
	git := CatalogEntityGitResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(git.AttrTypes())
	}

	defaultObjOptions := getDefaultObjectOptions()

	ghm := CatalogEntityGithubResourceModel{}
	o.Github.As(ctx, &ghm, defaultObjOptions)
	git.Github = ghm.FromApiModel(ctx, diagnostics, &entity.Github)

	glm := CatalogEntityGitlabResourceModel{}
	o.Gitlab.As(ctx, &glm, defaultObjOptions)
	git.Gitlab = glm.FromApiModel(ctx, diagnostics, &entity.Gitlab)

	azm := CatalogEntityAzureResourceModel{}
	o.Azure.As(ctx, &azm, defaultObjOptions)
	git.Azure = azm.FromApiModel(ctx, diagnostics, &entity.Azure)

	bbm := CatalogEntityBitbucketResourceModel{}
	o.Bitbucket.As(ctx, &bbm, defaultObjOptions)
	git.Bitbucket = bbm.FromApiModel(ctx, diagnostics, &entity.BitBucket)

	obj, d := types.ObjectValueFrom(ctx, git.AttrTypes(), &git)
	diagnostics.Append(d...)
	return obj
}

// Github

type CatalogEntityGithubResourceModel struct {
	Repository types.String `tfsdk:"repository"`
	BasePath   types.String `tfsdk:"base_path"`
}

func (o *CatalogEntityGithubResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository": types.StringType,
		"base_path":  types.StringType,
	}
}

func (o *CatalogEntityGithubResourceModel) ToApiModel() cortex.CatalogEntityGitGithub {
	return cortex.CatalogEntityGitGithub{
		Repository: o.Repository.ValueString(),
		BasePath:   o.BasePath.ValueString(),
	}
}

func (o *CatalogEntityGithubResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityGitGithub) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}
	basePath := types.StringValue(entity.BasePath)
	if entity.BasePath == "" {
		basePath = types.StringNull()
	}
	ghm := CatalogEntityGithubResourceModel{
		Repository: types.StringValue(entity.Repository),
		BasePath:   basePath,
	}
	obj, d := types.ObjectValueFrom(ctx, ghm.AttrTypes(), &ghm)
	diagnostics.Append(d...)
	return obj
}

// Gitlab

type CatalogEntityGitlabResourceModel struct {
	Repository types.String `tfsdk:"repository"`
	BasePath   types.String `tfsdk:"base_path"`
}

func (o *CatalogEntityGitlabResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository": types.StringType,
		"base_path":  types.StringType,
	}
}

func (o *CatalogEntityGitlabResourceModel) ToApiModel() cortex.CatalogEntityGitGitlab {
	return cortex.CatalogEntityGitGitlab{
		Repository: o.Repository.ValueString(),
		BasePath:   o.BasePath.ValueString(),
	}
}

func (o *CatalogEntityGitlabResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityGitGitlab) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	basePath := types.StringValue(entity.BasePath)
	if entity.BasePath == "" {
		basePath = types.StringNull()
	}
	ob := CatalogEntityGitlabResourceModel{
		Repository: types.StringValue(entity.Repository),
		BasePath:   basePath,
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

// AzureOps

type CatalogEntityAzureResourceModel struct {
	Project    types.String `tfsdk:"project"`
	Repository types.String `tfsdk:"repository"`
	BasePath   types.String `tfsdk:"base_path"`
}

func (o *CatalogEntityAzureResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project":    types.StringType,
		"repository": types.StringType,
		"base_path":  types.StringType,
	}
}

func (o *CatalogEntityAzureResourceModel) ToApiModel() cortex.CatalogEntityGitAzureDevOps {
	return cortex.CatalogEntityGitAzureDevOps{
		Project:    o.Project.ValueString(),
		Repository: o.Repository.ValueString(),
		BasePath:   o.BasePath.ValueString(),
	}
}

func (o *CatalogEntityAzureResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityGitAzureDevOps) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	basePath := types.StringValue(entity.BasePath)
	if entity.BasePath == "" {
		basePath = types.StringNull()
	}
	ob := CatalogEntityAzureResourceModel{
		Project:    types.StringValue(entity.Project),
		Repository: types.StringValue(entity.Repository),
		BasePath:   basePath,
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

// Bitbucket

type CatalogEntityBitbucketResourceModel struct {
	Repository types.String `tfsdk:"repository"`
}

func (o *CatalogEntityBitbucketResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository": types.StringType,
	}
}

func (o *CatalogEntityBitbucketResourceModel) ToApiModel() cortex.CatalogEntityGitBitBucket {
	return cortex.CatalogEntityGitBitBucket{
		Repository: o.Repository.ValueString(),
	}
}

func (o *CatalogEntityBitbucketResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityGitBitBucket) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityBitbucketResourceModel{
		Repository: types.StringValue(entity.Repository),
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

/***********************************************************************************************************************
 * Issues
 ***********************************************************************************************************************/

type CatalogEntityIssuesResourceModel struct {
	Jira CatalogEntityIssuesJiraResourceModel `tfsdk:"jira"`
}

func (o *CatalogEntityIssuesResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"jira": types.ObjectType{AttrTypes: o.Jira.AttrTypes()},
	}
}

func (o *CatalogEntityIssuesResourceModel) ToApiModel() cortex.CatalogEntityIssues {
	return cortex.CatalogEntityIssues{
		Jira: o.Jira.ToApiModel(),
	}
}

func (o *CatalogEntityIssuesResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityIssues) types.Object {
	iss := CatalogEntityIssuesResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(iss.AttrTypes())
	}

	jira := CatalogEntityIssuesJiraResourceModel{}
	iss.Jira = jira.FromApiModel(ctx, diagnostics, &entity.Jira)

	obj, d := types.ObjectValueFrom(ctx, iss.AttrTypes(), &iss)
	diagnostics.Append(d...)
	return obj
}

// JIRA

type CatalogEntityIssuesJiraResourceModel struct {
	DefaultJQL types.String `tfsdk:"default_jql"`
	Projects   types.List   `tfsdk:"projects"`
	Labels     types.List   `tfsdk:"labels"`
	Components types.List   `tfsdk:"components"`
}

func (o *CatalogEntityIssuesJiraResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"default_jql": types.StringType,
		"projects":    types.ListType{ElemType: types.StringType},
		"labels":      types.ListType{ElemType: types.StringType},
		"components":  types.ListType{ElemType: types.StringType},
	}
}

func (o *CatalogEntityIssuesJiraResourceModel) ToApiModel() cortex.CatalogEntityIssuesJira {
	var projects = make([]string, len(o.Projects.Elements()))
	for i, e := range o.Projects.Elements() {
		projects[i] = e.(types.String).ValueString()
	}
	var labels = make([]string, len(o.Labels.Elements()))
	for i, e := range o.Labels.Elements() {
		labels[i] = e.(types.String).ValueString()
	}
	var components = make([]string, len(o.Components.Elements()))
	for i, e := range o.Components.Elements() {
		components[i] = e.(types.String).ValueString()
	}
	return cortex.CatalogEntityIssuesJira{
		DefaultJQL: o.DefaultJQL.ValueString(),
		Projects:   slices.Reject(projects, func(i string) bool { return i == "" }),
		Labels:     slices.Reject(labels, func(i string) bool { return i == "" }),
		Components: slices.Reject(components, func(i string) bool { return i == "" }),
	}
}

func (o *CatalogEntityIssuesJiraResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityIssuesJira) CatalogEntityIssuesJiraResourceModel {
	obj := CatalogEntityIssuesJiraResourceModel{}
	if entity.DefaultJQL != "" {
		obj.DefaultJQL = types.StringValue(entity.DefaultJQL)
	} else {
		obj.DefaultJQL = types.StringNull()
	}
	if len(entity.Projects) > 0 {
		projects, d := types.ListValueFrom(ctx, types.StringType, slices.Reject(entity.Projects, func(i string) bool { return i == "" }))
		diagnostics.Append(d...)
		obj.Projects = projects
	} else {
		obj.Projects = types.ListNull(types.StringType)
	}
	if len(entity.Labels) > 0 {
		labels, d := types.ListValueFrom(ctx, types.StringType, slices.Reject(entity.Labels, func(i string) bool { return i == "" }))
		diagnostics.Append(d...)
		obj.Labels = labels
	} else {
		obj.Labels = types.ListNull(types.StringType)
	}
	if len(entity.Components) > 0 {
		components, d := types.ListValueFrom(ctx, types.StringType, slices.Reject(entity.Components, func(i string) bool { return i == "" }))
		diagnostics.Append(d...)
		obj.Components = components
	} else {
		obj.Components = types.ListNull(types.StringType)
	}
	return obj
}

/***********************************************************************************************************************
 * On-Call
 ***********************************************************************************************************************/

type CatalogEntityOnCallResourceModel struct {
	PagerDuty CatalogEntityOnCallPagerDutyResourceModel `tfsdk:"pager_duty"`
	OpsGenie  CatalogEntityOnCallOpsGenieResourceModel  `tfsdk:"ops_genie"`
	VictorOps CatalogEntityOnCallVictorOpsResourceModel `tfsdk:"victor_ops"`
}

func (o *CatalogEntityOnCallResourceModel) AttrTypes() map[string]attr.Type {
	pd := CatalogEntityOnCallPagerDutyResourceModel{}
	og := CatalogEntityOnCallOpsGenieResourceModel{}
	vo := CatalogEntityOnCallVictorOpsResourceModel{}
	return map[string]attr.Type{
		"pager_duty": types.ObjectType{AttrTypes: pd.AttrTypes()},
		"ops_genie":  types.ObjectType{AttrTypes: og.AttrTypes()},
		"victor_ops": types.ObjectType{AttrTypes: vo.AttrTypes()},
	}
}

func (o *CatalogEntityOnCallResourceModel) ToApiModel() cortex.CatalogEntityOnCall {
	return cortex.CatalogEntityOnCall{
		PagerDuty: o.PagerDuty.ToApiModel(),
		OpsGenie:  o.OpsGenie.ToApiModel(),
		VictorOps: o.VictorOps.ToApiModel(),
	}
}

func (o *CatalogEntityOnCallResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, onCall *cortex.CatalogEntityOnCall) types.Object {
	if !onCall.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	pd := CatalogEntityOnCallPagerDutyResourceModel{}
	og := CatalogEntityOnCallOpsGenieResourceModel{}
	vo := CatalogEntityOnCallVictorOpsResourceModel{}
	ob := CatalogEntityOnCallResourceModel{}

	if onCall.PagerDuty.ID != "" {
		ob.PagerDuty = pd.FromApiModel(&onCall.PagerDuty)
	}
	if onCall.OpsGenie.ID != "" {
		ob.OpsGenie = og.FromApiModel(&onCall.OpsGenie)
	}
	if onCall.VictorOps.ID != "" {
		ob.VictorOps = vo.FromApiModel(&onCall.VictorOps)
	}
	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

// PagerDuty

type CatalogEntityOnCallPagerDutyResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (o *CatalogEntityOnCallPagerDutyResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
}

func (o *CatalogEntityOnCallPagerDutyResourceModel) ToApiModel() cortex.CatalogEntityOnCallPagerDuty {
	return cortex.CatalogEntityOnCallPagerDuty{
		ID:   o.ID.ValueString(),
		Type: o.Type.ValueString(),
	}
}

func (o *CatalogEntityOnCallPagerDutyResourceModel) FromApiModel(pagerDuty *cortex.CatalogEntityOnCallPagerDuty) CatalogEntityOnCallPagerDutyResourceModel {
	return CatalogEntityOnCallPagerDutyResourceModel{
		ID:   types.StringValue(pagerDuty.ID),
		Type: types.StringValue(pagerDuty.Type),
	}
}

// OpsGenie

type CatalogEntityOnCallOpsGenieResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (o *CatalogEntityOnCallOpsGenieResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
}

func (o *CatalogEntityOnCallOpsGenieResourceModel) ToApiModel() cortex.CatalogEntityOnCallOpsGenie {
	return cortex.CatalogEntityOnCallOpsGenie{
		ID:   o.ID.ValueString(),
		Type: o.Type.ValueString(),
	}
}

func (o *CatalogEntityOnCallOpsGenieResourceModel) FromApiModel(genie *cortex.CatalogEntityOnCallOpsGenie) CatalogEntityOnCallOpsGenieResourceModel {
	return CatalogEntityOnCallOpsGenieResourceModel{
		ID:   types.StringValue(genie.ID),
		Type: types.StringValue(genie.Type),
	}
}

// VictorOps

type CatalogEntityOnCallVictorOpsResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (o *CatalogEntityOnCallVictorOpsResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
}

func (o *CatalogEntityOnCallVictorOpsResourceModel) ToApiModel() cortex.CatalogEntityOnCallVictorOps {
	return cortex.CatalogEntityOnCallVictorOps{
		ID:   o.ID.ValueString(),
		Type: o.Type.ValueString(),
	}
}

func (o *CatalogEntityOnCallVictorOpsResourceModel) FromApiModel(entity *cortex.CatalogEntityOnCallVictorOps) CatalogEntityOnCallVictorOpsResourceModel {
	return CatalogEntityOnCallVictorOpsResourceModel{
		ID:   types.StringValue(entity.ID),
		Type: types.StringValue(entity.Type),
	}
}

/***********************************************************************************************************************
 * BugSnag
 **********************************************************************************************************************/

type CatalogEntityBugSnagResourceModel struct {
	Project types.String `tfsdk:"project"`
}

func (o *CatalogEntityBugSnagResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.StringType,
	}
}

func (o *CatalogEntityBugSnagResourceModel) ToApiModel() cortex.CatalogEntityBugSnag {
	return cortex.CatalogEntityBugSnag{
		Project: o.Project.ValueString(),
	}
}

func (o *CatalogEntityBugSnagResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityBugSnag) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := &CatalogEntityBugSnagResourceModel{
		Project: types.StringValue(entity.Project),
	}
	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

/***********************************************************************************************************************
 * Checkmarx
 **********************************************************************************************************************/

type CatalogEntityCheckmarxResourceModel struct {
	Projects []CatalogEntityCheckmarxProjectResourceModel `tfsdk:"projects"`
}

func (o *CatalogEntityCheckmarxResourceModel) AttrTypes() map[string]attr.Type {
	ob := CatalogEntityCheckmarxProjectResourceModel{}
	return map[string]attr.Type{
		"projects": types.ListType{ElemType: types.ObjectType{AttrTypes: ob.AttrTypes()}},
	}
}

func (o *CatalogEntityCheckmarxResourceModel) ToApiModel() cortex.CatalogEntityCheckmarx {
	projects := make([]cortex.CatalogEntityCheckmarxProject, len(o.Projects))
	for i, p := range o.Projects {
		projects[i] = p.ToApiModel()
	}
	return cortex.CatalogEntityCheckmarx{
		Projects: projects,
	}
}

func (o *CatalogEntityCheckmarxResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityCheckmarx) types.Object {
	obj := CatalogEntityCheckmarxResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(obj.AttrTypes())
	}

	projects := make([]CatalogEntityCheckmarxProjectResourceModel, len(entity.Projects))
	for i, p := range entity.Projects {
		po := CatalogEntityCheckmarxProjectResourceModel{}
		projects[i] = po.FromApiModel(p)
	}
	obj.Projects = projects

	objectValue, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue

}

type CatalogEntityCheckmarxProjectResourceModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (o *CatalogEntityCheckmarxProjectResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.Int64Type,
		"name": types.StringType,
	}
}

func (o *CatalogEntityCheckmarxProjectResourceModel) ToApiModel() cortex.CatalogEntityCheckmarxProject {
	entity := cortex.CatalogEntityCheckmarxProject{}
	if o.ID.ValueInt64() > 0 {
		entity.ID = o.ID.ValueInt64()
		entity.Name = ""
	} else if o.Name.ValueString() != "" {
		entity.ID = 0
		entity.Name = o.Name.ValueString()
	}
	return entity
}

func (o *CatalogEntityCheckmarxProjectResourceModel) FromApiModel(project cortex.CatalogEntityCheckmarxProject) CatalogEntityCheckmarxProjectResourceModel {
	ob := CatalogEntityCheckmarxProjectResourceModel{}
	if project.ID > 0 {
		ob.ID = types.Int64Value(project.ID)
		ob.Name = types.StringNull()
	} else {
		ob.ID = types.Int64Null()
		ob.Name = types.StringValue(project.Name)
	}
	return ob
}

/***********************************************************************************************************************
 * Rollbar
 **********************************************************************************************************************/

type CatalogEntityRollbarResourceModel struct {
	Project types.String `tfsdk:"project"`
}

func (o *CatalogEntityRollbarResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.StringType,
	}
}

func (o *CatalogEntityRollbarResourceModel) ToApiModel() cortex.CatalogEntityRollbar {
	if o.Project.ValueString() == "" {
		return cortex.CatalogEntityRollbar{}
	}

	return cortex.CatalogEntityRollbar{
		Project: o.Project.ValueString(),
	}
}

func (o *CatalogEntityRollbarResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityRollbar) types.Object {
	ob := CatalogEntityRollbarResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(ob.AttrTypes())
	}

	if entity.Project != "" {
		ob.Project = types.StringValue(entity.Project)
	} else {
		ob.Project = types.StringNull()
	}
	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), ob)
	diagnostics.Append(d...)
	return obj
}

/***********************************************************************************************************************
 * Sentry
 **********************************************************************************************************************/

type CatalogEntitySentryResourceModel struct {
	Project types.String `tfsdk:"project"`
}

func (o *CatalogEntitySentryResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.StringType,
	}
}

func (o *CatalogEntitySentryResourceModel) ToApiModel() cortex.CatalogEntitySentry {
	return cortex.CatalogEntitySentry{
		Project: o.Project.ValueString(),
	}
}

func (o *CatalogEntitySentryResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntitySentry) types.Object {
	ob := CatalogEntitySentryResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(ob.AttrTypes())
	}

	ob.Project = types.StringValue(entity.Project)
	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

/***********************************************************************************************************************
 * Snyk
 **********************************************************************************************************************/

type CatalogEntitySnykResourceModel struct {
	Projects []CatalogEntitySnykProjectResourceModel `tfsdk:"projects"`
}

func (o *CatalogEntitySnykResourceModel) AttrTypes() map[string]attr.Type {
	obp := CatalogEntitySnykProjectResourceModel{}
	return map[string]attr.Type{
		"projects": types.ListType{
			ElemType: types.ObjectType{AttrTypes: obp.AttrTypes()},
		},
	}
}

func (o *CatalogEntitySnykResourceModel) ToApiModel() cortex.CatalogEntitySnyk {
	var projects = make([]cortex.CatalogEntitySnykProject, len(o.Projects))
	for i, e := range o.Projects {
		projects[i] = e.ToApiModel()
	}
	return cortex.CatalogEntitySnyk{
		Projects: projects,
	}
}

func (o *CatalogEntitySnykResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntitySnyk) types.Object {
	obj := &CatalogEntitySnykResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(obj.AttrTypes())
	}

	var projects = make([]CatalogEntitySnykProjectResourceModel, len(entity.Projects))
	for i, e := range entity.Projects {
		ob := CatalogEntitySnykProjectResourceModel{}
		projects[i] = ob.FromApiModel(&e)
	}
	obj.Projects = projects

	objectValue, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

type CatalogEntitySnykProjectResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	ProjectID    types.String `tfsdk:"project_id"`
	Source       types.String `tfsdk:"source"`
}

func (o *CatalogEntitySnykProjectResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"organization": types.StringType,
		"project_id":   types.StringType,
		"source":       types.StringType,
	}
}

func (o *CatalogEntitySnykProjectResourceModel) ToApiModel() cortex.CatalogEntitySnykProject {
	source := "OPEN_SOURCE"
	if o.Source.ValueString() != "" {
		source = o.Source.ValueString()
	}
	return cortex.CatalogEntitySnykProject{
		Organization: o.Organization.ValueString(),
		ProjectID:    o.ProjectID.ValueString(),
		Source:       source,
	}
}

func (o *CatalogEntitySnykProjectResourceModel) FromApiModel(entity *cortex.CatalogEntitySnykProject) CatalogEntitySnykProjectResourceModel {
	ob := CatalogEntitySnykProjectResourceModel{
		Organization: types.StringValue(entity.Organization),
		ProjectID:    types.StringValue(entity.ProjectID),
	}
	if entity.Source != "OPEN_SOURCE" {
		ob.Source = types.StringValue(entity.Source)
	} else {
		ob.Source = types.StringNull()
	}
	return ob
}

/***********************************************************************************************************************
 * APM Configuration
 **********************************************************************************************************************/

type CatalogEntityApmResourceModel struct {
	DataDog   CatalogEntityApmDataDogResourceModel   `tfsdk:"data_dog"`
	Dynatrace CatalogEntityApmDynatraceResourceModel `tfsdk:"dynatrace"`
	NewRelic  CatalogEntityApmNewRelicResourceModel  `tfsdk:"new_relic"`
}

func (o *CatalogEntityApmResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"data_dog":  types.ObjectType{AttrTypes: o.DataDog.AttrTypes()},
		"dynatrace": types.ObjectType{AttrTypes: o.Dynatrace.AttrTypes()},
		"new_relic": types.ObjectType{AttrTypes: o.NewRelic.AttrTypes()},
	}
}

func (o *CatalogEntityApmResourceModel) ToApiModel() cortex.CatalogEntityApm {
	return cortex.CatalogEntityApm{
		DataDog:   o.DataDog.ToApiModel(),
		Dynatrace: o.Dynatrace.ToApiModel(),
		NewRelic:  o.NewRelic.ToApiModel(),
	}
}

func (o *CatalogEntityApmResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityApm) types.Object {
	obj := CatalogEntityApmResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	if entity.DataDog.Enabled() {
		obj.DataDog = o.DataDog.FromApiModel(ctx, diagnostics, &entity.DataDog)
	}
	if entity.Dynatrace.Enabled() {
		obj.Dynatrace = o.Dynatrace.FromApiModel(ctx, diagnostics, &entity.Dynatrace)
	}
	if entity.NewRelic.Enabled() {
		obj.NewRelic = o.NewRelic.FromApiModel(&entity.NewRelic)
	}

	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

// DataDog

type CatalogEntityApmDataDogResourceModel struct {
	Monitors types.List `tfsdk:"monitors"`
}

func (o *CatalogEntityApmDataDogResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"monitors": types.ListType{ElemType: types.Int64Type},
	}
}

func (o *CatalogEntityApmDataDogResourceModel) ToApiModel() cortex.CatalogEntityApmDataDog {
	var monitors = make([]int64, len(o.Monitors.Elements()))
	for i, e := range o.Monitors.Elements() {
		monitors[i] = e.(types.Int64).ValueInt64()
	}
	return cortex.CatalogEntityApmDataDog{
		Monitors: slices.Reject(monitors, func(i int64) bool { return i == 0 }),
	}
}

func (o *CatalogEntityApmDataDogResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityApmDataDog) CatalogEntityApmDataDogResourceModel {
	obj := CatalogEntityApmDataDogResourceModel{}
	monitors := slices.Reject(entity.Monitors, func(i int64) bool { return i == 0 })
	monitorList, d := types.ListValueFrom(ctx, types.Int64Type, monitors)
	diagnostics.Append(d...)
	obj.Monitors = monitorList
	return obj
}

// Dynatrace

type CatalogEntityApmDynatraceResourceModel struct {
	EntityIDs          types.List `tfsdk:"entity_ids"`
	EntityNameMatchers types.List `tfsdk:"entity_name_matchers"`
}

func (o *CatalogEntityApmDynatraceResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"entity_ids":           types.ListType{ElemType: types.StringType},
		"entity_name_matchers": types.ListType{ElemType: types.StringType},
	}
}

func (o *CatalogEntityApmDynatraceResourceModel) ToApiModel() cortex.CatalogEntityApmDynatrace {
	var entityIds = make([]string, len(o.EntityIDs.Elements()))
	for i, e := range o.EntityIDs.Elements() {
		entityIds[i] = e.(types.String).ValueString()
	}
	var entityNameMatchers = make([]string, len(o.EntityNameMatchers.Elements()))
	for i, e := range o.EntityNameMatchers.Elements() {
		entityNameMatchers[i] = e.(types.String).ValueString()
	}
	return cortex.CatalogEntityApmDynatrace{
		EntityIDs:          slices.Reject(entityIds, func(i string) bool { return i == "" }),
		EntityNameMatchers: slices.Reject(entityNameMatchers, func(i string) bool { return i == "" }),
	}
}

func (o *CatalogEntityApmDynatraceResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityApmDynatrace) CatalogEntityApmDynatraceResourceModel {
	obj := CatalogEntityApmDynatraceResourceModel{}

	if len(entity.EntityIDs) > 0 {
		entityIds := slices.Reject(entity.EntityIDs, func(i string) bool { return i == "" })
		eis, d := types.ListValueFrom(ctx, types.StringType, entityIds)
		diagnostics.Append(d...)
		obj.EntityIDs = eis
	} else {
		obj.EntityIDs = types.ListNull(o.AttrTypes()["entity_ids"])
	}

	if len(entity.EntityNameMatchers) > 0 {
		entityNameMatchers := slices.Reject(entity.EntityNameMatchers, func(i string) bool { return i == "" })
		ems, d := types.ListValueFrom(ctx, types.StringType, entityNameMatchers)
		diagnostics.Append(d...)
		obj.EntityNameMatchers = ems
	} else {
		obj.EntityNameMatchers = types.ListNull(o.AttrTypes()["entity_name_matchers"])
	}

	return obj
}

// New Relic

type CatalogEntityApmNewRelicResourceModel struct {
	ApplicationID types.Int64  `tfsdk:"application_id"`
	Alias         types.String `tfsdk:"alias"`
}

func (o *CatalogEntityApmNewRelicResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"application_id": types.Int64Type,
		"alias":          types.StringType,
	}
}

func (o *CatalogEntityApmNewRelicResourceModel) ToApiModel() cortex.CatalogEntityApmNewRelic {
	ob := cortex.CatalogEntityApmNewRelic{
		ApplicationID: o.ApplicationID.ValueInt64(),
	}
	if !o.Alias.IsNull() && !o.Alias.IsUnknown() {
		ob.Alias = o.Alias.ValueString()
	}
	return ob
}

func (o *CatalogEntityApmNewRelicResourceModel) FromApiModel(entity *cortex.CatalogEntityApmNewRelic) CatalogEntityApmNewRelicResourceModel {
	obj := CatalogEntityApmNewRelicResourceModel{}
	if entity.ApplicationID > 0 {
		obj.ApplicationID = types.Int64Value(entity.ApplicationID)
	} else {
		obj.ApplicationID = types.Int64Null()
	}
	if entity.Alias != "" {
		obj.Alias = types.StringValue(entity.Alias)
	} else {
		obj.Alias = types.StringNull()
	}
	return obj
}

/***********************************************************************************************************************
 * Dashboards
 **********************************************************************************************************************/

type CatalogEntityDashboardResourceModel struct {
	Embeds []CatalogEntityDashboardEmbedResourceModel `tfsdk:"embeds"`
}

func (o *CatalogEntityDashboardResourceModel) AttrTypes() map[string]attr.Type {
	em := CatalogEntityDashboardEmbedResourceModel{}
	return map[string]attr.Type{
		"embeds": types.ListType{ElemType: types.ObjectType{AttrTypes: em.AttrTypes()}},
	}
}

func (o *CatalogEntityDashboardResourceModel) ToApiModel() cortex.CatalogEntityDashboards {
	var embeds = make([]cortex.CatalogEntityDashboardsEmbed, len(o.Embeds))
	for i, e := range o.Embeds {
		embeds[i] = e.ToApiModel()
	}
	return cortex.CatalogEntityDashboards{
		Embeds: embeds,
	}
}

func (o *CatalogEntityDashboardResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityDashboards) types.Object {
	obj := CatalogEntityDashboardResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(obj.AttrTypes())
	}

	var embeds = make([]CatalogEntityDashboardEmbedResourceModel, len(entity.Embeds))
	for i, e := range entity.Embeds {
		em := CatalogEntityDashboardEmbedResourceModel{}
		embeds[i] = em.FromApiModel(&e)
	}
	obj.Embeds = embeds
	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

type CatalogEntityDashboardEmbedResourceModel struct {
	Type types.String `tfsdk:"type"`
	URL  types.String `tfsdk:"url"`
}

func (o *CatalogEntityDashboardEmbedResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type": types.StringType,
		"url":  types.StringType,
	}
}

func (o *CatalogEntityDashboardEmbedResourceModel) ToApiModel() cortex.CatalogEntityDashboardsEmbed {
	return cortex.CatalogEntityDashboardsEmbed{
		Type: o.Type.ValueString(),
		URL:  o.URL.ValueString(),
	}
}

func (o *CatalogEntityDashboardEmbedResourceModel) FromApiModel(entity *cortex.CatalogEntityDashboardsEmbed) CatalogEntityDashboardEmbedResourceModel {
	return CatalogEntityDashboardEmbedResourceModel{
		Type: types.StringValue(entity.Type),
		URL:  types.StringValue(entity.URL),
	}
}

/***********************************************************************************************************************
 * SLOs
 **********************************************************************************************************************/

type CatalogEntitySLOsResourceModel struct {
	DataDog    []CatalogEntitySLODataDogResourceModel    `tfsdk:"data_dog"`
	Dynatrace  []CatalogEntitySLODynatraceResourceModel  `tfsdk:"dynatrace"`
	Lightstep  CatalogEntitySLOLightstepResourceModel    `tfsdk:"lightstep"`
	Prometheus []CatalogEntitySLOPrometheusResourceModel `tfsdk:"prometheus"`
	SignalFX   []CatalogEntitySLOSignalFxResourceModel   `tfsdk:"signal_fx"`
	SumoLogic  []CatalogEntitySLOSumoLogicResourceModel  `tfsdk:"sumo_logic"`
}

func (o *CatalogEntitySLOsResourceModel) AttrTypes() map[string]attr.Type {
	dd := CatalogEntitySLODataDogResourceModel{}
	dt := CatalogEntitySLODynatraceResourceModel{}
	ls := CatalogEntitySLOLightstepResourceModel{}
	p := CatalogEntitySLOPrometheusResourceModel{}
	sf := CatalogEntitySLOSignalFxResourceModel{}
	sl := CatalogEntitySLOSumoLogicResourceModel{}
	return map[string]attr.Type{
		"data_dog":   types.ListType{ElemType: types.ObjectType{AttrTypes: dd.AttrTypes()}},
		"dynatrace":  types.ListType{ElemType: types.ObjectType{AttrTypes: dt.AttrTypes()}},
		"lightstep":  types.ObjectType{AttrTypes: ls.AttrTypes()},
		"prometheus": types.ListType{ElemType: types.ObjectType{AttrTypes: p.AttrTypes()}},
		"signal_fx":  types.ListType{ElemType: types.ObjectType{AttrTypes: sf.AttrTypes()}},
		"sumo_logic": types.ListType{ElemType: types.ObjectType{AttrTypes: sl.AttrTypes()}},
	}
}

func (o *CatalogEntitySLOsResourceModel) ToApiModel() cortex.CatalogEntitySLOs {
	var dataDog = make([]cortex.CatalogEntitySLODataDog, len(o.DataDog))
	for i, e := range o.DataDog {
		dataDog[i] = e.ToApiModel()
	}
	var dynatrace = make([]cortex.CatalogEntitySLODynatrace, len(o.Dynatrace))
	for i, e := range o.Dynatrace {
		dynatrace[i] = e.ToApiModel()
	}
	var prometheusQueries = make([]cortex.CatalogEntitySLOPrometheusQuery, len(o.Prometheus))
	for i, e := range o.Prometheus {
		prometheusQueries[i] = e.ToApiModel()
	}
	var signalFx = make([]cortex.CatalogEntitySLOSignalFX, len(o.SignalFX))
	for i, e := range o.SignalFX {
		signalFx[i] = e.ToApiModel()
	}
	var sumoLogic = make([]cortex.CatalogEntitySLOSumoLogic, len(o.SumoLogic))
	for i, e := range o.SumoLogic {
		sumoLogic[i] = e.ToApiModel()
	}
	return cortex.CatalogEntitySLOs{
		DataDog:    dataDog,
		Dynatrace:  dynatrace,
		Lightstep:  o.Lightstep.ToApiModel(),
		Prometheus: prometheusQueries,
		SignalFX:   signalFx,
		SumoLogic:  sumoLogic,
	}
}

func (o *CatalogEntitySLOsResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntitySLOs) types.Object {
	obj := CatalogEntitySLOsResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(obj.AttrTypes())
	}

	if len(entity.DataDog) > 0 {
		obj.DataDog = make([]CatalogEntitySLODataDogResourceModel, len(entity.DataDog))
		for i, e := range entity.DataDog {
			m := CatalogEntitySLODataDogResourceModel{}
			obj.DataDog[i] = m.FromApiModel(&e)
		}
	}
	if len(entity.Dynatrace) > 0 {
		obj.Dynatrace = make([]CatalogEntitySLODynatraceResourceModel, len(entity.Dynatrace))
		for i, e := range entity.Dynatrace {
			m := CatalogEntitySLODynatraceResourceModel{}
			obj.Dynatrace[i] = m.FromApiModel(&e)
		}
	}
	obj.Lightstep = obj.Lightstep.FromApiModel(entity.Lightstep)
	if len(entity.Prometheus) > 0 {
		obj.Prometheus = make([]CatalogEntitySLOPrometheusResourceModel, len(entity.Prometheus))
		for i, e := range entity.Prometheus {
			m := CatalogEntitySLOPrometheusResourceModel{}
			obj.Prometheus[i] = m.FromApiModel(&e)
		}
	}
	if len(entity.SignalFX) > 0 {
		obj.SignalFX = make([]CatalogEntitySLOSignalFxResourceModel, len(entity.SignalFX))
		for i, e := range entity.SignalFX {
			m := CatalogEntitySLOSignalFxResourceModel{}
			obj.SignalFX[i] = m.FromApiModel(&e)
		}
	}
	if len(entity.SumoLogic) > 0 {
		obj.SumoLogic = make([]CatalogEntitySLOSumoLogicResourceModel, len(entity.SumoLogic))
		for i, e := range entity.SumoLogic {
			m := CatalogEntitySLOSumoLogicResourceModel{}
			obj.SumoLogic[i] = m.FromApiModel(&e)
		}
	}
	objectType, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectType
}

// DataDog

type CatalogEntitySLODataDogResourceModel struct {
	ID types.String `tfsdk:"id"`
}

func (o *CatalogEntitySLODataDogResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": types.StringType,
	}
}

func (o *CatalogEntitySLODataDogResourceModel) ToApiModel() cortex.CatalogEntitySLODataDog {
	return cortex.CatalogEntitySLODataDog{
		ID: o.ID.ValueString(),
	}
}

func (o *CatalogEntitySLODataDogResourceModel) FromApiModel(entity *cortex.CatalogEntitySLODataDog) CatalogEntitySLODataDogResourceModel {
	return CatalogEntitySLODataDogResourceModel{
		ID: types.StringValue(entity.ID),
	}
}

// Dynatrace

type CatalogEntitySLODynatraceResourceModel struct {
	ID types.String `tfsdk:"id"`
}

func (o *CatalogEntitySLODynatraceResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": types.StringType,
	}
}

func (o *CatalogEntitySLODynatraceResourceModel) ToApiModel() cortex.CatalogEntitySLODynatrace {
	return cortex.CatalogEntitySLODynatrace{
		ID: o.ID.ValueString(),
	}
}

func (o *CatalogEntitySLODynatraceResourceModel) FromApiModel(entity *cortex.CatalogEntitySLODynatrace) CatalogEntitySLODynatraceResourceModel {
	return CatalogEntitySLODynatraceResourceModel{
		ID: types.StringValue(entity.ID),
	}
}

// LightStep

type CatalogEntitySLOLightstepResourceModel struct {
	Streams []CatalogEntitySLOLightstepStreamResourceModel `tfsdk:"streams"`
}

func (o *CatalogEntitySLOLightstepResourceModel) AttrTypes() map[string]attr.Type {
	s := CatalogEntitySLOLightstepStreamResourceModel{}
	return map[string]attr.Type{
		"streams": types.ListType{ElemType: types.ObjectType{AttrTypes: s.AttrTypes()}},
	}
}

func (o *CatalogEntitySLOLightstepResourceModel) ToApiModel() []cortex.CatalogEntitySLOLightstepStream {
	var arr = make([]cortex.CatalogEntitySLOLightstepStream, len(o.Streams))
	for i, e := range o.Streams {
		arr[i] = e.ToApiModel()
	}
	return arr
}

func (o *CatalogEntitySLOLightstepResourceModel) FromApiModel(streams []cortex.CatalogEntitySLOLightstepStream) CatalogEntitySLOLightstepResourceModel {
	var arr = make([]CatalogEntitySLOLightstepStreamResourceModel, len(streams))
	for i, e := range streams {
		m := CatalogEntitySLOLightstepStreamResourceModel{}
		arr[i] = m.FromApiModel(&e)
	}
	return CatalogEntitySLOLightstepResourceModel{
		Streams: arr,
	}
}

type CatalogEntitySLOLightstepStreamResourceModel struct {
	StreamID types.String                                       `tfsdk:"stream_id"`
	Targets  CatalogEntitySLOLightstepStreamTargetResourceModel `tfsdk:"targets"`
}

func (o *CatalogEntitySLOLightstepStreamResourceModel) AttrTypes() map[string]attr.Type {
	t := CatalogEntitySLOLightstepStreamTargetResourceModel{}
	return map[string]attr.Type{
		"stream_id": types.StringType,
		"targets":   types.ObjectType{AttrTypes: t.AttrTypes()},
	}
}

func (o *CatalogEntitySLOLightstepStreamResourceModel) ToApiModel() cortex.CatalogEntitySLOLightstepStream {
	return cortex.CatalogEntitySLOLightstepStream{
		StreamID: o.StreamID.ValueString(),
		Targets:  o.Targets.ToApiModel(),
	}
}

func (o *CatalogEntitySLOLightstepStreamResourceModel) FromApiModel(entity *cortex.CatalogEntitySLOLightstepStream) CatalogEntitySLOLightstepStreamResourceModel {
	targets := CatalogEntitySLOLightstepStreamTargetResourceModel{}
	return CatalogEntitySLOLightstepStreamResourceModel{
		StreamID: types.StringValue(entity.StreamID),
		Targets:  targets.FromApiModel(&entity.Targets),
	}
}

type CatalogEntitySLOLightstepStreamTargetResourceModel struct {
	Latencies []CatalogEntitySLOLightstepStreamTargetLatencyResourceModel `tfsdk:"latencies"`
}

func (o *CatalogEntitySLOLightstepStreamTargetResourceModel) AttrTypes() map[string]attr.Type {
	l := CatalogEntitySLOLightstepStreamTargetLatencyResourceModel{}
	return map[string]attr.Type{
		"latencies": types.ListType{ElemType: types.ObjectType{AttrTypes: l.AttrTypes()}},
	}
}

func (o *CatalogEntitySLOLightstepStreamTargetResourceModel) ToApiModel() cortex.CatalogEntitySLOLightstepTargets {
	var latencies = make([]cortex.CatalogEntitySLOLightstepTargetLatency, len(o.Latencies))
	for i, e := range o.Latencies {
		latencies[i] = e.ToApiModel()
	}
	return cortex.CatalogEntitySLOLightstepTargets{
		Latencies: latencies,
	}
}

func (o *CatalogEntitySLOLightstepStreamTargetResourceModel) FromApiModel(entity *cortex.CatalogEntitySLOLightstepTargets) CatalogEntitySLOLightstepStreamTargetResourceModel {
	latencies := make([]CatalogEntitySLOLightstepStreamTargetLatencyResourceModel, len(entity.Latencies))
	for i, e := range entity.Latencies {
		m := CatalogEntitySLOLightstepStreamTargetLatencyResourceModel{}
		latencies[i] = m.FromApiModel(&e)
	}
	return CatalogEntitySLOLightstepStreamTargetResourceModel{
		Latencies: latencies,
	}
}

type CatalogEntitySLOLightstepStreamTargetLatencyResourceModel struct {
	Percentile types.Float64 `tfsdk:"percentile"`
	Target     types.Int64   `tfsdk:"target"`
	SLO        types.Float64 `tfsdk:"slo"`
}

func (o *CatalogEntitySLOLightstepStreamTargetLatencyResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"percentile": types.Float64Type,
		"target":     types.Int64Type,
		"slo":        types.Float64Type,
	}
}

func (o *CatalogEntitySLOLightstepStreamTargetLatencyResourceModel) ToApiModel() cortex.CatalogEntitySLOLightstepTargetLatency {
	return cortex.CatalogEntitySLOLightstepTargetLatency{
		Percentile: o.Percentile.ValueFloat64(),
		Target:     o.Target.ValueInt64(),
		SLO:        o.SLO.ValueFloat64(),
	}
}

func (o *CatalogEntitySLOLightstepStreamTargetLatencyResourceModel) FromApiModel(entity *cortex.CatalogEntitySLOLightstepTargetLatency) CatalogEntitySLOLightstepStreamTargetLatencyResourceModel {
	return CatalogEntitySLOLightstepStreamTargetLatencyResourceModel{
		Percentile: types.Float64Value(entity.Percentile),
		Target:     types.Int64Value(entity.Target),
		SLO:        types.Float64Value(entity.SLO),
	}
}

// Prometheus

type CatalogEntitySLOPrometheusResourceModel struct {
	ErrorQuery types.String  `tfsdk:"error_query"`
	TotalQuery types.String  `tfsdk:"total_query"`
	SLO        types.Float64 `tfsdk:"slo"`
	Name       types.String  `tfsdk:"name"`
	Alias      types.String  `tfsdk:"alias"`
}

func (o *CatalogEntitySLOPrometheusResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"error_query": types.StringType,
		"total_query": types.StringType,
		"slo":         types.Float64Type,
		"name":        types.StringType,
		"alias":       types.StringType,
	}
}

func (o *CatalogEntitySLOPrometheusResourceModel) ToApiModel() cortex.CatalogEntitySLOPrometheusQuery {
	ob := cortex.CatalogEntitySLOPrometheusQuery{
		ErrorQuery: o.ErrorQuery.ValueString(),
		TotalQuery: o.TotalQuery.ValueString(),
		SLO:        o.SLO.ValueFloat64(),
	}
	if o.Name.ValueString() != "" {
		ob.Name = o.Name.ValueString()
	}
	if o.Alias.ValueString() != "" {
		ob.Alias = o.Alias.ValueString()
	}
	return ob
}

func (o *CatalogEntitySLOPrometheusResourceModel) FromApiModel(entity *cortex.CatalogEntitySLOPrometheusQuery) CatalogEntitySLOPrometheusResourceModel {
	ob := CatalogEntitySLOPrometheusResourceModel{
		ErrorQuery: types.StringValue(entity.ErrorQuery),
		TotalQuery: types.StringValue(entity.TotalQuery),
		SLO:        types.Float64Value(entity.SLO),
	}
	if entity.Name != "" {
		ob.Name = types.StringValue(entity.Name)
	} else {
		ob.Name = types.StringNull()
	}
	if entity.Alias != "" {
		ob.Alias = types.StringValue(entity.Alias)
	} else {
		ob.Alias = types.StringNull()
	}
	return ob
}

// SignalFX

type CatalogEntitySLOSignalFxResourceModel struct {
	Query     types.String `tfsdk:"query"`
	Rollup    types.String `tfsdk:"rollup"`
	Target    types.Int64  `tfsdk:"target"`
	Lookback  types.String `tfsdk:"lookback"`
	Operation types.String `tfsdk:"operation"`
}

func (o *CatalogEntitySLOSignalFxResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"query":     types.StringType,
		"rollup":    types.StringType,
		"target":    types.Int64Type,
		"lookback":  types.StringType,
		"operation": types.StringType,
	}
}

func (o *CatalogEntitySLOSignalFxResourceModel) ToApiModel() cortex.CatalogEntitySLOSignalFX {
	return cortex.CatalogEntitySLOSignalFX{
		Query:     o.Query.ValueString(),
		Rollup:    o.Rollup.ValueString(),
		Target:    o.Target.ValueInt64(),
		Lookback:  o.Lookback.ValueString(),
		Operation: o.Operation.ValueString(),
	}
}

func (o *CatalogEntitySLOSignalFxResourceModel) FromApiModel(entity *cortex.CatalogEntitySLOSignalFX) CatalogEntitySLOSignalFxResourceModel {
	return CatalogEntitySLOSignalFxResourceModel{
		Query:     types.StringValue(entity.Query),
		Rollup:    types.StringValue(entity.Rollup),
		Target:    types.Int64Value(entity.Target),
		Lookback:  types.StringValue(entity.Lookback),
		Operation: types.StringValue(entity.Operation),
	}
}

// SumoLogic

type CatalogEntitySLOSumoLogicResourceModel struct {
	ID types.String `tfsdk:"id"`
}

func (o *CatalogEntitySLOSumoLogicResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": types.StringType,
	}
}

func (o *CatalogEntitySLOSumoLogicResourceModel) ToApiModel() cortex.CatalogEntitySLOSumoLogic {
	return cortex.CatalogEntitySLOSumoLogic{
		ID: o.ID.ValueString(),
	}
}

func (o *CatalogEntitySLOSumoLogicResourceModel) FromApiModel(entity *cortex.CatalogEntitySLOSumoLogic) CatalogEntitySLOSumoLogicResourceModel {
	return CatalogEntitySLOSumoLogicResourceModel{
		ID: types.StringValue(entity.ID),
	}
}

/***********************************************************************************************************************
 * Static Analysis
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisResourceModel struct {
	CodeCov   CatalogEntityStaticAnalysisCodeCovResourceModel   `tfsdk:"code_cov"`
	Mend      CatalogEntityStaticAnalysisMendResourceModel      `tfsdk:"mend"`
	SonarQube CatalogEntityStaticAnalysisSonarQubeResourceModel `tfsdk:"sonar_qube"`
	Veracode  CatalogEntityStaticAnalysisVeracodeResourceModel  `tfsdk:"veracode"`
}

func (o *CatalogEntityStaticAnalysisResourceModel) AttrTypes() map[string]attr.Type {
	cc := CatalogEntityStaticAnalysisCodeCovResourceModel{}
	me := CatalogEntityStaticAnalysisMendResourceModel{}
	sq := CatalogEntityStaticAnalysisSonarQubeResourceModel{}
	vc := CatalogEntityStaticAnalysisVeracodeResourceModel{}
	return map[string]attr.Type{
		"code_cov":   types.ObjectType{AttrTypes: cc.AttrTypes()},
		"mend":       types.ObjectType{AttrTypes: me.AttrTypes()},
		"sonar_qube": types.ObjectType{AttrTypes: sq.AttrTypes()},
		"veracode":   types.ObjectType{AttrTypes: vc.AttrTypes()},
	}
}

func (o *CatalogEntityStaticAnalysisResourceModel) ToApiModel() cortex.CatalogEntityStaticAnalysis {
	return cortex.CatalogEntityStaticAnalysis{
		CodeCov:   o.CodeCov.ToApiModel(),
		Mend:      o.Mend.ToApiModel(),
		SonarQube: o.SonarQube.ToApiModel(),
		Veracode:  o.Veracode.ToApiModel(),
	}
}

func (o *CatalogEntityStaticAnalysisResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityStaticAnalysis) types.Object {
	ob := CatalogEntityStaticAnalysisResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(ob.AttrTypes())
	}

	if entity.CodeCov.Enabled() {
		cc := CatalogEntityStaticAnalysisCodeCovResourceModel{}
		ob.CodeCov = cc.FromApiModel(&entity.CodeCov)
	}

	if entity.Mend.Enabled() {
		me := CatalogEntityStaticAnalysisMendResourceModel{}
		ob.Mend = me.FromApiModel(&entity.Mend)
	}

	if entity.SonarQube.Enabled() {
		sq := CatalogEntityStaticAnalysisSonarQubeResourceModel{}
		ob.SonarQube = sq.FromApiModel(&entity.SonarQube)
	}

	if entity.Veracode.Enabled() {
		vc := CatalogEntityStaticAnalysisVeracodeResourceModel{}
		ob.Veracode = vc.FromApiModel(&entity.Veracode)
	}

	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), ob)
	diagnostics.Append(d...)
	return obj
}

// CodeCov

type CatalogEntityStaticAnalysisCodeCovResourceModel struct {
	Repository types.String `tfsdk:"repository"`
	Provider   types.String `tfsdk:"provider"`
}

func (o *CatalogEntityStaticAnalysisCodeCovResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"provider":   types.StringType,
		"repository": types.StringType,
	}
}

func (o *CatalogEntityStaticAnalysisCodeCovResourceModel) ToApiModel() cortex.CatalogEntityStaticAnalysisCodeCov {
	return cortex.CatalogEntityStaticAnalysisCodeCov{
		Repository: o.Repository.ValueString(),
		Provider:   o.Provider.ValueString(),
	}
}

func (o *CatalogEntityStaticAnalysisCodeCovResourceModel) FromApiModel(cov *cortex.CatalogEntityStaticAnalysisCodeCov) CatalogEntityStaticAnalysisCodeCovResourceModel {
	return CatalogEntityStaticAnalysisCodeCovResourceModel{
		Repository: types.StringValue(cov.Repository),
		Provider:   types.StringValue(cov.Provider),
	}
}

// Mend

type CatalogEntityStaticAnalysisMendResourceModel struct {
	ApplicationIDs []types.String `tfsdk:"application_ids"`
	ProjectIDs     []types.String `tfsdk:"project_ids"`
}

func (o *CatalogEntityStaticAnalysisMendResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"application_ids": types.ListType{ElemType: types.StringType},
		"project_ids":     types.ListType{ElemType: types.StringType},
	}
}

func (o *CatalogEntityStaticAnalysisMendResourceModel) ToApiModel() cortex.CatalogEntityStaticAnalysisMend {
	applicationIds := make([]string, len(o.ApplicationIDs))
	for i, e := range o.ApplicationIDs {
		applicationIds[i] = e.ValueString()
	}
	projectIds := make([]string, len(o.ProjectIDs))
	for i, e := range o.ProjectIDs {
		projectIds[i] = e.ValueString()
	}
	return cortex.CatalogEntityStaticAnalysisMend{
		ApplicationIDs: applicationIds,
		ProjectIDs:     projectIds,
	}
}

func (o *CatalogEntityStaticAnalysisMendResourceModel) FromApiModel(entity *cortex.CatalogEntityStaticAnalysisMend) CatalogEntityStaticAnalysisMendResourceModel {
	applicationIds := make([]types.String, len(entity.ApplicationIDs))
	for i, e := range entity.ApplicationIDs {
		applicationIds[i] = types.StringValue(e)
	}
	projectIds := make([]types.String, len(entity.ProjectIDs))
	for i, e := range entity.ProjectIDs {
		projectIds[i] = types.StringValue(e)
	}
	return CatalogEntityStaticAnalysisMendResourceModel{
		ApplicationIDs: applicationIds,
		ProjectIDs:     projectIds,
	}
}

// SonarQube

type CatalogEntityStaticAnalysisSonarQubeResourceModel struct {
	Project types.String `tfsdk:"project"`
	Alias   types.String `tfsdk:"alias"`
}

func (o *CatalogEntityStaticAnalysisSonarQubeResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.StringType,
		"alias":   types.StringType,
	}
}

func (o *CatalogEntityStaticAnalysisSonarQubeResourceModel) ToApiModel() cortex.CatalogEntityStaticAnalysisSonarQube {
	entity := cortex.CatalogEntityStaticAnalysisSonarQube{
		Project: o.Project.ValueString(),
	}
	if !o.Alias.IsNull() && !o.Alias.IsUnknown() && o.Alias.ValueString() != "" {
		entity.Alias = o.Alias.ValueString()
	}
	return entity
}

func (o *CatalogEntityStaticAnalysisSonarQubeResourceModel) FromApiModel(entity *cortex.CatalogEntityStaticAnalysisSonarQube) CatalogEntityStaticAnalysisSonarQubeResourceModel {
	ob := CatalogEntityStaticAnalysisSonarQubeResourceModel{
		Project: types.StringValue(entity.Project),
	}
	if entity.Alias != "" {
		ob.Alias = types.StringValue(entity.Alias)
	} else {
		ob.Alias = types.StringNull()
	}
	return ob
}

// Veracode

type CatalogEntityStaticAnalysisVeracodeResourceModel struct {
	ApplicationNames []types.String                                            `tfsdk:"application_names"`
	Sandboxes        []CatalogEntityStaticAnalysisVeracodeSandboxResourceModel `tfsdk:"sandboxes"`
}

func (o *CatalogEntityStaticAnalysisVeracodeResourceModel) AttrTypes() map[string]attr.Type {
	sb := CatalogEntityStaticAnalysisVeracodeSandboxResourceModel{}
	return map[string]attr.Type{
		"application_names": types.ListType{ElemType: types.StringType},
		"sandboxes":         types.ListType{ElemType: types.ObjectType{AttrTypes: sb.AttrTypes()}},
	}
}

func (o *CatalogEntityStaticAnalysisVeracodeResourceModel) ToApiModel() cortex.CatalogEntityStaticAnalysisVeracode {
	var sandboxes = make([]cortex.CatalogEntityStaticAnalysisVeracodeSandbox, len(o.Sandboxes))
	for i, e := range o.Sandboxes {
		sandboxes[i] = e.ToApiModel()
	}
	applicationNames := make([]string, len(o.ApplicationNames))
	for i, e := range o.ApplicationNames {
		applicationNames[i] = e.ValueString()
	}
	return cortex.CatalogEntityStaticAnalysisVeracode{
		ApplicationNames: applicationNames,
		Sandboxes:        sandboxes,
	}
}

func (o *CatalogEntityStaticAnalysisVeracodeResourceModel) FromApiModel(entity *cortex.CatalogEntityStaticAnalysisVeracode) CatalogEntityStaticAnalysisVeracodeResourceModel {
	var sandboxes = make([]CatalogEntityStaticAnalysisVeracodeSandboxResourceModel, len(entity.Sandboxes))
	for i, e := range entity.Sandboxes {
		ob := CatalogEntityStaticAnalysisVeracodeSandboxResourceModel{}
		sandboxes[i] = ob.FromApiModel(&e)
	}
	applicationNames := make([]types.String, len(entity.ApplicationNames))
	for i, e := range entity.ApplicationNames {
		applicationNames[i] = types.StringValue(e)
	}
	return CatalogEntityStaticAnalysisVeracodeResourceModel{
		ApplicationNames: applicationNames,
		Sandboxes:        sandboxes,
	}
}

type CatalogEntityStaticAnalysisVeracodeSandboxResourceModel struct {
	ApplicationName types.String `tfsdk:"application_name"`
	SandboxName     types.String `tfsdk:"sandbox_name"`
}

func (o *CatalogEntityStaticAnalysisVeracodeSandboxResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"application_name": types.StringType,
		"sandbox_name":     types.StringType,
	}
}

func (o *CatalogEntityStaticAnalysisVeracodeSandboxResourceModel) ToApiModel() cortex.CatalogEntityStaticAnalysisVeracodeSandbox {
	return cortex.CatalogEntityStaticAnalysisVeracodeSandbox{
		ApplicationName: o.ApplicationName.ValueString(),
		SandboxName:     o.SandboxName.ValueString(),
	}
}

func (o *CatalogEntityStaticAnalysisVeracodeSandboxResourceModel) FromApiModel(entity *cortex.CatalogEntityStaticAnalysisVeracodeSandbox) CatalogEntityStaticAnalysisVeracodeSandboxResourceModel {
	return CatalogEntityStaticAnalysisVeracodeSandboxResourceModel{
		ApplicationName: types.StringValue(entity.ApplicationName),
		SandboxName:     types.StringValue(entity.SandboxName),
	}
}
