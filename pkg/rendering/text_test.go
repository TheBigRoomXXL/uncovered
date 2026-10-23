package rendering

import (
	"strconv"
	"testing"
	"uncovered/pkg/shared"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderPatchTextBasic(t *testing.T) {
	p := shared.Patch{
		File:         "something.go",
		Start:        10,
		End:          10,
		SourceBefore: []string{"bidule", "machine"},
		Source:       []string{"uncovered"},
		SourceAfter:  []string{"machin", "truc"},
	}
	got := RenderPatchText(p, false)

	expected := `something.go:10
 8 |   bidule
 9 |   machine
10 | ✗ uncovered
11 |   machin
12 |   truc
`
	assert.Equal(t, expected, got)
}

func TestRenderPatchTextWithColor(t *testing.T) {
	p := shared.Patch{
		File:         "something.go",
		Start:        10,
		End:          10,
		SourceBefore: []string{"bidule", "machine"},
		Source:       []string{"uncovered"},
		SourceAfter:  []string{"machin", "truc"},
	}
	got := RenderPatchText(p, true)

	c := ANSIColors
	expected := "something.go:10\n" +
		c.Dim + " 8 | bidule" + c.Reset + "\n" +
		c.Dim + " 9 | machine" + c.Reset + "\n" +
		c.Orange + "10 | uncovered" + c.Reset + "\n" +
		c.Dim + "11 | machin" + c.Reset + "\n" +
		c.Dim + "12 | truc" + c.Reset + "\n"

	assert.Equal(t, expected, got)
}

func TestCountDigit(t *testing.T) {

	testCases := []struct {
		number int
		expect int
	}{
		{
			number: 1,
			expect: 1,
		},
		{
			number: 10,
			expect: 2,
		},
		{
			number: 100,
			expect: 3,
		},
		{
			number: 42424324,
			expect: 8,
		},
	}

	for _, tc := range testCases {
		t.Run(strconv.Itoa(tc.number), func(t *testing.T) {
			count := countDigits(tc.number)
			require.Equal(t, tc.expect, count)
		})
	}
}

func TestCountPrecondition(t *testing.T) {

	shared.AssertViolationKind(t, shared.KindPreCondition, func() {
		countDigits(-1)
	})
}
