package autocloud_provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceScaffolding2 = `
data "autocloud_github_repos" "repos" {}
`

func TestAccDataSourceScaffolding2(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		// ref -> https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests/teststep
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceScaffolding2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.autocloud_github_repos.repos", "data.#"),
				),
			},
		},
	})
}
