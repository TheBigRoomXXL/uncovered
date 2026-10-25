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
	c := shared.NewConfig()
	c.Ouput = os.Stdout

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(io.Discard) // Parse now prints nothing on error
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "usage: %s [options] <cobertura-report.xml>\n\noptions:\n", args[0])
		flags.PrintDefaults()
	}

	renderName := "text"

	flags.Var((*stringList)(&c.Includes), "include", "only report on this file; repeat for several files")
	flags.Var((*stringList)(&c.Excludes), "exclude", "never report on this file; repeat for several files")
	flags.StringVar(&c.Base, "root", c.Base, "base used for relative paths in the report")
	flags.StringVar(&renderName, "render", renderName, "rendering format to use (text, markdown)")
	flags.IntVar(&c.ContextSize, "context-size", c.ContextSize, "how many line to display before and after each patch")
	flags.BoolVar(&c.Color, "color", c.Color, "colored output or not")

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
	c.ReportPath = flags.Arg(0)

	switch renderName {
	case "text":
		c.Render = rendering.NewRendererText(c.ContextSize, c.Color)
	case "markdown":
		c.Render = rendering.NewRendererMarkdown(c.ContextSize)
	default:
		return nil, fmt.Errorf("report can be text or markdown, got %s", renderName)
	}

	err := shared.ValidateConfig(c)
	if err != nil {
		return nil, err
	}

	return c, nil
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
