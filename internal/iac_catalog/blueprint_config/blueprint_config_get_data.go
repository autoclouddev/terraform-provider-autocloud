package blueprint_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils/interpolation_utils"
)

func GetBlueprintConfigSources(v interface{}, bp *BluePrintConfig, aliases blueprint_config_references.Data) error {
	bp.Children = make(map[string]*BluePrintConfig, 0)
	sources := v.(map[string]interface{})
	for key, value := range sources {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		log.Printf("SOURCE_INPUT_KEY: %v\n", key)
		formattedValue, _ := utils.PrettyString(strValue)
		log.Printf("SOURCE_INPUT_VALUE: %v", formattedValue)
		bc := BluePrintConfig{}
		err := json.Unmarshal([]byte(strValue), &bc)
		if err != nil {
			return errors.New("invalid conversion to BluePrintConfig")
		}
		aliasKey := fmt.Sprintf("%s#%s", strKey, bp.Id)
		aliases.SetValue(aliasKey, bc.Id)

		bp.Children[strKey] = &bc
	}
	return nil
}

func GetBlueprintConfigOmitVariables(v interface{}, bp *BluePrintConfig, aliases blueprint_config_references.Data) error {
	omit_variables := v.([]interface{})
	omit := make([]string, len(omit_variables))
	for i, optionValue := range omit_variables {
		omit[i] = optionValue.(string)
	}
	bp.OmitVariables = omit
	return nil
}

func GetBlueprintConfigDisplayOrder(v interface{}, bp *BluePrintConfig, aliases blueprint_config_references.Data) error {
	display_order := v.([]interface{})
	for _, currentVar := range display_order {
		displayOrder := DisplayOrder{}
		var values []string
		varOverrideMap := currentVar.(map[string]interface{})
		displayOrder.Priority = varOverrideMap["priority"].(int)
		valuesList := varOverrideMap["values"].([]interface{})
		for _, value := range valuesList {
			valueStr := value.(string)
			displayValue := ""
			// valueStr can be <alias>.variables.<variable_name> or <module>.<variable_name>.
			// If it's built with an alias, we need to convert it to <module>.<variable_name>
			if utils.HasReference(valueStr) {
				// path[0] => alias
				// path[1] => "variables"
				// path[2] => variable_name
				// convert alias to module name
				moduleName := GetModuleNameFromVariable(valueStr, aliases, bp)
				paths := strings.Split(valueStr, ".")
				if len(moduleName) > 0 {
					displayValue = fmt.Sprintf("%s.%s", moduleName, paths[2])
				}
				//referenceName := aliases.GetValue(paths[0])
				if len(displayValue) == 0 {
					// if there isn't any module name for the alias we just use the variable name
					displayValue = fmt.Sprintf("%s.%s", "generic", paths[2]) //paths[2]
				}
			}

			values = append(values, displayValue)
		}
		displayOrder.Values = values
		bp.DisplayOrder = displayOrder
	}
	return nil
}

func GetBlueprintConfigOverrideVariables(v interface{}, bp *BluePrintConfig) error {
	override_variables := v.([]interface{})
	aliases := blueprint_config_references.GetInstance()
	for _, currentVar := range override_variables {
		varOverrideMap := currentVar.(map[string]interface{})
		//create variable
		varName := varOverrideMap["name"].(string)
		vc, err := BuildVariableFromSchema(varOverrideMap, bp)
		if err != nil {
			return err
		}
		variables := varOverrideMap["variables"].(map[string]interface{})
		varInterpolation := make(map[string]string)
		for key, value := range vc.Variables {
			val := ApplyRefenceValue(value, aliases, bp)
			varInterpolation[key] = val
		}
		for key, value := range variables {
			/*
				We have to translate from <alias>.variables.<variable_name> to <module>.<variable_name>
				We check if a value of a variable interpolation element IS a reference (alias.variables.variable_name)
				If it is, we have to convert it to data_autocloud_blueprintconfig.<module>.<variable_name>
				data_autocloud_blueprint_config is an string reference used in the Mutation to get the value from other variable values
			*/
			//check if value is a reference
			val := ApplyRefenceValue(value.(string), aliases, bp)
			varInterpolation[key] = val
		}

		if len(vc.Value) > 0 {
			err := interpolation_utils.DetectInterpolation(vc.Value, varInterpolation)
			if err != nil {
				return err
			}
		}

		if len(vc.Default) > 0 {
			err := interpolation_utils.DetectInterpolation(vc.Default, varInterpolation)
			if err != nil {
				return err
			}
		}
		// check if the user defined variables for interpolation with an empty template
		if len(varInterpolation) > 0 {
			err := interpolation_utils.DetectInterpolation("", varInterpolation)
			if err != nil {
				return err
			}
		}

		bp.OverrideVariables[varName] = OverrideVariable{
			VariableName:      varName,
			VariableContent:   *vc,
			dirty:             false,
			InterpolationVars: varInterpolation,
			//UsedInHCL:       true,
		}

		// Conditionals
		conditionals, conditionalExists := varOverrideMap["conditional"].(*schema.Set)
		log.Printf("CONDITIONALS: %v \n", conditionals)
		if conditionalExists {
			if entry, ok := bp.OverrideVariables[varName]; ok {
				conditionals, err := getConditionals(conditionals, bp)
				if err != nil {
					return errors.New("GetBlueprintConfigFromSchema: Error accessing bp")
				}
				entry.Conditionals = conditionals
				bp.OverrideVariables[varName] = entry
			} else {
				return errors.New("GetBlueprintConfigFromSchema: Error accessing bp")
			}
		}
	}
	return nil
}

func BuildVariableFromSchema(rawSchema map[string]interface{}, bp *BluePrintConfig) (*VariableContent, error) {
	content := &VariableContent{}
	var requiredValues string
	requiredValuesInput, requiredValuesInputExist := rawSchema["required_values"]
	if requiredValuesInputExist {
		requiredValues = requiredValuesInput.(string)
	}

	content.DisplayName = rawSchema["display_name"].(string)
	content.HelperText = rawSchema["helper_text"].(string)
	content.Default = rawSchema["default"].(string)
	content.Variables = make(map[string]string, 0)
	content.RequiredValues = requiredValues

	if val, ok := rawSchema["variables"]; ok {
		var variablesMap = val.(map[string]interface{})
		content.Variables = utils.ConvertMap(variablesMap)
	}

	// Note: if it has a value, then it can NOT have form options "options"
	valueIsDefined := false
	value, valueExist := rawSchema["value"]
	valueStr, valueIsString := value.(string)
	valueIsDefined = valueStr != "" // NOTE: if the value is empty, we consider it as 'not defined'

	variableType := rawSchema["type"].(string)

	content.FormConfig.Type = variableType
	if valueExist && valueIsString && valueIsDefined {
		aliases := blueprint_config_references.GetInstance()
		//check if value is a reference
		content.Value = ApplyRefenceValue(valueStr, aliases, bp)
		return content, nil
	}

	content.Value = ""
	// variableContent with form options

	optionsFromSchema := rawSchema["options"].(*schema.Set)
	if len(optionsFromSchema.List()) > 1 {
		// it should be caught at schema check level - adding the check here to enforce it in case the schema changes
		return nil, errors.New("exactly one \"options\" must be defined")
	}

	if variableType == "shortText" && len(optionsFromSchema.List()) > 0 {
		return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrShortTextCantHaveOptions)
	}

	if variableType == "" && len(optionsFromSchema.List()) > 0 {
		variableType = LIST_TYPE
	}

	content.FormConfig = FormConfig{
		Type:            variableType,
		ValidationRules: make([]ValidationRule, 0),
		FieldOptions:    make([]FieldOption, 0),
	}
	if variableType == RADIO_TYPE || variableType == CHECKBOX_TYPE || variableType == LIST_TYPE {
		rawOptionsCluster := optionsFromSchema.List() // "options key in schema options {}" should always have 1 elem
		if len(optionsFromSchema.List()) == 1 {
			rawOptions := rawOptionsCluster[0].(map[string]interface{})
			optionSchema := rawOptions["option"].(*schema.Set)

			for _, option := range optionSchema.List() {
				options := option.(map[string]interface{})
				fieldOption := FieldOption{
					Label:   options["label"].(string),
					Value:   options["value"].(string),
					Checked: options["checked"].(bool),
				}
				content.FormConfig.FieldOptions = append(content.FormConfig.FieldOptions, fieldOption)
			}
		}
	}

	validationRulesList := rawSchema["validation_rule"].(*schema.Set).List()

	for _, validationRule := range validationRulesList {
		validationRuleMap := validationRule.(map[string]interface{})

		rule := validationRuleMap["rule"].(string)
		ruleValue := validationRuleMap["value"].(string)
		ruleScope := validationRuleMap["scope"].(string)

		if rule == "isRequired" && ruleValue != "" {
			return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrIsRequiredCantHaveValue)
		}
		if rule != "regex" && ruleScope != "" {
			return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrRegexOnlyCanHaveScope)
		}

		vr := ValidationRule{
			Rule:         rule,
			Value:        ruleValue,
			Scope:        ruleScope,
			ErrorMessage: validationRuleMap["error_message"].(string),
		}
		content.FormConfig.ValidationRules = append(content.FormConfig.ValidationRules, vr)
	}
	return content, nil
}

func getConditionals(varOverrideMap *schema.Set, bp *BluePrintConfig) ([]ConditionalConfig, error) {
	conditionalLen := varOverrideMap.Len()

	conditionals := make([]ConditionalConfig, 0)
	if conditionalLen == 0 {
		return conditionals, nil
	}
	conditionalsList := varOverrideMap.List()
	for _, conditional := range conditionalsList {
		conditionalMap := conditional.(map[string]interface{})
		conditionalContentMapList := conditionalMap["content"].(*schema.Set).List()
		if len(conditionalContentMapList) < 1 {
			continue
		}
		conditionalContentMap := conditionalContentMapList[0].(map[string]interface{})
		vc, err := BuildVariableFromSchema(conditionalContentMap, bp)
		if err != nil {
			return nil, err
		}
		c := ConditionalConfig{
			Source:          conditionalMap["source"].(string),
			Condition:       conditionalMap["condition"].(string),
			VariableContent: *vc,
		}

		conditionals = append(conditionals, c)
	}

	return conditionals, nil
}
