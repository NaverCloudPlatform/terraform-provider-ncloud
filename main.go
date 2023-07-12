package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Sset to true to run the provider with support for debuggers like delve")
	flag.Parse()

	serverFactory, _, err := provider.ProtoV5ProviderServerFactory(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	if *debugFlag {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/NaverCloudPlatform/ncloud",
		serverFactory,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
