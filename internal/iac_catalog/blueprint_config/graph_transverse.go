package blueprint_config

import (
	"fmt"
	"sort"
	"strings"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
)

func Transverse(root *BluePrintConfig) []generator.FormShape {
	emptyParent := BluePrintConfig{}
	root.Variables = bottomUpTraversal(root, &emptyParent)
	fmt.Println("parent variables:", emptyParent.Variables)
	return root.Variables
}

func bottomUpTraversal(node *BluePrintConfig, parent *BluePrintConfig) []generator.FormShape {
	if node.Id == "" {
		return []generator.FormShape{}
	}

	for overrideName, overrideVariable := range node.OverrideVariables {
		parts := strings.Split(overrideName, ".")
		if len(parts) == 3 {
			childName := parts[0]
			stripedOverrideName := parts[2]
			child, found := node.Children[childName]
			if !found {
				fmt.Println("we should warn here, child not found, incorrect usage of childName.variables.varName")
				continue
			}
			if len(child.OverrideVariables) == 0 {
				child.OverrideVariables = make(map[string]OverrideVariable)
			}
			child.OverrideVariables[stripedOverrideName] = overrideVariable
			node.Children[childName] = child
			delete(node.OverrideVariables, overrideName)
		}
		// if len(parts) != 3 or != 1 { //print warning
		// if len(parts) == 1 {  pass this override to all children
	}

	removedOmitsIndexes := make([]int, 0)
	for idx, ommitName := range node.OmitVariables {
		parts := strings.Split(ommitName, ".")
		if len(parts) == 3 {
			childName := parts[0]
			stripedOmmitName := parts[2]
			child, found := node.Children[childName]
			if !found {
				fmt.Println("we should warn here, child not found, incorrect usage of childName.variables.varName")
				continue
			}
			child.OmitVariables = append(child.OmitVariables, stripedOmmitName)
			node.Children[childName] = child
			removedOmitsIndexes = append(removedOmitsIndexes, idx)
		}
		// if len(parts) != 3 or != 1 { //print warning

	}
	node.OmitVariables = removeIndexesFromOmits(node.OmitVariables, removedOmitsIndexes)

	for _, childName := range node.ChildrenOrder {
		child := node.Children[childName]
		bottomUpTraversal(&child, node)
	}

	// Apply overrides and omits from parent to current node
	finalVariables := processNodeWithParent(*node)
	// here is where we should deal with ordering

	//TBD: sort variables by order

	//Merge variables with existing ones in node.Variables
	for _, variable := range finalVariables {
		varIdx, found := findVariableIndex(parent.Variables, variable.ID)
		// varIdx, found := findIndex(parent.Variables, func(v generator.FormShape) bool {
		// 	currentVarName := strings.Split(v.ID, ".")[1]
		// 	varName := strings.Split(variable.ID, ".")[1]
		// 	return currentVarName == varName
		// })
		if found {
			// Merge variables if found
			existingVariable := parent.Variables[varIdx]
			mergedVariable := mergeVariables(existingVariable, variable)
			parent.Variables = updateVariable(parent.Variables, mergedVariable)
		} else {
			// Add new variable if not found
			parent.Variables = append(parent.Variables, variable)
		}
	}

	return parent.Variables
}

func processNodeWithParent(node BluePrintConfig) []generator.FormShape {
	finalVariables := node.Variables

	for overrideVariableName, overrideVariable := range node.OverrideVariables {
		if variableIdexes, found := findVariableIndex2(node.Variables, overrideVariableName); found {
			// if variableIdx, found := findIndex(node.Variables, func(v generator.FormShape) bool {
			// 	localVarName := strings.Split(v.ID, ".")[1]
			// 	return localVarName == overrideVariableName
			// }); found {
			// Apply overrides using BuildOverridenVariable function
			for _, variableIdx := range variableIdexes {
				variable := node.Variables[variableIdx]
				variable = BuildOverridenVariable(variable, overrideVariable)
				finalVariables[variableIdx] = variable
			}
		} else {
			// Generate new variable using BuildGenericVariable function

			overrideVariable.VariableName = overrideVariableName
			newVariable, err := BuildGenericVariable(overrideVariable)
			if err != nil {
				fmt.Println("Error building generic variable:", err)
				continue
			}
			finalVariables = append(finalVariables, newVariable)
		}
	}

	// Exclude variables specified in omitVariables
	for _, omitVariableID := range node.OmitVariables {
		if variableIdexes, found := findVariableIndex2(finalVariables, omitVariableID); found {
			//if variableIdx, found := findIndex(node.Variables, func(v generator.FormShape) bool {
			// 	localVarName := strings.Split(v.ID, ".")[1]
			// 	return localVarName == omitVariableID
			// }); found {
			for _, variableIdx := range variableIdexes {
				finalVariables[variableIdx].IsHidden = true
				finalVariables[variableIdx].UsedInHCL = false
				if finalVariables[variableIdx].IsOverriden {
					finalVariables[variableIdx].UsedInHCL = true
				}
			}

			// if the blueprint config overrides an omitted variable, then it's an admitted var as we have to modify its behavior
		}
	}

	return finalVariables
}

func findIndex(variables []generator.FormShape, condition func(generator.FormShape) bool) (int, bool) {
	for idx, variable := range variables {
		if condition(variable) {
			return idx, true
		}
	}
	return -1, false
}

// TODO: remove this function
func findVariableIndex(variables []generator.FormShape, id string) (int, bool) {
	//varName := strings.Split(id, ".")[1]
	for idx, variable := range variables {
		//currentVarName := strings.Split(variable.ID, ".")[1]
		if variable.ID == id {
			return idx, true
		}
	}
	return -1, false
}

// TODO: remove this function
func findVariableIndex2(variables []generator.FormShape, overrideVariableName string) ([]int, bool) {
	foundIndexes := make([]int, 0)
	for idx, variable := range variables {
		localVarName := strings.Split(variable.ID, ".")[1]
		if localVarName == overrideVariableName {
			foundIndexes = append(foundIndexes, idx)
			//return idx, true
		}
	}
	if len(foundIndexes) > 0 {
		return foundIndexes, true
	}
	return []int{}, false
}

func mergeVariables(existing, newVariable generator.FormShape) generator.FormShape {
	// Merge the fields from newVariable into existing
	existing.ID = newVariable.ID
	existing.Module = newVariable.Module
	existing.ModuleID = newVariable.ModuleID
	existing.FormQuestion = newVariable.FormQuestion
	existing.FieldDataType = newVariable.FieldDataType
	existing.FieldDefaultValue = newVariable.FieldDefaultValue
	existing.FieldValue = newVariable.FieldValue
	existing.AllowConsumerToEdit = newVariable.AllowConsumerToEdit
	existing.IsHidden = newVariable.IsHidden
	existing.UsedInHCL = newVariable.UsedInHCL
	existing.Conditionals = newVariable.Conditionals
	existing.RequiredValues = newVariable.RequiredValues

	// Merge InterpolationVars from newVariable into existing
	if existing.InterpolationVars == nil {
		existing.InterpolationVars = make(map[string]string)
	}
	for key, value := range newVariable.InterpolationVars {
		existing.InterpolationVars[key] = value
	}

	existing.IsOverriden = newVariable.IsOverriden
	return existing
}

func updateVariable(variables []generator.FormShape, updatedVariable generator.FormShape) []generator.FormShape {
	for i, variable := range variables {
		if variable.ID == updatedVariable.ID {
			variables[i] = updatedVariable
			return variables
		}
	}
	return append(variables, updatedVariable)
}

func removeIndexesFromOmits(s []string, indexes []int) []string {
	// Sort indexes in descending order to prevent out-of-bounds access
	sort.Sort(sort.Reverse(sort.IntSlice(indexes)))

	for _, idx := range indexes {
		s = removeElementAtIndex(s, idx)
	}
	return s
}

func removeElementAtIndex(s []string, idx int) []string {
	// Check if the index is out of bounds
	if idx < 0 || idx >= len(s) {
		return s
	}

	// Remove the element at idx by re-slicing the slice
	return append(s[:idx], s[idx+1:]...)
}

func BFS(node *BluePrintConfig) {
	queue := make([]*BluePrintConfig, 0)
	queue = append(queue, node)

	for len(queue) > 0 {
		// Pop the first element from the queue
		currentNode := queue[0]
		queue = queue[1:]

		// Process the current node
		fmt.Println(currentNode.DisplayOrder)

		// Add the children of the current node to the queue
		for _, child := range currentNode.Children {
			queue = append(queue, &child)
		}
	}
}

func DFS(node *BluePrintConfig) {
	stack := make([]*BluePrintConfig, 0)
	stack = append(stack, node)

	for len(stack) > 0 {
		// Pop the last element from the stack
		currentNode := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Process the current node
		fmt.Println(currentNode.DisplayOrder)

		// Add the children of the current node to the stack
		for _, child := range currentNode.Children {
			stack = append(stack, &child)
		}
	}
}

// func processDisplayOrder(bp *BluePrintConfig) {
// 	variables := bp.DisplayOrder.Values
// 	for idx, variable := range variables {
// 		parts := strings.Split(variable, ".")
// 		if len(parts) == 1 {
// 			// get all module names from all children
// 			moduleNames := getModuleNamesFromBlueprint(bp) // check which ones to sort first
// 			//remove idx from variables
// 			//add moduleNames + variable to variables
// 		}
// 	}
// }

func getModuleNamesFromBlueprint(bp *BluePrintConfig) []string {
	moduleNames := make([]string, 0)
	for _, child := range bp.Children {
		moduleNames = getModuleNamesFromFormShape(child.Variables)
	}
	return moduleNames
}

func getModuleNamesFromFormShape(formShape []generator.FormShape) []string {
	moduleNamesSet := make(map[string]bool)
	moduleNames := make([]string, 0)
	for _, variable := range formShape {
		moduleName := strings.Split(variable.ID, ".")[0]
		moduleNamesSet[moduleName] = true
	}

	for moduleName := range moduleNamesSet {
		moduleNames = append(moduleNames, moduleName)
	}
	return moduleNames
}
