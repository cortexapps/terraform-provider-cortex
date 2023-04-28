package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CatalogEntityDataSource{}

func NewCatalogEntityDataSource() datasource.DataSource {
	return &CatalogEntityDataSource{}
}

// CatalogEntityDataSource defines the data source implementation.
type CatalogEntityDataSource struct {
	client *http.Client
}

// CatalogEntityDataSourceModel describes the data source data model.
type CatalogEntityDataSourceModel struct {
	Tag types.String `tfsdk:"tag"`
}

func (d *CatalogEntityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (d *CatalogEntityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Catalog Entity data source",

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
				MarkdownDescription: "Description of the entity visible in the Service or Resource Catalog. Markdown is supported.",
				Computed:            true,
			},
		},
	}
}

func (d *CatalogEntityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Tag = types.StringValue("tag")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
