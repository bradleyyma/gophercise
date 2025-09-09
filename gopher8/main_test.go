package main

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"(123) 456-7890", "1234567890"},
		{"123-456-7890", "1234567890"},
		{"1234567890", "1234567890"},
		{"(123)456-7890", "1234567890"},
		{"123 456 7890", "1234567890"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalize(tt.input)
			if got != tt.want {
				t.Errorf("normalize(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}
