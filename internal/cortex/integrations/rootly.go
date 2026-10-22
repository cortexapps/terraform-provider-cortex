package integrations

import "github.com/cortexapps/terraform-provider-cortex/internal/cortex"

type CreateRootlyConfigurationRequest struct {
	Alias     string `json:"alias"`
	IsDefault bool   `json:"isDefault"`
	ApiKey    string `json:"apiKey"`
}

type UpdateRootlyConfigurationRequest struct {
	Alias     string `json:"alias"`
	IsDefault bool   `json:"isDefault"`
	ApiKey    string `json:"apiKey,omitempty"`
}

// RootlyConfiguration is a Rootly configuration. The API returns only the last four characters of the key.
type RootlyConfiguration struct {
	Alias     string `json:"alias"`
	IsDefault bool   `json:"isDefault"`
	LastFour  string `json:"lastFour"`
}

func (c RootlyConfiguration) GetAlias() string { return c.Alias }

// Rootly returns the client for the Rootly configurations.
func Rootly(c *cortex.HttpClient) MultiInstanceClientInterface[CreateRootlyConfigurationRequest, UpdateRootlyConfigurationRequest, RootlyConfiguration] {
	return &MultiInstanceClient[CreateRootlyConfigurationRequest, UpdateRootlyConfigurationRequest, RootlyConfiguration]{client: c, base: "/api/v1/rootly/", name: "Rootly"}
}
