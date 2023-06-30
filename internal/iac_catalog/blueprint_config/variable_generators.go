package blueprint_config

import (
	"encoding/json"
	"fmt"
	"strings"

	//"log"
	"regexp"

	"github.com/apex/log"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/logger"
)

// it creates an iac question format from override var data
/*
TODO: create a test over this function, perhaps it is worth it to rethink the inputs,
We need as an output a FormShape
*/
func BuildOverridenVariable(iacModuleVar generator.FormShape, overrideData OverrideVariable) generator.FormShape {
	var log = logger.Create(log.Fields{"fn": "BuildOverridenVariable()"})
	fieldID := iacModuleVar.ID

	// map validation rules
	validationRules := make([]generator.ValidationRule, len(overrideData.FormConfig.ValidationRules))
	for i, vr := range overrideData.FormConfig.ValidationRules {
		validationRules[i] = generator.ValidationRule{
			Rule:         vr.Rule,
			Scope:        vr.Scope,
			Value:        vr.Value,
			ErrorMessage: vr.ErrorMessage,
		}
	}

	fieldLabel := iacModuleVar.FormQuestion.FieldLabel
	if overrideData.DisplayName != "" {
		fieldLabel = overrideData.DisplayName
	}

	explainingText := iacModuleVar.FormQuestion.ExplainingText
	if overrideData.HelperText != "" {
		explainingText = overrideData.HelperText
	}

	variableType := iacModuleVar.FormQuestion.FieldType
	if overrideData.FormConfig.Type != "" {
		variableType = overrideData.FormConfig.Type
	}

	// check if there's a "default" value override, otherwise use the IacModule var default value
	fieldDefaultValue := iacModuleVar.FieldDefaultValue
	fieldValue := iacModuleVar.FieldValue
	if overrideData.VariableContent.Default != "" {
		fieldDefaultValue = overrideData.VariableContent.Default
		fieldValue = overrideData.VariableContent.Default

		// if the list is empty, populate the field options with the variable's default values
		if len(overrideData.FormConfig.FieldOptions) == 0 && (variableType == RADIO_TYPE || variableType == CHECKBOX_TYPE || variableType == LIST_TYPE) {
			fieldOptions, err := ToFieldOptions(fieldValue)
			if err == nil {
				overrideData.FormConfig.FieldOptions = fieldOptions
			}
		}
	}

	newIacModuleVar := generator.FormShape{
		ID:       iacModuleVar.ID,
		Module:   iacModuleVar.Module,
		ModuleID: iacModuleVar.ModuleID,
		FormQuestion: generator.FormQuestion{
			FieldID:         fieldID,
			FieldType:       variableType,
			FieldLabel:      fieldLabel,
			ExplainingText:  explainingText,
			ValidationRules: validationRules,
		},
		FieldDataType:       iacModuleVar.FieldDataType,
		FieldDefaultValue:   fieldDefaultValue,
		FieldValue:          fieldValue,
		RequiredValues:      overrideData.RequiredValues,
		AllowConsumerToEdit: true,
		IsHidden:            overrideData.IsHidden,
		UsedInHCL:           true, //if a user overrides, then it is used,
		Conditionals:        iacModuleVar.Conditionals,
		IsOverriden:         true,
		InterpolationVars:   iacModuleVar.InterpolationVars, //OVERRIDE THIS
	}

	if variableType == RADIO_TYPE || variableType == CHECKBOX_TYPE || variableType == LIST_TYPE {
		// if there's an override with a value for this variable, we replace its value and mark it as uneditable
		switch {
		case len(overrideData.FormConfig.FieldOptions) == 0: // if the list is empty, set a default value
			// Use module default values
			if len(iacModuleVar.FormQuestion.FieldOptions) > 0 {
				newIacModuleVar.FormQuestion.FieldOptions = make([]generator.FieldOption, len(iacModuleVar.FormQuestion.FieldOptions))

				isChecked := false
				for i, option := range iacModuleVar.FormQuestion.FieldOptions {
					isChecked = option.Value == overrideData.Value || isChecked
					newIacModuleVar.FormQuestion.FieldOptions[i] = generator.FieldOption{
						Label:   option.Label,
						FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
						Value:   option.Value,
						Checked: option.Value == overrideData.Value,
					}
				}
				newIacModuleVar.AllowConsumerToEdit = !isChecked
			} else {
				newIacModuleVar.FormQuestion.FieldOptions = make([]generator.FieldOption, 0)
			}
		default:
			newIacModuleVar.FormQuestion.FieldOptions = make([]generator.FieldOption, len(overrideData.FormConfig.FieldOptions))

			for i, option := range overrideData.FormConfig.FieldOptions {
				newIacModuleVar.FormQuestion.FieldOptions[i] = generator.FieldOption{
					Label:   option.Label,
					FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
					Value:   option.Value,
					Checked: option.Checked,
				}
			}
		}
	}
	log.Debugf("conditionalLen: %v \n", len(overrideData.Conditionals))

	// add conditionals to the pre-existent conditionals
	for _, conditional := range overrideData.Conditionals {
		conditionalSource := conditional.Source

		// if it's not a multipart id, attach the conditionl to the current question
		if len(strings.Split(conditional.Source, ".")) == 1 {
			conditionalSource = fmt.Sprintf("%s.%s", iacModuleVar.Module, conditional.Source)
		}

		newConditional := generator.ConditionalConfig{
			Source:         conditionalSource,
			Condition:      conditional.Condition,
			Value:          conditional.Value,
			Type:           conditional.FormConfig.Type,
			RequiredValues: conditional.RequiredValues,
			Options:        make([]generator.FieldOption, 0), //conditional.FormConfig.FieldOptions,
		}
		for _, c := range conditional.FormConfig.FieldOptions {
			ao := generator.FieldOption{
				FieldID: fmt.Sprintf("%s-%s", fieldID, c.Value),
				Label:   c.Label,
				Value:   c.Value,
				Checked: c.Checked,
			}
			newConditional.Options = append(newConditional.Options, ao)
		}
		newIacModuleVar.Conditionals = append(newIacModuleVar.Conditionals, newConditional)
		str, _ := json.MarshalIndent(newConditional, "", "    ")
		log.Debugf("created conditional: %s", string(str))
	}
	if overrideData.Value != "" {
		// starting with the naive approach to see if an string is a module
		// we will replace this introducing the notion of all outputs involved at the API process
		r := regexp.MustCompile("module[.]([A-Za-z0-9_]+)[.]outputs[.]([A-Za-z0-9_]+)")
		newIacModuleVar.FieldValue = overrideData.Value
		newIacModuleVar.FieldDefaultValue = overrideData.Value
		newIacModuleVar.AllowConsumerToEdit = false
		//newIacModuleVar.IsHidden = overrideData.IsHidden
		newIacModuleVar.UsedInHCL = true

		if r.MatchString(overrideData.Value) || variableType == RAW_TYPE {
			newIacModuleVar.FieldDataType = "hcl-expression"
			newIacModuleVar.IsHidden = true
		}
	}

	if overrideData.InterpolationVars != nil || len(overrideData.InterpolationVars) > 0 {
		newIacModuleVar.InterpolationVars = overrideData.InterpolationVars
	}

	str, _ := json.MarshalIndent(newIacModuleVar, "", "    ")
	log.Debugf("New var result: %s", str)
	return newIacModuleVar
}

func BuildGenericVariable(ov OverrideVariable) (generator.FormShape, error) {
	fieldID := fmt.Sprintf("%s.%s", GENERIC, ov.VariableName)

	validationRules := make([]generator.ValidationRule, len(ov.FormConfig.ValidationRules))
	for i, vr := range ov.FormConfig.ValidationRules {
		validationRules[i] = generator.ValidationRule{
			Rule:         vr.Rule,
			Scope:        vr.Scope,
			Value:        vr.Value,
			ErrorMessage: vr.ErrorMessage,
		}
	}

	if ov.FormConfig.Type == "" && ov.Value == "" {
		//return generator.FormShape{}, fmt.Errorf("cant initialize generic variable %s without  a type", ov.VariableName)
		log.Debugf("cant initialize generic variable %s without  a type", ov.VariableName)
	}

	fieldLabel := ov.VariableName
	if ov.DisplayName != "" {
		fieldLabel = ov.DisplayName
	}

	fieldValue := ov.Value
	if ov.FormConfig.Type == MAP_TYPE {
		fieldValue = "{}" // empty map
		if ov.Value != "" {
			fieldValue = ov.Value
		}
	}

	formVariable := generator.FormShape{
		ID:         fieldID,
		Module:     GENERIC,
		FieldValue: fieldValue,
		FormQuestion: generator.FormQuestion{
			FieldID:         fieldID,
			FieldType:       ov.FormConfig.Type,
			ValidationRules: validationRules,
			FieldLabel:      fieldLabel,
			ExplainingText:  ov.HelperText,
		},
		AllowConsumerToEdit: len(fieldValue) == 0,
		IsHidden:            ov.IsHidden,
		UsedInHCL:           ov.UsedInHCL,
		RequiredValues:      ov.RequiredValues,
		Conditionals:        make([]generator.ConditionalConfig, len(ov.Conditionals)),
	}

	if ov.FormConfig.Type == RADIO_TYPE || ov.FormConfig.Type == CHECKBOX_TYPE || ov.FormConfig.Type == LIST_TYPE {
		// if the list is empty, set a default value
		if len(ov.FormConfig.FieldOptions) == 0 {
			formVariable.FormQuestion.FieldOptions = make([]generator.FieldOption, 0)
		} else {
			formVariable.FormQuestion.FieldOptions = make([]generator.FieldOption, len(ov.FormConfig.FieldOptions))

			for i, option := range ov.FormConfig.FieldOptions {
				formVariable.FormQuestion.FieldOptions[i] = generator.FieldOption{
					Label:   option.Label,
					FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
					Value:   option.Value,
					Checked: option.Checked,
				}
			}
		}

		// this is inside the if, which contains all keys of this map
		formTypesToTerraformTypes := map[string]string{
			"radio":    "string",
			"checkbox": "list(string)",
			"list":     "list(string)",
		}

		// in the frontend we allow the user to add more options if the type is a list(string)
		if ov.FormConfig.Type == LIST_TYPE {
			formVariable.FieldDataType = formTypesToTerraformTypes[ov.FormConfig.Type]
		}
	}

	for i, conditional := range ov.Conditionals {
		formVariable.Conditionals[i] = generator.ConditionalConfig{
			Source:         conditional.Source,
			Condition:      conditional.Condition,
			Value:          conditional.Value,
			Type:           conditional.FormConfig.Type,
			RequiredValues: conditional.RequiredValues,
		}
	}
	//str, _ := json.MarshalIndent(formVariable, "", "    ")
	//log.Printf("formVariable: %s", string(str))
	return formVariable, nil
}

func ToFieldOptions(options string) ([]FieldOption, error) {
	var defaultOptions []interface{} = make([]interface{}, 0)
	var fieldOptions []FieldOption = make([]FieldOption, 0)
	err := json.Unmarshal([]byte(options), &defaultOptions)
	if err != nil {
		return fieldOptions, err
	}

	fieldOptions = make([]FieldOption, len(defaultOptions))

	for i, option := range defaultOptions {
		strOption := fmt.Sprintf("%v", option)

		fieldOptions[i] = FieldOption{
			Label:   strOption,
			Value:   strOption,
			Checked: true,
		}
	}

	return fieldOptions, nil
}
