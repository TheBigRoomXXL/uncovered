package rendering

import (
	"testing"
	"uncovered/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestRenderPatchMarkdownBasic(t *testing.T) {
	p := models.Patch{
		File:         "something.go",
		Start:        10,
		End:          10,
		SourceBefore: []string{"bidule", "machine"},
		Source:       []string{"uncovered"},
		SourceAfter:  []string{"machin", "truc"},
	}
	got := RenderPatchMarkdown(p)

	expected := "/something.go:10\n" +
		"```diff\n" +
		"  bidule\n" +
		"  machine\n" +
		"- uncovered\n" +
		"  machin\n" +
		"  truc\n" +
		"```\n"
	assert.Equal(t, expected, got)
}
