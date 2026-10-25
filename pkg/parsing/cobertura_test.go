package parsing

import (
	"iter"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"uncovered/pkg/shared"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

const DIR_REPORTS_COBERTURA = "./test-data/cobertura"

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
			coverage, err := ParseCobertura(file, "/base", FilterNothing)

			assert.Nil(t, err)
			snaps.MatchStandaloneJSON(t, Transform(coverage))
		})
	}
}

func TestParseCoberturaMergeCollidingFiles(t *testing.T) {
	t.Parallel()
	// Both classes point to the same file: a line is covered if any class hits it
	report := strings.NewReader(`<coverage><packages><package><classes>
		<class filename="a.go"><lines>
			<line number="1" hits="0"/>
			<line number="2" hits="1"/>
		</lines></class>
		<class filename="./a.go"><lines>
			<line number="1" hits="1"/>
			<line number="2" hits="0"/>
			<line number="3" hits="0"/>
		</lines></class>
	</classes></package></packages></coverage>`)

	coverage, err := ParseCobertura(report, "/base", FilterNothing)

	assert.Nil(t, err)
	assert.Equal(t, shared.Coverage{
		"base/a.go": {File: "a.go", Lines: map[int]bool{1: true, 2: true, 3: false}},
	}, coverage)
}

func TestParseCoberturaFilterOnReportPath(t *testing.T) {
	t.Parallel()
	report := strings.NewReader(`<coverage><packages><package><classes>
		<class filename="keep.go"><lines><line number="1" hits="0"/></lines></class>
		<class filename="drop.go"><lines><line number="1" hits="0"/></lines></class>
	</classes></package></packages></coverage>`)

	coverage, err := ParseCobertura(report, "/base", func(path string) bool { return path == "keep.go" })

	assert.Nil(t, err)
	assert.Equal(t, shared.Coverage{
		"base/keep.go": {File: "keep.go", Lines: map[int]bool{1: false}},
	}, coverage)
}
