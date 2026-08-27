package filtering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FilterForFiles(t *testing.T) {
	var tests = []struct {
		name, filepath     string
		includes, excludes []string
		expected           bool
	}{
		{
			name:     "no-filters",
			filepath: "something/truc.go",
			includes: []string{},
			excludes: []string{},
			expected: true,
		},
		{
			name:     "include-the-filename",
			filepath: "something/truc.go",
			includes: []string{"truc.go"},
			excludes: []string{},
			expected: true,
		},
		{
			name:     "include-the-directory",
			filepath: "something/truc.go",
			includes: []string{"something"},
			excludes: []string{},
			expected: true,
		},
		{
			name:     "include-the-directory-with-slash",
			filepath: "something/truc.go",
			includes: []string{"something/"},
			excludes: []string{},
			expected: true,
		},
		{
			name:     "include-full-path",
			filepath: "something/truc.go",
			includes: []string{"something/truc.go"},
			excludes: []string{},
			expected: true,
		},
		{
			name:     "include-other-file",
			filepath: "something/truc.go",
			includes: []string{"bidule.go"},
			excludes: []string{},
			expected: false,
		},
		{
			name:     "include-other-directory",
			filepath: "something/truc.go",
			includes: []string{"machin/"},
			excludes: []string{},
			expected: false,
		},
		{
			name:     "include-other-full-path",
			filepath: "something/truc.go",
			includes: []string{"machin/bidule.go"},
			excludes: []string{},
			expected: false,
		},
		{
			name:     "exclude-the-filename",
			filepath: "something/truc.go",
			includes: []string{},
			excludes: []string{"truc.go"},
			expected: false,
		},
		{
			name:     "exclude-the-directory",
			filepath: "something/truc.go",
			includes: []string{},
			excludes: []string{"something"},
			expected: false,
		},
		{
			name:     "exclude-the-directory-with-slash",
			filepath: "something/truc.go",
			includes: []string{},
			excludes: []string{"something/"},
			expected: false,
		},
		{
			name:     "exclude-full-path",
			filepath: "something/truc.go",
			includes: []string{},
			excludes: []string{"something/truc.go"},
			expected: false,
		},
		{
			name:     "exclude-other-full-path",
			filepath: "something/truc.go",
			includes: []string{},
			excludes: []string{"machin/bidule.go"},
			expected: true,
		},
		{
			name:     "exclude-other-filename",
			filepath: "something/truc.go",
			includes: []string{},
			excludes: []string{"bidule.go"},
			expected: true,
		},
		{
			name:     "exclude-other-directory",
			filepath: "something/truc.go",
			includes: []string{},
			excludes: []string{"machin/"},
			expected: true,
		},
		{
			name:     "include-exclude-filename",
			filepath: "something/truc.go",
			includes: []string{"truc.go"},
			excludes: []string{"truc.go"},
			expected: false,
		},
		{
			name:     "include-exclude-directory",
			filepath: "something/truc.go",
			includes: []string{"/something"},
			excludes: []string{"/something"},
			expected: false,
		},
		{
			name:     "include-exclude-full-path",
			filepath: "something/truc.go",
			includes: []string{"/something/truc.go"},
			excludes: []string{"/something/truc.go"},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filefunc := FilterForFiles(tc.includes, tc.excludes)

			result := filefunc(tc.filepath)
			assert.Equal(t, tc.expected, result)
		})
	}
}
