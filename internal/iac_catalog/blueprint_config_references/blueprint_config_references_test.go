package blueprint_config_references_test

import (
	"testing"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
)

func TestBlueprintConfigReferencesSingleton(t *testing.T) {
	// Initialize the singleton instance
	data := blueprint_config_references.GetInstance()

	// Set some values and verify that they are set
	data.SetValue("key1", "value1")
	if data.GetValue("key1") != "value1" {
		t.Errorf("Error: expected value1, but got %s", data.GetValue("key1"))
	}

	data.SetValue("key2", "value2")
	if data.GetValue("key2") != "value2" {
		t.Errorf("Error: expected value2, but got %s", data.GetValue("key2"))
	}

	// Convert the data to string and verify that it is correct
	expectedJson := `{"key1":"value1","key2":"value2"}`
	jsonData := data.ToString()

	if jsonData != expectedJson {
		t.Errorf("Error: expected %s, but got %s", expectedJson, jsonData)
	}
}
