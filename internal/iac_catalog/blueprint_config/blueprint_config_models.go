package blueprint_config

import (
	"errors"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

type BluePrintConfig struct {
	Id                string                      `json:"id"`
	RefName           string                      `json:"refName"`
	OmitVariables     []string                    `json:"ommitVariables"`
	OverrideVariables map[string]OverrideVariable `json:"overrideVariables"`
	Variables         []autocloudsdk.FormShape    `json:"variables"`
	Children          map[string]BluePrintConfig  `json:"children"`
}

type OverrideVariable struct {
	VariableName   string              `json:"variableName" faker:"word"`
	Value          string              `json:"value" faker:"word"`
	RequiredValues []string            `json:"requiredValues" faker:"slice_len=2"`
	DisplayName    string              `json:"displayName" faker:"word"`
	HelperText     string              `json:"helperText" faker:"word"`
	FormConfig     FormConfig          `json:"formConfig"`
	Conditionals   []ConditionalConfig `json:"conditionals"`
}

type ConditionalConfig struct {
	Source         string        `json:"source" faker:"word"`
	Condition      string        `json:"condition" faker:"word"`
	Type           string        `json:"type" faker:"word"`
	Options        []FieldOption `json:"options"`
	Value          *string       `json:"value" faker:"word"`
	RequiredValues []string      `json:"requiredValues" faker:"slice_len=2"`
}
type FormConfig struct {
	Type            string           `json:"type" faker:"oneof: checkbox, radio"`
	FieldOptions    []FieldOption    `json:"fieldOptions"`
	ValidationRules []ValidationRule `json:"validationRules"`
}

type ValidationRule struct {
	Rule         string `json:"rule" faker:"word"`
	Value        string `json:"value" faker:"word"`
	ErrorMessage string `json:"errorMessage" faker:"word"`
}

type FieldOption struct {
	Label   string `json:"label" faker:"word"`
	Value   string `json:"value" faker:"word"`
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
var ErrVariableNotFound = errors.New("ERROR: no variable ID found")
