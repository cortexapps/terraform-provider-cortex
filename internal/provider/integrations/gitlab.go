package integrations

import (
	"context"
	"regexp"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	api "github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// gitlabDefinition maps credentials.token to personalAccessToken. The API ignores host on update, so a host change
// replaces the configuration.
type gitlabDefinition struct{}

type gitlabSettingsModel struct {
	Host                 types.String `tfsdk:"host"`
	GroupNames           types.List   `tfsdk:"group_names"`
	HidePersonalProjects types.Bool   `tfsdk:"hide_personal_projects"`
}

var _ multiInstanceDefinition = gitlabDefinition{}

func (gitlabDefinition) Name() string  { return "gitlab" }
func (gitlabDefinition) Title() string { return "GitLab" }
func (gitlabDefinition) CredentialKinds() []credentialKind {
	return []credentialKind{credentialToken}
}

func (gitlabDefinition) SettingsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "GitLab settings. Use `credentials.token`: `value` is the GitLab personal access token. " +
			"The Cortex API cannot change `host` in place, so a change replaces the configuration.",
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "URL of a self-managed GitLab instance. Not set means gitlab.com.",
				Optional:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"group_names": schema.ListAttribute{
				MarkdownDescription: "GitLab groups to include. Defaults to an empty list. Names must not be blank.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				// The API drops blank names, so they would make the apply inconsistent.
				Validators: []validator.List{listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(`\S`), "must not be blank"),
				)},
			},
			"hide_personal_projects": schema.BoolAttribute{
				MarkdownDescription: "Whether to hide personal projects. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (gitlabDefinition) SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics) {
	if settingsUnknown(plan, state) {
		return true, nil
	}
	p, s, diags := asSettings[gitlabSettingsModel](ctx, plan, state)
	if p == nil {
		return false, diags
	}
	return !p.Host.Equal(s.Host), diags
}

func (d gitlabDefinition) List(ctx context.Context, c *cortex.HttpClient, prior types.Object) ([]configurationState, error) {
	configurations, err := api.Gitlab(c).List(ctx)
	if err != nil {
		return nil, err
	}
	states := make([]configurationState, 0, len(configurations))
	for _, cfg := range configurations {
		st, err := d.state(ctx, cfg)
		if err != nil {
			return nil, err
		}
		states = append(states, st)
	}
	return states, nil
}

func (d gitlabDefinition) Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error) {
	s, groups, err := d.settings(ctx, in)
	if err != nil {
		return configurationState{}, err
	}
	cfg, err := api.Gitlab(c).Create(ctx, in.Alias, api.CreateGitlabConfigurationRequest{
		Alias:                in.Alias,
		IsDefault:            in.IsDefault,
		Host:                 s.Host.ValueString(),
		PersonalAccessToken:  in.Credentials.Parts["value"],
		HidePersonalProjects: s.HidePersonalProjects.ValueBool(),
		GroupNames:           groups,
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (d gitlabDefinition) Update(ctx context.Context, c *cortex.HttpClient, currentAlias string, in configurationInput) (configurationState, error) {
	s, groups, err := d.settings(ctx, in)
	if err != nil {
		return configurationState{}, err
	}
	cfg, err := api.Gitlab(c).Update(ctx, currentAlias, in.Alias, api.UpdateGitlabConfigurationRequest{
		Alias:                in.Alias,
		IsDefault:            in.IsDefault,
		HidePersonalProjects: s.HidePersonalProjects.ValueBool(),
		GroupNames:           groups,
		PersonalAccessToken:  in.Credentials.Parts["value"],
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (gitlabDefinition) Delete(ctx context.Context, c *cortex.HttpClient, alias string) error {
	return api.Gitlab(c).Delete(ctx, alias)
}

func (d gitlabDefinition) settings(ctx context.Context, in configurationInput) (gitlabSettingsModel, []string, error) {
	s, err := settingsFrom[gitlabSettingsModel](ctx, in.Settings)
	if err != nil {
		return s, nil, settingsError(d, err)
	}
	groups, diags := stringList(ctx, s.GroupNames)
	return s, groups, diagsError(diags)
}

func (d gitlabDefinition) state(ctx context.Context, cfg api.GitlabConfiguration) (configurationState, error) {
	groups := cfg.GroupNames
	if groups == nil {
		groups = []string{}
	}
	list, diags := types.ListValueFrom(ctx, types.StringType, groups)
	if err := diagsError(diags); err != nil {
		return configurationState{}, err
	}
	obj, err := settingsObject(ctx, d, gitlabSettingsModel{
		Host:                 optionalString(cfg.Host),
		GroupNames:           list,
		HidePersonalProjects: types.BoolValue(cfg.HidePersonalProjects),
	})
	return configurationState{
		Alias:     cfg.Alias,
		IsDefault: cfg.IsDefault,
		Settings:  obj,
		LastFour:  map[string]string{"value": cfg.LastFour},
	}, err
}
