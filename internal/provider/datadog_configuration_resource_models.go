package provider

import (
	"context"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/***********************************************************************************************************************
 * Models
 **********************************************************************************************************************/

// DatadogConfigurationResourceModel describes the Datadog integration configuration data model within Terraform.
type DatadogConfigurationResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Alias           types.String `tfsdk:"alias"`
	IsDefault       types.Bool   `tfsdk:"is_default"`
	ApiKey          types.String `tfsdk:"api_key"`
	AppKey          types.String `tfsdk:"app_key"`
	Region          types.String `tfsdk:"region"`
	Environments    types.List   `tfsdk:"environments"`
	CustomSubdomain types.String `tfsdk:"custom_subdomain"`
	ApiKeyLastFour  types.String `tfsdk:"api_key_last_four"`
	AppKeyLastFour  types.String `tfsdk:"app_key_last_four"`
}

// FromApiModel maps the API response into the model. The API never returns the keys, so ApiKey and AppKey keep
// their current values.
func (r *DatadogConfigurationResourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity cortex.DatadogConfiguration) {
	r.Id = types.StringValue(entity.Alias)
	r.Alias = types.StringValue(entity.Alias)
	r.IsDefault = types.BoolValue(entity.IsDefault)
	r.Region = types.StringValue(entity.Region)
	r.ApiKeyLastFour = types.StringValue(entity.LastFourApiKey)
	r.AppKeyLastFour = types.StringValue(entity.LastFourAppKey)
	if entity.CustomSubdomain != "" {
		r.CustomSubdomain = types.StringValue(entity.CustomSubdomain)
	} else {
		r.CustomSubdomain = types.StringNull()
	}

	environments := entity.Environments
	if environments == nil {
		environments = []string{}
	}
	list, d := types.ListValueFrom(ctx, types.StringType, environments)
	diagnostics.Append(d...)
	r.Environments = list
}

func (r *DatadogConfigurationResourceModel) ToCreateRequest(ctx context.Context, diagnostics *diag.Diagnostics) cortex.CreateDatadogConfigurationRequest {
	return cortex.CreateDatadogConfigurationRequest{
		Alias:           r.Alias.ValueString(),
		IsDefault:       r.IsDefault.ValueBool(),
		ApiKey:          r.ApiKey.ValueString(),
		AppKey:          r.AppKey.ValueString(),
		Region:          r.Region.ValueString(),
		Environments:    r.environmentsToApi(ctx, diagnostics),
		CustomSubdomain: r.CustomSubdomain.ValueString(),
	}
}

func (r *DatadogConfigurationResourceModel) ToUpdateRequest(ctx context.Context, diagnostics *diag.Diagnostics) cortex.UpdateDatadogConfigurationRequest {
	return cortex.UpdateDatadogConfigurationRequest{
		Alias:        r.Alias.ValueString(),
		IsDefault:    r.IsDefault.ValueBool(),
		Environments: r.environmentsToApi(ctx, diagnostics),
	}
}

// environmentsToApi always returns a non-nil slice, because the API rejects a null environments list.
func (r *DatadogConfigurationResourceModel) environmentsToApi(ctx context.Context, diagnostics *diag.Diagnostics) []string {
	environments := []string{}
	if r.Environments.IsNull() || r.Environments.IsUnknown() {
		return environments
	}
	diagnostics.Append(r.Environments.ElementsAs(ctx, &environments, false)...)
	return environments
}

/***********************************************************************************************************************
 * Key drift
 **********************************************************************************************************************/

// ClearDriftedKeys removes a key from the model when its last four characters no longer match the API value, for
// example after a rotation outside Terraform. Terraform replaces a resource only when a RequiresReplace value
// changes, so the null value makes the next plan show a diff and a replace. Call it only on Read, after
// FromApiModel.
func (r *DatadogConfigurationResourceModel) ClearDriftedKeys() {
	if keyDrifted(r.ApiKey, r.ApiKeyLastFour) {
		r.ApiKey = types.StringNull()
	}
	if keyDrifted(r.AppKey, r.AppKeyLastFour) {
		r.AppKey = types.StringNull()
	}
}

func keyDrifted(key types.String, apiLastFour types.String) bool {
	if key.IsNull() || key.IsUnknown() || apiLastFour.IsNull() || apiLastFour.IsUnknown() {
		return false
	}
	return lastFour(key.ValueString()) != apiLastFour.ValueString()
}

// datadogKeyRequiresReplace tells if a change to api_key or app_key needs a new configuration. The public API
// cannot update keys in place, and it returns only their last four characters. A replace is necessary when the
// configured key differs from the key in state, or when its last four characters differ from the API value. The
// second check catches a key rotated outside Terraform, and a key that does not match the API after an import.
func datadogKeyRequiresReplace(stateKey types.String, configKey types.String, stateLastFour types.String) bool {
	if configKey.IsUnknown() {
		return true
	}
	if !stateKey.IsNull() && !stateKey.IsUnknown() && stateKey.ValueString() != configKey.ValueString() {
		return true
	}
	if !stateLastFour.IsNull() && !stateLastFour.IsUnknown() && lastFour(configKey.ValueString()) != stateLastFour.ValueString() {
		return true
	}
	return false
}

func lastFour(value string) string {
	runes := []rune(value)
	if len(runes) <= 4 {
		return value
	}
	return string(runes[len(runes)-4:])
}
