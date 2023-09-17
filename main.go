package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf6server"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Sset to true to run the provider with support for debuggers like delve")
	flag.Parse()

	serverFactory, _, err := provider.ProtoV6ProviderServerFactory(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if *debugFlag {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/NaverCloudPlatform/ncloud",
		serverFactory,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
