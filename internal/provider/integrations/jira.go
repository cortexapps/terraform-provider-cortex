package integrations

import (
	"context"
	"errors"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	api "github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// jiraDefinition maps credentials.basic to email and apiToken (cloud variants) or username and password (on-prem).
// The variant sets the API "type". The API never returns cloudId, so Read keeps it from state. The API ignores host,
// frontendHost, and cloudId on update, so a change to one of them replaces the configuration.
type jiraDefinition struct{}

type jiraSettingsModel struct {
	Cloud       *jiraCloudModel       `tfsdk:"cloud"`
	CloudScoped *jiraCloudScopedModel `tfsdk:"cloud_scoped"`
	OnPrem      *jiraOnPremModel      `tfsdk:"on_prem"`
}

type jiraCloudModel struct {
	Subdomain types.String `tfsdk:"subdomain"`
	BaseUrl   types.String `tfsdk:"base_url"`
}

type jiraCloudScopedModel struct {
	Subdomain types.String `tfsdk:"subdomain"`
	CloudId   types.String `tfsdk:"cloud_id"`
	BaseUrl   types.String `tfsdk:"base_url"`
}

type jiraOnPremModel struct {
	Host         types.String `tfsdk:"host"`
	FrontendHost types.String `tfsdk:"frontend_host"`
}

var jiraBaseUrls = []string{"jira.com", "atlassian.net", "api.atlassian.com/ex/jira"}

var _ multiInstanceDefinition = jiraDefinition{}

func (jiraDefinition) Name() string  { return "jira" }
func (jiraDefinition) Title() string { return "Jira" }
func (jiraDefinition) CredentialKinds() []credentialKind {
	return []credentialKind{credentialBasic}
}
func (jiraDefinition) CredentialsUpdatable() bool { return true }

func (jiraDefinition) SettingsAttribute() schema.SingleNestedAttribute {
	required := func(description string) schema.StringAttribute {
		return schema.StringAttribute{MarkdownDescription: description, Required: true, Validators: []validator.String{stringvalidator.LengthAtLeast(1)}}
	}
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Jira settings. Set exactly one variant. Use `credentials.basic`: for the cloud variants, " +
			"`username` is the email and `password` is the API token; for `on_prem`, they are the username and password. " +
			"The Cortex API cannot change the variant, `host`, `frontend_host`, or `cloud_id` in place, so a change to one " +
			"of them replaces the configuration.",
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"cloud": schema.SingleNestedAttribute{
				MarkdownDescription: "Jira Cloud with basic authentication.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"subdomain": required("Subdomain of the Jira site, for example `acme`."),
					"base_url": schema.StringAttribute{
						MarkdownDescription: "Base URL. One of `jira.com`, `atlassian.net`, `api.atlassian.com/ex/jira`.",
						Required:            true,
						Validators:          []validator.String{stringvalidator.OneOf(jiraBaseUrls...)},
					},
				},
				// The validator counts the attribute it is on, so it sits on one variant and names the other two.
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("cloud_scoped"),
						path.MatchRelative().AtParent().AtName("on_prem"),
					),
				},
			},
			"cloud_scoped": schema.SingleNestedAttribute{
				MarkdownDescription: "Jira Cloud with a scoped API token.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"subdomain": required("Subdomain of the Jira site."),
					"cloud_id":  required("Atlassian cloud ID of the site. The Cortex API does not return it, so Terraform cannot detect a change outside Terraform."),
					"base_url": schema.StringAttribute{
						MarkdownDescription: "Base URL. Defaults to `api.atlassian.com/ex/jira`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("api.atlassian.com/ex/jira"),
						Validators:          []validator.String{stringvalidator.OneOf(jiraBaseUrls...)},
					},
				},
			},
			"on_prem": schema.SingleNestedAttribute{
				MarkdownDescription: "Jira Data Center or Server with basic authentication.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"host": required("URL of the Jira server."),
					"frontend_host": schema.StringAttribute{
						MarkdownDescription: "URL for links in Cortex, when it differs from `host`.",
						Optional:            true,
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
		},
	}
}

func (jiraDefinition) SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics) {
	p, s, diags := asSettings[jiraSettingsModel](ctx, plan, state)
	if p == nil {
		return false, diags
	}
	switch {
	case jiraVariant(*p) != jiraVariant(*s):
		return true, diags
	case p.OnPrem != nil:
		return !p.OnPrem.Host.Equal(s.OnPrem.Host) || !p.OnPrem.FrontendHost.Equal(s.OnPrem.FrontendHost), diags
	case p.CloudScoped != nil:
		return !p.CloudScoped.CloudId.Equal(s.CloudScoped.CloudId), diags
	}
	return false, diags
}

func jiraVariant(s jiraSettingsModel) string {
	switch {
	case s.Cloud != nil:
		return api.JiraTypeCloudBasic
	case s.CloudScoped != nil:
		return api.JiraTypeCloudScoped
	case s.OnPrem != nil:
		return api.JiraTypeOnPremBasic
	}
	return ""
}

// List skips configurations of types the provider does not support, for example OAuth, so that they do not break
// Read of other configurations.
func (d jiraDefinition) List(ctx context.Context, c *cortex.HttpClient, prior types.Object) ([]configurationState, error) {
	configurations, err := api.Jira(c).List(ctx)
	if err != nil {
		return nil, err
	}
	states := make([]configurationState, 0, len(configurations))
	for _, cfg := range configurations {
		if cfg.Type != api.JiraTypeCloudBasic && cfg.Type != api.JiraTypeCloudScoped && cfg.Type != api.JiraTypeOnPremBasic {
			continue
		}
		st, err := d.state(ctx, cfg, prior)
		if err != nil {
			return nil, err
		}
		states = append(states, st)
	}
	return states, nil
}

func (d jiraDefinition) Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error) {
	req, err := d.request(ctx, in)
	if err != nil {
		return configurationState{}, err
	}
	cfg, err := api.Jira(c).Create(ctx, in.Alias, req)
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg, in.Settings)
}

func (d jiraDefinition) Update(ctx context.Context, c *cortex.HttpClient, currentAlias string, in configurationInput) (configurationState, error) {
	req, err := d.request(ctx, in)
	if err != nil {
		return configurationState{}, err
	}
	cfg, err := api.Jira(c).Update(ctx, currentAlias, in.Alias, req)
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg, in.Settings)
}

func (jiraDefinition) Delete(ctx context.Context, c *cortex.HttpClient, alias string) error {
	return api.Jira(c).Delete(ctx, alias)
}

func (d jiraDefinition) request(ctx context.Context, in configurationInput) (api.JiraConfigurationRequest, error) {
	s, err := settingsFrom[jiraSettingsModel](ctx, in.Settings)
	if err != nil {
		return api.JiraConfigurationRequest{}, settingsError(d, err)
	}
	req := api.JiraConfigurationRequest{Alias: in.Alias, IsDefault: in.IsDefault, Type: jiraVariant(s)}
	username, password := in.Credentials.Parts["username"], in.Credentials.Parts["password"]
	switch {
	case s.Cloud != nil:
		req.Subdomain, req.BaseUrl = s.Cloud.Subdomain.ValueString(), s.Cloud.BaseUrl.ValueString()
		req.Email, req.ApiToken = username, password
	case s.CloudScoped != nil:
		req.Subdomain, req.BaseUrl, req.CloudId = s.CloudScoped.Subdomain.ValueString(), s.CloudScoped.BaseUrl.ValueString(), s.CloudScoped.CloudId.ValueString()
		req.Email, req.ApiToken = username, password
	case s.OnPrem != nil:
		req.Host, req.FrontendHost = s.OnPrem.Host.ValueString(), s.OnPrem.FrontendHost.ValueString()
		req.Username, req.Password = username, password
	default:
		return req, settingsError(d, errors.New("set one of cloud, cloud_scoped, on_prem"))
	}
	return req, nil
}

func (d jiraDefinition) state(ctx context.Context, cfg api.JiraConfiguration, prior types.Object) (configurationState, error) {
	s := jiraSettingsModel{}
	readable := map[string]string{}
	switch cfg.Type {
	case api.JiraTypeCloudBasic:
		s.Cloud = &jiraCloudModel{Subdomain: types.StringValue(cfg.Subdomain), BaseUrl: types.StringValue(cfg.BaseUrl)}
		readable["username"] = cfg.Email
	case api.JiraTypeCloudScoped:
		cloudId := types.StringNull()
		if !prior.IsNull() {
			if p, err := settingsFrom[jiraSettingsModel](ctx, prior); err == nil && p.CloudScoped != nil {
				cloudId = p.CloudScoped.CloudId
			}
		}
		s.CloudScoped = &jiraCloudScopedModel{Subdomain: types.StringValue(cfg.Subdomain), CloudId: cloudId, BaseUrl: types.StringValue(cfg.BaseUrl)}
		readable["username"] = cfg.Email
	case api.JiraTypeOnPremBasic:
		s.OnPrem = &jiraOnPremModel{Host: types.StringValue(cfg.Host), FrontendHost: optionalString(cfg.FrontendHost)}
		readable["username"] = cfg.Username
	}
	obj, err := settingsObject(ctx, d, s)
	return configurationState{
		Alias:     cfg.Alias,
		IsDefault: cfg.IsDefault,
		Settings:  obj,
		Readable:  readable,
		LastFour:  map[string]string{"password": cfg.LastFour},
	}, err
}
