package main

import "testing"

func TestIsUnique(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"abcd", true},
		{"abca", false},
		{"", true},
		{"a", true},
		{"aa", false},
		{"ab", true},
		{"abc", true},
		{"abbc", false},
		{"Aa", true},
	}

	for _, test := range tests {
		got := isUnique(test.input)
		if got != test.want {
			t.Errorf("isUnique(%q) = %v", test.input, got)
		}
	}
}
