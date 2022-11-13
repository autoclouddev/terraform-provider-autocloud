package autocloud_provider

import (
	"context"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// entry point
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("AUTOCLOUD_TOKEN", nil),
				},
				"apihost": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SDK_API_HOST", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"autocloud_blueprint": autocloudBlueprint(),
				"autocloud_module":    autocloudModule(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"autocloud_github_repos":        dataSourceRepositories(),
				"autocloud_module":              dataSourceAutocloudModule(),
				"autocloud_terraform_processor": dataSourceTerraformProcessor(),
			},
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
		apiHost := d.Get("apihost").(string)
		token := d.Get("token").(string)
		c, err := autocloudsdk.NewClient(&apiHost, &token)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}
}
