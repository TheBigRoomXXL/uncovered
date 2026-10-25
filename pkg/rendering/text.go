package rendering

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"uncovered/pkg/shared"
)

type ANSI struct {
	Reset  string
	Orange string
	Red    string
	Dim    string
}

var ANSIColors = ANSI{
	Reset:  "\033[0m",
	Orange: "\033[38;5;208m",
	Red:    "\033[31m",
	Dim:    "\033[38;5;244m",
}

var ANSINone = ANSI{
	Reset:  "",
	Orange: "",
	Red:    "",
	Dim:    "",
}

func NewRendererText(contextSize int, color bool) shared.Renderer {
	// FS is always the full file system, we only use it to facilitate tests of
	// WithSources using an interface
	fileSystem := os.DirFS("/")

	return func(stream io.Writer, patchs []shared.Patch) {
		streamBuffered := bufio.NewWriter(stream)
		for patch := range WithSources(fileSystem, patchs, contextSize) {
			RenderPatchText(streamBuffered, patch, color)
			io.WriteString(stream, "\n")
			streamBuffered.Flush()
		}
	}
}

// Display a patch a text, color enable ANSI coloring
func RenderPatchText(stream io.Writer, patch shared.Patch, color bool) {
	ansi := ANSINone
	hint := "✗ "
	pad := "  "
	if color {
		ansi = ANSIColors
		hint = ""
		pad = ""
	}
	width := countDigits(patch.End + len(patch.SourceAfter))

	io.WriteString(stream, patch.Identifier())
	io.WriteString(stream, "\n")

	if patch.SourceError != nil {
		io.WriteString(stream, fmt.Sprintf("%serror retriving source: %s%s\n", ansi.Red, patch.SourceError, ansi.Reset))
		return
	}

	for i, lineBefore := range patch.SourceBefore {
		lineNo := patch.Start - len(patch.SourceBefore) + i
		io.WriteString(stream, fmt.Sprintf("%s%*d | %s%s%s\n", ansi.Dim, width, lineNo, pad, lineBefore, ansi.Reset))
	}

	for i, line := range patch.Source {
		lineNo := patch.Start + i
		io.WriteString(stream, fmt.Sprintf("%s%*d | %s%s%s\n", ansi.Orange, width, lineNo, hint, line, ansi.Reset))
	}

	for i, lineAfter := range patch.SourceAfter {
		lineNo := patch.End + i + 1
		io.WriteString(stream, fmt.Sprintf("%s%*d | %s%s%s\n", ansi.Dim, width, lineNo, pad, lineAfter, ansi.Reset))
	}
}

func countDigits(n int) int {
	shared.PreCondition(n > 0, "line number should always be geater than 0")
	if n == 0 {
		return 1
	}

	count := 0
	for n > 0 {
		count++
		n /= 10
	}
	return count
}
