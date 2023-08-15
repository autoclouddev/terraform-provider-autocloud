package blueprint_config

import (
	"fmt"
	"strings"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

/*
Contains utility functions for translating references names to moudle names
reference name :=> <blueprintConfigChildName>.variables.<variableName>
module name :=> <ModuleName>.<variableName>
*/

// ApplyRefenceValue updates the value of a variable to an internal reference format.
// It checks if the variable value contains a reference <childName.variables.variableName> and extracts the reference components.
// The module name is looked up in the aliases data structure and, if not found, defaults to "generic".
// The function constructs a new value by combining the reference components with a predefined data reference prefix.
// "data_autocloud_blueprint_config.<moduleName>.<variableName>" is the data reference prefix, this is what the appy understand as a reference
// The updated value is returned, ensuring a standardized format for variable references.
func ApplyRefenceValue(variableValue string, aliases *blueprint_config_references.Data, bp *BluePrintConfig) string {
	if !utils.HasReference(variableValue) {
		return variableValue
	}
	moduleName := GetModuleNameFromVariable(variableValue, *aliases, bp)
	if len(moduleName) == 0 {
		moduleName = "generic"
	}
	dataRef := "data_autocloud_blueprint_config"
	reference := strings.Split(variableValue, ".")
	return fmt.Sprintf("%s.%s.%s", dataRef, moduleName, reference[2])
}

func GetModuleNameFromVariable(referenceName string, aliases blueprint_config_references.Data, bp *BluePrintConfig) string {
	if !utils.HasReference(referenceName) {
		return referenceName
	}
	name := findChildReferenceName(referenceName, aliases, bp)
	return name
}

// will look for the FIRST variable name that matches and returns its module name
// variables inside bp.Variables consist on <module_name>.<variable_name>
func findModuleNameInBlueprint(bp *BluePrintConfig, varName string) string {
	moduleName := findVariableNameInBlueprintVariables(bp, varName)
	if len(moduleName) > 0 {
		return moduleName
	}
	for _, c := range bp.Children {
		c := c
		moduleName := findVariableNameInBlueprintVariables(&c, varName)
		if len(moduleName) > 0 {
			return moduleName
		}
	}
	return ""
}

func findVariableNameInBlueprintVariables(bp *BluePrintConfig, varName string) string {
	for _, v := range bp.Variables {
		moduleAndVarNames := strings.Split(v.ID, ".")
		if moduleAndVarNames[1] == varName {
			return moduleAndVarNames[0]
		}
	}
	return ""
}

/*
Looks in the aliases MAP for the child name that matches the reference name
The map has the following structure

<childName>#<parentId> : <childId>
*/
func findChildReferenceName(referenceName string, aliases blueprint_config_references.Data, bp *BluePrintConfig) string {
	paths := strings.Split(referenceName, ".")
	strKey := paths[0]
	varName := paths[2]
	aliasKey := fmt.Sprintf("%s#%s", strKey, bp.Id)
	childrenId := aliases.GetValue(aliasKey)

	if len(childrenId) == 0 {
		return ""
	}
	var child *BluePrintConfig

	for _, chi := range bp.Children {
		chi := chi
		if chi.Id == childrenId {
			child = &chi
			break
		}
	}
	if child == nil {
		return ""
	}
	return findModuleNameInBlueprint(child, varName)
}

func GetVariableName(variableName string) string {
	parts := strings.Split(variableName, ".")
	if utils.HasReference(variableName) {
		return parts[0] + "." + parts[2]
	}
	// if len(parts) == 2 {
	// 	return parts[1]
	// }
	// return strings.Split(variableName, ".")[2]
	return variableName
}
