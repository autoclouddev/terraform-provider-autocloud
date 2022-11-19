package autocloud_provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

func autocloudBlueprint() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "main resource to create an IAC blueprint.",

		CreateContext: autocloudBlueprintCreate,
		ReadContext:   autocloudBlueprintRead,
		UpdateContext: autocloudBlueprintUpdate,
		DeleteContext: autocloudBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Name of the blueprint.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"author": {
				Description: "author",
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
						"git_url_options": {
							Description: "git_url_options",
							Type:        schema.TypeList,
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"git_url_default": {
							Description: "git_url_default",
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
				},
			},
			"file": {
				Description: "file",
				Type:        schema.TypeSet,
				Required:    true,
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
						"modules": {
							Description: "modules, array containing the names of the modules included in this file",
							Type:        schema.TypeList,
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"autocloud_module": {
				Description: "autocloud_module",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "autocloud module id",
							Type:        schema.TypeString,
							Required:    true,
						},
						"form_config": {
							Description: "form config",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"template_config": {
							Description: "template config",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func autocloudBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	generator := GetSdkIacCatalog(d)
	c := meta.(*autocloudsdk.Client)
	fmt.Println("SENDING DATA, CREATE")
	fmt.Println(generator)
	o, err := c.CreateGenerator(generator)
	if err != nil {
		resp := autocloudsdk.GetSdkHttpError(err)
		if resp != nil {
			if resp.Status == 400 {
				return diag.Errorf(resp.Message)
			}
		}
		return diag.FromErr(err)
	}

	d.SetId(o.ID)
	tflog.Trace(ctx, "created a resource")
	autocloudBlueprintRead(ctx, d, meta)
	return diags
}

func autocloudBlueprintRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*autocloudsdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	generatorID := d.Id()

	generator, err := c.GetGenerator(generatorID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", generator.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("author", generator.Author)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("description", generator.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("instructions", generator.Instructions)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("labels", generator.Labels)
	if err != nil {
		return diag.FromErr(err)
	}

	files := lowercaseFileDefs(generator.FileDefinitions)
	err = d.Set("file", files)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func autocloudBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*autocloudsdk.Client)

	updatedGen := GetSdkIacCatalog(d)
	updatedGen.ID = d.Id()
	fmt.Println("UPDATING GENERATOR REQUEST")
	fmt.Println(updatedGen)
	_, err := c.UpdateGenerator(updatedGen)
	if err != nil {
		resp := autocloudsdk.GetSdkHttpError(err)
		if resp != nil {
			if resp.Status == 400 {
				return diag.Errorf(resp.Message)
			}
		}

		return diag.FromErr(err)
	}
	return autocloudBlueprintRead(ctx, d, meta)
}

func autocloudBlueprintDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*autocloudsdk.Client)

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

func lowercaseFileDefs(files []autocloudsdk.IacCatalogFile) []interface{} {
	var out = make([]interface{}, 0)
	for _, file := range files {
		m := make(map[string]interface{})
		m["action"] = file.Action
		m["path_from_root"] = file.PathFromRoot
		m["filename_template"] = file.FilenameTemplate
		m["filename_vars"] = file.FilenameVars
		m["modules"] = file.Modules
		out = append(out, m)
	}

	return out
}
