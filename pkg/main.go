package pkg

import (
	"fmt"
	"log"
	"os"
	"uncovered/pkg/filtering"
	"uncovered/pkg/parsing"
	"uncovered/pkg/rendering"
	"uncovered/pkg/shared"
)

func Run(config *shared.Config) {
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
