package blueprint_config

import (
	"encoding/json"
	"log"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func GetFormShape(root BluePrintConfig) []autocloudsdk.FormShape {
	str, _ := json.MarshalIndent(root, "", "    ")
	log.Printf("root bc: %s", string(str))
	return postOrderTransversal(&root)
}

// transverses the tree from leaves to root,
// passing current level variables to its parent after
// processing the current level overrides and generics
func postOrderTransversal(root *BluePrintConfig) []autocloudsdk.FormShape {
	var vars []autocloudsdk.FormShape = root.Variables
	for _, v := range root.Children {
		v := v                                   // avoid implicit memory aliasing
		childrenvars := postOrderTransversal(&v) // this &v now the address of the inner v
		vars = append(vars, childrenvars...)
	}
	log.Printf("current node omit vars, %s", root.OmitVariables)
	admittedVars := OmitVars(vars, root.OmitVariables)
	log.Printf("the [%v] addmited vars", admittedVars)
	log.Printf("current override vars, %v", root.OverrideVariables)
	return overrideVariables(admittedVars, root.OverrideVariables)
}

func OmitVars(vars []autocloudsdk.FormShape, omitts []string) []autocloudsdk.FormShape {
	addmittedVars := vars
	for _, omit := range omitts {
		idx := findIdx(addmittedVars, omit)
		if idx == -1 {
			continue
		}
		addmittedVars = remove(addmittedVars, idx)
	}
	return addmittedVars
}

func findIdx(vars []autocloudsdk.FormShape, varname string) int {
	for i, v := range vars {
		varName, err := utils.GetVariableID(v.ID)
		if err != nil {
			log.Printf("the [%s] variable not found", varName)
			return -1
		}
		if varName == varname {
			log.Printf("the [%s] variable was omitted", varName)
			return i
		}
	}
	log.Printf("the [%s] omitted value not found in vars", varname)
	return -1
}

func remove(slice []autocloudsdk.FormShape, s int) []autocloudsdk.FormShape {
	return append(slice[:s], slice[s+1:]...)
}

func overrideVariables(vars []autocloudsdk.FormShape, overrides map[string]OverrideVariable) []autocloudsdk.FormShape {
	// transform all original Variables to its overrides
	for i, iacVar := range vars {
		varName, err := utils.GetVariableID(iacVar.ID)

		if err != nil {
			log.Printf("WARNING: no variable ID found -> %v, evaluated value : %v", err, iacVar)
			return make([]autocloudsdk.FormShape, 0)
		}
		if overrideVariableData, ok := overrides[varName]; ok {
			vars[i] = buildOverridenVariable(iacVar, overrideVariableData)
			delete(overrides, varName)
		}
	}
	// on this point only generics remain, no original variables
	for _, ov := range overrides {
		formVar := buildGenericVariable(ov)
		vars = append(vars, formVar)
	}
	// sort questions to keep ordering consistent ??
	/*sort.Slice(vars, func(i, j int) bool {
		return vars[i].ID < vars[j].ID
	})*/
	return vars
}
