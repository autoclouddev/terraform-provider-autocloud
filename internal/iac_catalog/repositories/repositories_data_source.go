package repositories

import (
	"context"
	"strconv"
	"time"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	gitrepository "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/git_repository"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositories() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "github_repos data source, to check the granted repositories.",

		ReadContext: dataSourceRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenData(repositories *[]gitrepository.GitRepository) []interface{} {
	if repositories != nil {
		data := make([]interface{}, len(*repositories))

		for i, repository := range *repositories {
			repo := make(map[string]interface{})

			repo["id"] = repository.ID
			repo["name"] = repository.Name
			repo["url"] = repository.HTMLURL
			repo["description"] = repository.Description
			data[i] = repo
		}

		return data
	}

	return make([]interface{}, 0)
}

func dataSourceRepositoriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*autocloudsdk.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	repositories, err := c.GitRepository.List()

	if err != nil {
		resp := autocloudsdk.GetSdkHttpError(err)
		if resp != nil {
			if resp.Status == 400 {
				return diag.Errorf(resp.Message)
			}
		}

		return diag.FromErr(err)
	}

	tflog.Trace(ctx, "getting the repositories")

	data := flattenData(repositories)
	err = d.Set("data", data)
	if err != nil {
		//fmt.Println(err)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
