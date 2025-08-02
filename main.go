package main

import (
	"context"
	"flag"
	"log"

	"github.com/cysp/terraform-provider-censusworkspace/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// set by goreleaser.
var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/cysp/censusworkspace",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.Factory(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
