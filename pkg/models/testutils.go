package models

import (
	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

type T interface {
	// Cleanup(func())
	// Error(args ...interface{})
	Errorf(format string, args ...interface{})
	// Fail()
	// FailNow()
	// Failed() bool
	// Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	// Log(args ...interface{})
	// Logf(format string, args ...interface{})
	// Name() string
	// Parallel()
	// Setenv(key string, value string)
	// Skip(args ...interface{})
	// SkipNow()
	// Skipf(format string, args ...interface{})
	// Skipped() bool
	// TempDir() string
}

// Generate Patch struct using rapid for PBT tests
var PatchGenerator = rapid.Custom(func(t *rapid.T) Patch {
	start := rapid.IntMin(0).Draw(t, "start")
	end := rapid.IntMin(start).Draw(t, "end")
	return Patch{
		File:  rapid.String().Draw(t, "file"),
		Start: start,
		End:   end,
	}
})

func AssertViolation(t T, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		_, ok := r.(*Violation)
		assert.Truef(t, ok, "expected *models.Violation panic, got %#v", r)
	}()
	f()
}

func AssertViolationKind(t T, kind ContractKind, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		v, ok := r.(*Violation)
		if assert.Truef(t, ok, "expected *models.Violation panic, got %#v", r) {
			assert.Equal(t, kind, v.Kind, v.Msg)
		}
	}()
	f()
}
