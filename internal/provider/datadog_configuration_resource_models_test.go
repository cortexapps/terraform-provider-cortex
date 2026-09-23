package provider

import (
	"context"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatadogConfigurationResourceModel_FromApiModel(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}
	data := DatadogConfigurationResourceModel{
		ApiKey: types.StringValue("secret-api-key-1234"),
		AppKey: types.StringValue("secret-app-key-5678"),
	}

	data.FromApiModel(ctx, &diagnostics, cortex.DatadogConfiguration{
		Alias:           "datadog-prod",
		IsDefault:       true,
		Environments:    []string{"prod"},
		Region:          "EU1",
		CustomSubdomain: "acme",
		LastFourApiKey:  "1234",
		LastFourAppKey:  "5678",
	})

	require.False(t, diagnostics.HasError())
	assert.Equal(t, "datadog-prod", data.Id.ValueString())
	assert.Equal(t, "datadog-prod", data.Alias.ValueString())
	assert.True(t, data.IsDefault.ValueBool())
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("prod")}), data.Environments)
	assert.Equal(t, "EU1", data.Region.ValueString())
	assert.Equal(t, "acme", data.CustomSubdomain.ValueString())
	assert.Equal(t, "1234", data.ApiKeyLastFour.ValueString())
	assert.Equal(t, "5678", data.AppKeyLastFour.ValueString())
	// The API never returns the keys, so the model keeps the values it had.
	assert.Equal(t, "secret-api-key-1234", data.ApiKey.ValueString())
	assert.Equal(t, "secret-app-key-5678", data.AppKey.ValueString())
}

func TestDatadogConfigurationResourceModel_FromApiModel_EmptyOptionals(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}
	data := DatadogConfigurationResourceModel{}

	data.FromApiModel(ctx, &diagnostics, cortex.DatadogConfiguration{
		Alias:  "datadog-prod",
		Region: "US1",
	})

	require.False(t, diagnostics.HasError())
	assert.True(t, data.CustomSubdomain.IsNull())
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{}), data.Environments)
	assert.True(t, data.ApiKey.IsNull())
	assert.True(t, data.AppKey.IsNull())
}

func TestDatadogConfigurationResourceModel_ToCreateRequest(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}
	data := DatadogConfigurationResourceModel{
		Alias:           types.StringValue("datadog-prod"),
		IsDefault:       types.BoolUnknown(),
		ApiKey:          types.StringValue("api-key"),
		AppKey:          types.StringValue("app-key"),
		Region:          types.StringValue("US5"),
		Environments:    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("prod"), types.StringValue("staging")}),
		CustomSubdomain: types.StringNull(),
	}

	req := data.ToCreateRequest(ctx, &diagnostics)

	require.False(t, diagnostics.HasError())
	assert.Equal(t, cortex.CreateDatadogConfigurationRequest{
		Alias:        "datadog-prod",
		IsDefault:    false,
		ApiKey:       "api-key",
		AppKey:       "app-key",
		Region:       "US5",
		Environments: []string{"prod", "staging"},
	}, req)
}

func TestDatadogConfigurationResourceModel_ToCreateRequest_NullEnvironments(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}
	data := DatadogConfigurationResourceModel{
		Alias:        types.StringValue("datadog-prod"),
		Region:       types.StringValue("US1"),
		Environments: types.ListNull(types.StringType),
	}

	req := data.ToCreateRequest(ctx, &diagnostics)

	require.False(t, diagnostics.HasError())
	// The API requires a non-null environments list.
	assert.NotNil(t, req.Environments)
	assert.Empty(t, req.Environments)
}

func TestDatadogConfigurationResourceModel_ToUpdateRequest(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}
	data := DatadogConfigurationResourceModel{
		Alias:        types.StringValue("datadog-renamed"),
		IsDefault:    types.BoolValue(true),
		Environments: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("prod")}),
	}

	req := data.ToUpdateRequest(ctx, &diagnostics)

	require.False(t, diagnostics.HasError())
	assert.Equal(t, cortex.UpdateDatadogConfigurationRequest{
		Alias:        "datadog-renamed",
		IsDefault:    true,
		Environments: []string{"prod"},
	}, req)
}

func TestDatadogKeyRequiresReplace(t *testing.T) {
	tests := []struct {
		name          string
		stateKey      types.String
		configKey     types.String
		stateLastFour types.String
		want          bool
	}{
		{"unchanged key", types.StringValue("old-key-1234"), types.StringValue("old-key-1234"), types.StringValue("1234"), false},
		{"changed key", types.StringValue("old-key-1234"), types.StringValue("new-key-9999"), types.StringValue("1234"), true},
		{"changed key with the same last four", types.StringValue("old-key-1234"), types.StringValue("new-key-1234"), types.StringValue("1234"), true},
		{"key rotated outside terraform", types.StringValue("old-key-1234"), types.StringValue("old-key-1234"), types.StringValue("9999"), true},
		{"imported, config matches api", types.StringNull(), types.StringValue("old-key-1234"), types.StringValue("1234"), false},
		{"imported, config does not match api", types.StringNull(), types.StringValue("new-key-9999"), types.StringValue("1234"), true},
		{"unknown config key", types.StringValue("old-key-1234"), types.StringUnknown(), types.StringValue("1234"), true},
		{"no last four in state", types.StringValue("old-key-1234"), types.StringValue("old-key-1234"), types.StringNull(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datadogKeyRequiresReplace(tt.stateKey, tt.configKey, tt.stateLastFour))
		})
	}
}

func TestDatadogConfigurationResourceModel_ClearDriftedKeys(t *testing.T) {
	data := DatadogConfigurationResourceModel{
		ApiKey:         types.StringValue("old-key-1234"),
		AppKey:         types.StringValue("app-key-5678"),
		ApiKeyLastFour: types.StringValue("9999"),
		AppKeyLastFour: types.StringValue("5678"),
	}

	data.ClearDriftedKeys()

	// A key rotated outside Terraform no longer matches the API, so it leaves state and the next plan replaces it.
	assert.True(t, data.ApiKey.IsNull())
	assert.Equal(t, "app-key-5678", data.AppKey.ValueString())
}

func TestDatadogConfigurationResourceModel_ClearDriftedKeys_NullKeys(t *testing.T) {
	data := DatadogConfigurationResourceModel{
		ApiKey:         types.StringNull(),
		AppKey:         types.StringNull(),
		ApiKeyLastFour: types.StringValue("1234"),
		AppKeyLastFour: types.StringValue("5678"),
	}

	data.ClearDriftedKeys()

	assert.True(t, data.ApiKey.IsNull())
	assert.True(t, data.AppKey.IsNull())
}

func TestDatadogConfigurationResourceModel_ToCreateRequest_NeverDefault(t *testing.T) {
	ctx := context.Background()
	diagnostics := diag.Diagnostics{}
	data := DatadogConfigurationResourceModel{
		Alias:     types.StringValue("datadog-prod"),
		IsDefault: types.BoolValue(true),
		Region:    types.StringValue("US1"),
	}

	req := data.ToCreateRequest(ctx, &diagnostics)

	require.False(t, diagnostics.HasError())
	// The API rejects a new default while another default exists, so Create sets the default with a later update.
	assert.False(t, req.IsDefault)
}

func TestDatadogUnsetsDefault(t *testing.T) {
	tests := []struct {
		name        string
		stateValue  types.Bool
		configValue types.Bool
		want        bool
	}{
		{"default set to false", types.BoolValue(true), types.BoolValue(false), true},
		{"default kept", types.BoolValue(true), types.BoolValue(true), false},
		{"default not configured", types.BoolValue(true), types.BoolNull(), false},
		{"non-default set to false", types.BoolValue(false), types.BoolValue(false), false},
		{"non-default set to true", types.BoolValue(false), types.BoolValue(true), false},
		{"no state", types.BoolNull(), types.BoolValue(false), false},
		{"unknown config", types.BoolValue(true), types.BoolUnknown(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datadogUnsetsDefault(tt.stateValue, tt.configValue))
		})
	}
}
