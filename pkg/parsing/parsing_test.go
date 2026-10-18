package parsing

import (
	"testing"
	"uncovered/pkg/shared"

	"github.com/stretchr/testify/assert"
)

func TestComputePath(t *testing.T) {
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
			assert.Equal(t, tc.expected, ComputePath(tc.base, tc.file))
		})
	}
}

func TestComputePathViolation(t *testing.T) {
	t.Parallel()

	shared.AssertViolationKind(t, shared.KindPreCondition, func() {
		ComputePath("./someting/relative", "main.go")
	})
}

func TestTransform(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		coverage shared.Coverage
		expected []shared.Patch
	}{
		{"nil coverage", nil, []shared.Patch{}},
		{"empty coverage", shared.Coverage{}, []shared.Patch{}},
		{
			"fully covered file",
			shared.Coverage{"a.go": {File: "a.go", Lines: map[int]bool{1: true, 2: true}}},
			[]shared.Patch{},
		},
		{
			"consecutive uncovered lines are grouped",
			shared.Coverage{"a.go": {File: "a.go", Lines: map[int]bool{1: false, 2: false, 3: false}}},
			[]shared.Patch{{File: "a.go", FileAbs: "a.go", Start: 1, End: 3}},
		},
		{
			"missing lines do not break a patch",
			shared.Coverage{"a.go": {File: "a.go", Lines: map[int]bool{2: false, 5: false, 9: false}}},
			[]shared.Patch{{File: "a.go", FileAbs: "a.go", Start: 2, End: 9}},
		},
		{
			"covered line breaks a patch",
			shared.Coverage{"a.go": {File: "a.go", Lines: map[int]bool{1: false, 2: false, 3: true, 4: false}}},
			[]shared.Patch{
				{File: "a.go", FileAbs: "a.go", Start: 1, End: 2},
				{File: "a.go", FileAbs: "a.go", Start: 4, End: 4},
			},
		},
		{
			"display path and absolute path are kept",
			shared.Coverage{"base/a.go": {File: "./a.go", Lines: map[int]bool{7: false}}},
			[]shared.Patch{{File: "./a.go", FileAbs: "base/a.go", Start: 7, End: 7}},
		},
		{
			"multiple files are sorted",
			shared.Coverage{
				"b.go": {File: "b.go", Lines: map[int]bool{1: false}},
				"a.go": {File: "a.go", Lines: map[int]bool{10: false, 11: true, 12: false}},
			},
			[]shared.Patch{
				{File: "a.go", FileAbs: "a.go", Start: 10, End: 10},
				{File: "a.go", FileAbs: "a.go", Start: 12, End: 12},
				{File: "b.go", FileAbs: "b.go", Start: 1, End: 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, Transform(tc.coverage))
		})
	}
}
