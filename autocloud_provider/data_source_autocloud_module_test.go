package autocloud_provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

const moduleName = "CloudfrontTest"
const source = "terraform-aws-modules/cloudfront/aws"
const version = "3.0.0"

const testAccDataSourceAutocloudModule = `
data "autocloud_module" "cloudfront" {
	filter  {
		name ="` + moduleName + `"
		source ="` + source + `"
		version ="` + version + `"
	}
}
`

func setupSdk() *autocloudsdk.Client {
	graphql := os.Getenv("SDK_GRAPHQL_HOST")
	appClient := os.Getenv("SDK_COGNITO_APP_CLIENT_ID")
	host := os.Getenv("SDK_API_HOST")
	username := os.Getenv("AUTOCLOUD_USERNAME")
	password := os.Getenv("AUTOCLOUD_PASSWORD")
	c, err := autocloudsdk.NewClient(&graphql, &host, &appClient)
	if err != nil {
		panic("sdk not initialized")
	}
	err = c.Login(&username, &password)
	if err != nil {
		panic("sdk not authentificated")
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
	fmt.Println(err)
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
			testAccPreCheck(t)
			setupSdk()
			createModule()
		},
		CheckDestroy: func(s *terraform.State) error {
			setupSdk()
			deleteModule()
			return nil
		},
		ProviderFactories: providerFactories,
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
						"data.autocloud_module.cloudfront", "form_config"),
				),
			},
		},
	})
}
