package cortex

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dghubble/sling"
)

type DatadogConfigurationsClientInterface interface {
	List(ctx context.Context) ([]DatadogConfiguration, error)
	Get(ctx context.Context, alias string) (DatadogConfiguration, error)
	Create(ctx context.Context, req CreateDatadogConfigurationRequest) (DatadogConfiguration, error)
	Update(ctx context.Context, alias string, req UpdateDatadogConfigurationRequest) (DatadogConfiguration, error)
	Delete(ctx context.Context, alias string) error
}

type DatadogConfigurationsClient struct {
	client *HttpClient
}

var _ DatadogConfigurationsClientInterface = &DatadogConfigurationsClient{}

func (c *DatadogConfigurationsClient) Client() *sling.Sling {
	return c.client.Client()
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

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

// DatadogConfigurationsResponse wraps the list, create, and update responses. Create and update also return every
// configuration in the tenant, not only the one that changed.
type DatadogConfigurationsResponse struct {
	Configurations []DatadogConfiguration `json:"configurations"`
}

// findByAlias returns the configuration with the alias. Returns ApiErrorNotFound when no configuration has it.
func (r DatadogConfigurationsResponse) findByAlias(alias string) (DatadogConfiguration, error) {
	for _, configuration := range r.Configurations {
		if configuration.Alias == alias {
			return configuration, nil
		}
	}
	return DatadogConfiguration{}, fmt.Errorf("datadog configuration %s: %w", alias, ApiErrorNotFound)
}

/***********************************************************************************************************************
 * GET /api/v1/datadog/configurations
 **********************************************************************************************************************/

func (c *DatadogConfigurationsClient) List(ctx context.Context) ([]DatadogConfiguration, error) {
	response := DatadogConfigurationsResponse{}
	apiError := ApiError{}

	body, err := c.Client().Get(Route("datadog", "configurations")).Receive(&response, &apiError)
	if err != nil {
		return nil, fmt.Errorf("failed listing datadog configurations: %+v", err)
	}

	err = c.client.handleResponseStatus(body, &apiError)
	if err != nil {
		return nil, fmt.Errorf("failed listing datadog configurations: %w", err)
	}
	return response.Configurations, nil
}

// Get finds a configuration by alias in the list response. GET /api/v1/datadog/configuration/:alias returns 400,
// not 404, for an unknown alias, so the list endpoint gives a reliable not-found signal. Returns ApiErrorNotFound
// when no configuration has the alias.
func (c *DatadogConfigurationsClient) Get(ctx context.Context, alias string) (DatadogConfiguration, error) {
	configurations, err := c.List(ctx)
	if err != nil {
		return DatadogConfiguration{}, err
	}
	return DatadogConfigurationsResponse{Configurations: configurations}.findByAlias(alias)
}

/***********************************************************************************************************************
 * POST /api/v1/datadog/configuration
 **********************************************************************************************************************/

type CreateDatadogConfigurationRequest struct {
	Alias           string   `json:"alias"`
	IsDefault       bool     `json:"isDefault"`
	ApiKey          string   `json:"apiKey"`
	AppKey          string   `json:"appKey"`
	Region          string   `json:"region"`
	Environments    []string `json:"environments"`
	CustomSubdomain string   `json:"customSubdomain,omitempty"`
}

func (c *DatadogConfigurationsClient) Create(ctx context.Context, req CreateDatadogConfigurationRequest) (DatadogConfiguration, error) {
	response := DatadogConfigurationsResponse{}
	apiError := ApiError{}

	body, err := c.Client().Post(Route("datadog", "configuration")).BodyJSON(&req).Receive(&response, &apiError)
	if err != nil {
		return DatadogConfiguration{}, fmt.Errorf("failed creating datadog configuration: %+v", err)
	}

	err = c.client.handleResponseStatus(body, &apiError)
	if err != nil {
		return DatadogConfiguration{}, err
	}
	return response.findByAlias(req.Alias)
}

/***********************************************************************************************************************
 * PUT /api/v1/datadog/configuration/:alias
 **********************************************************************************************************************/

// UpdateDatadogConfigurationRequest holds the only fields the public API can change. Alias is the new alias; the
// current alias goes in the path. The API keys, region, and custom subdomain cannot change after creation.
type UpdateDatadogConfigurationRequest struct {
	Alias        string   `json:"alias"`
	IsDefault    bool     `json:"isDefault"`
	Environments []string `json:"environments"`
}

func (c *DatadogConfigurationsClient) Update(ctx context.Context, alias string, req UpdateDatadogConfigurationRequest) (DatadogConfiguration, error) {
	response := DatadogConfigurationsResponse{}
	apiError := ApiError{}

	body, err := c.Client().Put(Route("datadog", "configuration/"+url.PathEscape(alias))).BodyJSON(&req).Receive(&response, &apiError)
	if err != nil {
		return DatadogConfiguration{}, fmt.Errorf("failed updating datadog configuration: %+v", err)
	}

	err = c.client.handleResponseStatus(body, &apiError)
	if err != nil {
		return DatadogConfiguration{}, err
	}
	return response.findByAlias(req.Alias)
}

/***********************************************************************************************************************
 * DELETE /api/v1/datadog/configuration/:alias
 **********************************************************************************************************************/

func (c *DatadogConfigurationsClient) Delete(ctx context.Context, alias string) error {
	apiError := ApiError{}

	body, err := c.Client().Delete(Route("datadog", "configuration/"+url.PathEscape(alias))).Receive(nil, &apiError)
	if err != nil {
		return fmt.Errorf("failed deleting datadog configuration: %+v", err)
	}

	return c.client.handleResponseStatus(body, &apiError)
}
