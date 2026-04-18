package provider

import (
	"context"
	"fmt"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AwsIntegrationResource{}
var _ resource.ResourceWithImportState = &AwsIntegrationResource{}

func NewAwsIntegrationResource() resource.Resource {
	return &AwsIntegrationResource{}
}

func NewAwsIntegrationResourceModel() AwsIntegrationResourceModel {
	return AwsIntegrationResourceModel{}
}

type AwsIntegrationResource struct {
	client *cortex.HttpClient
}

func (r *AwsIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "AWS Integration Configuration",

		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				MarkdownDescription: "The AWS Account ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"iam_role": schema.StringAttribute{
				MarkdownDescription: "The IAM role Cortex would be assuming.",
				Required:            true,
			},
			"account_alias": schema.StringAttribute{
				MarkdownDescription: "The account alias for the AWS account.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AwsIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_integration_configuration" // cortex_aws_integration_configuration
}

func (r *AwsIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AwsIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NewAwsIntegrationResourceModel()

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entity, err := r.client.Aws().GetConfiguration(ctx, data.AccountId.ValueString())
	if err != nil {
		if err == cortex.ApiErrorNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read AWS integration configuration %s, got error: %s", data.AccountId.ValueString(), err))
		return
	}

	data.FromApiModel(entity)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NewAwsIntegrationResourceModel()

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel()
	
	entity, err := r.client.Aws().CreateConfiguration(ctx, clientEntity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create AWS integration configuration, got error: %s", err))
		return
	}

	data.FromApiModel(entity)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := NewAwsIntegrationResourceModel()

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientEntity := data.ToApiModel()
	
	entity, err := r.client.Aws().UpdateConfiguration(ctx, clientEntity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update AWS integration configuration, got error: %s", err))
		return
	}

	data.FromApiModel(entity)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AwsIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := NewAwsIntegrationResourceModel()

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Aws().DeleteConfiguration(ctx, data.AccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete AWS integration configuration, got error: %s", err))
		return
	}
}

func (r *AwsIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("account_id"), req, resp)
}
