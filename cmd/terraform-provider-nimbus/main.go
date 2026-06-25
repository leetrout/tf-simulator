// Command terraform-provider-nimbus is a Terraform provider that manages
// resources in a running tfsim cloud via its HTTP API.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/leetrout/terraform-sim/cmd/terraform-provider-nimbus/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// version is overridable at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/leetrout/nimbus",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
