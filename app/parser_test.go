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
		"PreservesSpecialCharacters": {
			input: "'$%!~\\'",
			want:  "$%!~\\",
		},
		"WithEmptyQuotesIgnored": {
			input: "hello''world",
			want:  "helloworld",
		},
		"WithAdjacentQuotedStrings": {
			input: "'hello''world'",
			want:  "helloworld",
		},
		"WithGapAndAdjacentQuotedStrings": {
			input: "'hello'   'world'",
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

func TestClosingIndex(t *testing.T) {
	testCases := map[string]struct {
		start  int
		sep    byte
		args   string
		want   int
		exists bool
	}{
		"WithSingleQuote": {
			start:  0,
			sep:    '\'',
			args:   "'hello'",
			want:   6,
			exists: true,
		},
		"ReturnsFirstWord": {
			start:  0,
			sep:    '\'',
			args:   "'hello' 'world'",
			want:   6,
			exists: true,
		},
		"WithNonQuotes": {
			start:  0,
			sep:    ' ',
			args:   "hello world",
			want:   5,
			exists: true,
		},
		"WithStartAtNonZeroIndex": {
			start:  5,
			sep:    ' ',
			args:   "this is amazing",
			want:   7,
			exists: true,
		},
		"WithNonExistentSeparatorReturnsEndOfLine": {
			start:  0,
			sep:    '\'',
			args:   "'hello world",
			want:   12,
			exists: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, found := getClosingIndex(tc.start, tc.sep, tc.args)

			if found != tc.exists {
				t.Fatalf("should exist %v: found %v", tc.exists, found)
			}

			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}
