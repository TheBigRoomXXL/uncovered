package fakerepos

import (
	"errors"
	"strings"
)

// Truc is a small stack of strings.
type Truc struct {
	items []string
}

var ErrEmpty = errors.New("truc is empty")

// Push is fully covered.
func (t *Truc) Push(s string) {
	t.items = append(t.items, s)
}

// Pop is partially covered: the empty branch is never tested.
func (t *Truc) Pop() (string, error) {
	if len(t.items) == 0 {
		return "", ErrEmpty
	}
	last := t.items[len(t.items)-1]
	t.items = t.items[:len(t.items)-1]
	return last, nil
}

// Len is fully covered.
func (t *Truc) Len() int {
	return len(t.items)
}

// Join is never called by the tests.
func (t *Truc) Join(sep string) string {
	if len(t.items) == 0 {
		return ""
	}
	return strings.Join(t.items, sep)
}

// Reset is never called by the tests.
func (t *Truc) Reset() {
	t.items = nil
}
