package provider_plugin

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/autocloud_module"
)

type autoProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

// autocloudProvider is the provider implementation.
type autocloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &autocloudProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &autocloudProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *autocloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "autocloud"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *autocloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"endpoint": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *autocloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config autoProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"API Endpoint is required",
			"The provider cannot create a client without an endpoint",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Token is required",
			"The provider cannot create a client without a token",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("AUTOCLOUD_API")
	token := os.Getenv("AUTOCLOUD_TOKEN")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if endpoint == "" {
		endpoint = "https://api.autocloud.io/api/v.0.0.1"
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Token is required",
			"The provider cannot work without a token, add it from autocloud.io",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := autocloudsdk.NewClient(&endpoint, &token)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create client",
			"An unexpected error occurred while creating the client, please contact the provider delevelores.\n"+
				"Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *autocloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *autocloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		autocloud_module.NewModuleResource,
	}
}
