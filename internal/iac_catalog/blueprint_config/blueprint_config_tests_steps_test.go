package blueprint_config_test

/*
 THIS FILE CONTAINS CUSTOM TF ACC CHECK STEPS
*/

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func testAccCheckBlueprintConfigExist(resourceName string, blueprintConfig *blueprint_config.BluePrintConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		rawConf := rs.Primary.Attributes["blueprint_config"]

		err := json.Unmarshal([]byte(rawConf), blueprintConfig)
		if err != nil {
			return fmt.Errorf("not a valid blueprint config: %s", rawConf)
		}
		return nil
	}
}

func testAccCheckCorrectVariablesLength(resourceName string, formVariables *[]autocloudsdk.FormShape) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		variables := rs.Primary.Attributes["config"]
		err := json.Unmarshal([]byte(variables), formVariables)

		if err != nil {
			return fmt.Errorf("config variables: %s", variables)
		}
		/*
			if len(*formVariables) != 14 {
				return fmt.Errorf("form variables len: %v", len(*formVariables))
			}*/
		return nil
	}
}

func testAccCheckOmitCorrectness(omitted []string, formVars *[]autocloudsdk.FormShape) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, v := range *formVars {
			id, err := utils.GetVariableID(v.ID)
			if err != nil {
				return fmt.Errorf("error getting varID: %s", err)
			}
			if utils.Contains(omitted, id) && v.IsHidden != true {
				return fmt.Errorf("var: %s not omitted correctly", id)
			}
			fmt.Printf("included variable: %s\n", v.ID)
		}
		return nil
	}
}

// variableConf string, formVars *[]autocloudsdk.FormShape
func testAccCheckOverrides(resourceName string, overideVars []string) resource.TestCheckFunc {
	type ValidationRule struct {
		Rule         string `hcl:"rule"`
		Value        string `hcl:"value,optional"`
		ErrorMessage string `hcl:"error_message,optional"`
	}

	type FieldOption struct {
		Option []struct {
			Label   string `hcl:"label"`
			Value   string `hcl:"value"`
			Checked bool   `hcl:"checked,optional"`
		} `hcl:"option,block"`
	}
	type FormConfig struct {
		Type            string           `hcl:"type"`
		FieldOptions    []FieldOption    `hcl:"options,block"`
		ValidationRules []ValidationRule `hcl:"validation_rule,block"`
	}

	type Variable struct {
		Name        string     `hcl:"name"`
		DisplayName string     `hcl:"display_name"`
		HelperText  string     `hcl:"helper_text"`
		FormConfig  FormConfig `hcl:"form_config,block"`
	}

	type Config struct {
		Variable Variable `hcl:"variable,block"`
	}

	variable := overideVars[0]

	var got Config
	file, diags := hclsyntax.ParseConfig([]byte(variable), "hello.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags != nil {
		fmt.Println("cant parse config")
		fmt.Print(diags)
	}

	diags = gohcl.DecodeBody(file.Body, nil, &got)
	if diags != nil {
		fmt.Println("cant parse body")
		fmt.Print(diags)
	}

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		refVarType := reflect.TypeOf(got.Variable)
		refVarValue := reflect.ValueOf(got.Variable)
		if refVarType.Kind() == reflect.Struct {
			for i := 0; i < refVarType.NumField(); i++ {
				if refVarValue.Field(i).Type() == reflect.TypeOf("") { //if type is string
					fieldType := refVarType.Field(i)
					keyPrefix := "variable.0"
					key := createTerraformStateKey(fieldType.Tag, keyPrefix)
					storedState := rs.Primary.Attributes[key]
					input := refVarValue.FieldByName(fieldType.Name).String()
					stateAndInputMatch := storedState == input
					if !stateAndInputMatch {
						return fmt.Errorf("State file didnt save correct values, statefile: %s, input: %s", storedState, input)
					}
				}
				if refVarValue.Field(i).Type() == reflect.TypeOf(got.Variable.FormConfig) {
					refFormType := reflect.TypeOf(got.Variable.FormConfig)
					refFormValue := reflect.ValueOf(got.Variable.FormConfig)

					for i := 0; i < refFormType.NumField(); i++ {
						fieldType := refFormType.Field(i)
						if refFormValue.Field(i).Type() == reflect.TypeOf("") { //string
							keyPrefix := "variable.0.form_config"
							key := createTerraformStateKey(fieldType.Tag, keyPrefix)
							storedState := rs.Primary.Attributes[key]
							input := refVarValue.FieldByName(fieldType.Name).String()
							stateAndInputMatch := storedState == input
							fmt.Println(key)
							fmt.Println(input)
							fmt.Println(stateAndInputMatch)
							/*if !stateAndInputMatch {
								return fmt.Errorf("State file didnt save correct values, statefile: %s, input: %s", storedState, input)
							}*/
						}
					}
				}
			}
		}

		return nil
	}
}

func createTerraformStateKey(hclTag reflect.StructTag, keyPrefix string) string {
	s := strings.ReplaceAll(string(hclTag), "hcl:", "")
	s = s[1 : len(s)-1] // remove quotation marks
	key := fmt.Sprintf("%s.%s", keyPrefix, s)
	return key
}
