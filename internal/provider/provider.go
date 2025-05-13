package provider

import (
	"context"
	"fmt"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
)

// Ensure CortexProvider satisfies various provider interfaces.
var _ provider.Provider = &CortexProvider{}

const DefaultBaseApiUrl = "https://api.getcortexapp.com"

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
				Optional:            true,
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
	baseApiUrl := DefaultBaseApiUrl
	if data.BaseApiUrl.IsNull() || data.BaseApiUrl.ValueString() == "" {
		envApiUrl := os.Getenv("CORTEX_API_URL")
		if envApiUrl != "" {
			baseApiUrl = envApiUrl
		}
	} else {
		baseApiUrl = data.BaseApiUrl.ValueString()
	}
	if data.Token.IsUnknown() || data.Token.ValueString() == "" {
		token := os.Getenv("CORTEX_API_TOKEN")
		if token == "" {
			resp.Diagnostics.AddAttributeError(path.Root("token"), "token is required", "Please specify an API token for the Cortex API")
			return
		}
		data.Token = types.StringValue(token)
	}

	// Creating a new GitLab Client from the provider configuration
	client, err := cortex.NewClient(
		cortex.WithContext(ctx),
		cortex.WithURL(baseApiUrl),
		cortex.WithToken(data.Token.ValueString()),
		cortex.WithVersion(p.version),
	)

	if err != nil {
		resp.Diagnostics.AddError("Failed to create Cortex API Client from provider configuration", fmt.Sprintf("The provider failed to create a new Cortex API Client from the given configuration: %+v", err))
		return
	}

	// Example client configuration for data sources and resources
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CortexProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCatalogEntityResource,
		NewDepartmentResource,
		NewScorecardResource,
		NewResourceDefinitionResource,
		NewCatalogEntityCustomDataResource,
		NewCatalogEntityOpenAPIResource,
	}
}

func (p *CortexProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCatalogEntityDataSource,
		NewTeamDataSource,
		NewDepartmentDataSource,
		NewScorecardDataSource,
		NewResourceDefinitionDataSource,
		NewCatalogEntityCustomDataDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CortexProvider{
			version: version,
		}
	}
}
