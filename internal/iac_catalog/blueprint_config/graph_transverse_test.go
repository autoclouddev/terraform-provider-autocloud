package blueprint_config_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

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

func TestFindPathsToNodesWithZeroOutdegree(t *testing.T) {
	blueprint := loadBlueprintFromJsonFile("./finalConfig.json")
	variables := blueprint_config.Transverse(&blueprint)
	jsonString, _ := json.MarshalIndent(variables, "", "  ")
	fmt.Println(string(jsonString))
}
