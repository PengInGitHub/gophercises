package main

import "testing"

func TestNormalize(t *testing.T) {
	testCases := []struct {
		input  string
		answer string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567890"},
		{"(123) 456 7892", "1234567890"},
		{"(123) 456-7893", "1234567890"},
		{"123-456-7894", "1234567890"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567890"},
		{"(123)456-7892", "1234567890"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.answer {
				t.Errorf("Mistake: got %s but answer is %s", actual, tc.answer)
			}
		})
	}
}
