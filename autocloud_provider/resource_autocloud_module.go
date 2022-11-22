package autocloud_provider

import (
	"context"
	"fmt"
	"regexp"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"

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
	"template": {
		Description: "tf source code from registry",
		Type:        schema.TypeString,
		Computed:    true,
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
	"form_config": {
		Description: "Form config",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"tags_variable": {
		Description: "Tags variable name",
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "tags",
	},
}

func autocloudModule() *schema.Resource {
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

	iacModule := GetSdkIacModule(d)
	c := meta.(*autocloudsdk.Client)
	o, err := c.CreateModule(&iacModule)
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

	iacModule, err := c.GetModule(iacModuleID)
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

	err = d.Set("template", iacModule.Template)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("template_config", iacModule.Template)
	if err != nil {
		return diag.FromErr(err)
	}

	varsMap, err := GetVariablesIdMap(iacModule.Variables)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("variables", varsMap)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("form_config", iacModule.Variables)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("tags_variable", iacModule.TagsVariable)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func autocloudModuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloudsdk.Client)

	updatedIacModule := GetSdkIacModule(d)
	updatedIacModule.ID = d.Id()

	_, err := c.UpdateModule(&updatedIacModule)
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

	err := c.DeleteModule(iacModuleID)
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
