package blueprint_config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
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
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{SHORTTEXT_TYPE, RADIO_TYPE, CHECKBOX_TYPE, MAP_TYPE, RAW_TYPE, EDITOR_TYPE}, false),
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
					"variables": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
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
						MinItems: 1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"references": {
			Type:     schema.TypeString,
			Computed: true,
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

func dataSourceBlueprintConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	blueprintConfig, err := GetBlueprintConfigFromSchema(d)

	if err != nil {
		return diag.FromErr(err)
	}
	pretty, err := utils.PrettyStruct(blueprintConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("INPUT BLUEPRINTCONFIG->", pretty)

	// Save references
	aliases := blueprint_config_references.GetInstance()

	err = d.Set("references", aliases.ToString())
	if err != nil {
		return diag.FromErr(err)
	}

	// new form variables (as JSON)
	formVariables, err := GetFormShape(*blueprintConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	prioritiesDuplicated, _ := getDisplayOrderDuplicated(*blueprintConfig)
	if len(prioritiesDuplicated) > 0 {
		diags = append(diags, diag.Diagnostic{
			Detail:   prioritiesDuplicated,
			Severity: diag.Warning,
			Summary:  "Display order priorities duplicated",
		})
	}
	err = validateConditionals(formVariables)
	if err != nil {
		return diag.FromErr(err)
	}
	vars := GetVariablesInBlueprint(formVariables)
	err = d.Set("variables", vars)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonFormShape, err := utils.ToJsonString(formVariables)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("config", jsonFormShape)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("blueprint_config", pretty)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(blueprintConfig.Id)
	return diags
}

// maps tf declaration to object
func GetBlueprintConfigFromSchema(d *schema.ResourceData) (*BluePrintConfig, error) {
	bp := BluePrintConfig{}
	id, err := gonanoid.New()
	if err != nil {
		return nil, err
	}
	//if id is not set, anyways the id is always generated, it is not saved in the statefile between executions
	if len(d.Id()) == 0 {
		bp.Id = id
	} else {
		bp.Id = d.Id()
	}

	bp.OverrideVariables = make(map[string]OverrideVariable, 0)
	aliasToModuleNameMap := blueprint_config_references.GetInstance()

	// get sources
	var cerrors []error // collect data errors
	if v, ok := d.GetOk("source"); ok {
		cerrors = append(cerrors, GetBlueprintConfigSources(v, &bp, *aliasToModuleNameMap))
	}

	// get omit_variables
	if v, ok := d.GetOk("omit_variables"); ok {
		log.Printf("omit_vars get.ok is ok, %v\n", v)
		cerrors = append(cerrors, GetBlueprintConfigOmitVariables(v.(*schema.Set).List(), &bp, *aliasToModuleNameMap))
		log.Printf("the [%v] are the omitted vars", bp.OmitVariables)
	} else {
		log.Printf("omit_vars get.ok not ok, no variables were added\n")
	}

	// get display_order
	if v, ok := d.GetOk("display_order"); ok {
		log.Printf("display_order get.ok is ok, %v\n", v)
		cerrors = append(cerrors, GetBlueprintConfigDisplayOrder(v, &bp, *aliasToModuleNameMap))
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

func GetVariablesInBlueprint(formVariables []generator.FormShape) map[string]string {
	var outputVars = make(map[string]string)
	for _, form := range formVariables {
		questionName := strings.Split(form.ID, ".")
		variableName := questionName[1]
		outputVars[variableName] = form.ID
	}
	return outputVars
}

func validateConditionals(variables []generator.FormShape) error {
	// vars to map
	var varsMap = make(map[string]generator.FormShape, len(variables))
	for _, variable := range variables {
		varsMap[variable.ID] = variable
	}

	// validate conditionals
	for _, variable := range variables {
		for _, conditional := range variable.Conditionals {
			dependencyVariable, dependecyExist := varsMap[conditional.Source]
			if dependecyExist && dependencyVariable.FormQuestion.FieldType != RADIO_TYPE {
				return fmt.Errorf("the conditional's source variable can only be of 'radio' type [variable: %v, source variable: %v, source variable type: %v]", variable.ID, conditional.Source, dependencyVariable.FormQuestion.FieldType)
			}
		}
	}

	return nil
}

func hasError(errors []error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

func getDisplayOrderDuplicated(root BluePrintConfig) (string, error) {
	order, err := postDisplayOrderTransversal(&root)
	if err != nil {
		return "", err
	}
	var displayOrderByPriority = make(map[int][]string, 0)
	var prioritiesDuplicatedMessage = ""
	// we build a map of displayOrder by priority
	for _, displayOrder := range order {
		priority := displayOrder.Priority
		if len(displayOrder.Values) > 0 {
			// to do more pretty the message
			formatMessage := "\n%+v"
			if len(displayOrderByPriority[priority]) > 0 {
				formatMessage = "\n%+v\n"
			}
			displayOrderByPriority[priority] = append(displayOrderByPriority[priority], fmt.Sprintf(formatMessage, displayOrder))
		}
	}
	// then we add to the duplicated list those which has more than one entry
	for _, displayOrders := range displayOrderByPriority {
		if len(displayOrders) > 1 {
			prioritiesDuplicatedMessage = fmt.Sprintf("%v\n%v", prioritiesDuplicatedMessage, displayOrders)
		}
	}
	return prioritiesDuplicatedMessage, nil
}
