package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

/***********************************************************************************************************************
 * Models
 **********************************************************************************************************************/

// CatalogEntityDataSourceModel describes the data source data model.
type CatalogEntityDataSourceModel struct {
	Id            types.String                 `tfsdk:"id"`
	Tag           types.String                 `tfsdk:"tag"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Type          types.String                 `tfsdk:"type"`
	IsArchived    types.Bool                   `tfsdk:"is_archived"`
	LastUpdated   types.String                 `tfsdk:"last_updated"`
	Groups        []types.String               `tfsdk:"groups"`
	Links         []CatalogEntityLinkModel     `tfsdk:"links"`
	Metadata      []CatalogEntityMetadataModel `tfsdk:"metadata"`
	Ownership     types.Object                 `tfsdk:"ownership"`
	SlackChannels []types.Object               `tfsdk:"slack_channels"`
	Git           types.Object                 `tfsdk:"git"`
}

type CatalogEntityLinkModel struct {
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Url         types.String `tfsdk:"url"`
	Description types.String `tfsdk:"description"`
}

type CatalogEntityMetadataModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type CatalogEntityOwnershipModel struct {
	Emails []types.Object `tfsdk:"emails"`
	Groups []types.Object `tfsdk:"groups"`
}

type CatalogEntityOwnershipEmailModel struct {
	Email       types.String `tfsdk:"email"`
	Description types.String `tfsdk:"description"`
	Inheritance types.String `tfsdk:"inheritance"`
}

type CatalogEntityOwnershipGroupModel struct {
	GroupName   types.String `tfsdk:"group_name"`
	Description types.String `tfsdk:"description"`
	Provider    types.String `tfsdk:"provider"`
	Inheritance types.String `tfsdk:"inheritance"`
}

type CatalogEntitySlackChannelModel struct {
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

type CatalogEntityGitModel struct {
	Provider      types.String `tfsdk:"provider"`
	Repository    types.String `tfsdk:"repository"`
	Alias         types.String `tfsdk:"alias"`
	Basepath      types.String `tfsdk:"basepath"`
	RepositoryUrl types.String `tfsdk:"repository_url"`
}

// AttrTypes for nested objects.
func (o *CatalogEntityOwnershipModel) AttrTypes() map[string]attr.Type {
	email := CatalogEntityOwnershipEmailModel{}
	group := CatalogEntityOwnershipGroupModel{}
	return map[string]attr.Type{
		"emails": types.ListType{ElemType: types.ObjectType{AttrTypes: email.AttrTypes()}},
		"groups": types.ListType{ElemType: types.ObjectType{AttrTypes: group.AttrTypes()}},
	}
}

func (o *CatalogEntityOwnershipEmailModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"email":       types.StringType,
		"description": types.StringType,
		"inheritance": types.StringType,
	}
}

func (o *CatalogEntityOwnershipGroupModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"group_name":  types.StringType,
		"description": types.StringType,
		"provider":    types.StringType,
		"inheritance": types.StringType,
	}
}

func (o *CatalogEntitySlackChannelModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                  types.StringType,
		"description":           types.StringType,
		"notifications_enabled": types.BoolType,
	}
}

func (o *CatalogEntityGitModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"provider":       types.StringType,
		"repository":     types.StringType,
		"alias":          types.StringType,
		"basepath":       types.StringType,
		"repository_url": types.StringType,
	}
}

func (o *CatalogEntityDataSourceModel) FromApiModel(ctx context.Context, diagnostics *diag.Diagnostics, entity *cortex.CatalogEntity) {
	o.Id = types.StringValue(entity.Tag)
	o.Tag = types.StringValue(entity.Tag)
	o.Name = types.StringValue(entity.Name)
	o.Description = types.StringValue(entity.Description)
	o.Type = types.StringValue(entity.Type)
	o.IsArchived = types.BoolValue(entity.IsArchived)
	o.LastUpdated = types.StringValue(entity.LastUpdated)

	o.Groups = make([]types.String, len(entity.Groups))
	for i, g := range entity.Groups {
		o.Groups[i] = types.StringValue(g)
	}

	o.Links = make([]CatalogEntityLinkModel, len(entity.Links))
	for i, l := range entity.Links {
		o.Links[i] = CatalogEntityLinkModel{
			Name:        types.StringValue(l.Name),
			Type:        types.StringValue(l.Type),
			Url:         types.StringValue(l.Url),
			Description: types.StringValue(l.Description),
		}
	}

	o.Metadata = make([]CatalogEntityMetadataModel, len(entity.Metadata))
	for i, m := range entity.Metadata {
		valueBytes, err := json.Marshal(m.Value)
		if err != nil {
			diagnostics.AddWarning("Metadata serialization", fmt.Sprintf("Could not serialize metadata key %s: %s", m.Key, err))
			continue
		}
		o.Metadata[i] = CatalogEntityMetadataModel{
			Key:   types.StringValue(m.Key),
			Value: types.StringValue(string(valueBytes)),
		}
	}

	// Handle Ownership
	ownershipModel := CatalogEntityOwnershipModel{}
	emails := make([]types.Object, len(entity.Ownership.Emails))
	for i, e := range entity.Ownership.Emails {
		emailModel := CatalogEntityOwnershipEmailModel{
			Email:       types.StringValue(e.Email),
			Description: types.StringValue(e.Description),
			Inheritance: types.StringValue(e.Inheritance),
		}
		emailObj, d := types.ObjectValueFrom(ctx, emailModel.AttrTypes(), emailModel)
		diagnostics.Append(d...)
		emails[i] = emailObj
	}

	groups := make([]types.Object, len(entity.Ownership.Groups))
	for i, g := range entity.Ownership.Groups {
		groupModel := CatalogEntityOwnershipGroupModel{
			GroupName:   types.StringValue(g.GroupName),
			Description: types.StringValue(g.Description),
			Provider:    types.StringValue(g.Provider),
			Inheritance: types.StringValue(g.Inheritance),
		}
		groupObj, d := types.ObjectValueFrom(ctx, groupModel.AttrTypes(), groupModel)
		diagnostics.Append(d...)
		groups[i] = groupObj
	}

	ownershipModel.Emails = emails
	ownershipModel.Groups = groups
	ownershipObj, d := types.ObjectValueFrom(ctx, ownershipModel.AttrTypes(), ownershipModel)
	diagnostics.Append(d...)
	o.Ownership = ownershipObj

	// Handle SlackChannels
	o.SlackChannels = make([]types.Object, len(entity.SlackChannels))
	for i, sc := range entity.SlackChannels {
		scModel := CatalogEntitySlackChannelModel{
			Name:                 types.StringValue(sc.Name),
			Description:          types.StringValue(sc.Description),
			NotificationsEnabled: types.BoolValue(sc.NotificationsEnabled),
		}
		scObj, d := types.ObjectValueFrom(ctx, scModel.AttrTypes(), scModel)
		diagnostics.Append(d...)
		o.SlackChannels[i] = scObj
	}

	// Handle Git
	gitModel := CatalogEntityGitModel{
		Provider:      types.StringValue(entity.Git.Provider),
		Repository:    types.StringValue(entity.Git.Repository),
		Alias:         types.StringValue(entity.Git.Alias),
		Basepath:      types.StringValue(entity.Git.Basepath),
		RepositoryUrl: types.StringValue(entity.Git.RepositoryUrl),
	}
	gitObj, d := types.ObjectValueFrom(ctx, gitModel.AttrTypes(), gitModel)
	diagnostics.Append(d...)
	o.Git = gitObj
}

/***********************************************************************************************************************
 * Data Source implementation
 **********************************************************************************************************************/

func (d *CatalogEntityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entity"
}

func (d *CatalogEntityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the details of a single catalog entity by its tag.",

		Attributes: map[string]schema.Attribute{
			// Required
			"tag": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (tag) of the catalog entity.",
				Required:            true,
			},

			// Computed
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name for the entity.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the entity visible in the Service or Resource Catalog. Markdown is supported.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the entity (e.g. `service`, `resource`, `domain`).",
				Computed:            true,
			},
			"is_archived": schema.BoolAttribute{
				MarkdownDescription: "Whether the entity is archived.",
				Computed:            true,
			},
			"last_updated": schema.StringAttribute{
				MarkdownDescription: "ISO 8601 timestamp of when the entity was last updated.",
				Computed:            true,
			},
			"groups": schema.ListAttribute{
				MarkdownDescription: "Groups the entity belongs to. Corresponds to `x-cortex-groups`.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"links": schema.ListNestedAttribute{
				MarkdownDescription: "Links associated with the entity.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":        schema.StringAttribute{Computed: true},
						"type":        schema.StringAttribute{Computed: true, MarkdownDescription: "Link type (e.g. `runbook`, `documentation`, `dashboard`)."},
						"url":         schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
					},
				},
			},
			"metadata": schema.ListNestedAttribute{
				MarkdownDescription: "Custom metadata key/value pairs associated with the entity. Values are JSON-encoded strings.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key":   schema.StringAttribute{Computed: true},
						"value": schema.StringAttribute{Computed: true, MarkdownDescription: "JSON-encoded value."},
					},
				},
			},
			"ownership": schema.SingleNestedAttribute{
				MarkdownDescription: "Ownership details for the entity.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"emails": schema.ListNestedAttribute{
						MarkdownDescription: "Individual email owners.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"email":       schema.StringAttribute{Computed: true},
								"description": schema.StringAttribute{Computed: true},
								"inheritance": schema.StringAttribute{Computed: true},
							},
						},
					},
					"groups": schema.ListNestedAttribute{
						MarkdownDescription: "Group owners.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"group_name":  schema.StringAttribute{Computed: true},
								"description": schema.StringAttribute{Computed: true},
								"provider":    schema.StringAttribute{Computed: true},
								"inheritance": schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"slack_channels": schema.ListNestedAttribute{
				MarkdownDescription: "Slack channels associated with the entity.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":                  schema.StringAttribute{Computed: true},
						"description":           schema.StringAttribute{Computed: true},
						"notifications_enabled": schema.BoolAttribute{Computed: true},
					},
				},
			},
			"git": schema.SingleNestedAttribute{
				MarkdownDescription: "Git repository information associated with the entity.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"provider":       schema.StringAttribute{Computed: true, MarkdownDescription: "Git provider (e.g. `github`, `gitlab`)."},
					"repository":     schema.StringAttribute{Computed: true},
					"alias":          schema.StringAttribute{Computed: true},
					"basepath":       schema.StringAttribute{Computed: true},
					"repository_url": schema.StringAttribute{Computed: true},
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
			fmt.Sprintf("Expected *cortex.HttpClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	entity, err := d.client.CatalogEntities().Get(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read catalog entity, got error: %s", err))
		return
	}

	data.FromApiModel(ctx, &resp.Diagnostics, entity)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
