package provider

import (
	"context"
	"fmt"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CatalogEntitiesDataSource{}

func NewCatalogEntitiesDataSource() datasource.DataSource {
	return &CatalogEntitiesDataSource{}
}

// CatalogEntitiesDataSource defines the data source implementation.
type CatalogEntitiesDataSource struct {
	client *cortex.HttpClient
}

// CatalogEntitiesDataSourceModel describes the data source data model.
type CatalogEntitiesDataSourceModel struct {
	Id              types.String                       `tfsdk:"id"`
	Query           types.String                       `tfsdk:"query"`
	Groups          []types.String                     `tfsdk:"groups"`
	Owners          []types.String                     `tfsdk:"owners"`
	Types           []types.String                     `tfsdk:"types"`
	GitRepositories []types.String                     `tfsdk:"git_repositories"`
	IncludeArchived types.Bool                         `tfsdk:"include_archived"`
	IncludeOwners   types.Bool                         `tfsdk:"include_owners"`
	Entities        []CatalogEntityDataSourceItemModel `tfsdk:"entities"`
}

// CatalogEntityDataSourceItemModel represents a single entity in the list.
type CatalogEntityDataSourceItemModel struct {
	Tag         types.String                          `tfsdk:"tag"`
	Name        types.String                          `tfsdk:"name"`
	Description types.String                          `tfsdk:"description"`
	Type        types.String                          `tfsdk:"type"`
	Git         *CatalogEntityGitItemModel            `tfsdk:"git"`
	Ownership   *CatalogEntityOwnershipItemModel      `tfsdk:"ownership"`
}

// CatalogEntityGitItemModel holds the flat git information returned by the catalog list endpoint.
type CatalogEntityGitItemModel struct {
	Provider      types.String `tfsdk:"provider"`
	Repository    types.String `tfsdk:"repository"`
	RepositoryUrl types.String `tfsdk:"repository_url"`
}

// CatalogEntityOwnershipItemModel holds ownership data for a catalog entity item.
type CatalogEntityOwnershipItemModel struct {
	Groups      []CatalogEntityOwnershipGroupItemModel      `tfsdk:"groups"`
	Individuals []CatalogEntityOwnershipIndividualItemModel `tfsdk:"individuals"`
}

// CatalogEntityOwnershipGroupItemModel represents a single owning team for a catalog entity.
type CatalogEntityOwnershipGroupItemModel struct {
	GroupName   types.String `tfsdk:"group_name"`
	Description types.String `tfsdk:"description"`
	Provider    types.String `tfsdk:"provider"`
}

// CatalogEntityOwnershipIndividualItemModel represents a single individual owner for a catalog entity.
type CatalogEntityOwnershipIndividualItemModel struct {
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
}

func (d *CatalogEntitiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entities"
}

func (d *CatalogEntitiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Catalog Entities data source - returns a list of catalog entities that match the given search criteria",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier for this data source",
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "Filter based on a search query. This will search across entity properties. If provided, results will be sorted by relevance.",
				Optional:            true,
			},
			"groups": schema.ListAttribute{
				MarkdownDescription: "Filter based on groups, which correspond to the x-cortex-groups field in the Catalog Descriptor",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"owners": schema.ListAttribute{
				MarkdownDescription: "Filter based on owner group names, which correspond to the x-cortex-owners field in the Catalog Descriptor",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"types": schema.ListAttribute{
				MarkdownDescription: "Filter the response to specific types of entities (e.g., service, resource, domain)",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"git_repositories": schema.ListAttribute{
				MarkdownDescription: "Filter by GitHub repositories in the 'org/repo' format",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"include_archived": schema.BoolAttribute{
				MarkdownDescription: "Whether to include archived entities in the response",
				Optional:            true,
			},
			"include_owners": schema.BoolAttribute{
				MarkdownDescription: "When true, each entity in the response will include ownership information (teams and individuals). Corresponds to the `includeOwners` API parameter.",
				Optional:            true,
			},
			"entities": schema.ListNestedAttribute{
				MarkdownDescription: "List of catalog entities that match the search criteria",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the catalog entity",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Human-readable name for the entity",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the entity visible in the Service or Resource Catalog",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the entity (e.g., service, resource, domain)",
							Computed:            true,
						},
						"git": schema.SingleNestedAttribute{
							MarkdownDescription: "Git repository information for the entity. Populated when a git integration is configured.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"provider": schema.StringAttribute{
									MarkdownDescription: "Git provider (e.g., GITHUB, GITLAB)",
									Computed:            true,
								},
								"repository": schema.StringAttribute{
									MarkdownDescription: "Repository name in org/repo format",
									Computed:            true,
								},
								"repository_url": schema.StringAttribute{
									MarkdownDescription: "Full URL to the repository",
									Computed:            true,
								},
							},
						},
						"ownership": schema.SingleNestedAttribute{
							MarkdownDescription: "Ownership information for the entity. Populated when `include_owners` is true.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"groups": schema.ListNestedAttribute{
									MarkdownDescription: "List of owning teams for the entity",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"group_name": schema.StringAttribute{
												MarkdownDescription: "Tag of the owning team",
												Computed:            true,
											},
											"description": schema.StringAttribute{
												MarkdownDescription: "Description of this ownership entry",
												Computed:            true,
											},
											"provider": schema.StringAttribute{
												MarkdownDescription: "Provider for this group (e.g., GITHUB, OKTA)",
												Computed:            true,
											},
										},
									},
								},
								"individuals": schema.ListNestedAttribute{
									MarkdownDescription: "List of individual owners for the entity",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"email": schema.StringAttribute{
												MarkdownDescription: "Email address of the individual owner",
												Computed:            true,
											},
											"description": schema.StringAttribute{
												MarkdownDescription: "Description of this ownership entry",
												Computed:            true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *CatalogEntitiesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cortex.HttpClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cortex.HttpClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *CatalogEntitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CatalogEntitiesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build API parameters
	params := &cortex.CatalogEntityListParams{}

	if !data.Query.IsNull() && !data.Query.IsUnknown() {
		params.Query = data.Query.ValueString()
	}

	if len(data.Groups) > 0 {
		params.Groups = make([]string, len(data.Groups))
		for i, g := range data.Groups {
			params.Groups[i] = g.ValueString()
		}
	}

	if len(data.Owners) > 0 {
		params.Owners = make([]string, len(data.Owners))
		for i, o := range data.Owners {
			params.Owners[i] = o.ValueString()
		}
	}

	if len(data.Types) > 0 {
		params.Types = make([]string, len(data.Types))
		for i, t := range data.Types {
			params.Types[i] = t.ValueString()
		}
	}

	if len(data.GitRepositories) > 0 {
		params.GitRepositories = make([]string, len(data.GitRepositories))
		for i, gr := range data.GitRepositories {
			params.GitRepositories[i] = gr.ValueString()
		}
	}

	if !data.IncludeArchived.IsNull() && !data.IncludeArchived.IsUnknown() {
		params.IncludeArchived = data.IncludeArchived.ValueBool()
	}

	if !data.IncludeOwners.IsNull() && !data.IncludeOwners.IsUnknown() {
		params.IncludeOwners = data.IncludeOwners.ValueBool()
	}

	// Fetch all pages of results
	allEntities := []cortex.CatalogEntity{}
	page := 0
	pageSize := 250 // Default page size from API

	params.PageSize = pageSize
	params.Page = page

	for {
		// Issue API request
		entitiesResponse, err := d.client.CatalogEntities().List(ctx, params)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog entities, got error: %s", err))
			return
		}

		// Accumulate entities from this page
		allEntities = append(allEntities, entitiesResponse.Entities...)

		// Check if there are more pages
		if entitiesResponse.Page >= entitiesResponse.TotalPages-1 || len(entitiesResponse.Entities) == 0 {
			break
		}

		// Move to next page
		page++
		params.Page = page
	}

	// Map response to state
	data.Entities = make([]CatalogEntityDataSourceItemModel, len(allEntities))
	for i, entity := range allEntities {
		item := CatalogEntityDataSourceItemModel{
			Tag:         types.StringValue(entity.Tag),
			Name:        types.StringValue(entity.Name),
			Description: types.StringValue(entity.Description),
			Type:        types.StringValue(entity.Type),
		}

		// Git is returned by default from the list endpoint when a git integration is configured.
		if entity.Git.Repository != "" {
			item.Git = &CatalogEntityGitItemModel{
				Provider:      types.StringValue(entity.Git.Provider),
				Repository:    types.StringValue(entity.Git.Repository),
				RepositoryUrl: types.StringValue(entity.Git.RepositoryUrl),
			}
		}

		// Ownership is populated when include_owners is true. The list endpoint returns
		// owners.teams[] and owners.individuals[] (not ownership.groups[]).
		if len(entity.Owners.Teams) > 0 || len(entity.Owners.Individuals) > 0 {
			groups := make([]CatalogEntityOwnershipGroupItemModel, len(entity.Owners.Teams))
			for j, t := range entity.Owners.Teams {
				groups[j] = CatalogEntityOwnershipGroupItemModel{
					GroupName:   types.StringValue(t.Tag),
					Description: types.StringValue(t.Description),
					Provider:    types.StringValue(t.Provider),
				}
			}
			individuals := make([]CatalogEntityOwnershipIndividualItemModel, len(entity.Owners.Individuals))
			for j, ind := range entity.Owners.Individuals {
				individuals[j] = CatalogEntityOwnershipIndividualItemModel{
					Email:       types.StringValue(ind.Email),
					Description: types.StringValue(ind.Description),
				}
			}
			item.Ownership = &CatalogEntityOwnershipItemModel{Groups: groups, Individuals: individuals}
		}

		data.Entities[i] = item
	}

	// Set ID as a hash of the query parameters
	data.Id = types.StringValue("catalog_entities")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
