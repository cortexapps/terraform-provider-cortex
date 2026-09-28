package integrations

import "github.com/cortexapps/terraform-provider-cortex/internal/cortex"

type CreateIncidentIoConfigurationRequest struct {
	Alias     string `json:"alias"`
	IsDefault bool   `json:"isDefault"`
	ApiKey    string `json:"apiKey"`
}

type UpdateIncidentIoConfigurationRequest struct {
	Alias     string `json:"alias"`
	IsDefault bool   `json:"isDefault"`
	ApiKey    string `json:"apiKey,omitempty"`
}

// IncidentIoConfiguration is an incident.io configuration. The API returns only the last four characters of the key.
type IncidentIoConfiguration struct {
	Alias     string `json:"alias"`
	IsDefault bool   `json:"isDefault"`
	LastFour  string `json:"lastFour"`
}

func (c IncidentIoConfiguration) GetAlias() string { return c.Alias }

// IncidentIo returns the client for the incident.io configurations.
func IncidentIo(c *cortex.HttpClient) MultiInstanceClientInterface[CreateIncidentIoConfigurationRequest, UpdateIncidentIoConfigurationRequest, IncidentIoConfiguration] {
	return &MultiInstanceClient[CreateIncidentIoConfigurationRequest, UpdateIncidentIoConfigurationRequest, IncidentIoConfiguration]{client: c, base: "/api/v1/incidentio/", name: "incident.io"}
}
