package parsing

import (
	"testing"
	"uncovered/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestComputePaths(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		file     string
		expected string
	}{
		{"relative path joined to base", "/home/user/project", "src/main.go", "home/user/project/src/main.go"},
		{"absolute path ignores base", "/home/user/project", "/etc/app/main.go", "etc/app/main.go"},
		{"dot prefix is cleaned", "/home/user/project", "./src/main.go", "home/user/project/src/main.go"},
		{"parent dir is resolved", "/home/user/project", "../other/main.go", "home/user/other/main.go"},
		{"absolute path is cleaned", "/", "/a//b/../c/main.go", "a/c/main.go"},
		{"root base", "/", "main.go", "main.go"},
		{"base with trailing slash", "/home/user/", "main.go", "home/user/main.go"},
		{"escaping above root stays at root", "/home", "../../../main.go", "main.go"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			patchs := []models.Patch{{File: tc.file}}
			result := ComputePaths(patchs, tc.base)
			assert.Equal(t, tc.file, result[0].File)
			assert.Equal(t, tc.expected, result[0].FileAbs)
		})
	}
}

func TestComputePathsMultiplePatchs(t *testing.T) {
	patchs := []models.Patch{
		{File: "a.go", Start: 1, End: 2},
		{File: "/abs/b.go", Start: 3, End: 4},
	}

	result := ComputePaths(patchs, "/base")

	assert.Equal(t, []models.Patch{
		{File: "a.go", FileAbs: "base/a.go", Start: 1, End: 2},
		{File: "/abs/b.go", FileAbs: "abs/b.go", Start: 3, End: 4},
	}, result)
}

func TestComputePathsEmpty(t *testing.T) {
	assert.Empty(t, ComputePaths([]models.Patch{}, "/base"))
	assert.Empty(t, ComputePaths(nil, "/base"))
}

func TestComputePathsViolation(t *testing.T) {
	t.Parallel()

	models.AssertViolationKind(t, models.KindPreCondition, func() {
		ComputePaths([]models.Patch{}, "./someting/relative")
	})
}
