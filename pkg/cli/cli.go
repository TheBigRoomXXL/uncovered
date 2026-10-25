package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"uncovered/pkg"
	"uncovered/pkg/rendering"
	"uncovered/pkg/shared"
)

func RunFromCLI() {
	config, err := NewConfigFromCLI(os.Args)
	if errors.Is(err, flag.ErrHelp) {
		// -h/-help is not an error: the flag set already printed the usage.
		return
	}
	if err != nil {
		fmt.Println("usage error: ", err)
		os.Exit(2)
	}

	err = pkg.Run(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Build the config from cli options
func NewConfigFromCLI(args []string) (*shared.Config, error) {
	// Default values
	config := shared.NewConfig()
	config.Ouput = os.Stdout

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(io.Discard) // Parse now prints nothing on error
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "usage: %s [options] <cobertura-report.xml>\n\noptions:\n", args[0])
		flags.PrintDefaults()
	}

	flags.Var((*stringList)(&config.Includes), "include", "only report on this file; repeat for several files")
	flags.Var((*stringList)(&config.Excludes), "exclude", "never report on this file; repeat for several files")
	flags.IntVar(&config.ContextSize, "context-size", config.ContextSize, "how many line to display before and after each patch")
	flags.StringVar(&config.Base, "root", config.Base, "base used for relative paths in the report")

	if err := flags.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			flags.SetOutput(os.Stderr)
			flags.Usage()
			return nil, err
		}
		return nil, err
	}

	if flags.NArg() != 1 {
		return nil, fmt.Errorf("expected exactly one <report-path> as argument, got %s", flags.Args())
	}
	config.ReportPath = flags.Arg(0)

	config.Render = rendering.NewRendererText(config.ContextSize, true)

	err := shared.ValidateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

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
