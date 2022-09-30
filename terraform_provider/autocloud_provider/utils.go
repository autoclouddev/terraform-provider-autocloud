package autocloud_provider

import (
	"autocloud_sdk"
	"fmt"

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

func GetSdkIacCatalog(d *schema.ResourceData) autocloud_sdk.IacCatalog {

	var labels = []string{}
	if labelValues, isLabelValuesOk := d.GetOk("labels"); isLabelValuesOk {
		list := labelValues.([]interface{})
		labels = make([]string, len(list))
		for i, labelValues := range list {
			labels[i] = labelValues.(string)
		}
	}

	generator := autocloud_sdk.IacCatalog{
		Name:            d.Get("name").(string),
		Author:          d.Get("author").(string),
		Slug:            d.Get("slug").(string),
		Description:     d.Get("description").(string),
		Instructions:    d.Get("instructions").(string),
		Version:         d.Get("version").(string),
		Source:          d.Get("source").(string),
		Template:        d.Get("template").(string),
		Labels:          labels,
		FileDefinitions: GetSdkIacCatalogFileDefinitions(d),
	}

	return generator
}

func GetSdkIacCatalogFileDefinitions(d *schema.ResourceData) []autocloud_sdk.IacCatalogFile {

	var fileDefinitions []autocloud_sdk.IacCatalogFile
	if fileDefinitionsValues, ok := d.GetOk("file"); ok {

		list := fileDefinitionsValues.(*schema.Set).List()
		fileDefinitions = make([]autocloud_sdk.IacCatalogFile, len(list))
		for i, fileDefinitionsValue := range list {
			var fileDefinitionMap = fileDefinitionsValue.(map[string]interface{})

			var fileDefinition = autocloud_sdk.IacCatalogFile{}

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
