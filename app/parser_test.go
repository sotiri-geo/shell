package main

import (
	"errors"
	"slices"
	"testing"
)

func TestParserUnquoted(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  []string
	}{
		"MultipleWords": {
			input: "hello world",
			want:  []string{"hello", "world"},
		},
		"SingleWord": {
			input: "hello",
			want:  []string{"hello"},
		},
		"WithWhitespaceCollapsed": {
			input: "hello    world",
			want:  []string{"hello", "world"},
		},
		"OnlySeparators": {
			input: "   ",
			want:  []string{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := parseArgs(tc.input)
			if err != nil {
				t.Fatalf("should not error: %v", err)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParserSingleQuoted(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  []string
	}{
		"RemovesQuotes": {
			input: "'hello world'",
			want:  []string{"hello world"},
		},
		"PreservesWhitespace": {
			input: "'hello    world'",
			want:  []string{"hello    world"},
		},
		"PreservesWhitespaceAtQuotedEdges": {
			input: "' \tx\t '",
			want:  []string{" \tx\t "},
		},
		"SingleEmptyQuotes": {
			input: "''",
			want:  []string{},
		},
		"IgnoresStandaloneEmptyQuotesBetweenArgs": {
			input: "left '' right",
			want:  []string{"left", "right"},
		},
		"PreservesSpecialCharacters": {
			input: "'$%!~\\'",
			want:  []string{"$%!~\\"},
		},
		"WithEmptyQuotesIgnored": {
			input: "hello''world",
			want:  []string{"helloworld"},
		},
		"WithAdjacentQuotedStrings": {
			input: "'hello''world'",
			want:  []string{"helloworld"},
		},
		"WithGapAndAdjacentQuotedStrings": {
			input: "'hello'   'world'",
			want:  []string{"hello", "world"},
		},
		"CombinesQuotedAndUnquotedFragments": {
			input: "pre'x y'post next",
			want:  []string{"prex ypost", "next"},
		},
		"WithEmptyQuotesIgnoredAtEnd": {
			input: "helloworld''",
			want:  []string{"helloworld"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := parseArgs(tc.input)
			if err != nil {
				t.Fatalf("should not error: %v", err)
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParserUnclosedSingleQuotes(t *testing.T) {
	testCases := map[string]string{
		"LoneOpeningQuote": "'",
		"AfterEmptyQuotes": "hello'''world",
	}
	for name, input := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := parseArgs(input)
			if !errors.Is(err, ErrUnclosedSingleQuote) {
				t.Errorf("got error: %v, want: %v", err, ErrUnclosedSingleQuote)
			}
			if got != nil {
				t.Errorf("got %v, want %v", got, nil)
			}
		})
	}
}

func TestParserDoubleQuoted(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  []string
	}{
		"WithMulipleWhiteSpacePreserved": {
			// GIVEN arg with consecutive whitespace
			input: "\"hello   world\"",
			want:  []string{"hello   world"},
		},
		"WithAdjacentQuotesConcatenated": {
			// GIVEN arg with double quotes next to eachother
			input: "hello\"\"world",
			want:  []string{"helloworld"},
		},
		"WithQuotedAndUnquotedAdjacent": {
			// GIVEN arg with double quoted text next to unquoted text
			input: "\"hello\"world",
			want:  []string{"helloworld"},
		},
		"WithUnquotedAndQuotedAdjacent": {
			// GIVEN arg with unquoted text next to doulbe quoted text
			input: "hello\"world\"",
			want:  []string{"helloworld"},
		},
		"WithQuotedSeparatedByWhitespace": {
			// GIVEN arg with double quoted text separated by whitespace
			input: "\"hello\" \"world\"",
			want:  []string{"hello", "world"},
		},
		"ContainingSingleQuotes": {
			// GIVEN arg with single quote contained inside double quote
			input: "\"sotiri's test\"",
			want:  []string{"sotiri's test"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// WHEN invoking parseArgs with a1
			got, err := parseArgs(tc.input)
			// THEN no error with correct output
			if err != nil {
				t.Fatalf("should not error: %v", err)
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParserBackslashEscape(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  []string
	}{
		// GIVEN arg with 3 escaped whitespaces are presere
		"WithMultilpeWhitespacePreserved": {
			input: `three\ \ \ spaces`,
			want:  []string{"three   spaces"},
		},
		"WithSingleWhitespaceNotEscapedAndRemainingEscaped": {
			input: `three \ \ \ spaces`,
			want:  []string{"three", "   spaces"},
		},
		"WithSingleWhitespacePreservedAndRemainingCollapsed": {
			input: `three\   spaces`,
			want:  []string{"three ", "spaces"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// WHEN invoking parseArgs with a1
			got, err := parseArgs(tc.input)
			// THEN no error with correct output
			if err != nil {
				t.Fatalf("should not error: %v", err)
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
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
		"WithSingleQuoteWhitespaceAndMultipleWords": {
			start:  0,
			sep:    '\'',
			args:   "'hello' world",
			want:   6,
			exists: true,
		},
		"WithSingleQuoteAndMultipleWords": {
			start:  0,
			sep:    '\'',
			args:   "'hello'world",
			want:   6,
			exists: true,
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
