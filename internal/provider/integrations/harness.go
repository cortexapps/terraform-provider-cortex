package integrations

import (
	"context"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	api "github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// harnessDefinition maps credentials.token to apiKey. The API changes the account and the host in place, but it keeps
// the current host when an update omits it, so only removing a set host replaces the configuration.
type harnessDefinition struct{}

type harnessSettingsModel struct {
	AccountId types.String `tfsdk:"account_id"`
	Host      types.String `tfsdk:"host"`
}

var _ multiInstanceDefinition = harnessDefinition{}

func (harnessDefinition) Name() string  { return "harness" }
func (harnessDefinition) Title() string { return "Harness" }
func (harnessDefinition) CredentialKinds() []credentialKind {
	return []credentialKind{credentialToken}
}

func (harnessDefinition) SettingsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Harness settings. Use `credentials.token`: `value` is the Harness API key. A change of " +
			"`account_id` or `host` updates the configuration in place. The Cortex API cannot remove a host, so " +
			"removing `host` replaces the configuration. A `host` that is unknown until apply can resolve to null, so " +
			"it also replaces a configuration that has a host.",
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				MarkdownDescription: "Harness account identifier. Must not be blank.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.RegexMatches(nonBlank, "must not be blank")},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "URL of the Harness instance, for a self-managed or dedicated Harness. Not set means " +
					"`https://app.harness.io`. Must not be blank.",
				Optional: true,
				// The API keeps the current host on a blank update, so a blank host would make the apply inconsistent.
				Validators: []validator.String{stringvalidator.RegexMatches(nonBlank, "must not be blank")},
			},
		},
	}
}

// SettingsRequireReplace replaces only to remove a set host: the API keeps the current host when an update omits it.
// A host in an unknown settings block is unknown too.
func (harnessDefinition) SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics) {
	var planHost, stateHost types.String
	var diags diag.Diagnostics
	if settingsUnknown(plan, state) {
		s, err := settingsFrom[harnessSettingsModel](ctx, state)
		if err != nil {
			return false, diag.Diagnostics{diag.NewErrorDiagnostic("Invalid Harness settings", err.Error())}
		}
		planHost, stateHost = types.StringUnknown(), s.Host
	} else {
		var p, s *harnessSettingsModel
		p, s, diags = asSettings[harnessSettingsModel](ctx, plan, state)
		if p == nil {
			return false, diags
		}
		planHost, stateHost = p.Host, s.Host
	}
	return !stateHost.IsNull() && (planHost.IsNull() || planHost.IsUnknown()), diags
}

func (d harnessDefinition) List(ctx context.Context, c *cortex.HttpClient, prior types.Object) ([]configurationState, error) {
	configurations, err := api.Harness(c).List(ctx)
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

func (d harnessDefinition) Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error) {
	s, err := settingsFrom[harnessSettingsModel](ctx, in.Settings)
	if err != nil {
		return configurationState{}, settingsError(d, err)
	}
	cfg, err := api.Harness(c).Create(ctx, in.Alias, api.CreateHarnessConfigurationRequest{
		Alias:     in.Alias,
		IsDefault: in.IsDefault,
		ApiKey:    in.Credentials.Parts["value"],
		AccountId: s.AccountId.ValueString(),
		Host:      s.Host.ValueString(),
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (d harnessDefinition) Update(ctx context.Context, c *cortex.HttpClient, currentAlias string, in configurationInput) (configurationState, error) {
	s, err := settingsFrom[harnessSettingsModel](ctx, in.Settings)
	if err != nil {
		return configurationState{}, settingsError(d, err)
	}
	cfg, err := api.Harness(c).Update(ctx, currentAlias, in.Alias, api.UpdateHarnessConfigurationRequest{
		Alias:     in.Alias,
		IsDefault: in.IsDefault,
		ApiKey:    in.Credentials.Parts["value"],
		AccountId: s.AccountId.ValueString(),
		Host:      s.Host.ValueString(),
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(ctx, cfg)
}

func (harnessDefinition) Delete(ctx context.Context, c *cortex.HttpClient, alias string) error {
	return api.Harness(c).Delete(ctx, alias)
}

func (d harnessDefinition) state(ctx context.Context, cfg api.HarnessConfiguration) (configurationState, error) {
	obj, err := settingsObject(ctx, d, harnessSettingsModel{
		AccountId: types.StringValue(cfg.AccountId),
		Host:      optionalString(cfg.Host),
	})
	return configurationState{
		Alias:     cfg.Alias,
		IsDefault: cfg.IsDefault,
		Settings:  obj,
		LastFour:  map[string]string{"value": cfg.LastFourApiKey},
	}, err
}
