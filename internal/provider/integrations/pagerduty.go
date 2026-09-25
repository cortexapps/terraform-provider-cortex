package integrations

import (
	"context"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	api "github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// pagerDutyDefinition maps credentials.token to token. A tenant has at most one PagerDuty configuration, and PUT
// replaces it, so every change updates in place.
type pagerDutyDefinition struct{}

type pagerDutySettingsModel struct {
	IsTokenReadonly types.Bool `tfsdk:"is_token_readonly"`
}

var _ singleInstanceDefinition = pagerDutyDefinition{}

func (pagerDutyDefinition) Name() string  { return "pagerduty" }
func (pagerDutyDefinition) Title() string { return "PagerDuty" }
func (pagerDutyDefinition) CredentialKinds() []credentialKind {
	return []credentialKind{credentialToken}
}

func (pagerDutyDefinition) SettingsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "PagerDuty settings. Use `credentials.token`: `value` is the PagerDuty API token. A tenant " +
			"has one PagerDuty configuration, so `alias` and `is_default` are not allowed. When Cortex already has a " +
			"PagerDuty configuration, import it with ID `pagerduty`.",
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"is_token_readonly": schema.BoolAttribute{
				MarkdownDescription: "Whether the token is read-only. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (pagerDutyDefinition) SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics) {
	return false, nil
}

func (d pagerDutyDefinition) Get(ctx context.Context, c *cortex.HttpClient, prior types.Object) (configurationState, error) {
	cfg, err := api.PagerDuty(c).Get(ctx)
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (d pagerDutyDefinition) Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error) {
	req, err := d.request(ctx, in)
	if err != nil {
		return configurationState{}, err
	}
	cfg, err := api.PagerDuty(c).Create(ctx, req)
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (d pagerDutyDefinition) Replace(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error) {
	req, err := d.request(ctx, in)
	if err != nil {
		return configurationState{}, err
	}
	cfg, err := api.PagerDuty(c).Replace(ctx, req)
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (pagerDutyDefinition) Delete(ctx context.Context, c *cortex.HttpClient) error {
	return api.PagerDuty(c).Delete(ctx)
}

func (d pagerDutyDefinition) request(ctx context.Context, in configurationInput) (api.PagerDutyConfigurationRequest, error) {
	s, err := settingsFrom[pagerDutySettingsModel](ctx, in.Settings)
	if err != nil {
		return api.PagerDutyConfigurationRequest{}, settingsError(d, err)
	}
	return api.PagerDutyConfigurationRequest{Token: in.Credentials.Parts["value"], IsTokenReadonly: s.IsTokenReadonly.ValueBool()}, nil
}

func (d pagerDutyDefinition) state(ctx context.Context, cfg api.PagerDutyConfiguration) (configurationState, error) {
	obj, err := settingsObject(ctx, d, pagerDutySettingsModel{IsTokenReadonly: types.BoolValue(cfg.IsTokenReadonly)})
	return configurationState{Settings: obj, LastFour: map[string]string{"value": cfg.LastFour}}, err
}
