package integrations_test

import (
	"context"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/stretchr/testify/assert"
)

// datadogConfigurationsResponse is the wrapper of list, create, and update responses.
type datadogConfigurationsResponse struct {
	Configurations []integrations.DatadogConfiguration `json:"configurations"`
}

var testDatadogConfiguration = integrations.DatadogConfiguration{
	Alias:           "test-datadog",
	IsDefault:       true,
	Environments:    []string{"prod", "staging"},
	Region:          "US1",
	CustomSubdomain: "acme",
	LastFourApiKey:  "iKey",
	LastFourAppKey:  "pKey",
}

var testDatadogConfigurationsResponse = datadogConfigurationsResponse{
	Configurations: []integrations.DatadogConfiguration{
		{Alias: "other-datadog", Region: "EU1", Environments: []string{}, LastFourApiKey: "aaaa", LastFourAppKey: "bbbb"},
		testDatadogConfiguration,
	},
}

func TestListDatadogConfigurations(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/datadog/configurations",
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := integrations.Datadog(c).List(context.Background())
	assert.Nil(t, err, "error listing datadog configurations")
	assert.Len(t, res, 2)
	assert.Equal(t, testDatadogConfiguration, res[1])
}

func TestGetDatadogConfiguration(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/datadog/configurations",
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := integrations.Datadog(c).Get(context.Background(), testDatadogConfiguration.Alias)
	assert.Nil(t, err, "error retrieving a datadog configuration")
	assert.Equal(t, testDatadogConfiguration, res)
}

func TestGetDatadogConfigurationNotFound(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/datadog/configurations",
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	_, err = integrations.Datadog(c).Get(context.Background(), "missing-datadog")
	assert.ErrorIs(t, err, cortex.ApiErrorNotFound)
}

func TestCreateDatadogConfiguration(t *testing.T) {
	req := integrations.CreateDatadogConfigurationRequest{
		Alias:           testDatadogConfiguration.Alias,
		IsDefault:       true,
		ApiKey:          "fake-apiKey",
		AppKey:          "fake-appKey",
		Region:          "US1",
		Environments:    []string{"prod", "staging"},
		CustomSubdomain: "acme",
	}
	c, teardown, err := setupClient(
		"/api/v1/datadog/configuration",
		// The API returns every configuration in the tenant, not only the new one.
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := integrations.Datadog(c).Create(context.Background(), req.Alias, req)
	assert.Nil(t, err, "error creating a datadog configuration")
	assert.Equal(t, testDatadogConfiguration, res)
}

func TestUpdateDatadogConfiguration(t *testing.T) {
	oldAlias := "old-datadog"
	req := integrations.UpdateDatadogConfigurationRequest{
		Alias:           testDatadogConfiguration.Alias,
		IsDefault:       true,
		Environments:    []string{"prod", "staging"},
		ApiKey:          "fake-api-key-e5f6",
		AppKey:          "fake-app-key-g7h8",
		Region:          "EU1",
		CustomSubdomain: "acme",
	}
	c, teardown, err := setupClient(
		"/api/v1/datadog/configuration/"+oldAlias,
		// The API returns every configuration in the tenant, not only the updated one.
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := integrations.Datadog(c).Update(context.Background(), oldAlias, req.Alias, req)
	assert.Nil(t, err, "error updating a datadog configuration")
	assert.Equal(t, testDatadogConfiguration, res)
}

func TestCreateDatadogConfigurationMissingFromResponse(t *testing.T) {
	req := integrations.CreateDatadogConfigurationRequest{Alias: "missing-datadog", Region: "US1", Environments: []string{}}
	c, teardown, err := setupClient(
		"/api/v1/datadog/configuration",
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "POST"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	_, err = integrations.Datadog(c).Create(context.Background(), req.Alias, req)
	assert.ErrorContains(t, err, "missing-datadog")
}

func TestDeleteDatadogConfiguration(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/datadog/configuration/"+testDatadogConfiguration.Alias,
		map[string]interface{}{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = integrations.Datadog(c).Delete(context.Background(), testDatadogConfiguration.Alias)
	assert.Nil(t, err, "error deleting a datadog configuration")
}
