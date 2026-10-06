package cortex_test

import (
	"fmt"
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
      monitors: [123, 456]
`)
	require.NoError(t, err)
	assert.Equal(t, []int64{123, 456}, entity.Apm.DataDog.Monitors)
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
			assert.Contains(t, err.Error(), fmt.Sprintf("alias %q", alias))
		})
	}
}

func TestYamlToEntityReadsDatadogMonitorIDAsString(t *testing.T) {
	entity, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors: ["123", {id: "456"}]
`)
	require.NoError(t, err)
	assert.Equal(t, []int64{123, 456}, entity.Apm.DataDog.Monitors)
}

func TestYamlToEntityRejectsDatadogMonitorWithoutIntegerID(t *testing.T) {
	for name, tc := range map[string]struct{ monitors, message string }{
		"null":          {`[null]`, "has no id"},
		"object no id":  {`[{}]`, "has no id"},
		"object null":   {`[{id: null}]`, "has no id"},
		"string":        {`["abc"]`, "is not an integer"},
		"object bad id": {`[{id: abc}]`, "is not an integer"},
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
	for name, datadog := range map[string]string{
		"datadog not a map":   `datadog: [1]`,
		"monitors not a list": `datadog: {monitors: 123}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    `+datadog+`
`)
			assert.Error(t, err)
		})
	}
}
