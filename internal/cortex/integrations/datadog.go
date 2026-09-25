package integrations

import (
	"encoding/json"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
)

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

func (c DatadogConfiguration) GetAlias() string { return c.Alias }

type CreateDatadogConfigurationRequest struct {
	Alias           string   `json:"alias"`
	IsDefault       bool     `json:"isDefault"`
	ApiKey          string   `json:"apiKey"`
	AppKey          string   `json:"appKey"`
	Region          string   `json:"region"`
	Environments    []string `json:"environments"`
	CustomSubdomain string   `json:"customSubdomain,omitempty"`
}

// UpdateDatadogConfigurationRequest is the body of an update. Alias is the new alias; the current alias goes in the
// path. The API keeps the current value of each optional field that the request omits, so an update cannot remove
// the custom subdomain.
type UpdateDatadogConfigurationRequest struct {
	Alias           string   `json:"alias"`
	IsDefault       bool     `json:"isDefault"`
	Environments    []string `json:"environments"`
	ApiKey          string   `json:"apiKey,omitempty"`
	AppKey          string   `json:"appKey,omitempty"`
	Region          string   `json:"region,omitempty"`
	CustomSubdomain string   `json:"customSubdomain,omitempty"`
}

// MarshalJSON sends an empty list for nil environments, because the API rejects null.
func (r CreateDatadogConfigurationRequest) MarshalJSON() ([]byte, error) {
	type plain CreateDatadogConfigurationRequest
	if r.Environments == nil {
		r.Environments = []string{}
	}
	return json.Marshal(plain(r))
}

// MarshalJSON sends an empty list for nil environments, because the API rejects null.
func (r UpdateDatadogConfigurationRequest) MarshalJSON() ([]byte, error) {
	type plain UpdateDatadogConfigurationRequest
	if r.Environments == nil {
		r.Environments = []string{}
	}
	return json.Marshal(plain(r))
}

// Datadog returns the client for the Datadog configurations.
func Datadog(c *cortex.HttpClient) MultiInstanceClientInterface[CreateDatadogConfigurationRequest, UpdateDatadogConfigurationRequest, DatadogConfiguration] {
	return &MultiInstanceClient[CreateDatadogConfigurationRequest, UpdateDatadogConfigurationRequest, DatadogConfiguration]{client: c, base: "/api/v1/datadog/", name: "datadog"}
}
