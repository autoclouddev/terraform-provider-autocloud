package blueprint_config_test

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"path"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	acctest "gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func fakeFormShape() generator.FormShape {
	moduleName := "s3"
	a := generator.FormShape{}
	err := faker.FakeData(&a)
	if err != nil {
		log.Fatalln(err)
	}
	id := fmt.Sprintf("%s.%s", moduleName, faker.Word())
	a.ID = id
	return a
}

func fakeFormShapeWithID(id string) generator.FormShape {
	a := fakeFormShape()
	a.ID = id
	return a
}

func TestTreeTransversal(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "INFO", // to see the DEBUG logs
	})
	//     tree
	//      A
	//    /   \
	//   B     D
	//  /
	// C
	AVariables := []generator.FormShape{fakeFormShape()}
	BVariables := []generator.FormShape{fakeFormShape(), fakeFormShape()}
	CVariables := []generator.FormShape{fakeFormShape(), fakeFormShape()}
	DVariables := []generator.FormShape{fakeFormShape()}

	cVar := CVariables[1]
	cVar.ID = "cloudfront.name"
	overrideInB, err := CreateFakeOverride()
	if err != nil {
		t.Fatalf(err.Error())
	}
	overrideInB.DisplayName = "cloudfront.name is overrided!"
	overrideInB.VariableName = "C.variables.name"
	CVariables[1] = cVar
	tree := blueprint_config.BluePrintConfig{
		Id:        "A",
		Variables: AVariables,
		Children: []blueprint_config.BluePrintConfig{
			{
				Id:        "B.1",
				Variables: BVariables,
				OverrideVariables: map[string]blueprint_config.OverrideVariable{
					overrideInB.VariableName: *overrideInB,
				},
				Children: []blueprint_config.BluePrintConfig{
					{
						Id:        "C.1",
						Variables: CVariables,
					}},
			},
			{
				Id:        "D.1",
				Variables: DVariables,
				Children:  []blueprint_config.BluePrintConfig{},
			},
		},
	}

	form, _ := blueprint_config.GetFormShape(tree)
	//expectedOrder := []string{"A", "B", "C", "D"}
	expectedOrder := append(append(append(AVariables, BVariables...), CVariables...), DVariables...)
	expectedOrder = blueprint_config.SortVariableAlphanumeric(expectedOrder)

	if len(form) != len(expectedOrder) {
		t.Fatalf("expected form has different length, len(form): %v  len(expectedOrder): %v", len(form), len(expectedOrder))
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
		if form[i].ID == overrideInB.VariableName && form[i].FormQuestion.FieldLabel != overrideInB.DisplayName {
			t.Fatalf("Override by referece did not work")
		}
	}

	t.Cleanup(closer)
}

func printFormShapeVarsIds(form []generator.FormShape) string {
	r := ""
	for _, v := range form {
		r += fmt.Sprintf(" %s", v.ID)
	}
	return r
}

func TestOmitVars(t *testing.T) {
	vars := make([]generator.FormShape, 0)
	varCount := 10
	pickedCount := 3
	for i := 0; i < varCount; i++ {
		fk := fakeFormShape()
		fk.IsHidden = false
		vars = append(vars, fk)
	}
	omits := []string{}
	for i := 0; i < pickedCount; i++ {
		pick, err := rand.Int(rand.Reader, big.NewInt(int64(len(vars)-i)))
		if err != nil {
			panic(err)
		}
		varname, _ := utils.GetVariableID(vars[pick.Int64()].ID)
		omits = append(omits, varname)
	}
	result := blueprint_config.OmitVars(vars, omits, &(map[string]blueprint_config.OverrideVariable{}), &blueprint_config.BluePrintConfig{})
	fmt.Printf("original: %s\n", printFormShapeVarsIds(vars))
	fmt.Print("OMITTED\n")
	fmt.Println(omits)
	fmt.Printf("result: %s\n", printFormShapeVarsIds(result))

	hiddenVars := 0
	for _, v := range result {
		fmt.Println(v.ID, v.IsHidden)
		if !v.IsHidden {
			hiddenVars++
		}
	}

	if hiddenVars != varCount-pickedCount {
		t.Fatalf("vars were not omitted")
	}
}

func TestOmitReferenceVars(t *testing.T) {
	vars := make([]generator.FormShape, 0)

	s3Tags := fakeFormShape()
	s3Tags.ID = "s3.tags"
	vars = append(vars, s3Tags)

	cfTags := fakeFormShape()
	cfTags.ID = "cloudfront.tags"
	vars = append(vars, cfTags)

	omits := []string{"cloudfront.variables.tags"}
	omits = append(omits, "cloudfront.variables.tags")
	result := blueprint_config.OmitVars(vars, omits, &(map[string]blueprint_config.OverrideVariable{}), &blueprint_config.BluePrintConfig{})

	hiddenVars := 0
	for _, v := range result {
		if !v.IsHidden {
			hiddenVars++
		}
	}

	assert.Equal(t, 0, hiddenVars) // one omitted variable
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
	}
	ov.Value = "Bucket name"
	ov.DisplayName = "display"
	ov.HelperText = "help"
	ov.FormConfig = fmConfig
	bp.OverrideVariables["bucket"] = ov
	return bp
}

func TestGetDisplayOrder(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "INFO", // to see the DEBUG logs
	})
	//     tree
	//      A
	//    /   \
	//   B     D
	//  /
	// C
	AVariables := []generator.FormShape{fakeFormShape()}
	BVariables := []generator.FormShape{fakeFormShape(), fakeFormShape()}
	CVariables := []generator.FormShape{fakeFormShape(), fakeFormShape()}
	DVariables := []generator.FormShape{fakeFormShape()}

	tree := blueprint_config.BluePrintConfig{
		Id:        "A",
		Variables: AVariables,
		DisplayOrder: blueprint_config.DisplayOrder{
			Priority: 1,
			Values:   []string{"name"},
		},
		Children: []blueprint_config.BluePrintConfig{
			{
				Id:        "B.1",
				Variables: BVariables,
				DisplayOrder: blueprint_config.DisplayOrder{
					Priority: 0,
					Values:   []string{"name", BVariables[0].ID},
				},
				Children: []blueprint_config.BluePrintConfig{
					{
						Id:        "C.1",
						Variables: CVariables,
						DisplayOrder: blueprint_config.DisplayOrder{
							Priority: 2,
							Values:   []string{"description", CVariables[1].ID},
						},
					}},
			},
			{
				Id:        "D.1",
				Variables: DVariables,
				Children:  []blueprint_config.BluePrintConfig{},
			},
		},
	}

	displayOrder, _ := blueprint_config.GetDisplayOrder(tree)
	expectedDisplayOrder := []string{"name", BVariables[0].ID, "description", CVariables[1].ID}
	fmt.Printf("got: %s\n", displayOrder)

	if len(displayOrder) != len(expectedDisplayOrder) {
		t.Fatalf("expected display order has different length, len(displayOrder): %v  len(expectedDisplayOrder): %v", len(displayOrder), len(expectedDisplayOrder))
	}

	for i, v := range expectedDisplayOrder {
		if displayOrder[i] != v {
			t.Fatalf("expected displayOrder has different order")
		}
	}

	t.Cleanup(closer)
}

// No display_order has been defined in a blueprint
func TestFormShapeSortedCase1(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "INFO", // to see the DEBUG logs
	})

	AVariables := []generator.FormShape{
		fakeFormShapeWithID("s3.namespaces"),
		fakeFormShapeWithID("s3.name"),
		fakeFormShapeWithID("s3.description"),
		fakeFormShapeWithID("s3.amount"),
		fakeFormShapeWithID("s3.owner"),
	}

	tree := blueprint_config.BluePrintConfig{
		Id:        "A",
		Variables: AVariables,
	}

	form, _ := blueprint_config.GetFormShape(tree)
	expectedOrder := []string{"s3.amount", "s3.description", "s3.name", "s3.namespaces", "s3.owner"}

	for i, v := range expectedOrder {
		if form[i].ID != v {
			t.Fatalf("expected form has different order")
		}
	}

	t.Cleanup(closer)
}

// One display_order has been defined and blueprint has 1 module
func TestFormShapeSortedCase2(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "INFO", // to see the DEBUG logs
	})

	AVariables := []generator.FormShape{
		fakeFormShapeWithID("s3.namespaces"),
		fakeFormShapeWithID("s3.name"),
		fakeFormShapeWithID("s3.description"),
		fakeFormShapeWithID("s3.amount"),
		fakeFormShapeWithID("s3.owner"),
	}

	tree := blueprint_config.BluePrintConfig{
		Id:        "A",
		Variables: AVariables,
		DisplayOrder: blueprint_config.DisplayOrder{
			Priority: 0,
			Values:   []string{"name", "s3.description"},
		},
	}

	form, _ := blueprint_config.GetFormShape(tree)
	expectedOrder := []string{"s3.name", "s3.description", "s3.amount", "s3.namespaces", "s3.owner"}

	for i, v := range expectedOrder {
		if form[i].ID != v {
			t.Fatalf("expected form has different order")
		}
	}

	t.Cleanup(closer)
}

// One display_order has been defined and blueprint has 2 modules
func TestFormShapeSortedCase3(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "INFO", // to see the DEBUG logs
	})
	//     tree
	//      A
	//    /
	//   B
	AVariables := []generator.FormShape{
		fakeFormShapeWithID("s3.namespaces"),
		fakeFormShapeWithID("s3.name"),
		fakeFormShapeWithID("s3.description"),
		fakeFormShapeWithID("s3.amount"),
		fakeFormShapeWithID("s3.owner"),
	}
	BVariables := []generator.FormShape{
		fakeFormShapeWithID("kms.organization"),
		fakeFormShapeWithID("kms.name"),
		fakeFormShapeWithID("kms.description"),
		fakeFormShapeWithID("kms.tags"),
	}

	tree := blueprint_config.BluePrintConfig{
		Id:        "A",
		Variables: AVariables,
		DisplayOrder: blueprint_config.DisplayOrder{
			Priority: 0,
			Values:   []string{"name", "kms.description"},
		},
		Children: []blueprint_config.BluePrintConfig{
			{
				Id:        "B1",
				Variables: BVariables,
			},
		},
	}

	form, _ := blueprint_config.GetFormShape(tree)
	expectedOrder := []string{"kms.name", "s3.name", "kms.description", "s3.amount", "s3.description", "s3.namespaces", "kms.organization", "s3.owner", "kms.tags"}

	for i, v := range expectedOrder {
		if form[i].ID != v {
			t.Fatalf("expected form has different order")
		}
	}

	t.Cleanup(closer)
}

// Multiple display_order properties are set in multiple blueprint_configs
func TestFormShapeSortedCase4(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "INFO", // to see the DEBUG logs
	})
	//     tree
	//      A
	//    /
	//   B
	AVariables := []generator.FormShape{
		fakeFormShapeWithID("s3.namespaces"),
		fakeFormShapeWithID("s3.name"),
		fakeFormShapeWithID("s3.description"),
		fakeFormShapeWithID("s3.amount"),
		fakeFormShapeWithID("s3.owner"),
	}
	BVariables := []generator.FormShape{
		fakeFormShapeWithID("kms.organization"),
		fakeFormShapeWithID("kms.name"),
		fakeFormShapeWithID("kms.description"),
		fakeFormShapeWithID("kms.tags"),
	}

	tree := blueprint_config.BluePrintConfig{
		Id:        "A",
		Variables: AVariables,
		DisplayOrder: blueprint_config.DisplayOrder{
			Priority: 0,
			Values:   []string{"name", "s3.description"},
		},
		Children: []blueprint_config.BluePrintConfig{
			{
				Id:        "B1",
				Variables: BVariables,
				DisplayOrder: blueprint_config.DisplayOrder{
					Priority: 1,
					Values:   []string{"name", "kms.description", "organization"},
				},
			},
		},
	}

	form, _ := blueprint_config.GetFormShape(tree)
	expectedOrder := []string{"kms.name", "s3.name", "s3.description", "kms.description", "kms.organization", "s3.amount", "s3.namespaces", "s3.owner", "kms.tags"}

	for i, v := range expectedOrder {
		if form[i].ID != v {
			t.Fatalf("expected form has different order")
		}
	}

	t.Cleanup(closer)
}

func TestGetFormShape(t *testing.T) {
	testData, err := utils.LoadData[blueprint_config.BluePrintConfig](path.Join(FIXTURES_FOLDER, "blueprint_config.json"))
	assert.Nil(t, err)

	aliases := blueprint_config_references.GetInstance()
	aliases.SetValue("s3", "s3bucket")
	aliases.SetValue("policy", "s3bucketpolicy")
	aliases.SetValue("route53", "route53")
	aliases.SetValue("cloudfront", "cloudfront")
	aliases.SetValue("s3_bucket", "s3bucket")
	aliases.SetValue("s3_bucket_policy", "s3bucketpolicy")
	aliases.SetValue("kms", "kmskey")

	_, err = blueprint_config.GetFormShape(testData)
	assert.Nil(t, err)
}

func TestGetFormShapeError(t *testing.T) {
	testData, err := utils.LoadData[blueprint_config.BluePrintConfig](path.Join(FIXTURES_FOLDER, "blueprint_config.json"))
	assert.Nil(t, err)

	aliases := blueprint_config_references.GetInstance()
	aliases.SetValue("s3", "s3bucket")
	aliases.SetValue("policy", "s3bucketpolicy")
	aliases.SetValue("route53", "route53")
	aliases.SetValue("cloudfront", "cloudfront")
	aliases.SetValue("s3_bucket", "s3bucket")
	aliases.SetValue("s3_bucket_policy", "s3bucketpolicy")
	aliases.SetValue("kms", "kmskey")

	_, err = blueprint_config.GetFormShape(testData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Variable Reference is not matching any children variable:")
}
