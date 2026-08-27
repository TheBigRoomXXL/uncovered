package main

import (
	"bufio"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// ANSI escapes: the patch itself is orange, context keeps the default font
// color and only its line numbers are dimmed.
const (
	ansiReset  = "\033[0m"
	ansiOrange = "\033[38;5;208m"
	ansiDim    = "\033[38;5;244m"
)

func main() {
	config, err := NewConfigFromCLI(os.Args)
	if errors.Is(err, flag.ErrHelp) {
		// -h/-help is not an error: the flag set already printed the usage.
		return
	}
	if err != nil {
		log.Fatal("failed to parse arguments: ", err)
	}

	report, err := os.Open(config.reportPath)
	if err != nil {
		log.Fatal("failed to open report: ", err)
	}
	defer report.Close()

	untested, err := getUntestedPatchs(report, config.includes, config.excludes)
	if err != nil {
		log.Fatal("failed to parse report: ", err)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for _, patch := range untested {
		fmt.Fprintln(out, sPrintPatch(patch, config.contextSize))
	}
}

// Represent a continuous sequence of source code
type Patch struct {
	file  string
	start int
	end   int
}

// coberturaReport mirrors the subset of the cobertura XML schema we need:
// every class carries a filename and the per-line hit counts.
type coberturaReport struct {
	XMLName  xml.Name `xml:"coverage"`
	Packages []struct {
		Classes []struct {
			Filename string `xml:"filename,attr"`
			Lines    []struct {
				Number int `xml:"number,attr"`
				Hits   int `xml:"hits,attr"`
			} `xml:"lines>line"`
		} `xml:"classes>class"`
	} `xml:"packages>package"`
}

type Config struct {
	reportPath  string
	includes    []string
	excludes    []string
	contextSize int
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

// Build the config from cli options
func NewConfigFromCLI(args []string) (*Config, error) {
	// Default values
	config := &Config{
		includes:    []string{},
		excludes:    []string{},
		contextSize: 3,
	}

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "usage: %s [options] <cobertura-report.xml>\n\noptions:\n", args[0])
		flags.PrintDefaults()
	}
	flags.Var((*stringList)(&config.includes), "include", "only report on this file; repeat for several files")
	flags.Var((*stringList)(&config.excludes), "exclude", "never report on this file; repeat for several files")
	flags.IntVar(&config.contextSize, "context-size", 3, "how many line to display before and after each patch")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}

	// The report path is mandatory, the context size is optional.
	if flags.NArg() < 1 {
		flags.Usage()
		return nil, errors.New("missing <cobertura-report.xml>")
	}
	config.reportPath = flags.Arg(0)

	if flags.NArg() > 1 {
		contextSize, err := strconv.Atoi(flags.Arg(1))
		if err != nil || contextSize < 0 {
			flags.Usage()
			return nil, fmt.Errorf("context-size must be a positive integer, got %q", flags.Arg(1))
		}
		config.contextSize = contextSize
	}

	return config, nil
}

// Parse a cobertura XMl coverage report to extract all untest code patchs
func getUntestedPatchs(report io.Reader, includes, excludes []string) ([]Patch, error) {
	var parsed coberturaReport
	if err := xml.NewDecoder(report).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("cannot decode coverage report: %w", err)
	}

	patchs := []Patch{}
	for _, pkg := range parsed.Packages {
		for _, class := range pkg.Classes {
			if len(includes) > 0 && !slices.Contains(includes, class.Filename) {
				continue
			}
			if len(excludes) > 0 && slices.Contains(excludes, class.Filename) {
				continue
			}

			// A class may be reported in several chunks and lines are not
			// guaranteed to be ordered, so work on a sorted copy.
			lines := make([]struct {
				Number int
				Hits   int
			}, 0, len(class.Lines))
			for _, line := range class.Lines {
				lines = append(lines, struct {
					Number int
					Hits   int
				}{line.Number, line.Hits})
			}
			sort.Slice(lines, func(i, j int) bool { return lines[i].Number < lines[j].Number })

			// Group consecutive uncovered lines. Line numbers missing from the
			// report (blanks, comments, declarations) are not coverable, so they
			// do not break a patch: a patch ends on the first covered line.
			var current *Patch
			for _, line := range lines {
				if line.Hits > 0 {
					if current != nil {
						patchs = append(patchs, *current)
						current = nil
					}
					continue
				}
				if current == nil {
					current = &Patch{file: class.Filename, start: line.Number, end: line.Number}
				} else {
					current.end = line.Number
				}
			}
			if current != nil {
				patchs = append(patchs, *current)
			}
		}
	}

	return patchs, nil
}

// Display a patch in ANSI
// context determine how many lines befor and after the actual patch
// Context is default grey
// Patch is Orange
func sPrintPatch(patch Patch, contextSize int) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s:%d", patch.file, patch.start))
	if patch.start != patch.end {
		result.WriteString(fmt.Sprintf("-%d", patch.end))
	}
	result.WriteString("\n")

	file, err := os.Open(patch.file)
	if err != nil {
		result.WriteString(fmt.Sprintf("error reading file: %s\n", err))
		return result.String()
	}
	defer file.Close()

	from := patch.start - contextSize
	if from < 1 {
		from = 1
	}
	to := patch.end + contextSize
	width := len(strconv.Itoa(to))

	scanner := bufio.NewScanner(file)
	// PHP classes hold some very long generated lines, well past the 64KB
	// the scanner allows by default.
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	for lineNo := 1; lineNo <= to && scanner.Scan(); lineNo++ {
		if lineNo < from {
			continue
		}
		// Source is read as raw bytes and written back untouched: SSM PHP
		// files are ISO-8859-1 and must not be re-encoded.
		if lineNo >= patch.start && lineNo <= patch.end {
			result.WriteString(fmt.Sprintf("%s%*d | %s%s\n", ansiOrange, width, lineNo, scanner.Text(), ansiReset))
		} else {
			result.WriteString(fmt.Sprintf("%s%*d |%s %s\n", ansiDim, width, lineNo, scanner.Text(), ansiReset))
		}
	}
	if err := scanner.Err(); err != nil {
		result.WriteString(fmt.Sprintf("error reading file: %s\n", err))
	}

	return result.String()
}
