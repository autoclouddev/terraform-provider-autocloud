package blueprint_config_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-faker/faker/v4"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func TestOverrideVariable(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "DEBUG", // to see the DEBUG logs
	})
	originalVars := []autocloudsdk.FormShape{
		{
			ID:                  "module1.Id",
			Module:              "module1",
			FormQuestion:        autocloudsdk.FormQuestion{},
			FieldDataType:       "string",
			FieldDefaultValue:   "string",
			FieldValue:          "string",
			AllowConsumerToEdit: true,
			Conditionals:        []autocloudsdk.ConditionalConfig{},
		},
		{
			ID:                  "module1.OtherId",
			Module:              "module1",
			FormQuestion:        autocloudsdk.FormQuestion{},
			FieldDataType:       "string",
			FieldDefaultValue:   "string",
			FieldValue:          "string",
			AllowConsumerToEdit: true,
			Conditionals:        []autocloudsdk.ConditionalConfig{},
		},
		{
			ID:                  "module2.Id2",
			Module:              "module2",
			FormQuestion:        autocloudsdk.FormQuestion{},
			FieldDataType:       "string",
			FieldDefaultValue:   "string",
			FieldValue:          "string",
			AllowConsumerToEdit: true,
			Conditionals:        []autocloudsdk.ConditionalConfig{},
		},
	}

	overrides := make(map[string]blueprint_config.OverrideVariable)
	for _, original := range originalVars {
		ov, err := CreateFakeOverride()
		if err != nil {
			t.Fatalf(err.Error())
		}

		varName, err := utils.GetVariableID(original.ID)
		if err != nil {
			t.Fatalf(err.Error())
		}
		ov.VariableName = varName
		overrides[varName] = *ov
	}
	//override map is modified during blueprint_config.OverrideVariables
	overridesCopy := make(map[string]blueprint_config.OverrideVariable)
	for k, v := range overrides {
		overridesCopy[k] = v
	}

	newVars, err := blueprint_config.OverrideVariables(originalVars, overrides)

	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, newVar := range newVars {
		varName, err := utils.GetVariableID(newVar.ID)
		if err != nil {
			t.Fatalf(err.Error())
		}

		checkFormConfig(newVar, overridesCopy[varName], t)
		checkConditionals(newVar, overridesCopy[varName], t)
		checkAllowConsumerToEdit(newVar, overridesCopy[varName], t)
	}
	t.Cleanup(closer)
}

func TestOverrideVariableError(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "DEBUG", // to see the DEBUG logs
	})
	var expectedError = blueprint_config.ErrVariableNotFound
	originalVars := []autocloudsdk.FormShape{
		{
			ID:                  "module1Id",
			Module:              "module1",
			FormQuestion:        autocloudsdk.FormQuestion{},
			FieldDataType:       "string",
			FieldDefaultValue:   "string",
			FieldValue:          "string",
			AllowConsumerToEdit: true,
			Conditionals:        []autocloudsdk.ConditionalConfig{},
		},
	}

	overrides := make(map[string]blueprint_config.OverrideVariable)
	for i := 0; i < len(originalVars); i++ {
		ov, err := CreateFakeOverride()
		if err != nil {
			t.Fatalf(err.Error())
		}
		overrides[ov.VariableName] = *ov
	}

	_, err := blueprint_config.OverrideVariables(originalVars, overrides)

	if !errors.Is(err, expectedError) {
		t.Fatalf("no error was detected")
	}

	t.Cleanup(closer)
}

func TestBuildOverride(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "DEBUG", // to see the DEBUG logs
	})
	ov, err := CreateFakeOverride()
	if err != nil {
		t.Fatalf(err.Error())
	}

	original := autocloudsdk.FormShape{
		ID:                  "s3_module.someId",
		Module:              "s3_module",
		FormQuestion:        autocloudsdk.FormQuestion{},
		FieldDataType:       "string",
		FieldDefaultValue:   "string",
		FieldValue:          "string",
		AllowConsumerToEdit: true,
		Conditionals:        []autocloudsdk.ConditionalConfig{},
	}

	newVar := blueprint_config.BuildOverridenVariable(original, *ov)

	if newVar.ID != original.ID {
		t.Fatalf("id was not properly generated")
	}

	checkFormConfig(newVar, *ov, t)
	checkConditionals(newVar, *ov, t)
	checkAllowConsumerToEdit(newVar, *ov, t)
	t.Cleanup(closer)
}

func TestBuildGenericVar(t *testing.T) {
	closer := acctest.EnvSetter(map[string]string{
		"TF_LOG": "DEBUG", // to see the DEBUG logs
	})
	ov, err := CreateFakeOverride()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fs := blueprint_config.BuildGenericVariable(*ov)

	// check form shape is created
	genericId := "generic." + ov.VariableName
	if fs.ID != genericId {
		t.Fatalf("names do not match, overrideId: %s, formSape.ID: %s", ov.VariableName, fs.ID)
	}
	if fs.Module != "generic" {
		t.Fatalf("this is not a generic module")
	}

	checkFormConfig(fs, *ov, t)
	checkConditionals(fs, *ov, t)
	checkAllowConsumerToEdit(fs, *ov, t)
	t.Cleanup(closer)
}

func checkFormConfig(formShape autocloudsdk.FormShape, override blueprint_config.OverrideVariable, t *testing.T) {
	// FormConfig.FieldOptions
	for i, option := range formShape.FormQuestion.FieldOptions {
		if option.FieldID != fmt.Sprintf("%s-%s", formShape.ID, override.FormConfig.FieldOptions[i].Value) {
			t.Fatalf("formOptions is not correct")
		}
	}
	// FormConfig.ValidationRules
	for i, rule := range formShape.FormQuestion.ValidationRules {
		strRule, _ := utils.PrettyStruct(rule)
		strFormRule, _ := utils.PrettyStruct(override.FormConfig.ValidationRules[i])
		if strRule != strFormRule {
			t.Fatalf("validation rules not being formed correctly")
		}
	}
}

func checkConditionals(formShape autocloudsdk.FormShape, override blueprint_config.OverrideVariable, t *testing.T) {
	// CONDITIONALS

	if len(formShape.Conditionals) != len(override.Conditionals) {
		t.Fatalf(
			"conditions do not match, len(formShape.Conditionas): %v len(override.Conditinals): %v\n",
			len(formShape.Conditionals),
			len(override.Conditionals),
		)
	}
	for i, cond := range formShape.Conditionals {
		for j, option := range cond.Options {
			if formShape.Conditionals[i].Options[j].FieldID != fmt.Sprintf("%s-%s", formShape.ID, option.Value) {
				t.Fatalf("option.FieldId was not formed ok")
			}
		}
	}
}

func checkAllowConsumerToEdit(formShape autocloudsdk.FormShape, override blueprint_config.OverrideVariable, t *testing.T) {
	isValueSet := len(override.Value) != 0
	if isValueSet {
		UserCANNOTEdit := false
		if !(formShape.AllowConsumerToEdit == UserCANNOTEdit) { // I want to be explicit about hte value
			t.Fatalf(
				"\"AllowConsumerToEdit\" should be false when override has a value set, override.Value: %v, formShape.AllowConsumerToEdit: %v",
				override.Value,
				formShape.AllowConsumerToEdit,
			)
		}
	}
}

func CreateFakeOverride() (*blueprint_config.OverrideVariable, error) {
	var formConfig blueprint_config.FormConfig
	err := faker.FakeData(&formConfig)
	if err != nil {
		return nil, err
	}
	// faker cant create arrays of specific length
	formConfig.FieldOptions = make([]blueprint_config.FieldOption, 0)
	for i := 0; i < 2; i++ {
		var fieldOption blueprint_config.FieldOption
		err := faker.FakeData(&fieldOption)
		if err != nil {
			return nil, err
		}
		formConfig.FieldOptions = append(formConfig.FieldOptions, fieldOption)
	}
	var ValidationRule blueprint_config.ValidationRule
	err = faker.FakeData(&ValidationRule)
	if err != nil {
		return nil, err
	}
	formConfig.ValidationRules = []blueprint_config.ValidationRule{ValidationRule}

	var conditional blueprint_config.ConditionalConfig
	err = faker.FakeData(&conditional)
	if err != nil {
		fmt.Print(err)
	}
	/*conditional.Options = make([]blueprint_config.FieldOption, 0)
	for i := 0; i < 2; i++ {
		var fieldOption blueprint_config.FieldOption
		err := faker.FakeData(&fieldOption)
		if err != nil {
			return nil, err
		}
		conditional.Options = append(conditional.Options, fieldOption)
	}*/
	var ov blueprint_config.OverrideVariable
	err = faker.FakeData(&ov)
	if err != nil {
		return nil, err
	}
	ov.FormConfig = formConfig
	ov.Conditionals = []blueprint_config.ConditionalConfig{conditional}
	return &ov, nil
}
