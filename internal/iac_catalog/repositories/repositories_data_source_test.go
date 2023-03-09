package repositories_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
)

const testAccDataSourceScaffolding2 = `
data "autocloud_github_repos" "repos" {}
`

func TestAccDataSourceScaffolding2(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
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
