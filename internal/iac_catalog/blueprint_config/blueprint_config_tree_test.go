package blueprint_config_test

import (
	"fmt"
	"testing"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
)

func TestTreeTransversal(t *testing.T) {
	//     tree
	//      A
	//    /   \
	//   B     D
	//  /
	// C
	tree := blueprint_config.BluePrintConfig{
		Id: "root",
		Variables: []autocloudsdk.FormShape{
			{
				ID: "A",
			}},
		Children: []blueprint_config.BluePrintConfig{
			{
				Id: "root.1",
				Variables: []autocloudsdk.FormShape{{
					ID: "B",
				}},
				Children: []blueprint_config.BluePrintConfig{{
					Id: "root.1.1",
					Variables: []autocloudsdk.FormShape{{
						ID: "C",
					}},
				}},
			},
			{
				Id: "root.2",
				Variables: []autocloudsdk.FormShape{{
					ID: "D",
				}},
				Children: []blueprint_config.BluePrintConfig{},
			},
		},
	}
	form := blueprint_config.GetFormShape(tree)
	fmt.Println(form)
	expectedOrder := []string{"C", "B", "D", "A"}
	if len(form) != len(expectedOrder) {
		t.Fatalf("expected form has different length")
	}
	for i, v := range expectedOrder {
		if form[i].ID != v {
			t.Fatalf("expected form has different order")
		}
	}
	fmt.Println(form)
}
