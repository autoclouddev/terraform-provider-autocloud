package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider_go"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider_plugin"
)

var (
	// Version can be updated by goreleaser on release
	version string = "dev"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	//debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
	experimental := flag.Bool("experimental", false, "allows experimental capabilities")
	flag.Parse()

	ctx := context.Background()
	providers_v5 := []func() tfprotov5.ProviderServer{
		// Example terraform-plugin-sdk/v2 providers
		provider.New("dev", *experimental)().GRPCProvider,
		provider_go.WithFlagGate(*experimental),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers_v5...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		context.Background(),
		muxServer.ProviderServer,
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	providers_v6 := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
	}

	// Example terraform-plugin-framework provider
	if *experimental {
		providers_v6 = append(providers_v6, providerserver.NewProtocol6(provider_plugin.New(version)()))
	}

	muxServer6, err := tf6muxserver.NewMuxServer(ctx, providers_v6...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = tf6server.Serve(
		"autocloud.io/autoclouddev/autocloud",
		func() tfprotov6.ProviderServer {
			return muxServer6
		},
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
