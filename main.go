package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/myklst/terraform-provider-st-utilities/utilities"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name st-utilities

// Version info injected by goreleaser ldflags.
var (
	version string = "dev"
	commit  string = "none"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Start provider in debug mode")
	flag.Parse()

	err := providerserver.Serve(context.Background(), utilities.New, providerserver.ServeOpts{
		Address:         "registry.terraform.io/myklst/st-utilities",
		Debug:           *debugFlag,
		ProtocolVersion: 6,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("st-utilities provider %s (commit %s)", version, commit)
}
