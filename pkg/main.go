package pkg

import (
	"fmt"
	"os"
	"uncovered/pkg/filtering"
	"uncovered/pkg/parsing"
	"uncovered/pkg/shared"
)

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

	config.Render(config.Ouput, patchs)

	return nil
}
