package autocloud_provider

import (
	"autocloud_sdk"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RepositoriesSource struct {
	id          int
	name        string
	url         string
	description string
}

func dataSourceRepositories() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "github_repos data source, to check the granted repositories.",

		ReadContext: dataSourceRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"data": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenData(repositories *[]autocloud_sdk.Repository) []interface{} {
	if repositories != nil {
		data := make([]interface{}, len(*repositories), len(*repositories))

		for i, repository := range *repositories {
			repo := make(map[string]interface{})

			repo["id"] = repository.ID
			repo["name"] = repository.Name
			repo["url"] = repository.HtmlUrl
			repo["description"] = repository.Description
			data[i] = repo
		}

		return data
	}

	return make([]interface{}, 0)
}

func dataSourceRepositoriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*autocloud_sdk.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	repositories, err := c.GetRepositories()
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Trace(ctx, "geting the repositories")

	data := flattenData(&repositories)
	d.Set("data", data)
	if err != nil {
		fmt.Println(err)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
