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
			input: " ",
			expected: []string{},
		},
		{
			input: "  hello  ",
			expected: []string{"hello"},
		},
		{
			input: "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Fatalf("cleanInput(%v): Output count mismatch. Expected: %v | Actual: %v", c.input, c.expected, actual)
		}

		for i, word := range actual {
			if word != c.expected[i] {
				t.Errorf("cleanInput(%v): Output string mismatch. Expected: %v | Actual: %v", c.input, c.expected[i], word)
			}
		}
	}
}
