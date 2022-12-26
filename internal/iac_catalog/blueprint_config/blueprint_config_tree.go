package blueprint_config

import (
	"log"
	"sort"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func GetFormShape(root BluePrintConfig) []autocloudsdk.FormShape {
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
	admittedVars := omitVars(vars, root.OmitVariables)
	return overrideVariables(admittedVars, root.OverrideVariables)
}

func omitVars(vars []autocloudsdk.FormShape, omitts []string) []autocloudsdk.FormShape {
	addmittedVars := make([]autocloudsdk.FormShape, 0)
	for _, iacModuleVar := range vars {
		varName, err := utils.GetVariableID(iacModuleVar.ID)

		if err != nil {
			log.Printf("WARNING: no variable ID found -> %v, evaluated value : %v", err, iacModuleVar)
			return make([]autocloudsdk.FormShape, 0)
		}

		// omit vars
		if !utils.Contains(omitts, varName) {
			addmittedVars = append(addmittedVars, iacModuleVar)
			log.Printf("the [%s] variable was addmitted", varName)
		} else {
			log.Printf("the [%s] variable was omitted", varName)
		}
	}

	return addmittedVars
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
	// sort questions to keep ordering consistent
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].ID < vars[j].ID
	})
	return vars
}
