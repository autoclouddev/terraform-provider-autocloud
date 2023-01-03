package blueprint_config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

type BluePrintConfig struct {
	Id        string                   `json:"id"`
	RefName   string                   `json:"refName"`
	Variables []autocloudsdk.FormShape `json:"variables"`
	Children  []BluePrintConfig        `json:"children"`
}

type FormBuilder struct {
	//sourceModuleID    string
	//source            map[string]string
	OmitVariables     []string                    `json:"omitVariables"`
	OverrideVariables map[string]OverrideVariable `json:"overrideVariable"`
	BluePrintConfig   BluePrintConfig
}

type OverrideVariable struct {
	VariableName string              `json:"variableName"`
	Value        *string             `json:"value"`
	DisplayName  string              `json:"displayName"`
	HelperText   string              `json:"helperText"`
	FormConfig   *FormConfig         `json:"formConfig"`
	Conditionals []ConditionalConfig `json:"conditionals"`
}

type FormConfig struct {
	Type            string           `json:"type"`
	FieldOptions    []FieldOption    `json:"fieldOptions"`
	ValidationRules []ValidationRule `json:"validationRules"`
}
type ConditionalConfig struct {
	Source    string        `json:"source"`
	Condition string        `json:"condition"`
	Type      string        `json:"type"`
	Options   []FieldOption `json:"options"`
	Value     *string       `json:"value"`
}

type ValidationRule struct {
	Rule         string `json:"rule"`
	Value        string `json:"value"`
	ErrorMessage string `json:"errorMessage"`
}

type FieldOption struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Checked bool   `json:"checked"`
}

const GENERIC = "generic"
const RADIO_TYPE = "radio"
const CHECKBOX_TYPE = "checkbox"
const SHORTTEXT_TYPE = "shortText"

func DataSourceBlueprintConfig() *schema.Resource {
	setOfStringSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	optionItemSchema := &schema.Schema{
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
	}

	fieldOptionsSchema := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"option": optionItemSchema,
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
						ValidateFunc: validation.StringInSlice([]string{SHORTTEXT_TYPE, RADIO_TYPE, CHECKBOX_TYPE}, false),
					},
					"options":         fieldOptionsSchema,
					"validation_rule": validationRulesSchema,
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
					"conditional": conditionalSchema,
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
	formBuilder, err := GetFormBuilder(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: Refactor this section to its own method
	v, ok := d.GetOk("source")
	var blueprint BluePrintConfig
	if v != nil && ok {
		blueprint = mapVariables(formBuilder)
	} else {
		blueprint = mapModuleVariables(formBuilder)
	}
	// ENDS HERE

	// new form variables (as JSON)
	formVariables := GetFormShape(blueprint)
	if err != nil {
		return diag.FromErr(err)
	}

	// validate variables conditionals (by now, we're supporting conditionals only referencing radio fields)
	err = validateConditionals(formVariables)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonVariables, err := utils.ToJsonString(formVariables)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("config", jsonVariables)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("\nJSON config %v\n\n", jsonVariables)
	// TODO: end of deprecation
	configString, err := utils.ToJsonString(blueprint)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("blueprint_config", configString)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(blueprint.Id)

	return diags
}

// maps tf declaration to object
func GetFormBuilder(d *schema.ResourceData) (*FormBuilder, error) {
	formBuilder := &FormBuilder{}
	formBuilder.BluePrintConfig.Id = strconv.FormatInt(time.Now().Unix(), 10)
	if v, ok := d.GetOk("source"); ok {
		mapString := make(map[string]BluePrintConfig)

		for key, value := range v.(map[string]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			fmt.Println(strValue)
			bc := BluePrintConfig{}
			err := json.Unmarshal([]byte(strValue), &bc)
			if err != nil {
				return nil, errors.New("invalid conversion to BluePrintConfig")
			}
			mapString[strKey] = bc
			formBuilder.BluePrintConfig.Children = append(formBuilder.BluePrintConfig.Children, bc)
		}
	}

	if v, ok := d.GetOk("omit_variables"); ok {
		list := v.(*schema.Set).List()
		omit := make([]string, len(list))
		for i, optionValue := range list {
			omit[i] = optionValue.(string)
		}
		formBuilder.OmitVariables = omit
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

			// form_config
			var formConfig *FormConfig = nil
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
					if conditionalsInput, conditionalsOk := varOverrideMap["conditionals"]; conditionalsOk {
						if len(conditionalsInput.(*schema.Set).List()) > 0 {
							// if we don't have any conditionals, we should have at least 1 form_config
							return nil, fmt.Errorf("a form_config must be defined for the variable [%s]", varName)
						}
					}
				}
				if formConfigListLen > 1 {
					// it should be caught at schema check level - adding the check here to enforce it in case the schema changes
					return nil, errors.New("exactly one form_config must be defined")
				}

				if formConfigListLen == 1 {
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
					} else if variableType == SHORTTEXT_TYPE {
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
					formConfig = &FormConfig{
						Type:            variableType,
						FieldOptions:    fieldOptions,
						ValidationRules: validationRules,
					}
				}
			}

			// build the override variable wrapper object
			overrideVariables[varName] = OverrideVariable{
				VariableName: varName,
				DisplayName:  varOverrideMap["display_name"].(string),
				HelperText:   varOverrideMap["helper_text"].(string),
				FormConfig:   formConfig,
				Conditionals: getConditionals(varOverrideMap),
			}
		}

		formBuilder.OverrideVariables = overrideVariables
	}

	return formBuilder, nil
}

func getConditionals(varOverrideMap map[string]interface{}) []ConditionalConfig {
	var conditionals []ConditionalConfig

	if conditionalsInput, conditionalsOk := varOverrideMap["conditional"]; conditionalsOk {
		conditionalsList := conditionalsInput.(*schema.Set).List()
		conditionals = make([]ConditionalConfig, len(conditionalsList))

		for i, conditional := range conditionalsList {
			conditionalMap := conditional.(map[string]interface{})
			conditionalContentMap := conditionalMap["content"].(*schema.Set).List()[0].(map[string]interface{}) // length validated at schema level
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
				if len(conditionalValueMap) > 0 {
					fieldOptionList := conditionalValueMap[0].(map[string]interface{})["option"].([]interface{})

					fieldOptions = make([]FieldOption, len(fieldOptionList))

					optionCount := 0
					for _, vOption := range fieldOptionList {
						optionMap := vOption.(map[string]interface{})

						fieldOptions[optionCount] = FieldOption{
							Label:   optionMap["label"].(string),
							Value:   optionMap["value"].(string),
							Checked: optionMap["checked"].(bool),
						}
						optionCount++
					}
				}
			}

			conditionals[i] = ConditionalConfig{
				Source:    conditionalMap["source"].(string),
				Condition: conditionalMap["condition"].(string),
				Type:      conditionalMap["type"].(string),
				Options:   fieldOptions,
				Value:     staticValue,
			}
		}
	}

	return conditionals
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
			if dependencyVariable, ok := varsMap[conditional.Source]; ok {
				if dependencyVariable.FormQuestion.FieldType != RADIO_TYPE {
					return fmt.Errorf("the conditional's source variable can only be of 'radio' type [variable: %v, source variable: %v, source variable type: %v]", variable.ID, conditional.Source, dependencyVariable.FormQuestion.FieldType)
				}
			}
		}
	}

	return nil
}

// TODO: Refactor to share logic with mapVariables
// mapVariables
func mapModuleVariables(formBuilder *FormBuilder) BluePrintConfig {
	newForm := BluePrintConfig{
		Id: strconv.FormatInt(time.Now().Unix(), 10),
	}

	for _, variable := range formBuilder.OverrideVariables {
		fieldID := fmt.Sprintf("%s.%s", GENERIC, variable.VariableName)

		var varType = ""
		var validationRules []autocloudsdk.ValidationRule
		if variable.FormConfig != nil {
			varType = variable.FormConfig.Type
			validationRules = make([]autocloudsdk.ValidationRule, len(variable.FormConfig.ValidationRules))
			for i, vr := range variable.FormConfig.ValidationRules {
				validationRules[i] = autocloudsdk.ValidationRule{
					Rule:         vr.Rule,
					Value:        vr.Value,
					ErrorMessage: vr.ErrorMessage,
				}
			}
		}

		fieldLabel := variable.VariableName
		if variable.DisplayName != "" {
			fieldLabel = variable.DisplayName
		}

		newVariable := autocloudsdk.FormShape{
			ID:     fieldID,
			Type:   varType,
			Module: GENERIC,
			FormQuestion: autocloudsdk.FormQuestion{
				FieldID:         fieldID,
				FieldType:       varType,
				ValidationRules: validationRules,
				FieldLabel:      fieldLabel,
				ExplainingText:  variable.HelperText,
			},
			AllowConsumerToEdit: true,
			Conditionals:        mapToSdkConditionals(fieldID, variable.Conditionals),
		}

		if variable.FormConfig != nil && (variable.FormConfig.Type == RADIO_TYPE || variable.FormConfig.Type == CHECKBOX_TYPE) {
			// if the list is empty, set a default value
			if len(variable.FormConfig.FieldOptions) == 0 {
				value := "default"
				newVariable.FormQuestion.FieldOptions = []autocloudsdk.FieldOption{
					{
						Label:   "Autogenerated Option. Please update this value",
						FieldID: fmt.Sprintf("%s-%s", fieldID, value),
						Value:   value,
						Checked: false,
					},
				}
			} else {
				newVariable.FormQuestion.FieldOptions = make([]autocloudsdk.FieldOption, len(variable.FormConfig.FieldOptions))

				for i, option := range variable.FormConfig.FieldOptions {
					newVariable.FormQuestion.FieldOptions[i] = autocloudsdk.FieldOption{
						Label:   option.Label,
						FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
						Value:   option.Value,
						Checked: option.Checked,
					}
				}
			}
		}

		newForm.Variables = append(newForm.Variables, newVariable)
	}

	return newForm
}

// TODO: Update this variables to remove reference to module resource
// it omits and overrides the iacModule variables based on the formBuilder definition
func mapVariables(formBuilder *FormBuilder) BluePrintConfig {
	newForm := BluePrintConfig{
		Id: strconv.FormatInt(time.Now().Unix(), 10),
		// Variables: here should be all the variables,
	}

	for _, config := range formBuilder.BluePrintConfig.Children {
		// to store processed vars (without omitted vars, with overridden vars).
		newConfig := BluePrintConfig{
			Id:      config.Id,
			RefName: config.RefName,
		}

		for _, iacModuleVar := range config.Variables {
			varName, err := utils.GetVariableID(iacModuleVar.ID)

			if err != nil {
				return newForm
			}

			// omit vars
			if utils.Contains(formBuilder.OmitVariables, varName) {
				log.Printf("the [%s] variable was omitted", varName)
				continue
			}

			// if no override statement, then keep the original var without changes
			updatedIacModuleVar := iacModuleVar
			if overrideVariableData, ok := formBuilder.OverrideVariables[varName]; ok {
				updatedIacModuleVar = buildOverriddenVariable(iacModuleVar, overrideVariableData)
			}

			newConfig.Variables = append(newConfig.Variables, updatedIacModuleVar)
		}
		// sort questions to keep ordering consistent
		sort.Slice(newConfig.Variables, func(i, j int) bool {
			return newConfig.Variables[i].ID < newConfig.Variables[j].ID
		})
		newForm.Children = append(newForm.Children, newConfig)
	}

	return newForm
}

// it creates an iac question format from override var data
/*
TODO: create a test over this function, perhaps it is worth it to rethink the inputs,
We need as an output a FormShape
*/
//nolint:golint,unused
func buildOverriddenVariable(iacModuleVar autocloudsdk.FormShape, overrideData OverrideVariable) autocloudsdk.FormShape {
	fieldID := iacModuleVar.ID

	// map validation rules
	var validationRules []autocloudsdk.ValidationRule
	if overrideData.FormConfig != nil {
		validationRules := make([]autocloudsdk.ValidationRule, len(overrideData.FormConfig.ValidationRules))
		for i, vr := range overrideData.FormConfig.ValidationRules {
			validationRules[i] = autocloudsdk.ValidationRule{
				Rule:         vr.Rule,
				Value:        vr.Value,
				ErrorMessage: vr.ErrorMessage,
			}
		}
	}

	fieldLabel := overrideData.VariableName
	if overrideData.DisplayName != "" {
		fieldLabel = overrideData.DisplayName
	}

	explainingText := iacModuleVar.FormQuestion.ExplainingText
	if overrideData.HelperText != "" {
		explainingText = overrideData.HelperText
	}

	variableType := iacModuleVar.FormQuestion.FieldType
	if overrideData.FormConfig != nil && overrideData.FormConfig.Type != "" {
		variableType = overrideData.FormConfig.Type
	}

	newIacModuleVar := autocloudsdk.FormShape{
		ID:     iacModuleVar.ID,
		Type:   variableType,
		Module: iacModuleVar.Module,
		FormQuestion: autocloudsdk.FormQuestion{
			FieldID:         fieldID,
			FieldType:       variableType,
			FieldLabel:      fieldLabel,
			ExplainingText:  explainingText,
			ValidationRules: validationRules,
		},
		FieldDataType:       iacModuleVar.FieldDataType,
		FieldDefaultValue:   iacModuleVar.FieldDefaultValue,
		FieldValue:          iacModuleVar.FieldValue,
		AllowConsumerToEdit: true,
		Conditionals:        mapToSdkConditionals(fieldID, overrideData.Conditionals),
	}

	if overrideData.Value != nil {
		// starting with the naive approach to see if an string is a module
		// we will replace this introducing the notion of all outputs involved at the API process
		r := regexp.MustCompile("module[.]([A-Za-z0-9_]+)[.]outputs[.]([A-Za-z0-9_]+)")
		newIacModuleVar.FieldValue = *overrideData.Value
		newIacModuleVar.FieldDefaultValue = *overrideData.Value
		newIacModuleVar.AllowConsumerToEdit = false
		if r.MatchString(*overrideData.Value) {
			newIacModuleVar.FieldDataType = "hcl-expression"
		} else {
			newIacModuleVar.FieldDataType = "string"
		}
	}

	if overrideData.FormConfig != nil && (overrideData.FormConfig.Type == RADIO_TYPE || overrideData.FormConfig.Type == CHECKBOX_TYPE) {
		// if the list is empty, set a default value
		if len(overrideData.FormConfig.FieldOptions) == 0 {
			value := "default"
			newIacModuleVar.FormQuestion.FieldOptions = []autocloudsdk.FieldOption{
				{
					Label:   "Autogenerated Option. Please update this value",
					FieldID: fmt.Sprintf("%s-%s", fieldID, value),
					Value:   value,
					Checked: false,
				},
			}
		} else {
			newIacModuleVar.FormQuestion.FieldOptions = make([]autocloudsdk.FieldOption, len(overrideData.FormConfig.FieldOptions))

			for i, option := range overrideData.FormConfig.FieldOptions {
				newIacModuleVar.FormQuestion.FieldOptions[i] = autocloudsdk.FieldOption{
					Label:   option.Label,
					FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
					Value:   option.Value,
					Checked: option.Checked,
				}
			}
		}
	}

	return newIacModuleVar
}

func mapToSdkConditionals(fieldID string, conditionals []ConditionalConfig) []autocloudsdk.ConditionalConfig {
	sdkConditionals := make([]autocloudsdk.ConditionalConfig, len(conditionals))
	for i, conditional := range conditionals {
		fieldOptions := make([]autocloudsdk.FieldOption, len(conditional.Options))
		for j, option := range conditional.Options {
			fieldOptions[j] = autocloudsdk.FieldOption{
				FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
				Label:   option.Label,
				Value:   option.Value,
				Checked: option.Checked,
			}
		}
		sdkConditionals[i] = autocloudsdk.ConditionalConfig{
			Source:    conditional.Source,
			Condition: conditional.Condition,
			Options:   fieldOptions,
			Value:     conditional.Value,
			Type:      conditional.Type,
		}
	}

	return sdkConditionals
}
