package blueprint_config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
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
					ValidateFunc: validation.StringInSlice([]string{"isRequired", "regex"}, false),
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
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
				"type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"content": {
					Type:     schema.TypeSet,
					Required: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type:     schema.TypeSet,
								Optional: true,
								MinItems: 1,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"option": optionItemSchema,
									},
								},
							},
							"static": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
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
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
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
						ValidateFunc: validation.StringInSlice([]string{SHORTTEXT_TYPE, RADIO_TYPE, CHECKBOX_TYPE, MAP_TYPE}, false),
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
						Description: "required_values",
						Type:        schema.TypeMap,
						Optional:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"conditional":     conditionalSchema,
					"validation_rule": validationRulesSchema,
				},
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
	}

	return &schema.Resource{
		Description: "terraform form processor (form builder)",
		ReadContext: dataSourceBlueprintConfigRead,
		Schema:      bluePrintConfigSchema,
	}
}

func dataSourceBlueprintConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// map the resource to a FormBuilder object
	blueprintConfig, err := GetBlueprintConfigFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	pretty, err := utils.PrettyStruct(blueprintConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("INPUT BLUEPRINTCONFIG->", pretty)
	// new form variables (as JSON)
	formVariables := GetFormShape(*blueprintConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	err = validateConditionals(formVariables)
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

	/*bpOutput, err := utils.ToJsonString(blueprintConfig)
	//log.Println(bpOutput)
	if err != nil {
		return diag.FromErr(err)
	}*/

	err = d.Set("blueprint_config", pretty)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(blueprintConfig.Id)

	return diags
}

func ConvertMap(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)

	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}

	return mapString
}

// maps tf declaration to object
func GetBlueprintConfigFromSchema(d *schema.ResourceData) (*BluePrintConfig, error) {
	bp := &BluePrintConfig{}
	bp.Id = strconv.FormatInt(time.Now().Unix(), 10)
	bp.OverrideVariables = make(map[string]OverrideVariable, 0)
	if v, ok := d.GetOk("source"); ok {
		mapString := make(map[string]BluePrintConfig)
		for key, value := range v.(map[string]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			log.Printf("SOURCE_INPUT_KEY: %v/n", key)
			formattedValue, _ := utils.PrettyString(strValue)
			log.Printf("SOURCE_INPUT_VALUE: %v", formattedValue)
			bc := BluePrintConfig{}
			err := json.Unmarshal([]byte(strValue), &bc)
			if err != nil {
				return nil, errors.New("invalid conversion to BluePrintConfig")
			}
			mapString[strKey] = bc
			bp.Children = append(bp.Children, bc)
		}
	}
	if v, ok := d.GetOk("omit_variables"); ok {
		log.Printf("omit_vars get.ok is ok, %v\n", v)
		list := v.(*schema.Set).List()
		omit := make([]string, len(list))
		for i, optionValue := range list {
			omit[i] = optionValue.(string)
		}
		bp.OmitVariables = omit
		log.Printf("the [%v] are the omitted vars", bp.OmitVariables)
	} else {
		log.Printf("omit_vars get.ok not ok, no variables were added\n")
	}

	if v, ok := d.GetOk("variable"); ok {
		varsList := v.([]interface{})
		log.Println(varsList...)
		// override vars loop
		for _, currentVar := range varsList {
			varOverrideMap := currentVar.(map[string]interface{})
			varName := varOverrideMap["name"].(string)
			bp.OverrideVariables[varName] = OverrideVariable{
				VariableName: varName,
				DisplayName:  varOverrideMap["display_name"].(string),
				HelperText:   varOverrideMap["helper_text"].(string),
			}
			// Note: if it has a value, then it can NOT have form options "options"
			valueIsDefined := false
			value, valueExist := varOverrideMap["value"]
			valueStr, valueIsString := value.(string)
			valueIsDefined = valueStr != "" // NOTE: if the value is empty, we consider it as 'not defined'
			if valueExist && valueIsString && valueIsDefined {
				// refactor op, if entry, ok could be a function
				if entry, ok := bp.OverrideVariables[varName]; ok {
					entry.Value = valueStr

					bp.OverrideVariables[varName] = entry
				} else {
					return nil, errors.New("cant define blueprint config")
				}
				continue // no options or validation rules required, validation rules are for webform
			}

			optionsFromSchema := varOverrideMap["options"].(*schema.Set)
			//conditionalsInput, conditionalsExist := varOverrideMap["conditionals"]
			/*
				// I think this code is not reacheable
				if valueIsDefined {
					if len(optionsFromSchema.List()) != 0 {
						return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w Var name: [%s], var value [%v]", ErrSetValueInForm, varName, value)
					}

					// if the value is set and there's no form options, then we're ok to continue processing the next variable override
					continue
				}
			*/
			if len(optionsFromSchema.List()) > 1 {
				// it should be caught at schema check level - adding the check here to enforce it in case the schema changes
				return nil, errors.New("exactly one \"options\" must be defined")
			}
			/*
				// I think this conditional is wrong
				conditionalsDefined := conditionalsExist && len(conditionalsInput.(*schema.Set).List()) > 0
				if len(optionsFromSchema.List()) == 0 {

					if conditionalsDefined {
						// if we don't have any conditionals, we should have at least 1 form options
						return nil, fmt.Errorf("a options must be defined for the variable [%s]", varName)
					}
				}
			*/
			rawVariableType := varOverrideMap["type"]
			variableType := rawVariableType.(string)

			if variableType == "shortText" && len(optionsFromSchema.List()) > 0 {
				return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrShortTextCantHaveOptions)
			}

			if entry, ok := bp.OverrideVariables[varName]; ok {
				entry.FormConfig = FormConfig{
					Type:            variableType,
					ValidationRules: make([]ValidationRule, 0),
					FieldOptions:    make([]FieldOption, 0),
				}
				bp.OverrideVariables[varName] = entry
			} else {
				return nil, errors.New("GetBlueprintConfigFromSchema: Error accessing bp")
			}

			if variableType == RADIO_TYPE || variableType == CHECKBOX_TYPE {
				if len(optionsFromSchema.List()) != 1 {
					return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrOneBlockOptionsRequied)
				}
				rawOptionsCluster := optionsFromSchema.List()[0] // "options key in schema options {}"
				rawOptions := rawOptionsCluster.(map[string]interface{})

				optionSchema := rawOptions["option"].(*schema.Set)

				for _, option := range optionSchema.List() {
					options := option.(map[string]interface{})
					fieldOption := FieldOption{
						Label:   options["label"].(string),
						Value:   options["value"].(string),
						Checked: options["checked"].(bool),
					}
					// refactor op, if entry, ok could be a function
					if entry, ok := bp.OverrideVariables[varName]; ok {
						entry.FormConfig.FieldOptions = append(entry.FormConfig.FieldOptions, fieldOption)
						bp.OverrideVariables[varName] = entry
					} else {
						return nil, errors.New("GetBlueprintConfigFromSchema: Error accessing bp")
					}
				}
			}

			if variableType == MAP_TYPE {
				if val, ok := varOverrideMap["required_values"]; ok {
					var variablesMap = val.(map[string]interface{})
					ccc := ConvertMap(variablesMap)

					pairs := make([]autocloudsdk.KeyValue, 0)

					for key, value := range ccc {
						pair := autocloudsdk.KeyValue{Key: key, Value: value}
						pairs = append(pairs, pair)
					}
					mapValue, err := json.Marshal(pairs)
					if err != nil {
						fmt.Println(err)
					}

					if entry, ok := bp.OverrideVariables[varName]; ok {
						entry.Value = string(mapValue)

						bp.OverrideVariables[varName] = entry
					} else {
						return nil, errors.New("cant define blueprint config")
					}
				}
			}

			// validation rules
			validationRulesList := varOverrideMap["validation_rule"].(*schema.Set).List()

			for _, validationRule := range validationRulesList {
				validationRuleMap := validationRule.(map[string]interface{})

				rule := validationRuleMap["rule"].(string)
				ruleValue := validationRuleMap["value"].(string)

				if rule == "isRequired" && ruleValue != "" {
					return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrIsRequiredCantHaveValue)
				}
				vr := ValidationRule{
					Rule:         rule,
					Value:        ruleValue,
					ErrorMessage: validationRuleMap["error_message"].(string),
				}

				if entry, ok := bp.OverrideVariables[varName]; ok {
					entry.FormConfig.ValidationRules = append(entry.FormConfig.ValidationRules, vr)
					bp.OverrideVariables[varName] = entry
				} else {
					return nil, errors.New("GetBlueprintConfigFromSchema: Error accessing bp")
				}
			}
			// Conditionals

			conditionals, conditionalExists := varOverrideMap["conditional"].(*schema.Set)
			log.Printf("CONDITIONALS: %v \n", conditionals)
			if conditionalExists {
				if entry, ok := bp.OverrideVariables[varName]; ok {
					entry.Conditionals = getConditionals(conditionals)
					bp.OverrideVariables[varName] = entry
				} else {
					return nil, errors.New("GetBlueprintConfigFromSchema: Error accessing bp")
				}
			}
		}
	}
	str, err := json.MarshalIndent(bp, "", "    ")
	if err != nil {
		return nil, errors.New("invalid conversion to BluePrintConfig")
	}
	log.Printf("final bc: %s", string(str))
	return bp, nil
}

func validateConditionals(variables []autocloudsdk.FormShape) error {
	// vars to map
	var varsMap = make(map[string]autocloudsdk.FormShape, len(variables))
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

func getConditionals(varOverrideMap *schema.Set) []ConditionalConfig {
	conditionalLen := varOverrideMap.Len()
	log.Printf("conditionalLen %v, \n", conditionalLen)

	conditionals := make([]ConditionalConfig, 0) //make([]ConditionalConfig, len(conditionalsList))
	if conditionalLen == 0 {
		return conditionals
	}
	conditionalsList := varOverrideMap.List()
	log.Printf("conditionalList: %v \n", conditionalsList...)

	for _, conditional := range conditionalsList {
		log.Printf("conditional from list: %v \n", conditional)
		conditionalMap := conditional.(map[string]interface{})
		log.Printf("conditionalMap from list: %v \n", conditionalMap)
		conditionalContentMap := conditionalMap["content"].(*schema.Set).List()[0].(map[string]interface{})
		log.Printf("conditionalContentMap: %v \n", conditionalContentMap)

		// length validated at schema level
		var fieldOptions []FieldOption = make([]FieldOption, 0)
		var staticValue *string

		if staticVal, ok := conditionalContentMap["static"]; ok && staticVal != "" {
			str, castOk := staticVal.(string)
			if castOk {
				staticValue = &str
			}
		}

		if staticValue == nil {
			conditionalValueMap := conditionalContentMap["value"].(*schema.Set).List()
			if len(conditionalValueMap) == 0 {
				continue
			}
			log.Printf("conditionalValueMap: %v \n", conditionalValueMap)

			fieldOptionList := conditionalValueMap[0].(map[string]interface{})["option"].(*schema.Set).List()

			fieldOptions = make([]FieldOption, 0)
			for _, vOption := range fieldOptionList {
				optionMap := vOption.(map[string]interface{})
				fo := FieldOption{
					Label:   optionMap["label"].(string),
					Value:   optionMap["value"].(string),
					Checked: optionMap["checked"].(bool),
				}
				fieldOptions = append(fieldOptions, fo)
			}
		}

		c := ConditionalConfig{
			Source:    conditionalMap["source"].(string),
			Condition: conditionalMap["condition"].(string),
			Type:      conditionalMap["type"].(string),
			Options:   fieldOptions,
			Value:     staticValue,
		}
		conditionals = append(conditionals, c)
	}

	return conditionals
}
