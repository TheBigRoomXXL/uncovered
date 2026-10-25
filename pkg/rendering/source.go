package rendering

import (
	"bufio"
	"io/fs"
	"iter"
	"uncovered/pkg/shared"
)

// WithSources yields every patch with its source and surrounding context set.
// Patchs must be sorted by file then start (see parsing.SortPatch).
// TODO: addd suport for filterPatchs
func WithSources(sys fs.FS, patchs []shared.Patch, contextSize int) iter.Seq[shared.Patch] {
	shared.PreCondition(contextSize >= 0, "contex size must be positive")

	return func(yield func(shared.Patch) bool) {
		for low := 0; low < len(patchs); {
			if low > 0 {
				shared.PreConditionf(
					patchs[low-1].File < patchs[low].File,
					"patchs are not sorted by file: %s:%d-%d before %s:%d-%d",
					patchs[low-1].File, patchs[low-1].Start, patchs[low-1].End,
					patchs[low].File, patchs[low].Start, patchs[low].End,
				)
			}

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
						patchs[high-1].Start < patchs[high].Start,
						"patchs are not sorted by start: %s:%d-%d before %s:%d-%d",
						patchs[high-1].File, patchs[high-1].Start, patchs[high-1].End,
						patchs[high].File, patchs[high].Start, patchs[high].End,
					)
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

			if !withSourcesFile(sys, patchs[low:high], contextSize, yield) {
				return // consumer stopped iterating
			}
			low = high
		}
	}
}

// withSourcesFile reads one file and dispatches each line to every patch that wants it.
// Each patch is yielded as soon as it is complete.
// All patchs must belong to the same file, be sorted and be disjoint.
// It returns false if the consumer stopped iterating.
func withSourcesFile(sys fs.FS, patchs []shared.Patch, contextSize int, yield func(shared.Patch) bool) bool {
	shared.PreCondition(len(patchs) > 0, "cannot process empty list")

	file, err := sys.Open(patchs[0].FileAbs)
	if err != nil {
		setSourceError(patchs, err)
		return yieldAll(patchs, yield)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	head := 0 // only patch after head still expects lines

	for head < len(patchs) && scanner.Scan() {
		lineNumber += 1

		// Once the head patch has all of it's line, it is done: we yield it
		// and incread the head to avoid re-processing done patchs.
		for head < len(patchs) && lineNumber > patchs[head].End+contextSize {
			if !yield(patchs[head]) {
				return false
			}
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
		setSourceError(patchs[head:], err)
	}

	// Patchs still pending at EOF (context cut short by the end of file)
	return yieldAll(patchs[head:], yield)
}

// yieldAll yields every patch, it returns false if the consumer stopped iterating
func yieldAll(patchs []shared.Patch, yield func(shared.Patch) bool) bool {
	for _, patch := range patchs {
		if !yield(patch) {
			return false
		}
	}
	return true
}

// setSourceError report a file level failure on every patch of that file
func setSourceError(patchs []shared.Patch, err error) {
	for i := range patchs {
		patchs[i].SourceError = err
	}
}
