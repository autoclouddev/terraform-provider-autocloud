package blueprint_config_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/go-faker/faker/v4"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func fakeFormShape(moduleName string) autocloudsdk.FormShape {
	a := autocloudsdk.FormShape{}
	err := faker.FakeData(&a)
	if err != nil {
		log.Fatalln(err)
	}
	id := fmt.Sprintf("%s.%s", moduleName, faker.Word())
	a.ID = id
	return a
}

func TestTreeTransversal(t *testing.T) {
	//     tree
	//      A
	//    /   \
	//   B     D
	//  /
	// C
	AVariables := []autocloudsdk.FormShape{fakeFormShape("s3")}
	BVariables := []autocloudsdk.FormShape{fakeFormShape("s3"), fakeFormShape("s3")}
	CVariables := []autocloudsdk.FormShape{fakeFormShape("s3"), fakeFormShape("s3")}
	DVariables := []autocloudsdk.FormShape{fakeFormShape("s3")}

	tree := blueprint_config.BluePrintConfig{
		Id:        "root",
		Variables: AVariables,
		Children: []blueprint_config.BluePrintConfig{
			{
				Id:        "root.1",
				Variables: BVariables,
				Children: []blueprint_config.BluePrintConfig{{
					Id:        "root.1.1",
					Variables: CVariables,
				}},
			},
			{
				Id:        "root.2",
				Variables: DVariables,
				Children:  []blueprint_config.BluePrintConfig{},
			},
		},
	}
	form := blueprint_config.GetFormShape(tree)
	//expectedOrder := []string{"A", "B", "C", "D"}
	expectedOrder := append(append(append(AVariables, BVariables...), CVariables...), DVariables...)

	if len(form) != len(expectedOrder) {
		t.Fatalf("expected form has different length")
	}

	fmt.Printf("Avariables: %s\n", printFormShapeVarsIds(AVariables))
	fmt.Printf("Bvariables: %s\n", printFormShapeVarsIds(BVariables))
	fmt.Printf("Cvariables: %s\n", printFormShapeVarsIds(CVariables))
	fmt.Printf("Dvariables: %s\n", printFormShapeVarsIds(DVariables))

	fmt.Printf("expected: %s\n", printFormShapeVarsIds(expectedOrder))
	fmt.Printf("got: %s\n", printFormShapeVarsIds(form))

	for i, v := range expectedOrder {
		if form[i].ID != v.ID {
			t.Fatalf("expected form has different order")
		}
	}
}

func printFormShapeVarsIds(form []autocloudsdk.FormShape) string {
	r := ""
	for _, v := range form {
		r += fmt.Sprintf(" %s", v.ID)
	}
	return r
}

func TestOmitVars(t *testing.T) {
	vars := make([]autocloudsdk.FormShape, 0)
	varCount := 10
	pickedCount := 3
	for i := 0; i < varCount; i++ {
		vars = append(vars, fakeFormShape("s3"))
	}
	omits := []string{}
	for i := 0; i < pickedCount; i++ {
		pick := rand.Intn(len(vars) - i)
		varname, _ := utils.GetVariableID(vars[pick].ID)
		omits = append(omits, varname)
	}
	result := blueprint_config.OmitVars(vars, omits)
	fmt.Printf("original: %s\n", printFormShapeVarsIds(vars))
	fmt.Print("OMITTED\n")
	fmt.Println(omits)
	fmt.Printf("result: %s\n", printFormShapeVarsIds(result))

	if len(result) != varCount-pickedCount {
		t.Fatalf("vars were not omitted")
	}
}

func TestJsonUnmarshallOverride(t *testing.T) {
	bp := createBp()

	v, err := json.Marshal(bp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(v))
}

func createBp() *blueprint_config.BluePrintConfig {
	//map[bucket:{bucket  Bucket name Set the bucket name {shortText [] [{isRequired  invalid}]}}]
	bp := &blueprint_config.BluePrintConfig{
		OverrideVariables: make(map[string]blueprint_config.OverrideVariable, 0),
	}
	fieldOptions := make([]blueprint_config.FieldOption, 0)
	fmConfig := blueprint_config.FormConfig{
		Type:         "shortText",
		FieldOptions: fieldOptions,
	}
	ov := blueprint_config.OverrideVariable{
		VariableName: "bucket",
		Value:        "Bucket name",
		DisplayName:  "display",
		HelperText:   "help",
		FormConfig:   fmConfig,
	}
	bp.OverrideVariables["bucket"] = ov
	return bp
}
