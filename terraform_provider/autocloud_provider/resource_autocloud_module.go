package autocloud_provider

import (
	"autocloud_sdk"
	"context"
	"fmt"

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
			"git_config": {
				Description: "git_config",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_branch": {
							Description: "destination_branch",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"pull_request": {
							Description: "pull_request",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": {
										Description: "title",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"commit_message_template": {
										Description: "commit_message_template",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"body": {
										Description: "body",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"variables": {
										Description: "variables",
										Type:        schema.TypeSet,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"author_name": {
													Description: "author_name",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"cluster_name": {
													Description: "cluster_name",
													Type:        schema.TypeString,
													Optional:    true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"file": {
				Description: "file",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Description: "action",
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								isValidAction := Contains([]string{"CREATE", "EDIT", "HCLEDIT"}, val.(string))
								if !isValidAction {
									errs = append(errs, fmt.Errorf("%q must be a value in [CREATE, EDIT, HCLEDIT], got: %s", key, val))
								}
								return
							},
						},
						"path_from_root": {
							Description: "path_from_root",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"filename_template": {
							Description: "filename_template",
							Type:        schema.TypeString,
							Required:    true,
						},
						"filename_vars": {
							Description: "filename_vars",
							Type:        schema.TypeMap,
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func autocloudModuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	var diags diag.Diagnostics

	generator := GetSdkIacCatalog(d)
	c := meta.(*autocloud_sdk.Client)
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
	d.Set("fileDefinitions", generator.FileDefinitions)

	return diags
}

func autocloudModuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*autocloud_sdk.Client)

	updatedGen := GetSdkIacCatalog(d)
	updatedGen.ID = d.Id()

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
