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
	Links             []TeamLinkResourceModel             `tfsdk:"links"`
	SlackChannels     []TeamSlackChannelResourceModel     `tfsdk:"slack_channels"`
	AdditionalMembers []TeamAdditionalMemberResourceModel `tfsdk:"additional_members"`
}

type TeamLinkResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Url         types.String `tfsdk:"url"`
	Description types.String `tfsdk:"description,omitempty"`
}

type TeamSlackChannelResourceModel struct {
	Name                 types.String `tfsdk:"name"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

type TeamAdditionalMemberResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description,omitempty"`
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
		},
		Blocks: map[string]schema.Block{
			"link": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
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
			"slack_channel": schema.ListNestedBlock{
				MarkdownDescription: "A list of Slack channels related to the team.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the Slack channel.",
							Required:            true,
						},
						"notifications_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether notifications are enabled for the Slack channel.",
							Optional:            true,
						},
					},
				},
			},
			"additional_member": schema.ListNestedBlock{
				MarkdownDescription: "A list of additional members, outside of the IdP group. Use this field to add members like managers, PMs, etc.",
				NestedObject: schema.NestedBlockObject{
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
	teamResponse, err := r.client.Teams().Get(ctx, data.Tag.String())
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

	var slackChannels []TeamSlackChannelResourceModel
	for _, channel := range teamResponse.SlackChannels {
		slackChannels = append(slackChannels, TeamSlackChannelResourceModel{
			Name:                 types.StringValue(channel.Name),
			NotificationsEnabled: types.BoolValue(channel.NotificationsEnabled),
		})
	}
	data.SlackChannels = slackChannels

	var additionalMembers []TeamAdditionalMemberResourceModel
	for _, channel := range teamResponse.AdditionalMembers {
		additionalMembers = append(additionalMembers, TeamAdditionalMemberResourceModel{
			Name:        types.StringValue(channel.Name),
			Email:       types.StringValue(channel.Email),
			Description: types.StringValue(channel.Description),
		})
	}
	data.AdditionalMembers = additionalMembers

	var links []TeamLinkResourceModel
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

	var links []cortex.TeamLink
	for _, link := range data.Links {
		links = append(links, cortex.TeamLink{
			Name:        link.Name.String(),
			Type:        link.Type.String(),
			Url:         link.Url.String(),
			Description: link.Description.String(),
		})
	}

	var slackChannels []cortex.TeamSlackChannel
	for _, channel := range data.SlackChannels {
		slackChannels = append(slackChannels, cortex.TeamSlackChannel{
			Name:                 channel.Name.String(),
			NotificationsEnabled: channel.NotificationsEnabled.ValueBool(),
		})
	}

	var additionalMembers []cortex.TeamMember
	for _, member := range data.AdditionalMembers {
		additionalMembers = append(additionalMembers, cortex.TeamMember{
			Name:        member.Name.String(),
			Email:       member.Email.String(),
			Description: member.Description.String(),
		})
	}

	clientRequest := cortex.CreateTeamRequest{
		TeamTag:    data.Tag.String(),
		IsArchived: data.Archived.ValueBool(),
		Metadata: cortex.TeamMetadata{
			Name:        data.Name.String(),
			Description: data.Description.String(),
			Summary:     data.Summary.String(),
		},
		SlackChannels:     slackChannels,
		AdditionalMembers: additionalMembers,
		Links:             links,
		IdpGroup:          cortex.TeamIdpGroup{}, // TODO: unsure what to do with this yet, since it isn't present in the response
	}
	team, err := r.client.Teams().Create(ctx, clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = data.Tag
	data.Tag = types.StringValue(team.TeamTag)

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

	var links []cortex.TeamLink
	for _, link := range data.Links {
		links = append(links, cortex.TeamLink{
			Name:        link.Name.String(),
			Type:        link.Type.String(),
			Url:         link.Url.String(),
			Description: link.Description.String(),
		})
	}

	var slackChannels []cortex.TeamSlackChannel
	for _, channel := range data.SlackChannels {
		slackChannels = append(slackChannels, cortex.TeamSlackChannel{
			Name:                 channel.Name.String(),
			NotificationsEnabled: channel.NotificationsEnabled.ValueBool(),
		})
	}

	var additionalMembers []cortex.TeamMember
	for _, member := range data.AdditionalMembers {
		additionalMembers = append(additionalMembers, cortex.TeamMember{
			Name:        member.Name.String(),
			Email:       member.Email.String(),
			Description: member.Description.String(),
		})
	}

	clientRequest := cortex.UpdateTeamRequest{
		Metadata: cortex.TeamMetadata{
			Name:        data.Name.String(),
			Description: data.Description.String(),
			Summary:     data.Summary.String(),
		},
		SlackChannels:     slackChannels,
		AdditionalMembers: additionalMembers,
		Links:             links,
	}
	team, err := r.client.Teams().Update(ctx, data.Tag.String(), clientRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team, got error: %s", err))
		return
	}

	// Set the ID in state based on the tag
	data.Id = data.Tag
	data.Tag = types.StringValue(team.TeamTag)

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

	err := r.client.Teams().Archive(ctx, data.Tag.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete team, got error: %s", err))
		return
	}
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tag"), req, resp)
}
