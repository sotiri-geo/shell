package main

import (
	"strings"
)

func parseArgs(args string) string {
	output := make([]string, 0, len(args))

	// Identify each word within a single quote
	// when we identify a single quote, we essentially need to traverse until we find
	// a closing quote, whilst preserving the characters found.

	i := 0
	// Behaves like a while loop
	for i < len(args) {
		// Handle single quotes in a single batch
		if args[i] == '\'' {
			if end, found := getClosingIndex(i, '\'', args); found {
				output = append(output, args[i+1:end])
				// fast-forward
				i = end + 1
			}
			// handle error
			continue
		}
		// We now need to handle non-quoted chacacters that are not whitespace
		if args[i] != ' ' {
			if end, found := getClosingIndex(i, ' ', args); found {
				output = append(output, args[i+1:end])
				i = end + 1
			}
		}
		i++
	}

	return strings.Join(output, " ")
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

