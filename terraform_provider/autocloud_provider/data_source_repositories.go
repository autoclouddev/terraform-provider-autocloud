package autocloud_provider

import (
	"autocloud_sdk"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRepositories() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "me data source, to check current user (do not publish this).",

		ReadContext: dataSourceRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRepositoriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*autocloud_sdk.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	_, err := c.GetRepositories("auto-cloud")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Trace(ctx, "geting the repositories")

	// err = d.Set("id", me.Me.ID)
	if err != nil {
		fmt.Println(err)
		return diag.FromErr(err)

	}
	// err = d.Set("email", me.Me.Email)
	if err != nil {
		fmt.Println(err)
		return diag.FromErr(err)

	}
	// err = d.Set("name", me.Me.Name)
	if err != nil {
		fmt.Println(err)
		return diag.FromErr(err)

	}
	//strconv.FormatInt(time.Now().Unix(), 10)
	// d.SetId(me.Me.ID)

	return diags
}
