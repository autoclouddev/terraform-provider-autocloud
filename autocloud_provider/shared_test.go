package autocloud_provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// test function wrapper to verify errors are shown when there's an invalid configurations
// if the error is shown, the test is skipped. Otherwise, if the error is not shown, the test will fail
func validateErrors(t *testing.T, expectedError string, terraform string) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		ErrorCheck: func(err error) error {
			// if regex matches, do t.Skip instead of just passing the error through or something...
			re := regexp.MustCompile(expectedError)

			if re.MatchString(err.Error()) {
				t.Skipf("skipping test - The error was catched")
			}

			return err
		},
		Steps: []resource.TestStep{
			{
				Config: terraform,
			},
		},
	})
}
