package provider_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

// fakeIntegrationSpec describes how the fake stores and returns one integration.
type fakeIntegrationSpec struct {
	masks map[string]string // secret request field -> response field with its last four characters
}

// fakeSpecs follow the backend contract. The key is the path segment under /api/v1.
var fakeSpecs = map[string]fakeIntegrationSpec{
	"datadog": {masks: map[string]string{"apiKey": "lastFourApiKey", "appKey": "lastFourAppKey"}},
}

// fakeCortexApi is an in-memory Cortex API for the integration configuration routes. It applies the same alias and
// default rules as the backend.
type fakeCortexApi struct {
	mu      sync.Mutex
	configs map[string][]map[string]any
	creates map[string]int
	// deleteNotFound makes deletes answer 404, like a configuration that another client deleted after the refresh.
	deleteNotFound bool
}

func newFakeCortexApi(t *testing.T) (*fakeCortexApi, string) {
	f := &fakeCortexApi{configs: map[string][]map[string]any{}, creates: map[string]int{}}
	server := httptest.NewServer(f)
	t.Cleanup(server.Close)
	return f, server.URL
}

func (f *fakeCortexApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	seg, route, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, "/api/v1/"), "/")
	if _, ok := fakeSpecs[seg]; !ok {
		fail(w, http.StatusNotFound, "unexpected request "+r.URL.Path)
		return
	}
	var body map[string]any
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			fail(w, http.StatusBadRequest, "bad body")
			return
		}
	}
	f.serveMulti(w, r, seg, route, body)
}

func (f *fakeCortexApi) serveMulti(w http.ResponseWriter, r *http.Request, seg, route string, body map[string]any) {
	alias, _ := url.PathUnescape(strings.TrimPrefix(route, "configuration/"))
	hasAlias := strings.HasPrefix(route, "configuration/")
	switch {
	case r.Method == http.MethodGet && route == "configurations":
		f.writeAll(w, seg)
	case r.Method == http.MethodPost && route == "configuration":
		f.creates[seg]++
		newAlias, _ := body["alias"].(string)
		if f.find(seg, newAlias) >= 0 {
			fail(w, http.StatusBadRequest, "Configuration exists with that alias")
			return
		}
		isDefault, _ := body["isDefault"].(bool)
		hasDefault := f.defaultIndex(seg) >= 0
		if isDefault && hasDefault {
			fail(w, http.StatusBadRequest, "A default configuration is already set")
			return
		}
		body["isDefault"] = isDefault || !hasDefault
		f.configs[seg] = append(f.configs[seg], body)
		f.writeAll(w, seg)
	case r.Method == http.MethodPut && hasAlias:
		i := f.find(seg, alias)
		if i < 0 {
			fail(w, http.StatusNotFound, "Alias did not match any configuration")
			return
		}
		newAlias, _ := body["alias"].(string)
		if newAlias != alias && f.find(seg, newAlias) >= 0 {
			fail(w, http.StatusBadRequest, "Configuration exists with that alias")
			return
		}
		isDefault, _ := body["isDefault"].(bool)
		if isDefault {
			for j := range f.configs[seg] {
				f.configs[seg][j]["isDefault"] = j == i
			}
		} else if f.configs[seg][i]["isDefault"] == true {
			fail(w, http.StatusBadRequest, "You must change the default to a different configuration before updating")
			return
		}
		for k, v := range body {
			if v == nil || k == "isDefault" {
				continue
			}
			f.configs[seg][i][k] = v
		}
		f.writeAll(w, seg)
	case r.Method == http.MethodDelete && hasAlias:
		i := f.find(seg, alias)
		if i < 0 || f.deleteNotFound {
			fail(w, http.StatusNotFound, "Unable to find configuration")
			return
		}
		if f.configs[seg][i]["isDefault"] == true && len(f.configs[seg]) > 1 {
			fail(w, http.StatusBadRequest, "You must change the default to a different configuration before deleting")
			return
		}
		f.configs[seg] = append(f.configs[seg][:i], f.configs[seg][i+1:]...)
		w.WriteHeader(http.StatusOK)
	default:
		fail(w, http.StatusNotFound, "unexpected request")
	}
}

// render returns the configuration as the API does: no secrets, and the last four characters of each secret.
func (f *fakeCortexApi) render(seg string, cfg map[string]any) map[string]any {
	spec := fakeSpecs[seg]
	out := map[string]any{}
	for k, v := range cfg {
		if _, secret := spec.masks[k]; secret || v == nil {
			continue
		}
		out[k] = v
	}
	for field, mask := range spec.masks {
		if v, ok := cfg[field].(string); ok && v != "" {
			out[mask] = takeLastFour(v)
		}
	}
	return out
}

func (f *fakeCortexApi) writeAll(w http.ResponseWriter, seg string) {
	all := []any{}
	for _, c := range f.configs[seg] {
		all = append(all, f.render(seg, c))
	}
	writeJSON(w, map[string]any{"configurations": all})
}

func (f *fakeCortexApi) find(seg, alias string) int {
	for i, c := range f.configs[seg] {
		if c["alias"] == alias {
			return i
		}
	}
	return -1
}

func (f *fakeCortexApi) defaultIndex(seg string) int {
	for i, c := range f.configs[seg] {
		if c["isDefault"] == true {
			return i
		}
	}
	return -1
}

/***********************************************************************************************************************
 * Test controls: change the fake outside Terraform
 **********************************************************************************************************************/

// get returns a copy of a stored configuration.
func (f *fakeCortexApi) get(seg, alias string) (map[string]any, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	i := f.find(seg, alias)
	if i < 0 {
		return nil, false
	}
	out := map[string]any{}
	for k, v := range f.configs[seg][i] {
		out[k] = v
	}
	return out, true
}

func (f *fakeCortexApi) setField(seg, alias, field string, value any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.configs[seg][f.find(seg, alias)][field] = value
}

func (f *fakeCortexApi) setDefault(seg, alias string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	i := f.find(seg, alias)
	for j := range f.configs[seg] {
		f.configs[seg][j]["isDefault"] = j == i
	}
}

func (f *fakeCortexApi) remove(seg, alias string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if i := f.find(seg, alias); i >= 0 {
		f.configs[seg] = append(f.configs[seg][:i], f.configs[seg][i+1:]...)
	}
}

func (f *fakeCortexApi) answerDeletesWithNotFound() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.deleteNotFound = true
}

func (f *fakeCortexApi) seed(seg string, cfg map[string]any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.configs[seg] = append(f.configs[seg], cfg)
}

func (f *fakeCortexApi) createCount(seg string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.creates[seg]
}

/***********************************************************************************************************************
 * Helpers
 **********************************************************************************************************************/

func fail(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// takeLastFour matches the backend, which keeps the last four characters of a key, or all of a shorter key.
func takeLastFour(key string) string {
	if len(key) <= 4 {
		return key
	}
	return key[len(key)-4:]
}

// unitConfig adds a provider block that points at the fake API.
func unitConfig(url, body string) string {
	return fmt.Sprintf(`
provider "cortex" {
  base_api_url = %q
  token        = "test"
}
%s`, url, body)
}
