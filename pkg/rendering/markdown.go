package rendering

import (
	"fmt"
	"strings"
	"uncovered/pkg/shared"
)

// Display a patch a markdown
func RenderPatchMarkdown(patch shared.Patch) string {
	var result strings.Builder
	result.WriteString(patch.Identifier())

	if patch.SourceError != nil {
		result.WriteString(fmt.Sprintf("error retriving source: %s\n", patch.SourceError))
		return result.String()
	}

	result.WriteString("\n```diff\n")
	for _, lineBefore := range patch.SourceBefore {
		result.WriteString(fmt.Sprintf("  %s\n", lineBefore))
	}

	for _, line := range patch.Source {
		result.WriteString(fmt.Sprintf("- %s\n", line))
	}

	for _, lineAfter := range patch.SourceAfter {
		result.WriteString(fmt.Sprintf("  %s\n", lineAfter))
	}
	result.WriteString("```\n")

	return result.String()
}
