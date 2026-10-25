package fakerepos

import "testing"

func TestTrucPushPop(t *testing.T) {
	var tr Truc
	tr.Push("a")
	tr.Push("b")
	if tr.Len() != 2 {
		t.Fatalf("expected 2 items, got %d", tr.Len())
	}
	got, err := tr.Pop()
	if err != nil || got != "b" {
		t.Fatalf("expected b, got %q (%v)", got, err)
	}
}
