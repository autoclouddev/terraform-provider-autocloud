package blueprint_config

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/apex/log"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/logger"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func FormShapeToMap(formShape []generator.FormShape) (map[string]string, error) {
	varsMap := make(map[string]string)

	for _, form := range formShape {
		varName, err := utils.GetVariableID(form.ID)
		if err != nil {
			return varsMap, err
		}
		varsMap[varName] = fmt.Sprintf("data_autocloud_blueprint_config.%s", form.ID)
	}
	return varsMap, nil
}

func GetFormShape(root BluePrintConfig) ([]generator.FormShape, error) {
	var log = logger.Create(log.Fields{"fn": "GetFormShape()"})
	str, _ := json.MarshalIndent(root, "", "    ")
	log.Debugf("root bc: %s", string(str))
	formShape, err := postOrderTransversal(&root)
	if err != nil {
		return []generator.FormShape{}, err
	}
	return sortFormShape(root, formShape), nil
}

func hasReference(ov string) bool {
	keyValue := strings.Split(ov, ".")
	if len(keyValue) != 3 {
		return false
	}
	if keyValue[1] != "variables" {
		return false
	}
	return true
}

// transverses the tree from leaves to root,
// passing current level variables to its parent after
// processing the current level overrides and generics
func postOrderTransversal(root *BluePrintConfig) ([]generator.FormShape, error) {
	var vars []generator.FormShape = root.Variables
	// first, make sure we override all variables that have reference
	// a reference consist in the following code
	// variable = {
	// 	name = "s3.variables.kmsKeyName"
	// 	...
	//   }
	// on name, if you split it by the point, the first part is the child name, the second the variable name
	// analyze overrides that have "."  <source>.<varname>
	keys := make([]string, 0)
	for k := range root.OverrideVariables {
		keys = append(keys, k)
	}
	overridesWithReference := filter(keys, hasReference)
	// look for the variable in root.Children[source].variables
	for _, key := range overridesWithReference {
		keyValue := strings.Split(key, ".")
		child := keyValue[0]
		varName := keyValue[2]
		idx := findIdx(root.Children[child].Variables, varName)
		if idx < 0 {
			return []generator.FormShape{}, fmt.Errorf("Variable Reference is not matching any children variable: %s", key)
		}
		// build override in place
		root.Children[child].Variables[idx] = BuildOverridenVariable(root.Children[child].Variables[idx], root.OverrideVariables[key])
		// delete from overrides
		delete(root.OverrideVariables, key)
		// remove from omits
		for i, omitName := range root.OmitVariables {
			if varName == omitName {
				root.OmitVariables = append(root.OmitVariables[:i], root.OmitVariables[i+1:]...)
				break
			}
		}
	}

	for _, v := range root.Children {
		v := v                                        // avoid implicit memory aliasing
		childrenvars, err := postOrderTransversal(&v) // this &v now the address of the inner v
		if err != nil {
			return []generator.FormShape{}, err
		}
		vars = append(vars, childrenvars...)
	}
	log.Debugf("current node omit vars, %s", root.OmitVariables)
	admittedVars := OmitVars(vars, root.OmitVariables, &root.OverrideVariables)
	log.Debugf("the [%v] addmited vars", admittedVars)
	log.Debugf("current override vars, %v", root.OverrideVariables)
	return OverrideVariables(admittedVars, root.OverrideVariables)
}

// vars => variables coming from leaves (for example: a s3 autocloud_module variables)
// omits => current blueprint config vars to omit (a var will be discarded in case there are no overrides in the current blueprint config)
// overrideVariables ==> current blueprint config var overrides
func OmitVars(vars []generator.FormShape, omits []string, overrideVariables *map[string]OverrideVariable) []generator.FormShape {
	addmittedVars := vars
	for _, omit := range omits {
		idx := findIdx(addmittedVars, omit)
		if idx == -1 {
			continue
		}
		omittedVar := addmittedVars[idx]
		omittedVar.IsHidden = true
		omittedVar.UsedInHCL = false

		if omittedVar.IsOverriden {
			omittedVar.UsedInHCL = true
		}
		addmittedVars[idx] = omittedVar
		//addmittedVars = remove(addmittedVars, idx)
		// if the blueprint config overrides an omitted variable, then it's an admitted var as we have to modify its behavior
		if overrideVariable, isVarOverriden := (*overrideVariables)[omit]; isVarOverriden {
			overrideVariable.IsHidden = true // we don't want to show omitted vars
			overrideVariable.UsedInHCL = true
			(*overrideVariables)[omit] = overrideVariable
			continue
		}
	}
	return addmittedVars
}

func findIdx(vars []generator.FormShape, refname string) int {
	for i, v := range vars {
		varName, varname := "", refname
		if hasReference(refname) {
			keyValue := strings.Split(refname, ".")
			varname = fmt.Sprintf("%v.%v", keyValue[0], keyValue[2])
			varName = v.ID
		} else {
			varId, err := utils.GetVariableID(v.ID)
			if err != nil {
				log.Debugf("the [%s] variable not found\n", varId)
				return -1
			}
			varName = varId
		}
		if varName == varname {
			log.Debugf("the [%s] variable was omitted\n", varName)
			return i
		}
	}
	log.Debugf("the [%s] omitted value not found in vars\n", refname)
	return -1
}

// vars => form shapes coming from leaves
// overrides => new vars definitions + modifications to vars from leaves
func OverrideVariables(vars []generator.FormShape, overrides map[string]OverrideVariable) ([]generator.FormShape, error) {
	var log = logger.Create(log.Fields{"fn": "OverrideVariables()"})
	usedOverrides := make(map[string][]string, 0)
	// transform all original Variables to its overrides
	for i, iacVar := range vars {
		varName, err := utils.GetVariableID(iacVar.ID)

		if err != nil {
			log.Debugf("WARNING: no variable ID found -> %v, evaluated value : %v", err, iacVar)
			// consider returning an error instead
			return []generator.FormShape{}, fmt.Errorf("%w -> %v, evaluated value : %v", ErrVariableNotFound, err, iacVar)
		}
		overrideVariableData, ok := overrides[varName]
		if ok {
			str, _ := json.MarshalIndent(overrideVariableData, "", "    ")
			log.Debugf("data -> %s", string(str))
			vars[i] = BuildOverridenVariable(iacVar, overrideVariableData)
		}

		// check if we already have overridden a variable
		if _, isAlreadyOverridden := usedOverrides[varName]; !isAlreadyOverridden {
			usedOverrides[varName] = make([]string, 0)
		}
		if !utils.Contains(usedOverrides[varName], iacVar.ID) {
			usedOverrides[varName] = append(usedOverrides[varName], iacVar.ID)
		}
	}
	for varName, overridenVarIds := range usedOverrides {
		log.Debugf("the [%v] variable overrides %d question(s): [%v]", varName, len(overridenVarIds), overridenVarIds)
		delete(overrides, varName)
	}
	// on this point only generics remain, no original variables
	for _, ov := range overrides {
		formVar := BuildGenericVariable(ov)
		vars = append(vars, formVar)
	}
	// sort questions to keep ordering consistent ??
	/*sort.Slice(vars, func(i, j int) bool {
		return vars[i].ID < vars[j].ID
	})*/
	return vars, nil
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func sortFormShape(root BluePrintConfig, formShape []generator.FormShape) []generator.FormShape {
	var log = logger.Create(log.Fields{"fn": "sortFormShape()"})
	// sort according display_order configuration
	// first we sort formShape alphanumeric
	formShape = SortVariableAlphanumeric(formShape)
	varsOrder, err := GetDisplayOrder(root)
	if err != nil {
		log.Errorf("Error to merge display order: %v", err)
		return formShape
	}
	// reverse varsOrder so we can move the variables to the front of the slice
	// and they will be sorted as we need
	varsOrder = reverseStringSlice(varsOrder)
	for _, varName := range varsOrder {
		isGenericVariable := !strings.Contains(varName, ".")
		// the order is reversed so we move to the front any match
		formShape = moveVariablesToFront(varName, formShape, isGenericVariable)
	}

	return formShape
}

func GetDisplayOrder(root BluePrintConfig) ([]string, error) {
	var log = logger.Create(log.Fields{"fn": "GetDisplayOrder()"})
	order, err := postDisplayOrderTransversal(&root)
	if err != nil {
		return []string{}, err
	}
	// we need to sort by priority asc
	sort.Slice(order, func(i, j int) bool {
		return order[i].Priority < order[j].Priority
	})
	var orderMerged = make([]string, 0)
	for _, displayOrder := range order {
		for _, varName := range displayOrder.Values {
			if !contains(orderMerged, varName) {
				orderMerged = append(orderMerged, varName)
			}
		}
	}
	log.Debugf("order: %v", orderMerged)
	return orderMerged, nil
}

func SortVariableAlphanumeric(formShape []generator.FormShape) []generator.FormShape {
	sort.Slice(formShape, func(i, j int) bool {
		pathsVar1 := strings.Split(formShape[i].ID, ".")
		pathsVar2 := strings.Split(formShape[j].ID, ".")
		return pathsVar1[len(pathsVar1)-1] < pathsVar2[len(pathsVar2)-1]
	})
	return formShape
}

// transverses the tree from leaves to root
func postDisplayOrderTransversal(root *BluePrintConfig) ([]DisplayOrder, error) {
	var order []DisplayOrder = make([]DisplayOrder, 0)

	order = append(order, root.DisplayOrder)

	for _, v := range root.Children {
		v := v                                                // avoid implicit memory aliasing
		childrenOrder, err := postDisplayOrderTransversal(&v) // this &v now the address of the inner v
		if err != nil {
			return []DisplayOrder{}, err
		}
		order = append(order, childrenOrder...)
	}
	return order, nil
}

func moveVariablesToFront(varName string, formShape []generator.FormShape, compareOnlyVariableName bool) []generator.FormShape {
	newFormShape := make([]generator.FormShape, 0)
	for _, variable := range formShape {
		currentVarName := variable.ID
		if compareOnlyVariableName {
			paths := strings.Split(variable.ID, ".")
			currentVarName = paths[len(paths)-1]
		}
		if currentVarName == varName {
			newFormShape = append([]generator.FormShape{variable}, newFormShape...)
		} else {
			newFormShape = append(newFormShape, variable)
		}
	}
	return newFormShape
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func reverseStringSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
