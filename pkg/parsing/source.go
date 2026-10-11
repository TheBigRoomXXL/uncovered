package parsing

import (
	"bufio"
	"io"
	"io/fs"
	"uncovered/pkg/models"
)

func ParseSources(sys fs.FS, patchs []models.Patch, contextSize int) []models.Patch {
	// Sorting ensure we process patchs file by file and only increase line number
	SortPatch(patchs)

	var filepath string
	var file io.ReadCloser
	var scanner *bufio.Scanner
	var lineNumber int
	var err error
	var goToNextFile bool
	for i := range len(patchs) {
		// We mutate the array directly, no copy
		patch := &patchs[i]
		models.PreConditionf(patch.Start <= patch.End, "patch has an empty range: %s:%d-%d", patch.File, patch.Start, patch.End)

		// In case of error on the file, we skip every patch in that file
		if goToNextFile && filepath == patch.File {
			patch.SourceError = err
			continue
		}

		// Open the file it's a new one
		if filepath != patch.File {
			filepath = patch.File
			file, err = sys.Open(filepath)

			if err != nil {
				patch.SourceError = err
				goToNextFile = true
				continue
			}

			defer file.Close()
			scanner = bufio.NewScanner(file)
			lineNumber = 0
			goToNextFile = false
		}

		models.Assertf(lineNumber < patch.Start, "line number greater than current patch start: %d > %d", lineNumber, patch.Start)

		// Iterate over lines until we reach the part that interest us
		for scanner.Scan() {
			lineNumber += 1

			if lineNumber < patch.Start-contextSize {
				continue
			} else if lineNumber < patch.Start {
				patch.SourceBefore = append(patch.SourceBefore, scanner.Text())
			} else if lineNumber < patch.End {
				patch.Source = append(patch.Source, scanner.Text())
			} else if lineNumber == patch.End {
				patch.Source = append(patch.Source, scanner.Text())
				if contextSize == 0 {
					break // Process next patch
				}
			} else if lineNumber < patch.End+contextSize {
				patch.SourceAfter = append(patch.SourceAfter, scanner.Text())
			} else if lineNumber == patch.End+contextSize {
				patch.SourceAfter = append(patch.SourceAfter, scanner.Text())
				break // Process next patch
			} else {
				models.Unreachable("Current line should not be greater than patch.End+ContextSize")
			}
		}

		err = scanner.Err()
		if err != nil {
			patch.SourceError = err
			goToNextFile = true
		}
	}

	return patchs
}
