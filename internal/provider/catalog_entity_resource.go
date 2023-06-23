package provider

import (
	"context"
	"encoding/json"
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

func (r *CatalogEntityResource) toUpsertRequest(data *CatalogEntityResourceModel) cortex.UpsertCatalogEntityRequest {
	return cortex.UpsertCatalogEntityRequest{
		Info: data.ToApiModel(),
	}
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

			// Optional attributes
			"owners": schema.ListNestedAttribute{
				MarkdownDescription: "List of owners for the entity. Owners can be users, groups, or Slack channels.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of owner. Valid values are `EMAIL`, `GROUP`, `OKTA`, or `SLACK`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("EMAIL", "GROUP", "OKTA", "SLACK"),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the owner. Only required for `user` or `group` types.",
							Optional:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Email of the owner. Only allowed if `type` is `user`.",
							Optional:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the owner. Optional.",
							Optional:            true,
						},
						"provider": schema.StringAttribute{
							MarkdownDescription: "Provider of the owner. Only allowed if `type` is `group`.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ACTIVE_DIRECTORY", "BAMBOO_HR", "CORTEX", "GITHUB", "GOOGLE", "OKTA", "OPSGENIE", "WORKDAY"),
							},
						},
						"channel": schema.StringAttribute{
							MarkdownDescription: "Channel of the owner. Only allowed if `type` is `slack`. Omit the #.",
							Optional:            true,
						},
						"notifications_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether Slack notifications are enabled for all owners of this service. Only allowed if `type` is `slack`.",
							Optional:            true,
						},
					},
				},
			},
			"groups": schema.ListAttribute{
				MarkdownDescription: "List of groups related to the entity.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"links": schema.ListNestedAttribute{
				MarkdownDescription: "List of links related to the entity.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the link.",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the link. Valid values are `runbook`, `documentation`, `logs`, `dashboard`, `metrics`, `healthcheck`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("runbook", "documentation", "logs", "dashboard", "metrics", "healthcheck"),
							},
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "URL of the link.",
							Required:            true,
						},
					},
				},
			},
			"metadata": schema.StringAttribute{
				MarkdownDescription: "Custom metadata for the entity, in JSON format in a string. (Use the `jsonencode` function to convert a JSON object to a string.)",
				Optional:            true,
			},
			"dependencies": schema.ListNestedAttribute{
				MarkdownDescription: "List of dependencies for the entity.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the dependency.",
							Required:            true,
						},
						"method": schema.StringAttribute{
							MarkdownDescription: "HTTP method if depending on a specific endpoint.",
							Optional:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "The actual endpoint this dependency refers to.",
							Optional:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the dependency.",
							Optional:            true,
						},
						"metadata": schema.StringAttribute{
							MarkdownDescription: "Custom metadata for the dependency, in JSON format in a string. (Use the `jsonencode` function to convert a JSON object to a string.)",
							Optional:            true,
						},
					},
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
	upsertRequest := r.toUpsertRequest(data)
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

	// Set attributes from API response
	data.Id = types.StringValue(entity.Tag)
	data.Tag = types.StringValue(entity.Tag)

	// coerce map of unknown types into string
	if entity.Metadata != nil {
		metadata, err := json.Marshal(entity.Metadata)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metadata, got error: %s", err))
			return
		}
		data.Metadata = types.StringValue(string(metadata))
	}

	if data.Dependencies != nil {
		for _, dependency := range entity.Dependencies {
			depMetadata := []byte("")
			if dependency.Metadata != nil {
				depMetadata, err = json.Marshal(dependency.Metadata)
				if err != nil {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dependency metadata, got error: %s", err))
					return
				}
			}
			depResourceModel := CatalogEntityDependencyResourceModel{
				Tag:      types.StringValue(dependency.Tag),
				Method:   types.StringValue(dependency.Method),
				Path:     types.StringValue(dependency.Path),
				Metadata: types.StringValue(string(depMetadata)),
			}
			data.Dependencies = append(data.Dependencies, depResourceModel)
		}
	}

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
	upsertRequest := r.toUpsertRequest(data)
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
