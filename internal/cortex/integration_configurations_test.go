package cortex_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
)

func TestSingleInstanceGet(t *testing.T) {
	c, teardown, err := setupClient(
		cortex.Route("pagerduty", "default-configuration"),
		cortex.PagerDutyConfiguration{LastFour: "a1b2", IsTokenReadonly: true},
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := c.PagerDutyConfiguration().Get(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, cortex.PagerDutyConfiguration{LastFour: "a1b2", IsTokenReadonly: true}, res)
}

func TestSingleInstanceGetNotFound(t *testing.T) {
	c := errorClient(t, http.StatusNotFound, `{"message":"No configuration"}`)

	_, err := c.PagerDutyConfiguration().Get(context.Background())
	assert.ErrorIs(t, err, cortex.ApiErrorNotFound)
}

func TestSingleInstanceCreate(t *testing.T) {
	req := cortex.PagerDutyConfigurationRequest{Token: "fake-token-a1b2", IsTokenReadonly: true}
	c, teardown, err := setupClient(
		cortex.Route("pagerduty", "configuration"),
		map[string]any{"configurations": []cortex.PagerDutyConfiguration{{LastFour: "a1b2", IsTokenReadonly: true}}},
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := c.PagerDutyConfiguration().Create(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "a1b2", res.LastFour)
}

func TestSingleInstanceReplace(t *testing.T) {
	req := cortex.PagerDutyConfigurationRequest{Token: "fake-token-c3d4", IsTokenReadonly: false}
	c, teardown, err := setupClient(
		cortex.Route("pagerduty", "configuration"),
		map[string]any{"configurations": []cortex.PagerDutyConfiguration{{LastFour: "c3d4"}}},
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := c.PagerDutyConfiguration().Replace(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "c3d4", res.LastFour)
}

func TestSingleInstanceCreateEmptyResponse(t *testing.T) {
	c, teardown, err := setupClient(
		cortex.Route("pagerduty", "configuration"),
		map[string]any{"configurations": []cortex.PagerDutyConfiguration{}},
		AssertRequestMethod(t, "POST"),
	)
	assert.Nil(t, err)
	defer teardown()

	_, err = c.PagerDutyConfiguration().Create(context.Background(), cortex.PagerDutyConfigurationRequest{Token: "x"})
	assert.ErrorContains(t, err, "no configuration")
}

func TestSingleInstanceDelete(t *testing.T) {
	c, teardown, err := setupClient(
		cortex.Route("pagerduty", "configurations"),
		map[string]any{"configurations": []any{}},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err)
	defer teardown()

	assert.Nil(t, c.PagerDutyConfiguration().Delete(context.Background()))
}

func TestSingleInstanceGetDecodeErrorReturnsZeroValue(t *testing.T) {
	c := errorClient(t, http.StatusOK, `{"lastFour":"a1b2","isTokenReadonly":"not-a-bool"}`)

	res, err := c.PagerDutyConfiguration().Get(context.Background())
	assert.NotNil(t, err)
	assert.Equal(t, cortex.PagerDutyConfiguration{}, res)
}

// errorClient returns a client whose server answers every request with the status and body.
func errorClient(t *testing.T, status int, body string) *cortex.HttpClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	c, err := cortex.NewClient(cortex.WithURL(server.URL), cortex.WithToken("test"), cortex.WithVersion("test"))
	assert.Nil(t, err)
	return c
}

func TestMultiInstanceApiErrors(t *testing.T) {
	ctx := context.Background()
	rejected := errorClient(t, http.StatusBadRequest, `{"message":"Configuration exists with that alias"}`)

	_, err := rejected.DatadogConfigurations().Create(ctx, "dd", cortex.CreateDatadogConfigurationRequest{Alias: "dd"})
	assert.ErrorContains(t, err, "Configuration exists with that alias")
	assert.NotErrorIs(t, err, cortex.ApiErrorNotFound)

	_, err = rejected.DatadogConfigurations().Update(ctx, "dd", "dd", cortex.UpdateDatadogConfigurationRequest{Alias: "dd"})
	assert.ErrorContains(t, err, "Configuration exists with that alias")

	_, err = errorClient(t, http.StatusInternalServerError, `{"message":"boom"}`).DatadogConfigurations().List(ctx)
	assert.ErrorContains(t, err, "500")

	err = errorClient(t, http.StatusNotFound, `{"message":"Unable to find configuration"}`).DatadogConfigurations().Delete(ctx, "dd")
	assert.ErrorIs(t, err, cortex.ApiErrorNotFound)
}
