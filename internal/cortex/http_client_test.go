package cortex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RequestTest func(req *http.Request)

func setupClient(requestPath string, mockedResponse interface{}, requestTests ...RequestTest) (*HttpClient, func(), error) {
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

	c, err := NewClient(
		WithURL(ts.URL),
		WithToken("test"),
		WithVersion("test"),
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
		w.Write([]byte(pingResponseJSON))
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	testToken := "testing-123"
	c, err := NewClient(WithURL(ts.URL), WithToken(testToken))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	err = c.Ping(context.Background())
	if err != nil {
		t.Fatalf("Received error hitting Ping endpoint: %s", err.Error())
	}

	if expected := "Bearer " + testToken; expected != token {
		t.Fatalf("Expected %s, Got: %s for bearer token", expected, token)
	}
}

func AssertRequestBody(t *testing.T, src interface{}) RequestTest {
	return func(req *http.Request) {
		t.Run("AssertRequestBody", func(t *testing.T) {
			body := io.NopCloser(req.Body)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(src)
			if err != nil {
				t.Fatalf("could not encode JSON: %s", err)
			}

			b, err := io.ReadAll(body)
			if err != nil {
				t.Fatalf("could not read request body: %s", err)
			}

			if !bytes.Equal(buf.Bytes(), b) {
				t.Fatalf("expected request body to be %s, got %s", buf.String(), string(b))
			}
		})
	}
}

func AssertRequestMethod(t *testing.T, method string) RequestTest {
	return func(req *http.Request) {
		t.Run("AssertRequestMethod", func(t *testing.T) {
			if method != req.Method {
				t.Fatalf("expected request method to be %s, got %s", method, req.Method)
			}
		})
	}
}
