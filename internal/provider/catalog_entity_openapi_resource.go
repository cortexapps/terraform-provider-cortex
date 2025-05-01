package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CatalogEntityOpenAPIResource{}
	_ resource.ResourceWithConfigure   = &CatalogEntityOpenAPIResource{}
	_ resource.ResourceWithImportState = &CatalogEntityOpenAPIResource{}
)

type CatalogEntityOpenAPIResource struct {
	client *cortex.HttpClient
}

type CatalogEntityOpenAPIResourceModel struct {
	EntityTag types.String `tfsdk:"entity_tag"`
	Spec      types.String `tfsdk:"spec"`
	Id        types.String `tfsdk:"id"`
}

func NewCatalogEntityOpenAPIResource() resource.Resource {
	return &CatalogEntityOpenAPIResource{}
}

func (r *CatalogEntityOpenAPIResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cortex.HttpClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cortex.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *CatalogEntityOpenAPIResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entity_openapi"
}

func (r *CatalogEntityOpenAPIResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages OpenAPI specifications for Cortex catalog entities.",
		Attributes: map[string]schema.Attribute{
			"entity_tag": schema.StringAttribute{
				Description: "The tag or ID of the catalog entity that the OpenAPI specification will be associated with.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"spec": schema.StringAttribute{
				Description: "The OpenAPI specification in YAML or JSON format.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the OpenAPI specification.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CatalogEntityOpenAPIResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CatalogEntityOpenAPIResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	openAPISpec, err := r.client.CatalogEntityOpenAPI().Upsert(ctx, plan.EntityTag.ValueString(), cortex.UpsertCatalogEntityOpenAPIRequest{
		Spec: plan.Spec.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating OpenAPI specification",
			fmt.Sprintf("Could not create OpenAPI specification: %s", err.Error()),
		)
		return
	}

	plan.Id = types.StringValue(openAPISpec.ID())
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CatalogEntityOpenAPIResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CatalogEntityOpenAPIResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	openAPISpec, err := r.client.CatalogEntityOpenAPI().Get(ctx, state.EntityTag.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading OpenAPI specification",
			fmt.Sprintf("Could not read OpenAPI specification: %s", err.Error()),
		)
		return
	}

	state.Spec = types.StringValue(openAPISpec.Spec)
	state.Id = types.StringValue(openAPISpec.ID())

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *CatalogEntityOpenAPIResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CatalogEntityOpenAPIResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	openAPISpec, err := r.client.CatalogEntityOpenAPI().Upsert(ctx, plan.EntityTag.ValueString(), cortex.UpsertCatalogEntityOpenAPIRequest{
		Spec: plan.Spec.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating OpenAPI specification",
			fmt.Sprintf("Could not update OpenAPI specification: %s", err.Error()),
		)
		return
	}

	plan.Id = types.StringValue(openAPISpec.ID())
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CatalogEntityOpenAPIResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CatalogEntityOpenAPIResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CatalogEntityOpenAPI().Delete(ctx, state.EntityTag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting OpenAPI specification",
			fmt.Sprintf("Could not delete OpenAPI specification: %s", err.Error()),
		)
		return
	}
}

func (r *CatalogEntityOpenAPIResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("entity_tag"), req, resp)
}
