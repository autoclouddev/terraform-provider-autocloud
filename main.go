package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider_go"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
	experimental := flag.Bool("experimental", false, "allows experimental capabilities")
	flag.Parse()

	ctx := context.Background()
	providers := []func() tfprotov5.ProviderServer{
		// Example terraform-plugin-sdk/v2 providers
		provider.New("dev", *experimental)().GRPCProvider,
		provider_go.WithFlagGate(*experimental),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var serveOpts []tf5server.ServeOpt

	if *debugFlag {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"autocloud.io/autoclouddev/autocloud",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
