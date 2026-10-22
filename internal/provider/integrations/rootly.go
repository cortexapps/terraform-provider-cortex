package integrations

import (
	"context"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	api "github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// rootlyDefinition maps credentials.token to apiKey. Rootly has no settings.
type rootlyDefinition struct{}

var _ multiInstanceDefinition = rootlyDefinition{}

func (rootlyDefinition) Name() string  { return "rootly" }
func (rootlyDefinition) Title() string { return "Rootly" }
func (rootlyDefinition) CredentialKinds() []credentialKind {
	return []credentialKind{credentialToken}
}

func (rootlyDefinition) SettingsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Rootly settings. Rootly has no settings, so set `rootly = {}`. Use " +
			"`credentials.token`: `value` is the Rootly API key.",
		Optional:   true,
		Attributes: map[string]schema.Attribute{},
	}
}

func (rootlyDefinition) SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics) {
	return false, nil
}

func (d rootlyDefinition) List(ctx context.Context, c *cortex.HttpClient, prior types.Object) ([]configurationState, error) {
	configurations, err := api.Rootly(c).List(ctx)
	if err != nil {
		return nil, err
	}
	states := make([]configurationState, 0, len(configurations))
	for _, cfg := range configurations {
		states = append(states, d.state(cfg))
	}
	return states, nil
}

func (d rootlyDefinition) Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error) {
	cfg, err := api.Rootly(c).Create(ctx, in.Alias, api.CreateRootlyConfigurationRequest{
		Alias: in.Alias, IsDefault: in.IsDefault, ApiKey: in.Credentials.Parts["value"],
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(cfg), nil
}

func (d rootlyDefinition) Update(ctx context.Context, c *cortex.HttpClient, currentAlias string, in configurationInput) (configurationState, error) {
	cfg, err := api.Rootly(c).Update(ctx, currentAlias, in.Alias, api.UpdateRootlyConfigurationRequest{
		Alias: in.Alias, IsDefault: in.IsDefault, ApiKey: in.Credentials.Parts["value"],
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(cfg), nil
}

func (rootlyDefinition) Delete(ctx context.Context, c *cortex.HttpClient, alias string) error {
	return api.Rootly(c).Delete(ctx, alias)
}

func (rootlyDefinition) state(cfg api.RootlyConfiguration) configurationState {
	return configurationState{
		Alias:     cfg.Alias,
		IsDefault: cfg.IsDefault,
		Settings:  types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}),
		LastFour:  map[string]string{"value": cfg.LastFour},
	}
}
