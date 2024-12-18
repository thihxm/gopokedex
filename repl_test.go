package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "baTatiNha Quando Nasce",
			expected: []string{"batatinha", "quando", "nasce"},
		},
		{
			input:    "hello   World 2 ",
			expected: []string{"hello", "world", "2"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) == %q, expected %q", c.input, actual, c.expected)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("cleanInput(%q)[%d] == %q, expected %q", c.input, i, word, expectedWord)
			}
		}
	}
}
