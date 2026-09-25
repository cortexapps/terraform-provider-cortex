package provider

import (
	"context"
	"fmt"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &AwsIntegrationTypesResource{}
var _ resource.ResourceWithImportState = &AwsIntegrationTypesResource{}

func NewAwsIntegrationTypesResource() resource.Resource {
	return &AwsIntegrationTypesResource{}
}

type AwsIntegrationTypesResource struct {
	client *cortex.HttpClient
}

func (r *AwsIntegrationTypesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "AWS Integration Types Configuration",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"types": schema.ListNestedAttribute{
				MarkdownDescription: "List of AWS types and whether they're configured.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "The AWS type.",
							Required:            true,
						},
						"configured": schema.BoolAttribute{
							MarkdownDescription: "Whether the type is configured.",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *AwsIntegrationTypesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_integration_types"
}

func (r *AwsIntegrationTypesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cortex.HttpClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *AwsIntegrationTypesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AwsIntegrationTypesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.Aws().GetTypes(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read AWS types, got error: %s", err))
		return
	}

	data.FromApiModel(entity)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsIntegrationTypesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AwsIntegrationTypesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel()
	
	entity, err := r.client.Aws().UpdateTypes(ctx, clientEntity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create AWS types, got error: %s", err))
		return
	}

	data.FromApiModel(entity)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsIntegrationTypesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AwsIntegrationTypesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel()
	
	entity, err := r.client.Aws().UpdateTypes(ctx, clientEntity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update AWS types, got error: %s", err))
		return
	}

	data.FromApiModel(entity)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsIntegrationTypesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AwsIntegrationTypesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Aws().UpdateTypes(ctx, []cortex.AwsType{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete AWS types, got error: %s", err))
		return
	}
}

func (r *AwsIntegrationTypesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "aws_types")...)
}
