package blueprint_config

import (
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

func GetFormShape(root BluePrintConfig) []autocloudsdk.FormShape {
	variables := make([]autocloudsdk.FormShape, 0)
	postOrderTransversal(&root, &variables)
	return variables
}

func postOrderTransversal(root *BluePrintConfig, vars *[]autocloudsdk.FormShape) {
	if root == nil {
		return
	}
	for _, v := range root.Children {
		v := v                         // avoid implicit memory aliasing
		postOrderTransversal(&v, vars) // this &v now the address of the inner v
	}
	*vars = append(*vars, root.Variables...)
}
