package autocloud_provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutocloudModule(t *testing.T) {
	t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAutocloudModule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"autocloud_module.foo", "name", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccAutocloudModule = `
resource "autocloud_module" "foo" {
  name = "bar"
}
`
