package parsing

import (
	"testing"
	"testing/fstest"
	"uncovered/pkg/models"

	"github.com/stretchr/testify/assert"
)

var mockFS = fstest.MapFS{
	"fixture.go": &fstest.MapFile{
		Data: []byte(`1 blablabla
2 truc
3 bidule
4 machin
5 foo
6 bar`),
	},
}

func TestParseSources(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		patch        models.Patch
		contextSize  int
		SourceBefore []string
		Source       []string
		SourceAfter  []string
	}{
		{
			name: "basic",
			patch: models.Patch{
				File:  "fixture.go",
				Start: 3,
				End:   3,
			},
			contextSize:  1,
			SourceBefore: []string{"2 truc"},
			Source:       []string{"3 bidule"},
			SourceAfter:  []string{"4 machin"},
		},
		{
			name: "no-context",
			patch: models.Patch{
				File:  "fixture.go",
				Start: 3,
				End:   3,
			},
			contextSize:  0,
			SourceBefore: nil,
			Source:       []string{"3 bidule"},
			SourceAfter:  nil,
		},
		{
			name: "multi-line-patch",
			patch: models.Patch{
				File:  "fixture.go",
				Start: 2,
				End:   3,
			},
			contextSize:  0,
			SourceBefore: nil,
			Source:       []string{"2 truc", "3 bidule"},
			SourceAfter:  nil,
		},
		{
			name: "first-line",
			patch: models.Patch{
				File:  "fixture.go",
				Start: 1,
				End:   1,
			},
			contextSize:  3,
			SourceBefore: nil,
			Source:       []string{"1 blablabla"},
			SourceAfter:  []string{"2 truc", "3 bidule", "4 machin"},
		},
		{
			name: "last-line",
			patch: models.Patch{
				File:  "fixture.go",
				Start: 6,
				End:   6,
			},
			contextSize:  3,
			SourceBefore: []string{"3 bidule", "4 machin", "5 foo"},
			Source:       []string{"6 bar"},
			SourceAfter:  nil,
		},
		{
			name: "context-bigger-than-file",
			patch: models.Patch{
				File:  "fixture.go",
				Start: 3,
				End:   3,
			},
			contextSize:  100,
			SourceBefore: []string{"1 blablabla", "2 truc"},
			Source:       []string{"3 bidule"},
			SourceAfter:  []string{"4 machin", "5 foo", "6 bar"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := ParseSources(mockFS, []models.Patch{tc.patch}, tc.contextSize)
			assert.Equal(t, len(results), 1)
			assert.Equal(t, results[0].File, tc.patch.File)
			assert.Equal(t, results[0].Start, tc.patch.Start)
			assert.Equal(t, results[0].End, tc.patch.End)
			assert.Equal(t, results[0].SourceBefore, tc.SourceBefore)
			assert.Equal(t, results[0].Source, tc.Source)
			assert.Equal(t, results[0].SourceAfter, tc.SourceAfter)
			assert.Equal(t, results[0].SourceError, nil)
		})
	}
}

func TestParseSourcesContextOverlap(t *testing.T) {

	patchs := []models.Patch{
		{
			File:  "fixture.go",
			Start: 3,
			End:   3,
		},
		{
			File:  "fixture.go",
			Start: 5,
			End:   5,
		},
	}

	results := ParseSources(mockFS, patchs, 3)

	assert.Len(t, results, 2)

	assert.Equal(t, results[0], models.Patch{
		File:         "fixture.go",
		Start:        3,
		End:          3,
		SourceBefore: []string{"1 blablabla", "2 truc"},
		Source:       []string{"3 bidule"},
		SourceAfter:  []string{"4 machin", "5 foo", "6 bar"},
		SourceError:  nil,
	})

	assert.Equal(t, results[1], models.Patch{
		File:         "fixture.go",
		Start:        5,
		End:          5,
		SourceBefore: []string{"2 truc", "3 bidule", "4 machin"},
		Source:       []string{"5 foo"},
		SourceAfter:  []string{"6 bar"},
		SourceError:  nil,
	})
}

func TestParseSourcesPanicOnColidingPatchs(t *testing.T) {
	// The following patch share the line 4. They should have been merge instead.
	// We should detect the malformation and panic
	patchs := []models.Patch{
		{
			File:  "fixture.go",
			Start: 3,
			End:   4,
		}, {
			File:  "fixture.go",
			Start: 4,
			End:   5,
		},
	}

	assert.Panics(t, func() {
		ParseSources(mockFS, patchs, 1)
	})
}

func TestParseSourcesPanicOnEmptyPatchs(t *testing.T) {
	// Start is greater than End
	// We should detect the malformation and panic
	patchs := []models.Patch{
		{
			File:  "fixture.go",
			Start: 5,
			End:   4,
		},
	}

	assert.Panics(t, func() {
		ParseSources(mockFS, patchs, 1)
	})
}
