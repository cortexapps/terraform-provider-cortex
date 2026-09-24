package integrations_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
)

// The helpers below are copies of the ones in the cortex_test package, which other packages cannot import.

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
		cortex.WithContext(context.Background()),
		cortex.WithURL(ts.URL),
		cortex.WithToken("test"),
		cortex.WithVersion("test"),
	)
	if err != nil {
		ts.Close()
		return nil, nil, fmt.Errorf("could not build client: %w", err)
	}
	return c, ts.Close, nil
}

func AssertRequestBody(t *testing.T, src interface{}) RequestTest {
	return func(req *http.Request) {
		t.Run("AssertRequestBody", func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(src)
			assert.Nil(t, err, "could not encode JSON")

			b, err := io.ReadAll(req.Body)
			assert.Nil(t, err, "could not read request body")

			assert.True(t, bytes.Equal(buf.Bytes(), b), "expected request body to be %s, got %s", buf.String(), string(b))
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
