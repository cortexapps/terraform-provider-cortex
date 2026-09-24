package integrations_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/stretchr/testify/assert"
)

func TestSingleInstanceGet(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/pagerduty/default-configuration",
		integrations.PagerDutyConfiguration{LastFour: "a1b2", IsTokenReadonly: true},
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := integrations.PagerDuty(c).Get(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, integrations.PagerDutyConfiguration{LastFour: "a1b2", IsTokenReadonly: true}, res)
}

func TestSingleInstanceGetNotFound(t *testing.T) {
	c := errorClient(t, http.StatusNotFound, `{"message":"No configuration"}`)

	_, err := integrations.PagerDuty(c).Get(context.Background())
	assert.ErrorIs(t, err, cortex.ApiErrorNotFound)
}

func TestSingleInstanceCreate(t *testing.T) {
	req := integrations.PagerDutyConfigurationRequest{Token: "fake-token-a1b2", IsTokenReadonly: true}
	c, teardown, err := setupClient(
		"/api/v1/pagerduty/configuration",
		map[string]any{"configurations": []integrations.PagerDutyConfiguration{{LastFour: "a1b2", IsTokenReadonly: true}}},
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := integrations.PagerDuty(c).Create(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "a1b2", res.LastFour)
}

func TestSingleInstanceReplace(t *testing.T) {
	req := integrations.PagerDutyConfigurationRequest{Token: "fake-token-c3d4", IsTokenReadonly: false}
	c, teardown, err := setupClient(
		"/api/v1/pagerduty/configuration",
		map[string]any{"configurations": []integrations.PagerDutyConfiguration{{LastFour: "c3d4"}}},
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := integrations.PagerDuty(c).Replace(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "c3d4", res.LastFour)
}

func TestSingleInstanceCreateEmptyResponse(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/pagerduty/configuration",
		map[string]any{"configurations": []integrations.PagerDutyConfiguration{}},
		AssertRequestMethod(t, "POST"),
	)
	assert.Nil(t, err)
	defer teardown()

	_, err = integrations.PagerDuty(c).Create(context.Background(), integrations.PagerDutyConfigurationRequest{Token: "x"})
	assert.ErrorContains(t, err, "no configuration")
}

func TestSingleInstanceDelete(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/pagerduty/configurations",
		map[string]any{"configurations": []any{}},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err)
	defer teardown()

	assert.Nil(t, integrations.PagerDuty(c).Delete(context.Background()))
}

func TestSingleInstanceGetDecodeErrorReturnsZeroValue(t *testing.T) {
	c := errorClient(t, http.StatusOK, `{"lastFour":"a1b2","isTokenReadonly":"not-a-bool"}`)

	res, err := integrations.PagerDuty(c).Get(context.Background())
	assert.NotNil(t, err)
	assert.Equal(t, integrations.PagerDutyConfiguration{}, res)
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

	_, err := integrations.Datadog(rejected).Create(ctx, "dd", integrations.CreateDatadogConfigurationRequest{Alias: "dd"})
	assert.ErrorContains(t, err, "Configuration exists with that alias")
	assert.NotErrorIs(t, err, cortex.ApiErrorNotFound)

	_, err = integrations.Datadog(rejected).Update(ctx, "dd", "dd", integrations.UpdateDatadogConfigurationRequest{Alias: "dd"})
	assert.ErrorContains(t, err, "Configuration exists with that alias")

	_, err = integrations.Datadog(errorClient(t, http.StatusInternalServerError, `{"message":"boom"}`)).List(ctx)
	assert.ErrorContains(t, err, "500")

	err = integrations.Datadog(errorClient(t, http.StatusNotFound, `{"message":"Unable to find configuration"}`)).Delete(ctx, "dd")
	assert.ErrorIs(t, err, cortex.ApiErrorNotFound)
}
