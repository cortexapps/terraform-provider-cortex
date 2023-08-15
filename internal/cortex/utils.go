package cortex

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"reflect"
	"strings"
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
		return val
	}
	return defaultValue
}

func MapFetchToString(m map[string]interface{}, key string) string {
	return MapFetch(m, key, "").(string)
}

func InterfaceToString(v interface{}) (string, error) {
	value := ""
	// see if the value is a map or a scalar
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Map || rt.Kind() == reflect.Struct || rt.Kind() == reflect.Slice {
		// this is a map/struct/slice, let's convert to JSON
		sv, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		value = string(sv)
	} else {
		// otherwise it's a scalar (int/string) so we can just convert to string
		value = fmt.Sprintf("%v", v)
	}
	return value, nil
}

func StringToInterface(v string) (interface{}, error) {
	value := interface{}(nil)
	if strings.Contains(v, "{") && strings.Contains(v, "}") { // hacky stupid way
		err := json.Unmarshal([]byte(v), &value)
		if err != nil {
			return nil, err
		}
	} else {
		value = v
	}
	return value, nil
}
