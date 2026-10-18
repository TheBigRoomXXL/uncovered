package rendering

import (
	"testing"
	"uncovered/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestRenderPatchBasic(t *testing.T) {
	p := models.Patch{
		File:         "something.go",
		Start:        10,
		End:          10,
		SourceBefore: []string{"bidule", "machine"},
		Source:       []string{"uncovered"},
		SourceAfter:  []string{"machin", "truc"},
	}
	got := RenderPatchText(p, false)

	expected := `something.go:10
 8 |   bidule
 9 |   machine
10 | - uncovered
11 |   machin
12 |   truc
`
	assert.Equal(t, expected, got)
}

func TestRenderPatchWithColor(t *testing.T) {
	p := models.Patch{
		File:         "something.go",
		Start:        10,
		End:          10,
		SourceBefore: []string{"bidule", "machine"},
		Source:       []string{"uncovered"},
		SourceAfter:  []string{"machin", "truc"},
	}
	got := RenderPatchText(p, true)

	c := ANSIColors
	expected := "something.go:10\n" +
		c.Dim + " 8 |   bidule" + c.Reset + "\n" +
		c.Dim + " 9 |   machine" + c.Reset + "\n" +
		c.Orange + "10 | - uncovered" + c.Reset + "\n" +
		c.Dim + "11 |   machin" + c.Reset + "\n" +
		c.Dim + "12 |   truc" + c.Reset + "\n"

	assert.Equal(t, expected, got)
}
