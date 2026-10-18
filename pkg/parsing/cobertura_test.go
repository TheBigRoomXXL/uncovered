package parsing

import (
	"iter"
	"os"
	"path/filepath"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

const DIR_TEST_DATA = "../../testdata/"
const DIR_REPORTS_COBERTURA = DIR_TEST_DATA + "reports/cobertura"

// ForFilesIn yields the filename and an opened *os.File for each non-directory entry in dirPath.
// The yielded file is closed when the test (including parallel subtests) completes.
// Does not recurse
func ForFilesIn(t *testing.T, dirPath string) iter.Seq2[string, *os.File] {
	t.Helper()

	return func(yield func(string, *os.File) bool) {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			path := filepath.Join(dirPath, entry.Name())

			file, err := os.Open(path)
			if err != nil {
				continue
			}
			t.Cleanup(func() { file.Close() })

			ok := yield(entry.Name(), file)
			if !ok {
				break
			}
		}
	}
}

func FilterNothing(path string) bool {
	return true
}

func Test_Cobertura_Parsing(t *testing.T) {
	t.Parallel()
	for filename, file := range ForFilesIn(t, DIR_REPORTS_COBERTURA) {
		t.Run(filename, func(t *testing.T) {
			t.Parallel()
			result, err := ParseCobertura(file, FilterNothing)

			assert.Nil(t, err)
			snaps.MatchStandaloneJSON(t, result)
		})
	}
}
