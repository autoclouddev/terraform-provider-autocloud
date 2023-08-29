package blueprint_config_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
)

const FIXTURES_FOLDER = "fixtures"

func loadTestData[T any](testCase string) (out T) {
	var testData T
	file, err := os.ReadFile(filepath.Join(FIXTURES_FOLDER, testCase))
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	err = json.Unmarshal(file, &testData)
	if err != nil {
		fmt.Println("Error loading file:", err)
	}
	return testData
}

func TestGetBlueprintConfigSources(t *testing.T) {
	bp := blueprint_config.BluePrintConfig{}
	bp.Id = "sources"
	testData := loadTestData[interface{}]("sources.json")
	err := blueprint_config.GetBlueprintConfigSources(testData, &bp)
	assert.Nil(t, err)
}

func TestGetBlueprintConfigOmitVariable(t *testing.T) {
	bp := blueprint_config.BluePrintConfig{}
	bp.Id = "omit_variables"
	testData := loadTestData[interface{}]("omit_variables.json")
	err := blueprint_config.GetBlueprintConfigOmitVariables(testData, &bp)
	assert.Nil(t, err)
}

func TestGetBlueprintConfigDisplayOrder(t *testing.T) {
	bp := blueprint_config.BluePrintConfig{}
	bp.Id = "display_order"
	aliases := blueprint_config_references.GetInstance()
	testData := loadTestData[interface{}]("display_order.json")

	aliases.SetValue("ec2", "ec2_instance")
	err := blueprint_config.GetBlueprintConfigDisplayOrder(testData, &bp, *aliases)
	assert.Nil(t, err)
}
