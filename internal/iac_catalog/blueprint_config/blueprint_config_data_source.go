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

	fieldOptionsSchema := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		// MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"option": {
					Type:     schema.TypeList,
					Optional: true,
					MinItems: 1,
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
				},
			},
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

	formConfigSchema :=
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"shortText", "radio", "checkbox"}, false),
					},
					"options":         fieldOptionsSchema,
					"validation_rule": validationRulesSchema,
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
					"form_config": formConfigSchema,
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
				},
			},
		},
		"builder": { // it keeps the form builder (omit vars, override vars, ...) as json
			Description: "Form builder JSON (it keeps the parsed form builder as json)",
			Type:        schema.TypeString,
			Computed:    true,
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

	// new form variables (as JSON)
	formVariables := GetFormShape(*blueprintConfig)
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

	bpOutput, err := utils.ToJsonString(blueprintConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("blueprint_config", bpOutput)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(blueprintConfig.Id)

	return diags
}

// maps tf declaration to object
func GetBlueprintConfigFromSchema(d *schema.ResourceData) (*BluePrintConfig, error) {
	bp := &BluePrintConfig{}
	bp.Id = strconv.FormatInt(time.Now().Unix(), 10)
	if v, ok := d.GetOk("source"); ok {
		mapString := make(map[string]BluePrintConfig)
		for key, value := range v.(map[string]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			log.Printf("INPUT: %v/n", key)
			log.Printf("VALUE: %v", strValue)
			bc := BluePrintConfig{}
			err := json.Unmarshal([]byte(strValue), &bc)
			if err != nil {
				return nil, errors.New("invalid conversion to BluePrintConfig")
			}
			mapString[strKey] = bc
			bp.Children = append(bp.Children, bc)
		}
		str, err := utils.ToJsonString(bp)
		if err != nil {
			return nil, errors.New("invalid conversion to BluePrintConfig")
		}
		log.Printf("final bc: %s", str)
	}

	if v, ok := d.GetOk("omit_variables"); ok {
		list := v.(*schema.Set).List()
		omit := make([]string, len(list))
		for i, optionValue := range list {
			omit[i] = optionValue.(string)
		}
		bp.OmitVariables = omit
	}

	if v, ok := d.GetOk("variable"); ok {
		varsList := v.([]interface{})
		overrideVariables := make(map[string]OverrideVariable, 0)

		// override vars loop
		for _, f := range varsList {
			varOverrideMap := f.(map[string]interface{})
			varName := varOverrideMap["name"].(string)

			// Note: if it has a value, then it can NOT have a form_config
			isValueDefined := false
			value, ok := varOverrideMap["value"]
			if ok {
				valueStr, ok := value.(string)
				if ok {
					isValueDefined = valueStr != "" // NOTE: if the value is empty, we consider it as 'not defined'
					if isValueDefined {
						overrideVariables[varName] = OverrideVariable{
							VariableName: varName,
							DisplayName:  varOverrideMap["display_name"].(string),
							HelperText:   varOverrideMap["helper_text"].(string),
							Value:        &valueStr,
						}
					}
				}
			}

			// for config (we should only have 1 form_config by var)
			if formConfigInput, ok := varOverrideMap["form_config"]; ok {
				formConfigList := formConfigInput.(*schema.Set).List()

				formConfigListLen := len(formConfigList)

				if isValueDefined {
					if formConfigListLen != 0 {
						return nil, fmt.Errorf("a form_config can not be added when setting the variable's value. Var name: [%s], var value [%v]", varName, value)
					}

					// if the value is set and there's no form_config, then we're ok to continue processing the next variable override
					continue
				}

				if formConfigListLen == 0 {
					return nil, fmt.Errorf("a form_config must be defined for variable [%s]", varName)
				}
				if formConfigListLen > 1 {
					// it should be caught at schema check level - adding the check here to enforce it in case the schema changes
					return nil, errors.New("exactly one form_config must be defined")
				}

				formConfigMap := formConfigList[0].(map[string]interface{})
				variableType := formConfigMap["type"].(string)

				// field options
				var fieldOptions []FieldOption

				fieldOptionList := formConfigMap["options"].([]interface{})
				if variableType == RADIO_TYPE || variableType == CHECKBOX_TYPE {
					if len(fieldOptionList) != 1 {
						return nil, errors.New("one options block is required")
					}

					options := fieldOptionList[0].(map[string]interface{})["option"].([]interface{})
					fieldOptions = make([]FieldOption, len(options))
					optionCount := 0
					for _, vOption := range options {
						optionMap := vOption.(map[string]interface{})

						fieldOptions[optionCount] = FieldOption{
							Label:   optionMap["label"].(string),
							Value:   optionMap["value"].(string),
							Checked: optionMap["checked"].(bool),
						}
						optionCount++
					}
				} else if variableType == "shortText" {
					if len(fieldOptionList) > 0 {
						return nil, errors.New("ShortText variables can not have options")
					}

					// nothing should be done here.
				}

				// validation rules
				validationRulesList := formConfigMap["validation_rule"].(*schema.Set).List()
				validationRules := make([]ValidationRule, len(validationRulesList))

				for iValidationRule, validationRule := range validationRulesList {
					validationRuleMap := validationRule.(map[string]interface{})

					rule := validationRuleMap["rule"].(string)
					ruleValue := validationRuleMap["value"].(string)

					if rule == "isRequired" && ruleValue != "" {
						return nil, errors.New("'isRequired' validation rule can not have a value")
					}

					validationRules[iValidationRule] = ValidationRule{
						Rule:         rule,
						Value:        ruleValue,
						ErrorMessage: validationRuleMap["error_message"].(string),
					}
				}

				// build var config
				formConfig := FormConfig{
					Type:            variableType,
					FieldOptions:    fieldOptions,
					ValidationRules: validationRules,
				}

				// build the override variable wrapper object

				overrideVariables[varName] = OverrideVariable{
					VariableName: varName,
					DisplayName:  varOverrideMap["display_name"].(string),
					HelperText:   varOverrideMap["helper_text"].(string),
					FormConfig:   formConfig,
				}
			}
		}
		bp.OverrideVariables = overrideVariables
	}

	return bp, nil
}
