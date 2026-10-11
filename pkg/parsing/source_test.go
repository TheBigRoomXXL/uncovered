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

func TestParseSourcesBasic(t *testing.T) {

	patchs := []models.Patch{{
		File:  "fixture.go",
		Start: 3,
		End:   3,
	}}

	results := ParseSources(mockFS, patchs, 1)

	assert.Len(t, results, 1)

	assert.Equal(t, results[0], models.Patch{
		File:         "fixture.go",
		Start:        3,
		End:          3,
		SourceBefore: []string{"2 truc"},
		Source:       []string{"3 bidule"},
		SourceAfter:  []string{"4 machin"},
		SourceError:  nil,
	})
}

func TestParseSourcesNoContext(t *testing.T) {

	patchs := []models.Patch{{
		File:  "fixture.go",
		Start: 3,
		End:   3,
	}}

	results := ParseSources(mockFS, patchs, 0)

	assert.Len(t, results, 1)

	assert.Equal(t, results[0], models.Patch{
		File:         "fixture.go",
		Start:        3,
		End:          3,
		SourceBefore: nil,
		Source:       []string{"3 bidule"},
		SourceAfter:  nil,
		SourceError:  nil,
	})
}

func TestParseSourcesLine1(t *testing.T) {

	patchs := []models.Patch{{
		File:  "fixture.go",
		Start: 1,
		End:   1,
	}}

	results := ParseSources(mockFS, patchs, 3)

	assert.Len(t, results, 1)

	assert.Equal(t, results[0], models.Patch{
		File:         "fixture.go",
		Start:        1,
		End:          1,
		SourceBefore: nil,
		Source:       []string{"1 blablabla"},
		SourceAfter:  []string{"2 truc", "3 bidule", "4 machin"},
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
