package parsing

import (
	"testing"
	"uncovered/pkg/models"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

func AssertPatchsOrder(t models.T, patchs []models.Patch) {
	t.Helper()

	if len(patchs) <= 1 {
		return
	}

	for i := 1; i < len(patchs); i++ {
		if patchs[i].File < patchs[i-1].File {
			t.Fatalf(
				"patchs are not sorted correct: filepath %s (index %d) < %s (index %d)",
				patchs[i].File, i, patchs[i-1].File, i-1,
			)
		}

		if patchs[i].File == patchs[i-1].File && patchs[i].Start < patchs[i-1].Start {
			t.Fatalf(
				"patchs are not sorted correct: start %d (index %d) < %d (index %d)",
				patchs[i].Start, i, patchs[i-1].Start, i-1,
			)
		}
	}
}

func Test_SortPatchs(t *testing.T) {
	tests := []struct {
		name   string
		patchs []models.Patch
	}{
		{
			name:   "0 patch",
			patchs: []models.Patch{},
		},
		{
			name:   "1 patch",
			patchs: []models.Patch{{File: "truc.go", Start: 10, End: 12}},
		},
		{
			name:   "unsorted by name",
			patchs: []models.Patch{{File: "truc.go", Start: 10, End: 12}, {File: "bidule.go", Start: 10, End: 12}},
		},
		{
			name:   "unsorted by start",
			patchs: []models.Patch{{File: "truc.go", Start: 10, End: 12}, {File: "truc.go", Start: 7, End: 8}},
		},
		{
			name: "already sorted list",
			patchs: []models.Patch{
				{File: "a.go", Start: 1, End: 5},
				{File: "a.go", Start: 10, End: 15},
				{File: "b.go", Start: 2, End: 4},
			},
		},
		{
			name: "multiple files with mixed start positions",
			patchs: []models.Patch{
				{File: "c.go", Start: 5, End: 10},
				{File: "a.go", Start: 20, End: 25},
				{File: "a.go", Start: 5, End: 10},
				{File: "b.go", Start: 1, End: 3},
			},
		},
		{
			name: "same start lines across different files",
			patchs: []models.Patch{
				{File: "zebra.go", Start: 10, End: 20},
				{File: "alpha.go", Start: 10, End: 15},
			},
		},
		{
			name: "empty string filenames",
			patchs: []models.Patch{
				{File: "b.go", Start: 10, End: 12},
				{File: "", Start: 5, End: 8},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SortPatch(tc.patchs)

			assert.Len(t, result, len(tc.patchs))
			AssertPatchsOrder(t, result)
		})
	}
}

func Test_SortPatchs_PBT(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		patchs := rapid.SliceOf(models.PatchGenerator).Draw(t, "patchs")

		result := SortPatch(patchs)
		assert.Len(t, result, len(patchs))
		AssertPatchsOrder(t, result)
	})
}
