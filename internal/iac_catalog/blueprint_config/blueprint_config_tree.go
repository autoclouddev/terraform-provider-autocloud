package blueprint_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/apex/log"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/logger"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func GetFormShape(root BluePrintConfig) ([]generator.FormShape, error) {
	var log = logger.Create(log.Fields{"fn": "GetFormShape()"})
	str, _ := json.MarshalIndent(root, "", "    ")
	log.Debugf("root bc: %s", string(str))
	formShape, err := PostOrderTransversal(&root)
	if err != nil {
		return []generator.FormShape{}, err
	}
	return sortFormShape(root, formShape), nil
}

// transverses the tree from leaves to root,
// passing current level variables to its parent after
// processing the current level overrides and generics
func PostOrderTransversal(root *BluePrintConfig) ([]generator.FormShape, error) {
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
	overridesWithReference := filter(keys, utils.HasReference)
	// look for the variable in root.Children[source].variables
	for _, key := range overridesWithReference {
		keyValue := strings.Split(key, ".")
		varName := keyValue[2]
		for cindex := 0; cindex <= len(root.Children)-1; cindex++ {
			matches := FindIdx(root.Children[cindex].Variables, key, root)
			// if len(matches) < 1 {
			// 	return []generator.FormShape{}, fmt.Errorf("Variable Reference is not matching any children variable: %s", key)
			// }
			for _, idx := range matches {
				// build override in place
				root.Children[cindex].Variables[idx] = BuildOverridenVariable(root.Children[cindex].Variables[idx], root.OverrideVariables[key])
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
		}
	}

	for _, v := range root.Children {
		v := v                                        // avoid implicit memory aliasing
		childrenvars, err := PostOrderTransversal(&v) // this &v now the address of the inner v
		if err != nil {
			return []generator.FormShape{}, err
		}
		vars = append(vars, childrenvars...)
	}

	log.Debugf("current node omit vars, %s", root.OmitVariables)
	admittedVars := OmitVars(vars, root.OmitVariables, &root.OverrideVariables, root)
	log.Debugf("the [%v] addmited vars", admittedVars)
	log.Debugf("current override vars, %v", root.OverrideVariables)
	return OverrideVariables(admittedVars, root.OverrideVariables, root)
}

// vars => variables coming from leaves (for example: a s3 autocloud_module variables)
// omits => current blueprint config vars to omit (a var will be discarded in case there are no overrides in the current blueprint config)
// overrideVariables ==> current blueprint config var overrides
func OmitVars(vars []generator.FormShape, omits []string, overrideVariables *map[string]OverrideVariable, bp *BluePrintConfig) []generator.FormShape {
	addmittedVars := vars
	for _, omit := range omits {
		matches := FindIdx(addmittedVars, omit, bp)
		for _, idx := range matches {
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
	}
	return addmittedVars
}

// vars => form shapes coming from leaves
// overrides => new vars definitions + modifications to vars from leaves
func OverrideVariables(vars []generator.FormShape, overrides map[string]OverrideVariable, bp *BluePrintConfig) ([]generator.FormShape, error) {
	var log = logger.Create(log.Fields{"fn": "OverrideVariables()"})
	usedOverrides := make(map[string][]string, 0)
	// transform all original Variables to its overrides
	for overrideName, overrideData := range overrides {
		matches := FindIdx(vars, overrideName, bp)
		for _, idx := range matches {
			str, _ := json.MarshalIndent(overrideData, "", "    ")
			log.Debugf("data -> %s", string(str))
			vars[idx] = BuildOverridenVariable(vars[idx], overrideData)

			// check if we already have overridden a variable
			if _, isAlreadyOverridden := usedOverrides[overrideName]; !isAlreadyOverridden {
				usedOverrides[overrideName] = make([]string, 0)
			}
			usedOverrides[overrideName] = append(usedOverrides[overrideName], overrideName)
		}
	}
	for varName, overridenVarIds := range usedOverrides {
		log.Debugf("the [%v] variable overrides %d question(s): [%v]", varName, len(overridenVarIds), overridenVarIds)
		delete(overrides, varName)
	}
	// on this point only generics remain, no original variables
	for _, ov := range overrides {
		formVar, err := BuildGenericVariable(ov)
		if err != nil {
			return nil, err
		}
		vars = append(vars, formVar)
	}

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
	fmt.Println("VARS ORDER ->", varsOrder)
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
		// check against the alias
		fmt.Println("currentVarName ->", currentVarName)
		fmt.Println("varName ->", varName)
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

// Find index for a variable given its name
func FindIdx(vars []generator.FormShape, refname string, bp *BluePrintConfig) []int {
	matches := make([]int, 0)
	for i, v := range vars {
		varName, varname := "", refname
		var err error
		if utils.HasReference(refname) {
			varname, err = GetVariableReferenceID(refname, bp)
			varName = v.ID
			if err != nil {
				log.Debugf("the [%s] reference variable not found\n", varName)
			}
		} else {
			varId, err := utils.GetVariableID(v.ID)
			if err != nil {
				log.Debugf("the [%s] variable not found\n", varId)
			}
			varName = varId
		}
		if varName == varname {
			log.Debugf("the [%s] variable was omitted\n", varName)
			matches = append(matches, i)
		}
	}
	log.Debugf("the [%s] omitted value not found in vars\n", refname)
	return matches
}

// variables id follow the pattern "<alias>.variables.<variable name>""
func GetVariableReferenceID(variableKey string, bp *BluePrintConfig) (string, error) {
	var aliases = blueprint_config_references.GetInstance()
	fmt.Println("ALIASES: ", aliases.ToString())
	keyValue := strings.Split(variableKey, ".")
	moduleName := GetModuleNameFromVariable(variableKey, *aliases, bp)
	if utils.HasReference(variableKey) && len(keyValue) == 3 && len(moduleName) > 0 {
		fmt.Println("found moduleName: ", moduleName)
		return fmt.Sprintf("%v.%v", moduleName, keyValue[2]), nil
	}
	return "", errors.New("Invalid Key")
}
