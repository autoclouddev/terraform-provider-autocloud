package utils_test

import (
	"errors"
	"fmt"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func TestParseAndMapVariables(t *testing.T) {
	json := `
	[
  {
    "id": "s3bucket.attach_require_latest_tls_policy",
    "type": "string",
    "module": "S3bucket",
    "formQuestion": {
      "fieldId": "s3bucket.attach_require_latest_tls_policy",
      "fieldType": "radio",
      "fieldLabel": "attach_require_latest_tls_policy",
      "fieldOptions": [
        {
          "label": "Yes",
          "value": "true",
          "checked": false,
          "fieldId": "s3bucket.attach_require_latest_tls_policy-true"
        },
        {
          "label": "No",
          "value": "false",
          "checked": true,
          "fieldId": "s3bucket.attach_require_latest_tls_policy-false"
        }
      ],
      "explainingText": "Controls if S3 bucket should require the latest version of TLS",
      "validationRules": null
    }
  },
  {
    "id": "s3bucket.attach_elb_log_delivery_policy",
    "type": "string",
    "module": "S3bucket",
    "formQuestion": {
      "fieldId": "s3bucket.attach_elb_log_delivery_policy",
      "fieldType": "radio",
      "fieldLabel": "attach_elb_log_delivery_policy",
      "fieldOptions": [
        {
          "label": "Yes",
          "value": "true",
          "checked": false,
          "fieldId": "s3bucket.attach_elb_log_delivery_policy-true"
        },
        {
          "label": "No",
          "value": "false",
          "checked": true,
          "fieldId": "s3bucket.attach_elb_log_delivery_policy-false"
        }
      ],
      "explainingText": "Controls if S3 bucket should have ELB log delivery policy attached",
      "validationRules": null
    }
  }
]
`
	vars, err := utils.ParseVariables(json)
	assert.Nil(t, err)
	assert.NotNil(t, vars)
	assert.Equal(t, vars[1].ID, "s3bucket.attach_elb_log_delivery_policy")

	varsMap, err := utils.GetVariablesIdMap(json)
	assert.Nil(t, err)
	assert.Equal(t, "s3bucket.attach_require_latest_tls_policy", varsMap["attach_require_latest_tls_policy"])
	assert.Equal(t, "s3bucket.attach_elb_log_delivery_policy", varsMap["attach_elb_log_delivery_policy"])
}

func TestHasReference(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{input: "foo.variables.bar", expected: true},
		{input: "foo.variables_", expected: false},
		{input: "foo.bar.variables", expected: false},
		{input: "foo.bar", expected: false},
		{input: "", expected: false},
	}

	for _, tc := range testCases {
		actual := utils.HasReference(tc.input)
		if actual != tc.expected {
			t.Errorf("HasReference(%s) = %v, expected %v", tc.input, actual, tc.expected)
		}
	}
}

func TestIsValidId(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid input",
			input:    "module.id",
			expected: true,
		},
		{
			name:     "Invalid input - empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid input - missing variable name",
			input:    "module.",
			expected: false,
		},
		{
			name:     "Invalid input - missing module name",
			input:    ".variable",
			expected: false,
		},
		{
			name:     "Invalid input - too many components",
			input:    "module.id.extra",
			expected: false,
		},
		{
			name:     "Invalid input - invalid characters",
			input:    "module!id",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.IsValidId(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestGetVariableID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "Valid input",
			input:    "module.variable",
			expected: "variable",
			err:      nil,
		},
		{
			name:     "Invalid input - empty string",
			input:    "",
			expected: "",
			err:      errors.New("Invalid Key"),
		},
		{
			name:     "Invalid input - missing variable name",
			input:    "module.",
			expected: "",
			err:      errors.New("Invalid Key"),
		},
		{
			name:     "Invalid input - missing module name",
			input:    ".variable",
			expected: "",
			err:      errors.New("Invalid Key"),
		},
		{
			name:     "Invalid input - too many components",
			input:    "module.variable.extra",
			expected: "",
			err:      errors.New("Invalid Key"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.GetVariableID(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, but got %q", tc.expected, result)
			}
			if err == nil && tc.err != nil || err != nil && tc.err == nil || err != nil && tc.err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected error %v, but got %v", tc.err, err)
			}
		})
	}
}

func TestGetVariableReferenceID(t *testing.T) {
	aliases := blueprint_config_references.GetInstance()
	aliases.SetValue("alias1", "module1")
	aliases.SetValue("alias2", "module2")

	tests := []struct {
		input string
		want  string
		err   error
	}{
		{"alias1.variables.variable1", "module1.variable1", nil},
		{"alias2.variables.variable2", "module2.variable2", nil},
		{"alias1.variables.", "", errors.New("Invalid Key")},
		{"alias1.variables", "", errors.New("Invalid Key")},
		{"alias1", "", errors.New("Invalid Key")},
		{"", "", errors.New("Invalid Key")},
	}

	for _, tt := range tests {
		got, err := blueprint_config.GetVariableReferenceID(tt.input, &blueprint_config.BluePrintConfig{})

		assert.Equal(t, tt.err, err, fmt.Sprintf("Error mismatch for input %q", tt.input))
		assert.Equal(t, tt.want, got, fmt.Sprintf("Output mismatch for input %q", tt.input))
	}
}

func TestFindIdx(t *testing.T) {
	// Create test data
	vars := []generator.FormShape{
		{
			ID: "test.test1",
		},
		{
			ID: "test.test2",
		},
		{
			ID: "test.test3",
		},
	}

	// Test case 1: reference variable not found
	reference := "alias.variables.test4"
	expected := []int{}
	result := blueprint_config.FindIdx(vars, reference, &blueprint_config.BluePrintConfig{})
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Test case 1 failed: expected %v, but got %v", expected, result)
	}

	// Test case 2: regular variable not found
	reference = "test4"
	expected = []int{}
	result = blueprint_config.FindIdx(vars, reference, &blueprint_config.BluePrintConfig{})
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Test case 2 failed: expected %v, but got %v", expected, result)
	}

	// Test case 3: reference variable found
	aliasToModuleNameMap := blueprint_config_references.GetInstance()
	aliasToModuleNameMap.SetValue("alias", "test")
	reference = "alias.variables.test1"
	expected = []int{0}
	result = blueprint_config.FindIdx(vars, reference, &blueprint_config.BluePrintConfig{})
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Test case 3 failed: expected %v, but got %v", expected, result)
	}

	// Test case 4: regular variable found
	reference = "test1"
	expected = []int{0}
	result = blueprint_config.FindIdx(vars, reference, &blueprint_config.BluePrintConfig{})
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Test case 4 failed: expected %v, but got %v", expected, result)
	}

	// Test case 5: multiple matches found
	reference = "test1"
	expected = []int{0, 3}
	vars = append(vars, generator.FormShape{
		ID: "test2.test1",
	}, generator.FormShape{
		ID: "test.test5",
	})
	result = blueprint_config.FindIdx(vars, reference, &blueprint_config.BluePrintConfig{})
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Test case 5 failed: expected %v, but got %v", expected, result)
	}
}

func TestMergeMaps(t *testing.T) {
	// Test case 1: Empty map 1 and empty map 2
	map1 := make(map[string]string)
	map2 := make(map[string]string)
	const VAL1 = "value1"
	const VAL3 = "value3"
	utils.MergeMaps(&map1, &map2)

	if len(map1) != 0 {
		t.Errorf("Expected empty map but got %v", map1)
	}

	// Test case 2: Non-empty map 1 and empty map 2
	map1 = map[string]string{"key1": VAL1, "key2": "value2"}
	map2 = make(map[string]string)
	utils.MergeMaps(&map1, &map2)

	if len(map1) != 2 || map1["key1"] != VAL1 || map1["key2"] != "value2" {
		t.Errorf("Expected map {key1: value1, key2: value2} but got %v", map1)
	}

	// Test case 3: Empty map 1 and non-empty map 2
	map1 = make(map[string]string)
	map2 = map[string]string{"key3": VAL3, "key4": "value4"}
	utils.MergeMaps(&map1, &map2)

	if len(map1) != 2 || map1["key3"] != VAL3 || map1["key4"] != "value4" {
		t.Errorf("Expected map {key3: value3, key4: value4} but got %v", map1)
	}

	// Test case 4: Non-empty map 1 and non-empty map 2 with no overlapping keys
	map1 = map[string]string{"key1": VAL1, "key2": "value2"}
	map2 = map[string]string{"key3": VAL3, "key4": "value4"}
	utils.MergeMaps(&map1, &map2)

	if len(map1) != 4 || map1["key1"] != VAL1 || map1["key2"] != "value2" || map1["key3"] != VAL3 || map1["key4"] != "value4" {
		t.Errorf("Expected map {key1: value1, key2: value2, key3: value3, key4: value4} but got %v", map1)
	}

	// Test case 5: Non-empty map 1 and non-empty map 2 with overlapping keys
	map1 = map[string]string{"key1": VAL1, "key2": "value2"}
	map2 = map[string]string{"key2": "new_value2", "key3": VAL3}
	utils.MergeMaps(&map1, &map2)

	if len(map1) != 3 || map1["key1"] != VAL1 || map1["key2"] != "new_value2" || map1["key3"] != "value3" {
		t.Errorf("Expected map {key1: value1, key2: new_value2, key3: value3} but got %v", map1)
	}
}

func TestLoadData(t *testing.T) {
	const FIXTURES_FOLDER = "fixtures"

	type ExpectedType struct {
		FieldA string `json:"fieldA"`
		FieldB int    `json:"fieldB"`
	}
	// Test case where file exists and is valid
	expected := ExpectedType{
		FieldA: "test",
		FieldB: 123,
	}
	filePath := path.Join(FIXTURES_FOLDER, "valid_test.json")
	actual, err := utils.LoadData[ExpectedType](filePath)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Actual value %v did not match expected value %v", actual, expected)
	}

	// Test case where file does not exist
	filePath = path.Join(FIXTURES_FOLDER, "nonexistent_test.json")
	_, err = utils.LoadData[ExpectedType](filePath)
	if err == nil {
		t.Errorf("Expected error when reading non-existent file, but received none")
	}
}
