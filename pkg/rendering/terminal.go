package rendering

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"uncovered/pkg/models"
)

// ANSI escapes: the patch itself is orange, context keeps the default font
// color and only its line numbers are dimmed.
const (
	ansiReset  = "\033[0m"
	ansiOrange = "\033[38;5;208m"
	ansiDim    = "\033[38;5;244m"
)

// Display a patch in ANSI
// context determine how many lines befor and after the actual patch
// Context is default grey
// Patch is Orange
func RenderTerminal(patch models.Patch, contextSize int) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s:%d", patch.File, patch.Start))
	if patch.Start != patch.End {
		result.WriteString(fmt.Sprintf("-%d", patch.End))
	}
	result.WriteString("\n")

	file, err := os.Open(patch.File)
	if err != nil {
		result.WriteString(fmt.Sprintf("error reading file: %s\n", err))
		return result.String()
	}
	defer file.Close()

	from := patch.Start - contextSize
	if from < 1 {
		from = 1
	}
	to := patch.End + contextSize
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
		if lineNo >= patch.Start && lineNo <= patch.End {
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
