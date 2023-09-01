package blueprint_config_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
)

func loadBlueprintFromJsonFile(path string) blueprint_config.BluePrintConfig {
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
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

func TestFindPathsToNodesWithZeroOutdegree(t *testing.T) {
	createMockResult := false
	blueprint := loadBlueprintFromJsonFile("./fullExample.json")
	variables := blueprint_config.Transverse(&blueprint)
	jsonGraph, _ := json.MarshalIndent(&blueprint, "", "  ")
	writeFile("./graphAfterPrrocess.json", jsonGraph)
	if createMockResult {
		jsonString, _ := json.MarshalIndent(variables, "", "  ")
		writeFile("./result.json", jsonString)
	} else {
		expectedQuestions := loadResultFile("./result.json")
		for i, question := range variables {
			expected := expectedQuestions[i]
			if question.ID != expected.ID {
				t.Errorf("Expected %s but got %s", expected.ID, question.ID)
			}
		}
	}
}

func TestGetDisplayOrders(t *testing.T) {
	blueprint := loadBlueprintFromJsonFile("./fullExample.json")
	blueprint_config.Transverse(&blueprint)
	result := blueprint_config.GetAllDisplayOrdersByBFS(&blueprint)
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonResult))
}