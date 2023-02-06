package blueprint_config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func DataSourceBlueprintConfig() *schema.Resource {
	setOfStringSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	validationRulesSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rule": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"isRequired", "regex"}, false),
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"error_message": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
			},
		},
	}

	optionItemSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"label": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
				"checked": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}

	variableSchema := map[string]*schema.Schema{

		//"form_config": formConfigSchema,
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"helper_text": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{SHORTTEXT_TYPE, RADIO_TYPE, CHECKBOX_TYPE, MAP_TYPE}, false),
		},
		"options": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"option": optionItemSchema,
				},
			},
		},
		"required_values": {
			Type:     schema.TypeString,
			Optional: true,
		},
		//"conditional":     conditionalSchema,
		"validation_rule": validationRulesSchema,
	}

	conditionalSchema := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"source": {
					Type:     schema.TypeString,
					Required: true,
				},
				"condition": {
					Type:     schema.TypeString,
					Required: true,
				},
				"content": {
					Type:     schema.TypeSet,
					Required: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: utils.MergeSchemas(variableSchema, map[string]*schema.Schema{}),
					},
				},
			},
		},
	}

	bluePrintConfigSchema := map[string]*schema.Schema{
		"source": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"omit_variables": setOfStringSchema,
		"variable": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: utils.MergeSchemas(variableSchema, map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"conditional": conditionalSchema,
				}),
			},
		},
		"config": { // the form as json to replace the default variables
			Description: "Variables retrieved in the tree",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"blueprint_config": { // the form as json to replace the default variables
			Description: "Processed form variables JSON (to replace the default module variables variables)",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}

	return &schema.Resource{
		Description: "terraform form processor (form builder)",
		ReadContext: dataSourceBlueprintConfigRead,
		Schema:      bluePrintConfigSchema,
	}
}

func dataSourceBlueprintConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// map the resource to a FormBuilder object
	blueprintConfig, err := GetBlueprintConfigFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	pretty, err := utils.PrettyStruct(blueprintConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("INPUT BLUEPRINTCONFIG->", pretty)
	// new form variables (as JSON)
	formVariables, err := GetFormShape(*blueprintConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	err = validateConditionals(formVariables)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonFormShape, err := utils.ToJsonString(formVariables)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("config", jsonFormShape)
	if err != nil {
		return diag.FromErr(err)
	}

	/*bpOutput, err := utils.ToJsonString(blueprintConfig)
	//log.Println(bpOutput)
	if err != nil {
		return diag.FromErr(err)
	}*/

	err = d.Set("blueprint_config", pretty)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(blueprintConfig.Id)

	return diags
}

func ConvertMap(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)

	for key, value := range mapInterface {
		mapValue := value
		if value == nil {
			mapValue = ""
		}
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", mapValue)

		mapString[strKey] = strValue
	}

	return mapString
}

// maps tf declaration to object
//
//nolint:gocyclo
func GetBlueprintConfigFromSchema(d *schema.ResourceData) (*BluePrintConfig, error) {
	bp := &BluePrintConfig{}
	bp.Id = strconv.FormatInt(time.Now().Unix(), 10)
	bp.OverrideVariables = make(map[string]OverrideVariable, 0)
	bp.Children = make(map[string]BluePrintConfig)
	if v, ok := d.GetOk("source"); ok {
		for key, value := range v.(map[string]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			log.Printf("SOURCE_INPUT_KEY: %v\n", key)
			formattedValue, _ := utils.PrettyString(strValue)
			log.Printf("SOURCE_INPUT_VALUE: %v", formattedValue)
			bc := BluePrintConfig{}
			err := json.Unmarshal([]byte(strValue), &bc)
			if err != nil {
				return nil, errors.New("invalid conversion to BluePrintConfig")
			}
			bp.Children[strKey] = bc
		}
	}
	if v, ok := d.GetOk("omit_variables"); ok {
		log.Printf("omit_vars get.ok is ok, %v\n", v)
		list := v.(*schema.Set).List()
		omit := make([]string, len(list))
		for i, optionValue := range list {
			omit[i] = optionValue.(string)
		}
		bp.OmitVariables = omit
		log.Printf("the [%v] are the omitted vars", bp.OmitVariables)
	} else {
		log.Printf("omit_vars get.ok not ok, no variables were added\n")
	}

	if v, ok := d.GetOk("variable"); ok {
		varsList := v.([]interface{})
		// override vars loop
		for _, currentVar := range varsList {
			varOverrideMap := currentVar.(map[string]interface{})
			//create variable
			varName := varOverrideMap["name"].(string)
			vc, err := BuildVariableFromSchema(varOverrideMap)
			if err != nil {
				return nil, err
			}

			bp.OverrideVariables[varName] = OverrideVariable{
				VariableName:    varName,
				VariableContent: *vc,
			}

			// Conditionals
			conditionals, conditionalExists := varOverrideMap["conditional"].(*schema.Set)
			log.Printf("CONDITIONALS: %v \n", conditionals)
			if conditionalExists {
				if entry, ok := bp.OverrideVariables[varName]; ok {
					conditionals, err := getConditionals(conditionals)
					if err != nil {
						return nil, errors.New("GetBlueprintConfigFromSchema: Error accessing bp")
					}
					entry.Conditionals = conditionals
					bp.OverrideVariables[varName] = entry
				} else {
					return nil, errors.New("GetBlueprintConfigFromSchema: Error accessing bp")
				}
			}
		}
	}
	str, err := json.MarshalIndent(bp, "", "    ")
	if err != nil {
		return nil, errors.New("invalid conversion to BluePrintConfig")
	}
	log.Printf("final bc: %s", string(str))
	return bp, nil
}

func validateConditionals(variables []autocloudsdk.FormShape) error {
	// vars to map
	var varsMap = make(map[string]autocloudsdk.FormShape, len(variables))
	for _, variable := range variables {
		varsMap[variable.ID] = variable
	}

	// validate conditionals
	for _, variable := range variables {
		for _, conditional := range variable.Conditionals {
			dependencyVariable, dependecyExist := varsMap[conditional.Source]
			if dependecyExist && dependencyVariable.FormQuestion.FieldType != RADIO_TYPE {
				return fmt.Errorf("the conditional's source variable can only be of 'radio' type [variable: %v, source variable: %v, source variable type: %v]", variable.ID, conditional.Source, dependencyVariable.FormQuestion.FieldType)
			}
		}
	}

	return nil
}

func BuildVariableFromSchema(rawSchema map[string]interface{}) (*VariableContent, error) {
	content := &VariableContent{}
	var requiredValues string
	requiredValuesInput, requiredValuesInputExist := rawSchema["required_values"]
	if requiredValuesInputExist {
		requiredValues = requiredValuesInput.(string)
	}

	content.DisplayName = rawSchema["display_name"].(string)
	content.HelperText = rawSchema["helper_text"].(string)
	content.RequiredValues = requiredValues

	// Note: if it has a value, then it can NOT have form options "options"
	valueIsDefined := false
	value, valueExist := rawSchema["value"]
	valueStr, valueIsString := value.(string)
	valueIsDefined = valueStr != "" // NOTE: if the value is empty, we consider it as 'not defined'

	if valueExist && valueIsString && valueIsDefined {
		content.Value = valueStr
		return content, nil
	}
	// variableContent with form options

	optionsFromSchema := rawSchema["options"].(*schema.Set)
	if len(optionsFromSchema.List()) > 1 {
		// it should be caught at schema check level - adding the check here to enforce it in case the schema changes
		return nil, errors.New("exactly one \"options\" must be defined")
	}

	variableType := rawSchema["type"].(string)

	if variableType == "shortText" && len(optionsFromSchema.List()) > 0 {
		return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrShortTextCantHaveOptions)
	}
	content.FormConfig = FormConfig{
		Type:            variableType,
		ValidationRules: make([]ValidationRule, 0),
		FieldOptions:    make([]FieldOption, 0),
	}
	if variableType == RADIO_TYPE || variableType == CHECKBOX_TYPE {
		rawOptionsCluster := optionsFromSchema.List() // "options key in schema options {}" should always have 1 elem
		if len(optionsFromSchema.List()) == 1 {
			rawOptions := rawOptionsCluster[0].(map[string]interface{})
			optionSchema := rawOptions["option"].(*schema.Set)

			for _, option := range optionSchema.List() {
				options := option.(map[string]interface{})
				fieldOption := FieldOption{
					Label:   options["label"].(string),
					Value:   options["value"].(string),
					Checked: options["checked"].(bool),
				}
				content.FormConfig.FieldOptions = append(content.FormConfig.FieldOptions, fieldOption)
			}
		}
	}

	if variableType == MAP_TYPE {
		var variablesMap map[string]interface{}
		err := json.Unmarshal([]byte(requiredValues), &variablesMap)
		if err != nil {
			return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrMapCantBeParsed)
		}

		ccc := ConvertMap(variablesMap)

		pairs := make([]autocloudsdk.KeyValue, 0)

		for key, value := range ccc {
			pair := autocloudsdk.KeyValue{Key: key, Value: value}
			pairs = append(pairs, pair)
		}
		mapValue, err := json.Marshal(pairs)
		if err != nil {
			return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrMapCantBeParsed)
		}
		content.RequiredValues = string(mapValue)
	}
	validationRulesList := rawSchema["validation_rule"].(*schema.Set).List()

	for _, validationRule := range validationRulesList {
		validationRuleMap := validationRule.(map[string]interface{})

		rule := validationRuleMap["rule"].(string)
		ruleValue := validationRuleMap["value"].(string)

		if rule == "isRequired" && ruleValue != "" {
			return nil, fmt.Errorf("GetBlueprintConfigFromSchema: %w", ErrIsRequiredCantHaveValue)
		}
		vr := ValidationRule{
			Rule:         rule,
			Value:        ruleValue,
			ErrorMessage: validationRuleMap["error_message"].(string),
		}
		content.FormConfig.ValidationRules = append(content.FormConfig.ValidationRules, vr)
	}
	return content, nil
}

func getConditionals(varOverrideMap *schema.Set) ([]ConditionalConfig, error) {
	conditionalLen := varOverrideMap.Len()

	conditionals := make([]ConditionalConfig, 0)
	if conditionalLen == 0 {
		return conditionals, nil
	}
	conditionalsList := varOverrideMap.List()

	for _, conditional := range conditionalsList {
		conditionalMap := conditional.(map[string]interface{})
		conditionalContentMapList := conditionalMap["content"].(*schema.Set).List()
		if len(conditionalContentMapList) < 1 {
			continue
		}
		conditionalContentMap := conditionalContentMapList[0].(map[string]interface{})
		vc, err := BuildVariableFromSchema(conditionalContentMap)
		if err != nil {
			return nil, err
		}
		c := ConditionalConfig{
			Source:          conditionalMap["source"].(string),
			Condition:       conditionalMap["condition"].(string),
			VariableContent: *vc,
		}

		conditionals = append(conditionals, c)
	}

	return conditionals, nil
}
