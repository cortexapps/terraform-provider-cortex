package cortex

// DatadogConfiguration is one Datadog integration configuration. The API never returns the full API key or
// application key, only their last four characters.
type DatadogConfiguration struct {
	Alias           string   `json:"alias"`
	IsDefault       bool     `json:"isDefault"`
	Environments    []string `json:"environments"`
	Region          string   `json:"region"`
	CustomSubdomain string   `json:"customSubdomain,omitempty"`
	LastFourApiKey  string   `json:"lastFourApiKey"`
	LastFourAppKey  string   `json:"lastFourAppKey"`
}

func (c DatadogConfiguration) GetAlias() string   { return c.Alias }
func (c DatadogConfiguration) GetIsDefault() bool { return c.IsDefault }

// DatadogConfigurationsResponse is the list response. Tests use it to build fake responses.
type DatadogConfigurationsResponse struct {
	Configurations []DatadogConfiguration `json:"configurations"`
}

type CreateDatadogConfigurationRequest struct {
	Alias           string   `json:"alias"`
	IsDefault       bool     `json:"isDefault"`
	ApiKey          string   `json:"apiKey"`
	AppKey          string   `json:"appKey"`
	Region          string   `json:"region"`
	Environments    []string `json:"environments"`
	CustomSubdomain string   `json:"customSubdomain,omitempty"`
}

// UpdateDatadogConfigurationRequest holds the only fields the public API can change. Alias is the new alias; the
// current alias goes in the path. The API keys, region, and custom subdomain cannot change after creation.
type UpdateDatadogConfigurationRequest struct {
	Alias        string   `json:"alias"`
	IsDefault    bool     `json:"isDefault"`
	Environments []string `json:"environments"`
}

func (c *HttpClient) DatadogConfigurations() *MultiInstanceClient[CreateDatadogConfigurationRequest, UpdateDatadogConfigurationRequest, DatadogConfiguration] {
	return &MultiInstanceClient[CreateDatadogConfigurationRequest, UpdateDatadogConfigurationRequest, DatadogConfiguration]{client: c, domain: "datadog", name: "datadog"}
}
