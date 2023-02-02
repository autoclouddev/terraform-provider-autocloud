package blueprint_config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apex/log"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/logger"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func GetFormShape(root BluePrintConfig) ([]autocloudsdk.FormShape, error) {
	var log = logger.Create(log.Fields{"fn": "GetFormShape()"})
	str, _ := json.MarshalIndent(root, "", "    ")
	log.Debugf("root bc: %s", string(str))
	formShape, err := postOrderTransversal(&root)
	if err != nil {
		return []autocloudsdk.FormShape{}, err
	}
	return formShape, nil
}

// transverses the tree from leaves to root,
// passing current level variables to its parent after
// processing the current level overrides and generics
func postOrderTransversal(root *BluePrintConfig) ([]autocloudsdk.FormShape, error) {
	var vars []autocloudsdk.FormShape = root.Variables
	// first, make sure we override all variables that have reference
	// a reference consist in the following code
	// variable = {
	// 	name = "s3.variables.kmsKeyName"
	// 	...
	//   }
	// on name, if you split it by the point, the first part is the child name, the second the variable name
	// analyze overrides that have "."  <source>.<varname>
	hasReference := func(ov string) bool {
		keyValue := strings.Split(ov, ".")
		if len(keyValue) != 3 {
			return false
		}
		if keyValue[1] != "variables" {
			return false
		}
		return true
	}
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
			return []autocloudsdk.FormShape{}, fmt.Errorf("Variable Reference is not matching any children variable: %s", key)
		}
		// build override in place
		root.Children[child].Variables[idx] = BuildOverridenVariable(root.Children[child].Variables[idx], root.OverrideVariables[key])
		// delete from overrides
		delete(root.OverrideVariables, key)
	}

	for _, v := range root.Children {
		v := v                                        // avoid implicit memory aliasing
		childrenvars, err := postOrderTransversal(&v) // this &v now the address of the inner v
		if err != nil {
			return []autocloudsdk.FormShape{}, err
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
func OmitVars(vars []autocloudsdk.FormShape, omits []string, overrideVariables *map[string]OverrideVariable) []autocloudsdk.FormShape {
	addmittedVars := vars
	for _, omit := range omits {
		idx := findIdx(addmittedVars, omit)
		if idx == -1 {
			continue
		}
		omittedVar := addmittedVars[idx]
		omittedVar.IsHidden = true
		omittedVar.UsedInHCL = false
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

func findIdx(vars []autocloudsdk.FormShape, varname string) int {
	for i, v := range vars {
		varName, err := utils.GetVariableID(v.ID)
		if err != nil {
			log.Debugf("the [%s] variable not found\n", varName)
			return -1
		}
		if varName == varname {
			log.Debugf("the [%s] variable was omitted\n", varName)
			return i
		}
	}
	log.Debugf("the [%s] omitted value not found in vars\n", varname)
	return -1
}

// vars => form shapes coming from leaves
// overrides => new vars definitions + modifications to vars from leaves
func OverrideVariables(vars []autocloudsdk.FormShape, overrides map[string]OverrideVariable) ([]autocloudsdk.FormShape, error) {
	var log = logger.Create(log.Fields{"fn": "OverrideVariables()"})
	usedOverrides := make(map[string][]string, 0)
	// transform all original Variables to its overrides
	for i, iacVar := range vars {
		varName, err := utils.GetVariableID(iacVar.ID)

		if err != nil {
			log.Debugf("WARNING: no variable ID found -> %v, evaluated value : %v", err, iacVar)
			// consider returning an error instead
			return []autocloudsdk.FormShape{}, fmt.Errorf("%w -> %v, evaluated value : %v", ErrVariableNotFound, err, iacVar)
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
		log.Debugf("ok? %v", ok)
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
