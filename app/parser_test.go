package main

import (
	"errors"
	"slices"
	"testing"
)

func TestParseArgsUnquoted(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "SplitsAtSpaces",
			input: "hello world",
			want:  []string{"hello", "world"},
		},
		{
			name:  "KeepsSingleWord",
			input: "hello",
			want:  []string{"hello"},
		},
		{
			name:  "CollapsesConsecutiveSpaces",
			input: "hello    world",
			want:  []string{"hello", "world"},
		},
		{
			name:  "ReturnsNoArgumentsForSpacesOnly",
			input: "   ",
			want:  []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertParsedArgs(t, tc.input, tc.want)
		})
	}
}

func TestParseArgsSingleQuoted(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "RemovesQuotes",
			input: "'hello world'",
			want:  []string{"hello world"},
		},
		{
			name:  "PreservesWhitespace",
			input: "'hello    world'",
			want:  []string{"hello    world"},
		},
		{
			name:  "PreservesWhitespaceAtQuotedEdges",
			input: "' \tx\t '",
			want:  []string{" \tx\t "},
		},
		{
			name:  "ReturnsNoArgumentsForEmptyQuotes",
			input: "''",
			want:  []string{},
		},
		{
			name:  "IgnoresStandaloneEmptyQuotesBetweenArgs",
			input: "left '' right",
			want:  []string{"left", "right"},
		},
		{
			name:  "PreservesSpecialCharacters",
			input: "'$%!~\\'",
			want:  []string{"$%!~\\"},
		},
		{
			name:  "JoinsFragmentsAroundEmptyQuotes",
			input: "hello''world",
			want:  []string{"helloworld"},
		},
		{
			name:  "ConcatenatesAdjacentQuotedFragments",
			input: "'hello''world'",
			want:  []string{"helloworld"},
		},
		{
			name:  "SplitsQuotedFragmentsSeparatedBySpaces",
			input: "'hello'   'world'",
			want:  []string{"hello", "world"},
		},
		{
			name:  "CombinesQuotedAndUnquotedFragments",
			input: "pre'x y'post next",
			want:  []string{"prex ypost", "next"},
		},
		{
			name:  "IgnoresTrailingEmptyQuotes",
			input: "helloworld''",
			want:  []string{"helloworld"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertParsedArgs(t, tc.input, tc.want)
		})
	}
}

func TestParseArgsDoubleQuoted(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			// GIVEN consecutive spaces inside double quotes
			name:  "PreservesSpaces",
			input: `"hello   world"`,
			want:  []string{"hello   world"},
		},
		{
			// GIVEN adjacent double quotes between two unquoted fragments
			name:  "JoinsFragmentsAroundEmptyQuotes",
			input: `hello""world`,
			want:  []string{"helloworld"},
		},
		{
			// GIVEN a double-quoted prefix next to unquoted text
			name:  "JoinsQuotedPrefixToUnquotedSuffix",
			input: `"hello"world`,
			want:  []string{"helloworld"},
		},
		{
			// GIVEN an unquoted prefix next to double-quoted text
			name:  "JoinsUnquotedPrefixToQuotedSuffix",
			input: `hello"world"`,
			want:  []string{"helloworld"},
		},
		{
			// GIVEN double-quoted fragments separated by an unquoted space
			name:  "SplitsQuotedFragmentsSeparatedBySpaces",
			input: `"hello" "world"`,
			want:  []string{"hello", "world"},
		},
		{
			// GIVEN a single quote inside a double-quoted argument
			name:  "KeepsLiteralSingleQuote",
			input: `"sotiri's test"`,
			want:  []string{"sotiri's test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertParsedArgs(t, tc.input, tc.want)
		})
	}
}

func TestParseArgsBackslashEscapes(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			// GIVEN an unquoted argument with three escaped spaces
			name:  "KeepsEscapedSpacesInOneArgument",
			input: `three\ \ \ spaces`,
			want:  []string{"three   spaces"},
		},
		{
			name:  "SplitsAtUnescapedSpaceBeforeEscapedSpaces",
			input: `three \ \ \ spaces`,
			want:  []string{"three", "   spaces"},
		},
		{
			name:  "SplitsAtUnescapedSpacesAfterEscapedSpace",
			input: `three\   spaces`,
			want:  []string{"three ", "spaces"},
		},
		{
			name:  "RemovesBackslashBeforeUnquotedLetter",
			input: `test\nexample`,
			want:  []string{"testnexample"},
		},
		{
			name:  "UnquotedBackslashEscapesBackslash",
			input: `hello\\world`,
			want:  []string{`hello\world`},
		},
		{
			name:  "KeepsEscapedSpecialCharactersLiteral",
			input: `\'\"\$\*\&\?`,
			want:  []string{`'"$*&?`},
		},
		{
			name:  "KeepsBackslashBeforeNewlineInsideSingleQuotes",
			input: "'shell\\\nscript'",
			want:  []string{"shell\\\nscript"},
		},
		{
			name:  "KeepsBackslashBeforeDoubleQuoteInsideSingleQuotes",
			input: `'example\"test'`,
			want:  []string{`example\"test`},
		},
		{
			name:  "JoinsSingleQuotedTextAndEscapedSingleQuote",
			input: `test'example'\'hello`,
			want:  []string{`testexample'hello`},
		},
		{
			name:  "DoubleQuotedBackslashEscapesBackslash",
			input: `"A \\ escapes itself"`,
			want:  []string{`A \ escapes itself`},
		},
		{
			name:  "DoubleQuotedBackslashEscapesDoubleQuote",
			input: `"A \" inside double quotes"`,
			want:  []string{`A " inside double quotes`},
		},
		{
			name:  "KeepsBackslashBeforeOrdinaryCharacterInsideDoubleQuotes",
			input: `"pre\qpost"`,
			want:  []string{`pre\qpost`},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertParsedArgs(t, tc.input, tc.want)
		})
	}
}

func assertParsedArgs(t *testing.T, input string, want []string) {
	t.Helper()

	// WHEN invoking the parser with the input
	got, err := parseArgs(input)

	// THEN it returns the expected arguments without an error
	if err != nil {
		t.Fatalf("parseArgs(%q): %v", input, err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("parseArgs(%q) = %q, want %q", input, got, want)
	}
}

func TestParseArgsErrors(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  error
	}{
		{
			// GIVEN unquoted input ending in a backslash
			name:  "RejectsTrailingUnquotedBackslash",
			input: `hello\`,
			want:  ErrEndingInBackslash,
		},
		{
			// GIVEN a single quote with no matching closing quote
			name:  "RejectsLoneOpeningSingleQuote",
			input: "'",
			want:  ErrUnclosedSingleQuote,
		},
		{
			// GIVEN an unclosed single quote after a pair of empty quotes
			name:  "RejectsUnclosedSingleQuoteAfterEmptyQuotes",
			input: "hello'''world",
			want:  ErrUnclosedSingleQuote,
		},
		{
			// GIVEN an open double quote and a backslash at the end of input
			name:  "RejectsTrailingBackslashInsideDoubleQuotes",
			input: `"hello\`,
			want:  ErrUnclosedDoubleQuote,
		},
		{
			// GIVEN a backslash that escapes the final double quote
			name:  "RejectsEscapedFinalDoubleQuote",
			input: `"hello\"`,
			want:  ErrUnclosedDoubleQuote,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// WHEN invoking the parser with the input
			got, err := parseArgs(tc.input)

			// THEN it returns the expected error and nil arguments
			if !errors.Is(err, tc.want) {
				t.Errorf("parseArgs(%q) error = %v, want %v", tc.input, err, tc.want)
			}
			if got != nil {
				t.Errorf("parseArgs(%q) args = %q, want nil", tc.input, got)
			}
		})
	}
}

func TestGetClosingIndex(t *testing.T) {
	testCases := []struct {
		name      string
		start     int
		sep       byte
		input     string
		want      int
		wantFound bool
	}{
		{
			name:      "FindsClosingSingleQuote",
			start:     0,
			sep:       '\'',
			input:     "'hello'",
			want:      6,
			wantFound: true,
		},
		{
			name:      "StopsAtFirstClosingQuote",
			start:     0,
			sep:       '\'',
			input:     "'hello' 'world'",
			want:      6,
			wantFound: true,
		},
		{
			name:      "FindsSpaceSeparator",
			start:     0,
			sep:       ' ',
			input:     "hello world",
			want:      5,
			wantFound: true,
		},
		{
			name:      "SearchesAfterStartIndex",
			start:     5,
			sep:       ' ',
			input:     "this is amazing",
			want:      7,
			wantFound: true,
		},
		{
			name:      "ReturnsEndOfInputWhenSeparatorIsMissing",
			start:     0,
			sep:       '\'',
			input:     "'hello world",
			want:      12,
			wantFound: false,
		},
		{
			name:      "StopsAtClosingQuoteBeforeSpace",
			start:     0,
			sep:       '\'',
			input:     "'hello' world",
			want:      6,
			wantFound: true,
		},
		{
			name:      "StopsAtClosingQuoteBeforeUnquotedText",
			start:     0,
			sep:       '\'',
			input:     "'hello'world",
			want:      6,
			wantFound: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, found := getClosingIndex(tc.start, tc.sep, tc.input)
			if got != tc.want || found != tc.wantFound {
				t.Errorf("getClosingIndex(%d, %q, %q) = (%d, %t), want (%d, %t)",
					tc.start, tc.sep, tc.input, got, found, tc.want, tc.wantFound)
			}
		})
	}
}
