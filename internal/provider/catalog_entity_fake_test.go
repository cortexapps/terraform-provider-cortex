package provider_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"gopkg.in/yaml.v3"
)

// fakeCatalogApi is an in-memory Cortex API for the cortex_catalog_entity routes. Like the backend, it stores the
// descriptor as submitted and returns it on read, so a block the provider omits on write is absent on read.
type fakeCatalogApi struct {
	mu          sync.Mutex
	descriptors map[string]map[string]any
}

func newFakeCatalogApi(t *testing.T) (*fakeCatalogApi, string) {
	f := &fakeCatalogApi{descriptors: map[string]map[string]any{}}
	server := httptest.NewServer(f)
	t.Cleanup(server.Close)
	return f, server.URL
}

func (f *fakeCatalogApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	path := r.URL.Path
	switch {
	case r.Method == http.MethodPost && path == "/api/v1/open-api":
		body, _ := io.ReadAll(r.Body)
		var descriptor map[string]any
		if err := yaml.Unmarshal(body, &descriptor); err != nil {
			fakeCatalogFail(w, http.StatusBadRequest, "bad descriptor")
			return
		}
		info, _ := descriptor["info"].(map[string]any)
		tag, _ := info["x-cortex-tag"].(string)
		f.descriptors[tag] = descriptor
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"violations":[]}`))
	case r.Method == http.MethodGet && strings.HasSuffix(path, "/openapi") && r.URL.Query().Get("yaml") == "true":
		descriptor, ok := f.descriptors[strings.TrimSuffix(strings.TrimPrefix(path, "/api/v1/catalog/"), "/openapi")]
		if !ok {
			fakeCatalogFail(w, http.StatusNotFound, "entity not found")
			return
		}
		out, _ := yaml.Marshal(descriptor)
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(out)
	case r.Method == http.MethodDelete && strings.HasPrefix(path, "/api/v1/catalog/"):
		tag := strings.TrimPrefix(path, "/api/v1/catalog/")
		if _, ok := f.descriptors[tag]; !ok {
			fakeCatalogFail(w, http.StatusNotFound, "entity not found")
			return
		}
		delete(f.descriptors, tag)
		w.WriteHeader(http.StatusOK)
	default:
		fakeCatalogFail(w, http.StatusNotFound, "unexpected request "+r.Method+" "+path)
	}
}

// setInfo changes one info key of a stored descriptor, as if a client other than Terraform had changed it.
func (f *fakeCatalogApi) setInfo(tag, key string, value any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.descriptors[tag]["info"].(map[string]any)[key] = value
}

func fakeCatalogFail(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"type": "error", "message": message})
}

func catalogUnitConfig(url, body string) string {
	return fmt.Sprintf(`
provider "cortex" {
  base_api_url = %q
  token        = "test"
}
%s`, url, body)
}
