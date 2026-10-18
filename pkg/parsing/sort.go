package parsing

import (
	"sort"
	"uncovered/pkg/shared"
)

// slice of patch are sorted by filename and start so that they can be prosseced continuously
func SortPatch(patchs []shared.Patch) []shared.Patch {
	// TODO: repalce with slices.SortStableFunc()
	sort.Slice(patchs, func(i, j int) bool {
		if patchs[i].File != patchs[j].File {
			return patchs[i].File < patchs[j].File
		} else {
			return patchs[i].Start < patchs[j].Start
		}
	})
	return patchs
}
