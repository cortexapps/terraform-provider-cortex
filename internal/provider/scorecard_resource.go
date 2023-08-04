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
			"ladder": schema.SingleNestedAttribute{
				MarkdownDescription: "Ladder of the scorecard.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"levels": schema.ListNestedAttribute{
						MarkdownDescription: "Levels of the scorecard.",
						Required:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								// Required level attributes
								"name": schema.StringAttribute{
									MarkdownDescription: "Name of the level.",
									Required:            true,
								},
								"rank": schema.Int64Attribute{
									MarkdownDescription: "Rank of the level. 1 is the lowest.",
									Required:            true,
								},
								"color": schema.StringAttribute{
									MarkdownDescription: "Color of the level.",
									Required:            true,
								},

								// Optional level attributes
								"description": schema.StringAttribute{
									MarkdownDescription: "Description of the level.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"rules": schema.ListNestedAttribute{
				MarkdownDescription: "Rules of the scorecard.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// Required rule attributes
						"title": schema.StringAttribute{
							MarkdownDescription: "Title of the rule.",
							Required:            true,
						},
						"expression": schema.StringAttribute{
							MarkdownDescription: "Expression of the rule.",
							Required:            true,
						},
						"weight": schema.Int64Attribute{
							MarkdownDescription: "Weight of the rule.",
							Required:            true,
						},
						"level": schema.StringAttribute{
							MarkdownDescription: "Level of the rule for the ladder.",
							Required:            true,
						},

						// Optional rule attributes
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the rule.",
							Optional:            true,
						},
						"failure_message": schema.StringAttribute{
							MarkdownDescription: "Failure message of the rule.",
							Optional:            true,
						},
					},
				},
			},

			// Optional attributes
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the scorecard.",
				Optional:            true,
			},
			"draft": schema.BoolAttribute{
				MarkdownDescription: "Whether the scorecard is a draft.",
				Optional:            true,
			},
			"filter": schema.SingleNestedAttribute{
				MarkdownDescription: "Filter of the scorecard.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"category": schema.StringAttribute{
						MarkdownDescription: "By default, Scorecards are evaluated against all services. You can specify the category as RESOURCE to evaluate a Scorecard against resources or DOMAIN to evaluate a Scorecard against domains.",
						Optional:            true,
					},
					"query": schema.StringAttribute{
						MarkdownDescription: "A CQL query that is run against the category; only entities matching this query will be evaluated by the Scorecard.",
						Optional:            true,
					},
				},
			},
			"evaluation": schema.SingleNestedAttribute{
				MarkdownDescription: "Evaluation of the scorecard.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"window": schema.Int64Attribute{
						MarkdownDescription: "In hours. By default, Scorecards are evaluated every 4 hours. If you would like to evaluate Scorecards less frequently, you can override the evaluation window. This can help with rate limits. Note that Scorecards cannot be evaluated more than once per 4 hours.",
						Optional:            true,
					},
				},
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
	data := NewScorecardResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	entity, err := r.client.Scorecards().Get(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scorecard, got error: %s", err))
		return
	}

	// Map data from the API response to the model
	data.FromApiModel(ctx, &resp.Diagnostics, entity)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create Creates a new scorecard.
func (r *ScorecardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewScorecardResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel(ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	scorecard, err := r.client.Scorecards().Upsert(ctx, clientEntity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create scorecard, got error: %s", err))
		return
	}

	data.FromApiModel(ctx, &resp.Diagnostics, scorecard)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScorecardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewScorecardResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel(ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	scorecard, err := r.client.Scorecards().Upsert(ctx, clientEntity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update scorecard, got error: %s", err))
		return
	}

	data.FromApiModel(ctx, &resp.Diagnostics, scorecard)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScorecardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewScorecardResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Scorecards().Delete(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete scorecard, got error: %s", err))
		return
	}
}

func (r *ScorecardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
