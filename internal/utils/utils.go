package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"

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

func GetSdkIacCatalogFileDefinitions(d *schema.ResourceData) []autocloudsdk.IacCatalogFile {
	var fileDefinitions []autocloudsdk.IacCatalogFile
	if fileDefinitionsValues, ok := d.GetOk("file"); ok {
		list := fileDefinitionsValues.(*schema.Set).List()
		fileDefinitions = make([]autocloudsdk.IacCatalogFile, len(list))
		for i, fileDefinitionsValue := range list {
			var fileDefinitionMap = fileDefinitionsValue.(map[string]interface{})

			var fileDefinition = autocloudsdk.IacCatalogFile{}

			if val, ok := fileDefinitionMap["action"]; ok {
				fileDefinition.Action = val.(string)
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

			fileDefinitions[i] = fileDefinition
		}
	}

	return fileDefinitions
}

func GetSdkIacCatalogGitConfigPR(pullRequestConfigValues interface{}) autocloudsdk.IacCatalogGitConfigPR {
	var pullRequestConfig autocloudsdk.IacCatalogGitConfigPR
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

func GetSdkIacCatalogGitConfig(d *schema.ResourceData) autocloudsdk.IacCatalogGitConfig {
	var gitConfig autocloudsdk.IacCatalogGitConfig
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
				gitConfig.PullRequest = GetSdkIacCatalogGitConfigPR(val)
			}
		}
	}

	return gitConfig
}

func GetSdkIacModule(d *schema.ResourceData) autocloudsdk.IacModule {
	var displayOrder = []string{}
	if orderValues, isorderValuesOk := d.GetOk("display_order"); isorderValuesOk {
		list := orderValues.([]interface{})
		displayOrder = ToStringSlice(list)
	}

	// note: the Template and Variables fields are calculated by the SDK
	iacModule := autocloudsdk.IacModule{
		Name:         d.Get("name").(string),
		Source:       d.Get("source").(string),
		Version:      d.Get("version").(string),
		TagsVariable: d.Get("tags_variable").(string),
		DisplayOrder: displayOrder,
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

func ParseVariables(str string) ([]autocloudsdk.FormShape, error) {
	vars := []autocloudsdk.FormShape{}
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
	if len(keyValue) == 2 {
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
