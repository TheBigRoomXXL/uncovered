// Intergrations test for uncovered
package tests

import (
	"strings"
	"testing"
	"uncovered/pkg"
	"uncovered/pkg/rendering"
	"uncovered/pkg/shared"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestRunWithRenderer(t *testing.T) {
	c := shared.NewConfig()
	c.Base = "./fakerepos"
	c.ReportPath = "./fakerepos/coverage.xml"

	testCases := []struct {
		name   string
		render shared.Renderer
	}{
		{
			name:   "markdown-0-context-lines",
			render: rendering.NewRendererMarkdown(0),
		},
		{
			name:   "markdown-3-context-lines",
			render: rendering.NewRendererMarkdown(3),
		},
		{
			name:   "text-3-context-lines-colorful",
			render: rendering.NewRendererText(3, true),
		},
		{
			name:   "text-5-context-lines-no-color",
			render: rendering.NewRendererText(3, false),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			c.Ouput = &b
			c.Render = tc.render

			pkg.Run(c)
			snaps.MatchStandaloneSnapshot(t, b.String())
		})
	}
}
