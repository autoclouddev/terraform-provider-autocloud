package blueprint_config

import (
	"errors"

	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
)

type BluePrintConfig struct {
	Id                string                      `json:"id"`
	RefName           string                      `json:"refName"`
	OmitVariables     []string                    `json:"omitVariables"`
	OverrideVariables map[string]OverrideVariable `json:"overrideVariables"`
	Variables         []generator.FormShape       `json:"variables"`
	DisplayOrder      DisplayOrder                `json:"displayOrder"`
	Children          map[string]*BluePrintConfig `json:"children"`
}

type VariableContent struct {
	Value          string            `json:"value" faker:"word"`
	RequiredValues string            `json:"requiredValues" faker:"slice_len=2"`
	DisplayName    string            `json:"displayName" faker:"word"`
	HelperText     string            `json:"helperText" faker:"word"`
	Default        string            `json:"default" faker:"word"`
	Variables      map[string]string `json:"variables"`
	FormConfig     FormConfig        `json:"formConfig"`
}
type OverrideVariable struct {
	VariableName string `json:"variableName" faker:"word"`
	VariableContent
	Conditionals      []ConditionalConfig `json:"conditionals"`
	IsHidden          bool                `json:"isHidden"` // based on omit variables
	UsedInHCL         bool                `json:"usedInHCL"`
	InterpolationVars map[string]string   `json:"interpolationVars"`
	dirty             bool
}

type ConditionalConfig struct {
	Source    string `json:"source" faker:"word"`
	Condition string `json:"condition" faker:"word"`
	VariableContent
}

type FormConfig struct {
	Type            string           `json:"type" faker:"oneof: checkbox, radio"`
	FieldOptions    []FieldOption    `json:"fieldOptions"`
	ValidationRules []ValidationRule `json:"validationRules"`
}

type ValidationRule struct {
	Rule         string `json:"rule" faker:"word"`
	Value        string `json:"value" faker:"word"`
	Scope        string `json:"scope" faker:"word"`
	ErrorMessage string `json:"errorMessage" faker:"word"`
}

type FieldOption struct {
	Label   string `json:"label" faker:"word"`
	Value   string `json:"value" faker:"word"`
	Checked bool   `json:"checked"`
}

type DisplayOrder struct {
	Priority int      `json:"priority" faker:"0"`
	Values   []string `json:"values"`
}

const GENERIC = "generic"
const RADIO_TYPE = "radio"
const CHECKBOX_TYPE = "checkbox"
const LIST_TYPE = "list"
const SHORTTEXT_TYPE = "shortText"
const MAP_TYPE = "map"
const RAW_TYPE = "raw"
const EDITOR_TYPE = "editor"

const STATE_FILE = "terraform.tfstate"

var ErrSetValueInForm = errors.New("A form options can not be added when setting the variable's value.")
var ErrOneBlockOptionsRequied = errors.New("No more than 1 \"options\" blocks are allowed")
var ErrShortTextCantHaveOptions = errors.New("ShortText variables can not have options")
var ErrIsRequiredCantHaveValue = errors.New("'isRequired' validation rule can not have a value")
var ErrRegexOnlyCanHaveScope = errors.New("Only 'regex' validation rule can have a scope")
var ErrVariableNotFound = errors.New("ERROR: no variable ID found")
var ErrMapCantBeParsed = errors.New("Map type cant be created")
