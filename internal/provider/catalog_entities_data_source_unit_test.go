package provider

import (
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlattenCatalogEntitySlackChannels(t *testing.T) {
	channels := []cortex.CatalogEntitySlackChannel{
		{
			Name:                 "C0123456789",
			Description:          "x-monitoring-alerts-channel",
			NotificationsEnabled: true,
		},
	}

	result := flattenCatalogEntitySlackChannels(channels)

	require.Len(t, result, 1)
	assert.Equal(t, "C0123456789", result[0].Name.ValueString())
	assert.Equal(t, "x-monitoring-alerts-channel", result[0].Description.ValueString())
	assert.True(t, result[0].NotificationsEnabled.ValueBool())
}
