package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"uncovered/pkg"
	"uncovered/pkg/cli"
)

func main() {
	config, err := cli.NewConfigFromCLI(os.Args)
	if errors.Is(err, flag.ErrHelp) {
		// -h/-help is not an error: the flag set already printed the usage.
		return
	}
	if err != nil {
		log.Fatal("failed to parse arguments: ", err)
	}

	pkg.Run(config)
}
