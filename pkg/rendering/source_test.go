package rendering

import (
	"slices"
	"testing"
	"testing/fstest"
	"uncovered/pkg/shared"

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

func TestWithSources(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		patch        shared.Patch
		contextSize  int
		SourceBefore []string
		Source       []string
		SourceAfter  []string
	}{
		{
			name: "basic",
			patch: shared.Patch{
				File:    "fixture.go",
				FileAbs: "fixture.go",
				Start:   3,
				End:     3,
			},
			contextSize:  1,
			SourceBefore: []string{"2 truc"},
			Source:       []string{"3 bidule"},
			SourceAfter:  []string{"4 machin"},
		},
		{
			name: "no-context",
			patch: shared.Patch{
				File:    "fixture.go",
				FileAbs: "fixture.go",
				Start:   3,
				End:     3,
			},
			contextSize:  0,
			SourceBefore: nil,
			Source:       []string{"3 bidule"},
			SourceAfter:  nil,
		},
		{
			name: "multi-line-patch",
			patch: shared.Patch{
				File:    "fixture.go",
				FileAbs: "fixture.go",
				Start:   2,
				End:     3,
			},
			contextSize:  0,
			SourceBefore: nil,
			Source:       []string{"2 truc", "3 bidule"},
			SourceAfter:  nil,
		},
		{
			name: "first-line",
			patch: shared.Patch{
				File:    "fixture.go",
				FileAbs: "fixture.go",
				Start:   1,
				End:     1,
			},
			contextSize:  3,
			SourceBefore: nil,
			Source:       []string{"1 blablabla"},
			SourceAfter:  []string{"2 truc", "3 bidule", "4 machin"},
		},
		{
			name: "last-line",
			patch: shared.Patch{
				File:    "fixture.go",
				FileAbs: "fixture.go",
				Start:   6,
				End:     6,
			},
			contextSize:  3,
			SourceBefore: []string{"3 bidule", "4 machin", "5 foo"},
			Source:       []string{"6 bar"},
			SourceAfter:  nil,
		},
		{
			name: "context-bigger-than-file",
			patch: shared.Patch{
				File:    "fixture.go",
				FileAbs: "fixture.go",
				Start:   3,
				End:     3,
			},
			contextSize:  100,
			SourceBefore: []string{"1 blablabla", "2 truc"},
			Source:       []string{"3 bidule"},
			SourceAfter:  []string{"4 machin", "5 foo", "6 bar"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			results := slices.Collect(WithSources(mockFS, []shared.Patch{tc.patch}, tc.contextSize))
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

func TestWithSourcesContextOverlap(t *testing.T) {
	t.Parallel()

	patchs := []shared.Patch{
		{
			File:    "fixture.go",
			FileAbs: "fixture.go",
			Start:   3,
			End:     3,
		},
		{
			File:    "fixture.go",
			FileAbs: "fixture.go",
			Start:   5,
			End:     5,
		},
	}

	results := slices.Collect(WithSources(mockFS, patchs, 3))

	assert.Len(t, results, 2)

	assert.Equal(t, results[0], shared.Patch{
		File:         "fixture.go",
		FileAbs:      "fixture.go",
		Start:        3,
		End:          3,
		SourceBefore: []string{"1 blablabla", "2 truc"},
		Source:       []string{"3 bidule"},
		SourceAfter:  []string{"4 machin", "5 foo", "6 bar"},
		SourceError:  nil,
	})

	assert.Equal(t, results[1], shared.Patch{
		File:         "fixture.go",
		FileAbs:      "fixture.go",
		Start:        5,
		End:          5,
		SourceBefore: []string{"2 truc", "3 bidule", "4 machin"},
		Source:       []string{"5 foo"},
		SourceAfter:  []string{"6 bar"},
		SourceError:  nil,
	})
}

func TestWithSourcesPanicOnColidingPatchs(t *testing.T) {
	t.Parallel()
	// The following patch share the line 4. They should have been merge instead.
	// We should detect the malformation and panic
	patchs := []shared.Patch{
		{
			File:    "fixture.go",
			FileAbs: "fixture.go",
			Start:   3,
			End:     4,
		}, {
			File:    "fixture.go",
			FileAbs: "fixture.go",
			Start:   4,
			End:     5,
		},
	}

	shared.AssertViolationKind(t, shared.KindPreCondition, func() {
		for range WithSources(mockFS, patchs, 1) {
		}
	})
}

func TestWithSourcesPanicOnAdjacentPatchs(t *testing.T) {
	t.Parallel()
	// The following patch share the line 4. They should have been merge instead.
	// We should detect the malformation and panic
	patchs := []shared.Patch{
		{
			File:    "fixture.go",
			FileAbs: "fixture.go",
			Start:   3,
			End:     3,
		}, {
			File:    "fixture.go",
			FileAbs: "fixture.go",
			Start:   4,
			End:     4,
		},
	}

	shared.AssertViolationKind(t, shared.KindPreCondition, func() {
		for range WithSources(mockFS, patchs, 1) {
		}
	})
}

func TestWithSourcesPanicOnEmptyPatchs(t *testing.T) {
	t.Parallel()
	// Start is greater than End
	// We should detect the malformation and panic
	patchs := []shared.Patch{
		{
			File:    "fixture.go",
			FileAbs: "fixture.go",
			Start:   5,
			End:     4,
		},
	}

	shared.AssertViolationKind(t, shared.KindPreCondition, func() {
		for range WithSources(mockFS, patchs, 1) {
		}
	})
}

func TestWithSourcesPanicOnNegativeContextSize(t *testing.T) {
	t.Parallel()
	// Start is greater than End
	// We should detect the malformation and panic
	patchs := []shared.Patch{
		{
			File:    "fixture.go",
			FileAbs: "fixture.go",
			Start:   5,
			End:     4,
		},
	}

	shared.AssertViolationKind(t, shared.KindPreCondition, func() {
		for range WithSources(mockFS, patchs, -1) {
		}
	})
}

func TestWithSourcesEarlyBreak(t *testing.T) {
	t.Parallel()

	patchs := []shared.Patch{
		{File: "fixture.go", FileAbs: "fixture.go", Start: 1, End: 1},
		{File: "fixture.go", FileAbs: "fixture.go", Start: 3, End: 3},
		{File: "missing.go", FileAbs: "missing.go", Start: 1, End: 1},
	}

	count := 0
	for patch := range WithSources(mockFS, patchs, 0) {
		count++
		assert.Equal(t, patch.Source, []string{"1 blablabla"})
		break
	}
	assert.Equal(t, count, 1)
}

func TestWithSourcesPanicOnUnsortedStart(t *testing.T) {
	t.Parallel()
	// Patchs of the same file must be sorted by start
	patchs := []shared.Patch{
		{File: "fixture.go", FileAbs: "fixture.go", Start: 5, End: 5},
		{File: "fixture.go", FileAbs: "fixture.go", Start: 1, End: 1},
	}

	shared.AssertViolationKind(t, shared.KindPreCondition, func() {
		for range WithSources(mockFS, patchs, 1) {
		}
	})
}

func TestWithSourcesPanicOnUnsortedFile(t *testing.T) {
	t.Parallel()
	// Patchs must be sorted by file, a file cannot come back once its run is over
	patchs := []shared.Patch{
		{File: "fixture.go", FileAbs: "fixture.go", Start: 1, End: 1},
		{File: "other.go", FileAbs: "other.go", Start: 1, End: 1},
		{File: "fixture.go", FileAbs: "fixture.go", Start: 3, End: 3},
	}

	shared.AssertViolationKind(t, shared.KindPreCondition, func() {
		for range WithSources(mockFS, patchs, 1) {
		}
	})
}
