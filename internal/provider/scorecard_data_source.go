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
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"draft": schema.BoolAttribute{
				Computed: true,
			},
			"ladder": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"levels": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Computed: true,
								},
								"rank": schema.Int64Attribute{
									Computed: true,
								},
								"color": schema.StringAttribute{
									Computed: true,
								},
								"description": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
			"rules": schema.SetNestedAttribute{
				MarkdownDescription: "Rules of the scorecard.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// Required rule attributes
						"title": schema.StringAttribute{
							Computed: true,
						},
						"expression": schema.StringAttribute{
							Computed: true,
						},
						"weight": schema.Int64Attribute{
							Computed: true,
						},
						"level": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"failure_message": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"filter": schema.SingleNestedAttribute{
				MarkdownDescription: "Filter of the scorecard.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"types": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"include": schema.SetAttribute{
								ElementType: types.StringType,
								Computed:    true,
							},
							"exclude": schema.SetAttribute{
								ElementType: types.StringType,
								Computed:    true,
							},
						},
					},
					"groups": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"include": schema.SetAttribute{
								ElementType: types.StringType,
								Computed:    true,
							},
							"exclude": schema.SetAttribute{
								ElementType: types.StringType,
								Computed:    true,
							},
						},
					},
					"query": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"evaluation": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"window": schema.Int64Attribute{
						Computed: true,
					},
				},
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
