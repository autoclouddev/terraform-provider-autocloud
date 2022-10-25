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
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("AUTOCLOUD_USERNAME", nil),
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("AUTOCLOUD_PASSWORD", nil),
				},
				"graphqlhost": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SDK_GRAPHQL_HOST", nil),
				},
				"apihost": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SDK_API_HOST", nil),
				},
				"appclientid": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SDK_COGNITO_APP_CLIENT_ID", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"autocloud_module":   autocloudModule(),
				"autocloud_module_1": autocloudModule1(), // TODO: rename it when done with the blueprint changes
			},
			DataSourcesMap: map[string]*schema.Resource{
				"autocloud_me":           dataSourceMe(),
				"autocloud_github_repos": dataSourceRepositories(),
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
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		graphql := d.Get("graphqlhost").(string)
		apiHost := d.Get("apihost").(string)
		appClientId := d.Get("appclientid").(string)
		c, err := autocloudsdk.NewClient(&graphql, &apiHost, &appClientId)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if err != nil {
			return nil, diag.FromErr(err)
		}

		if (username != "") && (password != "") {
			err := c.Login(&username, &password)
			if err != nil {
				return nil, diag.FromErr(err)
			}

			return c, diags
		}

		return c, diags
	}
}
