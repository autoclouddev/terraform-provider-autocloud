package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/iac_module"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}

func ConvertMap(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)

	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}

	return mapString
}

func ToStringSlice(sliceInterface []interface{}) []string {
	values := make([]string, len(sliceInterface))
	for idx, value := range sliceInterface {
		values[idx] = value.(string)
	}
	return values
}

func ToStringMap(str string) (map[string]string, error) {
	outputMap := map[string]string{}
	err := json.Unmarshal([]byte(str), &outputMap)
	if err != nil {
		return nil, err
	}

	return outputMap, nil
}

func GetSdkIacCatalogModules(d *schema.ResourceData) []autocloudsdk.IacCatalogModule {
	var iacModules []autocloudsdk.IacCatalogModule
	if autocloudModules, ok := d.GetOk("autocloud_module"); ok {
		list := autocloudModules.(*schema.Set).List()
		iacModules = make([]autocloudsdk.IacCatalogModule, len(list))
		for i, autocloudModuleData := range list {
			var autocloudModuleMap = autocloudModuleData.(map[string]interface{})

			autocloudModule := autocloudsdk.IacCatalogModule{}
			if val, ok := autocloudModuleMap["id"]; ok {
				autocloudModule.ID = val.(string)
			}
			if val, ok := autocloudModuleMap["form_config"]; ok {
				autocloudModule.Variables = val.(string)
			}
			if val, ok := autocloudModuleMap["template_config"]; ok {
				autocloudModule.Template = val.(string)
			}

			if val, ok := d.GetOk("tags_variable"); ok {
				autocloudModule.TagsVariable = val.(string)
			}

			if orderValues, isorderValuesOk := d.GetOk("display_order"); isorderValuesOk {
				list := orderValues.([]interface{})
				autocloudModule.DisplayOrder = ToStringSlice(list)
			}

			iacModules[i] = autocloudModule
		}
	}

	return iacModules
}

func GetSdkIacCatalogFileDefinitions(d *schema.ResourceData) ([]generator.IacCatalogFile, error) {
	var fileDefinitions []generator.IacCatalogFile
	if fileDefinitionsValues, ok := d.GetOk("file"); ok {
		list := fileDefinitionsValues.(*schema.Set).List()
		fileDefinitions = make([]generator.IacCatalogFile, 0)
		for _, fileDefinitionsValue := range list {
			var fileDefinitionMap = fileDefinitionsValue.(map[string]interface{})
			var fileDefinition = generator.IacCatalogFile{}

			if val, ok := fileDefinitionMap["action"]; ok {
				fileDefinition.Action = val.(string)
			}

			// prevent empty values
			if fileDefinition.Action == "" {
				continue
			}

			if val, ok := fileDefinitionMap["destination"]; ok {
				fileDefinition.Destination = val.(string)
			}

			if val, ok := fileDefinitionMap["variables"]; ok {
				var variablesMap = val.(map[string]interface{})
				fileDefinition.Variables = ConvertMap(variablesMap)
			}
			if val, ok := fileDefinitionMap["modules"]; ok {
				var data = val.([]interface{})
				fileDefinition.Modules = ToStringSlice(data)
			}

			if val, ok := fileDefinitionMap["content"]; ok {
				fileDefinition.Content = val.(string)
			}

			if val, ok := fileDefinitionMap["header"]; ok {
				fileDefinition.Header = val.(string)
			}

			if val, ok := fileDefinitionMap["footer"]; ok {
				fileDefinition.Footer = val.(string)
			}

			if len(fileDefinition.Modules) == 0 && (fileDefinition.Header != "" || fileDefinition.Footer != "") {
				return nil, errors.New("modules can not be empty when using header or footer attributes")
			}

			if len(fileDefinition.Modules) == 0 && fileDefinition.Content == "" {
				return nil, errors.New("file block should contain content or modules attributes")
			}

			if len(fileDefinition.Modules) > 0 && fileDefinition.Content != "" {
				return nil, errors.New("file block should contain content or modules attributes, but not both")
			}

			fileDefinitions = append(fileDefinitions, fileDefinition)
		}
	}

	return fileDefinitions, nil
}

func GetSdkIacCatalogGitConfigPR(pullRequestConfigValues interface{}) generator.IacCatalogGitConfigPR {
	var pullRequestConfig generator.IacCatalogGitConfigPR
	list := pullRequestConfigValues.(*schema.Set).List()
	for _, pullRequestConfigValue := range list {
		var pullRequestConfigMap, ok = pullRequestConfigValue.(map[string]interface{})

		// Prevent read the entire empty property
		if !ok {
			return pullRequestConfig
		}

		if val, ok := pullRequestConfigMap["title"]; ok {
			fmt.Print(pullRequestConfigMap["title"])
			pullRequestConfig.Title = val.(string)
		}

		if val, ok := pullRequestConfigMap["commit_message_template"]; ok {
			pullRequestConfig.CommitMessageTemplate = val.(string)
		}

		if val, ok := pullRequestConfigMap["body"]; ok {
			pullRequestConfig.Body = val.(string)
		}

		if val, ok := pullRequestConfigMap["variables"]; ok {
			var pullRequestConfigVariablesMap = val.(map[string]interface{})
			pullRequestConfig.Variables = ConvertMap(pullRequestConfigVariablesMap)
		}
	}

	return pullRequestConfig
}

func GetSdkIacCatalogGitConfig(d *schema.ResourceData) *generator.IacCatalogGitConfig {
	var gitConfig generator.IacCatalogGitConfig
	if gitConfigValues, ok := d.GetOk("git_config"); ok {
		list := gitConfigValues.(*schema.Set).List()
		for _, gitConfigValue := range list {
			var gitConfigMap = gitConfigValue.(map[string]interface{})

			if val, ok := gitConfigMap["destination_branch"]; ok {
				fmt.Print(gitConfigMap["destination_branch"])
				gitConfig.DestinationBranch = val.(string)
			}

			if val, ok := gitConfigMap["git_url_options"]; ok {
				list := val.([]interface{})
				options := make([]string, len(list))
				for i, optionValue := range list {
					options[i] = optionValue.(string)
				}
				gitConfig.GitURLOptions = options
				if len(options) == 1 {
					gitConfig.GitURLDefault = options[0]
				}
			}

			// if there is only one option, the default repo shouldn't be available
			if val, ok := gitConfigMap["git_url_default"]; ok && len(gitConfig.GitURLDefault) == 0 {
				gitConfig.GitURLDefault = val.(string)
			}

			if val, ok := gitConfigMap["pull_request"]; ok {
				pr := GetSdkIacCatalogGitConfigPR(val)
				gitConfig.PullRequest = &pr
			}
		}
	} else {
		return nil
	}

	return &gitConfig
}

func GetSdkIacModuleInput(d *schema.ResourceData) iac_module.ModuleInput {
	iacModule := iac_module.ModuleInput{
		Name:         d.Get("name").(string),
		Source:       d.Get("source").(string),
		Version:      d.Get("version").(string),
		TagsVariable: d.Get("tags_variable").(string),
		Header:       d.Get("header").(string),
		Footer:       d.Get("footer").(string),
	}

	return iacModule
}

func MergeSchemas(a, b map[string]*schema.Schema) map[string]*schema.Schema {
	merged := make(map[string]*schema.Schema)
	for k, v := range a {
		merged[k] = v
	}
	for k, v := range b {
		merged[k] = v
	}
	return merged
}

func ParseVariables(str string) ([]generator.FormShape, error) {
	vars := []generator.FormShape{}
	err := json.Unmarshal([]byte(str), &vars)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

func GetVariablesIdMap(str string) (map[string]string, error) {
	vars, err := ParseVariables(str)
	if err != nil {
		return nil, err
	}

	varsMap := make(map[string]string)
	for _, v := range vars {
		varName, err := GetVariableID(v.ID)
		if err == nil {
			varsMap[varName] = v.ID
		}
	}

	return varsMap, nil
}

// variables id follow the pattern "<source module>.<variable name>""
func GetVariableID(variableKey string) (string, error) {
	keyValue := strings.Split(variableKey, ".")
	if IsValidId(variableKey) && len(keyValue) == 2 {
		return keyValue[1], nil
	}

	return "", errors.New("Invalid Key")
}

// marshals and converts an object into a compacted JSON string
func ToJsonString(obj any) (string, error) {
	jsonDoc, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	jsonBuffer := &bytes.Buffer{}
	if err := json.Compact(jsonBuffer, jsonDoc); err != nil {
		return "", err
	}
	return jsonBuffer.String(), nil
}

func ToJsonStringNoError(obj any) string {
	jsonStr, _ := ToJsonString(obj)
	return jsonStr
}

// remove empty spaces from a JSON string
//
//nolint:golint,unused
func compactJson(jsonStr string) string {
	jsonBuffer := &bytes.Buffer{}
	if err := json.Compact(jsonBuffer, []byte(jsonStr)); err != nil {
		panic(err)
	}
	return jsonBuffer.String()
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

// Verify that variable matches with a module.id pattern
func IsValidId(ov string) bool {
	pattern := "^[a-zA-Z0-9_]+\\.{1}[a-zA-Z0-9_]+$"

	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)
	result := re.MatchString(ov)
	return result
}

// Verify that variable matches with a reference
func HasReference(ov string) bool {
	pattern := "[a-zA-Z0-9_]+\\.{1}variables\\.{1}[a-zA-Z0-9_]+"

	// Compile the regular expression pattern
	re := regexp.MustCompile(pattern)
	result := re.MatchString(ov)
	return result
}

func SaveLogs(name string, data interface{}) {
	fmt.Printf("[saving data at %s file] -> %s", name, data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error writing data ->r", err)
	}

	os.WriteFile(name, []byte(string(jsonData)), 0600)
}

// Iterate through the second map and add its key-value pairs to the first map
// If a key already exists in the first map, its value will be updated with the value from the second map
func MergeMaps(m1, m2 *map[string]string) {
	for k, v := range *m2 {
		(*m1)[k] = v
	}
}

func LoadData[T any](testCase string) (out T, err error) {
	var testData T
	file, err := os.ReadFile(testCase)
	if err != nil {
		return testData, errors.New("error reading file")
	}

	err = json.Unmarshal(file, &testData)

	if err != nil {
		return testData, errors.New("error loading file")
	}
	return testData, nil
}
