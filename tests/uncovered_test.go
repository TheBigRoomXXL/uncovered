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
			render: rendering.NewRendererText(5, false),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			c.Ouput = &b
			c.Render = tc.render

			if err := pkg.Run(c); err != nil {
				t.Fatal(err)
			}
			snaps.MatchStandaloneSnapshot(t, b.String())
		})
	}
}

func TestRunWithConfig(t *testing.T) {
	testCases := []struct {
		name   string
		update func(c *shared.Config)
	}{
		{
			name:   "no-filters",
			update: func(c *shared.Config) {},
		},
		{
			name: "include-one-file",
			update: func(c *shared.Config) {
				c.Includes = []string{"truc.go"}
			},
		},
		{
			name: "include-directory",
			update: func(c *shared.Config) {
				c.Includes = []string{"something.go/"}
			},
		},
		{
			name: "include-several-files",
			update: func(c *shared.Config) {
				c.Includes = []string{"truc.go", "this.go"}
			},
		},
		{
			name: "exclude-one-file",
			update: func(c *shared.Config) {
				c.Excludes = []string{"machin.go"}
			},
		},
		{
			name: "exclude-several-files",
			update: func(c *shared.Config) {
				c.Excludes = []string{"machin.go", "that.go"}
			},
		},
		{
			name: "include-directory-exclude-file",
			update: func(c *shared.Config) {
				c.Includes = []string{"something.go/"}
				c.Excludes = []string{"that.go"}
			},
		},
		{
			name: "exclude-everything",
			update: func(c *shared.Config) {
				c.Excludes = []string{".go"}
			},
		},
		{
			name: "include-nothing-matching",
			update: func(c *shared.Config) {
				c.Includes = []string{"does-not-exist.go"}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			c := shared.NewConfig()
			c.Base = "./fakerepos"
			c.ReportPath = "./fakerepos/coverage.xml"
			c.Ouput = &b
			c.Render = rendering.NewRendererMarkdown(1)
			tc.update(c)

			if err := pkg.Run(c); err != nil {
				t.Fatal(err)
			}
			snaps.MatchStandaloneSnapshot(t, b.String())
		})
	}
}

func TestRunErrors(t *testing.T) {
	testCases := []struct {
		name       string
		reportPath string
	}{
		{name: "missing-report", reportPath: "./fakerepos/does-not-exist.xml"},
		{name: "not-a-report", reportPath: "./fakerepos/go.mod"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			c := shared.NewConfig()
			c.Base = "./fakerepos"
			c.ReportPath = tc.reportPath
			c.Ouput = &b
			c.Render = rendering.NewRendererMarkdown(1)

			if err := pkg.Run(c); err == nil {
				t.Fatalf("expected an error, got output:\n%s", b.String())
			}
		})
	}
}
