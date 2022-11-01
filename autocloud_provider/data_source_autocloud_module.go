package autocloud_provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

var dataFilter = map[string]*schema.Schema{
	"filter": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"version": {
					Type:     schema.TypeString,
					Required: true,
				},

				"source": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
}

func dataSourceAutocloudModule() *schema.Resource {
	return &schema.Resource{
		Description: "autocloud module, this resource fetches a previously created module",
		ReadContext: dataSourceAutocloudModuleRead,
		Schema:      mergeSchemas(autocloudModuleSchema, dataFilter),
	}
}

func dataSourceAutocloudModuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*autocloudsdk.Client)
	var diags diag.Diagnostics
	name, version, source := "", "", ""
	if filter, ok := d.GetOk("filter"); ok {
		list := filter.(*schema.Set).List()
		for _, f := range list {
			var filterMap = f.(map[string]interface{})
			name = filterMap["name"].(string)
			version = filterMap["version"].(string)
			source = filterMap["source"].(string)
		}
		module, err := c.GetModuleByParams(name, version, source)
		fmt.Println(module)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(module.ID)
		err = d.Set("name", module.Name)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("version", module.Version)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("source", module.Source)
		if err != nil {
			return diag.FromErr(err)
		}
		varsMap, err := GetVariablesIdMap(module.Variables)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("variables", varsMap)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("form_config", module.Variables)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("template", module.Template)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("template_config", module.Template)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(errors.New("no fetched module"))
	}
	return diags
}
