package cli

import (
	"flag"
	"strings"
	"testing"
	"uncovered/pkg/shared"

	"github.com/stretchr/testify/require"
)

func TestNewConfigFromCLINoOption(t *testing.T) {
	t.Parallel()
	defaults := shared.NewConfig()
	defaults.ReportPath = "report.xml"

	cliConfig, err := NewConfigFromCLI([]string{"uncovered", "report.xml"})

	require.Nil(t, err)
	require.Equal(t, defaults.Base, cliConfig.Base)
	require.Equal(t, defaults.Includes, cliConfig.Includes)
	require.Equal(t, defaults.Excludes, cliConfig.Excludes)
	require.Equal(t, defaults.ContextSize, cliConfig.ContextSize)
	require.NotNil(t, cliConfig.Ouput)
	require.NotNil(t, cliConfig.Render)
}

func TestNewConfigFromCLIHelpOption(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"uncovered -h",
		"uncovered -help",
		"uncovered --help",
		"uncovered -h something blablabla",
		"uncovered -help something blablabla",
		"uncovered --help something blablabla",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			input := strings.Split(tc, " ")
			conf, err := NewConfigFromCLI(input)
			require.ErrorIs(t, err, flag.ErrHelp)
			require.Nil(t, conf)
		})
	}
}

func TestNewConfigOptionsIncludes(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		args   string
		expect []string
	}{
		{
			args:   "uncovered --include something report.xml",
			expect: []string{"something"},
		},
		{
			args:   "uncovered -include something report.xml",
			expect: []string{"something"},
		},
		{
			args:   "uncovered --include something --include else report.xml",
			expect: []string{"something", "else"},
		},
		{
			args:   "uncovered -include something -include else report.xml",
			expect: []string{"something", "else"},
		},
		{
			args:   "uncovered --include something -include else report.xml",
			expect: []string{"something", "else"},
		},
		{
			args:   "uncovered report.xml",
			expect: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.args, func(t *testing.T) {
			input := strings.Split(tc.args, " ")
			conf, err := NewConfigFromCLI(input)
			require.Nil(t, err)
			require.Equal(t, tc.expect, conf.Includes)
		})
	}
}

func TestNewConfigOptionsExcludes(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		args   string
		expect []string
	}{
		{
			args:   "uncovered --exclude something report.xml",
			expect: []string{"something"},
		},
		{
			args:   "uncovered -exclude something report.xml",
			expect: []string{"something"},
		},
		{
			args:   "uncovered --exclude something --exclude else report.xml",
			expect: []string{"something", "else"},
		},
		{
			args:   "uncovered -exclude something -exclude else report.xml",
			expect: []string{"something", "else"},
		},

		{
			args:   "uncovered --exclude something -exclude else report.xml",
			expect: []string{"something", "else"},
		},
		{
			args:   "uncovered report.xml",
			expect: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.args, func(t *testing.T) {
			input := strings.Split(tc.args, " ")
			conf, err := NewConfigFromCLI(input)
			require.Nil(t, err)
			require.Equal(t, tc.expect, conf.Excludes)
		})
	}
}

func TestNewConfigOptionsContextSize(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		args   string
		expect int
	}{
		{
			args:   "uncovered --context-size 9 report.xml",
			expect: 9,
		},
		{
			args:   "uncovered --context-size 1 report.xml",
			expect: 1,
		},
		{
			args:   "uncovered --context-size 256 report.xml",
			expect: 256,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.args, func(t *testing.T) {
			input := strings.Split(tc.args, " ")
			conf, err := NewConfigFromCLI(input)
			require.Nil(t, err)
			require.Equal(t, tc.expect, conf.ContextSize)
		})
	}
}

func TestNewConfigOptionsError(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		args   string
		expect string
	}{
		{
			args:   "uncovered --context-size 0 report.xml",
			expect: "context-size must be a positive integer, got 0",
		},
		{
			args:   "uncovered --context-size -1 report.xml",
			expect: "context-size must be a positive integer, got -1",
		},
		{
			args:   "uncovered --context-size 2^9999 report.xml",
			expect: "invalid value \"2^9999\" for flag -context-size: parse error",
		},
		{
			args:   "uncovered --context-size 2,53 report.xml",
			expect: "invalid value \"2,53\" for flag -context-size: parse error",
		},
		{
			args:   "uncovered --context-size 2,00 report.xml",
			expect: "invalid value \"2,00\" for flag -context-size: parse error",
		},
		{
			args:   "uncovered --exclude report.xml",
			expect: "expected exactly one <report-path> as argument, got []",
		},
		{
			args:   "uncovered --include report.xml",
			expect: "expected exactly one <report-path> as argument, got []",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.args, func(t *testing.T) {
			input := strings.Split(tc.args, " ")
			conf, err := NewConfigFromCLI(input)
			require.Equal(t, tc.expect, err.Error())
			require.Nil(t, conf)
		})
	}
}
