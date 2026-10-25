package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"uncovered/pkg/cli"
	"uncovered/pkg/filtering"
	"uncovered/pkg/parsing"
	"uncovered/pkg/rendering"
	"uncovered/pkg/shared"
)

func main() {
	config, err := cli.NewConfigFromCLI(os.Args)
	if errors.Is(err, flag.ErrHelp) {
		// -h/-help is not an error: the flag set already printed the usage.
		return
	}

	if err != nil {
		fmt.Println("usage error: ", err)
		os.Exit(2)
	}

	err = Run(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func Run(config *shared.Config) error {
	report, err := os.Open(config.ReportPath)
	if err != nil {
		return fmt.Errorf("failed to open report: %w", err)
	}
	defer report.Close()

	filter := filtering.Filter(config.Includes, config.Excludes)
	patchs, err := parsing.Parse(report, config.Base, filter)
	if err != nil {
		return fmt.Errorf("failed to parse report: %w", err)
	}

	fileSystem := os.DirFS("/")
	render := rendering.NewRendererText(fileSystem, config.ContextSize, true)
	render(os.Stdout, patchs)

	return nil
}
