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
