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
	mux := http.NewServeMux()
	mux.HandleFunc(cortex.Route("pagerduty", "default-configuration"), func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"No configuration"}`, http.StatusNotFound)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	c, err := cortex.NewClient(cortex.WithURL(server.URL), cortex.WithToken("test"), cortex.WithVersion("test"))
	assert.Nil(t, err)

	_, err = c.PagerDutyConfiguration().Get(context.Background())
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
	mux := http.NewServeMux()
	mux.HandleFunc(cortex.Route("pagerduty", "default-configuration"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"lastFour":"a1b2","isTokenReadonly":"not-a-bool"}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	c, err := cortex.NewClient(cortex.WithURL(server.URL), cortex.WithToken("test"), cortex.WithVersion("test"))
	assert.Nil(t, err)

	res, err := c.PagerDutyConfiguration().Get(context.Background())
	assert.NotNil(t, err)
	assert.Equal(t, cortex.PagerDutyConfiguration{}, res)
}
