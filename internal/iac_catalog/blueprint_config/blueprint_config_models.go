package blueprint_config

import autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"

type BluePrintConfig struct {
	Id                string                      `json:"id"`
	RefName           string                      `json:"refName"`
	OmitVariables     []string                    `json:"-"`
	OverrideVariables map[string]OverrideVariable `json:"-"`
	Variables         []autocloudsdk.FormShape    `json:"variables"`
	Children          []BluePrintConfig           `json:"children"`
	//omit
}

type OverrideVariable struct {
	VariableName string     `json:"variableName"`
	Value        *string    `json:"value"`
	DisplayName  string     `json:"displayName"`
	HelperText   string     `json:"helperText"`
	FormConfig   FormConfig `json:"formConfig"`
}

type FormConfig struct {
	Type            string           `json:"type"`
	FieldOptions    []FieldOption    `json:"fieldOptions"`
	ValidationRules []ValidationRule `json:"validationRules"`
}

type ValidationRule struct {
	Rule         string `json:"rule"`
	Value        string `json:"value"`
	ErrorMessage string `json:"errorMessage"`
}

type FieldOption struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Checked bool   `json:"checked"`
}

const GENERIC = "generic"
const RADIO_TYPE = "radio"
const CHECKBOX_TYPE = "checkbox"
