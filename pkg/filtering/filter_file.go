package filtering

import "slices"

func FilterForFiles(includes, excludes []string) func(filepath string) bool {
	return func(filepath string) bool {
		if len(includes) > 0 && !slices.Contains(includes, filepath) {
			return false
		}
		if len(excludes) > 0 && slices.Contains(excludes, filepath) {
			return false
		}
		return true
	}
}
