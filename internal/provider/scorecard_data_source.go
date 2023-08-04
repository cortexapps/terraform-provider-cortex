package provider

import (
	"context"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ScorecardDataSource{}

func NewScorecardDataSource() datasource.DataSource {
	return &ScorecardDataSource{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// ScorecardDataSource defines the data source implementation.
type ScorecardDataSource struct {
	client *cortex.HttpClient
}

/***********************************************************************************************************************
 * Functions
 **********************************************************************************************************************/

func (d *ScorecardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard"
}

func (d *ScorecardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Scorecard data source",

		Attributes: map[string]schema.Attribute{
			// Required
			"tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the scorecard",
				Required:            true,
			},

			// Computed
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *ScorecardDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScorecardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ScorecardDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := d.client.Scorecards().Get(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scorecard, got error: %s", err))
		return
	}
	data.FromApiModel(ctx, &resp.Diagnostics, &entity)

	// Write to TF state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
