package provider

import (
	"context"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CatalogEntityResource{}
var _ resource.ResourceWithImportState = &CatalogEntityResource{}

func NewCatalogEntityResource() resource.Resource {
	return &CatalogEntityResource{}
}

// CatalogEntityResource defines the resource implementation.
type CatalogEntityResource struct {
	client *cortex.Client
}

// CatalogEntityResourceModel describes the resource data model.
type CatalogEntityResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Tag         types.String `tfsdk:"tag"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

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
				//PlanModifiers: []planmodifier.String{
				//	stringplanmodifier.UseStateForUnknown(),
				//},
			},

			//Computed
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *CatalogEntityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cortex.Client)

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// Sample data for now
	data.Tag = types.StringValue("test")
	data.Name = types.StringValue("A Test Service")
	data.Description = types.StringValue("A test service for the Terraform provider")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a catalog entity")

	// Set the ID in state based on the tag
	//data.Id = data.Tag

	// Save data into Terraform state
	//resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.SetAttribute(ctx, path.Root("id"), data.Tag)
	resp.State.SetAttribute(ctx, path.Root("tag"), data.Tag)
	resp.State.SetAttribute(ctx, path.Root("name"), data.Name)
	resp.State.SetAttribute(ctx, path.Root("description"), data.Description)
}

func (r *CatalogEntityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CatalogEntityResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	//resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.SetAttribute(ctx, path.Root("id"), data.Tag)
	resp.State.SetAttribute(ctx, path.Root("tag"), data.Tag)
	resp.State.SetAttribute(ctx, path.Root("name"), data.Name)
	resp.State.SetAttribute(ctx, path.Root("description"), data.Description)
}

func (r *CatalogEntityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CatalogEntityResourceModel

	// Read Terraform plan data into the model
	//resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	//resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.SetAttribute(ctx, path.Root("id"), data.Tag)
	resp.State.SetAttribute(ctx, path.Root("tag"), data.Tag)
	resp.State.SetAttribute(ctx, path.Root("name"), data.Name)
	resp.State.SetAttribute(ctx, path.Root("description"), data.Description)
}

func (r *CatalogEntityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data *CatalogEntityResourceModel

	// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *CatalogEntityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
	resp.State.SetAttribute(ctx, path.Root("name"), "A Test Service")
	resp.State.SetAttribute(ctx, path.Root("description"), "A test service for the Terraform provider")
}
