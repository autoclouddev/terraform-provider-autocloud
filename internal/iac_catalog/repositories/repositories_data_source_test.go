package repositories_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
)

const testAccDataSourceScaffolding2 = `
data "autocloud_github_repos" "repos" {}
`

func TestAccDataSourceScaffolding2(t *testing.T) {
	resourceName := "data.autocloud_github_repos.repos"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		// ref -> https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests/teststep
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceScaffolding2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "data.#"),
					testAccLogAllEData(resourceName),
				),
			},
		},
	})
}

func testAccLogAllEData(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		attributes := rs.Primary.Attributes
		apiHost := os.Getenv("AUTOCLOUD_API")
		token := os.Getenv("AUTOCLOUD_TOKEN")

		c, err := autocloudsdk.NewClient(&apiHost, &token)
		if err != nil {
			return err
		}
		repos, errGettingList := c.GitRepository.List()
		if errGettingList != nil {
			return errGettingList
		}

		for index, repo := range *repos {
			keyId := fmt.Sprintf("data.%v.id", index)
			attributeIdValue := attributes[keyId]
			if fmt.Sprint(repo.ID) != attributeIdValue {
				return fmt.Errorf("Ids don't match. State id: %v, repo id: %v", attributeIdValue, repo.ID)
			}
		}

		return nil
	}
}
