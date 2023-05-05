package provider

import (
	"context"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DepartmentResource{}
var _ resource.ResourceWithImportState = &DepartmentResource{}

func NewDepartmentResource() resource.Resource {
	return &DepartmentResource{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// DepartmentResource defines the resource implementation.
type DepartmentResource struct {
	client *cortex.HttpClient
}

// DepartmentResourceModel describes the department data model within Terraform.
type DepartmentResourceModel struct {
	Id          types.String                    `tfsdk:"id"`
	Tag         types.String                    `tfsdk:"tag"`
	Name        types.String                    `tfsdk:"name"`
	Description types.String                    `tfsdk:"description"`
	Members     []DepartmentMemberResourceModel `tfsdk:"members"`
}

type DepartmentMemberResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description,omitempty"`
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
				//PlanModifiers: []planmodifier.String{
				//	stringplanmodifier.UseStateForUnknown(),
				//},
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

			// Computed attributes
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"member": schema.ListNestedBlock{
				MarkdownDescription: "A list of additional members.",
				NestedObject: schema.NestedBlockObject{
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
	var data *DepartmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	departmentResponse, err := r.client.Departments().Get(ctx, data.Tag.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read department, got error: %s", err))
		return
	}

	// Map data from the API response to the model
	data.Id = types.StringValue(departmentResponse.Tag)
	data.Tag = types.StringValue(departmentResponse.Tag)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create Creates a new team.
func (r *DepartmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DepartmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var members []cortex.DepartmentMember
	for _, member := range data.Members {
		members = append(members, cortex.DepartmentMember{
			Name:        member.Name.String(),
			Email:       member.Email.String(),
			Description: member.Description.String(),
		})
	}

	clientRequest := cortex.CreateDepartmentRequest{
		Tag:         data.Tag.String(),
		Name:        data.Name.String(),
		Description: data.Description.String(),
		Members:     members,
	}
	department, err := r.client.Departments().Create(ctx, clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create department, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = data.Tag
	data.Tag = types.StringValue(department.Tag)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	//resp.State.SetAttribute(ctx, path.Root("id"), data.TeamTag)
	//resp.State.SetAttribute(ctx, path.Root("team_tag"), data.TeamTag)
}

func (r *DepartmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DepartmentResourceModel

	// Read Terraform plan data into the model
	//resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var members []cortex.DepartmentMember
	for _, member := range data.Members {
		members = append(members, cortex.DepartmentMember{
			Name:        member.Name.String(),
			Email:       member.Email.String(),
			Description: member.Description.String(),
		})
	}

	clientRequest := cortex.UpdateDepartmentRequest{
		Name:        data.Name.String(),
		Description: data.Description.String(),
		Members:     members,
	}
	department, err := r.client.Departments().Update(ctx, data.Tag.String(), clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update department, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = data.Tag
	data.Tag = types.StringValue(department.Tag)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DepartmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DepartmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Departments().Delete(ctx, data.Tag.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete department, got error: %s", err))
		return
	}
}

func (r *DepartmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
