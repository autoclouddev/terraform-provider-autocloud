package autocloud_provider

import (
	"autocloud_sdk"
	"context"

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
			},
			ResourcesMap: map[string]*schema.Resource{
				"autocloud_module": autocloudModule(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"autocloud_me": dataSourceMe(),
			},
			//ConfigureContextFunc: providerConfigure,
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

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if (username != "") && (password != "") {
			c, err := autocloud_sdk.NewClient(nil, &username, &password)
			if err != nil {
				return nil, diag.FromErr(err)
			}

			return c, diags
		}

		c, err := autocloud_sdk.NewClient(nil, nil, nil)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}
}
