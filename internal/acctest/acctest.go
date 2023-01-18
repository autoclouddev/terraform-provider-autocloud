package acctest

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/joho/godotenv"
	provider "gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider_go"
	//provider_go "gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider_go"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.

//nolint:unparam // The error result is required, but intentionally always nil here
var ProviderFactories = map[string]func() (*schema.Provider, error){
	"autocloud": func() (*schema.Provider, error) {
		return provider.New("dev")(), nil
	},
	/*
		"tfe": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			mux, err := tfmux.NewMuxServer(
				ctx, provider.PluginProviderServer, testAccProvider.GRPCProvider,
			)
			if err != nil {
				return nil, err
			}

			return mux.ProviderServer(), nil
		},*/
}

func CreateMuxFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	prov := map[string]func() (tfprotov5.ProviderServer, error){
		"autocloud": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				// Example terraform-plugin-sdk/v2 providers
				provider.New("dev")().GRPCProvider,
				provider_go.PluginProviderServer,
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
	return prov
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
	os.Setenv("TF_LOG", "DEBUG")
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
