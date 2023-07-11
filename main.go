package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/myklst/terraform-provider-st-utilities/utilities"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name test

func main() {
	providerserver.Serve(context.Background(), utilities.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/myklst/st-utilities",
	})
}
