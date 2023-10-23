package blueprint_config

import (
	"fmt"
	"sort"
	"strings"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
)

// entry point for the transverse algorithm
func Transverse(root *BluePrintConfig) []generator.FormShape {
	emptyParent := BluePrintConfig{}
	root.Variables = bottomUpTraversal(root, &emptyParent)
	variables := root.Variables
	result := GetAllDisplayOrdersByBFS(root)
	sortedVariables := orderVariables(variables, result)
	return sortedVariables
}

// bottomUpTraversal traverses the BluePrintConfig tree in a bottom-up manner,
// applying overrides, omits, and merging variables from child nodes to parent nodes.
// it passes omits and overrides to children nodes if they are not applied to the current node.
// variables are processed in the node they are defined in.
// It returns a slice of FormShape variables representing the processed variables.
func bottomUpTraversal(node *BluePrintConfig, parent *BluePrintConfig) []generator.FormShape {
	if node.Id == "" {
		return []generator.FormShape{}
	}

	// Apply overrides for child nodes' variables to their respective nodes.
	// Remove these overrides from the current node's OverrideVariables map.
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

	// Handle omitted variables for child nodes and update node's OmitVariables accordingly.
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

	for childName := range node.Children {
		child := node.Children[childName]
		bottomUpTraversal(child, node)
	}

	// Apply overrides and omits from parent to current node
	finalVariables := processCurrentNode(*node)
	if finalVariables == nil {
		fmt.Println("finalVariables is nil, node id:", node.Id)
	}
	// here is where we should deal with ordering
	if node.Variables == nil {
		fmt.Println("node variables is nil, node id:", node.Id)
	}
	node.Variables = make([]generator.FormShape, 0)
	node.Variables = make([]generator.FormShape, len(finalVariables))
	copy(node.Variables, finalVariables)

	//Merge variables with existing ones in parent.Variables
	for _, variable := range finalVariables {
		varIdx, found := findVariableIndex(parent.Variables, variable.ID)
		if found {
			// Merge variables if found
			parentVariable := parent.Variables[varIdx]
			//originalVariable := parentVariable
			updatedVariable := variable
			// If existing variable  has less hops from its original node to the root, use it as original variable, the existing variable was processed before
			// it is important to check this, it avoids keeping track of the order of children in the graph
			if parentVariable.HopsFromNode < updatedVariable.HopsFromNode {
				parentVariable, updatedVariable = updatedVariable, parentVariable
			}
			mergedVariable := mergeVariables(parentVariable, updatedVariable)
			parent.Variables = updateVariable(parent.Variables, mergedVariable)
		} else {
			// Add new variable if not found
			parent.Variables = append(parent.Variables, variable)
		}
	}
	return parent.Variables
}

func processCurrentNode(node BluePrintConfig) []generator.FormShape {
	finalVariables := node.Variables

	for overrideVariableName, overrideVariable := range node.OverrideVariables {
		if variableIdexes, found := findVariableIndexes(node.Variables, overrideVariableName); found {
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
		if variableIdexes, found := findVariableIndexes(finalVariables, omitVariableID); found {
			for _, variableIdx := range variableIdexes {
				finalVariables[variableIdx].IsHidden = true
				finalVariables[variableIdx].UsedInHCL = false
				if finalVariables[variableIdx].IsOverriden {
					finalVariables[variableIdx].UsedInHCL = true
				}
			}
		}
	}

	//increment Hops in node
	for idx := range finalVariables {
		finalVariables[idx].HopsFromNode++
	}

	return finalVariables
}

func findVariableIndex(variables []generator.FormShape, id string) (int, bool) {
	for idx, variable := range variables {
		if variable.ID == id {
			return idx, true
		}
	}
	return -1, false
}

func findVariableIndexes(variables []generator.FormShape, overrideVariableName string) ([]int, bool) {
	foundIndexes := make([]int, 0)
	for idx, variable := range variables {
		localVarName := strings.Split(variable.ID, ".")[1]
		if localVarName == overrideVariableName {
			foundIndexes = append(foundIndexes, idx)
		}
	}
	if len(foundIndexes) > 0 {
		return foundIndexes, true
	}
	return []int{}, false
}

// Merge the fields from newVariable into existing
func mergeVariables(existing, newVariable generator.FormShape) generator.FormShape {
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
	existing.HopsFromNode = newVariable.HopsFromNode

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

func GetAllDisplayOrdersByBFS(node *BluePrintConfig) []DisplayOrder {
	result := make([]DisplayOrder, 0)
	queue := make([]*BluePrintConfig, 0)
	queue = append(queue, node)

	for len(queue) > 0 {
		// Pop the first element from the queue
		currentNode := queue[0]
		queue = queue[1:]

		// Process the current node
		processedDisplayOrder := processDisplayOrder(currentNode)
		result = append(result, processedDisplayOrder)
		//fmt.Println(processedDisplayOrder)

		// Add the children of the current node to the queue
		for _, child := range currentNode.Children {
			queue = append(queue, child)
		}
	}
	return result
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
			stack = append(stack, child)
		}
	}
}

func processDisplayOrder(bp *BluePrintConfig) DisplayOrder {
	variables := bp.DisplayOrder.Values
	priority := bp.DisplayOrder.Priority
	newVariables := make([]string, 0)
	for _, variable := range variables {
		parts := strings.Split(variable, ".")
		if len(parts) == 3 { //children.variables.variable
			childName := parts[0]
			variableName := parts[2]
			// find the child with the name childName
			childIdxes, found := findVariableIndexes(bp.Children[childName].Variables, variableName)
			if found {
				for _, childIdx := range childIdxes {
					varname := bp.Children[childName].Variables[childIdx].ID
					newVariables = append(newVariables, varname)
				}
			}
		}
		if len(parts) == 1 {
			for childName := range bp.Children {
				variableName := parts[0]
				// find the child with the name childName
				childIdxes, found := findVariableIndexes(bp.Children[childName].Variables, variableName)
				if found {
					for _, childIdx := range childIdxes {
						varname := bp.Children[childName].Variables[childIdx].ID
						newVariables = append(newVariables, varname)
					}
				}
			}
		}
	}

	return DisplayOrder{
		Values:   newVariables,
		Priority: priority,
	}
}

func orderVariables(variables []generator.FormShape, displayOrder []DisplayOrder) []generator.FormShape {
	//order variables alphanumerically
	sort.Slice(variables, func(i, j int) bool {
		return variables[i].ID < variables[j].ID
	})

	displayOrderDataWithValues := make([]DisplayOrder, 0)
	for _, displayOrderData := range displayOrder {
		if len(displayOrderData.Values) > 0 {
			displayOrderDataWithValues = append(displayOrderDataWithValues, displayOrderData)
		}
	}
	if len(displayOrderDataWithValues) == 0 {
		return variables
	}
	// Sort the displayOrderDataWithValues by priority, highest priority first
	sort.Slice(displayOrderDataWithValues, func(i, j int) bool {
		return displayOrderDataWithValues[i].Priority > displayOrderDataWithValues[j].Priority
	})

	// Order the variables according to the displayOrderDataWithValues
	sortedVariables := make([]generator.FormShape, 0)
	for _, displayOrderData := range displayOrderDataWithValues {
		sortedVariables = sortVariablesByIDOrder(variables, displayOrderData.Values)
	}

	return sortedVariables
}

func sortVariablesByIDOrder(variables []generator.FormShape, idOrder []string) []generator.FormShape {
	idIndex := make(map[string]int)

	// Create a map to store the indexes of IDs in the order list
	for idx, id := range idOrder {
		idIndex[id] = idx
	}

	// Separate the variables to be sorted from the rest
	var variablesToSort []generator.FormShape
	var otherVariables []generator.FormShape
	for _, variable := range variables {
		if _, found := idIndex[variable.ID]; found {
			variablesToSort = append(variablesToSort, variable)
		} else {
			otherVariables = append(otherVariables, variable)
		}
	}

	// Sort the variables to be sorted based on the given IDs order
	sort.SliceStable(variablesToSort, func(i, j int) bool {
		return idIndex[variablesToSort[i].ID] < idIndex[variablesToSort[j].ID]
	})

	// Merge the sorted variables with the rest of the variables
	return append(variablesToSort, otherVariables...)
}
