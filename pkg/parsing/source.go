package parsing

import (
	"bufio"
	"io/fs"
	"uncovered/pkg/models"
)

// ParseSources fills the source and its surrounding context of every patch
func ParseSources(sys fs.FS, patchs []models.Patch, contextSize int) []models.Patch {
	models.PreCondition(contextSize >= 0, "contexe size must be positive")

	// Sorting ensure we process patchs file by file and only increase line number
	SortPatch(patchs)

	for low := 0; low < len(patchs); {
		// Delimit the run of patchs belownging to the same file
		high := low
		for high < len(patchs) && patchs[high].File == patchs[low].File {
			models.PreConditionf(
				patchs[high].Start <= patchs[high].End,
				"patch has an empty range: %s:%d-%d",
				patchs[high].File, patchs[high].Start, patchs[high].End,
			)
			if high > low {
				models.PreConditionf(
					patchs[high-1].End < patchs[high].Start,
					"patchs overlap, they should have been merged: %s:%d-%d and %s:%d-%d",
					patchs[high-1].File, patchs[high-1].Start, patchs[high-1].End,
					patchs[high].File, patchs[high].Start, patchs[high].End,
				)
			}
			high++
		}

		parseSourcesFile(sys, patchs[low:high], contextSize)
		low = high
	}

	return patchs
}

// parseSourcesFile reads one file and dispatches each line to every patch that wants it.
// All patchs must belong to the same file, be sorted and be disjoint.
func parseSourcesFile(sys fs.FS, patchs []models.Patch, contextSize int) {
	models.PreCondition(len(patchs) > 0, "cannot process empty list")

	file, err := sys.Open(patchs[0].File)
	if err != nil {
		setSourceError(patchs, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	head := 0 // only patch after head still expects lines

	for head < len(patchs) && scanner.Scan() {
		lineNumber += 1

		// Once the head patch has all of it's line, we can incread the head
		// to avoid re-processing done patchs.
		for head < len(patchs) && lineNumber > patchs[head].End+contextSize {
			head += 1
		}

		for i := head; i < len(patchs); i++ {
			patch := &patchs[i]

			if lineNumber < patch.Start-contextSize {
				break // start is increasing, so if this patch is not ready, neither are the following ones
			}

			if lineNumber < patch.Start {
				patch.SourceBefore = append(patch.SourceBefore, scanner.Text())
			} else if lineNumber <= patch.End {
				patch.Source = append(patch.Source, scanner.Text())
			} else {
				patch.SourceAfter = append(patch.SourceAfter, scanner.Text())
			}
		}
	}

	if err := scanner.Err(); err != nil {
		setSourceError(patchs, err)
	}
}

// setSourceError report a file level failure on every patch of that file
func setSourceError(patchs []models.Patch, err error) {
	for i := range patchs {
		patchs[i].SourceError = err
	}
}
