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
		FieldDefaultValue:   iacModuleVar.FieldDefaultValue,
		FieldValue:          iacModuleVar.FieldValue,
		RequiredValues:      overrideData.RequiredValues,
		AllowConsumerToEdit: true,
		IsHidden:            overrideData.IsHidden,
		UsedInHCL:           true, //if a user overrides, then it is used,
		Conditionals:        iacModuleVar.Conditionals,
		IsOverriden:         true,
	}

	if variableType == RADIO_TYPE || variableType == CHECKBOX_TYPE || variableType == LIST_TYPE {
		// try to map the value to an array of strings (options)
		// var fieldOptions []string
		// useValueFieldOptions := false

		// if overrideData.Value != "" {
		// 	err := json.Unmarshal([]byte(overrideData.Value), &fieldOptions)
		// 	useValueFieldOptions = err == nil
		// }

		// if there's an override with a value for this variable, we replace its value and mark it as uneditable
		switch {
		// case useValueFieldOptions:
		// 	if len(fieldOptions) > 0 {
		// 		newIacModuleVar.FormQuestion.FieldOptions = make([]generator.FieldOption, len(fieldOptions))

		// 		for i, option := range fieldOptions {
		// 			newIacModuleVar.FormQuestion.FieldOptions[i] = generator.FieldOption{
		// 				Label:   option,
		// 				FieldID: fmt.Sprintf("%s-%s", fieldID, option),
		// 				Value:   option,
		// 				Checked: true,
		// 			}
		// 		}
		// 		newIacModuleVar.AllowConsumerToEdit = false
		// 		newIacModuleVar.IsHidden = overrideData.IsHidden
		// 	} else {
		// 		newIacModuleVar.FormQuestion.FieldOptions = getDefaultFieldOptions(fieldID)
		// 	}

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
				//newIacModuleVar.IsHidden = overrideData.IsHidden
			} else {
				newIacModuleVar.FormQuestion.FieldOptions = getDefaultFieldOptions(fieldID)
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
			//newIacModuleVar.IsHidden = false
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
		if r.MatchString(overrideData.Value) {
			newIacModuleVar.FieldDataType = "hcl-expression"
			newIacModuleVar.IsHidden = true
		}
	}
	str, _ := json.MarshalIndent(newIacModuleVar, "", "    ")
	log.Debugf("New var result: %s", str)
	return newIacModuleVar
}

func getDefaultFieldOptions(fieldID string) []generator.FieldOption {
	value := "default"
	return []generator.FieldOption{
		{
			Label:   "Autogenerated Option. Please update this value",
			FieldID: fmt.Sprintf("%s-%s", fieldID, value),
			Value:   value,
			Checked: false,
		},
	}
}

func BuildGenericVariable(ov OverrideVariable) generator.FormShape {
	fieldID := fmt.Sprintf("%s.%s", GENERIC, ov.VariableName)

	validationRules := make([]generator.ValidationRule, len(ov.FormConfig.ValidationRules))
	for i, vr := range ov.FormConfig.ValidationRules {
		validationRules[i] = generator.ValidationRule{
			Rule:         vr.Rule,
			Value:        vr.Value,
			ErrorMessage: vr.ErrorMessage,
		}
	}

	if ov.FormConfig.Type == "" && ov.Value == "" {
		log.Fatalf("cant initialize generic variable %s without type", ov.VariableName)
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
			formVariable.FormQuestion.FieldOptions = getDefaultFieldOptions(fieldID)
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
	return formVariable
}
