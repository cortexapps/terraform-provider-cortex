package integrations

import "github.com/cortexapps/terraform-provider-cortex/internal/cortex"

type CreateHarnessConfigurationRequest struct {
	Alias     string `json:"alias"`
	IsDefault bool   `json:"isDefault"`
	ApiKey    string `json:"apiKey"`
	AccountId string `json:"accountId"`
	Host      string `json:"host,omitempty"`
}

// UpdateHarnessConfigurationRequest omits an unset key or host: the API keeps the current value of a field an update omits.
type UpdateHarnessConfigurationRequest struct {
	Alias     string `json:"alias"`
	IsDefault bool   `json:"isDefault"`
	ApiKey    string `json:"apiKey,omitempty"`
	AccountId string `json:"accountId,omitempty"`
	Host      string `json:"host,omitempty"`
}

// HarnessConfiguration is a Harness configuration. The API returns only the last four characters of the key.
type HarnessConfiguration struct {
	Alias          string `json:"alias"`
	IsDefault      bool   `json:"isDefault"`
	AccountId      string `json:"accountId"`
	Host           string `json:"host,omitempty"`
	LastFourApiKey string `json:"lastFourApiKey"`
}

func (c HarnessConfiguration) GetAlias() string { return c.Alias }

// Harness returns the client for the Harness configurations.
func Harness(c *cortex.HttpClient) MultiInstanceClientInterface[CreateHarnessConfigurationRequest, UpdateHarnessConfigurationRequest, HarnessConfiguration] {
	return &MultiInstanceClient[CreateHarnessConfigurationRequest, UpdateHarnessConfigurationRequest, HarnessConfiguration]{client: c, base: "/api/v1/harness/", name: "harness"}
}
