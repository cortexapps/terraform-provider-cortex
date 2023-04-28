package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure CortexProvider satisfies various provider interfaces.
var _ provider.Provider = &CortexProvider{}

// CortexProvider defines the provider implementation.
type CortexProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CortexProviderModel describes the provider data model.
type CortexProviderModel struct {
	BaseApiUrl types.String `tfsdk:"base_api_url"`
	Token      types.String `tfsdk:"token"`
}

func (p *CortexProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cortex"
	resp.Version = p.version
}

func (p *CortexProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_api_url": schema.StringAttribute{
				MarkdownDescription: "Base URL to the Cortex API",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The API token used to authenticate with Cortex",
				Optional:            false,
				Sensitive:           true,
			},
		},
	}
}

func (p *CortexProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CortexProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.BaseApiUrl.IsNull() {
		data.BaseApiUrl = types.StringValue("https://api.getcortexapp.com/api/v1/")
	}
	if data.Token.IsNull() {
		resp.Diagnostics.AddError("token is required", "Please specify an API token for the Cortex API")
		return
	}

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CortexProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCatalogEntityResource,
	}
}

func (p *CortexProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCatalogEntityDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CortexProvider{
			version: version,
		}
	}
}
