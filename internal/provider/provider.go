package provider

import (
	"context"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/autocloud_module"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/repositories"
)

// entry point
func New(version string, experimental bool) func() *schema.Provider {
	return func() *schema.Provider {
		dataSourcesMap := make(map[string]*schema.Resource)
		dataSourcesMap["autocloud_github_repos"] = repositories.DataSourceRepositories()
		//dataSourcesMap["autocloud_module"] = autocloud_module.DataSourceAutocloudModule()

		if !experimental {
			dataSourcesMap["autocloud_blueprint_config"] = blueprint_config.DataSourceBlueprintConfig()
		}
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("AUTOCLOUD_TOKEN", nil),
				},
				"endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("AUTOCLOUD_API", "https://api.autocloud.io/api/v.0.0.1"),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"autocloud_blueprint": blueprint.ResourceAutocloudBlueprint(),
				"autocloud_module":    autocloud_module.ResourceAutocloudModule(),
			},
			DataSourcesMap: dataSourcesMap,
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func configure(version string, p *schema.Provider) func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	// Setup a User-Agent for your API client (replace the provider name for yours):
	// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
	// TODO: myClient.UserAgent = userAgent

	// in here we could setup multiple versions aka: "dev" "prod" and so on
	// sentry setup, etc

	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		err := blueprint_config.LoadReferencesFromState(ctx)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		apiHost := d.Get("endpoint").(string)
		if apiHost == "" {
			return nil, diag.Errorf("Autocloud API Endpoint is empty")
		}

		token := d.Get("token").(string)
		if token == "" {
			return nil, diag.Errorf("No AutoCloud auth token found")
		}
		c, err := autocloudsdk.NewClient(&apiHost, &token)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}
}
