package main

import "testing"

func TestRemoveProfane(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			input:    "hey kerfuffle sharbert fornax",
			expected: "hey **** **** ****",
		},
		{
			input:    "KERFUFFle",
			expected: "****",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "kerfuffle!",
			expected: "kerfuffle!",
		},
	}

	for _, c := range cases {
		got := RemoveProfane(c.input)
		if got != c.expected {
			t.Errorf("input does not match the expected. got: %s, expected: %s", got, c.expected)
		}
	}
}
