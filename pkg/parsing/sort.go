package parsing

import (
	"sort"
	"uncovered/pkg/models"
)

// slice of patch are sorted by filename and start so that they can be prosseced continuously
func SortPatch(patchs []models.Patch) []models.Patch {
	sort.Slice(patchs, func(i, j int) bool {
		switch {
		case patchs[i].File != patchs[j].File:
			return patchs[i].File < patchs[j].File
		default:
			return patchs[i].Start < patchs[j].Start
		}
	})
	return patchs
}
