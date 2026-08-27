package cli

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"uncovered/pkg/models"
)

// stringList collects a repeatable string option: every occurrence of the
// flag appends instead of overwriting the previous value.
type stringList []string

func (list *stringList) String() string {
	if list == nil {
		return ""
	}
	return strings.Join(*list, ",")
}

func (list *stringList) Set(value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	*list = append(*list, value)
	return nil
}

// Build the config from cli options
func NewConfigFromCLI(args []string) (*models.Config, error) {
	// Default values
	config := &models.Config{
		Includes:    []string{},
		Excludes:    []string{},
		ContextSize: 3,
	}

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "usage: %s [options] <cobertura-report.xml>\n\noptions:\n", args[0])
		flags.PrintDefaults()
	}
	flags.Var((*stringList)(&config.Includes), "include", "only report on this file; repeat for several files")
	flags.Var((*stringList)(&config.Excludes), "exclude", "never report on this file; repeat for several files")
	flags.IntVar(&config.ContextSize, "context-size", 3, "how many line to display before and after each patch")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}

	// The report path is mandatory, the context size is optional.
	if flags.NArg() < 1 {
		flags.Usage()
		return nil, errors.New("missing <cobertura-report.xml>")
	}
	config.ReportPath = flags.Arg(0)

	if flags.NArg() > 1 {
		contextSize, err := strconv.Atoi(flags.Arg(1))
		if err != nil || contextSize < 0 {
			flags.Usage()
			return nil, fmt.Errorf("context-size must be a positive integer, got %q", flags.Arg(1))
		}
		config.ContextSize = contextSize
	}

	return config, nil
}
