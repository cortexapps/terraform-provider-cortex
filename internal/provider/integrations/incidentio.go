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

// incidentIoDefinition maps credentials.token to apiKey. incident.io has no settings.
type incidentIoDefinition struct{}

var _ multiInstanceDefinition = incidentIoDefinition{}

func (incidentIoDefinition) Name() string  { return "incident_io" }
func (incidentIoDefinition) Title() string { return "incident.io" }
func (incidentIoDefinition) CredentialKinds() []credentialKind {
	return []credentialKind{credentialToken}
}
func (incidentIoDefinition) CredentialsUpdatable() bool { return true }

func (incidentIoDefinition) SettingsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "incident.io settings. incident.io has no settings, so set `incident_io = {}`. Use " +
			"`credentials.token`: `value` is the incident.io API key.",
		Optional:   true,
		Attributes: map[string]schema.Attribute{},
	}
}

func (incidentIoDefinition) SettingsRequireReplace(ctx context.Context, plan types.Object, state types.Object) (bool, diag.Diagnostics) {
	return false, nil
}

func (d incidentIoDefinition) List(ctx context.Context, c *cortex.HttpClient, prior types.Object) ([]configurationState, error) {
	configurations, err := api.IncidentIo(c).List(ctx)
	if err != nil {
		return nil, err
	}
	states := make([]configurationState, 0, len(configurations))
	for _, cfg := range configurations {
		states = append(states, d.state(cfg))
	}
	return states, nil
}

func (d incidentIoDefinition) Create(ctx context.Context, c *cortex.HttpClient, in configurationInput) (configurationState, error) {
	cfg, err := api.IncidentIo(c).Create(ctx, in.Alias, api.CreateIncidentIoConfigurationRequest{
		Alias: in.Alias, IsDefault: in.IsDefault, ApiKey: in.Credentials.Parts["value"],
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(cfg), nil
}

func (d incidentIoDefinition) Update(ctx context.Context, c *cortex.HttpClient, currentAlias string, in configurationInput) (configurationState, error) {
	cfg, err := api.IncidentIo(c).Update(ctx, currentAlias, in.Alias, api.UpdateIncidentIoConfigurationRequest{
		Alias: in.Alias, IsDefault: in.IsDefault, ApiKey: in.Credentials.Parts["value"],
	})
	if err != nil {
		return configurationState{}, err
	}
	return d.state(cfg), nil
}

func (incidentIoDefinition) Delete(ctx context.Context, c *cortex.HttpClient, alias string) error {
	return api.IncidentIo(c).Delete(ctx, alias)
}

func (incidentIoDefinition) state(cfg api.IncidentIoConfiguration) configurationState {
	return configurationState{
		Alias:     cfg.Alias,
		IsDefault: cfg.IsDefault,
		Settings:  types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}),
		LastFour:  map[string]string{"value": cfg.LastFour},
	}
}
