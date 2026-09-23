package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatadogConfigurationResource{}
var _ resource.ResourceWithImportState = &DatadogConfigurationResource{}

func NewDatadogConfigurationResource() resource.Resource {
	return &DatadogConfigurationResource{}
}

func NewDatadogConfigurationResourceModel() DatadogConfigurationResourceModel {
	return DatadogConfigurationResourceModel{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// DatadogConfigurationResource defines the resource implementation.
type DatadogConfigurationResource struct {
	client *cortex.HttpClient
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *DatadogConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Datadog integration configuration. Each configuration connects Cortex to one Datadog " +
			"organization and is identified by its alias.\n\n" +
			"The Cortex API cannot change `api_key`, `app_key`, `region`, or `custom_subdomain` in place, so a change to " +
			"one of them destroys the configuration and creates a new one. Cortex does not allow deleting the default " +
			"configuration while other configurations exist, so to replace the default configuration, first make a " +
			"different configuration the default.\n\n" +
			"Set `is_default = true` on one configuration only. Cortex keeps exactly one default, so two configurations " +
			"with `is_default = true` take the default from each other on every apply.",

		Attributes: map[string]schema.Attribute{
			// Required attributes
			"alias": schema.StringAttribute{
				MarkdownDescription: "Unique alias of the configuration. A change renames the configuration in place.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Datadog API key. The Cortex API never returns it, so Terraform detects a change " +
					"outside Terraform only through `api_key_last_four`.",
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					datadogKeyReplaceModifier("api_key_last_four"),
				},
			},
			"app_key": schema.StringAttribute{
				MarkdownDescription: "Datadog application key. The Cortex API never returns it, so Terraform detects a " +
					"change outside Terraform only through `app_key_last_four`.",
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					datadogKeyReplaceModifier("app_key_last_four"),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Datadog region (site) of the organization. One of `US1`, `US3`, `US5`, `EU1`, `US1_FED`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("US1", "US3", "US5", "EU1", "US1_FED"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Optional attributes
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "Whether this is the default Datadog configuration. When not set, Terraform keeps " +
					"the value from Cortex. Cortex always makes the first configuration the default, so `is_default = false` " +
					"fails on the first configuration. Cortex does not allow setting the current default to `false`: set " +
					"`is_default = true` on another configuration and remove `is_default` from this one, apply, and then set " +
					"it to `false` if necessary.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					datadogDefaultModifier{},
				},
			},
			"environments": schema.ListAttribute{
				MarkdownDescription: "Datadog environments to use for this configuration. Defaults to an empty list.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"custom_subdomain": schema.StringAttribute{
				MarkdownDescription: "Custom subdomain of the Datadog organization, if it uses one.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Computed attributes
			"id": schema.StringAttribute{
				MarkdownDescription: "Same as `alias`.",
				Computed:            true,
			},
			"api_key_last_four": schema.StringAttribute{
				MarkdownDescription: "Last four characters of the API key that Cortex stores.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_key_last_four": schema.StringAttribute{
				MarkdownDescription: "Last four characters of the application key that Cortex stores.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// datadogKeyReplaceModifier replaces the configuration when a key changes, because the API cannot update keys in
// place. After an import the key is null in state; a configured key whose last four characters match the API value
// then updates state in place instead of forcing a replace.
func datadogKeyReplaceModifier(lastFourAttribute string) planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			var stateLastFour types.String
			resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(lastFourAttribute), &stateLastFour)...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.RequiresReplace = datadogKeyRequiresReplace(req.StateValue, req.PlanValue, stateLastFour)
		},
		"A change to the key replaces the configuration.",
		"A change to the key replaces the configuration.",
	)
}

// datadogDefaultModifier fails the plan when the configuration sets is_default to false on the current default
// configuration, because the API rejects that update.
type datadogDefaultModifier struct{}

func (m datadogDefaultModifier) Description(ctx context.Context) string {
	return "Rejects is_default = false on the current default configuration."
}

func (m datadogDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m datadogDefaultModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if datadogUnsetsDefault(req.StateValue, req.ConfigValue) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Cannot unset the default Datadog configuration",
			"Cortex does not allow setting the current default configuration to is_default = false. Set is_default = true "+
				"on another configuration and remove is_default from this one, apply, and then set it to false if necessary.",
		)
	}
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *DatadogConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datadog_configuration"
}

func (r *DatadogConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cortex.HttpClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DatadogConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NewDatadogConfigurationResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	entity, err := r.client.DatadogConfigurations().Get(ctx, data.Alias.ValueString())
	if err != nil {
		if errors.Is(err, cortex.ApiErrorNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read datadog configuration %s, got error: %s", data.Alias.ValueString(), err))
		return
	}

	// Map entity to resource model
	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	data.ClearDriftedKeys()

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatadogConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewDatadogConfigurationResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := data.ToCreateRequest(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Cortex makes a new configuration the default when no default exists, whatever the request asks. Fail before
	// the create, so an explicit is_default = false does not leave a configuration that contradicts the plan.
	if !data.IsDefault.IsNull() && !data.IsDefault.IsUnknown() && !data.IsDefault.ValueBool() {
		existing, err := r.client.DatadogConfigurations().List(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list datadog configurations, got error: %s", err))
			return
		}
		if datadogCreateBecomesDefault(existing) {
			resp.Diagnostics.AddAttributeError(
				path.Root("is_default"),
				"Datadog configuration must be the default",
				fmt.Sprintf("No datadog configuration is the default, so Cortex makes %s the default. Remove "+
					"is_default or set it to true.", data.Alias.ValueString()),
			)
			return
		}
	}

	entity, err := r.client.DatadogConfigurations().Create(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create datadog configuration, got error: %s", err))
		return
	}

	// The create request never asks for the default, so an update sets it. The API unsets the previous default.
	wantDefault := !data.IsDefault.IsNull() && !data.IsDefault.IsUnknown() && data.IsDefault.ValueBool()
	if wantDefault && !entity.IsDefault {
		updateRequest := data.ToUpdateRequest(ctx, &resp.Diagnostics)
		updated, err := r.client.DatadogConfigurations().Update(ctx, entity.Alias, updateRequest)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Created datadog configuration %s, but unable to make it the default, got error: %s", entity.Alias, err))
		} else {
			entity = updated
		}
	}

	// Map entity to resource model. Save state also after an error, so Terraform keeps track of the configuration.
	data.FromApiModel(ctx, &resp.Diagnostics, entity)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatadogConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewDatadogConfigurationResourceModel()
	state := NewDatadogConfigurationResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	// Read the prior state for the current alias, because the plan can rename the configuration.
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := data.ToUpdateRequest(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.DatadogConfigurations().Update(ctx, state.Alias.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update datadog configuration %s, got error: %s", state.Alias.ValueString(), err))
		return
	}

	// Map entity to resource model
	data.FromApiModel(ctx, &resp.Diagnostics, entity)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatadogConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewDatadogConfigurationResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DatadogConfigurations().Delete(ctx, data.Alias.ValueString())
	if err != nil && !errors.Is(err, cortex.ApiErrorNotFound) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete datadog configuration %s, got error: %s", data.Alias.ValueString(), err))
		return
	}
}

func (r *DatadogConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("alias"), req, resp)
}
