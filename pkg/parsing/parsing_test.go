package parsing

import (
	"testing"
	"uncovered/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestCleanPaths(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patchs := []models.Patch{{File: tt.file}}
			result := CleanPaths(patchs, tt.base)
			assert.Equal(t, tt.expected, result[0].File)
		})
	}
}

func TestCleanPathsMultiplePatchs(t *testing.T) {
	patchs := []models.Patch{
		{File: "a.go", Start: 1, End: 2},
		{File: "/abs/b.go", Start: 3, End: 4},
	}

	result := CleanPaths(patchs, "/base")

	assert.Equal(t, []models.Patch{
		{File: "base/a.go", Start: 1, End: 2},
		{File: "abs/b.go", Start: 3, End: 4},
	}, result)
}

func TestCleanPathsEmpty(t *testing.T) {
	assert.Empty(t, CleanPaths([]models.Patch{}, "/base"))
	assert.Empty(t, CleanPaths(nil, "/base"))
}

func TestCleanPathsViolation(t *testing.T) {
	t.Parallel()

	models.AssertViolationKind(t, models.KindPreCondition, func() {
		CleanPaths([]models.Patch{}, "./someting/relative")
	})
}
