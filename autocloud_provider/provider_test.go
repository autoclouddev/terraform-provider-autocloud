package autocloud_provider

// ref -> https://www.terraform.io/plugin/sdkv2/testing

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"

	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
}

// NOTE: to make this test work, it requires to create an .env file in ../.env
// and set the value AUTOCLOUD_API to http://localhost:8080/api/v.0.0.1
func TestProviderEndpoint(t *testing.T) {
	emptyResourceConfig := terraform.NewResourceConfigRaw(map[string]interface{}{})

	// 1 - initialize the provider WITHOUT an endpoint value in the .env file or in the provider configuration
	autocloudProvider := New("dev")()
	diags := autocloudProvider.Configure(context.Background(), emptyResourceConfig)

	assert.NotNil(t, diags)
	assert.Equal(t, "Autocloud API Endpoint is empty", diags[0].Summary)

	// load .env
	testAccPreCheck(t)

	// 2 - initialize the provider WITHOUT a given endpoint but WITH a value in the .env
	autocloudProvider = New("dev")()
	diags = autocloudProvider.Configure(context.Background(), emptyResourceConfig)
	sdkClient := autocloudProvider.Meta().(*autocloudsdk.Client)

	assert.Nil(t, diags)
	assert.NotNil(t, sdkClient)
	assert.Equal(t, "http://localhost:8080/api/v.0.0.1", sdkClient.HostURL)

	// 3 - initialize the provider WITH a given endpoint (it shouldn't use the .env value)
	expectedEndpoint := "https://api.autocloud.domain.com/api/v.0.0.1"
	providerConfiguration := map[string]interface{}{
		"endpoint": expectedEndpoint,
	}
	resourceConfig := terraform.NewResourceConfigRaw(providerConfiguration)
	diags = autocloudProvider.Configure(context.Background(), resourceConfig)

	sdkClient = autocloudProvider.Meta().(*autocloudsdk.Client)

	assert.Nil(t, diags)
	assert.NotNil(t, sdkClient)
	assert.Equal(t, expectedEndpoint, sdkClient.HostURL)
}
