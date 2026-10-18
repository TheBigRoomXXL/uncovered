package parsing

import (
	"bufio"
	"io/fs"
	"uncovered/pkg/shared"
)

// AddSources fills the source and its surrounding context of every patch
// TODO: addd suport for filterPatchs
func AddSources(sys fs.FS, patchs []shared.Patch, contextSize int) []shared.Patch {
	shared.PreCondition(contextSize >= 0, "contex size must be positive")

	// Sorting ensure we process patchs file by file and only increase line number
	SortPatch(patchs)

	for low := 0; low < len(patchs); {
		// Delimit the run of patchs belownging to the same file
		high := low
		for high < len(patchs) && patchs[high].FileAbs == patchs[low].FileAbs {
			shared.PreConditionf(
				patchs[high].Start <= patchs[high].End,
				"patch has an empty range: %s:%d-%d",
				patchs[high].File, patchs[high].Start, patchs[high].End,
			)
			if high > low {
				shared.PreConditionf(
					patchs[high-1].End < patchs[high].Start,
					"patchs overlap, they should have been merged: %s:%d-%d and %s:%d-%d",
					patchs[high-1].File, patchs[high-1].Start, patchs[high-1].End,
					patchs[high].File, patchs[high].Start, patchs[high].End,
				)
				shared.PreConditionf(
					patchs[high-1].End+1 != patchs[high].Start,
					"patchs adjacent, they should have been merged: %s:%d-%d and %s:%d-%d",
					patchs[high-1].File, patchs[high-1].Start, patchs[high-1].End,
					patchs[high].File, patchs[high].Start, patchs[high].End,
				)
			}
			high++
		}

		AddSourcesFile(sys, patchs[low:high], contextSize)
		low = high
	}

	return patchs
}

// AddSourcesFile reads one file and dispatches each line to every patch that wants it.
// All patchs must belong to the same file, be sorted and be disjoint.
func AddSourcesFile(sys fs.FS, patchs []shared.Patch, contextSize int) {
	shared.PreCondition(len(patchs) > 0, "cannot process empty list")

	file, err := sys.Open(patchs[0].FileAbs)
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
func setSourceError(patchs []shared.Patch, err error) {
	for i := range patchs {
		patchs[i].SourceError = err
	}
}
