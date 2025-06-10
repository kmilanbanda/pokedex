package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: " eMiLiA BeSt gIrL",
			expected: []string{"emilia", "best", "girl"},
		},
		{
			input: "TTYL",
			expected: []string{"ttyl"},
		},
		{
			input: "hELP mE Help Me",
			expected: []string{"help", "me", "help", "me"},
		},
	}

	for _, c :=  range cases {
		actual := cleanInput(c.input)
		if len(c.expected) != len(actual) {
			t.Errorf("expected length (%v) != actual length (%v)", len(c.expected), len(actual))
		}
	}
}

