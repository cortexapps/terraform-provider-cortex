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
var _ datasource.DataSource = &TeamsDataSource{}

func NewTeamsDataSource() datasource.DataSource {
	return &TeamsDataSource{}
}

// TeamsDataSource defines the data source implementation.
type TeamsDataSource struct {
	client *cortex.HttpClient
}

// TeamsDataSourceModel describes the data source data model.
type TeamsDataSourceModel struct {
	Id    types.String              `tfsdk:"id"`
	Teams []TeamDataSourceItemModel `tfsdk:"teams"`
}

// TeamDataSourceItemModel represents a single team in the list.
type TeamDataSourceItemModel struct {
	Tag           types.String              `tfsdk:"tag"`
	Name          types.String              `tfsdk:"name"`
	Description   types.String              `tfsdk:"description"`
	IsArchived    types.Bool                `tfsdk:"is_archived"`
	Members       []cortex.TeamMember       `tfsdk:"members"`
	SlackChannels []cortex.TeamSlackChannel `tfsdk:"slack_channels"`
}

func (d *TeamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (d *TeamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Teams data source - returns a list of all teams",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier for this data source",
			},
			"teams": schema.ListNestedAttribute{
				MarkdownDescription: "List of all teams",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							MarkdownDescription: "Tag of the team",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the team",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the team",
							Computed:            true,
						},
						"is_archived": schema.BoolAttribute{
							MarkdownDescription: "Whether the team is archived",
							Computed:            true,
						},
						"members": schema.ListNestedAttribute{
							MarkdownDescription: "List of team members",
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
							MarkdownDescription: "List of Slack channels associated with the team",
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
				},
			},
		},
	}
}

func (d *TeamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamsResponse, err := d.client.Teams().List(ctx, &cortex.TeamListParams{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read teams, got error: %s", err))
		return
	}

	data.Id = types.StringValue("teams")

	// Convert API response to model
	teams := make([]TeamDataSourceItemModel, len(teamsResponse.Teams))
	for i, team := range teamsResponse.Teams {

		teams[i] = TeamDataSourceItemModel{
			Tag:           types.StringValue(team.TeamTag),
			Name:          types.StringValue(team.Metadata.Name),
			Description:   types.StringValue(team.Metadata.Description),
			IsArchived:    types.BoolValue(team.IsArchived),
			Members:       team.CortexTeam.Members,
			SlackChannels: team.SlackChannels,
		}
	}

	data.Teams = teams

	// Write to TF state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
