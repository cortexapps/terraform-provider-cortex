package provider

import (
	"context"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CatalogEntityResource{}
var _ resource.ResourceWithImportState = &CatalogEntityResource{}

func NewCatalogEntityResource() resource.Resource {
	return &CatalogEntityResource{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// CatalogEntityResource defines the resource implementation.
type CatalogEntityResource struct {
	client *cortex.HttpClient
}

// CatalogEntityResourceModel describes the resource data model.
type CatalogEntityResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Tag         types.String `tfsdk:"tag"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *CatalogEntityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entity"
}

func (r *CatalogEntityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Catalog Entity",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name for the entity",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the entity visible in the Service or Resource Catalog. Markdown is supported.",
				Optional:            true,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the entity. Corresponds to the x-cortex-tag field in the entity descriptor.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			//Computed
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CatalogEntityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CatalogEntityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CatalogEntityResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	upsertRequest := cortex.UpsertCatalogEntityRequest{
		Info: cortex.CatalogEntityData{
			Tag:         data.Tag.ValueString(),
			Title:       data.Name.ValueString(),
			Description: data.Description.ValueString(),
		},
	}
	ceResponse, err := r.client.CatalogEntities().Upsert(ctx, upsertRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog entity, got error: %s", err))
		return
	}

	// Set computed attributes
	data.Id = types.StringValue(ceResponse.Tag)
	data.Tag = types.StringValue(ceResponse.Tag)
	// TODO: Add other attributes, consolidate this into a shared method

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CatalogEntityResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	entity, err := r.client.CatalogEntities().GetFromDescriptor(ctx, data.Tag.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	// TODO: Good practice it seems here is to issue a Read call which then stores into state?
	// https://github.com/hashicorp/terraform-provider-consul/blob/master/consul/resource_consul_acl_token.go#L222
	// https://github.com/hashicorp/terraform-provider-consul/blob/master/consul/resource_consul_acl_token.go#L144
	// This could be just a v1 provider sdk thing though

	// or maybe see https://github.com/hashicorp/terraform-provider-random/blob/main/internal/provider/resource_pet.go#L125
	// which seems to have a batch state set?

	// Set attributes from API response
	data.Id = types.StringValue(entity.Tag)
	data.Tag = types.StringValue(entity.Tag)
	data.Name = types.StringValue(entity.Title)
	data.Description = types.StringValue(entity.Description)
	// TODO: Add other attributes, consolidate this into a shared method

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CatalogEntityResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	upsertRequest := cortex.UpsertCatalogEntityRequest{
		Info: cortex.CatalogEntityData{
			Tag:         data.Tag.ValueString(),
			Title:       data.Name.ValueString(),
			Description: data.Description.ValueString(),
		},
	}
	entity, err := r.client.CatalogEntities().Upsert(ctx, upsertRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog entity, got error: %s", err))
		return
	}

	// Set computed attributes
	data.Id = data.Tag
	data.Tag = types.StringValue(entity.Tag)
	data.Name = types.StringValue(entity.Title)
	data.Description = types.StringValue(entity.Description)
	// TODO: Add other attributes, consolidate this into a shared method

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CatalogEntityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CatalogEntityResourceModel

	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CatalogEntities().Delete(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete catalog entity, got error: %s", err))
		return
	}
}

func (r *CatalogEntityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
