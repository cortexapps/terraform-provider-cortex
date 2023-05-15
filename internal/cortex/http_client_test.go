package cortex_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RequestTest func(req *http.Request)

func setupClient(requestPath string, mockedResponse interface{}, requestTests ...RequestTest) (*cortex.HttpClient, func(), error) {
	mux := http.NewServeMux()
	mux.HandleFunc(requestPath, func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		for _, test := range requestTests {
			test(req)
		}
		if err := json.NewEncoder(w).Encode(mockedResponse); err != nil {
			panic(fmt.Errorf("could not encode JSON: %w", err))
		}
	})

	ts := httptest.NewServer(mux)

	c, err := cortex.NewClient(
		cortex.WithURL(ts.URL),
		cortex.WithToken("test"),
		cortex.WithVersion("test"),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("could not build client: %w", err)
	}

	teardown := func() {
		ts.Close()
	}

	return c, teardown, nil
}

var pingResponseJSON = `{}`

func TestClientInitialization(t *testing.T) {
	var token string

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token = req.Header.Get("Authorization")
		length, _ := w.Write([]byte(pingResponseJSON))
		assert.Greater(t, length, 0)
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	testToken := "testing-123"
	c, err := cortex.NewClient(
		cortex.WithURL(ts.URL),
		cortex.WithToken(testToken),
	)

	assert.Nil(t, err, "received error initializing API client")

	err = c.Ping(context.Background())
	assert.Nil(t, err, "Received error hitting Ping endpoint")

	expectedAuthString := fmt.Sprintf("Bearer %s", testToken)
	assert.Equal(t, expectedAuthString, token, "Expected auth string to be %s, got %s", expectedAuthString, token)
}

func AssertRequestBody(t *testing.T, src interface{}) RequestTest {
	return func(req *http.Request) {
		t.Run("AssertRequestBody", func(t *testing.T) {
			body := io.NopCloser(req.Body)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(src)
			assert.Nil(t, err, "could not encode JSON")

			b, err := io.ReadAll(body)
			assert.Nil(t, err, "could not read request body")

			assert.True(t, bytes.Equal(buf.Bytes(), b), "expected request body to be %s, got %s", buf.String(), string(b))
		})
	}
}

func AssertRequestBodyYaml(t *testing.T, src interface{}) RequestTest {
	return func(req *http.Request) {
		t.Run("AssertRequestBodyRaw", func(t *testing.T) {
			bodyIo := io.NopCloser(req.Body)

			buf := new(bytes.Buffer)
			err := yaml.NewEncoder(buf).Encode(src)
			assert.Nil(t, err, "could not encode YAML")

			b, err := io.ReadAll(bodyIo)
			assert.Nil(t, err, "could not read request body")
			assert.Equal(t, buf.String(), string(b), "expected request body to be %s, got %s", buf.String(), string(b))
		})
	}
}

func AssertRequestMethod(t *testing.T, method string) RequestTest {
	return func(req *http.Request) {
		t.Run("AssertRequestMethod", func(t *testing.T) {
			assert.Equal(t, method, req.Method, "expected request method to be %s, got %s", method, req.Method)
		})
	}
}
