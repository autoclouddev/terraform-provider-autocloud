package autocloud_provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceScaffolding2(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		// ref -> https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests/teststep
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceScaffolding2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.autocloud_github_repos", "repos", regexp.MustCompile("^terraform-generator-test")),
				),
			},
		},
	})
}

const testAccDataSourceScaffolding2 = `
data "autocloud_github_repos" "repos" {}
`
