package blueprint_config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func DataSourceBlueprintConfig() *schema.Resource {
	setOfStringSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	validationRulesSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rule": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"isRequired", "regex", "minLength", "maxLength"}, false),
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"scope": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "",
					ValidateFunc: validation.StringInSlice([]string{"value", "key"}, false),
				},
				"error_message": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
			},
		},
	}

	optionItemSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"label": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
					ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
						input := v.(string)

						isReference := len(strings.Split(input, ".")) == 3
						if isReference {
							return diag.Diagnostics{}
						}

						isOutput := strings.Contains(input, ".outputs.")
						if isOutput {
							return diag.Diagnostics{}
						}

						var diags diag.Diagnostics
						// this is needed to test if is a valid hcl code
						inputasHcl := "value = " + input

						hclInput, parsingDiag := hclsyntax.ParseConfig([]byte(inputasHcl), "", hcl.Pos{Line: 1, Column: 1})
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Input should be a valid hcl code",
							Detail:   fmt.Sprintf("%q is not a valid hcl code, please make sure you are calling it within Heredoc references, ex: value = <<-HCL %s\nHCL ", input, input),
						}
						if parsingDiag.HasErrors() {
							diags = append(diags, diag)
							return diags
						}
						att, err := hclInput.Body.JustAttributes()
						if err != nil {
							diags = append(diags, diag)
							return diags
						}
						for _, a := range att {
							_, err := a.Expr.Value(nil)
							if err != nil {
								diag.Summary = "Input should be a valid terraform data type"
								diag.Detail = fmt.Sprintf("%q is not a valid terraform data type, variables not allowed. Tip: Are you passing an unquoted string?   ", input)
								diags = append(diags, diag)
								return diags
							}
						}

						return diags
					},
				},
				"checked": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}

	variableSchema := map[string]*schema.Schema{

		//"form_config": formConfigSchema,
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"helper_text": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
				input := v.(string)

				isReference := len(strings.Split(input, ".")) == 3
				if isReference {
					return diag.Diagnostics{}
				}

				isOutput := strings.Contains(input, ".outputs.")
				if isOutput {
					return diag.Diagnostics{}
				}

				var diags diag.Diagnostics
				// this is needed to test if is a valid hcl code
				inputasHcl := "value = " + input

				hclInput, parsingDiag := hclsyntax.ParseConfig([]byte(inputasHcl), "", hcl.Pos{Line: 1, Column: 1})
				diag := diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Input should be a valid hcl code",
					Detail:   fmt.Sprintf("%q is not a valid hcl code, please make sure you are calling it within Heredoc references, ex: value = <<-HCL %s\nHCL ", input, input),
				}
				if parsingDiag.HasErrors() {
					diags = append(diags, diag)
					return diags
				}
				att, err := hclInput.Body.JustAttributes()
				if err != nil {
					diags = append(diags, diag)
					return diags
				}
				for _, a := range att {
					_, err := a.Expr.Value(nil)
					if err != nil {
						diag.Summary = "Input should be a valid terraform data type"
						diag.Detail = fmt.Sprintf("%q is not a valid terraform data type, variables not allowed. Tip: Are you passing an unquoted string?   ", input)
						diags = append(diags, diag)
						return diags
					}
				}

				return diags
			},
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			//ValidateFunc: validation.StringInSlice([]string{SHORTTEXT_TYPE, RADIO_TYPE, CHECKBOX_TYPE, MAP_TYPE, RAW_TYPE, EDITOR_TYPE}, false),
		},
		"options": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"option": optionItemSchema,
				},
			},
		},
		"required_values": {
			Type:     schema.TypeString,
			Optional: true,
		},
		//"conditional":     conditionalSchema,
		"validation_rule": validationRulesSchema,
		"default": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"variables": {
			Description: "A key value map of variables to be used in variable interpolation",
			Type:        schema.TypeMap,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	conditionalSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"source": {
					Type:     schema.TypeString,
					Required: true,
				},
				"condition": {
					Type:     schema.TypeString,
					Required: true,
					ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
						input := v.(string)

						var diags diag.Diagnostics
						// this is needed to test if is a valid hcl code
						inputasHcl := "value = " + input

						hclInput, parsingDiag := hclsyntax.ParseConfig([]byte(inputasHcl), "", hcl.Pos{Line: 1, Column: 1})
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Condition should be a valid hcl code",
							Detail:   fmt.Sprintf("%q is not a valid hcl code, please make sure you are calling it within Heredoc references, ex: value = <<-HCL %s\nHCL ", input, input),
						}
						if parsingDiag.HasErrors() {
							diags = append(diags, diag)
							return diags
						}
						att, err := hclInput.Body.JustAttributes()
						if err != nil {
							diags = append(diags, diag)
							return diags
						}
						for _, a := range att {
							_, err := a.Expr.Value(nil)
							if err != nil {
								diag.Summary = "Condition should be a valid terraform data type"
								diag.Detail = fmt.Sprintf("%q is not a valid terraform data type, variables not allowed. Tip: Are you passing an unquoted string?   ", input)
								diags = append(diags, diag)
								return diags
							}
						}

						return diags
					},
				},
				"content": {
					Type:     schema.TypeSet,
					Required: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: utils.MergeSchemas(variableSchema, map[string]*schema.Schema{}),
					},
				},
			},
		},
	}

	bluePrintConfigSchema := map[string]*schema.Schema{
		"source": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"omit_variables": setOfStringSchema,
		"variable": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: utils.MergeSchemas(variableSchema, map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"conditional": conditionalSchema,
				}),
			},
		},
		"config": { // the form as json to replace the default variables
			Description: "Variables retrieved in the tree",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"blueprint_config": { // the form as json to replace the default variables
			Description: "Processed form variables JSON (to replace the default module variables variables)",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"display_order": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"priority": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"values": {
						Type:     schema.TypeList,
						Required: true,
						MinItems: 0,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"variables": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			}},
	}

	return &schema.Resource{
		Description: "terraform form processor (form builder)",
		ReadContext: dataSourceBlueprintConfigRead,
		Schema:      bluePrintConfigSchema,
	}
}

// main function to read context from terraform
func dataSourceBlueprintConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	blueprintConfig, err := GetBlueprintConfigFromSchema(d)

	if err != nil {
		return diag.FromErr(err)
	}
	formattedBlueprint, err := utils.PrettyStruct(blueprintConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("INPUT BLUEPRINTCONFIG->", formattedBlueprint)

	err = d.Set("blueprint_config", formattedBlueprint)
	if err != nil {
		return diag.FromErr(err)
	}
	//formVariables := []generator.FormShape{}
	// err = validateConditionals(formVariables)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// here will apply overrides, omits and will order the questions by display_order definitions
	processedVars := Transverse(blueprintConfig)
	formattedFormVariables, err := utils.ToJsonString(processedVars)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("config", formattedFormVariables)
	if err != nil {
		return diag.FromErr(err)
	}

	variablesMap := CreateMapFromFormQuestions(processedVars)
	err = d.Set("variables", variablesMap)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(blueprintConfig.Id)
	return diags
}

// Bridges the gap between golang structs and terraform schema
func GetBlueprintConfigFromSchema(d *schema.ResourceData) (*BluePrintConfig, error) {
	id, err := gonanoid.New()
	if err != nil {
		return nil, err
	}
	bp := BluePrintConfig{
		Id:                id,
		Variables:         make([]generator.FormShape, 0),
		OmitVariables:     make([]string, 0),
		OverrideVariables: make(map[string]OverrideVariable, 0),
		Children:          make(map[string]*BluePrintConfig, 0),
	}

	// get sources
	var cerrors []error // collect data errors
	if v, ok := d.GetOk("source"); ok {
		cerrors = append(cerrors, GetBlueprintConfigSources(v, &bp))
	}

	// get omit_variables
	if v, ok := d.GetOk("omit_variables"); ok {
		log.Printf("omit_vars get.ok is ok, %v\n", v)
		cerrors = append(cerrors, GetBlueprintConfigOmitVariables(v.(*schema.Set).List(), &bp))
		log.Printf("the [%v] are the omitted vars", bp.OmitVariables)
	} else {
		log.Printf("omit_vars get.ok not ok, no variables were added\n")
	}

	// get display_order
	if v, ok := d.GetOk("display_order"); ok {
		log.Printf("display_order get.ok is ok, %v\n", v)
		bp.DisplayOrder = DisplayOrder{
			Priority: 0,
			Values:   []string{},
		}
		displayOrder := v.([]interface{})
		if len(displayOrder) == 1 {
			validDisplayOrder := displayOrder[0].(map[string]interface{})
			bp.DisplayOrder.Priority = validDisplayOrder["priority"].(int)
			for _, value := range validDisplayOrder["values"].([]interface{}) {
				strValue := value.(string)
				bp.DisplayOrder.Values = append(bp.DisplayOrder.Values, strValue)
			}
		} else if len(displayOrder) > 1 {
			return nil, errors.New("display_order should only be defined once")
		}

		log.Printf("the %v is the displayOrder", bp.DisplayOrder)
	} else {
		log.Printf("display_order get.ok not ok\n")
	}

	// get override vars
	if v, ok := d.GetOk("variable"); ok {
		cerrors = append(cerrors, GetBlueprintConfigOverrideVariables(v, &bp))
	}

	// Propagate error
	err = hasError(cerrors)
	if err != nil {
		return nil, err
	}

	str, err := json.MarshalIndent(bp, "", "    ")
	if err != nil {
		return nil, errors.New("invalid conversion to BluePrintConfig")
	}

	log.Printf("final bc: %s", string(str))
	return &bp, nil
}

func CreateMapFromFormQuestions(formVariables []generator.FormShape) map[string]string {
	var outputVars = make(map[string]string)
	for _, form := range formVariables {
		questionName := strings.Split(form.ID, ".")
		variableName := questionName[1]
		outputVars[variableName] = form.ID
	}
	return outputVars
}

// func validateConditionals(variables []generator.FormShape) error {
// 	// vars to map
// 	var varsMap = make(map[string]generator.FormShape, len(variables))
// 	for _, variable := range variables {
// 		varsMap[variable.ID] = variable
// 	}

// 	// validate conditionals
// 	for _, variable := range variables {
// 		for _, conditional := range variable.Conditionals {
// 			dependencyVariable, dependecyExist := varsMap[conditional.Source]
// 			if dependecyExist && dependencyVariable.FormQuestion.FieldType != RADIO_TYPE {
// 				return fmt.Errorf("the conditional's source variable can only be of 'radio' type [variable: %v, source variable: %v, source variable type: %v]", variable.ID, conditional.Source, dependencyVariable.FormQuestion.FieldType)
// 			}
// 		}
// 	}

// 	return nil
// }

func hasError(errors []error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
