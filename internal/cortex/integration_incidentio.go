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

func (c IncidentIoConfiguration) GetAlias() string { return c.Alias }
