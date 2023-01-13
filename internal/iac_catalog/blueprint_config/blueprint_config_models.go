package blueprint_config

import (
	"errors"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

type BluePrintConfig struct {
	Id                string                      `json:"id"`
	RefName           string                      `json:"refName"`
	OmitVariables     []string                    `json:"omitVariables"`
	OverrideVariables map[string]OverrideVariable `json:"overrideVariables"`
	Variables         []autocloudsdk.FormShape    `json:"variables"`
	Children          []BluePrintConfig           `json:"children"`
}

type OverrideVariable struct {
	VariableName   string              `json:"variableName"`
	Value          string              `json:"value"`
	RequiredValues []string            `json:"requiredValues"`
	DisplayName    string              `json:"displayName"`
	HelperText     string              `json:"helperText"`
	FormConfig     FormConfig          `json:"formConfig"`
	Conditionals   []ConditionalConfig `json:"conditionals"`
}

type ConditionalConfig struct {
	Source         string        `json:"source"`
	Condition      string        `json:"condition"`
	Type           string        `json:"type"`
	Options        []FieldOption `json:"options"`
	Value          *string       `json:"value"`
	RequiredValues []string      `json:"requiredValues"`
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
const SHORTTEXT_TYPE = "shortText"
const MAP_TYPE = "map"

var ErrSetValueInForm = errors.New("A form options can not be added when setting the variable's value.")
var ErrOneBlockOptionsRequied = errors.New("No more than 1 \"options\" blocks are allowed")
var ErrShortTextCantHaveOptions = errors.New("ShortText variables can not have options")
var ErrIsRequiredCantHaveValue = errors.New("'isRequired' validation rule can not have a value")
