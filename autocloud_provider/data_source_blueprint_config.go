package autocloud_provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

type FormBuilder struct {
	sourceModuleID    string
	OmitVariables     []string                    `json:"omitVariables"`
	OverrideVariables map[string]OverrideVariable `json:"overrideVariable"`
}

type OverrideVariable struct {
	VariableName string     `json:"variableName"`
	Value        *string    `json:"value"`
	DisplayName  string     `json:"displayName"`
	HelperText   string     `json:"helperText"`
	FormConfig   FormConfig `json:"formConfig"`
}

type FormConfig struct {
	Type            string           `json:"type"`
	FieldOptions    []FieldOption    `json:"fieldOptions"`
	ValidationRules []ValidationRule `json:"validationRules"`
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

func dataSourceBlueprintConfig() *schema.Resource {
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
					"field_options":   fieldOptionsSchema,
					"validation_rule": validationRulesSchema,
				},
			},
		}

	return &schema.Resource{
		Description: "terraform form processor (form builder)",
		ReadContext: dataSourceBlueprintConfigRead,
		Schema: map[string]*schema.Schema{
			"source_module_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"omit_variables": setOfStringSchema,
			"override_variable": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variable_name": {
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
			"form_config": { // the form as json to replace the default variables
				Description: "Processed form variables JSON (to replace the default module variables variables)",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceBlueprintConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// map the resource to a FormBuilder object
	formBuilder, err := getFormBuilder(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// build the form shape (iac questions' format)
	c := m.(*autocloudsdk.Client)
	iacModule, err := c.GetModule(formBuilder.sourceModuleID)
	if err != nil {
		return diag.FromErr(err)
	}

	// omit and override variables
	newForm, err := mapModuleVariables(formBuilder, iacModule)
	if err != nil {
		return diag.FromErr(err)
	}

	// builder (list of vars to override and omit, as JSON)
	jsonString, err := toJsonString(formBuilder)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("builder", jsonString)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("\nJSON builder %v\n", jsonString)

	// new form variables (as JSON)
	jsonString, err = toJsonString(newForm)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("form_config", jsonString)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("\nJSON form_config %v\n\n", jsonString)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

// maps tf declaration to object
func getFormBuilder(d *schema.ResourceData) (*FormBuilder, error) {
	formBuilder := &FormBuilder{}

	if v, ok := d.GetOk("source_module_id"); ok {
		formBuilder.sourceModuleID = v.(string)
	}

	if v, ok := d.GetOk("omit_variables"); ok {
		list := v.(*schema.Set).List()
		omit := make([]string, len(list))
		for i, optionValue := range list {
			omit[i] = optionValue.(string)
		}
		formBuilder.OmitVariables = omit
	}

	if v, ok := d.GetOk("override_variable"); ok {
		varsList := v.([]interface{})
		overrideVariables := make(map[string]OverrideVariable, 0)

		// override vars loop
		for _, f := range varsList {
			varOverrideMap := f.(map[string]interface{})
			varName := varOverrideMap["variable_name"].(string)

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
						return nil, fmt.Errorf("A form_config can not be added when setting the variable's value. Var name: [%s], var value [%v]", varName, value)
					}

					// if the value is set and there's no form_config, then we're ok to continue processing the next variable override
					continue
				}

				if formConfigListLen == 0 {
					return nil, fmt.Errorf("A form_config must be defined for variable [%s]", varName)
				}
				if formConfigListLen > 1 {
					// it should be caught at schema check level - adding the check here to enforce it in case the schema changes
					return nil, errors.New("Exactly one form_config must be defined")
				}

				formConfigMap := formConfigList[0].(map[string]interface{})
				variableType := formConfigMap["type"].(string)

				// field options
				var fieldOptions []FieldOption

				fieldOptionList := formConfigMap["field_options"].([]interface{})
				if variableType == "radio" || variableType == "checkbox" {
					if len(fieldOptionList) != 1 {
						return nil, errors.New("One field_options block is required")
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

		formBuilder.OverrideVariables = overrideVariables
	}

	return formBuilder, nil
}

// it omits and overrides the iacModule variables based on the formBuilder definition
func mapModuleVariables(formBuilder *FormBuilder, iacModule *autocloudsdk.IacModule) ([]autocloudsdk.FormShape, error) {
	// parse iacModule.Variables string into a []autocloudsdk.FormShape slice
	iacModuleVars, err := ParseVariables(iacModule.Variables)

	if err != nil {
		return nil, err
	}

	// to store processed vars (without omitted vars, with overridden vars)
	var overridenIacModuleVars = make([]autocloudsdk.FormShape, 0)

	for _, iacModuleVar := range iacModuleVars {
		varName, err := getVariableID(iacModuleVar.ID)

		if err != nil {
			return nil, err
		}

		// omit vars
		if Contains(formBuilder.OmitVariables, varName) {
			log.Printf("the [%s] variable was omitted", varName)
			continue
		}

		// if no override statement, then keep the original var without changes
		updatedIacModuleVar := iacModuleVar
		if overrideVariableData, ok := formBuilder.OverrideVariables[varName]; ok {
			updatedIacModuleVar = buildOverridenVariable(iacModuleVar, overrideVariableData)
		}

		overridenIacModuleVars = append(overridenIacModuleVars, updatedIacModuleVar)
	}

	// sort questions to keep ordering consistent
	sort.Slice(overridenIacModuleVars, func(i, j int) bool {
		return overridenIacModuleVars[i].ID < overridenIacModuleVars[j].ID
	})

	return overridenIacModuleVars, nil
}

// it creates an iac question format from override var data
func buildOverridenVariable(iacModuleVar autocloudsdk.FormShape, overrideData OverrideVariable) autocloudsdk.FormShape {
	fieldID := iacModuleVar.ID

	// map validation rules
	validationRules := make([]autocloudsdk.ValidationRule, len(overrideData.FormConfig.ValidationRules))
	for i, vr := range overrideData.FormConfig.ValidationRules {
		validationRules[i] = autocloudsdk.ValidationRule{
			Rule:         vr.Rule,
			Value:        vr.Value,
			ErrorMessage: vr.ErrorMessage,
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

	newIacModuleVar := autocloudsdk.FormShape{
		ID:     iacModuleVar.ID,
		Type:   overrideData.FormConfig.Type,
		Module: iacModuleVar.Module,
		FormQuestion: autocloudsdk.FormQuestion{
			FieldID:         fieldID,
			FieldType:       overrideData.FormConfig.Type,
			FieldLabel:      fieldLabel,
			ExplainingText:  explainingText,
			ValidationRules: validationRules,
		},
		FieldDataType:     iacModuleVar.FieldDataType,
		FieldDefaultValue: iacModuleVar.FieldDefaultValue,
		FieldValue:        iacModuleVar.FieldValue,
	}

	if overrideData.Value != nil {
		newIacModuleVar.FieldValue = *overrideData.Value
		newIacModuleVar.FieldDefaultValue = *overrideData.Value
		newIacModuleVar.FieldDataType = "hcl-expression"
	}

	if overrideData.FormConfig.Type == "radio" || overrideData.FormConfig.Type == "checkbox" {
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
