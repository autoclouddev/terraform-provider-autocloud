package acctest

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/joho/godotenv"
	provider "gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.

//nolint:unparam // The error result is required, but intentionally always nil here
var ProviderFactories = map[string]func() (*schema.Provider, error){
	"autocloud": func() (*schema.Provider, error) {
		return provider.New("dev")(), nil
	},
}

func TestAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	_, thisFileName, _, _ := runtime.Caller(0) //nolint:golint,dogsled
	thisFileDirectory := path.Join(path.Dir(thisFileName))
	err := godotenv.Load(path.Join(thisFileDirectory, "/../../.env"))
	if err != nil {
		t.Fatalf("cant load .env file, err: %s", err)
	}
}

// test function wrapper to verify errors are shown when there's an invalid configurations
// if the error is shown, the test is skipped. Otherwise, if the error is not shown, the test will fail
func ValidateErrors(t *testing.T, expectedError error, terraform string) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		ErrorCheck: func(err error) error {
			if strings.Contains(err.Error(), expectedError.Error()) {
				fmt.Println("successful test - The error was catched")
				return nil

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
