package main

import (
	"errors"
	"strings"
)

var ErrUnclosedSingleQuote = errors.New("closing single quote not found")

const (
	SingleQuote = '\''
	WhiteSpace  = ' '
)

// parseArgs preserves text inside single quotes and separates arguments at unquoted spaces.
func parseArgs(args string) ([]string, error) {
	output := []string{}
	var arg strings.Builder
	end := 0
	for end < len(args) {
		if args[end] == SingleQuote {
			closing, found := getClosingIndex(end, '\'', args)
			if !found {
				return nil, ErrUnclosedSingleQuote
			}

			arg.WriteString(args[end+1 : closing])
			end = closing + 1
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
				if args[end] != WhiteSpace {
					break
				}
				end++
			}
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
