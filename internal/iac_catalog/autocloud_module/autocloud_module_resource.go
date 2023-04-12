package autocloud_module

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var autocloudModuleSchema = map[string]*schema.Schema{
	"id": {
		Description: "id",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"name": {
		Description: "name",
		Type:        schema.TypeString,
		Required:    true,
		ValidateFunc: func(val any, key string) (warns []string, errs []error) {
			if len(val.(string)) == 0 {
				errs = append(errs, fmt.Errorf("name should not be empty"))
			}
			is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(val.(string))
			if !is_alphanumeric {
				errs = append(errs, fmt.Errorf("name should only contain alphanumeric characters, got: %s", val))
			}
			return
		},
	},
	"source": {
		Description: "terraform module source url from registry",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"version": {
		Description: "terraform module source url version from registry",
		Type:        schema.TypeString,
		Optional:    true,
	},

	"variables": {
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"template_config": {
		Description: "Template config",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"config": {
		Description: "Form config",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"tags_variable": {
		Description: "Tags variable name",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"outputs": {
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// encoding in json for now, but it would be nice to avoid that encoding
	"blueprint_config": {
		Description: "blueprint config ",
		Type:        schema.TypeString,
		Computed:    true,
	},
}

func ResourceAutocloudModule() *schema.Resource {
	return &schema.Resource{
		Description: "Create an IAC module.",

		CreateContext: autocloudModuleCreate,
		ReadContext:   autocloudModuleRead,
		UpdateContext: autocloudModuleUpdate,
		DeleteContext: autocloudModuleDelete,

		Schema: autocloudModuleSchema,
	}
}

func autocloudModuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	var diags diag.Diagnostics

	moduleInput := utils.GetSdkIacModuleInput(d)
	c := meta.(*autocloudsdk.Client)
	//o, err := c.CreateModule(&iacModule)

	o, err := c.Module.Create(moduleInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(o.ID)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")
	// this will populate the state after creating a new resource
	autocloudModuleRead(ctx, d, meta)
	return diags
}

func autocloudModuleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloudsdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	iacModuleID := d.Id()

	iacModule, err := c.Module.Get(iacModuleID)
	if err != nil {
		resp := autocloudsdk.GetSdkHttpError(err)
		if resp != nil {
			if resp.Status == 400 {
				return diag.Errorf(resp.Message)
			}
		}

		return diag.FromErr(err)
	}

	err = d.Set("name", iacModule.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("source", iacModule.Source)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("version", iacModule.Version)
	if err != nil {
		return diag.FromErr(err)
	}

	varsMap, err := utils.GetVariablesIdMap(iacModule.Variables)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("variables", varsMap)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("config", iacModule.Variables)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("tags_variable", iacModule.TagsVariable)
	if err != nil {
		return diag.FromErr(err)
	}

	outputsMap, err := utils.ToStringMap(iacModule.Outputs)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("outputs", outputsMap)
	if err != nil {
		return diag.FromErr(err)
	}
	variables := []generator.FormShape{}
	err = json.Unmarshal([]byte(iacModule.Variables), &variables)
	if err != nil {
		return diag.FromErr(err)
	}

	// populate vars' module id
	for i, v := range variables {
		if v.ModuleID == "" {
			variables[i].ModuleID = iacModuleID
		}
	}
	config := blueprint_config.BluePrintConfig{
		Id:        iacModule.ID,
		Variables: variables,
		Children:  make(map[string]blueprint_config.BluePrintConfig),
		DisplayOrder: blueprint_config.DisplayOrder{
			Priority: 1000,
			Values:   make([]string, 0),
		},
	}
	jsonconf, err := json.Marshal(config)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("blueprint_config", string(jsonconf))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func autocloudModuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloudsdk.Client)

	updatedIacModule := utils.GetSdkIacModuleInput(d)
	updatedIacModule.ID = d.Id()

	_, err := c.Module.Update(updatedIacModule)
	if err != nil {
		resp := autocloudsdk.GetSdkHttpError(err)
		if resp != nil {
			if resp.Status == 400 {
				return diag.Errorf(resp.Message)
			}
		}

		return diag.FromErr(err)
	}
	return autocloudModuleRead(ctx, d, meta)
}

func autocloudModuleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloudsdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	iacModuleID := d.Id()

	err := c.Module.Delete(iacModuleID)
	if err != nil {
		resp := autocloudsdk.GetSdkHttpError(err)
		if resp != nil {
			if resp.Status == 400 {
				return diag.Errorf(resp.Message)
			}
		}
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
