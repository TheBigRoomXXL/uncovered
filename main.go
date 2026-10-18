package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"uncovered/pkg/cli"
	"uncovered/pkg/filtering"
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

	filter := filtering.Filter(config.Includes, config.Excludes)
	patchs, err := parsing.Parse(report, config.Base, config.ContextSize, filter)
	if err != nil {
		log.Fatal("failed to parse report: ", err)
	}
	for _, patch := range patchs {
		// TODO: add color to config, based on env.
		fmt.Println(rendering.RenderPatchText(patch, true))
	}
}
