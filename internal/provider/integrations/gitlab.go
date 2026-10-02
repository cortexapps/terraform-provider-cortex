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

// gitlabDefinition maps credentials.token to personalAccessToken. The API changes the host in place, but it keeps the
// current host when an update omits it, so only removing a set host replaces the configuration.
type gitlabDefinition struct{}

type gitlabSettingsModel struct {
	Host                 types.String `tfsdk:"host"`
	GroupNames           types.List   `tfsdk:"group_names"`
	HidePersonalProjects types.Bool   `tfsdk:"hide_personal_projects"`
}

var _ multiInstanceDefinition = gitlabDefinition{}

var nonBlank = regexp.MustCompile(`\S`)

func (gitlabDefinition) Name() string  { return "gitlab" }
func (gitlabDefinition) Title() string { return "GitLab" }
func (gitlabDefinition) CredentialKinds() []credentialKind {
	return []credentialKind{credentialToken}
}

func (gitlabDefinition) SettingsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "GitLab settings. Use `credentials.token`: `value` is the GitLab personal access token. " +
			"A change of `host` updates the configuration in place. The Cortex API cannot remove a host, so removing " +
			"`host` replaces the configuration. A `host` that is unknown until apply can resolve to null, so it also " +
			"replaces a configuration that has a host.",
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "URL of a self-managed GitLab instance. Not set means gitlab.com. Must not be blank.",
				Optional:            true,
				// The API keeps the current host on a blank update, so a blank host would make the apply inconsistent.
				Validators: []validator.String{stringvalidator.RegexMatches(nonBlank, "must not be blank")},
			},
			"group_names": schema.ListAttribute{
				MarkdownDescription: "GitLab groups to include. Defaults to an empty list. Names must not be blank.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				// The API drops blank names, so they would make the apply inconsistent.
				Validators: []validator.List{listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(nonBlank, "must not be blank"),
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

// SettingsRequireReplace replaces only to remove a set host: the API keeps the current host when an update omits it.
// A host in an unknown settings block is unknown too.
func (gitlabDefinition) SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics) {
	var planHost, stateHost types.String
	var diags diag.Diagnostics
	if settingsUnknown(plan, state) {
		s, err := settingsFrom[gitlabSettingsModel](ctx, state)
		if err != nil {
			return false, diag.Diagnostics{diag.NewErrorDiagnostic("Invalid GitLab settings", err.Error())}
		}
		planHost, stateHost = types.StringUnknown(), s.Host
	} else {
		var p, s *gitlabSettingsModel
		p, s, diags = asSettings[gitlabSettingsModel](ctx, plan, state)
		if p == nil {
			return false, diags
		}
		planHost, stateHost = p.Host, s.Host
	}
	return !stateHost.IsNull() && (planHost.IsNull() || planHost.IsUnknown()), diags
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
		Host:                 s.Host.ValueString(),
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
