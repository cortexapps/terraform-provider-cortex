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
var _ datasource.DataSource = &DepartmentDataSource{}

func NewDepartmentDataSource() datasource.DataSource {
	return &DepartmentDataSource{}
}

// DepartmentDataSource defines the data source implementation.
type DepartmentDataSource struct {
	client *cortex.HttpClient
}

// DepartmentDataSourceModel describes the data source data model.
type DepartmentDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Tag         types.String `tfsdk:"tag"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (d *DepartmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_department"
}

func (d *DepartmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Department data source",

		Attributes: map[string]schema.Attribute{
			// Required
			"tag": schema.StringAttribute{
				MarkdownDescription: "Tag of the department",
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
		},
	}
}

func (d *DepartmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DepartmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DepartmentDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	department, err := d.client.Departments().Get(ctx, data.Tag.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read department, got error: %s", err))
		return
	}
	data.Id = types.StringValue(department.Tag)
	data.Tag = types.StringValue(department.Tag)
	data.Name = types.StringValue(department.Name)
	data.Description = types.StringValue(department.Description)

	// Write to TF state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
