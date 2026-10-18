package fakerepos

import "testing"

func TestClassify(t *testing.T) {
	if got := Classify(3); got != "positive" {
		t.Errorf("Classify(3) = %q", got)
	}
	if got := Classify(0); got != "zero" {
		t.Errorf("Classify(0) = %q", got)
	}
}

func TestFizzBuzz(t *testing.T) {
	cases := map[int]string{1: "1", 3: "Fizz", 5: "Buzz", 15: "FizzBuzz"}
	for n, want := range cases {
		if got := FizzBuzz(n); got != want {
			t.Errorf("FizzBuzz(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestDivide(t *testing.T) {
	got, err := Divide(10, 2)
	if err != nil || got != 5 {
		t.Fatalf("Divide(10, 2) = %d, %v", got, err)
	}
}
