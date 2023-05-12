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
var _ resource.Resource = &ScorecardResource{}
var _ resource.ResourceWithImportState = &ScorecardResource{}

func NewScorecardResource() resource.Resource {
	return &ScorecardResource{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// ScorecardResource defines the resource implementation.
type ScorecardResource struct {
	client *cortex.HttpClient
}

// ScorecardResourceModel describes the scorecard data model within Terraform.
type ScorecardResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Tag         types.String `tfsdk:"tag"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsDraft     types.Bool   `tfsdk:"is_draft"`
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *ScorecardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Scorecard Entity",

		Attributes: map[string]schema.Attribute{
			// Required attributes
			"tag": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the scorecard.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the scorecard.",
				Required:            true,
			},

			// Optional attributes
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the scorecard.",
				Optional:            true,
			},
			"is_draft": schema.BoolAttribute{
				MarkdownDescription: "Whether the scorecard is a draft.",
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
		Blocks: map[string]schema.Block{
			"level": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the link.",
							Required:            true,
						},
						"number": schema.Int64Attribute{
							MarkdownDescription: "Rank of the level. 1 is the highest.",
							Required:            true,
						},
					},
				},
			},
			"rule": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"title": schema.StringAttribute{
							MarkdownDescription: "Title of the rule.",
							Required:            true,
						},
						"expression": schema.StringAttribute{
							MarkdownDescription: "Expression of the rule.",
							Required:            true,
						},
						"number": schema.Int64Attribute{
							MarkdownDescription: "Rank of the level. 1 is the highest.",
							Required:            true,
						},
						"weight": schema.Int64Attribute{
							MarkdownDescription: "Numerical weight of the rule. When using levels, this defaults to 1.",
							Required:            true,
						},

						// Optional attributes
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the rule.",
							Optional:            true,
						},
						"level_name": schema.StringAttribute{
							MarkdownDescription: "Name of the level this rule is associated with, if applicable.",
							Optional:            true,
						},
						"failure_message": schema.StringAttribute{
							MarkdownDescription: "Message to display when the rule fails.",
							Optional:            true,
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

func (r *ScorecardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard"
}

func (r *ScorecardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ScorecardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ScorecardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	scorecardResponse, err := r.client.Scorecards().Get(ctx, data.Tag.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scorecard, got error: %s", err))
		return
	}

	// Map data from the API response to the model
	data.Id = types.StringValue(scorecardResponse.Tag)
	data.Tag = types.StringValue(scorecardResponse.Tag)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create Creates a new scorecard.
func (r *ScorecardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ScorecardResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientRequest := cortex.UpsertScorecardRequest{
		Tag:         data.Tag.String(),
		Name:        data.Name.String(),
		Description: data.Description.String(),
		IsDraft:     data.IsDraft.ValueBool(),
	}
	scorecard, err := r.client.Scorecards().Upsert(ctx, clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create scorecard, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = types.StringValue(scorecard.Tag)
	data.Tag = types.StringValue(scorecard.Tag)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScorecardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ScorecardResourceModel

	// Read Terraform plan data into the model
	//resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientRequest := cortex.UpsertScorecardRequest{
		Tag:         data.Tag.String(),
		Name:        data.Name.String(),
		Description: data.Description.String(),
		IsDraft:     data.IsDraft.ValueBool(),
	}
	scorecard, err := r.client.Scorecards().Upsert(ctx, clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update scorecard, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = types.StringValue(scorecard.Tag)
	data.Tag = types.StringValue(scorecard.Tag)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScorecardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ScorecardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Scorecards().Delete(ctx, data.Tag.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete scorecard, got error: %s", err))
		return
	}
}

func (r *ScorecardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
