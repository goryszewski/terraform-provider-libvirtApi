package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"terraform-provider-libvirtapi/internal/provider"
)

var version string = "1.0.0"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set true to run debuggers")
	opts := providerserver.ServeOpts{
		Address: "github.com/goryszewski/libvirtapi",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}

}
