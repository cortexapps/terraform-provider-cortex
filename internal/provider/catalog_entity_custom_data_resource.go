package provider

import (
	"context"
	"fmt"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CatalogEntityCustomDataResource{}
var _ resource.ResourceWithImportState = &CatalogEntityCustomDataResource{}

func NewCatalogEntityCustomDataResource() resource.Resource {
	return &CatalogEntityCustomDataResource{}
}

func NewCatalogEntityCustomDataResourceModel() CatalogEntityCustomDataResourceModel {
	return CatalogEntityCustomDataResourceModel{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

type CatalogEntityCustomDataResource struct {
	client *cortex.HttpClient
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *CatalogEntityCustomDataResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entity_custom_data"
}

func (r *CatalogEntityCustomDataResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Catalog Entity Custom Data",

		Attributes: map[string]schema.Attribute{
			// Required attributes
			"tag": schema.StringAttribute{
				MarkdownDescription: "The Catalog Entity tag for this custom data.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key of the custom data entry for the catalog entity.",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value for the custom data. Must be a JSON-encoded string.",
				Required:            true,
			},

			// Optional attributes
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team.",
				Optional:            true,
			},

			// Computed attributes
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *CatalogEntityCustomDataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CatalogEntityCustomDataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NewCatalogEntityCustomDataResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	entity, err := r.client.CatalogEntityCustomData().Get(ctx, data.Tag.ValueString(), data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog entity custom data, got error: %s", err))
		return
	}

	// Map data from the API response to the model
	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityCustomDataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewCatalogEntityCustomDataResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.CatalogEntityCustomData().Upsert(ctx, clientEntity.Tag, clientEntity.ToUpsertRequest())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create catalog entity custom data, got error: %s", err))
		return
	}

	// Set computed attributes
	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityCustomDataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewCatalogEntityCustomDataResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.CatalogEntityCustomData().Upsert(ctx, clientEntity.Tag, clientEntity.ToUpsertRequest())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update catalog entity custom data, got error: %s", err))
		return
	}

	// Set computed attributes
	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityCustomDataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewCatalogEntityCustomDataResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CatalogEntityCustomData().Delete(ctx, data.Tag.ValueString(), data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete catalog entity custom data, got error: %s", err))
		return
	}
}

func (r *CatalogEntityCustomDataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: tag:key. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tag"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), idParts[1])...)
}
