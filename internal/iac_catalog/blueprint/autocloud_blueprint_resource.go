package blueprint

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func ResourceAutocloudBlueprint() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "main resource to create an IAC blueprint.",

		CreateContext: autocloudBlueprintCreate,
		ReadContext:   autocloudBlueprintRead,
		UpdateContext: autocloudBlueprintUpdate,
		DeleteContext: autocloudBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "id",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Name of the blueprint, this will be used in the autocloud app",
				Type:        schema.TypeString,
				Required:    true,
			},
			"author": {
				Description: "Name of the author of the blueprint",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "What does this blueprint do?",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instructions": {
				Description: "Use this space to instruct the user on how to use this blueprint, markdown is supported",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"labels": {
				Description: "Place a valid cloud platform label here, valid options are AWS, GCP and AZURE",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"git_config": {
				Description: "THis block contains the git configuration for the destination repository",
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_branch": {
							Description: "To where the Pull Request will be sent, this branch should exist in the destination repository",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"git_url_options": {
							Description: "A list of valid git urls remotes. We currently support only Github",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"git_url_default": {
							Description: "The default Github url to use when the user does not provide one",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"pull_request": {
							Description: "Configuration of the pull request to be created",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": {
										Description: "Title of the PR, it supports template using {{notation}}",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"commit_message_template": {
										Description: "Commit message of the PR, it supports template using {{notation}}",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"body": {
										Description: "Body of the PR, it supports template using {{notation}}",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"variables": {
										Description: "A key value map of variables to be used in the templates",
										Type:        schema.TypeMap,
										Optional:    true,
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
				Description: "A list of files to be created in the destination repository, at least 1 is required",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Description: "Currently we support only CREATE",
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								isValidAction := true //Contains([]string{"CREATE", "EDIT", "HCLEDIT"}, val.(string))
								if !isValidAction {
									errs = append(errs, fmt.Errorf("%q must be a value in [CREATE, EDIT, HCLEDIT], got: %s", key, val))
								}
								return
							},
						},
						"destination": {
							Description: "A valid path in the destination repository, it should include the file name, it supports template using {{notation}}",
							Type:        schema.TypeString,
							Required:    true,
						},
						"variables": {
							Description: "A key value map of variables to be used in the templates",
							Type:        schema.TypeMap,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"modules": {
							Description: "An array containing the names of the modules included in this file, please use the name of the module as defined in the modules block or resource.autocloud_module.<module>.name",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"content": {
							Description: "You can use this field to provide the content of the file, it supports template using {{notation}}, if this is set, the modules field will be ignored",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"header": {
							Description: "You can add a header to the file, it supports template using {{notation}}, can be used with module or content",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"footer": {
							Description: "You can add a footer to the file, it supports template using {{notation}}, can be used with module or content",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"config": {
				Description: "A valid output from data.autocloud_blueprint_config.config",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func autocloudBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	generator, err := GetSdkIacCatalog(d)
	if err != nil {
		return diag.Errorf("Failed to get generator: %s", err.Error())
	}
	c := meta.(*autocloudsdk.Client)
	log.Printf("sending data: %v ", generator)
	o, err := c.IacGenerator.Create(*generator)
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
	_, err := GetSdkIacCatalog(d)
	if err != nil {
		return diag.Errorf("Failed to get generator: %s", err.Error())
	}

	c := meta.(*autocloudsdk.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	generatorID := d.Id()

	generator, err := c.IacGenerator.Get(generatorID)
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

	files, diags := lowercaseFileDefs(generator.FileDefinitions, diags)

	err = d.Set("file", files)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("generator.FormQuestions %v\n", generator.FormQuestions)
	formQuestions, err := json.Marshal(generator.FormQuestions)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("config", string(formQuestions))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func autocloudBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*autocloudsdk.Client)

	updatedGen, err := GetSdkIacCatalog(d)
	if err != nil {
		return diag.Errorf(err.Error())
	}
	updatedGen.ID = d.Id()
	log.Println("UPDATING GENERATOR REQUEST")
	log.Println(updatedGen)
	_, err = c.IacGenerator.Update(*updatedGen)
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

	err := c.IacGenerator.Delete(generatorID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func lowercaseFileDefs(files []generator.IacCatalogFile, diags diag.Diagnostics) ([]interface{}, diag.Diagnostics) {
	var out = make([]interface{}, 0)
	for _, file := range files {
		m := make(map[string]interface{})
		m["action"] = file.Action
		m["destination"] = file.Destination
		m["variables"] = file.Variables
		m["modules"] = file.Modules
		m["content"] = file.Content
		m["footer"] = file.Header
		m["header"] = file.Footer
		out = append(out, m)

		if file.Content != "" && (file.Footer != "" || file.Header != "" || len(file.Modules) > 0) {
			diags = append(diags, diag.Diagnostic{
				Detail:   "footer, header, or module properties defined on a file block will be omitted if the content property is defined",
				Severity: diag.Warning,
				Summary:  "content property overrides any other attribute",
			})
		}
	}

	return out, diags
}

func modulesBelongToConfig(generator *generator.IacCatalog) bool {
	files := generator.FileDefinitions
	// using a map as a set
	modulesInConfig := map[string]struct{}{}

	for _, question := range generator.FormQuestions {
		moduleName := strings.Split(question.ID, ".")[0]
		modulesInConfig[moduleName] = struct{}{}
	}
	fmt.Println("modulesInConfig", modulesInConfig)
	for _, file := range files {
		fileModules := file.Modules
		fmt.Println("fileModules", fileModules)
		if len(fileModules) == 0 {
			continue
		}
		for _, module := range fileModules {
			if _, ok := modulesInConfig[module]; !ok {
				return false
			}
		}
	}
	return true
}

// Reads the resource schema, validate inputs and returns a IacCatalog object
// if the inputs are invalid, it returns an error
// use this function before creating or updating a generator
func GetSdkIacCatalog(d *schema.ResourceData) (*generator.IacCatalog, error) {
	var labels = []string{}
	if labelValues, isLabelValuesOk := d.GetOk("labels"); isLabelValuesOk {
		list := labelValues.([]interface{})
		labels = utils.ToStringSlice(list)
	}
	var formShape []generator.FormShape
	if config, configExist := d.GetOk("config"); configExist {
		configstr := config.(string)
		err := json.Unmarshal([]byte(configstr), &formShape)
		if err != nil {
			log.Printf("error: %v", err)
			return nil, fmt.Errorf("Incompatible config format, are you using a valid .config attribute?")
		}
	}

	gc := utils.GetSdkIacCatalogGitConfig(d)

	fileDef, err := utils.GetSdkIacCatalogFileDefinitions(d)
	if err != nil {
		return nil, err
	}

	// TODO: convert tree to array
	// read from leaves to root all variables and make a huge array of variables
	// process overrides and conditionals
	generator := generator.IacCatalog{
		Name:            d.Get("name").(string),
		Author:          d.Get("author").(string),
		Description:     d.Get("description").(string),
		Instructions:    d.Get("instructions").(string),
		Labels:          labels,
		FileDefinitions: fileDef,
		GitConfig:       gc,
		FormQuestions:   formShape,
	}

	if !modulesBelongToConfig(&generator) {
		return nil, fmt.Errorf("modules in file definitions do not belong to config")
	}

	return &generator, nil
}
