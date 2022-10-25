package autocloud_provider

import (
	"fmt"

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

func GetSdkIacCatalog(d *schema.ResourceData) autocloudsdk.IacCatalog {
	var labels = []string{}
	if labelValues, isLabelValuesOk := d.GetOk("labels"); isLabelValuesOk {
		list := labelValues.([]interface{})
		labels = make([]string, len(list))
		for i, labelValues := range list {
			labels[i] = labelValues.(string)
		}
	}

	generator := autocloudsdk.IacCatalog{
		Name:            d.Get("name").(string),
		Author:          d.Get("author").(string),
		Description:     d.Get("description").(string),
		Instructions:    d.Get("instructions").(string),
		Labels:          labels,
		FileDefinitions: GetSdkIacCatalogFileDefinitions(d),
		GitConfig:       GetSdkIacCatalogGitConfig(d),
		IacModuleIds:    GetSdkIacCatalogModuleIDs(d),
	}

	return generator
}

func GetSdkIacCatalogModuleIDs(d *schema.ResourceData) []string {
	var iacModulesIDs []string
	if autocloudModules, ok := d.GetOk("autocloud_module_1"); ok {
		list := autocloudModules.(*schema.Set).List()
		iacModulesIDs = make([]string, len(list))
		for i, autocloudModule := range list {
			var autocloudModuleMap = autocloudModule.(map[string]interface{})

			if val, ok := autocloudModuleMap["id"]; ok {
				iacModulesIDs[i] = val.(string)
			}
		}
	}

	return iacModulesIDs
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

			if val, ok := fileDefinitionMap["path_from_root"]; ok {
				fileDefinition.PathFromRoot = val.(string)
			}

			if val, ok := fileDefinitionMap["filename_template"]; ok {
				fileDefinition.FilenameTemplate = val.(string)
			}

			if val, ok := fileDefinitionMap["filename_vars"]; ok {
				var filenamesValueMap = val.(map[string]interface{})
				fileDefinition.FilenameVars = ConvertMap(filenamesValueMap)
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

			if val, ok := gitConfigMap["git_url_default"]; ok {
				gitConfig.GitURLDefault = val.(string)
			}

			if val, ok := gitConfigMap["git_url_options"]; ok {
				list := val.([]interface{})
				options := make([]string, len(list))
				for i, optionValue := range list {
					options[i] = optionValue.(string)
				}
				gitConfig.GitURLOptions = options
			}

			if val, ok := gitConfigMap["pull_request"]; ok {
				gitConfig.PullRequest = GetSdkIacCatalogGitConfigPR(val)
			}
		}
	}

	return gitConfig
}

func GetSdkIacModule(d *schema.ResourceData) autocloudsdk.IacModule {
	// note: the Template and Variables fields are calculated by the SDK
	iacModule := autocloudsdk.IacModule{
		Name:    d.Get("name").(string),
		Source:  d.Get("source").(string),
		Version: d.Get("version").(string),
	}

	return iacModule
}
