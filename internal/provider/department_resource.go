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
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DepartmentResource{}
var _ resource.ResourceWithImportState = &DepartmentResource{}

func NewDepartmentResource() resource.Resource {
	return &DepartmentResource{}
}

func NewDepartmentResourceModel() DepartmentResourceModel {
	return DepartmentResourceModel{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// DepartmentResource defines the resource implementation.
type DepartmentResource struct {
	client *cortex.HttpClient
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *DepartmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Department Entity",

		Attributes: map[string]schema.Attribute{
			// Required attributes
			"tag": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the department.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the team.",
				Required:            true,
			},

			// Optional attributes
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team.",
				Optional:            true,
			},
			"members": schema.ListNestedAttribute{
				MarkdownDescription: "A list of additional members.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the member.",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A short description of the member.",
							Optional:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Email of the member.",
							Required:            true,
						},
					},
				},
			},

			// Computed attributes
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *DepartmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_department"
}

func (r *DepartmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DepartmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NewDepartmentResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	entity, err := r.client.Departments().Get(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read department %s, got error: %s", data.Tag.ValueString(), err))
		return
	}

	// Map entity to resource model
	data.FromApiModel(entity)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create Creates a new team.
func (r *DepartmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewDepartmentResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel()
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.Departments().Create(ctx, clientEntity.ToCreateRequest())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create department, got error: %s", err))
		return
	}

	// Map entity to resource model
	data.FromApiModel(entity)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DepartmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewDepartmentResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel()
	if resp.Diagnostics.HasError() {
		return
	}
	entity, err := r.client.Departments().Update(ctx, data.Tag.ValueString(), clientEntity.ToUpdateRequest())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update department, got error: %s", err))
		return
	}

	// Map entity to resource model
	data.FromApiModel(entity)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DepartmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewDepartmentResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Departments().Delete(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete department, got error: %s", err))
		return
	}
}

func (r *DepartmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
