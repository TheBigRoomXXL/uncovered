package parsing

import (
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"uncovered/pkg/shared"
)

func Parse(report io.Reader, base string, contextSize int, filterFile func(string) bool) ([]shared.Patch, error) {
	// TODO: accept a list of report, and merge them
	coverage, err := ParseCobertura(report, base, filterFile)
	if err != nil {
		return nil, err
	}

	patchs := Transform(coverage)

	fileSystem := os.DirFS("/")
	patchs = AddSources(fileSystem, patchs, contextSize)

	return patchs, nil
}

// Transform a coverage map into a list of uncovered patch
func Transform(coverage shared.Coverage) []shared.Patch {
	patchs := []shared.Patch{}
	for path, file := range coverage {
		numbers := slices.Sorted(maps.Keys(file.Lines))

		// Group consecutive uncovered lines. Line numbers missing from the
		// report (blanks, comments, declarations) are not coverable, so they
		// do not break a patch: a patch ends on the first covered line.
		var patch *shared.Patch
		for _, number := range numbers {
			if file.Lines[number] {
				// We reach a covered line so it is the end of the current patch
				if patch != nil {
					patchs = append(patchs, *patch)
					patch = nil
				}
				continue
			}

			if patch == nil {
				patch = &shared.Patch{File: file.File, FileAbs: path, Start: number, End: number}
			} else {
				patch.End = number
			}
		}
		if patch != nil {
			patchs = append(patchs, *patch)
		}
	}

	patchs = SortPatch(patchs)
	return patchs
}

// Internally we use fs.FS so all path must be relative to /
// relative pass are interpelated to a "base" wich can be the cwd or a user inputed path
func ComputePath(base, path string) string {
	shared.PreCondition(filepath.IsAbs(base), "base should be absolute")

	var abs string
	if filepath.IsAbs(path) {
		abs = filepath.Clean(path)
	} else {
		abs = filepath.Join(base, path) // Join also Cleans
	}

	return strings.TrimPrefix(abs, "/")
}
