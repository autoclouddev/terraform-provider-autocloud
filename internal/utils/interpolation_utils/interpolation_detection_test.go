package interpolation_utils_test

import (
	"fmt"
	"testing"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils/interpolation_utils"
)

func TestInterpolationDetectionPassedVars(t *testing.T) {
	template := `{{ Name }} is {{ Age }} years old `

	vars := map[string]string{
		"Name": "John",
	}

	err := interpolation_utils.DetectInterpolation(template, vars)
	if err != nil {
		if err.Error() != "Variable \"Age\" is not found in variables" {
			t.Errorf("Expected error message to be \"Variable \"Age\" is not found in variables\"")
		}
	} else {
		t.Errorf("Expected an error")
	}
}

func TestInterpolationDetectionEmptyVariables(t *testing.T) {
	template := `{{ Name }} is {{ Age }} years old`

	vars := map[string]string{}

	err := interpolation_utils.DetectInterpolation(template, vars)
	if err != nil {
		fmt.Println(err)
	} else {
		t.Errorf("Expected an error")
	}
}

func TestInterpolationDetectionEmptyTemplate(t *testing.T) {
	template := ``

	vars := map[string]string{}

	err := interpolation_utils.DetectInterpolation(template, vars)
	if err != nil {
		t.Errorf("Expected no error")
	}
}

func TestInterpolationDetectionInvalidTemplate(t *testing.T) {
	template := `{{}}}}`

	vars := map[string]string{}

	err := interpolation_utils.DetectInterpolation(template, vars)
	if err == nil {
		t.Errorf("Expected error")
	}
	if err.Error() != "Error parsing template: {{}}}}" {
		t.Errorf("Expected error message to be \"Error parsing template: {{}}}}\"")
	}
}
