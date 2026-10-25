package rendering

import (
	"strings"
	"testing"
	"uncovered/pkg/shared"

	"github.com/stretchr/testify/assert"
)

func TestRenderPatchMarkdownBasic(t *testing.T) {
	p := shared.Patch{
		File:         "something.go",
		Start:        10,
		End:          10,
		SourceBefore: []string{"bidule", "machine"},
		Source:       []string{"uncovered"},
		SourceAfter:  []string{"machin", "truc"},
	}

	var b strings.Builder
	RenderPatchMarkdown(&b, p)
	got := b.String()

	expected := "something.go:10\n" +
		"```diff\n" +
		"  bidule\n" +
		"  machine\n" +
		"- uncovered\n" +
		"  machin\n" +
		"  truc\n" +
		"```\n"
	assert.Equal(t, expected, got)
}
