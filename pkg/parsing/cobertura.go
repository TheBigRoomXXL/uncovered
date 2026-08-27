package parsing

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"uncovered/pkg/models"
)

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

// Parse a cobertura XMl coverage report to extract all untest code patchs
func ParseCobertura(report io.Reader, filterFile func(filtepath string) bool) ([]models.Patch, error) {
	var parsed coberturaReport
	if err := xml.NewDecoder(report).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("cannot decode coverage report: %w", err)
	}
	patchs := []models.Patch{}
	for _, pkg := range parsed.Packages {
		for _, class := range pkg.Classes {
			if filterFile(class.Filename) {
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
			var current *models.Patch
			for _, line := range lines {
				if line.Hits > 0 {
					if current != nil {
						patchs = append(patchs, *current)
						current = nil
					}
					continue
				}
				if current == nil {
					current = &models.Patch{File: class.Filename, Start: line.Number, End: line.Number}
				} else {
					current.End = line.Number
				}
			}
			if current != nil {
				patchs = append(patchs, *current)
			}
		}
	}

	return patchs, nil
}
