package parsing

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"uncovered/pkg/models"
)

func Parse(report io.Reader, base string, contextSize int, filterFile func(string) bool) ([]models.Patch, error) {
	patchs, err := ParseCobertura(report, filterFile)
	if err != nil {
		return nil, err
	}

	patchs = SortPatch(patchs)
	patchs = ComputePaths(patchs, base)

	fileSystem := os.DirFS("/")
	patchs = ParseSources(fileSystem, patchs, contextSize)

	return patchs, nil
}

// We fs.FS so all path must be relative to /
// relative pass are interpelated to a "base" wich can be the cwd or a user inputed path
func ComputePaths(patchs []models.Patch, base string) []models.Patch {
	models.PreCondition(filepath.IsAbs(base), "base should be absolute")
	for i := range patchs {
		patchs[i].FileAbs = strings.TrimPrefix(toAbsolute(base, patchs[i].File), "/")
	}
	return patchs
}

func toAbsolute(base, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(base, path) // Join also Cleans
}
