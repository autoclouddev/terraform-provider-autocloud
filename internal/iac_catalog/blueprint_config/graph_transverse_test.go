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
	blueprint := loadBlueprintFromJsonFile("./fullExample.json")
	variables := blueprint_config.Transverse(&blueprint)
	jsonGraph, _ := json.MarshalIndent(blueprint, "", "  ")
	writeFile("./graphAfterPrrocess.json", jsonGraph)
	jsonString, _ := json.MarshalIndent(variables, "", "  ")
	fmt.Println(string(jsonString))
	//writeFile("./result.json", jsonString)
	expectedOutput := loadResultFile("./result.json")
	fmt.Println("result")
	for _, variable := range variables {
		fmt.Printf("%s,", variable.ID)
	}
	fmt.Println("expected")
	for _, variable := range expectedOutput {
		fmt.Printf("%s,", variable.ID)
	}

}

func TestGetDisplayOrder(t *testing.T) {
	blueprint := loadBlueprintFromJsonFile("./fullExample.json")
	blueprint_config.BFS(&blueprint)
}
