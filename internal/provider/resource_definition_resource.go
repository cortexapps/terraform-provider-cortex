package provider

import (
	"context"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceDefinitionResource{}
var _ resource.ResourceWithImportState = &ResourceDefinitionResource{}

func NewResourceDefinitionResource() resource.Resource {
	return &ResourceDefinitionResource{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// ResourceDefinitionResource defines the resource implementation.
type ResourceDefinitionResource struct {
	client *cortex.HttpClient
}

// ResourceDefinitionResourceModel describes the department data model within Terraform.
type ResourceDefinitionResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Source      types.String `tfsdk:"source"`
	Schema      types.Map    `tfsdk:"schema"`
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *ResourceDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ResourceDefinition Entity",

		Attributes: map[string]schema.Attribute{
			// Required attributes
			"type": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the resource definition.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the team.",
				Required:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Source of the resource definition.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("cortex", "custom"),
				},
			},
			"schema": schema.MapAttribute{
				MarkdownDescription: "Schema for the resource definition.",
				ElementType:         types.StringType,
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

func (r *ResourceDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_definition"
}

func (r *ResourceDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceDefinitionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	response, err := r.client.ResourceDefinitions().Get(ctx, data.Type.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read resource definition, got error: %s", err))
		return
	}

	// Map data from the API response to the model
	data.Id = types.StringValue(response.Type)
	data.Type = types.StringValue(response.Type)
	data.Name = types.StringValue(response.Name)
	data.Description = types.StringValue(response.Description)
	data.Source = types.StringValue(response.Source)

	// TODO: Ensure this will actually work - unsure if this properly maintains schema types internally
	mv, diag := types.MapValueFrom(ctx, types.StringType, response.Schema)
	if diag.HasError() {
		resp.Diagnostics = diag
		return
	}
	data.Schema = mv

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create Creates a new team.
func (r *ResourceDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ResourceDefinitionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schemaReq := map[string]interface{}{}
	data.Schema.ElementsAs(ctx, schemaReq, true)

	clientRequest := cortex.CreateResourceDefinitionRequest{
		Type:        data.Type.String(),
		Name:        data.Name.String(),
		Description: data.Description.String(),
		Source:      data.Source.String(),
		Schema:      schemaReq,
	}
	resourceDefinition, err := r.client.ResourceDefinitions().Create(ctx, clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource definition, got error: %s", err))
		return
	}

	// Set the ID in state based on the Type
	data.Id = data.Type
	data.Type = types.StringValue(resourceDefinition.Type)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	//resp.State.SetAttribute(ctx, path.Root("id"), data.TeamTag)
	//resp.State.SetAttribute(ctx, path.Root("team_tag"), data.TeamTag)
}

func (r *ResourceDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ResourceDefinitionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schemaReq := map[string]interface{}{}
	data.Schema.ElementsAs(ctx, schemaReq, true)

	clientRequest := cortex.UpdateResourceDefinitionRequest{
		Name:        data.Name.String(),
		Description: data.Description.String(),
		Schema:      schemaReq,
	}
	resourceDefinition, err := r.client.ResourceDefinitions().Update(ctx, data.Type.String(), clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update resource definition, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = data.Type
	data.Type = types.StringValue(resourceDefinition.Type)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceDefinitionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ResourceDefinitions().Delete(ctx, data.Type.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete resource definition, got error: %s", err))
		return
	}
}

func (r *ResourceDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
