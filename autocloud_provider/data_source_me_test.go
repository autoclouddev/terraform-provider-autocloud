package autocloud_provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceMe = `
data "autocloud_me" "user" {}
`

func TestAccDataSourceMe(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		// ref -> https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests/teststep
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMe,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.autocloud_me.user", "email", regexp.MustCompile("^enrique.enciso")),
				),
			},
		},
	})
}
