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
var _ datasource.DataSource = &RelationshipTypeDataSource{}

func NewRelationshipTypeDataSource() datasource.DataSource {
	return &RelationshipTypeDataSource{}
}

// RelationshipTypeDataSource defines the data source implementation.
type RelationshipTypeDataSource struct {
	client *cortex.HttpClient
}

func (d *RelationshipTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_relationship_type"
}

func (d *RelationshipTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	filterNestedAttributes := map[string]schema.Attribute{
		"include": schema.BoolAttribute{
			MarkdownDescription: "Whether to include or exclude the listed types/providers.",
			Computed:            true,
		},
		"types": schema.SetAttribute{
			MarkdownDescription: "Entity types to filter on.",
			ElementType:         types.StringType,
			Computed:            true,
		},
		"providers": schema.SetAttribute{
			MarkdownDescription: "Entity providers to filter on.",
			ElementType:         types.StringType,
			Computed:            true,
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Relationship Type data source",

		Attributes: map[string]schema.Attribute{
			// Required for lookup
			"tag": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (tag) for the relationship type.",
				Required:            true,
			},

			// Computed
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"allow_cycles": schema.BoolAttribute{
				Computed: true,
			},
			"is_single_source": schema.BoolAttribute{
				Computed: true,
			},
			"is_single_destination": schema.BoolAttribute{
				Computed: true,
			},
			"definition_location": schema.StringAttribute{
				Computed: true,
			},
			"sources_filter": schema.SingleNestedAttribute{
				MarkdownDescription: "Filter for which entity types can be sources.",
				Computed:            true,
				Attributes:          filterNestedAttributes,
			},
			"destinations_filter": schema.SingleNestedAttribute{
				MarkdownDescription: "Filter for which entity types can be destinations.",
				Computed:            true,
				Attributes:          filterNestedAttributes,
			},
			"source_label_singular": schema.StringAttribute{
				Computed: true,
			},
			"source_label_plural": schema.StringAttribute{
				Computed: true,
			},
			"destination_label_singular": schema.StringAttribute{
				Computed: true,
			},
			"destination_label_plural": schema.StringAttribute{
				Computed: true,
			},
			"inheritances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"parameter_tag": schema.StringAttribute{
							Computed: true,
						},
						"inheritance_type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"is_cortex_managed": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *RelationshipTypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RelationshipTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := NewRelationshipTypeDataSourceModel()

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := d.client.RelationshipTypes().Get(ctx, data.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read relationship type, got error: %s", err))
		return
	}

	data.FromApiModel(ctx, &resp.Diagnostics, entity)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
