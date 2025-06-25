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
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// add more cases here
		{
			input:    "My TeST Strings hereE",
			expected: []string{"my", "test", "strings", "heree"},
		},
		{
			input:    "ANOTHER TEST STRING    IN    HERE",
			expected: []string{"another", "test", "string", "in", "here"},
		},
		{
			input:    "             leading  whitespace and trailing    whitespace test          ",
			expected: []string{"leading", "whitespace", "and", "trailing", "whitespace", "test"},
		},
		/*
				below is an expected failure due to expected slice being wrong size
			{
				input:    "             leading  whitespace and trailing    whitespace test          ",
				expected: []string{"leading", "whitespace", "and", "trailing", "whitespace"},
			},
		*/
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("FAIL: slices not same length | %v != %v |\n", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("FAIL: %s != %s\n", word, expectedWord)
			}
		}
	}
}
