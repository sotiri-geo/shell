package main

import (
	"strings"
)

func parseArgs(args string) string {
	var output strings.Builder
	// Identify each word within a single quote
	// when we identify a single quote, we essentially need to traverse until we find
	// a closing quote, whilst preserving the characters found.

	i := 0
	// Behaves like a while loop
	for i < len(args) {
		// Handle single quotes in a single batch
		if args[i] == '\'' {
			if end, found := getClosingIndex(i, '\'', args); found {
				output.WriteString(args[i+1 : end])
				// move to the char next to the closing quote
				i = end + 1
			}
			continue
		}
		// We now need to handle non-quoted chacacters that are not whitespace
		if args[i] != ' ' {
			// Walk up to a whitespace or a single quote
			start := i
			for i < len(args) {
				// Break out of current buffer and flush
				if args[i] == ' ' || args[i] == '\'' {
					output.WriteString(args[start:i])
					break
				}
				i++
			}

			// If exhausted all characters: flush to buffer
			if i == len(args) {
				output.WriteString(args[start:i])
			}

			continue
		}
		// Collapse whitespace from non-quoted characters
		if args[i] == ' ' {
			for i < len(args) {
				if args[i] == ' ' {
					i++
					continue
				}
				break
			}
			output.WriteString(" ")
			continue
		}
		i++
	}

	return output.String()
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
