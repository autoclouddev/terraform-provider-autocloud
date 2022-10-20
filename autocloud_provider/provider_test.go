package autocloud_provider

// ref -> https://www.terraform.io/plugin/sdkv2/testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joho/godotenv"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.

//nolint:unparam // The error result is required, but intentionally always nil here
var providerFactories = map[string]func() (*schema.Provider, error){
	"autocloud": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	fmt.Printf("Sdk url: %#v\n", os.Getenv("SDK_API_URL"))
}
