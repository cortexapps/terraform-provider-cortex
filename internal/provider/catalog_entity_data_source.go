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
var _ datasource.DataSource = &CatalogEntityDataSource{}

func NewCatalogEntityDataSource() datasource.DataSource {
	return &CatalogEntityDataSource{}
}

// CatalogEntityDataSource defines the data source implementation.
type CatalogEntityDataSource struct {
	client *cortex.HttpClient
}

type CatalogEntityChildSourceModel struct {
	Tag types.String `tfsdk:"tag"`
}

// CatalogEntityDataSourceModel describes the data source data model.
type CatalogEntityDataSourceModel struct {
	Id          types.String                    `tfsdk:"id"`
	Tag         types.String                    `tfsdk:"tag"`
	Name        types.String                    `tfsdk:"name"`
	Description types.String                    `tfsdk:"description"`
	Children    []CatalogEntityChildSourceModel `tfsdk:"children"`
}

func (o *CatalogEntityDataSourceModel) FromApiModel(entity cortex.CatalogEntityData) {
	o.Id = types.StringValue(entity.Tag)
	o.Tag = types.StringValue(entity.Tag)
	o.Name = types.StringValue(entity.Title)
	o.Description = types.StringValue(entity.Description)

	children := make([]CatalogEntityChildSourceModel, len(entity.Children))
	for i, child := range entity.Children {
		children[i] = CatalogEntityChildSourceModel{Tag: types.StringValue(child.Tag)}
	}

	o.Children = children
}

func (d *CatalogEntityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entity"
}

func (d *CatalogEntityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Catalog Entity data source",

		Attributes: map[string]schema.Attribute{
			// Required
			"tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the catalog entity",
				Required:            true,
			},

			// Computed
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name for the entity",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the entity visible in the Service or Resource Catalog. Markdown is supported.",
				Computed:            true,
			},
			"children": schema.ListNestedAttribute{
				MarkdownDescription: "List of child entities for the entity. Only used for entities of type `TEAM` or `DOMAIN`.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the child entity.",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *CatalogEntityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CatalogEntityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CatalogEntityDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	entity, err := d.client.CatalogEntities().GetFromDescriptor(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	// Map to state
	data.FromApiModel(entity)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
