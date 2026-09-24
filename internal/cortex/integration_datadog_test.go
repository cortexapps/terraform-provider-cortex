package cortex_test

import (
	"context"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
)

// datadogConfigurationsResponse is the wrapper of list, create, and update responses.
type datadogConfigurationsResponse struct {
	Configurations []cortex.DatadogConfiguration `json:"configurations"`
}

var testDatadogConfiguration = cortex.DatadogConfiguration{
	Alias:           "test-datadog",
	IsDefault:       true,
	Environments:    []string{"prod", "staging"},
	Region:          "US1",
	CustomSubdomain: "acme",
	LastFourApiKey:  "iKey",
	LastFourAppKey:  "pKey",
}

var testDatadogConfigurationsResponse = datadogConfigurationsResponse{
	Configurations: []cortex.DatadogConfiguration{
		{Alias: "other-datadog", Region: "EU1", Environments: []string{}, LastFourApiKey: "aaaa", LastFourAppKey: "bbbb"},
		testDatadogConfiguration,
	},
}

func TestListDatadogConfigurations(t *testing.T) {
	c, teardown, err := setupClient(
		cortex.Route("datadog", "configurations"),
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.DatadogConfigurations().List(context.Background())
	assert.Nil(t, err, "error listing datadog configurations")
	assert.Len(t, res, 2)
	assert.Equal(t, testDatadogConfiguration, res[1])
}

func TestGetDatadogConfiguration(t *testing.T) {
	c, teardown, err := setupClient(
		cortex.Route("datadog", "configurations"),
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.DatadogConfigurations().Get(context.Background(), testDatadogConfiguration.Alias)
	assert.Nil(t, err, "error retrieving a datadog configuration")
	assert.Equal(t, testDatadogConfiguration, res)
}

func TestGetDatadogConfigurationNotFound(t *testing.T) {
	c, teardown, err := setupClient(
		cortex.Route("datadog", "configurations"),
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	_, err = c.DatadogConfigurations().Get(context.Background(), "missing-datadog")
	assert.ErrorIs(t, err, cortex.ApiErrorNotFound)
}

func TestCreateDatadogConfiguration(t *testing.T) {
	req := cortex.CreateDatadogConfigurationRequest{
		Alias:           testDatadogConfiguration.Alias,
		IsDefault:       true,
		ApiKey:          "fake-apiKey",
		AppKey:          "fake-appKey",
		Region:          "US1",
		Environments:    []string{"prod", "staging"},
		CustomSubdomain: "acme",
	}
	c, teardown, err := setupClient(
		cortex.Route("datadog", "configuration"),
		// The API returns every configuration in the tenant, not only the new one.
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.DatadogConfigurations().Create(context.Background(), req.Alias, req)
	assert.Nil(t, err, "error creating a datadog configuration")
	assert.Equal(t, testDatadogConfiguration, res)
}

func TestUpdateDatadogConfiguration(t *testing.T) {
	oldAlias := "old-datadog"
	req := cortex.UpdateDatadogConfigurationRequest{
		Alias:        testDatadogConfiguration.Alias,
		IsDefault:    true,
		Environments: []string{"prod", "staging"},
	}
	c, teardown, err := setupClient(
		cortex.Route("datadog", "configuration/"+oldAlias),
		// The API returns every configuration in the tenant, not only the updated one.
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.DatadogConfigurations().Update(context.Background(), oldAlias, req.Alias, req)
	assert.Nil(t, err, "error updating a datadog configuration")
	assert.Equal(t, testDatadogConfiguration, res)
}

func TestCreateDatadogConfigurationMissingFromResponse(t *testing.T) {
	req := cortex.CreateDatadogConfigurationRequest{Alias: "missing-datadog", Region: "US1", Environments: []string{}}
	c, teardown, err := setupClient(
		cortex.Route("datadog", "configuration"),
		testDatadogConfigurationsResponse,
		AssertRequestMethod(t, "POST"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	_, err = c.DatadogConfigurations().Create(context.Background(), req.Alias, req)
	assert.ErrorContains(t, err, "missing-datadog")
}

func TestDeleteDatadogConfiguration(t *testing.T) {
	c, teardown, err := setupClient(
		cortex.Route("datadog", "configuration/"+testDatadogConfiguration.Alias),
		map[string]interface{}{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.DatadogConfigurations().Delete(context.Background(), testDatadogConfiguration.Alias)
	assert.Nil(t, err, "error deleting a datadog configuration")
}
