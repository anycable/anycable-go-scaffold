package main

import (
	"fmt"
	"os"

	acli "github.com/anycable/anycable-go/cli"
	"github.com/anycable/mycable/pkg/cli"
	"github.com/anycable/mycable/pkg/config"
	"github.com/anycable/mycable/pkg/version"

	"github.com/apex/log"

	_ "github.com/anycable/anycable-go/diagnostics"
)

func main() {
	conf := config.NewConfig()

	anyconf, err, ok := acli.NewConfigFromCLI(
		os.Args,
		acli.WithCLIName("mycable"),
		acli.WithCLIUsageHeader("MyCable, the custom AnyCable-Go build"),
		acli.WithCLIVersion(version.Version()),
		acli.WithCLICustomOptions(cli.CustomOptions(conf)),
	)

	if err != nil {
		log.Fatalf("%v", err)
	}

	if ok {
		os.Exit(0)
	}

	if err := cli.Run(conf, anyconf); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
