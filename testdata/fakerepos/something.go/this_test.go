package something

import "testing"

func TestReverse(t *testing.T) {
	if got := Reverse("héllo"); got != "olléh" {
		t.Errorf("Reverse = %q", got)
	}
}

func TestIsPalindrome(t *testing.T) {
	if !IsPalindrome("kayak", false) {
		t.Error("kayak should be a palindrome")
	}
	if IsPalindrome("gopher", false) {
		t.Error("gopher should not be a palindrome")
	}
}
