package autocloud_module_test

/*

// Commented until we have a definition for this feature, this is about the data resource, not the "resource" resource of iac_module

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
)

const moduleName = "CloudfrontTest"
const source = "terraform-aws-modules/cloudfront/aws"
const version = "3.0.0"

const testAccDataSourceAutocloudModule = `
data "autocloud_module" "cloudfront" {
	name = "cloudfrontExample"
	filter  {
		name ="` + moduleName + `"
		source ="` + source + `"
		version ="` + version + `"
	}
}
`

func setupSdk() *autocloudsdk.Client {
	host := os.Getenv("AUTOCLOUD_API")
	token := os.Getenv("AUTOCLOUD_TOKEN")
	c, err := autocloudsdk.NewClient(&host, &token)
	if err != nil {
		panic("sdk not initialized")
	}
	return c
}

func createModule() {
	c := setupSdk()
	iacModule := autocloudsdk.IacModule{
		Name:    moduleName,
		Source:  source,
		Version: version,
	}
	_, err := c.CreateModule(&iacModule)
	//fmt.Println(err)
	if err != nil {
		panic("module not created")
	}
}

func deleteModule() {
	c := setupSdk()
	modules, _ := c.GetModules()
	moduleId := ""
	for _, m := range modules {
		if m.Name == moduleName {
			moduleId = m.ID
			break
		}
	}
	err := c.DeleteModule(moduleId)
	if err != nil {
		panic("module not deleted")
	}
}

func TestAccDataSourceAutocloudModule(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			setupSdk()
			createModule()
		},
		CheckDestroy: func(s *terraform.State) error {
			setupSdk()
			deleteModule()
			return nil
		},
		ProviderFactories: acctest.ProviderFactories,
		// ref -> https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests/teststep
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutocloudModule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.autocloud_module.cloudfront", "name", regexp.MustCompile(moduleName)),
					resource.TestMatchResourceAttr(
						"data.autocloud_module.cloudfront", "version", regexp.MustCompile(version)),
					resource.TestMatchResourceAttr(
						"data.autocloud_module.cloudfront", "source", regexp.MustCompile(source)),
					resource.TestCheckResourceAttrSet(
						"data.autocloud_module.cloudfront", "template"),
					resource.TestCheckResourceAttrSet(
						"data.autocloud_module.cloudfront", "template_config"),
					resource.TestCheckResourceAttrSet(
						"data.autocloud_module.cloudfront", "variables.%"),
					resource.TestCheckResourceAttr(
						"data.autocloud_module.cloudfront", "variables.is_ipv6_enabled", "CloudfrontTest.is_ipv6_enabled"),
					resource.TestCheckResourceAttrSet(
						"data.autocloud_module.cloudfront", "blueprint_config"),
				),
			},
		},
	})
}
*/
