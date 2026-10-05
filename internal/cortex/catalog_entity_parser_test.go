package cortex_test

import (
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
	_, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors:
        - id: 123
          alias: second-account
`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "second-account")
}

func TestYamlToEntityRejectsDatadogMonitorWithoutNumericId(t *testing.T) {
	for name, monitors := range map[string]string{
		"string":        `["abc"]`,
		"object no id":  `[{alias: x}]`,
		"object bad id": `[{id: abc}]`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := parseDescriptor(t, `
info:
  x-cortex-tag: test
  x-cortex-apm:
    datadog:
      monitors: `+monitors+`
`)
			assert.Error(t, err)
		})
	}
}
