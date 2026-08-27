package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"uncovered/cli"
	"uncovered/pkg/parsing"
	"uncovered/pkg/rendering"
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

	report, err := os.Open(config.ReportPath)
	if err != nil {
		log.Fatal("failed to open report: ", err)
	}
	defer report.Close()

	untested, err := parsing.ParseCobertura(report, config.Includes, config.Excludes)
	if err != nil {
		log.Fatal("failed to parse report: ", err)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for _, patch := range untested {
		fmt.Fprintln(out, rendering.RenderTerminal(patch, config.ContextSize))
	}
}
