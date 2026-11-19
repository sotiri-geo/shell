package main

import (
	"errors"
	"strings"
)

var ErrUnclosedSingleQuote = errors.New("closing single quote not found")

// validateArgs verifies that the argument is parseable. Currently, this is defined that
// an open single quote is followed by a closing single quote.
func validateArgs(args string) error {
	count := 0

	for _, chr := range args {
		if chr == '\'' {
			count++
		}
	}

	if count%2 != 0 {
		return ErrUnclosedSingleQuote
	}
	return nil
}

// parseArgs will preseve the characters inside the single quotes as per the GNU docs. Args that
// are defined outside this boundary are separated by whitespace.
func parseArgs(args string) ([]string, error) {
	output := []string{}
	// Vaidate that for an open quote has a closing quote
	if err := validateArgs(args); err != nil {
		return nil, err
	}

	start, end := 0, 0
	// The reconstruction of an arg begins as an empty string buffer. We grow the buffer by the following rules:
	// 1. non-quoted characters
	// 2. non-whitespace characters
	//
	// The traversal continues until the following conditions are met, which we then flush to output slice:
	// 1. End of string
	// 2. Closing quote
	// h e l l o ' '
	var arg strings.Builder
	for end < len(args) {
		if args[end] == '\'' {
			// Start measuring from the moment we inspect the first open quote
			start = end
			end++
			for end < len(args) {
				if args[end] == '\'' {
					if end-start > 1 {
						arg.WriteString(args[start+1 : end])
					}
					// empty '' is effectively ignored
					break
				}
				end++
			}
			// move cursor to the next character after closing quote
			// e.g. 'abc'd
			//           ^ (start = end)
			end++
			start = end
			continue
		}

		if args[end] == ' ' {
			// if there is existing arg - flush to output
			if arg.Len() > 0 {
				output = append(output, arg.String())
				arg.Reset()
				continue
			}
			// Traverse until exhausted non-whitespace characters
			for end < len(args) {
				if args[end] != ' ' {
					break
				}
				end++
			}
			start = end
			continue
		}

		// Continue to grow buffer
		arg.WriteByte(args[end])
		end++
	}
	// Flush remaining arg
	if arg.Len() > 0 {
		output = append(output, arg.String())
	}

	return output, nil
}

func getClosingIndex(start int, sep byte, args string) (end int, found bool) {
	// We first inspect the next element in array
	end = start + 1
	for end < len(args) {
		if args[end] == sep {
			return end, true
		}
		end++
	}
	return end, false
}
