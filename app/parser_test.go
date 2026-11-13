package main

import "testing"

func TestParserUnquoted(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  string
	}{
		"MultipleWords": {
			input: "hello world",
			want:  "hello world",
		},
		"SingleWord": {
			input: "hello",
			want:  "hello",
		},
		"WithWhitespaceCollapsed": {
			input: "hello    world",
			want:  "hello world",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := parseArgs(tc.input)

			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParserSingleQuoted(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  string
	}{
		"RemovesQuotes": {
			input: "'hello world'",
			want:  "hello world",
		},
		"PreservesWhitespace": {
			input: "'hello    world'",
			want:  "hello    world",
		},
		"SingleQuotes": {
			input: "''",
			want:  "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := parseArgs(tc.input)

			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
