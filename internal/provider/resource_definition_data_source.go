package provider

import (
	"context"
	"fmt"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ResourceDefinitionDataSource{}

func NewResourceDefinitionDataSource() datasource.DataSource {
	return &ResourceDefinitionDataSource{}
}

// ResourceDefinitionDataSource defines the data source implementation.
type ResourceDefinitionDataSource struct {
	client *cortex.HttpClient
}

// ResourceDefinitionDataSourceModel describes the data source data model.
type ResourceDefinitionDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Source      types.String `tfsdk:"source"`
	Schema      types.Map    `tfsdk:"schema"`
}

func (d *ResourceDefinitionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_definition"
}

func (d *ResourceDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ResourceDefinition data source",

		Attributes: map[string]schema.Attribute{
			// Required
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of resource definition",
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
			"source": schema.StringAttribute{
				Computed: true,
			},
			"schema": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *ResourceDefinitionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ResourceDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ResourceDefinitionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceDefinition, err := d.client.ResourceDefinitions().Get(ctx, data.Type.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read resource definition, got error: %s", err))
		return
	}

	// attempt to extract data from schema
	schemaData, diags := types.MapValueFrom(ctx, types.StringType, resourceDefinition.Schema)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(resourceDefinition.Type)
	data.Type = types.StringValue(resourceDefinition.Type)
	data.Name = types.StringValue(resourceDefinition.Name)
	data.Description = types.StringValue(resourceDefinition.Description)
	data.Source = types.StringValue(resourceDefinition.Source)
	data.Schema = schemaData

	// Write to TF state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
