package provider

import (
	"context"
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
var _ resource.Resource = &TeamResource{}
var _ resource.ResourceWithImportState = &TeamResource{}

func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// TeamResource defines the resource implementation.
type TeamResource struct {
	client *cortex.HttpClient
}

// TeamResourceModel describes the team data model within Terraform.
type TeamResourceModel struct {
	Id                types.String                        `tfsdk:"id"`
	Tag               types.String                        `tfsdk:"tag"`
	Name              types.String                        `tfsdk:"name"`
	Description       types.String                        `tfsdk:"description"`
	Summary           types.String                        `tfsdk:"summary"`
	Archived          types.Bool                          `tfsdk:"archived"`
	SlackChannels     []TeamSlackChannelResourceModel     `tfsdk:"slack_channels"`
	Links             []TeamLinkResourceModel             `tfsdk:"links"`
	AdditionalMembers []TeamAdditionalMemberResourceModel `tfsdk:"additional_members"`
}

type TeamLinkResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Url         types.String `tfsdk:"url"`
	Description types.String `tfsdk:"description"`
}

type TeamSlackChannelResourceModel struct {
	Name                 types.String `tfsdk:"name"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

type TeamAdditionalMemberResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
}

/***********************************************************************************************************************
 * Schema
 **********************************************************************************************************************/

func (r *TeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Team Entity",

		Attributes: map[string]schema.Attribute{
			// Required attributes
			"tag": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the team.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the team.",
				Required:            true,
			},
			"slack_channels": schema.ListNestedAttribute{
				MarkdownDescription: "A list of Slack channels related to the team.",
				Required:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the Slack channel.",
							Required:            true,
						},
						"notifications_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether or not notifications are enabled for the Slack channel.",
							Required:            true,
						},
					},
				},
			},
			"links": schema.ListNestedAttribute{
				MarkdownDescription: "Links related to the team.",
				Required:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the link.",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of link.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("runbook", "documentation", "logs", "dashboard", "metrics", "healthcheck"),
							},
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "URL of the link.",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A short description of the link.",
							Optional:            true,
						},
					},
				},
			},
			"additional_members": schema.ListNestedAttribute{
				MarkdownDescription: "A list of additional members, outside of the IdP group. Use this field to add members like managers, PMs, etc.",
				Required:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the member.",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A short description of the member.",
							Optional:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Email of the member.",
							Required:            true,
						},
					},
				},
			},

			// Optional attributes
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team.",
				Optional:            true,
			},
			"summary": schema.StringAttribute{
				MarkdownDescription: "Summary of the team.",
				Optional:            true,
			},

			// Computed attributes
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"archived": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

/***********************************************************************************************************************
 * Methods
 **********************************************************************************************************************/

func (r *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *TeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TeamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Issue API request
	teamResponse, err := r.client.Teams().Get(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team, got error: %s", err))
		return
	}

	// Map data from the API response to the model
	data.Id = types.StringValue(teamResponse.TeamTag)
	data.Tag = types.StringValue(teamResponse.TeamTag)
	data.Name = types.StringValue(teamResponse.Metadata.Name)
	data.Description = types.StringValue(teamResponse.Metadata.Description)
	data.Summary = types.StringValue(teamResponse.Metadata.Summary)
	data.Archived = types.BoolValue(teamResponse.IsArchived)

	var slackChannels = make([]TeamSlackChannelResourceModel, 0)
	for _, channel := range teamResponse.SlackChannels {
		slackChannels = append(slackChannels, TeamSlackChannelResourceModel{
			Name:                 types.StringValue(channel.Name),
			NotificationsEnabled: types.BoolValue(channel.NotificationsEnabled),
		})
	}
	data.SlackChannels = slackChannels

	var additionalMembers = make([]TeamAdditionalMemberResourceModel, 0)
	for _, channel := range teamResponse.AdditionalMembers {
		additionalMembers = append(additionalMembers, TeamAdditionalMemberResourceModel{
			Name:        types.StringValue(channel.Name),
			Email:       types.StringValue(channel.Email),
			Description: types.StringValue(channel.Description),
		})
	}
	data.AdditionalMembers = additionalMembers

	var links = make([]TeamLinkResourceModel, 0)
	for _, link := range teamResponse.Links {
		links = append(links, TeamLinkResourceModel{
			Name:        types.StringValue(link.Name),
			Type:        types.StringValue(link.Type),
			Url:         types.StringValue(link.Url),
			Description: types.StringValue(link.Description),
		})
	}
	data.Links = links

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create Creates a new team.
func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TeamResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var links = make([]cortex.TeamLink, 0)
	for _, link := range data.Links {
		links = append(links, cortex.TeamLink{
			Name:        link.Name.ValueString(),
			Type:        link.Type.ValueString(),
			Url:         link.Url.ValueString(),
			Description: link.Description.ValueString(),
		})
	}

	var slackChannels = make([]cortex.TeamSlackChannel, 0)
	for _, channel := range data.SlackChannels {
		slackChannels = append(slackChannels, cortex.TeamSlackChannel{
			Name:                 channel.Name.ValueString(),
			NotificationsEnabled: channel.NotificationsEnabled.ValueBool(),
		})
	}

	var additionalMembers = make([]cortex.TeamMember, 0)
	for _, member := range data.AdditionalMembers {
		additionalMembers = append(additionalMembers, cortex.TeamMember{
			Name:        member.Name.ValueString(),
			Email:       member.Email.ValueString(),
			Description: member.Description.ValueString(),
		})
	}

	clientRequest := cortex.CreateTeamRequest{
		TeamTag:    data.Tag.ValueString(),
		IsArchived: data.Archived.ValueBool(),
		Metadata: cortex.TeamMetadata{
			Name:        data.Name.ValueString(),
			Description: data.Description.ValueString(),
			Summary:     data.Summary.ValueString(),
		},
		SlackChannels:     slackChannels,
		AdditionalMembers: additionalMembers,
		Links:             links,
	}
	if true { // TODO: key off of IDP group vs Cortex managed
		clientRequest.Type = "CORTEX"
		clientRequest.CortexTeam = cortex.TeamCortexManaged{
			Members: make([]cortex.TeamMember, 0),
		}
	} else {
		clientRequest.IdpGroup = cortex.TeamIdpGroup{
			Group:    data.Tag.ValueString(),
			Provider: "OKTA", // TODO make dynamic in TF syntax
			Members:  []cortex.TeamIdpGroupMember{},
		}
	}
	team, err := r.client.Teams().Create(ctx, clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = data.Tag
	data.Tag = types.StringValue(team.TeamTag)
	data.Archived = types.BoolValue(team.IsArchived)
	// TODO: Add other attributes, consolidate this into a shared method

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	//resp.State.SetAttribute(ctx, path.Root("id"), data.TeamTag)
	//resp.State.SetAttribute(ctx, path.Root("tag"), data.TeamTag)
}

func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TeamResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var links = make([]cortex.TeamLink, 0)
	for _, link := range data.Links {
		links = append(links, cortex.TeamLink{
			Name:        link.Name.ValueString(),
			Type:        link.Type.ValueString(),
			Url:         link.Url.ValueString(),
			Description: link.Description.ValueString(),
		})
	}

	var slackChannels = make([]cortex.TeamSlackChannel, 0)
	for _, channel := range data.SlackChannels {
		slackChannels = append(slackChannels, cortex.TeamSlackChannel{
			Name:                 channel.Name.ValueString(),
			NotificationsEnabled: channel.NotificationsEnabled.ValueBool(),
		})
	}
	//
	var additionalMembers = make([]cortex.TeamMember, 0)
	for _, member := range data.AdditionalMembers {
		additionalMembers = append(additionalMembers, cortex.TeamMember{
			Name:        member.Name.ValueString(),
			Email:       member.Email.ValueString(),
			Description: member.Description.ValueString(),
		})
	}

	clientRequest := cortex.UpdateTeamRequest{
		Metadata: cortex.TeamMetadata{
			Name:        data.Name.ValueString(),
			Description: data.Description.ValueString(),
			Summary:     data.Summary.ValueString(),
		},
		SlackChannels:     slackChannels,
		AdditionalMembers: additionalMembers,
		Links:             links,
	}
	team, err := r.client.Teams().Update(ctx, data.Tag.ValueString(), clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = data.Tag
	data.Tag = types.StringValue(team.TeamTag)
	data.Archived = types.BoolValue(team.IsArchived)
	// TODO: Add other attributes, consolidate this into a shared method

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TeamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Teams().Delete(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team, got error: %s", err))
		return
	}
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
