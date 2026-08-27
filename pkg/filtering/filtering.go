package filtering

import (
	"strings"
)

func FilterForFiles(includes, excludes []string) func(filepath string) bool {
	return func(filepath string) bool {
		if len(includes) > 0 && !matchOneOf(filepath, includes) {
			return false
		}
		if len(excludes) > 0 && matchOneOf(filepath, excludes) {
			return false
		}
		return true
	}
}

func matchOneOf(needle string, haystack []string) bool {
	for _, substr := range haystack {
		if strings.Contains(needle, substr) {
			return true
		}
	}
	return false
}
