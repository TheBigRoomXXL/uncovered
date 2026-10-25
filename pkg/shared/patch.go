package shared

import (
	"fmt"
	"io"
	"strings"
)

// Coverage of every file in the reports, keyed by the normalized absolute path
type Coverage = map[string]FileCoverage

// FileCoverage holds the coverage data of a single file
type FileCoverage struct {
	File  string       // path as written in the report, used for display
	Lines map[int]bool // key is line number, value is covered or not
}

// Represent a continuous sequence of source code
type Patch struct {
	File         string   `json:"file"`
	FileAbs      string   `json:"file_abs"`
	Start        int      `json:"start"` // Inclusive
	End          int      `json:"end"`   // Inclusive
	Source       []string `json:"source"`
	SourceBefore []string `json:"source_before"`
	SourceAfter  []string `json:"source_after"`
	SourceError  error    `json:"source_error"`
}

func (p *Patch) Identifier() string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s:%d", p.File, p.Start))
	if p.Start != p.End {
		result.WriteString(fmt.Sprintf("-%d", p.End))
	}
	return result.String()
}

// Renderer does not return an errors because it is already writing
// to ouput, so if an error occure, it is best place to cleanly incorporate
// the error into the programme output
// For example it can halde issue at the patch level rather than the entire batch
type Renderer = func(stream io.Writer, patchs []Patch)
