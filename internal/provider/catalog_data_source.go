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
var _ datasource.DataSource = &CatalogDataSource{}

func NewCatalogDataSource() datasource.DataSource {
	return &CatalogDataSource{}
}

// CatalogDataSource defines the data source implementation.
type CatalogDataSource struct {
	client *cortex.HttpClient
}

func (d *CatalogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog"
}

func (d *CatalogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Catalog data source",

		Attributes: map[string]schema.Attribute{
			// Required — lookup key
			"slug": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (slug) for the catalog.",
				Required:            true,
			},

			// Computed
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the catalog.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the catalog.",
				Computed:            true,
			},
			"icon_tag": schema.StringAttribute{
				MarkdownDescription: "Icon tag for the catalog.",
				Computed:            true,
			},
			"is_draft": schema.BoolAttribute{
				MarkdownDescription: "Whether the catalog is a draft.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the catalog. One of: FILTER, RELATIONSHIP_TYPE, DOMAIN.",
				Computed:            true,
			},
			"relationship_type_tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the relationship type associated with this catalog.",
				Computed:            true,
			},
			"filter": schema.SingleNestedAttribute{
				MarkdownDescription: "Filter configuration for the catalog.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"query": schema.StringAttribute{
						MarkdownDescription: "Query string for the filter.",
						Computed:            true,
					},
					"types": schema.SingleNestedAttribute{
						MarkdownDescription: "Entity type filter.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"include": schema.SetAttribute{
								MarkdownDescription: "Set of entity types to include.",
								Computed:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.SetAttribute{
								MarkdownDescription: "Set of entity types to exclude.",
								Computed:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"groups": schema.SingleNestedAttribute{
						MarkdownDescription: "Group filter.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"include": schema.SetAttribute{
								MarkdownDescription: "Set of groups to include.",
								Computed:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.SetAttribute{
								MarkdownDescription: "Set of groups to exclude.",
								Computed:            true,
								ElementType:         types.StringType,
							},
						},
					},
				},
			},
			"is_cortex_managed": schema.BoolAttribute{
				MarkdownDescription: "Whether this catalog is managed by Cortex.",
				Computed:            true,
			},
		},
	}
}

func (d *CatalogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cortex.HttpClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *CatalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := NewCatalogDataSourceModel()

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := d.client.Catalogs().Get(ctx, data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog, got error: %s", err))
		return
	}

	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write to TF state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
