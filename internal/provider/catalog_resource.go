package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CatalogResource{}
var _ resource.ResourceWithImportState = &CatalogResource{}

func NewCatalogResource() resource.Resource {
	return &CatalogResource{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// CatalogResource defines the resource implementation.
type CatalogResource struct {
	client *cortex.HttpClient
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *CatalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog"
}

func (r *CatalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Catalog Entity",

		Attributes: map[string]schema.Attribute{
			// Required attributes
			"slug": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (slug) for the catalog.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the catalog.",
				Required:            true,
			},
			"icon_tag": schema.StringAttribute{
				MarkdownDescription: "Icon tag for the catalog.",
				Required:            true,
			},
			"is_draft": schema.BoolAttribute{
				MarkdownDescription: "Whether the catalog is a draft.",
				Required:            true,
			},

			// Optional attributes
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the catalog.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the catalog. One of: FILTER, RELATIONSHIP_TYPE, DOMAIN.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("FILTER", "RELATIONSHIP_TYPE", "DOMAIN"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"relationship_type_tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the relationship type associated with this catalog.",
				Optional:            true,
			},
			"filter": schema.SingleNestedAttribute{
				MarkdownDescription: "Filter configuration for the catalog.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"query": schema.StringAttribute{
						MarkdownDescription: "Query string for the filter.",
						Optional:            true,
					},
					"types": schema.SingleNestedAttribute{
						MarkdownDescription: "Entity type filter.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include": schema.SetAttribute{
								MarkdownDescription: "Set of entity types to include. Mutually exclusive with exclude.",
								Optional:            true,
								ElementType:         types.StringType,
								Validators: []validator.Set{
									setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("exclude")),
								},
							},
							"exclude": schema.SetAttribute{
								MarkdownDescription: "Set of entity types to exclude. Mutually exclusive with include.",
								Optional:            true,
								ElementType:         types.StringType,
								Validators: []validator.Set{
									setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("include")),
								},
							},
						},
					},
					"groups": schema.SingleNestedAttribute{
						MarkdownDescription: "Group filter.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include": schema.SetAttribute{
								MarkdownDescription: "Set of groups to include. Mutually exclusive with exclude.",
								Optional:            true,
								ElementType:         types.StringType,
								Validators: []validator.Set{
									setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("exclude")),
								},
							},
							"exclude": schema.SetAttribute{
								MarkdownDescription: "Set of groups to exclude. Mutually exclusive with include.",
								Optional:            true,
								ElementType:         types.StringType,
								Validators: []validator.Set{
									setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("include")),
								},
							},
						},
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
			"is_cortex_managed": schema.BoolAttribute{
				MarkdownDescription: "Whether this catalog is managed by Cortex.",
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

func (r *CatalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CatalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NewCatalogResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	entity, err := r.client.Catalogs().Get(ctx, data.Slug.ValueString())
	if err != nil {
		// If the catalog was deleted outside of Terraform, remove it from state
		// so that a subsequent plan recreates it instead of erroring.
		if errors.Is(err, cortex.ApiErrorNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog, got error: %s", err))
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

// Create creates a new catalog.
func (r *CatalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewCatalogResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.Catalogs().Upsert(ctx, clientEntity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create catalog, got error: %s", err))
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

func (r *CatalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewCatalogResourceModel()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.Catalogs().Upsert(ctx, clientEntity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update catalog, got error: %s", err))
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

func (r *CatalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewCatalogResourceModel()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Cortex-managed catalogs cannot be deleted via the API.
	if data.IsCortexManaged.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot Delete Cortex-Managed Catalog",
			"Cannot delete a Cortex-managed catalog",
		)
		return
	}

	err := r.client.Catalogs().Delete(ctx, data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete catalog, got error: %s", err))
		return
	}
}

func (r *CatalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	entity, err := r.client.Catalogs().Get(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to read catalog %q: %s", req.ID, err))
		return
	}
	if entity.IsCortexManaged {
		resp.Diagnostics.AddError(
			"Cannot Import Cortex-Managed Catalog",
			fmt.Sprintf("Catalog %q is managed by Cortex and cannot be imported into Terraform state. "+
				"Cortex-managed catalogs cannot be updated or deleted via the API.", req.ID),
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, resp)
}
