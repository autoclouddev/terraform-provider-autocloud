package blueprint_config_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
)

func loadBlueprintFromJsonFile(path string) blueprint_config.BluePrintConfig {
	jsonFile, err := os.Open(path)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var blueprint blueprint_config.BluePrintConfig
	err = json.Unmarshal(byteValue, &blueprint)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return blueprint
}

func loadResultFile(path string) []generator.FormShape {
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var result []generator.FormShape
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return result
}

func writeFile(path string, content []byte) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	_, err2 := f.Write(content)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func TestGraphTransverse(t *testing.T) {
	testCases := []struct {
		inputFilePath    string
		expectedFilePath string
	}{
		{
			inputFilePath:    "./blueprint_fixtures/gcp-getting-started.json",
			expectedFilePath: "./blueprint_fixtures/expected_result_gcp-getting-started.json",
		},
		{
			inputFilePath:    "./blueprint_fixtures/fargate_ecs_cluster.json",
			expectedFilePath: "./blueprint_fixtures/expected_result_fargate_ecs_cluster.json",
		},
		// Add more test cases as needed
	}
	createMockResult := false
	for _, tc := range testCases {
		t.Run(tc.inputFilePath, func(t *testing.T) {
			blueprint := loadBlueprintFromJsonFile(tc.inputFilePath)
			variables := blueprint_config.Transverse(&blueprint)
			if createMockResult {
				jsonString, _ := json.MarshalIndent(variables, "", "  ")
				writeFile(tc.expectedFilePath, jsonString)
			} else {
				expectedQuestions := loadResultFile(tc.expectedFilePath)
				for i, question := range variables {
					expected := expectedQuestions[i]
					if question.ID != expected.ID {
						t.Errorf("Expected %s but got %s", expected.ID, question.ID)
					}
				}
			}
		})
	}
}

func TestGetDisplayOrders(t *testing.T) {
	blueprint := loadBlueprintFromJsonFile("./blueprint_fixtures/gcp-getting-started.json")
	blueprint_config.Transverse(&blueprint)
	result := blueprint_config.GetAllDisplayOrdersByBFS(&blueprint)

	expectedContent := []string{
		"generic.namespace",
		"generic.environment",
		"generic.name",
		"kmskey.keyring",
		"generic.bucket_name",
		"generic.location",
		"generic.project_id",
		"generic.labels",
	}

	if !reflect.DeepEqual(expectedContent, result[0].Values) {
		t.Errorf("Result slice does not match the expected content")
	}

	//uncomment to see the output
	//jsonResult, _ := json.MarshalIndent(result, "", "  ")
	//fmt.Println(string(jsonResult))
}
