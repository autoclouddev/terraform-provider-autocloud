package autocloud_provider

import (
	"context"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO: dev convention autocloudModule1 will be renamed to autocloudModule once we integrate all the blueprint changes
func autocloudModule1() *schema.Resource {
	return &schema.Resource{
		Description: "Create an IAC module.",

		CreateContext: autocloudModule1Create,
		ReadContext:   autocloudModule1Read,
		UpdateContext: autocloudModule1Update,
		DeleteContext: autocloudModule1Delete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "name",
				Type:        schema.TypeString,
				Optional:    true,
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
				Description: "variables form shape for this module",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func autocloudModule1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	autocloudModule1Read(ctx, d, meta)
	return diags
}

func autocloudModule1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloudsdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	iacModuleID := d.Id()

	iacModule, err := c.GetModule(iacModuleID)
	if err != nil {
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

	err = d.Set("variables", iacModule.Variables)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func autocloudModule1Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloudsdk.Client)

	updatedIacModule := GetSdkIacModule(d)
	updatedIacModule.ID = d.Id()

	_, err := c.UpdateModule(&updatedIacModule)
	if err != nil {
		return diag.FromErr(err)
	}
	return autocloudModule1Read(ctx, d, meta)
}

func autocloudModule1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloudsdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	iacModuleID := d.Id()

	err := c.DeleteModule(iacModuleID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
