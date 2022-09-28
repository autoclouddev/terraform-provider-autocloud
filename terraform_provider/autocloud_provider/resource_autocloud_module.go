package autocloud_provider

import (
	"autocloud_sdk"
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func autocloudModule() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Create an IAC generator.",

		CreateContext: autocloudModuleCreate,
		ReadContext:   autocloudModuleRead,
		UpdateContext: autocloudModuleUpdate,
		DeleteContext: autocloudModuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"author": {
				Description: "author",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"slug": {
				Description: "slug",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": {
				Description: "description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instructions": {
				Description: "instructions",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"labels": {
				Description: "labels",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func autocloudModuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	var diags diag.Diagnostics

	var labels = []string{}
	if v, ok := d.GetOk("labels"); ok {
		list := v.([]interface{})
		labels = make([]string, len(list))
		for i, v := range list {
			labels[i] = v.(string)
		}
	}
	c := meta.(*autocloud_sdk.Client)
	generator := autocloud_sdk.IacCatalog{
		Name:         d.Get("name").(string),
		Author:       d.Get("author").(string),
		Slug:         d.Get("slug").(string),
		Description:  d.Get("description").(string),
		Instructions: d.Get("instructions").(string),
		Labels:       labels,
	}
	o, err := c.CreateGenerator(generator)
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
	c := meta.(*autocloud_sdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	generatorID := d.Id()

	generator, err := c.GetGenerator(generatorID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", generator.Name)
	d.Set("author", generator.Author)
	d.Set("slug", generator.Slug)
	d.Set("description", generator.Description)
	d.Set("instructions", generator.Instructions)
	d.Set("labels", generator.Labels)

	return diags
}

func autocloudModuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloud_sdk.Client)
	generatorID := d.Id()

	var labels = []string{}
	if v, ok := d.GetOk("labels"); ok {
		list := v.([]interface{})
		labels = make([]string, len(list))
		for i, v := range list {
			labels[i] = v.(string)
		}
	}

	updatedGen := autocloud_sdk.IacCatalog{
		ID:           generatorID,
		Name:         d.Get("name").(string),
		Author:       d.Get("author").(string),
		Slug:         d.Get("slug").(string),
		Description:  d.Get("description").(string),
		Instructions: d.Get("instructions").(string),
		Labels:       labels,
	}
	_, err := c.UpdateGenerator(updatedGen)
	if err != nil {
		return diag.FromErr(err)
	}
	return autocloudModuleRead(ctx, d, meta)
}

func autocloudModuleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloud_sdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	generatorID := d.Id()

	err := c.DeleteGenerator(generatorID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
