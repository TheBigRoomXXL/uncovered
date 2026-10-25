package something

import "strings"

// Reverse is fully covered.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// IsPalindrome is partially covered: the case-insensitive path is never tested.
func IsPalindrome(s string, ignoreCase bool) bool {
	if ignoreCase {
		s = strings.ToLower(s)
	}
	return s == Reverse(s)
}

// CountVowels is never called by the tests.
func CountVowels(s string) int {
	count := 0
	for _, c := range strings.ToLower(s) {
		switch c {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}
