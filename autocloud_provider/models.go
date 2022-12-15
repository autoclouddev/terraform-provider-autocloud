package autocloud_provider

import autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"

type BluePrintConfig struct {
	Id        string                   `json:"id"`
	Variables []autocloudsdk.FormShape `json:"variables"`
	Children  []BluePrintConfig        `json:"children"`
}
