package rendering

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"uncovered/pkg/shared"
)

func NewRendererMarkdown(contextSize int) shared.Renderer {
	// FS is always the full file system, we only use it to facilitate tests of
	// WithSources using an interface
	fileSystem := os.DirFS("/")

	return func(stream io.Writer, patchs []shared.Patch) {
		streamBuffered := bufio.NewWriter(stream)
		for patch := range WithSources(fileSystem, patchs, contextSize) {
			RenderPatchMarkdown(streamBuffered, patch)
			io.WriteString(stream, "\n")
			streamBuffered.Flush()
		}
	}
}

// Display a patch as markdown
func RenderPatchMarkdown(stream io.Writer, patch shared.Patch) {
	io.WriteString(stream, patch.Identifier())

	if patch.SourceError != nil {
		fmt.Fprintf(stream, "error retriving source: %s\n", patch.SourceError)
		return
	}

	io.WriteString(stream, "\n```diff\n")
	for _, lineBefore := range patch.SourceBefore {
		fmt.Fprintf(stream, "  %s\n", lineBefore)
	}

	for _, line := range patch.Source {
		fmt.Fprintf(stream, "- %s\n", line)
	}

	for _, lineAfter := range patch.SourceAfter {
		fmt.Fprintf(stream, "  %s\n", lineAfter)
	}
	io.WriteString(stream, "```\n")
}
