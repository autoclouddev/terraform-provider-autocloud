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
	}
	node.OmitVariables = removeIndexesFromOmits(node.OmitVariables, removedOmitsIndexes)

	for _, child := range node.Children {
		child := child // memory aliasing
		bottomUpTraversal(&child, node)
	}

	// Apply overrides and omits from parent to current node
	finalVariables := processNodeWithParent(*node)
	// here is where we should deal with ordering

	//TBD: sort variables by order

	//Merge variables with existing ones in node.Variables
	for _, variable := range finalVariables {
		//varIdx, found := findVariableIndex(parent.Variables, variable.ID)
		varIdx, found := findIndex(parent.Variables, func(v generator.FormShape) bool {
			currentVarName := strings.Split(v.ID, ".")[1]
			varName := strings.Split(variable.ID, ".")[1]
			return currentVarName == varName
		})
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
		//if variableIdx, found := findVariableIndex2(node.Variables, overrideVariableName); found {
		if variableIdx, found := findIndex(node.Variables, func(v generator.FormShape) bool {
			localVarName := strings.Split(v.ID, ".")[1]
			return localVarName == overrideVariableName
		}); found {
			// Apply overrides using BuildOverridenVariable function
			variable := node.Variables[variableIdx]
			variable = BuildOverridenVariable(variable, overrideVariable)
			finalVariables[variableIdx] = variable
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
		//if variableIdx, found := findVariableIndex2(finalVariables, omitVariableID); found {
		if variableIdx, found := findIndex(node.Variables, func(v generator.FormShape) bool {
			localVarName := strings.Split(v.ID, ".")[1]
			return localVarName == omitVariableID
		}); found {
			finalVariables[variableIdx].IsHidden = true
			finalVariables[variableIdx].UsedInHCL = false
			if finalVariables[variableIdx].IsOverriden {
				finalVariables[variableIdx].UsedInHCL = true
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

func findVariableIndex(variables []generator.FormShape, id string) (int, bool) {
	varName := strings.Split(id, ".")[1]
	for idx, variable := range variables {
		currentVarName := strings.Split(variable.ID, ".")[1]
		if currentVarName == varName {
			return idx, true
		}
	}
	return -1, false
}

func findVariableIndex2(variables []generator.FormShape, overrideVariableName string) (int, bool) {
	for idx, variable := range variables {
		localVarName := strings.Split(variable.ID, ".")[1]
		if localVarName == overrideVariableName {
			return idx, true
		}
	}
	return -1, false
}

func mergeVariables(existing, newVariable generator.FormShape) generator.FormShape {
	// Implement logic to merge existing and newVariable, e.g., combine attributes
	// Return the merged variable
	mergedVariable := existing
	return mergedVariable
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
