package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RelationshipTypeResource{}
var _ resource.ResourceWithImportState = &RelationshipTypeResource{}

func NewRelationshipTypeResource() resource.Resource {
	return &RelationshipTypeResource{}
}

func NewRelationshipTypeResourceModel() RelationshipTypeResourceModel {
	return RelationshipTypeResourceModel{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// RelationshipTypeResource defines the resource implementation.
type RelationshipTypeResource struct {
	client *cortex.HttpClient
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *RelationshipTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_type"
}

func (r *RelationshipTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	filterNestedAttributes := map[string]schema.Attribute{
		"include": schema.BoolAttribute{
			MarkdownDescription: "Whether to include or exclude the listed types/providers.",
			Required:            true,
		},
		"types": schema.SetAttribute{
			MarkdownDescription: "Entity types to filter on.",
			ElementType:         types.StringType,
			Optional:            true,
		},
		"providers": schema.SetAttribute{
			MarkdownDescription: "Entity providers to filter on.",
			ElementType:         types.StringType,
			Optional:            true,
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Relationship Type resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			// Required attributes
			"tag": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (tag) for the relationship type.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name for the relationship type.",
				Required:            true,
			},
			"allow_cycles": schema.BoolAttribute{
				MarkdownDescription: "Whether cycles are allowed in this relationship type.",
				Required:            true,
			},
			"is_single_source": schema.BoolAttribute{
				MarkdownDescription: "Whether each destination can only have one source.",
				Required:            true,
			},
			"is_single_destination": schema.BoolAttribute{
				MarkdownDescription: "Whether each source can only have one destination.",
				Required:            true,
			},
			"definition_location": schema.StringAttribute{
				MarkdownDescription: "Where the relationship is defined (e.g. SOURCE or DESTINATION).",
				Required:            true,
			},
			"sources_filter": schema.SingleNestedAttribute{
				MarkdownDescription: "Filter for which entity types can be sources.",
				Required:            true,
				Attributes:          filterNestedAttributes,
			},
			"destinations_filter": schema.SingleNestedAttribute{
				MarkdownDescription: "Filter for which entity types can be destinations.",
				Required:            true,
				Attributes:          filterNestedAttributes,
			},

			// Optional attributes
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the relationship type.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"source_label_singular": schema.StringAttribute{
				MarkdownDescription: "Singular label for the source entity.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"source_label_plural": schema.StringAttribute{
				MarkdownDescription: "Plural label for the source entity.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"destination_label_singular": schema.StringAttribute{
				MarkdownDescription: "Singular label for the destination entity.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"destination_label_plural": schema.StringAttribute{
				MarkdownDescription: "Plural label for the destination entity.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"inheritances": schema.SetNestedAttribute{
				MarkdownDescription: "Parameter inheritance configurations for this relationship type. " +
					"To clear all inheritances, remove this block from your config. " +
					"Setting `inheritances = []` is not supported and will cause a plan diff.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"parameter_tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the parameter to inherit.",
							Required:            true,
						},
						"inheritance_type": schema.StringAttribute{
							MarkdownDescription: "Type of inheritance (e.g. FALLBACK).",
							Required:            true,
						},
					},
				},
			},

			"create_catalog": schema.BoolAttribute{
				MarkdownDescription: "Whether to generate a catalog page for this relationship type.",
				Optional:            true,
				Computed:            true,
			},

			// Computed attributes
			"is_cortex_managed": schema.BoolAttribute{
				MarkdownDescription: "Whether this relationship type is managed by Cortex.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *RelationshipTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RelationshipTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NewRelationshipTypeResourceModel()

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.RelationshipTypes().Get(ctx, data.Tag.ValueString())
	if err != nil {
		if errors.Is(err, cortex.ApiErrorNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read relationship type, got error: %s", err))
		return
	}

	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RelationshipTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewRelationshipTypeResourceModel()

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := data.ToCreateRequest()
	entity, err := r.client.RelationshipTypes().Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create relationship type, got error: %s", err))
		return
	}

	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RelationshipTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewRelationshipTypeResourceModel()

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := data.ToUpdateRequest()
	entity, err := r.client.RelationshipTypes().Update(ctx, data.Tag.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update relationship type, got error: %s", err))
		return
	}

	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RelationshipTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewRelationshipTypeResourceModel()

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RelationshipTypes().Delete(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete relationship type, got error: %s", err))
		return
	}
}

func (r *RelationshipTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	entity, err := r.client.RelationshipTypes().Get(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read relationship type for import, got error: %s", err))
		return
	}

	if entity.IsCortexManaged {
		resp.Diagnostics.AddError(
			"Cannot Import Cortex-Managed Relationship Type",
			fmt.Sprintf("%s is managed by Cortex and cannot be updated or deleted via the API. Use terraform state rm if needed.", req.ID),
		)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
