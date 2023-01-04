package blueprint_config

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

// it creates an iac question format from override var data
/*
TODO: create a test over this function, perhaps it is worth it to rethink the inputs,
We need as an output a FormShape
*/
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
		FieldDataType:       iacModuleVar.FieldDataType,
		FieldDefaultValue:   iacModuleVar.FieldDefaultValue,
		FieldValue:          iacModuleVar.FieldValue,
		AllowConsumerToEdit: true,
		Conditionals:        make([]autocloudsdk.ConditionalConfig, len(overrideData.Conditionals)),
	}

	if overrideData.Value != "" {
		// starting with the naive approach to see if an string is a module
		// we will replace this introducing the notion of all outputs involved at the API process
		r := regexp.MustCompile("module[.]([A-Za-z0-9_]+)[.]outputs[.]([A-Za-z0-9_]+)")
		newIacModuleVar.FieldValue = overrideData.Value
		newIacModuleVar.FieldDefaultValue = overrideData.Value
		newIacModuleVar.AllowConsumerToEdit = false
		if r.MatchString(overrideData.Value) {
			newIacModuleVar.FieldDataType = "hcl-expression"
		} else {
			newIacModuleVar.FieldDataType = "string"
		}
	}

	if overrideData.FormConfig.Type == RADIO_TYPE || overrideData.FormConfig.Type == CHECKBOX_TYPE {
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
	for i, conditional := range overrideData.Conditionals {
		log.Printf("buildOverridenVariable -> conditional: %v \n", conditional)
		fieldOptions := make([]autocloudsdk.FieldOption, 0)
		for _, option := range conditional.Options {
			fo := autocloudsdk.FieldOption{
				FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
				Label:   option.Label,
				Value:   option.Value,
				Checked: option.Checked,
			}
			fieldOptions = append(fieldOptions, fo)
		}
		newIacModuleVar.Conditionals[i] = autocloudsdk.ConditionalConfig{
			Source:    conditional.Source,
			Condition: conditional.Condition,
			Options:   fieldOptions,
			Value:     conditional.Value,
			Type:      conditional.Type,
		}
	}
	str, _ := json.MarshalIndent(newIacModuleVar, "", "    ")
	log.Printf("formVariable: %s", string(str))

	return newIacModuleVar
}

func buildGenericVariable(ov OverrideVariable) autocloudsdk.FormShape {
	fieldID := fmt.Sprintf("%s.%s", GENERIC, ov.VariableName)

	validationRules := make([]autocloudsdk.ValidationRule, len(ov.FormConfig.ValidationRules))
	for i, vr := range ov.FormConfig.ValidationRules {
		validationRules[i] = autocloudsdk.ValidationRule{
			Rule:         vr.Rule,
			Value:        vr.Value,
			ErrorMessage: vr.ErrorMessage,
		}
	}

	fieldLabel := ov.VariableName
	if ov.DisplayName != "" {
		fieldLabel = ov.DisplayName
	}

	formVariable := autocloudsdk.FormShape{
		ID:     fieldID,
		Type:   ov.FormConfig.Type,
		Module: GENERIC,
		FormQuestion: autocloudsdk.FormQuestion{
			FieldID:         fieldID,
			FieldType:       ov.FormConfig.Type,
			ValidationRules: validationRules,
			FieldLabel:      fieldLabel,
			ExplainingText:  ov.HelperText,
		},
		AllowConsumerToEdit: true,
		Conditionals:        make([]autocloudsdk.ConditionalConfig, len(ov.Conditionals)),
	}

	if ov.FormConfig.Type == RADIO_TYPE || ov.FormConfig.Type == CHECKBOX_TYPE {
		// if the list is empty, set a default value
		if len(ov.FormConfig.FieldOptions) == 0 {
			value := "default"
			formVariable.FormQuestion.FieldOptions = []autocloudsdk.FieldOption{
				{
					Label:   "Autogenerated Option. Please update this value",
					FieldID: fmt.Sprintf("%s-%s", fieldID, value),
					Value:   value,
					Checked: false,
				},
			}
		} else {
			formVariable.FormQuestion.FieldOptions = make([]autocloudsdk.FieldOption, len(ov.FormConfig.FieldOptions))

			for i, option := range ov.FormConfig.FieldOptions {
				formVariable.FormQuestion.FieldOptions[i] = autocloudsdk.FieldOption{
					Label:   option.Label,
					FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
					Value:   option.Value,
					Checked: option.Checked,
				}
			}
		}
	}

	for i, conditional := range ov.Conditionals {
		fieldOptions := make([]autocloudsdk.FieldOption, len(conditional.Options))
		for j, option := range conditional.Options {
			fieldOptions[j] = autocloudsdk.FieldOption{
				FieldID: fmt.Sprintf("%s-%s", fieldID, option.Value),
				Label:   option.Label,
				Value:   option.Value,
				Checked: option.Checked,
			}
		}
		formVariable.Conditionals[i] = autocloudsdk.ConditionalConfig{
			Source:    conditional.Source,
			Condition: conditional.Condition,
			Options:   fieldOptions,
			Value:     conditional.Value,
			Type:      conditional.Type,
		}
	}
	str, _ := json.MarshalIndent(formVariable, "", "    ")
	log.Printf("formVariable: %s", string(str))
	return formVariable
}
