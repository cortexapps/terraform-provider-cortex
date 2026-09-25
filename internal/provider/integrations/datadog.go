package integrations

import (
	"context"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	api "github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// datadogDefinition maps credentials.key_pair to apiKey (key) and appKey (secret). The public update endpoint can
// change only environments, alias, and the default.
type datadogDefinition struct{}

type datadogSettingsModel struct {
	Region          types.String `tfsdk:"region"`
	Environments    types.List   `tfsdk:"environments"`
	CustomSubdomain types.String `tfsdk:"custom_subdomain"`
}

var _ multiInstanceDefinition = datadogDefinition{}

func (datadogDefinition) Name() string  { return "datadog" }
func (datadogDefinition) Title() string { return "Datadog" }
func (datadogDefinition) CredentialKinds() []credentialKind {
	return []credentialKind{credentialKeyPair}
}
func (datadogDefinition) CredentialsUpdatable() bool { return true }

func (datadogDefinition) SettingsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Datadog settings. Use `credentials.key_pair`: `key` is the Datadog API key and `secret` is " +
			"the application key. A change of the keys, `region`, or `custom_subdomain` updates the configuration in place. " +
			"The Cortex API cannot remove a custom subdomain, so removing `custom_subdomain` replaces the configuration.",
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				MarkdownDescription: "Datadog region (site). One of `US1`, `US3`, `US5`, `EU1`, `US1_FED`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("US1", "US3", "US5", "EU1", "US1_FED")},
			},
			"environments": schema.ListAttribute{
				MarkdownDescription: "Datadog environments to use. Defaults to an empty list.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"custom_subdomain": schema.StringAttribute{
				MarkdownDescription: "Custom subdomain of the Datadog organization, if it uses one.",
				Optional:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		},
	}
}

// SettingsRequireReplace replaces only to remove a set custom subdomain: the API keeps the current custom subdomain
// when an update omits it.
func (datadogDefinition) SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics) {
	if settingsUnknown(plan, state) {
		s, err := settingsFrom[datadogSettingsModel](ctx, state)
		if err != nil {
			return false, diag.Diagnostics{diag.NewErrorDiagnostic("Invalid Datadog settings", err.Error())}
		}
		return !s.CustomSubdomain.IsNull(), nil
	}
	p, s, diags := asSettings[datadogSettingsModel](ctx, plan, state)
	if p == nil {
		return false, diags
	}
	return !s.CustomSubdomain.IsNull() && (p.CustomSubdomain.IsNull() || p.CustomSubdomain.IsUnknown()), diags
}

func (d datadogDefinition) List(ctx context.Context, c *cortex.HttpClient, prior types.Object) ([]configurationState, error) {
	configurations, err := api.Datadog(c).List(ctx)
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

func (d datadogDefinition) Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error) {
	s, environments, err := d.settings(ctx, in)
	if err != nil {
		return configurationState{}, err
	}
	cfg, err := api.Datadog(c).Create(ctx, in.Alias, api.CreateDatadogConfigurationRequest{
		Alias:           in.Alias,
		IsDefault:       in.IsDefault,
		ApiKey:          in.Credentials.Parts["key"],
		AppKey:          in.Credentials.Parts["secret"],
		Region:          s.Region.ValueString(),
		Environments:    environments,
		CustomSubdomain: s.CustomSubdomain.ValueString(),
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (d datadogDefinition) Update(ctx context.Context, c *cortex.HttpClient, currentAlias string, in configurationInput) (configurationState, error) {
	s, environments, err := d.settings(ctx, in)
	if err != nil {
		return configurationState{}, err
	}
	cfg, err := api.Datadog(c).Update(ctx, currentAlias, in.Alias, api.UpdateDatadogConfigurationRequest{
		Alias:           in.Alias,
		IsDefault:       in.IsDefault,
		Environments:    environments,
		ApiKey:          in.Credentials.Parts["key"],
		AppKey:          in.Credentials.Parts["secret"],
		Region:          s.Region.ValueString(),
		CustomSubdomain: s.CustomSubdomain.ValueString(),
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (datadogDefinition) Delete(ctx context.Context, c *cortex.HttpClient, alias string) error {
	return api.Datadog(c).Delete(ctx, alias)
}

func (d datadogDefinition) settings(ctx context.Context, in configurationInput) (datadogSettingsModel, []string, error) {
	s, err := settingsFrom[datadogSettingsModel](ctx, in.Settings)
	if err != nil {
		return s, nil, settingsError(d, err)
	}
	environments, diags := stringList(ctx, s.Environments)
	return s, environments, diagsError(diags)
}

func (d datadogDefinition) state(ctx context.Context, cfg api.DatadogConfiguration) (configurationState, error) {
	environments := cfg.Environments
	if environments == nil {
		environments = []string{}
	}
	list, diags := types.ListValueFrom(ctx, types.StringType, environments)
	if err := diagsError(diags); err != nil {
		return configurationState{}, err
	}
	obj, err := settingsObject(ctx, d, datadogSettingsModel{
		Region:          types.StringValue(cfg.Region),
		Environments:    list,
		CustomSubdomain: optionalString(cfg.CustomSubdomain),
	})
	return configurationState{
		Alias:     cfg.Alias,
		IsDefault: cfg.IsDefault,
		Settings:  obj,
		LastFour:  map[string]string{"key": cfg.LastFourApiKey, "secret": cfg.LastFourAppKey},
	}, err
}
