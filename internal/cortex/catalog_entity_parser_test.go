package cortex_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func parseDescriptor(t *testing.T, descriptor string) (cortex.CatalogEntityData, error) {
	t.Helper()
	var raw map[string]interface{}
	require.NoError(t, yaml.Unmarshal([]byte(descriptor), &raw))
	parser := cortex.CatalogEntityParser{}
	return parser.YamlToEntity(raw)
}

func TestYamlToEntityReadsDatadogMonitorsAsIntegers(t *testing.T) {
	entity, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors: [123, 9007199254740993]
`)
	require.NoError(t, err)
	assert.Equal(t, []int64{123, 9007199254740993}, entity.Apm.DataDog.Monitors)
}

func TestYamlToEntityReadsDatadogMonitorsAsObjects(t *testing.T) {
	entity, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors:
        - id: 123
        - 456
        - id: 789
          alias: null
`)
	require.NoError(t, err)
	assert.Equal(t, []int64{123, 456, 789}, entity.Apm.DataDog.Monitors)
}

func TestYamlToEntityRejectsDatadogMonitorWithAlias(t *testing.T) {
	for name, alias := range map[string]string{
		"string":  `second-account`,
		"number":  `2024`,
		"boolean": `true`,
		"empty":   `""`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors:
        - id: 123
          alias: `+alias+`
`)
			require.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("alias %q", strings.Trim(alias, `"`)))
		})
	}
}

func TestYamlToEntityReportsParsedIDForDatadogMonitorWithAlias(t *testing.T) {
	_, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors: [{id: "456", alias: x}]
`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "datadog monitor 456 uses alias")
}

func TestYamlToEntityReadsDatadogMonitorsAsInt64(t *testing.T) {
	parser := cortex.CatalogEntityParser{}
	entity, err := parser.YamlToEntity(map[string]interface{}{
		"info": map[string]interface{}{
			"x-cortex-tag": "test",
			"x-cortex-apm": map[string]interface{}{
				"datadog": map[string]interface{}{"monitors": []interface{}{int64(5), map[string]interface{}{"id": int64(6)}}},
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, []int64{5, 6}, entity.Apm.DataDog.Monitors)
}

func TestYamlToEntityReadsDatadogMonitorObjectIDAsString(t *testing.T) {
	entity, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors: [{id: "456"}]
`)
	require.NoError(t, err)
	assert.Equal(t, []int64{456}, entity.Apm.DataDog.Monitors)
}

func TestYamlToEntityRejectsDatadogMonitorWithoutIntegerID(t *testing.T) {
	for name, tc := range map[string]struct{ monitors, message string }{
		"null":          {`[null]`, "has no id"},
		"object no id":  {`[{}]`, "has no id"},
		"object null":   {`[{id: null}]`, "has no id"},
		"bare string":   {`["123"]`, "must be a number or an object"},
		"object string": {`[{id: abc}]`, "is not an integer"},
		"exponent":      {`[1e3]`, "is not an integer (float64)"},
		"fraction":      {`[1.5]`, "is not an integer"},
		"infinity":      {`[.inf]`, "is not an integer"},
		"too large":     {`[18446744073709551615]`, "is not an integer"},
		"list":          {`[[1]]`, "is not an integer"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors: `+tc.monitors+`
`)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.message)
		})
	}
}

func TestYamlToEntityRejectsMalformedDatadogApm(t *testing.T) {
	for name, tc := range map[string]struct{ datadog, message string }{
		"datadog not a map":   {`datadog: [1]`, "datadog apm is not an object"},
		"monitors not a list": {`datadog: {monitors: 123}`, "datadog monitors is not a list"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    `+tc.datadog+`
`)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.message)
		})
	}
}
