package rendering

import (
	"fmt"
	"strings"
	"uncovered/pkg/shared"
)

type ANSI struct {
	Reset  string
	Orange string
	Dim    string
}

var ANSIColors = ANSI{
	Reset:  "\033[0m",
	Orange: "\033[38;5;208m",
	Dim:    "\033[38;5;244m",
}

var ANSINone = ANSI{
	Reset:  "",
	Orange: "",
	Dim:    "",
}

// Display a patch a text, color enable ANSI coloring
func RenderPatchText(patch shared.Patch, color bool) string {
	ansi := ANSINone
	hint := "✗ "
	pad := "  "
	if color {
		ansi = ANSIColors
		hint = ""
		pad = ""
	}

	width := countDigits(patch.End + len(patch.SourceAfter))

	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s:%d", patch.File, patch.Start))
	if patch.Start != patch.End {
		result.WriteString(fmt.Sprintf("-%d", patch.End))
	}
	result.WriteString("\n")
	for i, lineBefore := range patch.SourceBefore {
		lineNo := patch.Start - len(patch.SourceBefore) + i
		result.WriteString(fmt.Sprintf("%s%*d | %s%s%s\n", ansi.Dim, width, lineNo, pad, lineBefore, ansi.Reset))
	}

	for i, line := range patch.Source {
		lineNo := patch.Start + i
		result.WriteString(fmt.Sprintf("%s%*d | %s%s%s\n", ansi.Orange, width, lineNo, hint, line, ansi.Reset))
	}

	for i, lineAfter := range patch.SourceAfter {
		lineNo := patch.End + i + 1
		result.WriteString(fmt.Sprintf("%s%*d | %s%s%s\n", ansi.Dim, width, lineNo, pad, lineAfter, ansi.Reset))
	}

	return result.String()
}

func countDigits(n int) int {
	if n == 0 {
		return 1
	}

	if n < 0 {
		n = -n
	}

	count := 0
	for n > 0 {
		count++
		n /= 10
	}
	return count
}
