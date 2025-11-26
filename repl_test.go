package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO  			world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello  WorLD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "	hello	world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   		",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		actL := len(actual)
		expL := len(c.expected)
		if actL != expL {
			t.Errorf("Expected length of a slice is: %d, actual is: %d", expL, actL)
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Actual word: %s when expected is: %s", word, expectedWord)
			}
		}
	}
}
