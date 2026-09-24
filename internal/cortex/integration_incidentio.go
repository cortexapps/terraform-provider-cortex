package cortex

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

func (c IncidentIoConfiguration) GetAlias() string   { return c.Alias }
func (c IncidentIoConfiguration) GetIsDefault() bool { return c.IsDefault }

func (c *HttpClient) IncidentIoConfigurations() *MultiInstanceClient[CreateIncidentIoConfigurationRequest, UpdateIncidentIoConfigurationRequest, IncidentIoConfiguration] {
	return &MultiInstanceClient[CreateIncidentIoConfigurationRequest, UpdateIncidentIoConfigurationRequest, IncidentIoConfiguration]{client: c, domain: "incident_io", name: "incident.io"}
}
