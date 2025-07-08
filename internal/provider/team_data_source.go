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
var _ datasource.DataSource = &TeamDataSource{}

func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

// TeamDataSource defines the data source implementation.
type TeamDataSource struct {
	client *cortex.HttpClient
}

// TeamDataSourceModel describes the data source data model.
type TeamDataSourceModel struct {
	Id            types.String              `tfsdk:"id"`
	Tag           types.String              `tfsdk:"tag"`
	Members       []cortex.TeamMember       `tfsdk:"members"`
	SlackChannels []cortex.TeamSlackChannel `tfsdk:"slack_channels"`
}

func (d *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Team data source",

		Attributes: map[string]schema.Attribute{
			// Required
			"tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the team",
				Required:            true,
			},

			// Computed
			"id": schema.StringAttribute{
				Computed: true,
			},
			"members": schema.ListNestedAttribute{
				MarkdownDescription: "List of team members.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"slack_channels": schema.ListNestedAttribute{
				MarkdownDescription: "List of Slack channels associated with the team.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"notifications_enabled": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *TeamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamResponse, err := d.client.Teams().Get(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team, got error: %s", err))
		return
	}

	data.Id = types.StringValue(teamResponse.TeamTag)
	data.Tag = types.StringValue(teamResponse.TeamTag)

	data.Members = teamResponse.CortexTeam.Members
	data.SlackChannels = teamResponse.SlackChannels

	// Write to TF state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
