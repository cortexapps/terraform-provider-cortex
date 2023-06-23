package cortex

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"net/http"
)

// yamlDecoder decodes http response YAML into a YAML-tagged struct value.
type yamlDecoder struct{}

// Decode decodes the Response Body into the value pointed to by v.
// Caller must provide a non-nil v and close the resp.Body.
func (d yamlDecoder) Decode(resp *http.Response, v interface{}) error {
	return yaml.NewDecoder(resp.Body).Decode(v)
}

// jsonDecoder decodes http response JSON into a JSON-tagged struct value.
type jsonDecoder struct{}

// Decode decodes the Response Body into the value pointed to by v.
// Caller must provide a non-nil v and close the resp.Body.
func (d jsonDecoder) Decode(resp *http.Response, v interface{}) error {
	return json.NewDecoder(resp.Body).Decode(v)
}

func MapFetch(m map[string]interface{}, key string, defaultValue any) any {
	if val, ok := m[key]; ok {
		return val.(string)
	}
	return defaultValue
}

func MapFetchToString(m map[string]interface{}, key string) string {
	return MapFetch(m, key, "").(string)
}
