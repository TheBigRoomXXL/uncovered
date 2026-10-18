package parsing

import (
	"encoding/xml"
	"fmt"
	"io"
	"uncovered/pkg/shared"
)

// coberturaReport mirrors the subset of the cobertura XML schema we need:
// every class carries a filename and the per-line hit counts.
// every tags of attributes we don't read is omited so that we cannot fail on parsing them
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

// Parse a cobertura XMl coverage report to extract the coverage of every file.
// Files are filtered on the path as written in the report, then keyed by their
// absolute path so that classes pointing to the same file are merged.
func ParseCobertura(report io.Reader, base string, filter func(filtepath string) bool) (shared.Coverage, error) {
	var parsed coberturaReport
	if err := xml.NewDecoder(report).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to decode coverage report: %w", err)
	}

	coverage := shared.Coverage{}
	for _, pkg := range parsed.Packages {
		for _, class := range pkg.Classes {
			if !filter(class.Filename) {
				continue
			}

			path := ComputePath(base, class.Filename)
			file, ok := coverage[path]
			if !ok {
				file = shared.FileCoverage{File: class.Filename, Lines: map[int]bool{}}
				coverage[path] = file
			}

			// A line reported several times is covered at least once
			for _, line := range class.Lines {
				file.Lines[line.Number] = file.Lines[line.Number] || line.Hits > 0
			}
		}
	}

	return coverage, nil
}
