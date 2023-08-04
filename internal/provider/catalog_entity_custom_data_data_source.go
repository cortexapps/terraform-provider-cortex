package provider

import (
	"context"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CatalogEntityCustomDataDataSource{}

func NewCatalogEntityCustomDataDataSource() datasource.DataSource {
	return &CatalogEntityCustomDataDataSource{}
}

// CatalogEntityCustomDataDataSource defines the data source implementation.
type CatalogEntityCustomDataDataSource struct {
	client *cortex.HttpClient
}

func (d *CatalogEntityCustomDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entity_custom_data"
}

func (d *CatalogEntityCustomDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Catalog Entity Custom Data data source",

		Attributes: map[string]schema.Attribute{
			// Required
			"tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the catalog entity",
				Required:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key of the custom data entry",
				Required:            true,
			},

			// Computed
			"id": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the entity visible in the Service or Resource Catalog. Markdown is supported.",
				Computed:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Source of the custom data.",
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value of the custom data attribute.",
				Computed:            true,
			},
		},
	}
}

func (d *CatalogEntityCustomDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CatalogEntityCustomDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := NewCatalogEntityCustomDataDataSourceModel()

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	entity, err := d.client.CatalogEntityCustomData().Get(ctx, data.Tag.ValueString(), data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	// Map to state
	data.FromApiModel(ctx, &resp.Diagnostics, entity)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
