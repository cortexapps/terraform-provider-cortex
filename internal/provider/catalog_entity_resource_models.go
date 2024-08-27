package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/life4/genesis/slices"
	"strings"
)

// CatalogEntityResourceModel describes the resource data model.
type CatalogEntityResourceModel struct {
	Id             types.String                             `tfsdk:"id"`
	Tag            types.String                             `tfsdk:"tag"`
	Name           types.String                             `tfsdk:"name"`
	Description    types.String                             `tfsdk:"description"`
	Type           types.String                             `tfsdk:"type"`
	Definition     types.String                             `tfsdk:"definition"`
	Owners         []CatalogEntityOwnerResourceModel        `tfsdk:"owners"`
	Children       []CatalogEntityChildResourceModel        `tfsdk:"children"`
	Parents  	   []CatalogEntityParentResourceModel       `tfsdk:"parents"`
	Groups         []types.String                           `tfsdk:"groups"`
	Links          []CatalogEntityLinkResourceModel         `tfsdk:"links"`
	IgnoreMetadata types.Bool                               `tfsdk:"ignore_metadata"`
	Metadata       types.String                             `tfsdk:"metadata"`
	Dependencies   []types.Object                           `tfsdk:"dependencies"`
	Alerts         []types.Object                           `tfsdk:"alerts"`
	Apm            types.Object                             `tfsdk:"apm"`
	Dashboards     types.Object                             `tfsdk:"dashboards"`
	Git            types.Object                             `tfsdk:"git"`
	Issues         types.Object                             `tfsdk:"issues"`
	OnCall         types.Object                             `tfsdk:"on_call"`
	SLOs           types.Object                             `tfsdk:"slos"`
	StaticAnalysis types.Object                             `tfsdk:"static_analysis"`
	CiCd           types.Object                             `tfsdk:"ci_cd"`
	BugSnag        types.Object                             `tfsdk:"bug_snag"`
	Checkmarx      types.Object                             `tfsdk:"checkmarx"`
	Coralogix      types.Object                             `tfsdk:"coralogix"`
	FireHydrant    types.Object                             `tfsdk:"firehydrant"`
	K8s            types.Object                             `tfsdk:"k8s"`
	LaunchDarkly   types.Object                             `tfsdk:"launch_darkly"`
	MicrosoftTeams []types.Object                           `tfsdk:"microsoft_teams"`
	Rollbar        types.Object                             `tfsdk:"rollbar"`
	Sentry         types.Object                             `tfsdk:"sentry"`
	ServiceNow     types.Object                             `tfsdk:"service_now"`
	Slack          types.Object                             `tfsdk:"slack"`
	Snyk           types.Object                             `tfsdk:"snyk"`
	Wiz            types.Object                             `tfsdk:"wiz"`
	Team           types.Object                             `tfsdk:"team"`
}

func getDefaultObjectOptions() basetypes.ObjectAsOptions {
	return basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}
}

func (o *CatalogEntityResourceModel) ToApiModel(ctx context.Context, diagnostics *diag.Diagnostics) cortex.CatalogEntityData {
	defaultObjOptions := getDefaultObjectOptions()

	definition := make(map[string]interface{})
	if !o.Definition.IsNull() && !o.Definition.IsUnknown() && o.Definition.ValueString() != "" {
		err := json.Unmarshal([]byte(o.Definition.ValueString()), &definition)
		if err != nil {
			diagnostics.AddError("error parsing x-cortex-definition", fmt.Sprintf("%+v", err))
			definition = make(map[string]interface{})
		}
	} else {
		definition = make(map[string]interface{})
	}
	owners := make([]cortex.CatalogEntityOwner, len(o.Owners))
	for i, owner := range o.Owners {
		owners[i] = owner.ToApiModel()
	}
	children := make([]cortex.CatalogEntityChild, len(o.Children))
	for i, child := range o.Children {
		children[i] = child.ToApiModel()
	}
	Parents := make([]cortex.CatalogEntityParent, len(o.Parents))
	for i, parent := range o.Parents {
		Parents[i] = parent.ToApiModel()
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
	if !o.IgnoreMetadata.ValueBool() {
		if !o.Metadata.IsNull() && !o.Metadata.IsUnknown() && o.Metadata.ValueString() != "" {
			err := json.Unmarshal([]byte(o.Metadata.ValueString()), &metadata)
			if err != nil {
				diagnostics.AddError("error parsing custom metadata", fmt.Sprintf("%+v", err))
				metadata = make(map[string]interface{})
			}
		} else {
			metadata = make(map[string]interface{})
		}
	}
	if !o.Metadata.IsNull() && !o.Metadata.IsUnknown() && o.Metadata.ValueString() != "" {
		err := json.Unmarshal([]byte(o.Metadata.ValueString()), &metadata)
		if err != nil {
			diagnostics.AddError("error parsing custom metadata", fmt.Sprintf("%+v", err))
			metadata = make(map[string]interface{})
		}
	} else {
		metadata = make(map[string]interface{})
	}
	dependencies := make([]cortex.CatalogEntityDependency, len(o.Dependencies))
	for i, dependency := range o.Dependencies {
		dep := CatalogEntityDependencyResourceModel{}
		err := dependency.As(ctx, &dep, defaultObjOptions)
		if err != nil {
			diagnostics.AddError("error parsing dependency", fmt.Sprintf("%+v", err))
		}
		dependencies[i] = dep.ToApiModel(diagnostics)
	}
	alerts := make([]cortex.CatalogEntityAlert, len(o.Alerts))
	for i, alert := range o.Alerts {
		al := CatalogEntityAlertResourceModel{}
		err := alert.As(ctx, &al, defaultObjOptions)
		if err != nil {
			diagnostics.AddError("error parsing alert", fmt.Sprintf("%+v", err))
		}
		alerts[i] = al.ToApiModel()
	}
	dashboards := CatalogEntityDashboardResourceModel{}
	err := o.Dashboards.As(ctx, &dashboards, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing dashboards configuration", fmt.Sprintf("%+v", err))
	}
	apm := CatalogEntityApmResourceModel{}
	err = o.Apm.As(ctx, &apm, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing APM configuration", fmt.Sprintf("%+v", err))
	}
	git := CatalogEntityGitResourceModel{}
	err = o.Git.As(ctx, &git, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing git configuration", fmt.Sprintf("%+v", err))
	}
	issues := CatalogEntityIssuesResourceModel{}
	err = o.Issues.As(ctx, &issues, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing issues configuration", fmt.Sprintf("%+v", err))
	}
	onCall := CatalogEntityOnCallResourceModel{}
	err = o.OnCall.As(ctx, &onCall, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing on-call configuration", fmt.Sprintf("%+v", err))
	}
	serviceLevelObjectives := CatalogEntitySLOsResourceModel{}
	err = o.SLOs.As(ctx, &serviceLevelObjectives, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing SLO configuration", fmt.Sprintf("%+v", err))
	}
	staticAnalysis := CatalogEntityStaticAnalysisResourceModel{}
	err = o.StaticAnalysis.As(ctx, &staticAnalysis, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing static analysis configuration", fmt.Sprintf("%+v", err))
	}
	ciCd := CatalogEntityCiCdResourceModel{}
	err = o.CiCd.As(ctx, &ciCd, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing ci/cd configuration", fmt.Sprintf("%+v", err))
	}
	bugSnag := CatalogEntityBugSnagResourceModel{}
	err = o.BugSnag.As(ctx, &bugSnag, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing bugsnag configuration", fmt.Sprintf("%+v", err))
	}
	checkmarx := CatalogEntityCheckmarxResourceModel{}
	err = o.Checkmarx.As(ctx, &checkmarx, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing Checkmarx configuration", fmt.Sprintf("%+v", err))
	}
	coralogix := CatalogEntityCoralogixResourceModel{}
	err = o.Coralogix.As(ctx, &coralogix, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing Coralogix configuration", fmt.Sprintf("%+v", err))
	}
	firehydrant := CatalogEntityFireHydrantResourceModel{}
	err = o.FireHydrant.As(ctx, &firehydrant, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing FireHydrant configuration", fmt.Sprintf("%+v", err))
	}
	k8s := CatalogEntityK8sResourceModel{}
	err = o.K8s.As(ctx, &k8s, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing K8s configuration", fmt.Sprintf("%+v", err))
	}
	launchDarkly := CatalogEntityLaunchDarklyResourceModel{}
	err = o.LaunchDarkly.As(ctx, &launchDarkly, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing LaunchDarkly configuration", fmt.Sprintf("%+v", err))
	}
	microsoftTeams := make([]cortex.CatalogEntityMicrosoftTeam, len(o.MicrosoftTeams))
	for i, team := range o.MicrosoftTeams {
		t := CatalogEntityMicrosoftTeamResourceModel{}
		err := team.As(ctx, &t, defaultObjOptions)
		if err != nil {
			diagnostics.AddError("error parsing microsoft team", fmt.Sprintf("%+v", err))
		}
		microsoftTeams[i] = t.ToApiModel()
	}
	rollbar := CatalogEntityRollbarResourceModel{}
	err = o.Rollbar.As(ctx, &rollbar, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing Rollbar configuration", fmt.Sprintf("%+v", err))
	}
	sentry := CatalogEntitySentryResourceModel{}
	err = o.Sentry.As(ctx, &sentry, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing Sentry configuration", fmt.Sprintf("%+v", err))
	}
	serviceNow := CatalogEntityServiceNowResourceModel{}
	err = o.ServiceNow.As(ctx, &serviceNow, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing ServiceNow configuration", fmt.Sprintf("%+v", err))
	}
	slack := CatalogEntitySlackResourceModel{}
	err = o.Slack.As(ctx, &slack, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing Slack configuration", fmt.Sprintf("%+v", err))
	}
	snyk := CatalogEntitySnykResourceModel{}
	err = o.Snyk.As(ctx, &snyk, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing Snyk configuration", fmt.Sprintf("%+v", err))
	}
	wiz := CatalogEntityWizResourceModel{}
	err = o.Wiz.As(ctx, &wiz, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing Wiz configuration", fmt.Sprintf("%+v", err))
	}
	team := CatalogEntityTeamResourceModel{}
	err = o.Team.As(ctx, &team, defaultObjOptions)
	if err != nil {
		diagnostics.AddError("error parsing Team configuration", fmt.Sprintf("%+v", err))
	}

	return cortex.CatalogEntityData{
		Tag:            o.Tag.ValueString(),
		Title:          o.Name.ValueString(),
		Description:    o.Description.ValueString(),
		Type:           o.Type.ValueString(),
		Definition:     definition,
		Owners:         owners,
		Children:       children,
		Parents: 	 	Parents,
		Groups:         groups,
		Links:          links,
		IgnoreMetadata: o.IgnoreMetadata.ValueBool(),
		Metadata:       metadata,
		Dependencies:   dependencies,
		Alerts:         alerts,
		Dashboards:     dashboards.ToApiModel(),
		Apm:            apm.ToApiModel(ctx),
		Git:            git.ToApiModel(ctx),
		Issues:         issues.ToApiModel(ctx),
		OnCall:         onCall.ToApiModel(ctx),
		SLOs:           serviceLevelObjectives.ToApiModel(ctx),
		StaticAnalysis: staticAnalysis.ToApiModel(ctx),
		CiCd:           ciCd.ToApiModel(ctx),
		BugSnag:        bugSnag.ToApiModel(),
		Checkmarx:      checkmarx.ToApiModel(),
		Coralogix:      coralogix.ToApiModel(),
		FireHydrant:    firehydrant.ToApiModel(),
		K8s:            k8s.ToApiModel(),
		LaunchDarkly:   launchDarkly.ToApiModel(),
		MicrosoftTeams: microsoftTeams,
		Rollbar:        rollbar.ToApiModel(),
		Sentry:         sentry.ToApiModel(),
		ServiceNow:     serviceNow.ToApiModel(),
		Slack:          slack.ToApiModel(),
		Snyk:           snyk.ToApiModel(),
		Wiz:            wiz.ToApiModel(),
		Team:           team.ToApiModel(),
	}
}

func (o *CatalogEntityResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.CatalogEntityData) {
	o.Id = types.StringValue(entity.Tag)
	o.Name = types.StringValue(entity.Title)
	if entity.Type == "service" || entity.Type == "" {
		o.Type = types.StringNull()
	} else {
		o.Type = types.StringValue(entity.Type)
	}
	if entity.Description != "" {
		o.Description = types.StringValue(entity.Description)
	} else {
		o.Description = types.StringNull()
	}

	// coerce map of unknown types into string
	if entity.Definition != nil && len(entity.Definition) > 0 {
		definition, err := json.Marshal(entity.Definition)
		if err != nil {
			diagnostics.AddError("Error parsing definition: %s", err.Error())
			return
		}
		o.Definition = types.StringValue(string(definition))
	} else {
		o.Definition = types.StringNull()
	}

	if len(entity.Owners) > 0 {
		o.Owners = make([]CatalogEntityOwnerResourceModel, len(entity.Owners))
		for i, owner := range entity.Owners {
			m := CatalogEntityOwnerResourceModel{}
			o.Owners[i] = m.FromApiModel(&owner)
		}
	} else {
		o.Owners = nil
	}

	if len(entity.Children) > 0 {
		o.Children = make([]CatalogEntityChildResourceModel, len(entity.Children))
		for i, child := range entity.Children {
			m := CatalogEntityChildResourceModel{}
			o.Children[i] = m.FromApiModel(&child)
		}
	} else {
		o.Children = nil
	}

	if len(entity.Parents) > 0 {
		o.Parents = make([]CatalogEntityParentResourceModel, len(entity.Parents))
		for i, parent := range entity.Parents {
			m := CatalogEntityParentResourceModel{}
			o.Parents[i] = m.FromApiModel(&parent)
		}
	} else {
		o.Parents = nil
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
		o.Alerts = make([]types.Object, len(entity.Alerts))
		for i, alert := range entity.Alerts {
			m := CatalogEntityAlertResourceModel{}
			o.Alerts[i] = m.FromApiModel(ctx, diagnostics, &alert)
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

	ciCd := CatalogEntityCiCdResourceModel{}
	o.CiCd = ciCd.FromApiModel(ctx, diagnostics, &entity.CiCd)

	bugSnag := CatalogEntityBugSnagResourceModel{}
	o.BugSnag = bugSnag.FromApiModel(ctx, diagnostics, &entity.BugSnag)

	checkmarx := CatalogEntityCheckmarxResourceModel{}
	o.Checkmarx = checkmarx.FromApiModel(ctx, diagnostics, &entity.Checkmarx)

	coralogix := CatalogEntityCoralogixResourceModel{}
	o.Coralogix = coralogix.FromApiModel(ctx, diagnostics, &entity.Coralogix)

	firehydrant := CatalogEntityFireHydrantResourceModel{}
	o.FireHydrant = firehydrant.FromApiModel(ctx, diagnostics, &entity.FireHydrant)

	k8s := CatalogEntityK8sResourceModel{}
	o.K8s = k8s.FromApiModel(ctx, diagnostics, &entity.K8s)

	launchDarkly := CatalogEntityLaunchDarklyResourceModel{}
	o.LaunchDarkly = launchDarkly.FromApiModel(ctx, diagnostics, &entity.LaunchDarkly)

	if len(entity.MicrosoftTeams) > 0 {
		o.MicrosoftTeams = make([]types.Object, len(entity.MicrosoftTeams))
		for i, team := range entity.MicrosoftTeams {
			m := CatalogEntityMicrosoftTeamResourceModel{}
			o.MicrosoftTeams[i] = m.FromApiModel(ctx, diagnostics, &team)
		}
	} else {
		o.MicrosoftTeams = nil
	}

	rollbar := CatalogEntityRollbarResourceModel{}
	o.Rollbar = rollbar.FromApiModel(ctx, diagnostics, &entity.Rollbar)

	sentry := CatalogEntitySentryResourceModel{}
	o.Sentry = sentry.FromApiModel(ctx, diagnostics, &entity.Sentry)

	slack := CatalogEntitySlackResourceModel{}
	o.Slack = slack.FromApiModel(ctx, diagnostics, &entity.Slack)

	serviceNow := CatalogEntityServiceNowResourceModel{}
	o.ServiceNow = serviceNow.FromApiModel(ctx, diagnostics, &entity.ServiceNow)

	snyk := CatalogEntitySnykResourceModel{}
	o.Snyk = snyk.FromApiModel(ctx, diagnostics, &entity.Snyk)

	wiz := CatalogEntityWizResourceModel{}
	o.Wiz = wiz.FromApiModel(ctx, diagnostics, &entity.Wiz)

	team := CatalogEntityTeamResourceModel{}
	o.Team = team.FromApiModel(ctx, diagnostics, &entity.Team)
}

/***********************************************************************************************************************
 * Owners
 ***********************************************************************************************************************/

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
	case "email":
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
	case "email":
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

/***********************************************************************************************************************
 * Children/Parents
 ***********************************************************************************************************************/

// CatalogEntityChildResourceModel describes a child of the catalog entity.
type CatalogEntityChildResourceModel struct {
	Tag types.String `tfsdk:"tag"`
}

func (o *CatalogEntityChildResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag": types.StringType,
	}
}

func (o *CatalogEntityChildResourceModel) ToApiModel() cortex.CatalogEntityChild {
	return cortex.CatalogEntityChild{
		Tag: o.Tag.ValueString(),
	}
}

func (o *CatalogEntityChildResourceModel) FromApiModel(entity *cortex.CatalogEntityChild) CatalogEntityChildResourceModel {
	return CatalogEntityChildResourceModel{
		Tag: types.StringValue(entity.Tag),
	}
}

// CatalogEntityParentResourceModel describes a parent domain of the catalog entity.
type CatalogEntityParentResourceModel struct {
	Tag types.String `tfsdk:"tag"`
}

func (o *CatalogEntityParentResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag": types.StringType,
	}
}

func (o *CatalogEntityParentResourceModel) ToApiModel() cortex.CatalogEntityParent {
	return cortex.CatalogEntityParent{
		Tag: o.Tag.ValueString(),
	}
}

func (o *CatalogEntityParentResourceModel) FromApiModel(entity *cortex.CatalogEntityParent) CatalogEntityParentResourceModel {
	return CatalogEntityParentResourceModel{
		Tag: types.StringValue(entity.Tag),
	}
}

/***********************************************************************************************************************
 * Links
 ***********************************************************************************************************************/

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

/***********************************************************************************************************************
 * Dependencies
 ***********************************************************************************************************************/

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

func (o *CatalogEntityDependencyResourceModel) ToApiModel(diagnostics *diag.Diagnostics) cortex.CatalogEntityDependency {
	metadata := make(map[string]interface{})
	if !o.Metadata.IsNull() && !o.Metadata.IsUnknown() && o.Metadata.ValueString() != "" {
		err := json.Unmarshal([]byte(o.Metadata.ValueString()), &metadata)
		if err != nil {
			diagnostics.AddError("error parsing dependency metadata", fmt.Sprintf("%+v", err))
			metadata = nil
		}
		if len(metadata) == 0 {
			metadata = nil
		}
	} else {
		metadata = nil
	}

	return cortex.CatalogEntityDependency{
		Tag:         o.Tag.ValueString(),
		Method:      o.Method.ValueString(),
		Path:        o.Path.ValueString(),
		Description: o.Description.ValueString(),
		Metadata:    metadata,
	}
}

func (o *CatalogEntityDependencyResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, dependency *cortex.CatalogEntityDependency) types.Object {
	obj := CatalogEntityDependencyResourceModel{
		Tag: types.StringValue(dependency.Tag),
	}
	if dependency.Description != "" {
		obj.Description = types.StringValue(dependency.Description)
	} else {
		obj.Description = types.StringNull()
	}
	if dependency.Path != "" {
		obj.Path = types.StringValue(dependency.Path)
	} else {
		obj.Path = types.StringNull()
	}
	if dependency.Method != "" {
		obj.Method = types.StringValue(dependency.Method)
	} else {
		obj.Method = types.StringNull()
	}
	if dependency.Metadata != nil && len(dependency.Metadata) > 0 {
		depMetadata, err := json.Marshal(dependency.Metadata)
		if err != nil {
			diagnostics.AddError("error marshalling dependency metadata", fmt.Sprintf("%+v", err))
			depMetadata = []byte{}
		}
		obj.Metadata = types.StringValue(string(depMetadata))
	} else {
		obj.Metadata = types.StringNull()
	}

	retObj, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return retObj
}

/***********************************************************************************************************************
 * Alerts
 ***********************************************************************************************************************/

type CatalogEntityAlertResourceModel struct {
	Type  types.String `tfsdk:"type"`
	Tag   types.String `tfsdk:"tag"`
	Value types.String `tfsdk:"value"`
}

func (o *CatalogEntityAlertResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":  types.StringType,
		"tag":   types.StringType,
		"value": types.StringType,
	}
}

func (o *CatalogEntityAlertResourceModel) ToApiModel() cortex.CatalogEntityAlert {
	return cortex.CatalogEntityAlert{
		Type:  o.Type.ValueString(),
		Tag:   o.Tag.ValueString(),
		Value: o.Value.ValueString(),
	}
}

func (o *CatalogEntityAlertResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, alert *cortex.CatalogEntityAlert) types.Object {
	if !alert.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityAlertResourceModel{
		Type:  types.StringValue(alert.Type),
		Tag:   types.StringValue(alert.Tag),
		Value: types.StringValue(alert.Value),
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
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
	Jira types.Object `tfsdk:"jira"`
}

func (o *CatalogEntityIssuesResourceModel) AttrTypes() map[string]attr.Type {
	ob := CatalogEntityIssuesJiraResourceModel{}
	return map[string]attr.Type{
		"jira": types.ObjectType{AttrTypes: ob.AttrTypes()},
	}
}

func (o *CatalogEntityIssuesResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityIssues {
	ob := CatalogEntityIssuesJiraResourceModel{}
	o.Jira.As(ctx, &ob, getDefaultObjectOptions())
	return cortex.CatalogEntityIssues{
		Jira: ob.ToApiModel(),
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

func (o *CatalogEntityIssuesJiraResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityIssuesJira) types.Object {
	obj := CatalogEntityIssuesJiraResourceModel{}

	if !entity.Enabled() {
		return types.ObjectNull(obj.AttrTypes())
	}

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

	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

/***********************************************************************************************************************
 * On-Call
 ***********************************************************************************************************************/

type CatalogEntityOnCallResourceModel struct {
	PagerDuty types.Object `tfsdk:"pager_duty"`
	OpsGenie  types.Object `tfsdk:"ops_genie"`
	VictorOps types.Object `tfsdk:"victor_ops"`
	XMatters  types.Object `tfsdk:"xmatters"`
}

func (o *CatalogEntityOnCallResourceModel) AttrTypes() map[string]attr.Type {
	pd := CatalogEntityOnCallPagerDutyResourceModel{}
	og := CatalogEntityOnCallOpsGenieResourceModel{}
	vo := CatalogEntityOnCallVictorOpsResourceModel{}
	xm := CatalogEntityOnCallXMattersResourceModel{}
	return map[string]attr.Type{
		"pager_duty": types.ObjectType{AttrTypes: pd.AttrTypes()},
		"ops_genie":  types.ObjectType{AttrTypes: og.AttrTypes()},
		"victor_ops": types.ObjectType{AttrTypes: vo.AttrTypes()},
		"xmatters":   types.ObjectType{AttrTypes: xm.AttrTypes()},
	}
}

func (o *CatalogEntityOnCallResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityOnCall {
	defaultObjOptions := getDefaultObjectOptions()

	pd := CatalogEntityOnCallPagerDutyResourceModel{}
	o.PagerDuty.As(ctx, &pd, defaultObjOptions)

	og := CatalogEntityOnCallOpsGenieResourceModel{}
	o.OpsGenie.As(ctx, &og, defaultObjOptions)

	vo := CatalogEntityOnCallVictorOpsResourceModel{}
	o.VictorOps.As(ctx, &vo, defaultObjOptions)

	xm := CatalogEntityOnCallXMattersResourceModel{}
	o.XMatters.As(ctx, &xm, defaultObjOptions)

	return cortex.CatalogEntityOnCall{
		PagerDuty: pd.ToApiModel(),
		OpsGenie:  og.ToApiModel(),
		VictorOps: vo.ToApiModel(),
		XMatters:  xm.ToApiModel(),
	}
}

func (o *CatalogEntityOnCallResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, onCall *cortex.CatalogEntityOnCall) types.Object {
	if !onCall.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	pd := CatalogEntityOnCallPagerDutyResourceModel{}
	og := CatalogEntityOnCallOpsGenieResourceModel{}
	vo := CatalogEntityOnCallVictorOpsResourceModel{}
	xm := CatalogEntityOnCallXMattersResourceModel{}

	ob := CatalogEntityOnCallResourceModel{
		PagerDuty: pd.FromApiModel(ctx, diagnostics, &onCall.PagerDuty),
		OpsGenie:  og.FromApiModel(ctx, diagnostics, &onCall.OpsGenie),
		VictorOps: vo.FromApiModel(ctx, diagnostics, &onCall.VictorOps),
		XMatters:  xm.FromApiModel(ctx, diagnostics, &onCall.XMatters),
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

func (o *CatalogEntityOnCallPagerDutyResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityOnCallPagerDuty) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityOnCallPagerDutyResourceModel{
		ID:   types.StringValue(entity.ID),
		Type: types.StringValue(entity.Type),
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
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

func (o *CatalogEntityOnCallOpsGenieResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityOnCallOpsGenie) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityOnCallOpsGenieResourceModel{
		ID:   types.StringValue(entity.ID),
		Type: types.StringValue(entity.Type),
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
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

func (o *CatalogEntityOnCallVictorOpsResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityOnCallVictorOps) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityOnCallVictorOpsResourceModel{
		ID:   types.StringValue(entity.ID),
		Type: types.StringValue(entity.Type),
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

// XMatters

type CatalogEntityOnCallXMattersResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (o *CatalogEntityOnCallXMattersResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
}

func (o *CatalogEntityOnCallXMattersResourceModel) ToApiModel() cortex.CatalogEntityOnCallXMatters {
	return cortex.CatalogEntityOnCallXMatters{
		ID:   o.ID.ValueString(),
		Type: o.Type.ValueString(),
	}
}

func (o *CatalogEntityOnCallXMattersResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityOnCallXMatters) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityOnCallXMattersResourceModel{
		ID:   types.StringValue(entity.ID),
		Type: types.StringValue(entity.Type),
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
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
 * Coralogix
 **********************************************************************************************************************/

type CatalogEntityCoralogixResourceModel struct {
	Applications []CatalogEntityCoralogixApplicationResourceModel `tfsdk:"applications"`
}

func (o *CatalogEntityCoralogixResourceModel) AttrTypes() map[string]attr.Type {
	ob := CatalogEntityCoralogixApplicationResourceModel{}
	return map[string]attr.Type{
		"applications": types.ListType{ElemType: types.ObjectType{AttrTypes: ob.AttrTypes()}},
	}
}

func (o *CatalogEntityCoralogixResourceModel) ToApiModel() cortex.CatalogEntityCoralogix {
	applications := make([]cortex.CatalogEntityCoralogixApplication, len(o.Applications))
	for i, p := range o.Applications {
		applications[i] = p.ToApiModel()
	}
	return cortex.CatalogEntityCoralogix{
		Applications: applications,
	}
}

func (o *CatalogEntityCoralogixResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityCoralogix) types.Object {
	obj := CatalogEntityCoralogixResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(obj.AttrTypes())
	}

	applications := make([]CatalogEntityCoralogixApplicationResourceModel, len(entity.Applications))
	for i, p := range entity.Applications {
		po := CatalogEntityCoralogixApplicationResourceModel{}
		applications[i] = po.FromApiModel(p)
	}
	obj.Applications = applications

	objectValue, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue

}

type CatalogEntityCoralogixApplicationResourceModel struct {
	Name  types.String `tfsdk:"name"`
	Alias types.String `tfsdk:"alias"`
}

func (o *CatalogEntityCoralogixApplicationResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"alias": types.StringType,
	}
}

func (o *CatalogEntityCoralogixApplicationResourceModel) ToApiModel() cortex.CatalogEntityCoralogixApplication {
	entity := cortex.CatalogEntityCoralogixApplication{}
	entity.Name = o.Name.ValueString()
	if o.Alias.ValueString() != "" {
		entity.Alias = o.Alias.ValueString()
	}
	return entity
}

func (o *CatalogEntityCoralogixApplicationResourceModel) FromApiModel(application cortex.CatalogEntityCoralogixApplication) CatalogEntityCoralogixApplicationResourceModel {
	ob := CatalogEntityCoralogixApplicationResourceModel{}
	ob.Name = types.StringValue(application.Name)
	if application.Alias == "" {
		ob.Alias = types.StringNull()
	} else {
		ob.Alias = types.StringValue(application.Alias)
	}
	return ob
}

/***********************************************************************************************************************
 * FireHydrant
 **********************************************************************************************************************/

type CatalogEntityFireHydrantResourceModel struct {
	Services []CatalogEntityFireHydrantServiceResourceModel `tfsdk:"services"`
}

func (o *CatalogEntityFireHydrantResourceModel) AttrTypes() map[string]attr.Type {
	ob := CatalogEntityFireHydrantServiceResourceModel{}
	return map[string]attr.Type{
		"services": types.ListType{ElemType: types.ObjectType{AttrTypes: ob.AttrTypes()}},
	}
}

func (o *CatalogEntityFireHydrantResourceModel) ToApiModel() cortex.CatalogEntityFireHydrant {
	services := make([]cortex.CatalogEntityFireHydrantService, len(o.Services))
	for i, s := range o.Services {
		services[i] = s.ToApiModel()
	}
	return cortex.CatalogEntityFireHydrant{
		Services: services,
	}
}

func (o *CatalogEntityFireHydrantResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityFireHydrant) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	obj := &CatalogEntityFireHydrantResourceModel{}

	services := make([]CatalogEntityFireHydrantServiceResourceModel, len(entity.Services))
	for i, s := range entity.Services {
		so := CatalogEntityFireHydrantServiceResourceModel{}
		services[i] = so.FromApiModel(s)
	}
	obj.Services = services

	objectValue, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

type CatalogEntityFireHydrantServiceResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (o *CatalogEntityFireHydrantServiceResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
}

func (o *CatalogEntityFireHydrantServiceResourceModel) ToApiModel() cortex.CatalogEntityFireHydrantService {
	entity := cortex.CatalogEntityFireHydrantService{
		ID:   o.ID.ValueString(),
		Type: o.Type.ValueString(),
	}
	return entity
}

func (o *CatalogEntityFireHydrantServiceResourceModel) FromApiModel(service cortex.CatalogEntityFireHydrantService) CatalogEntityFireHydrantServiceResourceModel {
	return CatalogEntityFireHydrantServiceResourceModel{
		ID:   types.StringValue(service.ID),
		Type: types.StringValue(service.Type),
	}
}

/***********************************************************************************************************************
 * Kubernetes
 **********************************************************************************************************************/

type CatalogEntityK8sResourceModel struct {
	Deployments  []CatalogEntityK8sDeploymentResourceModel  `tfsdk:"deployments"`
	ArgoRollouts []CatalogEntityK8sArgoRolloutResourceModel `tfsdk:"argo_rollouts"`
	StatefulSets []CatalogEntityK8sStatefulSetResourceModel `tfsdk:"stateful_sets"`
	CronJobs     []CatalogEntityK8sCronJobResourceModel     `tfsdk:"cron_jobs"`
}

func (o *CatalogEntityK8sResourceModel) AttrTypes() map[string]attr.Type {
	de := CatalogEntityK8sDeploymentResourceModel{}
	ar := CatalogEntityK8sArgoRolloutResourceModel{}
	ss := CatalogEntityK8sStatefulSetResourceModel{}
	cj := CatalogEntityK8sCronJobResourceModel{}
	return map[string]attr.Type{
		"deployments":   types.ListType{ElemType: types.ObjectType{AttrTypes: de.AttrTypes()}},
		"argo_rollouts": types.ListType{ElemType: types.ObjectType{AttrTypes: ar.AttrTypes()}},
		"stateful_sets": types.ListType{ElemType: types.ObjectType{AttrTypes: ss.AttrTypes()}},
		"cron_jobs":     types.ListType{ElemType: types.ObjectType{AttrTypes: cj.AttrTypes()}},
	}
}

func (o *CatalogEntityK8sResourceModel) ToApiModel() cortex.CatalogEntityK8s {
	deployments := make([]cortex.CatalogEntityK8sDeployment, len(o.Deployments))
	for i, c := range o.Deployments {
		deployments[i] = c.ToApiModel()
	}
	argoRollouts := make([]cortex.CatalogEntityK8sArgoRollout, len(o.ArgoRollouts))
	for i, c := range o.ArgoRollouts {
		argoRollouts[i] = c.ToApiModel()
	}
	statefulSets := make([]cortex.CatalogEntityK8sStatefulSet, len(o.StatefulSets))
	for i, c := range o.StatefulSets {
		statefulSets[i] = c.ToApiModel()
	}
	cronJobs := make([]cortex.CatalogEntityK8sCronJob, len(o.CronJobs))
	for i, c := range o.CronJobs {
		cronJobs[i] = c.ToApiModel()
	}

	return cortex.CatalogEntityK8s{
		Deployments:  deployments,
		ArgoRollouts: argoRollouts,
		StatefulSets: statefulSets,
		CronJobs:     cronJobs,
	}
}

func (o *CatalogEntityK8sResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityK8s) types.Object {
	ob := CatalogEntityK8sResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(ob.AttrTypes())
	}

	var deployments = make([]CatalogEntityK8sDeploymentResourceModel, len(entity.Deployments))
	for i, e := range entity.Deployments {
		ch := CatalogEntityK8sDeploymentResourceModel{}
		deployments[i] = ch.FromApiModel(&e)
	}
	ob.Deployments = deployments

	var argoRollouts = make([]CatalogEntityK8sArgoRolloutResourceModel, len(entity.ArgoRollouts))
	for i, e := range entity.ArgoRollouts {
		ch := CatalogEntityK8sArgoRolloutResourceModel{}
		argoRollouts[i] = ch.FromApiModel(&e)
	}
	ob.ArgoRollouts = argoRollouts

	var statefulSets = make([]CatalogEntityK8sStatefulSetResourceModel, len(entity.StatefulSets))
	for i, e := range entity.StatefulSets {
		ch := CatalogEntityK8sStatefulSetResourceModel{}
		statefulSets[i] = ch.FromApiModel(&e)
	}
	ob.StatefulSets = statefulSets

	var cronJobs = make([]CatalogEntityK8sCronJobResourceModel, len(entity.CronJobs))
	for i, e := range entity.CronJobs {
		ch := CatalogEntityK8sCronJobResourceModel{}
		cronJobs[i] = ch.FromApiModel(&e)
	}
	ob.CronJobs = cronJobs

	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

// Deployment

type CatalogEntityK8sDeploymentResourceModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Cluster    types.String `tfsdk:"cluster"`
}

func (o *CatalogEntityK8sDeploymentResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"identifier": types.StringType,
		"cluster":    types.StringType,
	}
}

func (o *CatalogEntityK8sDeploymentResourceModel) ToApiModel() cortex.CatalogEntityK8sDeployment {
	return cortex.CatalogEntityK8sDeployment{
		Identifier: o.Identifier.ValueString(),
		Cluster:    o.Cluster.ValueString(),
	}
}

func (o *CatalogEntityK8sDeploymentResourceModel) FromApiModel(entity *cortex.CatalogEntityK8sDeployment) CatalogEntityK8sDeploymentResourceModel {
	return CatalogEntityK8sDeploymentResourceModel{
		Identifier: types.StringValue(entity.Identifier),
		Cluster:    types.StringValue(entity.Cluster),
	}
}

// ArgoRollout

type CatalogEntityK8sArgoRolloutResourceModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Cluster    types.String `tfsdk:"cluster"`
}

func (o *CatalogEntityK8sArgoRolloutResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"identifier": types.StringType,
		"cluster":    types.StringType,
	}
}

func (o *CatalogEntityK8sArgoRolloutResourceModel) ToApiModel() cortex.CatalogEntityK8sArgoRollout {
	return cortex.CatalogEntityK8sArgoRollout{
		Identifier: o.Identifier.ValueString(),
		Cluster:    o.Cluster.ValueString(),
	}
}

func (o *CatalogEntityK8sArgoRolloutResourceModel) FromApiModel(entity *cortex.CatalogEntityK8sArgoRollout) CatalogEntityK8sArgoRolloutResourceModel {
	return CatalogEntityK8sArgoRolloutResourceModel{
		Identifier: types.StringValue(entity.Identifier),
		Cluster:    types.StringValue(entity.Cluster),
	}
}

// StatefulSet

type CatalogEntityK8sStatefulSetResourceModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Cluster    types.String `tfsdk:"cluster"`
}

func (o *CatalogEntityK8sStatefulSetResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"identifier": types.StringType,
		"cluster":    types.StringType,
	}
}

func (o *CatalogEntityK8sStatefulSetResourceModel) ToApiModel() cortex.CatalogEntityK8sStatefulSet {
	return cortex.CatalogEntityK8sStatefulSet{
		Identifier: o.Identifier.ValueString(),
		Cluster:    o.Cluster.ValueString(),
	}
}

func (o *CatalogEntityK8sStatefulSetResourceModel) FromApiModel(entity *cortex.CatalogEntityK8sStatefulSet) CatalogEntityK8sStatefulSetResourceModel {
	return CatalogEntityK8sStatefulSetResourceModel{
		Identifier: types.StringValue(entity.Identifier),
		Cluster:    types.StringValue(entity.Cluster),
	}
}

// CronJob

type CatalogEntityK8sCronJobResourceModel struct {
	Identifier types.String `tfsdk:"identifier"`
	Cluster    types.String `tfsdk:"cluster"`
}

func (o *CatalogEntityK8sCronJobResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"identifier": types.StringType,
		"cluster":    types.StringType,
	}
}

func (o *CatalogEntityK8sCronJobResourceModel) ToApiModel() cortex.CatalogEntityK8sCronJob {
	return cortex.CatalogEntityK8sCronJob{
		Identifier: o.Identifier.ValueString(),
		Cluster:    o.Cluster.ValueString(),
	}
}

func (o *CatalogEntityK8sCronJobResourceModel) FromApiModel(entity *cortex.CatalogEntityK8sCronJob) CatalogEntityK8sCronJobResourceModel {
	return CatalogEntityK8sCronJobResourceModel{
		Identifier: types.StringValue(entity.Identifier),
		Cluster:    types.StringValue(entity.Cluster),
	}
}

/***********************************************************************************************************************
 * LaunchDarkly
 **********************************************************************************************************************/

type CatalogEntityLaunchDarklyResourceModel struct {
	Projects []CatalogEntityLaunchDarklyProjectResourceModel `tfsdk:"projects"`
}

func (o *CatalogEntityLaunchDarklyResourceModel) AttrTypes() map[string]attr.Type {
	pp := CatalogEntityLaunchDarklyProjectResourceModel{}
	return map[string]attr.Type{
		"projects": types.ListType{ElemType: types.ObjectType{AttrTypes: pp.AttrTypes()}},
	}
}

func (o *CatalogEntityLaunchDarklyResourceModel) ToApiModel() cortex.CatalogEntityLaunchDarkly {
	projects := make([]cortex.CatalogEntityLaunchDarklyProject, len(o.Projects))
	for i, c := range o.Projects {
		projects[i] = c.ToApiModel()
	}
	return cortex.CatalogEntityLaunchDarkly{
		Projects: projects,
	}
}

func (o *CatalogEntityLaunchDarklyResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityLaunchDarkly) types.Object {
	ob := CatalogEntityLaunchDarklyResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(ob.AttrTypes())
	}

	var projects = make([]CatalogEntityLaunchDarklyProjectResourceModel, len(entity.Projects))
	for i, e := range entity.Projects {
		ch := CatalogEntityLaunchDarklyProjectResourceModel{}
		projects[i] = ch.FromApiModel(&e)
	}
	ob.Projects = projects

	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

// Projects

type CatalogEntityLaunchDarklyProjectResourceModel struct {
	ID           types.String                                               `tfsdk:"id"`
	Type         types.String                                               `tfsdk:"type"`
	Alias        types.String                                               `tfsdk:"alias"`
	Environments []CatalogEntityLaunchDarklyProjectEnvironmentResourceModel `tfsdk:"environments"`
}

func (o *CatalogEntityLaunchDarklyProjectResourceModel) AttrTypes() map[string]attr.Type {
	ob := CatalogEntityLaunchDarklyProjectEnvironmentResourceModel{}
	return map[string]attr.Type{
		"id":           types.StringType,
		"type":         types.StringType,
		"alias":        types.StringType,
		"environments": types.ListType{ElemType: types.ObjectType{AttrTypes: ob.AttrTypes()}},
	}
}

func (o *CatalogEntityLaunchDarklyProjectResourceModel) ToApiModel() cortex.CatalogEntityLaunchDarklyProject {
	environments := make([]cortex.CatalogEntityLaunchDarklyProjectEnvironment, len(o.Environments))
	for i, c := range o.Environments {
		environments[i] = c.ToApiModel()
	}
	return cortex.CatalogEntityLaunchDarklyProject{
		ID:           o.ID.ValueString(),
		Type:         o.Type.ValueString(),
		Alias:        o.Alias.ValueString(),
		Environments: environments,
	}
}

func (o *CatalogEntityLaunchDarklyProjectResourceModel) FromApiModel(entity *cortex.CatalogEntityLaunchDarklyProject) CatalogEntityLaunchDarklyProjectResourceModel {
	ob := CatalogEntityLaunchDarklyProjectResourceModel{
		ID:   types.StringValue(entity.ID),
		Type: types.StringValue(entity.Type),
	}
	ob.Environments = make([]CatalogEntityLaunchDarklyProjectEnvironmentResourceModel, len(entity.Environments))
	for i, e := range entity.Environments {
		ch := CatalogEntityLaunchDarklyProjectEnvironmentResourceModel{}
		ob.Environments[i] = ch.FromApiModel(&e)
	}
	if entity.Alias != "" {
		ob.Alias = types.StringValue(entity.Alias)
	} else {
		ob.Alias = types.StringNull()
	}
	return ob
}

// Project Environments

type CatalogEntityLaunchDarklyProjectEnvironmentResourceModel struct {
	Name types.String `tfsdk:"name"`
}

func (o *CatalogEntityLaunchDarklyProjectEnvironmentResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}

func (o *CatalogEntityLaunchDarklyProjectEnvironmentResourceModel) ToApiModel() cortex.CatalogEntityLaunchDarklyProjectEnvironment {
	return cortex.CatalogEntityLaunchDarklyProjectEnvironment{
		Name: o.Name.ValueString(),
	}
}

func (o *CatalogEntityLaunchDarklyProjectEnvironmentResourceModel) FromApiModel(entity *cortex.CatalogEntityLaunchDarklyProjectEnvironment) CatalogEntityLaunchDarklyProjectEnvironmentResourceModel {
	return CatalogEntityLaunchDarklyProjectEnvironmentResourceModel{
		Name: types.StringValue(entity.Name),
	}
}

/***********************************************************************************************************************
 * Microsoft Teams
 **********************************************************************************************************************/

type CatalogEntityMicrosoftTeamResourceModel struct {
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

func (o *CatalogEntityMicrosoftTeamResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                  types.StringType,
		"description":           types.StringType,
		"notifications_enabled": types.BoolType,
	}
}

func (o *CatalogEntityMicrosoftTeamResourceModel) ToApiModel() cortex.CatalogEntityMicrosoftTeam {
	return cortex.CatalogEntityMicrosoftTeam{
		Name:                 o.Name.ValueString(),
		Description:          o.Description.ValueString(),
		NotificationsEnabled: o.NotificationsEnabled.ValueBool(),
	}
}

func (o *CatalogEntityMicrosoftTeamResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityMicrosoftTeam) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityMicrosoftTeamResourceModel{
		Name:                 types.StringValue(entity.Name),
		NotificationsEnabled: types.BoolValue(entity.NotificationsEnabled),
	}
	if entity.Description != "" {
		ob.Description = types.StringValue(entity.Description)
	} else {
		ob.Description = types.StringNull()
	}
	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
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
 * ServiceNow
 **********************************************************************************************************************/

type CatalogEntityServiceNowResourceModel struct {
	Services []CatalogEntityServiceNowServiceResourceModel `tfsdk:"services"`
}

func (o *CatalogEntityServiceNowResourceModel) AttrTypes() map[string]attr.Type {
	ob := CatalogEntityServiceNowServiceResourceModel{}
	return map[string]attr.Type{
		"services": types.ListType{ElemType: types.ObjectType{AttrTypes: ob.AttrTypes()}},
	}
}

func (o *CatalogEntityServiceNowResourceModel) ToApiModel() cortex.CatalogEntityServiceNow {
	services := make([]cortex.CatalogEntityServiceNowService, len(o.Services))
	for i, c := range o.Services {
		services[i] = c.ToApiModel()
	}
	return cortex.CatalogEntityServiceNow{
		Services: services,
	}
}

func (o *CatalogEntityServiceNowResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityServiceNow) types.Object {
	ob := CatalogEntityServiceNowResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(ob.AttrTypes())
	}

	var services = make([]CatalogEntityServiceNowServiceResourceModel, len(entity.Services))
	for i, e := range entity.Services {
		ch := CatalogEntityServiceNowServiceResourceModel{}
		services[i] = ch.FromApiModel(&e)
	}
	ob.Services = services

	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

type CatalogEntityServiceNowServiceResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	TableName types.String `tfsdk:"table_name"`
}

func (o *CatalogEntityServiceNowServiceResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.Int64Type,
		"table_name": types.StringType,
	}
}

func (o *CatalogEntityServiceNowServiceResourceModel) ToApiModel() cortex.CatalogEntityServiceNowService {
	return cortex.CatalogEntityServiceNowService{
		ID:        o.ID.ValueInt64(),
		TableName: o.TableName.ValueString(),
	}
}

func (o *CatalogEntityServiceNowServiceResourceModel) FromApiModel(entity *cortex.CatalogEntityServiceNowService) CatalogEntityServiceNowServiceResourceModel {
	return CatalogEntityServiceNowServiceResourceModel{
		ID:        types.Int64Value(entity.ID),
		TableName: types.StringValue(entity.TableName),
	}
}

/***********************************************************************************************************************
 * Sentry
 **********************************************************************************************************************/

type CatalogEntitySlackResourceModel struct {
	Channels []CatalogEntitySlackChannelResourceModel `tfsdk:"channels"`
}

func (o *CatalogEntitySlackResourceModel) AttrTypes() map[string]attr.Type {
	obc := CatalogEntitySlackChannelResourceModel{}
	return map[string]attr.Type{
		"channels": types.SetType{ElemType: types.ObjectType{AttrTypes: obc.AttrTypes()}},
	}
}

func (o *CatalogEntitySlackResourceModel) ToApiModel() cortex.CatalogEntitySlack {
	channels := make([]cortex.CatalogEntitySlackChannel, len(o.Channels))
	for i, c := range o.Channels {
		channels[i] = c.ToApiModel()
	}
	return cortex.CatalogEntitySlack{
		Channels: channels,
	}
}

func (o *CatalogEntitySlackResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntitySlack) types.Object {
	ob := CatalogEntitySlackResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(ob.AttrTypes())
	}

	var channels = make([]CatalogEntitySlackChannelResourceModel, len(entity.Channels))
	for i, e := range entity.Channels {
		ch := CatalogEntitySlackChannelResourceModel{}
		channels[i] = ch.FromApiModel(&e)
	}
	ob.Channels = channels

	obj, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return obj
}

type CatalogEntitySlackChannelResourceModel struct {
	Name                 types.String `tfsdk:"name"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

func (o *CatalogEntitySlackChannelResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                  types.StringType,
		"notifications_enabled": types.BoolType,
	}
}

func (o *CatalogEntitySlackChannelResourceModel) ToApiModel() cortex.CatalogEntitySlackChannel {
	return cortex.CatalogEntitySlackChannel{
		Name:                 o.Name.ValueString(),
		NotificationsEnabled: o.NotificationsEnabled.ValueBool(),
	}
}

func (o *CatalogEntitySlackChannelResourceModel) FromApiModel(entity *cortex.CatalogEntitySlackChannel) CatalogEntitySlackChannelResourceModel {
	return CatalogEntitySlackChannelResourceModel{
		Name:                 types.StringValue(entity.Name),
		NotificationsEnabled: types.BoolValue(entity.NotificationsEnabled),
	}
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
 * Wiz
 **********************************************************************************************************************/

type CatalogEntityWizResourceModel struct {
	Projects []CatalogEntityWizProjectResourceModel `tfsdk:"projects"`
}

func (o *CatalogEntityWizResourceModel) AttrTypes() map[string]attr.Type {
	obp := CatalogEntityWizProjectResourceModel{}
	return map[string]attr.Type{
		"projects": types.ListType{
			ElemType: types.ObjectType{AttrTypes: obp.AttrTypes()},
		},
	}
}

func (o *CatalogEntityWizResourceModel) ToApiModel() cortex.CatalogEntityWiz {
	var projects = make([]cortex.CatalogEntityWizProject, len(o.Projects))
	for i, e := range o.Projects {
		projects[i] = e.ToApiModel()
	}
	return cortex.CatalogEntityWiz{
		Projects: projects,
	}
}

func (o *CatalogEntityWizResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityWiz) types.Object {
	obj := &CatalogEntityWizResourceModel{}
	if !entity.Enabled() {
		return types.ObjectNull(obj.AttrTypes())
	}

	var projects = make([]CatalogEntityWizProjectResourceModel, len(entity.Projects))
	for i, e := range entity.Projects {
		ob := CatalogEntityWizProjectResourceModel{}
		projects[i] = ob.FromApiModel(&e)
	}
	obj.Projects = projects

	objectValue, d := types.ObjectValueFrom(ctx, o.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

type CatalogEntityWizProjectResourceModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

func (o *CatalogEntityWizProjectResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project_id": types.StringType,
	}
}

func (o *CatalogEntityWizProjectResourceModel) ToApiModel() cortex.CatalogEntityWizProject {
	return cortex.CatalogEntityWizProject{
		ProjectID: o.ProjectID.ValueString(),
	}
}

func (o *CatalogEntityWizProjectResourceModel) FromApiModel(entity *cortex.CatalogEntityWizProject) CatalogEntityWizProjectResourceModel {
	return CatalogEntityWizProjectResourceModel{
		ProjectID: types.StringValue(entity.ProjectID),
	}
}

/***********************************************************************************************************************
 * APM Configuration
 **********************************************************************************************************************/

type CatalogEntityApmResourceModel struct {
	DataDog   types.Object   `tfsdk:"data_dog"`
	Dynatrace types.Object   `tfsdk:"dynatrace"`
	NewRelic  []types.Object `tfsdk:"new_relic"`
}

func (o *CatalogEntityApmResourceModel) AttrTypes() map[string]attr.Type {
	dd := CatalogEntityApmDataDogResourceModel{}
	dt := CatalogEntityApmDynatraceResourceModel{}
	nr := CatalogEntityApmNewRelicResourceModel{}

	return map[string]attr.Type{
		"data_dog":  types.ObjectType{AttrTypes: dd.AttrTypes()},
		"dynatrace": types.ObjectType{AttrTypes: dt.AttrTypes()},
		"new_relic": types.ListType{ElemType: types.ObjectType{AttrTypes: nr.AttrTypes()}},
	}
}

func (o *CatalogEntityApmResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityApm {
	defaultObjOptions := getDefaultObjectOptions()

	dd := CatalogEntityApmDataDogResourceModel{}
	o.DataDog.As(ctx, &dd, defaultObjOptions)

	dt := CatalogEntityApmDynatraceResourceModel{}
	o.Dynatrace.As(ctx, &dt, defaultObjOptions)

	entity := cortex.CatalogEntityApm{
		DataDog:   dd.ToApiModel(),
		Dynatrace: dt.ToApiModel(),
	}
	entity.NewRelic = make([]cortex.CatalogEntityApmNewRelic, len(o.NewRelic))
	for i, app := range o.NewRelic {
		nob := CatalogEntityApmNewRelicResourceModel{}
		app.As(ctx, &nob, defaultObjOptions)
		entity.NewRelic[i] = nob.ToApiModel()
	}
	return entity
}

func (o *CatalogEntityApmResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityApm) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	dd := CatalogEntityApmDataDogResourceModel{}
	dt := CatalogEntityApmDynatraceResourceModel{}
	obj := CatalogEntityApmResourceModel{
		DataDog:   dd.FromApiModel(ctx, diagnostics, &entity.DataDog),
		Dynatrace: dt.FromApiModel(ctx, diagnostics, &entity.Dynatrace),
	}
	//nr := CatalogEntityApmNewRelicResourceModel{}
	if len(entity.NewRelic) > 0 {
		obj.NewRelic = make([]types.Object, len(entity.NewRelic))
		for i, e := range entity.NewRelic {
			nob := CatalogEntityApmNewRelicResourceModel{}
			obj.NewRelic[i] = nob.FromApiModel(ctx, diagnostics, &e)
		}
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

func (o *CatalogEntityApmDataDogResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityApmDataDog) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	monitors := slices.Reject(entity.Monitors, func(i int64) bool { return i == 0 })
	monitorList, d := types.ListValueFrom(ctx, types.Int64Type, monitors)
	diagnostics.Append(d...)

	obj := CatalogEntityApmDataDogResourceModel{
		Monitors: monitorList,
	}
	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
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

func (o *CatalogEntityApmDynatraceResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityApmDynatrace) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

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

	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
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

func (o *CatalogEntityApmNewRelicResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityApmNewRelic) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

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

	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
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
	Lightstep  types.Object                              `tfsdk:"lightstep"`
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

func (o *CatalogEntitySLOsResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntitySLOs {
	var dataDog = make([]cortex.CatalogEntitySLODataDog, len(o.DataDog))
	for i, e := range o.DataDog {
		dataDog[i] = e.ToApiModel()
	}
	var dynatrace = make([]cortex.CatalogEntitySLODynatrace, len(o.Dynatrace))
	for i, e := range o.Dynatrace {
		dynatrace[i] = e.ToApiModel()
	}

	lightstep := CatalogEntitySLOLightstepResourceModel{}
	o.Lightstep.As(ctx, &lightstep, getDefaultObjectOptions())

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
		Lightstep:  lightstep.ToApiModel(ctx),
		Prometheus: prometheusQueries,
		SignalFX:   signalFx,
		SumoLogic:  sumoLogic,
	}
}

func (o *CatalogEntitySLOsResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntitySLOs) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	obj := CatalogEntitySLOsResourceModel{}
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

	ls := CatalogEntitySLOLightstepResourceModel{}
	obj.Lightstep = ls.FromApiModel(ctx, diagnostics, entity.Lightstep)

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
	Streams []types.Object `tfsdk:"streams"`
}

func (o *CatalogEntitySLOLightstepResourceModel) AttrTypes() map[string]attr.Type {
	s := CatalogEntitySLOLightstepStreamResourceModel{}
	return map[string]attr.Type{
		"streams": types.ListType{ElemType: types.ObjectType{AttrTypes: s.AttrTypes()}},
	}
}

func (o *CatalogEntitySLOLightstepResourceModel) ToApiModel(ctx context.Context) []cortex.CatalogEntitySLOLightstepStream {
	defaultObjectOptions := getDefaultObjectOptions()
	var arr = make([]cortex.CatalogEntitySLOLightstepStream, len(o.Streams))
	for i, stream := range o.Streams {
		lsr := CatalogEntitySLOLightstepStreamResourceModel{}
		stream.As(ctx, &lsr, defaultObjectOptions)
		arr[i] = lsr.ToApiModel(ctx)
	}
	return arr
}

func (o *CatalogEntitySLOLightstepResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, streams []cortex.CatalogEntitySLOLightstepStream) types.Object {
	if len(streams) == 0 {
		return types.ObjectNull(o.AttrTypes())
	}

	var arr = make([]types.Object, len(streams))
	for i, e := range streams {
		m := CatalogEntitySLOLightstepStreamResourceModel{}
		arr[i] = m.FromApiModel(ctx, diagnostics, &e)
	}
	obj := CatalogEntitySLOLightstepResourceModel{
		Streams: arr,
	}
	objectValue, d := types.ObjectValueFrom(ctx, obj.AttrTypes(), &obj)
	diagnostics.Append(d...)
	return objectValue
}

type CatalogEntitySLOLightstepStreamResourceModel struct {
	StreamID types.String `tfsdk:"stream_id"`
	Targets  types.Object `tfsdk:"targets"`
}

func (o *CatalogEntitySLOLightstepStreamResourceModel) AttrTypes() map[string]attr.Type {
	t := CatalogEntitySLOLightstepStreamTargetResourceModel{}
	return map[string]attr.Type{
		"stream_id": types.StringType,
		"targets":   types.ObjectType{AttrTypes: t.AttrTypes()},
	}
}

func (o *CatalogEntitySLOLightstepStreamResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntitySLOLightstepStream {
	ts := CatalogEntitySLOLightstepStreamTargetResourceModel{}
	o.Targets.As(ctx, &ts, getDefaultObjectOptions())

	return cortex.CatalogEntitySLOLightstepStream{
		StreamID: o.StreamID.ValueString(),
		Targets:  ts.ToApiModel(ctx),
	}
}

func (o *CatalogEntitySLOLightstepStreamResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntitySLOLightstepStream) types.Object {
	ts := CatalogEntitySLOLightstepStreamTargetResourceModel{}

	ob := CatalogEntitySLOLightstepStreamResourceModel{
		StreamID: types.StringValue(entity.StreamID),
		Targets:  ts.FromApiModel(ctx, diagnostics, &entity.Targets),
	}
	objectType, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return objectType
}

type CatalogEntitySLOLightstepStreamTargetResourceModel struct {
	Latencies []types.Object `tfsdk:"latencies"`
}

func (o *CatalogEntitySLOLightstepStreamTargetResourceModel) AttrTypes() map[string]attr.Type {
	l := CatalogEntitySLOLightstepStreamTargetLatencyResourceModel{}
	return map[string]attr.Type{
		"latencies": types.ListType{ElemType: types.ObjectType{AttrTypes: l.AttrTypes()}},
	}
}

func (o *CatalogEntitySLOLightstepStreamTargetResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntitySLOLightstepTargets {
	defaultObjectOptions := getDefaultObjectOptions()

	var latencies = make([]cortex.CatalogEntitySLOLightstepTargetLatency, len(o.Latencies))
	for i, e := range o.Latencies {
		ob := CatalogEntitySLOLightstepStreamTargetLatencyResourceModel{}
		e.As(ctx, &ob, defaultObjectOptions)
		latencies[i] = ob.ToApiModel()
	}

	return cortex.CatalogEntitySLOLightstepTargets{
		Latencies: latencies,
	}
}

func (o *CatalogEntitySLOLightstepStreamTargetResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntitySLOLightstepTargets) types.Object {
	latencies := make([]types.Object, len(entity.Latencies))
	for i, e := range entity.Latencies {
		m := CatalogEntitySLOLightstepStreamTargetLatencyResourceModel{}
		latencies[i] = m.FromApiModel(ctx, diagnostics, &e)
	}
	ob := CatalogEntitySLOLightstepStreamTargetResourceModel{
		Latencies: latencies,
	}
	objectValue, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return objectValue
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

func (o *CatalogEntitySLOLightstepStreamTargetLatencyResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntitySLOLightstepTargetLatency) types.Object {
	ob := CatalogEntitySLOLightstepStreamTargetLatencyResourceModel{
		Percentile: types.Float64Value(entity.Percentile),
		Target:     types.Int64Value(entity.Target),
		SLO:        types.Float64Value(entity.SLO),
	}
	objectType, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return objectType
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
	CodeCov   types.Object `tfsdk:"code_cov"`
	Mend      types.Object `tfsdk:"mend"`
	SonarQube types.Object `tfsdk:"sonar_qube"`
	Veracode  types.Object `tfsdk:"veracode"`
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

func (o *CatalogEntityStaticAnalysisResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityStaticAnalysis {
	doo := getDefaultObjectOptions()

	cc := CatalogEntityStaticAnalysisCodeCovResourceModel{}
	o.CodeCov.As(ctx, &cc, doo)

	me := CatalogEntityStaticAnalysisMendResourceModel{}
	o.Mend.As(ctx, &me, doo)

	sq := CatalogEntityStaticAnalysisSonarQubeResourceModel{}
	o.SonarQube.As(ctx, &sq, doo)

	vc := CatalogEntityStaticAnalysisVeracodeResourceModel{}
	o.Veracode.As(ctx, &vc, doo)

	return cortex.CatalogEntityStaticAnalysis{
		CodeCov:   cc.ToApiModel(),
		Mend:      me.ToApiModel(),
		SonarQube: sq.ToApiModel(),
		Veracode:  vc.ToApiModel(),
	}
}

func (o *CatalogEntityStaticAnalysisResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityStaticAnalysis) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	cc := CatalogEntityStaticAnalysisCodeCovResourceModel{}
	me := CatalogEntityStaticAnalysisMendResourceModel{}
	sq := CatalogEntityStaticAnalysisSonarQubeResourceModel{}
	vc := CatalogEntityStaticAnalysisVeracodeResourceModel{}

	ob := CatalogEntityStaticAnalysisResourceModel{
		CodeCov:   cc.FromApiModel(ctx, diagnostics, &entity.CodeCov),
		Mend:      me.FromApiModel(ctx, diagnostics, &entity.Mend),
		SonarQube: sq.FromApiModel(ctx, diagnostics, &entity.SonarQube),
		Veracode:  vc.FromApiModel(ctx, diagnostics, &entity.Veracode),
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), ob)
	diagnostics.Append(d...)
	return obj
}

// CodeCov

type CatalogEntityStaticAnalysisCodeCovResourceModel struct {
	Repository types.String `tfsdk:"repository"`
	Provider   types.String `tfsdk:"provider"`
	Owner      types.String `tfsdk:"owner"`
	Flag       types.String `tfsdk:"flag"`
}

func (o *CatalogEntityStaticAnalysisCodeCovResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"provider":   types.StringType,
		"repository": types.StringType,
		"owner":      types.StringType,
		"flag":       types.StringType,
	}
}

func (o *CatalogEntityStaticAnalysisCodeCovResourceModel) ToApiModel() cortex.CatalogEntityStaticAnalysisCodeCov {
	return cortex.CatalogEntityStaticAnalysisCodeCov{
		Repository: o.Repository.ValueString(),
		Provider:   o.Provider.ValueString(),
		Owner:      o.Owner.ValueString(),
		Flag:       o.Flag.ValueString(),
	}
}

func (o *CatalogEntityStaticAnalysisCodeCovResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, cov *cortex.CatalogEntityStaticAnalysisCodeCov) types.Object {
	if !cov.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityStaticAnalysisCodeCovResourceModel{
		Repository: types.StringValue(cov.Repository),
		Provider:   types.StringValue(cov.Provider),
	}
	if cov.Owner != "" {
		ob.Owner = types.StringValue(cov.Owner)
	} else {
		ob.Owner = types.StringNull()
	}
	if cov.Flag != "" {
		ob.Flag = types.StringValue(cov.Flag)
	} else {
		ob.Flag = types.StringNull()
	}
	objectValue, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return objectValue
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

func (o *CatalogEntityStaticAnalysisMendResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityStaticAnalysisMend) types.Object {
	applicationIds := make([]types.String, len(entity.ApplicationIDs))
	for i, e := range entity.ApplicationIDs {
		applicationIds[i] = types.StringValue(e)
	}
	projectIds := make([]types.String, len(entity.ProjectIDs))
	for i, e := range entity.ProjectIDs {
		projectIds[i] = types.StringValue(e)
	}

	ob := CatalogEntityStaticAnalysisMendResourceModel{
		ApplicationIDs: applicationIds,
		ProjectIDs:     projectIds,
	}
	objectValue, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return objectValue

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

func (o *CatalogEntityStaticAnalysisSonarQubeResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityStaticAnalysisSonarQube) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	ob := CatalogEntityStaticAnalysisSonarQubeResourceModel{
		Project: types.StringValue(entity.Project),
	}
	if entity.Alias != "" {
		ob.Alias = types.StringValue(entity.Alias)
	} else {
		ob.Alias = types.StringNull()
	}
	objectValue, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return objectValue
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

func (o *CatalogEntityStaticAnalysisVeracodeResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityStaticAnalysisVeracode) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	var sandboxes = make([]CatalogEntityStaticAnalysisVeracodeSandboxResourceModel, len(entity.Sandboxes))
	for i, e := range entity.Sandboxes {
		ob := CatalogEntityStaticAnalysisVeracodeSandboxResourceModel{}
		sandboxes[i] = ob.FromApiModel(&e)
	}
	applicationNames := make([]types.String, len(entity.ApplicationNames))
	for i, e := range entity.ApplicationNames {
		applicationNames[i] = types.StringValue(e)
	}
	ob := CatalogEntityStaticAnalysisVeracodeResourceModel{
		ApplicationNames: applicationNames,
		Sandboxes:        sandboxes,
	}
	objectValue, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return objectValue
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

/***********************************************************************************************************************
 * CI/CD
 **********************************************************************************************************************/

type CatalogEntityCiCdResourceModel struct {
	Buildkite types.Object `tfsdk:"buildkite"`
}

func (o *CatalogEntityCiCdResourceModel) AttrTypes() map[string]attr.Type {
	bk := CatalogEntityCiCdBuildkiteResourceModel{}
	return map[string]attr.Type{
		"buildkite": types.ObjectType{AttrTypes: bk.AttrTypes()},
	}
}

func (o *CatalogEntityCiCdResourceModel) ToApiModel(ctx context.Context) cortex.CatalogEntityCiCd {
	doo := getDefaultObjectOptions()

	bk := CatalogEntityCiCdBuildkiteResourceModel{}
	o.Buildkite.As(ctx, &bk, doo)

	return cortex.CatalogEntityCiCd{
		Buildkite: bk.ToApiModel(),
	}
}

func (o *CatalogEntityCiCdResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntityCiCd) types.Object {
	if !entity.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	bk := CatalogEntityCiCdBuildkiteResourceModel{}
	ob := CatalogEntityCiCdResourceModel{
		Buildkite: bk.FromApiModel(ctx, diagnostics, &entity.Buildkite),
	}
	obj, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), ob)
	diagnostics.Append(d...)
	return obj
}

// Buildkite

type CatalogEntityCiCdBuildkiteResourceModel struct {
	Pipelines []CatalogEntityCiCdBuildkitePipelineResourceModel `tfsdk:"pipelines"`
	Tags      []CatalogEntityCiCdBuildkiteTagResourceModel      `tfsdk:"tags"`
}

func (o *CatalogEntityCiCdBuildkiteResourceModel) AttrTypes() map[string]attr.Type {
	pp := CatalogEntityCiCdBuildkitePipelineResourceModel{}
	tg := CatalogEntityCiCdBuildkiteTagResourceModel{}
	return map[string]attr.Type{
		"pipelines": types.SetType{ElemType: types.ObjectType{AttrTypes: pp.AttrTypes()}},
		"tags":      types.SetType{ElemType: types.ObjectType{AttrTypes: tg.AttrTypes()}},
	}
}

func (o *CatalogEntityCiCdBuildkiteResourceModel) ToApiModel() cortex.CatalogEntityCiCdBuildkite {
	var pipelines = make([]cortex.CatalogEntityCiCdBuildkitePipeline, len(o.Pipelines))
	for i, e := range o.Pipelines {
		pipelines[i] = e.ToApiModel()
	}
	var tags = make([]cortex.CatalogEntityCiCdBuildkiteTag, len(o.Tags))
	for i, e := range o.Tags {
		tags[i] = e.ToApiModel()
	}
	return cortex.CatalogEntityCiCdBuildkite{
		Pipelines: pipelines,
		Tags:      tags,
	}
}

func (o *CatalogEntityCiCdBuildkiteResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, bk *cortex.CatalogEntityCiCdBuildkite) types.Object {
	if !bk.Enabled() {
		return types.ObjectNull(o.AttrTypes())
	}

	var pipelines = make([]CatalogEntityCiCdBuildkitePipelineResourceModel, len(bk.Pipelines))
	for i, e := range bk.Pipelines {
		pp := CatalogEntityCiCdBuildkitePipelineResourceModel{}
		pipelines[i] = pp.FromApiModel(&e)
	}
	var tags = make([]CatalogEntityCiCdBuildkiteTagResourceModel, len(bk.Tags))
	for i, e := range bk.Tags {
		tg := CatalogEntityCiCdBuildkiteTagResourceModel{}
		tags[i] = tg.FromApiModel(&e)
	}

	ob := CatalogEntityCiCdBuildkiteResourceModel{
		Pipelines: pipelines,
		Tags:      tags,
	}
	objectValue, d := types.ObjectValueFrom(ctx, ob.AttrTypes(), &ob)
	diagnostics.Append(d...)
	return objectValue
}

// Buildkite Pipeline

type CatalogEntityCiCdBuildkitePipelineResourceModel struct {
	Slug types.String `tfsdk:"slug"`
}

func (o *CatalogEntityCiCdBuildkitePipelineResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"slug": types.StringType,
	}
}

func (o *CatalogEntityCiCdBuildkitePipelineResourceModel) ToApiModel() cortex.CatalogEntityCiCdBuildkitePipeline {
	return cortex.CatalogEntityCiCdBuildkitePipeline{
		Slug: o.Slug.ValueString(),
	}
}

func (o *CatalogEntityCiCdBuildkitePipelineResourceModel) FromApiModel(pipeline *cortex.CatalogEntityCiCdBuildkitePipeline) CatalogEntityCiCdBuildkitePipelineResourceModel {
	return CatalogEntityCiCdBuildkitePipelineResourceModel{
		Slug: types.StringValue(pipeline.Slug),
	}
}

// Buildkite Tag

type CatalogEntityCiCdBuildkiteTagResourceModel struct {
	Tag types.String `tfsdk:"tag"`
}

func (o *CatalogEntityCiCdBuildkiteTagResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag": types.StringType,
	}
}

func (o *CatalogEntityCiCdBuildkiteTagResourceModel) ToApiModel() cortex.CatalogEntityCiCdBuildkiteTag {
	return cortex.CatalogEntityCiCdBuildkiteTag{
		Tag: o.Tag.ValueString(),
	}
}

func (o *CatalogEntityCiCdBuildkiteTagResourceModel) FromApiModel(pipeline *cortex.CatalogEntityCiCdBuildkiteTag) CatalogEntityCiCdBuildkiteTagResourceModel {
	return CatalogEntityCiCdBuildkiteTagResourceModel{
		Tag: types.StringValue(pipeline.Tag),
	}
}
